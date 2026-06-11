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
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	ackmetadata "github.com/aws-controllers-k8s/code-generator/pkg/metadata"
	"github.com/aws-controllers-k8s/code-generator/pkg/sdk"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

var (
	optTestConfigPath string
)

var e2eTestsCmd = &cobra.Command{
	Use:   "e2e-tests <service>",
	Short: "Generates Go e2e test files for an AWS service controller",
	Long: `Generates Go e2e test scaffolds that exercise the CRUD lifecycle
for each resource defined in the controller's testconfig.yaml file.

The generated tests use the shared test library from test-infra/pkg/e2e/
and follow the create → wait-synced → update → delete pattern.`,
	RunE: generateE2ETests,
}

func init() {
	e2eTestsCmd.PersistentFlags().StringVar(
		&optTestConfigPath, "test-config", "",
		"Path to testconfig.yaml file. Defaults to <output-path>/testconfig.yaml",
	)
	rootCmd.AddCommand(e2eTestsCmd)
}

func generateE2ETests(cmd *cobra.Command, args []string) error {
	cmdStart := time.Now()
	if len(args) != 1 {
		return fmt.Errorf("please specify the service alias for the AWS service API to generate")
	}
	svcAlias := strings.ToLower(args[0])
	if optOutputPath == "" {
		optOutputPath = filepath.Join(optServicesDir, svcAlias)
	}

	cfg, err := setupGenerator(svcAlias)
	if err != nil {
		return err
	}

	// Load testconfig.yaml
	testConfigPath := optTestConfigPath
	if testConfigPath == "" {
		testConfigPath = filepath.Join(optOutputPath, "testconfig.yaml")
	}
	testCfg, err := ackgenconfig.NewTestConfig(testConfigPath)
	if err != nil {
		return fmt.Errorf("loading test config: %w", err)
	}
	if testCfg == nil {
		return fmt.Errorf("testconfig.yaml not found at %s — create it to define test values for resources", testConfigPath)
	}

	// Load the AWS SDK model
	metadata, err := ackmetadata.NewServiceMetadata(optMetadataConfigPath)
	if err != nil {
		return err
	}
	m, err := loadModelWithLatestAPIVersion(svcAlias, metadata, cfg)
	if err != nil {
		return err
	}

	// Generate test templates
	ts, err := ackgenerate.E2ETests(m, optTemplateDirs, testCfg)
	if err != nil {
		return err
	}

	if err = ts.Execute(); err != nil {
		return err
	}

	for path, contents := range ts.Executed() {
		if optDryRun {
			fmt.Printf("============================= %s ======================================\n", path)
			fmt.Println(strings.TrimSpace(contents.String()))
			continue
		}
		outPath := filepath.Join(optOutputPath, path)
		outDir := filepath.Dir(outPath)
		if _, err := sdk.EnsureDir(outDir); err != nil {
			return err
		}
		if err = ioutil.WriteFile(outPath, contents.Bytes(), 0666); err != nil {
			return err
		}
	}

	util.Tracef("generateE2ETests total: %s\n", time.Since(cmdStart))
	return nil
}
