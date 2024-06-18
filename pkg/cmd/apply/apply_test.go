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
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	jpltesting "github.com/mia-platform/jpl/pkg/testing"
	jplutil "github.com/mia-platform/jpl/pkg/util"
	"github.com/mia-platform/vab/pkg/cmd/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	flowcontrolapi "k8s.io/api/flowcontrol/v1beta3"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest/fake"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	configFlags := util.NewConfigFlags()

	cmd := NewCommand(configFlags)
	assert.NotNil(t, cmd)
}

func TestFlagsToOptions(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	configFile := "path/to/file.yaml"

	tests := map[string]struct {
		flags           *Flags
		configFlags     *util.ConfigFlags
		args            []string
		expectedOptions *Options
		expectedError   string
	}{
		"empty arg return error": {
			flags:         &Flags{},
			configFlags:   util.NewConfigFlags(),
			args:          []string{},
			expectedError: fmt.Sprintf("at least %d arguments are needed", minArgs),
		},
		"only one arg return error": {
			flags:         &Flags{},
			configFlags:   util.NewConfigFlags(),
			args:          []string{"first"},
			expectedError: fmt.Sprintf("at least %d arguments are needed", minArgs),
		},
		"invalid context path return error": {
			flags:         &Flags{timeout: "invalid"},
			configFlags:   util.NewConfigFlags(),
			args:          []string{"first", filepath.Join("invalid", "path")},
			expectedError: "error locating files",
		},
		"invalid timeout return error": {
			flags:         &Flags{timeout: "invalid"},
			configFlags:   util.NewConfigFlags(),
			args:          []string{"first", tmpDir},
			expectedError: "failed to parse request timeout",
		},
		"two arguments": {
			flags:       &Flags{timeout: timeoutDefaultValue},
			configFlags: util.NewConfigFlags(),
			args:        []string{"first", tmpDir},
			expectedOptions: &Options{
				fieldManager: "vab",
				group:        "first",
				contextPath:  tmpDir,
				configPath:   "",
			},
		},
		"three arguments": {
			flags:       &Flags{timeout: timeoutDefaultValue, dryRun: true},
			configFlags: &util.ConfigFlags{ConfigPath: &configFile},
			args:        []string{"first", "second", tmpDir},
			expectedOptions: &Options{
				fieldManager: "vab",
				dryRun:       true,
				group:        "first",
				cluster:      "second",
				contextPath:  tmpDir,
				configPath:   configFile,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			opts, err := test.flags.ToOptions(test.configFlags, test.args)
			if len(test.expectedError) > 0 {
				assert.ErrorContains(t, err, test.expectedError)
				assert.Nil(t, opts)
				return
			}

			assert.NoError(t, err)
			// TODO: find a better way to comparing things?
			// check that factoryAndConfigFunc is not nil to avoid missing the assignment
			assert.NotNil(t, opts.factoryAndConfigFunc)
			// remove function to allow easy comparison between objects
			opts.factoryAndConfigFunc = nil
			assert.Equal(t, test.expectedOptions, opts)
		})
	}
}

func TestApplyRun(t *testing.T) {
	t.Parallel()

	testdata := "testdata"
	configPath := filepath.Join(testdata, "testconfig.yaml")

	tests := map[string]struct {
		options                  *Options
		client                   *fake.RESTClient
		expectedError            string
		returnErrorInLocalServer bool
	}{
		"missing group in config return error": {
			options: &Options{
				group:       "missing",
				contextPath: testdata,
				configPath:  configPath,
			},
			expectedError: `cannot find "missing" group in config at path "testdata/testconfig.yaml"`,
		},
		"no cluster inside a group return error": {
			options: &Options{
				group:       "no-clusters",
				contextPath: testdata,
				configPath:  configPath,
			},
			expectedError: `group "no-clusters" doesn't have any cluster`,
		},
		"group does't have the specified cluster return error": {
			options: &Options{
				group:       "test-group",
				cluster:     "missing",
				contextPath: testdata,
				configPath:  configPath,
			},
			expectedError: `group "test-group" doesn't have cluster "missing"`,
		},
		"missing context in cluster return error": {
			options: &Options{
				group:       "test-group",
				cluster:     "test-cluster",
				contextPath: testdata,
				configPath:  configPath,
			},
			expectedError: `error executing apply for "test-group/test-cluster": no context found`,
		},
		"error checking flowcontrol API return error": {
			options: &Options{
				group:       "test-group2",
				cluster:     "test-cluster",
				contextPath: testdata,
				configPath:  configPath,
			},
			returnErrorInLocalServer: true,
			expectedError:            `error executing apply for "test-group2/test-cluster": checking flowcontrol api`,
		},
		"invalid context path return error": {
			options: &Options{
				group:       "test-group2",
				cluster:     "test-cluster",
				contextPath: filepath.Join("invalid", "path"),
				configPath:  configPath,
			},
			expectedError: "must build at directory: not a valid directory",
		},
		"successful apply": {
			options: &Options{
				group:       "test-group2",
				cluster:     "test-cluster",
				contextPath: testdata,
				configPath:  configPath,
			},
			client: &fake.RESTClient{
				Client: fake.CreateHTTPClient(func(*http.Request) (*http.Response, error) {
					return &http.Response{StatusCode: http.StatusOK, Header: jpltesting.DefaultHeaders()}, nil
				}),
			},
		},
		// "": {},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()

			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(r http.ResponseWriter, _ *http.Request) {
				if test.returnErrorInLocalServer {
					server.CloseClientConnections()
					return
				}

				r.Header().Add(flowcontrolapi.ResponseHeaderMatchedFlowSchemaUID, "unused")
				r.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			factory := jpltesting.NewTestClientFactory()
			factory.Client = test.client
			restConfig, err := factory.ToRESTConfig()
			require.NoError(t, err)
			restConfig.Host = server.URL

			test.options.factoryAndConfigFunc = func(string) (jplutil.ClientFactory, *genericclioptions.ConfigFlags) {
				return factory, genericclioptions.NewConfigFlags(false)
			}

			err = test.options.Run(ctx)
			switch len(test.expectedError) {
			case 0:
				assert.NoError(t, err)
			default:
				assert.ErrorContains(t, err, test.expectedError)
			}
		})
	}
}
