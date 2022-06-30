# Commands

## `init <project-name>`

Creates the folder `project-name` with a preliminary directory structure:

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

And the skeleton of the `config.yaml` file:

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

Validates the configuration

## `sync`

Fetch new vendors and update the configuration after editing the `config.yaml` file.

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

Results in the project root:

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

## `apply <group-name> <cluster-name>`

Builds and applies the local configuration to the specified cluster.

