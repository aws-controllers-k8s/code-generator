// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package command

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	ackchainsaw "github.com/aws-controllers-k8s/code-generator/pkg/generate/chainsaw"
	ackmetadata "github.com/aws-controllers-k8s/code-generator/pkg/metadata"
)

var (
	optChainsawTestsDir     string
	optChainsawTemplatesDir string
	optChainsawNamePrefix   string
)

var chainsawCmd = &cobra.Command{
	Use:   "chainsaw <service>",
	Short: "Generate Chainsaw e2e tests for the resources of an AWS service controller",
	Long: `Generates a Chainsaw e2e test for each resource that has an inputs.yaml in
<tests-dir>/<resource>/, e.g. test/e2e/chainsaw/vpc/inputs.yaml. The
inputs give the values the resource's model can't: a spec to create, updates,
and the resources it depends on.`,
	RunE: generateChainsaw,
}

func init() {
	chainsawCmd.PersistentFlags().StringVar(
		&optChainsawTestsDir, "tests-dir", "", "directory holding one directory per resource with an inputs.yaml (default <output>/test/e2e/chainsaw)",
	)
	chainsawCmd.PersistentFlags().StringVar(
		&optChainsawTemplatesDir, "templates-dir", "", "directory to write the shared Chainsaw step templates to (default <tests-dir>/_templates)",
	)
	chainsawCmd.PersistentFlags().StringVar(
		&optChainsawNamePrefix, "name-prefix", "", "prefix for the names of the generated Chainsaw tests",
	)
	chainsawCmd.PersistentFlags().StringVar(
		&optGenVersion, "version", "v1alpha1", "the resource API Version of the CRDs under test",
	)
	rootCmd.AddCommand(chainsawCmd)
}

func generateChainsaw(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("please specify the service alias for the AWS service API to generate")
	}
	svcAlias := strings.ToLower(args[0])
	if optOutputPath == "" {
		optOutputPath = filepath.Join(optServicesDir, svcAlias)
	}
	if optChainsawTestsDir == "" {
		optChainsawTestsDir = filepath.Join(optOutputPath, "test", "e2e", "chainsaw")
	}
	if optChainsawTemplatesDir == "" {
		optChainsawTemplatesDir = filepath.Join(optChainsawTestsDir, ackchainsaw.TemplatesDirName)
	}
	// The tests refer to the templates by a path relative to their directory,
	// which needs both to be absolute, or both relative to the same place
	for _, dir := range []*string{&optChainsawTestsDir, &optChainsawTemplatesDir} {
		abs, err := filepath.Abs(*dir)
		if err != nil {
			return err
		}
		*dir = abs
	}

	cfg, err := setupGenerator(svcAlias)
	if err != nil {
		return err
	}
	metadata, err := ackmetadata.NewServiceMetadata(optMetadataConfigPath)
	if err != nil {
		return err
	}
	m, err := loadModelWithLatestAPIVersion(svcAlias, metadata, cfg)
	if err != nil {
		return err
	}
	crds, err := m.GetCRDs()
	if err != nil {
		return err
	}

	if _, err := os.Stat(optChainsawTestsDir); err == nil {
		if err := writeSuiteFiles(svcAlias); err != nil {
			return err
		}
	}

	// A resource whose test can't be generated doesn't stop the others; the
	// failures are reported, and fail the command, at the end
	var withoutInputs, failed []string
	for _, crd := range crds {
		testDir := filepath.Join(optChainsawTestsDir, strings.ReplaceAll(crd.Names.Snake, "_", "-"))
		inputsPath := filepath.Join(testDir, ackchainsaw.InputsFileName)
		if _, err := os.Stat(inputsPath); os.IsNotExist(err) {
			withoutInputs = append(withoutInputs, crd.Kind)
			continue
		}
		in, err := ackchainsaw.LoadInputs(inputsPath)
		if err != nil {
			failed = append(failed, fmt.Sprintf("%s: %v", crd.Kind, err))
			continue
		}
		templatesDir, err := filepath.Rel(testDir, optChainsawTemplatesDir)
		if err != nil {
			return err
		}
		files, report, err := ackchainsaw.Generate(m, crd, in, ackchainsaw.Options{
			Service:      svcAlias,
			APIVersion:   optGenVersion,
			TemplatesDir: filepath.ToSlash(templatesDir),
			NamePrefix:   optChainsawNamePrefix,
		})
		if err != nil {
			failed = append(failed, fmt.Sprintf("%s: %v", crd.Kind, err))
			continue
		}

		names := make([]string, 0, len(files))
		for name := range files {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			if optDryRun {
				fmt.Printf("============================= %s ======================================\n", filepath.Join(testDir, name))
				fmt.Println(strings.TrimSpace(string(files[name])))
				continue
			}
			if err := os.WriteFile(filepath.Join(testDir, name), files[name], 0o644); err != nil {
				return err
			}
		}
		if !optDryRun {
			if err := removeStaleFiles(testDir, files); err != nil {
				return err
			}
		}
		fmt.Printf("%s (%s):\n%s\n", crd.Kind, testDir, report)
	}
	if len(withoutInputs) != 0 {
		fmt.Printf("No inputs, so no test: %s\n", strings.Join(withoutInputs, ", "))
	}
	if len(failed) != 0 {
		return fmt.Errorf("couldn't generate tests for:\n  %s", strings.Join(failed, "\n  "))
	}
	return nil
}

// writeSuiteFiles writes the step templates the service's generated tests
// share, to the templates directory, keeping a hand-written file in place of
// one, and removes generated templates the generator no longer writes. It also
// writes a starting .chainsaw.yaml if the tests directory has none; an
// existing one is the controller's, and is left as it is.
func writeSuiteFiles(service string) error {
	files, err := ackchainsaw.SuiteFiles()
	if err != nil {
		return err
	}
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	templates := map[string][]byte{}
	for _, name := range names {
		dest := filepath.Join(optChainsawTemplatesDir, path.Base(name))
		templates[path.Base(name)] = files[name]
		if optDryRun {
			fmt.Printf("============================= %s ======================================\n", dest)
			fmt.Println(strings.TrimSpace(string(files[name])))
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		written, err := ackchainsaw.WriteGeneratedFile(dest, files[name])
		if err != nil {
			return err
		}
		if !written {
			fmt.Printf("keeping hand-written %s\n", dest)
		}
	}
	if optDryRun {
		return nil
	}
	if err := removeStaleFiles(optChainsawTemplatesDir, templates); err != nil {
		return err
	}

	config := filepath.Join(optChainsawTestsDir, ackchainsaw.ConfigFileName)
	if _, err := os.Stat(config); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	scaffold, err := ackchainsaw.ScaffoldConfig(service)
	if err != nil {
		return err
	}
	if err := os.WriteFile(config, scaffold, 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote a starting %s; it's this controller's to edit, and won't be regenerated\n", config)
	return nil
}

// removeStaleFiles removes the generated files in testDir that the generator
// no longer writes, e.g. the manifest of an update that's been taken out of
// the inputs. Only files starting with the generated header are removed, so
// inputs.yaml and hand-written files stay.
func removeStaleFiles(testDir string, files map[string][]byte) error {
	entries, err := os.ReadDir(testDir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		if _, ok := files[e.Name()]; ok {
			continue
		}
		path := filepath.Join(testDir, e.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.HasPrefix(string(content), ackchainsaw.GeneratedHeader) {
			if err := os.Remove(path); err != nil {
				return err
			}
			fmt.Printf("removed stale %s\n", path)
		}
	}
	return nil
}
