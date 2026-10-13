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

import "time"

// Test hooks for the external sdk_test package. This file is compiled only
// for tests.

// MaxFetchAttempts is the number of times a model file is requested before
// a retryable status is returned to the caller.
const MaxFetchAttempts = maxFetchAttempts

// ParseRetryAfter exposes parseRetryAfter.
var ParseRetryAfter = parseRetryAfter

// SetFetchRetryDelays replaces the model fetch backoff delays and returns a
// function that restores the previous ones.
func SetFetchRetryDelays(base, max time.Duration) func() {
	prevBase, prevMax := fetchRetryBaseDelay, fetchRetryMaxDelay
	fetchRetryBaseDelay, fetchRetryMaxDelay = base, max
	return func() {
		fetchRetryBaseDelay, fetchRetryMaxDelay = prevBase, prevMax
	}
}
