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

// Validate package is used for validating a configuration files
package validate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mia-platform/vab/internal/testutils"
	"github.com/mia-platform/vab/pkg/logger"
	"github.com/stretchr/testify/assert"
)

// Test parsing error returned from ReadConfig
func TestValidateParseError(t *testing.T) {
	targetPath := testutils.GetTestFile("utils", "invalid_yaml.yaml")
	buffer := new(bytes.Buffer)
	logger := logger.DisabledLogger{}

	err := ConfigurationFile(logger, targetPath, buffer)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "error while parsing the configuration file:")
	}
}

// Test validation of valid empty config
func TestValidateEmptySpec(t *testing.T) {
	targetPath := testutils.GetTestFile("utils", "empty_config.yaml")
	buffer := new(bytes.Buffer)
	logger := logger.DisabledLogger{}
	err := ConfigurationFile(logger, targetPath, buffer)
	if assert.NoError(t, err) {
		assert.Equal(t, buffer.String(), expectedOutput1)
	}
}

// Test validation of wrong Kind/APIVersion
func TestCheckTypeMeta(t *testing.T) {
	targetPath := testutils.GetTestFile("validate", "invalidkind.yaml")
	buffer := new(bytes.Buffer)
	logger := logger.DisabledLogger{}

	err := ConfigurationFile(logger, targetPath, buffer)
	if assert.Error(t, err) {
		assert.Equal(t, buffer.String(), expectedOutput2)
	}
}

// Test validate with ad-hoc invalid file for getting all warning, info and errors
func TestValidateOutput(t *testing.T) {
	targetPath := testutils.GetTestFile("validate", "all-check-config.yaml")
	buffer := new(bytes.Buffer)
	logger := logger.DisabledLogger{}
	err := ConfigurationFile(logger, targetPath, buffer)
	if assert.Error(t, err) {
		loggedLines := strings.Split(buffer.String(), "\n")
		expectedOutput3Array := strings.Split(expectedOutput3, "\n")
		assert.Equal(t, len(loggedLines), len(expectedOutput3Array), "Wrong log lines founded")
		for _, line := range loggedLines {
			assert.Contains(t, expectedOutput3Array, line, "Unexpected log line")
		}
	}
}

const (
	expectedOutput1 = `[warn][default] no module found: check the config file if this behavior is unexpected
[warn][default] no add-on found: check the config file if this behavior is unexpected
[warn] no group found: check the config file if this behavior is unexpected
The configuration is valid!
`
	expectedOutput2 = `[error] wrong kind: WrongKind - expected: ClustersConfiguration
[error] wrong version: wrong.version.io/v1 - expected: vab.mia-platform.eu/v1alpha1
[warn][default] no module found: check the config file if this behavior is unexpected
[warn][default] no add-on found: check the config file if this behavior is unexpected
[warn] no group found: check the config file if this behavior is unexpected
`
	expectedOutput3 = `[error][default] missing version of module module-1/flavor-1
[warn][default] missing weight of module module-1/flavor-1: setting default (0)
[info][default] disabling module module-2/flavor-2
[error][default] missing version of add-on addon-1
[info][default] disabling add-on addon-2
[error] please specify a valid name for each group
[error][undefined] missing cluster name in group: please specify a valid name for each cluster
[error][undefined/undefined] missing cluster context: please specify a valid context for each cluster
[warn][undefined/undefined] no module found: check the config file if this behavior is unexpected
[warn][undefined/undefined] no add-on found: check the config file if this behavior is unexpected
[warn][undefined/cluster-1] no module found: check the config file if this behavior is unexpected
[warn][undefined/cluster-1] no add-on found: check the config file if this behavior is unexpected
[warn][group-1] no cluster found in group: check the config file if this behavior is unexpected
`
)
