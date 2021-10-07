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
)

var (
	optReleaseOutputPath  string
	optImageRepository    string
	optServiceAccountName string
)

var releaseCmd = &cobra.Command{
	Use:   "release <service> <release_version>",
	Short: "Generates release artifacts for a specific service controller and release version",
	RunE:  generateRelease,
}

func init() {
	releaseCmd.PersistentFlags().StringVar(
		&optImageRepository, "image-repository", "", "the Docker image repository to use in release artifacts. Defaults to 'public.ecr.aws/aws-controllers-k8s/$service-controller'",
	)
	releaseCmd.PersistentFlags().StringVar(
		&optServiceAccountName, "service-account-name", "default", "The name of the ServiceAccount AND ClusterRole used for ACK service controller",
	)
	releaseCmd.PersistentFlags().StringVarP(
		&optReleaseOutputPath, "output", "o", "", "path to root directory to create generated files. Defaults to "+optServicesDir+"/$service",
	)
	rootCmd.AddCommand(releaseCmd)
}

// generateRelease generates the Helm charts and other release artifacts for a
// service controller and release version
func generateRelease(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("please specify the service alias and the release version to generate release artifacts for")
	}
	svcAlias := strings.ToLower(args[0])
	if optReleaseOutputPath == "" {
		optReleaseOutputPath = filepath.Join(optServicesDir, svcAlias)
	}
	if optImageRepository == "" {
		optImageRepository = fmt.Sprintf("public.ecr.aws/aws-controllers-k8s/%s-controller", svcAlias)
	}
	// TODO(jaypipes): We could do some git-fu here to verify that the release
	// version supplied hasn't been used (as a Git tag) before...
	releaseVersion := strings.ToLower(args[1])

	ctx, cancel := contextWithSigterm(context.Background())
	defer cancel()
	if err := ensureSDKRepo(ctx, optCacheDir, optRefreshCache); err != nil {
		return err
	}
	m, err := loadModel(svcAlias, "")
	if err != nil {
		return err
	}

	metadata, err := ackmetadata.NewServiceMetadata(optMetadataConfigPath)
	if err != nil {
		return err
	}

	ts, err := ackgenerate.Release(
		m, metadata, optTemplateDirs,
		releaseVersion, optImageRepository, optServiceAccountName,
	)
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
		outPath := filepath.Join(optReleaseOutputPath, path)
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
