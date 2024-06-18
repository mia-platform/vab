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

package util

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestWriteKustomizationData(t *testing.T) {
	t.Parallel()

	testdata := "testdata"
	tests := map[string]struct {
		path           string
		expectedString string
		expectedError  bool
	}{
		"read kustomization files": {
			path: filepath.Join(testdata, "kustomize"),
			expectedString: `apiVersion: v1
kind: Service
metadata:
  name: example
spec:
  ports:
  - port: 80
    targetPort: web
  selector:
    app: example
`,
		},
		"missing files": {
			path:          t.TempDir(),
			expectedError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buffer := new(bytes.Buffer)
			err := WriteKustomizationData(test.path, buffer)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.expectedString, buffer.String())
		})
	}
}

func TestGroupFromConfig(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.yaml")
	tests := map[string]struct {
		group         string
		path          string
		expectedGroup v1alpha1.Group
		expectedError string
	}{
		"invalid config path": {
			path:          filepath.Join(t.TempDir(), "missing"),
			expectedError: "cannot read config file",
		},
		"missing group in file": {
			path:          configPath,
			group:         "missing",
			expectedError: `cannot find "missing" group in config at path "testdata/config.yaml"`,
		},
		"group found": {
			path:  configPath,
			group: "test-group",
			expectedGroup: v1alpha1.Group{
				Name: "test-group",
				Clusters: []v1alpha1.Cluster{
					{
						Name:    "test-cluster",
						Modules: make(map[string]v1alpha1.Package),
						AddOns:  make(map[string]v1alpha1.Package),
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			group, err := GroupFromConfig(test.group, test.path)
			if len(test.expectedError) > 0 {
				assert.ErrorContains(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.expectedGroup, group)
		})
	}
}
