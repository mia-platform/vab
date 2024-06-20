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

package sync

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/mia-platform/vab/internal/git"
	"github.com/mia-platform/vab/pkg/cmd/util"
	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	configFlags := util.NewConfigFlags()

	cmd := NewCommand(configFlags)
	assert.NotNil(t, cmd)
}

func TestToOptions(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	testCases := map[string]struct {
		args            []string
		dryRun          bool
		configPath      string
		expectedOptions *Options
		expectedError   string
	}{
		"invalid context path": {
			args:          []string{filepath.Join("invalid", "path")},
			expectedError: "no such file or directory",
		},
		"return options": {
			args:       []string{tempDir},
			configPath: "custom.yaml",
			dryRun:     true,
			expectedOptions: &Options{
				contextPath: tempDir,
				dryRun:      true,
				configPath:  "custom.yaml",
			},
		},
		"no config path": {
			args: []string{tempDir},
			expectedOptions: &Options{
				contextPath: tempDir,
				dryRun:      false,
				configPath:  "",
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			flags := Flags{
				dryRun: testCase.dryRun,
			}
			configFlags := util.NewConfigFlags()
			configFlags.ConfigPath = &testCase.configPath

			options, err := flags.ToOptions(configFlags, testCase.args)
			switch len(testCase.expectedError) {
			case 0:
				assert.NoError(t, err)
				assert.NotNil(t, options.filesGetter)
				options.filesGetter = nil
				assert.Equal(t, testCase.expectedOptions, options)
			default:
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

func TestRun(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.yaml")
	tests := map[string]struct {
		options       *Options
		expectedError string
	}{
		"clone packages": {
			options: &Options{
				configPath:  configPath,
				contextPath: t.TempDir(),
				filesGetter: git.NewTestFilesGetter(t),
			},
		},
		"don't clone packages": {
			options: &Options{
				configPath:  configPath,
				contextPath: t.TempDir(),
				dryRun:      true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.options.Run(context.TODO())
			if len(test.expectedError) > 0 {
				assert.ErrorContains(t, err, test.expectedError)
				return
			}

			assert.NoError(t, err)
		})
	}
}
