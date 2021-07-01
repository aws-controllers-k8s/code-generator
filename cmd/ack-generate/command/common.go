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
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

const (
	sdkRepoURL = "https://github.com/aws/aws-sdk-go"
)

// ensureDir makes sure that a supplied directory exists and
// returns whether the directory already existed.
func ensureDir(fp string) (bool, error) {
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

// ensureSDKRepo ensures that we have a git clone'd copy of the aws-sdk-go
// repository, which we use model JSON files from. Upon successful return of
// this function, the sdkDir global variable will be set to the directory where
// the aws-sdk-go is found. It will also optionally fetch all the remote tags
// and checkout the given tag.
func ensureSDKRepo(
	cacheDir string,
	// A boolean instructing ensureSDKRepo whether to fetch the remote tags from
	// the upstream repository
	fetchTags bool,
) error {
	var err error
	srcPath := filepath.Join(cacheDir, "src")
	if err = os.MkdirAll(srcPath, os.ModePerm); err != nil {
		return err
	}

	// Clone repository if it doen't exist
	sdkDir = filepath.Join(srcPath, "aws-sdk-go")
	if _, err := os.Stat(sdkDir); os.IsNotExist(err) {
		err = util.CloneRepository(context.Background(), sdkDir, sdkRepoURL)
		if err != nil {
			return fmt.Errorf("canot clone repository: %v", err)
		}
	}

	// Fetch all tags
	if fetchTags {
		err = util.FetchRepositoryTags(context.Background(), sdkDir)
		if err != nil {
			return fmt.Errorf("cannot fetch tags: %v", err)
		}
	}

	// get sdkVersion and ensure it prefix
	// TODO(a-hilaly) Parse `ack-generate-metadata.yaml` and pass the aws-sdk-go
	// version here.
	sdkVersion, err := getSDKVersion("")
	if err != nil {
		return err
	}
	sdkVersion = ensureSemverPrefix(sdkVersion)

	repo, err := util.LoadRepository(sdkDir)
	if err != nil {
		return fmt.Errorf("cannot read local repository: %v", err)
	}

	// Now checkout the local repository.
	err = util.CheckoutRepositoryTag(repo, sdkVersion)
	if err != nil {
		return fmt.Errorf("cannot checkout tag: %v", err)
	}

	return err
}

// ensureSemverPrefix takes a semver string and tries to append the 'v'
// prefix if it's missing.
func ensureSemverPrefix(s string) string {
	// trim all leading 'v' runes (characters)
	s = strings.TrimLeft(s, "v")
	return fmt.Sprintf("v%s", s)
}

// getSDKVersion returns the github.com/aws/aws-sdk-go version to use. It
// first tries to get the version from the --aws-sdk-go-version flag, then
// from the ack-generate-metadata.yaml and finally look for the service
// go.mod controller.
func getSDKVersion(
	lastGenerationVersion string,
) (string, error) {
	// First try to get the version from --aws-sdk-go-version flag
	if optAWSSDKGoVersion != "" {
		return optAWSSDKGoVersion, nil
	}

	// then, try to use last generation version (from ack-generate-metadata.yaml)
	if lastGenerationVersion != "" {
		return lastGenerationVersion, nil
	}

	// then, try to parse the service controller go.mod file
	sdkVersion, err := getSDKVersionFromGoMod(filepath.Join(optOutputPath, "go.mod"))
	if err == nil {
		return sdkVersion, nil
	}

	return "", err
}

// getSDKVersionFromGoMod parses a given go.mod file and returns
// the aws-sdk-go version in the required modules.
func getSDKVersionFromGoMod(goModPath string) (string, error) {
	b, err := ioutil.ReadFile(goModPath)
	if err != nil {
		return "", err
	}
	goMod, err := modfile.Parse("", b, nil)
	if err != nil {
		return "", err
	}
	sdkModule := strings.TrimPrefix(sdkRepoURL, "https://")
	for _, require := range goMod.Require {
		if require.Mod.Path == sdkModule {
			return require.Mod.Version, nil
		}
	}
	return "", fmt.Errorf("couldn't find %s in the go.mod require block", sdkModule)
}
