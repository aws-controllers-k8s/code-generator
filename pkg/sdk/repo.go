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

package sdk

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/mod/modfile"
)

const (
	sdkGoV2Module = "github.com/aws/aws-sdk-go-v2"
)

func ContextWithSigterm(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	signalCh := make(chan os.Signal, 1)

	// recreate the context.CancelFunc
	cancelFunc := func() {
		signal.Stop(signalCh)
		cancel()
	}

	// notify on SIGINT or SIGTERM
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-signalCh:
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx, cancelFunc
}

// EnsureDir makes sure that a supplied directory exists and
// returns whether the directory already existed.
func EnsureDir(fp string) (bool, error) {
	fi, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return false, os.MkdirAll(fp, os.ModePerm)
		}
		return false, err
	}
	if !fi.IsDir() {
		return false, fmt.Errorf("expected %s to be a directory", fp)
	}
	if !isDirWriteable(fp) {
		return true, fmt.Errorf("%s is not a writeable directory", fp)
	}

	return true, nil
}

// isDirWriteable returns true if the supplied directory path is writeable,
// false otherwise
func isDirWriteable(fp string) bool {
	testPath := filepath.Join(fp, "test")
	f, err := os.Create(testPath)
	if err != nil {
		return false
	}
	f.Close()
	os.Remove(testPath)
	return true
}

// EnsureSemverPrefix takes a semver string and ensures the 'v' prefix is
// present.
func EnsureSemverPrefix(s string) string {
	// trim all leading 'v' runes (characters)
	s = strings.TrimLeft(s, "v")
	return fmt.Sprintf("v%s", s)
}

// GetSDKVersion returns the github.com/aws/aws-sdk-go-v2 version to use. It
// first tries to get the version from the --aws-sdk-go-version flag, then
// from the ack-generate-metadata.yaml and finally looks for the service
// go.mod controller.
func GetSDKVersion(
	awsSDKGoVersion string,
	lastGenerationVersion string,
	controllerRepoPath string,
) (string, error) {
	// First try to get the version from --aws-sdk-go-version flag
	if awsSDKGoVersion != "" {
		return awsSDKGoVersion, nil
	}

	// then, try to use last generation version (from ack-generate-metadata.yaml)
	if lastGenerationVersion != "" {
		return lastGenerationVersion, nil
	}

	// then, try to parse the service controller go.mod file
	sdkVersion, err := getSDKVersionFromGoMod(filepath.Join(controllerRepoPath, "go.mod"))
	if err == nil {
		return sdkVersion, nil
	}

	return "", err
}

// getSDKVersionFromGoMod parses a given go.mod file and returns
// the aws-sdk-go version in the required modules.
func getSDKVersionFromGoMod(goModPath string) (string, error) {
	b, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}
	goMod, err := modfile.Parse("", b, nil)
	if err != nil {
		return "", err
	}
	for _, require := range goMod.Require {
		if require.Mod.Path == sdkGoV2Module {
			return require.Mod.Version, nil
		}
	}
	return "", fmt.Errorf("couldn't find %s in the go.mod require block", sdkGoV2Module)
}
