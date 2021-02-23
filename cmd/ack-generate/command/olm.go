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

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"

	generate "github.com/aws-controllers-k8s/code-generator/pkg/generate"
	ackgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	olmgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/olm"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
)

const (
	olmConfigFileSuffix = "olmconfig.yaml"
)

// optOLMConfigPath
var optOLMConfigPath string
var optDisableCommonLinks bool
var optDisableCommonKeywords bool

// olmCmd is the command that generates a service ClusterServiceVersion base
// for generating an operator lifecycle manager bundle.
var olmCmd = &cobra.Command{
	Use: "olm <service> <version>",
	Short: `Generate an Operator Lifecycle Manager's ClusterServiceVersion
resources for AWS service API. Expects a configuration file to be passed in
via the --olm-config option, or otherwise reads the config at the root of the
project with the filename <servicealias>-olmconfig.yaml`,
	RunE: generateOLMAssets,
}

func init() {
	olmCmd.PersistentFlags().StringVar(
		&optOLMConfigPath, "olm-config", "", "the OLM configuration file to inform how OLM assets are generated.",
	)
	// optionally disable common links in the resulting manifest
	olmCmd.PersistentFlags().BoolVar(
		&optDisableCommonLinks, "no-common-links", false, "does not include common links in the rendered cluster service version.",
	)
	// optionally disable common keywords in the resulting manifest
	olmCmd.PersistentFlags().BoolVar(
		&optDisableCommonKeywords, "no-common-keywords", false, "does not include common keywords in the rendered cluster service version",
	)

	rootCmd.AddCommand(olmCmd)
}

// generateOLMAssets generates all assets necessary for delivering the
// service controllers via operator lifecycle manager ("OLM").
func generateOLMAssets(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf(
			"please specify the service alias and version " +
				"for the AWS service API to generate",
		)
	}
	svcAlias := strings.ToLower(args[0])
	if optOutputPath == "" {
		optOutputPath = filepath.Join(optServicesDir, svcAlias)
	}

	version := args[1]

	// get the generator inputs
	if err := ensureSDKRepo(optCacheDir); err != nil {
		return err
	}
	sdkHelper := ackmodel.NewSDKHelper(sdkDir)
	sdkAPI, err := sdkHelper.API(svcAlias)
	if err != nil {
		newSvcAlias, err := FallBackFindServiceID(sdkDir, svcAlias)
		if err != nil {
			return err
		}
		sdkAPI, err = sdkHelper.API(newSvcAlias) // retry with serviceID
		if err != nil {
			return fmt.Errorf("service %s not found", svcAlias)
		}
	}

	latestAPIVersion, err = getLatestAPIVersion()
	if err != nil {
		return err
	}
	g, err := generate.New(
		sdkAPI, latestAPIVersion, optGeneratorConfigPath, ackgenerate.DefaultConfig,
	)
	if err != nil {
		return err
	}

	if optOLMConfigPath == "" {
		optOLMConfigPath = strings.Join([]string{svcAlias, olmConfigFileSuffix}, "-")
	}

	// read the configuration from file
	svcConfigYAML, err := ioutil.ReadFile(optOLMConfigPath)
	if err != nil {
		fmt.Println("unable to read configuration file at path:", optOLMConfigPath)
		return err
	}

	// set the base metadata and then override values as
	// defined by the service config.
	svcMeta := ackmodel.NewOLMMetadata() // TODO this can take in variables
	if err = yaml.Unmarshal(svcConfigYAML, &svcMeta); err != nil {
		fmt.Println("unable to convert olm configuration file to data instance")
		return err
	}

	// prepare the common metadata
	commonMeta := ackmodel.OLMCommonMetadata{}
	if !optDisableCommonLinks {
		commonMeta.Links = ackmodel.CommonLinks
	}

	if !optDisableCommonKeywords {
		commonMeta.Keywords = ackmodel.CommonKeywords
	}

	// generate templates
	ts, err := olmgenerate.BundleAssets(g, commonMeta, svcMeta, version, optTemplatesDir)
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
		if _, err := ensureDir(outDir); err != nil {
			return err
		}
		if err = ioutil.WriteFile(outPath, contents.Bytes(), 0666); err != nil {
			return err
		}
	}

	return nil
}
