// Copyright 2022 Mia-Platform

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutils

import (
	"path"
)

const (
	// Invalid resources names
	InvalidFolderPath  = "/invalid/path"
	InvalidFileName    = "invalid.yaml"
	InvalidGroupName   = "invalid-group"
	InvalidClusterName = "invalid-cluster"

	// Valid resources names
	TestGroupName1       = "test-group"
	TestGroupName2       = "test-group2"
	TestClusterName1     = "test-cluster"
	TestClusterName2     = "test-cluster2"
	KustomizeTestDirName = "kustomize-test"
)

func GetTestFile(module string, args ...string) string {
	combinedElements := append([]string{
		"..",
		"..",
		"tests",
		module,
	},
		args...,
	)
	return path.Join(combinedElements...)
}
