// Copyright Mia srl
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apply

import (
	"fmt"
	"io/fs"
	"testing"

	jpl "github.com/mia-platform/jpl/deploy"
	"github.com/mia-platform/vab/internal/testutils"
	"github.com/mia-platform/vab/pkg/logger"
	"github.com/stretchr/testify/assert"
)

const (
	crdDefaultRetries  = 10
	testBuildFolder    = "apply-test"
	testConfigFileName = "testconfig.yaml"
)

func TestWrongContextPath(t *testing.T) {
	log := logger.DisabledLogger{}
	configPath := testutils.GetTestFile("apply", testBuildFolder, testConfigFileName)
	options := jpl.NewOptions()
	err := Apply(log, configPath, "", "", testutils.InvalidFolderPath, options, crdDefaultRetries)

	if assert.Error(t, err) {
		assert.ErrorIs(t, err, fs.ErrNotExist)
	}

	err = Apply(log, configPath, "", "", configPath, options, crdDefaultRetries)
	if assert.Error(t, err) {
		assert.Equal(t, err.Error(), fmt.Sprintf("the target path %s is not a directory", configPath))
	}
}

func TestBuildInvalidConfigPath(t *testing.T) {
	log := logger.DisabledLogger{}
	contextPath := testutils.GetTestFile("apply", testBuildFolder)
	options := jpl.NewOptions()
	err := Apply(log, testutils.InvalidFileName, "", "", contextPath, options, crdDefaultRetries)

	if assert.Error(t, err) {
		assert.ErrorIs(t, err, fs.ErrNotExist)
	}
}

func TestBuildInvalidKustomization(t *testing.T) {
	log := logger.DisabledLogger{}
	configPath := testutils.GetTestFile("apply", testBuildFolder, testConfigFileName)
	options := jpl.NewOptions()
	err := Apply(log, configPath, testutils.TestGroupName1, testutils.TestClusterName1, testutils.GetTestFile("apply", testBuildFolder), options, crdDefaultRetries)
	assert.Error(t, err)
}

func TestGetContextError(t *testing.T) {
	configPath := testutils.GetTestFile("apply", testBuildFolder, testConfigFileName)
	_, err := getContext(configPath, "notExistent", "test-cluster")
	assert.Error(t, err)

	_, err = getContext(configPath, "test-group2", "notExistent")
	assert.Error(t, err)

	configPathError := "notExistent"
	_, err = getContext(configPathError, "test-group", "test-cluster")
	assert.Error(t, err)
}
