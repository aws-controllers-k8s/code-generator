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
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	ackgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	ackmetadata "github.com/aws-controllers-k8s/code-generator/pkg/metadata"
	"github.com/aws-controllers-k8s/code-generator/pkg/sdk"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

var controllerCmd = &cobra.Command{
	Use:   "controller <service>",
	Short: "Generates a fully-built service controller from the current directory",
	Long: `Runs the full generation pipeline for a controller repo. By default uses
the current working directory; use -o to specify the controller path.
This includes APIs, deepcopy, CRDs, controller code, RBAC, formatting,
and boilerplate file copying.

Usage:
  ack-generate controller s3 -o /path/to/s3-controller
  # or from within the controller repo:
  ack-generate controller s3`,
	RunE: generateController,
}

func init() {
	rootCmd.AddCommand(controllerCmd)
}

func generateController(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("please specify the service alias for the AWS service API to generate")
	}
	svcAlias := strings.ToLower(args[0])

	ctx, cancel := sdk.ContextWithSigterm(context.Background())
	defer cancel()

	controllerCtx, err := resolveControllerContext(svcAlias)
	if err != nil {
		return err
	}

	pipelineOpts := ackgenerate.BuildControllerOptions{
		SvcAlias:             controllerCtx.SvcAlias,
		ControllerSourcePath: controllerCtx.ControllerPath,
		APIVersion:           controllerCtx.APIVersion,
		RBACRoleName:         controllerCtx.RBACRoleName,
		RuntimeVersion:       controllerCtx.RuntimeVersion,
		ControllerGenVersion: optControllerGenVersion,
		CacheDir:             optCacheDir,
		BoilerplateFS:        embeddedBoilerplateFS,
		TemplatesFS:          embeddedTemplatesFS,
	}

	// Step 1: Generate and write all templates (APIs + controller code)
	fmt.Printf("Building service controller for %s\n", svcAlias)
	ts, err := ackgenerate.Controller(
		controllerCtx.Model, optTemplateDirs, optServiceAccountName,
		controllerCtx.APIVersion, embeddedTemplatesFS,
	)
	if err != nil {
		return err
	}
	if err := writeTemplateSet(ts, optOutputPath); err != nil {
		return err
	}

	// Save generation metadata and copy generator.yaml
	apisPath := filepath.Join(optOutputPath, "apis")
	if err := ackmetadata.CreateGenerationMetadata(
		controllerCtx.APIVersion, apisPath, ackmetadata.UpdateReasonAPIGeneration,
		sdkVersion, optGeneratorConfigPath,
	); err != nil {
		return fmt.Errorf("creating generation metadata: %w", err)
	}
	if err := util.CopyFile(optGeneratorConfigPath, filepath.Join(apisPath, controllerCtx.APIVersion, "generator.yaml")); err != nil {
		return fmt.Errorf("copying generator.yaml: %w", err)
	}

	// Step 2: Pre-codegen pipeline (runtime CRDs, deepcopy, CRDs)
	if err := ackgenerate.BuildControllerPreCodegen(ctx, pipelineOpts); err != nil {
		return err
	}

	// Step 3: Post-codegen pipeline (go mod tidy, RBAC, formatting, boilerplate)
	return ackgenerate.BuildControllerPostCodegen(pipelineOpts)
}
