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
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	DefaultConfigPath = "config.yaml"

	configPathFlagName      = "config"
	configPathFlagShortName = "c"
	configPathUsage         = "Path to the configuration file to use."
)

var (
	configValidExtensions = []string{"yaml", "yml"}
)

type ConfigFlags struct {
	ConfigPath *string
}

func NewConfigFlags() *ConfigFlags {
	stringPointer := func(str string) *string {
		return &str
	}
	return &ConfigFlags{
		ConfigPath: stringPointer(""),
	}
}

func (f *ConfigFlags) AddFlags(flags *pflag.FlagSet) {
	if f.ConfigPath != nil {
		flags.StringVarP(f.ConfigPath, configPathFlagName, configPathFlagShortName, *f.ConfigPath, configPathUsage)
		if err := cobra.MarkFlagFilename(flags, configPathFlagName, configValidExtensions...); err != nil {
			panic(err)
		}
	}
}
