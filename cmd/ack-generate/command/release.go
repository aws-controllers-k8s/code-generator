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
	"strings"

	"github.com/spf13/cobra"

	ackgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	"github.com/aws-controllers-k8s/code-generator/pkg/sdk"
)

var releaseCmd = &cobra.Command{
	Use:   "release <service> <release_version>",
	Short: "Generates release artifacts for a specific service controller and release version",
	Long: `Runs the full release pipeline for a controller repo. By default uses
the current working directory; use -o to specify the controller path.
This includes template generation, CRDs, RBAC, and Helm template
post-processing.

Usage:
  ack-generate release s3 v1.2.3 -o /path/to/s3-controller
  # or from within the controller repo:
  ack-generate release s3 v1.2.3`,
	RunE: generateRelease,
}

func init() {
	rootCmd.AddCommand(releaseCmd)
}

func generateRelease(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("please specify the service alias and the release version to generate release artifacts for")
	}
	svcAlias := strings.ToLower(args[0])
	releaseVersion := strings.ToLower(args[1])

	ctx, cancel := sdk.ContextWithSigterm(context.Background())
	defer cancel()

	controllerCtx, err := resolveControllerContext(svcAlias)
	if err != nil {
		return err
	}

	if optImageRepository == "" {
		optImageRepository = fmt.Sprintf("public.ecr.aws/aws-controllers-k8s/%s-controller", svcAlias)
	}

	fmt.Printf("Building release artifacts for %s-%s\n", svcAlias, releaseVersion)
	ts, err := ackgenerate.Release(
		controllerCtx.Model, controllerCtx.Metadata, optTemplateDirs,
		releaseVersion, optImageRepository, optServiceAccountName,
		embeddedTemplatesFS,
	)
	if err != nil {
		return err
	}
	if err := writeTemplateSet(ts, optOutputPath); err != nil {
		return err
	}

	releaseOpts := ackgenerate.BuildReleaseOptions{
		SvcAlias:             controllerCtx.SvcAlias,
		ControllerSourcePath: controllerCtx.ControllerPath,
		APIVersion:           controllerCtx.APIVersion,
		RBACRoleName:         controllerCtx.RBACRoleName,
		RuntimeVersion:       controllerCtx.RuntimeVersion,
		ControllerGenVersion: optControllerGenVersion,
		CacheDir:             optCacheDir,
	}
	return ackgenerate.BuildRelease(ctx, releaseOpts)
}
