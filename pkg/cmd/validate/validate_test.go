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

package validate

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mia-platform/vab/pkg/cmd/util"
	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "valid.yaml")

	configFlags := util.NewConfigFlags()
	configFlags.ConfigPath = &configPath
	cmd := NewCommand(configFlags)
	assert.NotNil(t, cmd)

	buffer := new(bytes.Buffer)
	cmd.SetOut(buffer)
	assert.NoError(t, cmd.Execute())
	t.Log(buffer.String())
}

func TestValidationTextGeneration(t *testing.T) {
	t.Parallel()
	testdata := "testdata"

	tests := map[string]struct {
		options        *Options
		expectedString string
		expectedError  string
	}{
		"invalid config file": {
			options: &Options{
				configPath: filepath.Join(testdata, "invalid_yaml.yaml"),
			},
			expectedString: "",
			expectedError:  "parsing configuration file",
		},
		"empty config file": {
			options: &Options{
				configPath: filepath.Join(testdata, "empty_config.yaml"),
			},
			expectedString: `[warn][default] no module found: check the config file if this behavior is unexpected
[warn][default] no addon found: check the config file if this behavior is unexpected
[warn] no group found: check the config file if this behavior is unexpected
The configuration is valid!
`,
		},
		"invalind kind": {
			options: &Options{
				configPath: filepath.Join(testdata, "invalidkind.yaml"),
			},
			expectedString: `[error] wrong kind: WrongKind - expected: ClustersConfiguration
[error] wrong version: wrong.version.io/v1 - expected: vab.mia-platform.eu/v1alpha1
[warn][default] no module found: check the config file if this behavior is unexpected
[warn][default] no addon found: check the config file if this behavior is unexpected
[warn] no group found: check the config file if this behavior is unexpected
`,
			expectedError: "configuration is invalid",
		},
		"all checks": {
			options: &Options{
				configPath: filepath.Join(testdata, "all-check-config.yaml"),
			},
			expectedString: `[error][default] missing version of module category/module-1
[info][default] disabling module category/module-2
[error][default] missing version of addon category/addon-1
[info][default] disabling addon category/addon-2
[error] please specify a valid name for each group
[error][undefined] missing cluster name in group: please specify a valid name for each cluster
[error][undefined/undefined] missing cluster context: please specify a valid context for each cluster
[warn][undefined/undefined] no module found: check the config file if this behavior is unexpected
[warn][undefined/undefined] no addon found: check the config file if this behavior is unexpected
[warn][undefined/cluster-1] no module found: check the config file if this behavior is unexpected
[warn][undefined/cluster-1] no addon found: check the config file if this behavior is unexpected
[warn][group-1] no cluster found in group: check the config file if this behavior is unexpected
`,
			expectedError: "configuration is invalid",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buffer := new(bytes.Buffer)
			test.options.writer = buffer

			err := test.options.Run(context.TODO())
			if len(test.expectedError) > 0 {
				assert.ErrorContains(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			loggedLines := strings.Split(buffer.String(), "\n")
			expectedOutputArray := strings.Split(test.expectedString, "\n")
			assert.Equal(t, len(loggedLines), len(expectedOutputArray), "Wrong log lines founded")
			for _, line := range loggedLines {
				assert.Contains(t, expectedOutputArray, line, "Unexpected log line")
			}
		})
	}
}
