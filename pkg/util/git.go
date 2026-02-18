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

package util

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

// CloneRepository clones a git repository into a given directory.
//
// Equivalent to: git clone $repositoryURL $path
func CloneRepository(ctx context.Context, path, repositoryURL string) error {
	cmd := exec.CommandContext(ctx, "git", "clone", repositoryURL, path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(out))
	}
	return nil
}

// HasTag checks if a tag exists in the local repository.
func HasTag(path string, tag string) bool {
	cmd := exec.Command("git", "-C", path, "rev-parse", "--verify", fmt.Sprintf("refs/tags/%s", tag))
	return cmd.Run() == nil
}

// FetchRepositoryTag fetches a single tag from the remote repository.
//
// Equivalent to: git -C $path fetch origin tag $tag
func FetchRepositoryTag(ctx context.Context, path string, tag string) error {
	cmd := exec.CommandContext(ctx, "git", "-C", path, "fetch", "origin", "tag", tag)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(out))
	}
	return nil
}

// CheckoutRepositoryTag checks out a repository tag.
//
// Equivalent to: git -C $path checkout tags/$tag -f
func CheckoutRepositoryTag(path string, tag string) error {
	cmd := exec.Command("git", "-C", path, "checkout", fmt.Sprintf("tags/%s", tag), "-f")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}
	return nil
}
