# Commands

## `init <project-name>`

Creates the folder `project-name` with a preliminary directory structure, together with the skeleton of the configuration file.

The project directory will include the `vendors` folder (either empty or containing the `modules`/`add-ons` folders), the `overrides` directory (either empty or including a folder with a minimal default configuration), and the configuration file.

**Directory structure:**

    ./project-name
    ├── vendors
    |   ├── modules
    |   |   └── ..
    |   └── add-ons
    |       └── ..
    ├── overrides
    |   └── all-clusters
    |       └── kustomization.yaml
    └── config.yaml

**Configuration file:**

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: ClustersConfiguration
metadata:
  name: project-name
spec:
  modules: {}
    # example-module:
    #   version: 1.20.1
    #   weight: 10
  addons: {}
    # example-addon:
    #   version: 1.20.1
    #   weight: 1
  groups: []
    # - group: group-name
    #   clusters:
    #     - cluster1-name: context1-name
    #       addons:
    #         example-addon: --- OVERRIDE ADD-ON
    #           version: 1.20.100
    #     - cluster2-name: context2-name
    #       modules:
    #         example-module: --- DISABLE MODULE
    #           enabled: false
    #         other-module: --- ADD NEW MODULE
    #           version: 1.20.20 
    #           weight: 1
    #     - cluster3-name: context3-name
```

## `validate <config-file-path>`

Validates the configuration.

It returns an error if the `config.yaml` is malformed or includes resources that do not exist in our catalogue.

## `sync`

Fetches new vendors and updates the clusters configuration locally to the latest changes of the configuration file.

After the execution, the `vendors` folder will include the new modules/add-ons (if not already present), and the directory structure inside the `overrides` folder will be updated according to the current configuration.

Example:

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: ClustersConfiguration
metadata:
  name: project-name
spec:
  modules:
    ingress/traefik:
      version: 1.20.1
      weight: 10
    cni/cilium:
      version: 1.20.1
      weight: 1
  groups: 
    - group: group-name
      clusters:
        - cluster1-name: context1-name
          modules:
            cni/cilium:
              enabled: false
            cni/calico:
              version: 1.20.20 
              weight: 1
        - cluster2-name: context2-name
```

The command execution will build the following directory structure:

    ./project-name
    ├── vendors
    |   ├── modules
    |   |   ├── ingress
    |   |   |   └── traefik
    |   |   |       ├── [...]
    |   |   |       └── kustomization.yaml
    |   |   └── cni
    |   |       ├── cilium
    |   |       |   ├── [...]
    |   |       |   └── kustomization.yaml
    |   |       └── calico
    |   |           ├── [...]
    |   |           └── kustomization.yaml
    |   └── add-ons
    |       └── ..
    ├── overrides
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

## `build <group-name> <cluster-name>`

Runs `kustomize build` for the specified cluster.

It returns the full configuration locally without applying it to the cluster, allowing the user to check it beforehand.

## `apply <group-name> <cluster-name>`

Builds and applies the local configuration to the specified cluster.
