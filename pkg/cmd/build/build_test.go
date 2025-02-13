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

package build

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/mia-platform/vab/pkg/cmd/util"
	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	testdata := "testdata"
	configPath := filepath.Join(testdata, "config.yaml")

	configFlags := util.NewConfigFlags()
	configFlags.ConfigPath = &configPath

	cmd := NewCommand(configFlags)
	assert.NotNil(t, cmd)

	buffer := new(bytes.Buffer)
	cmd.SetArgs([]string{"test-group2", testdata})
	cmd.SetOut(buffer)
	assert.NoError(t, cmd.Execute())
	t.Log(buffer.String())
}

func TestBuildRun(t *testing.T) {
	t.Parallel()

	testdata := "testdata"
	configFile := filepath.Join(testdata, "config.yaml")

	tests := map[string]struct {
		options        *Options
		expectedOutput string
		expectedError  string
	}{
		"missing files in cluster folder": {
			options: &Options{
				group:       "test-group",
				cluster:     "test-cluster",
				contextPath: testdata,
				configPath:  configFile,
			},
			expectedError: `building resources for "test-group/test-cluster":`,
		},
		"missing configuration file": {
			options: &Options{
				group:       "test-group",
				configPath:  filepath.Join(t.TempDir(), "missing.yaml"),
				contextPath: testdata,
			},
			expectedError: "reading config file:",
		},
		"build single cluster": {
			options: &Options{
				group:       "test-group2",
				cluster:     "test-cluster",
				contextPath: testdata,
				configPath:  configFile,
			},
			expectedOutput: `---
### BUILD RESULTS FOR: "test-group2/test-cluster" ###
apiVersion: v1
kind: Service
metadata:
  name: test
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: test
  type: ClusterIP
`,
		},
		"build entire group": {
			options: &Options{
				group:       "test-group2",
				contextPath: testdata,
				configPath:  configFile,
			},
			expectedOutput: `---
### BUILD RESULTS FOR: "test-group2/test-cluster" ###
apiVersion: v1
kind: Service
metadata:
  name: test
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: test
  type: ClusterIP
---
### BUILD RESULTS FOR: "test-group2/test-cluster2" ###
apiVersion: v1
kind: Service
metadata:
  name: test
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: test
  type: ClusterIP
`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buffer := new(bytes.Buffer)
			test.options.writer = buffer

			err := test.options.Run(t.Context())
			if len(test.expectedError) > 0 {
				assert.ErrorContains(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expectedOutput, buffer.String())
		})
	}
}
