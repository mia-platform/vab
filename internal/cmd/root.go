// Copyright 2022 Mia-Platform

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

type FlagPole struct {
	Name   string
	Config string
}

var flags = &FlagPole{}

// NewRootCommand returns a new cobra.Command for vab's root command
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "vab",
		Version: versionString(),
		Short:   "A tool for installing the Mia-Platform distro on your clusters",
	}

	rootCmd.AddCommand(NewInitCommand())
	rootCmd.AddCommand(NewValidateCommand())
	rootCmd.AddCommand(NewSyncCommand())
	rootCmd.AddCommand(NewBuildCommand())
	rootCmd.AddCommand(NewApplyCommand())

	return rootCmd
}

// Version is dynamically set by the ci or overridden by the Makefile.
var Version = "DEV"

// BuildDate is dynamically set at build time by the cli or overridden in the Makefile.
var BuildDate = "" // YYYY-MM-DD

func versionString() string {
	version := Version

	if BuildDate != "" {
		version = fmt.Sprintf("%s (%s)", version, BuildDate)
	}

	// don't return GoVersion during a test run for consistent test output
	if flag.Lookup("test.v") != nil {
		return version
	}

	return fmt.Sprintf("%s, Go Version: %s", version, runtime.Version())
}
