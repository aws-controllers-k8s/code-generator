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
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// LoadRepository loads a repository from the local file system.
// TODO(a-hilaly): load repository into a memory filesystem (needs go1.16
// migration or use something like https://github.com/spf13/afero
func LoadRepository(path string) (*git.Repository, error) {
	return git.PlainOpen(path)
}

// CloneRepository clones a git repository into a given directory.
//
// Calling his function is equivalent to executing `git clone $repositoryURL $path`
func CloneRepository(ctx context.Context, path, repositoryURL string) error {
	_, err := git.PlainCloneContext(ctx, path, false, &git.CloneOptions{
		URL:      repositoryURL,
		Progress: nil,
		// Clone and fetch all tags
		Tags: git.AllTags,
	})
	return err
}

// FetchRepositoryTags fetches a repository remote tags.
//
// Calling this function is equivalent to executing `git -C $path fetch --all --tags`
func FetchRepositoryTags(ctx context.Context, path string) error {
	// PlainOpen will make the git commands run against the local
	// repository and directly make changes to it. So no need to
	// save/rewrite the refs
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	err = repo.FetchContext(ctx, &git.FetchOptions{
		Progress: nil,
		Tags:     git.AllTags,
	})
	// weirdly go-git returns a error "Already up to date" when all tags
	// are already fetched. We should ignore this error.
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}

// getRepositoryTagRef returns the git reference (commit hash) of a given tag.
// NOTE: It is not possible to checkout a tag without knowing it's reference.
//
// Calling this function is equivalent to executing `git rev-list -n 1 $tagName`
func getRepositoryTagRef(repo *git.Repository, tagName string) (*plumbing.Reference, error) {
	tagRefs, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	for {
		tagRef, err := tagRefs.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error finding tag reference: %v", err)
		}
		if tagRef.Name().Short() == tagName {
			return tagRef, nil
		}
	}
	return nil, errors.New("tag reference not found")
}

// CheckoutRepositoryTag checkouts a repository tag by looking for the tag
// reference then calling the checkout function.
//
// Calling This function is equivalent to executing `git checkout tags/$tag`
func CheckoutRepositoryTag(repo *git.Repository, tag string) error {

	tagRef, err := getRepositoryTagRef(repo, tag)

	if err != nil {
		return err
	}
	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	// AWS-SDK-GO-V2 - Hash value for tag not found, hence use tagName to checkout
	err = wt.Checkout(&git.CheckoutOptions{
		// Checkout only take hashes or branch names.
		//Hash:   tagRef.Hash(),
		Branch: tagRef.Name(),
		Force:  true,
	})
	return err
}
