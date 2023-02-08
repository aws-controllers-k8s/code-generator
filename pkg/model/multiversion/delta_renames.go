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

package multiversion

import (
	"errors"
	"fmt"
)

var (
	// ErrProhibitedRename is returned by ComputeRenamesDelta when a prohibited renaming
	// pattern is detected.
	//
	// One example of prohibited renaming is: renaming X to Y in v1 and Z to Y in v2.
	ErrProhibitedRename = errors.New("prohibited rename")
)

// ComputeRenamesDelta returns a map representing the field renames map between two
// distinguished api versions.
//
// Examples:
//
//	if we rename X to Y in v1 and X to Y in v2 the map of renames is {}
//	if we rename X to Y in v1 and X to Z in v2 the map of renames is {Y: Z}
//	if we don't rename any field in v1 and we rename X to Y in v2 the map of renames if {X: Y}
func ComputeRenamesDelta(srcRenames, dstRenames map[string]string) (map[string]string, error) {
	// returns an error if we find any prohibited or unsupported renaming pattern.
	for dstOld, dstNew := range dstRenames {
		for srcOld, srcNew := range srcRenames {
			if dstOld != srcOld && dstNew == srcNew {
				errMsg := fmt.Sprintf("found conflicting renames %s:%s and %s:%s", srcOld, srcNew, dstOld, dstNew)
				return nil, fmt.Errorf("%v: %s", ErrAPIVersionRemoved, errMsg)
			}
		}
	}

	renames := make(map[string]string)
	// loop over the destination renames and check if the same rename exists in the
	// source renames.
	for dstOld, dstNew := range dstRenames {
		srcNew, ok := srcRenames[dstOld]
		if ok {
			// if src and dst have the same renames - that's not a rename.
			if srcNew != dstNew {
				renames[srcNew] = dstNew
			}
			continue
		}
		renames[dstOld] = dstNew
	}

	// loop over the source renames and check the ones that don't exist in destination renames.
	for srcOld, srcNew := range srcRenames {
		_, ok := dstRenames[srcOld]
		if !ok {
			renames[srcNew] = srcOld
		}
	}
	return renames, nil
}
