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
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	ackgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	ackmetadata "github.com/aws-controllers-k8s/code-generator/pkg/metadata"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

type contentType int

const (
	ctUnknown contentType = iota
	ctJSON
	ctYAML
)

var (
	optGenVersion   string
	apisVersionPath string
)

// apiCmd is the command that generates service API types
var apisCmd = &cobra.Command{
	Use:      "apis <service>",
	Short:    "Generate Kubernetes API type definitions for an AWS service API",
	RunE:     generateAPIs,
	PostRunE: saveGeneratedMetadata,
}

func init() {
	apisCmd.PersistentFlags().StringVar(
		&optGenVersion, "version", "v1alpha1", "the resource API Version to use when generating API infrastructure and type definitions",
	)
	rootCmd.AddCommand(apisCmd)
}

// saveGeneratedMetadata saves the parameters used to generate APIs and checksum
// of the generated code.
func saveGeneratedMetadata(cmd *cobra.Command, args []string) error {
	err := ackmetadata.CreateGenerationMetadata(
		optGenVersion,
		filepath.Join(optOutputPath, "apis"),
		ackmetadata.UpdateReasonAPIGeneration,
		optAWSSDKGoVersion,
		optGeneratorConfigPath,
	)
	if err != nil {
		return fmt.Errorf("cannot create generation metadata file: %v", err)
	}

	copyDest := filepath.Join(
		optOutputPath, "apis", optGenVersion, "generator.yaml",
	)
	err = util.CopyFile(optGeneratorConfigPath, copyDest)
	if err != nil {
		return fmt.Errorf("cannot copy generator configuration file: %v", err)
	}

	return nil
}

// generateAPIs generates the Go files for each resource in the AWS service
// API.
func generateAPIs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("please specify the service alias for the AWS service API to generate")
	}
	svcPackage := strings.ToLower(args[0])
	if optOutputPath == "" {
		optOutputPath = filepath.Join(optServicesDir, svcPackage)
	}
	ctx, cancel := contextWithSigterm(context.Background())
	defer cancel()
	if err := ensureSDKRepo(ctx, optCacheDir, optRefreshCache); err != nil {
		return err
	}
	if optModelName == "" {
		optModelName = svcPackage
	}
	sdkHelper := ackmodel.NewSDKHelper(sdkDir)
	sdkAPI, err := sdkHelper.API(optModelName)
	if err != nil {
		newSvcAlias, err := FallBackFindServiceID(sdkDir, optModelName)
		if err != nil {
			return err
		}
		sdkAPI, err = sdkHelper.API(newSvcAlias) // retry with serviceID
		if err != nil {
			return fmt.Errorf("service %s not found", svcPackage)
		}
	}
	model, err := ackmodel.New(
		sdkAPI, svcPackage, optGenVersion, optGeneratorConfigPath, ackgenerate.DefaultConfig,
	)
	if err != nil {
		return err
	}
	ts, err := ackgenerate.APIs(model, optTemplateDirs)
	if err != nil {
		return err
	}

	if err = ts.Execute(); err != nil {
		return err
	}

	apisVersionPath = filepath.Join(optOutputPath, "apis", optGenVersion)
	for path, contents := range ts.Executed() {
		if optDryRun {
			fmt.Printf("============================= %s ======================================\n", path)
			fmt.Println(strings.TrimSpace(contents.String()))
			continue
		}
		outPath := filepath.Join(apisVersionPath, path)
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
