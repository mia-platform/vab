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

package create

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	t.Parallel()
	testFolder := t.TempDir()

	cmd := NewCommand()
	assert.NotNil(t, cmd)

	cmd.SetArgs([]string{testFolder})
	assert.NoError(t, cmd.Execute())
}

func TestToOptions(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	testCases := map[string]struct {
		args          []string
		expectedError string
		expectedPath  string
	}{
		"flags with one args": {
			args:         []string{tempDir},
			expectedPath: tempDir,
		},
		"flags with more than one args": {
			args:         []string{tempDir, "/another/path"},
			expectedPath: tempDir,
		},
		"flags with . path": {
			args: []string{"."},
			expectedPath: func() string {
				dir, err := os.Getwd()
				require.NoError(t, err)
				return dir
			}(),
		},
		"missing path": {
			args: []string{filepath.Join("invalid", "path")},
			expectedPath: func() string {
				dir, err := os.Getwd()
				require.NoError(t, err)
				return filepath.Join(dir, "invalid", "path")
			}(),
		},
		"path is not a directory": {
			args: func() []string {
				filePath := filepath.Join(tempDir, "filename")
				_, err := os.Create(filePath)
				require.NoError(t, err)
				return []string{filePath}
			}(),
			expectedError: "is not a directory",
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			flags := Flags{}
			options, err := flags.ToOptions(testCase.args)
			if testCase.expectedError != "" {
				assert.ErrorContains(t, err, testCase.expectedError)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedPath, options.path)
		})
	}
}

func TestCreateValidArgs(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		args                []string
		expectedCompletions []string
		expectedDirective   cobra.ShellCompDirective
	}{
		"no argument provided, return project path completion": {
			expectedCompletions: cobra.AppendActiveHelp([]string{}, pathArgHelpText),
			expectedDirective:   cobra.ShellCompDirectiveDefault,
		},
		"single argument provided, return no more argument error": {
			args:                []string{"argument"},
			expectedCompletions: cobra.AppendActiveHelp([]string{}, tooManyArgsHelpText),
			expectedDirective:   cobra.ShellCompDirectiveNoFileComp,
		},
		"more than one argument provided, return no more argument error": {
			args:                []string{"argument1", "argument2"},
			expectedCompletions: cobra.AppendActiveHelp([]string{}, tooManyArgsHelpText),
			expectedDirective:   cobra.ShellCompDirectiveNoFileComp,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			completions, directive := validArgs(nil, test.args, "")
			assert.Equal(t, test.expectedCompletions, completions)
			assert.Equal(t, test.expectedDirective, directive)
		})
	}
}
