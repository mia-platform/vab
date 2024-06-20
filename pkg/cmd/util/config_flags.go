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
	configPathFlagName      = "config"
	configPathFlagShortName = "c"
	configPathUsage         = "path to the configuration file to use"

	verboseFlagName      = "verbose"
	verboseFlagShortName = "v"
	verboseUsage         = "setting logging verbosity; use number between 0 and 10"
)

var (
	configValidExtensions = []string{"yaml", "yml"}
)

type ConfigFlags struct {
	ConfigPath *string
	Verbose    *int
}

func NewConfigFlags() *ConfigFlags {
	stringPointer := func(str string) *string {
		return &str
	}
	intPointer := func(number int) *int {
		return &number
	}
	return &ConfigFlags{
		ConfigPath: stringPointer(""),
		Verbose:    intPointer(0),
	}
}

func (f *ConfigFlags) AddFlags(flags *pflag.FlagSet) {
	if f.ConfigPath != nil {
		flags.StringVarP(f.ConfigPath, configPathFlagName, configPathFlagShortName, *f.ConfigPath, configPathUsage)
		if err := cobra.MarkFlagFilename(flags, configPathFlagName, configValidExtensions...); err != nil {
			panic(err)
		}
	}

	if f.Verbose != nil {
		flags.IntVarP(f.Verbose, verboseFlagName, verboseFlagShortName, *f.Verbose, verboseUsage)
	}
}
