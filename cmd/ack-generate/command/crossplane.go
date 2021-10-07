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
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	ackgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/generate/config"
	cpgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/crossplane"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
)

// crossplaneCmd is the command that generates Crossplane API types
var crossplaneCmd = &cobra.Command{
	Use:   "crossplane <service>",
	Short: "Generate Crossplane Provider",
	RunE:  generateCrossplane,
}

var providerDir string

func init() {
	crossplaneCmd.PersistentFlags().StringVar(
		&providerDir, "provider-dir", ".", "the directory of the Crossplane provider",
	)
	rootCmd.AddCommand(crossplaneCmd)
}

func generateCrossplane(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("please specify the service alias for the AWS service API to generate")
	}
	ctx, cancel := contextWithSigterm(context.Background())
	defer cancel()
	if err := ensureSDKRepo(ctx, optCacheDir, optRefreshCache); err != nil {
		return err
	}
	svcAlias := strings.ToLower(args[0])
	cfgPath := filepath.Join(providerDir, "apis", svcAlias, optGenVersion, "generator-config.yaml")
	_, err := os.Stat(cfgPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if os.IsNotExist(err) {
		cfgPath = ""
	}
	cfg, err := ackgenconfig.New(cfgPath, ackgenerate.DefaultConfig)
	if err != nil {
		return err
	}
	sdkHelper := model.NewSDKHelper(sdkDir)
	sdkHelper.APIGroupSuffix = "aws.crossplane.io"
	sdkAPI, err := sdkHelper.API(svcAlias)
	if err != nil {
		retryModelName, err := FallBackFindServiceID(sdkDir, svcAlias)
		if err != nil {
			return err
		}
		// Retry using path found by querying service ID
		sdkAPI, err = sdkHelper.API(retryModelName)
		if err != nil {
			return fmt.Errorf("cannot get the API model for service %s", svcAlias)
		}
	}
	m, err := ackmodel.New(
		sdkAPI, svcAlias, optGenVersion, cfg,
	)
	if err != nil {
		return err
	}

	ts, err := cpgenerate.Crossplane(m, optTemplateDirs)
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
		outPath := filepath.Join(providerDir, path)
		outDir := filepath.Dir(outPath)
		if _, err := ensureDir(outDir); err != nil {
			return err
		}
		if err = ioutil.WriteFile(outPath, contents.Bytes(), 0666); err != nil {
			return err
		}
	}
	apiPath := filepath.Join(providerDir, "apis", svcAlias, optGenVersion)
	controllerPath := filepath.Join(providerDir, "pkg", "controller", svcAlias)
	// TODO(muvaf): goimports don't allow to be included as a library. Make sure
	// goimports binary exists.
	if err := exec.Command("goimports", "-w", apiPath, controllerPath).Run(); err != nil {
		return errors.Wrap(err, "cannot run goimports")
	}
	return nil
}
