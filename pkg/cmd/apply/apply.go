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
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/go-logr/logr"
	jplclient "github.com/mia-platform/jpl/pkg/client"
	"github.com/mia-platform/jpl/pkg/event"
	"github.com/mia-platform/jpl/pkg/flowcontrol"
	"github.com/mia-platform/jpl/pkg/inventory"
	"github.com/mia-platform/jpl/pkg/resourcereader"
	jplutil "github.com/mia-platform/jpl/pkg/util"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/mia-platform/vab/pkg/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
)

const (
	shortCmd = "Build and apply the local configuration"
	longCmd  = `Builds and applies the local configuration to the specified cluster or group,
	or to all of them`
	cmdUsage = "apply GROUP [CLUSTER] CONTEXT"

	dryRunDefaultValue = false
	dryRunFlagName     = "dry-run"
	dryRunUsage        = "if true does not apply the configurations"

	timeoutDefaultValue = "0s"
	timeoutFlagName     = "timeout"
	timeoutFlagUsage    = `the length of time to wait before giving up.
	Non-zero values should contain a corresponding
	time unit (e.g. 1s, 2m, 3h). A value of zero means
	don't timeout requests.`

	applyErrorFormat = "applying resources for %q: %w"

	minArgs = 2
	maxArgs = 3
)

// Flags contains all the flags for the `apply` command. They will be converted to Options
// that contains all runtime options for the command.
type Flags struct {
	dryRun  bool
	timeout string
}

// AddFlags set the connection between Flags property to command line flags
func (f *Flags) AddFlags(flags *pflag.FlagSet) {
	flags.BoolVar(&f.dryRun, dryRunFlagName, dryRunDefaultValue, heredoc.Doc(dryRunUsage))
	flags.StringVar(&f.timeout, timeoutFlagName, timeoutDefaultValue, heredoc.Doc(timeoutFlagUsage))
}

type factoryAndConfigFunc func(context string) (jplutil.ClientFactory, *genericclioptions.ConfigFlags)

func defaultFactoryAndConfigfunc(context string) (jplutil.ClientFactory, *genericclioptions.ConfigFlags) {
	config := genericclioptions.NewConfigFlags(true)
	config.Context = &context
	factory := jplutil.NewFactory(config)
	return factory, config
}

// Options have the data required to perform the apply operation
type Options struct {
	dryRun               bool
	timeout              time.Duration
	fieldManager         string
	group                string
	cluster              string
	contextPath          string
	configPath           string
	factoryAndConfigFunc factoryAndConfigFunc
	logger               logr.Logger
}

func NewCommand(cf *util.ConfigFlags) *cobra.Command {
	flags := &Flags{}

	cmd := &cobra.Command{
		Use:   cmdUsage,
		Short: heredoc.Doc(shortCmd),
		Long:  heredoc.Doc(longCmd),

		Args: cobra.RangeArgs(minArgs, maxArgs),
		Run: func(cmd *cobra.Command, args []string) {
			options, err := flags.ToOptions(cf, args)
			cobra.CheckErr(err)
			cobra.CheckErr(options.Run(cmd.Context()))
		},
	}

	flags.AddFlags(cmd.Flags())
	return cmd
}

// ToOptions transform the command flags in command runtime arguments
func (f *Flags) ToOptions(cf *util.ConfigFlags, args []string) (*Options, error) {
	group := args[0]
	cluster := ""
	contextPath := args[len(args)-1]
	if len(args) >= maxArgs {
		cluster = args[1]
	}

	cleanedContextPath, err := util.ValidateContextPath(contextPath)
	if err != nil {
		return nil, err
	}

	var timeout time.Duration
	if timeout, err = time.ParseDuration(f.timeout); err != nil {
		return nil, fmt.Errorf("failed to parse request timeout: %w", err)
	}

	configPath := ""
	if cf.ConfigPath != nil && len(*cf.ConfigPath) > 0 {
		configPath = filepath.Clean(*cf.ConfigPath)
	}

	return &Options{
		dryRun:               f.dryRun,
		timeout:              timeout,
		fieldManager:         "vab",
		group:                group,
		cluster:              cluster,
		contextPath:          cleanedContextPath,
		configPath:           configPath,
		factoryAndConfigFunc: defaultFactoryAndConfigfunc,
	}, nil
}

// Run execute the apply command
func (o *Options) Run(ctx context.Context) error {
	o.logger = logr.FromContextOrDiscard(ctx)

	group, err := util.GroupFromConfig(o.group, o.configPath)
	if err != nil {
		return err
	}

	found := false
	for _, cluster := range group.Clusters {
		clusterName := cluster.Name
		if o.cluster != "" && clusterName != o.cluster {
			continue
		}

		found = true

		applyCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		clusterID := util.ClusterID(o.group, clusterName)
		clusterLogger := o.logger.WithName(clusterID)
		clusterLogger.V(2).Info("applying files")

		eventCh, err := o.apply(applyCtx, cluster)
		if err != nil {
			return err
		}

		for {
			select {
			case event, open := <-eventCh:
				if !open {
					clusterLogger.V(2).Info("finish applying files")
					return nil
				}

				fmt.Fprintf(os.Stderr, "%s: %s\n", clusterID, event.String())
			case <-applyCtx.Done():
				return applyCtx.Err()
			}
		}
	}

	switch {
	case !found && len(o.cluster) == 0:
		return fmt.Errorf("group %q doesn't have any cluster", o.group)
	case !found && len(o.cluster) != 0:
		return fmt.Errorf("group %q doesn't have cluster %q", o.group, o.cluster)
	}
	return nil
}

func (o *Options) apply(ctx context.Context, cluster v1alpha1.Cluster) (<-chan event.Event, error) {
	clusterID := util.ClusterID(o.group, cluster.Name)
	if len(cluster.Context) == 0 {
		return nil, fmt.Errorf(applyErrorFormat, clusterID, fmt.Errorf("no context found"))
	}

	factory, err := o.factoryFor(clusterID, cluster.Context)
	if err != nil {
		return nil, fmt.Errorf(applyErrorFormat, clusterID, err)
	}

	return o.applyManifests(ctx, factory, cluster.Name)
}

// factoryFor return a rest.Config for connecting to the clusterID with context name
func (o *Options) factoryFor(clusterID, kubeContext string) (jplutil.ClientFactory, error) {
	factory, config := o.factoryAndConfigFunc(kubeContext)
	restConfig, err := factory.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	clusterLogger := o.logger.WithName(clusterID)
	var enabled bool
	clusterLogger.V(5).Info("checking flowcontrol APIs availability")
	if enabled, err = flowcontrol.IsEnabled(context.Background(), restConfig); err != nil {
		return nil, fmt.Errorf("flowcontrol api: %w", err)
	}

	clusterLogger.V(5).Info("flowcontrol APIs status", "enabled", enabled)
	if enabled {
		config.WrapConfigFn = func(c *rest.Config) *rest.Config {
			c.QPS = -1
			c.Burst = -1
			return c
		}
	}

	return factory, nil
}

func (o *Options) applyManifests(ctx context.Context, factory jplutil.ClientFactory, clusterName string) (<-chan event.Event, error) {
	path := filepath.Join(o.contextPath, util.ClusterPath(o.group, clusterName))
	clusterLogger := o.logger.WithName(util.ClusterID(o.group, clusterName))

	clusterLogger.V(2).Info("reading manifests", "path", path)
	manifests, err := readManifests(factory, path)
	if err != nil {
		return nil, err
	}

	clusterLogger.V(2).Info("finish reading manifests", "path", path)
	inventoryName := "eu.mia-platform.vab.resourcestorage"
	inventory, err := inventory.NewConfigMapStore(factory, inventoryName, metav1.NamespaceSystem, o.fieldManager)
	if err != nil {
		return nil, err
	}

	applier, err := jplclient.NewBuilder().
		WithFactory(factory).
		WithInventory(inventory).
		Build()
	if err != nil {
		return nil, err
	}

	return applier.Run(ctx, manifests, jplclient.ApplierOptions{
		DryRun:       o.dryRun,
		FieldManager: o.fieldManager,
		Timeout:      o.timeout,
	}), nil
}

// readManifests return the manifests array that are read at path
func readManifests(factory jplutil.ClientFactory, path string) ([]*unstructured.Unstructured, error) {
	buffer := new(bytes.Buffer)
	if err := util.WriteKustomizationData(path, buffer); err != nil {
		return nil, err
	}

	reader, err := resourcereader.
		NewResourceReaderBuilder(factory).
		ResourceReader(buffer, resourcereader.StdinPath)
	if err != nil {
		return nil, err
	}

	return reader.Read()
}
