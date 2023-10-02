# Commands

## `init <project-name>`

Creates the folder `project-name` with a preliminary directory structure,
together with the skeleton of the configuration file.

The project directory will include the `clusters` directory (including the `all-group` folder with a minimal
kustomize configuration), and the configuration file.

### Directory structure

```txt
./project-name
├── clusters
|   └── all-groups
|       └── kustomization.yaml
└── config.yaml
```

### Configuration file

```yaml
apiVersion: vab.mia-platform.eu/v1alpha1
kind: ClustersConfiguration
name: project-name
spec:
  modules: {}
    # example-module:
    #   version: 1.20.1
  addons: {}
    # example-addon:
    #   version: 1.20.1
  groups: []
    # - name: group-name
    #   clusters:
    #     - name: cluster1-name
    #       context: context1-name
    #       addons:
    #         example-addon: --- OVERRIDE ADD-ON
    #           version: 1.20.100
    #     - name: cluster2-name
    #       context: context2-name
    #       modules:
    #         example-module: --- DISABLE MODULE
    #           disable: true
    #         other-module: --- ADD NEW MODULE
    #           version: 1.20.20
    #     - name: cluster3-name
    #       context: context3-name
```

## `validate <config-file-path>`

Validates the configuration contained in `<config-file-path>`.

It returns an error if the `<config-file-path>` is malformed or includes resources that do not exist in our catalogue.

## `sync`

Fetches new and updated vendor versions and updates the clusters configuration locally to the latest changes
of the configuration file.

After the execution, the `vendors` folder will include the new and updated  modules/add-ons (if not already present),
and the directory structure inside the `clusters` folder will be updated according to the current configuration.

Example:

```yaml
apiVersion: vab.mia-platform.eu/v1alpha1
kind: ClustersConfiguration
name: project-name
spec:
  modules:
    ingress/traefik/base:
      version: 1.20.1
    cni/cilium/base:
      version: 1.20.1
  groups:
    - group: group-name
      clusters:
        - cluster1-name: context1-name
          modules:
            cni/cilium/base:
              disable: true
            cni/calico/base:
              version: 1.20.20
        - cluster2-name: context2-name
```

The command execution will build the following directory structure:

```txt
./project-name
├── vendors
|   ├── modules
|   |   ├── ingress
|   |   |   └── traefik
|   |   |       └── base
|   |   |           ├── [...]
|   |   |           └── kustomization.yaml
|   |   └── cni
|   |       ├── cilium
|   |       |   └── base
|   |       |       ├── [...]
|   |       |       └── kustomization.yaml
|   |       └── calico
|   |           └── base
|   |               ├── [...]
|   |               └── kustomization.yaml
|   └── add-ons
|       └── ..
├── clusters
|   ├── all-groups
|   |   ├── [...]
|   |   └── kustomization.yaml
|   └── group-name
|       ├── all-clusters
|       |   ├── ingress
|       |   |   └── traefik
|       |   ├── cni
|       |   |   └── cilium
|       |   └── kustomization.yaml
|       ├── cluster1-name
|       |   ├── cni
|       |   |   └── calico
|       |   └── kustomization.yaml  
|       └── cluster2-name
|           └── kustomization.yaml
└── config.yaml
```

## `build <group-name> <cluster-name> <config-file>`

Runs `kustomize build` for the specified cluster or group.

It returns the full configuration locally without applying it to the cluster, allowing the user to check
if all the resources are generated correctly for the target cluster.

## `apply <group-name> <cluster-name> <config-file>`

Builds and applies the local configuration to the specified cluster, group, or to all of them.
The command builds the configurations using the same function as the `build` command.
The apply uses the `kubectl apply` command, creating a yaml file with the resources and then applying it to a
KUBECONFIG context with the same name of the `ClusterName` specified in the main configuration file.

If the cluster has `ClusterName` that satisfies the regex `test*`, the configuration file will skip
the `kubectl apply` step.
