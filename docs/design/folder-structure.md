# Folder structure

At the first level inside the root directory, there will be:

- The configuration file `config.yaml`
- The **vendors** folder, containing the Kustomize bases for modules and add-ons
- The **clusters** folder, including the Kustomize patches divided by clusters and groups

## Vendors

The **vendors** folder will include two main subfolders: **modules** and **addons**.
The `module` foldder will contain the downloaded resources divided by flavors. Modules and add-ons will be structured
as Kustomize bases, hence each element will come with its own `kustomization.yaml` files.

For example:

```txt
./project-name
├── config.yaml
├── clusters
|   └── [...]
└── vendors
    ├── modules
    |   ├── ingress-1.0.0
    |   |   ├── traefik
    |   |   |   ├── [...]
    |   |   |   └── kustomization.yaml
    |   |   └── nginx
    |   |       ├── [...]
    |   |       └── kustomization.yaml
    |   └── cni-1.0.0
    |       ├── cilium
    |       |   ├── [...]
    |       |   └── kustomization.yaml
    |       └── calico
    |           ├── [...]
    |           └── kustomization.yaml
    └── add-ons
        └── ingress-traefik-monitoring-1.0.0
            ├── [...]
            └── kustomization.yaml
```

## Clusters

The structure of the folder `clusters` depends on the cluster configuration provided by the user.
Customizations will be grouped by cluster, which will be collected in their respective cluster group.
Consequently, at the first level, there will be as many folders as the number of cluster groups, each of which
will contain both the default (applied to all the clusters in the group, included in the `all-clusters` directory)
and the individual cluster customizations. Since Kustomize is managing the configurations, the overrides of the clusters
will be structured as Kustomize patches, hence including their own `kustomization.yaml`.
The `all-groups` directory will eventually contain customizations common to all groups.

Assuming that modules and add-ons will not have the same names, there will be no need to create a further nesting levels
to distinguish modules and add-ons.

For example:

```txt
./project-name
├── config.yaml
├── vendors
|   └── [...]
└── clusters
    ├── all-groups
    |   ├── [...]
    |   └── kustomization.yaml
    ├── group-1
    |   ├── all-clusters
    |   |   ├── traefik
    |   |   |   └── patch.yaml
    |   |   ├── cilium
    |   |   |   └── patch.yaml
    |   |   └── kustomization.yaml
    |   ├── cluster-1
    |   |   ├── traefik
    |   |   |   └── patch-1.yaml
    |   |   ├── calico
    |   |   |   └── patch-1.yaml
    |   |   └── kustomization.yaml
    |   └── cluster-2
    |       ├── cilium
    |       |   └── patch-3.yaml
    |       └── kustomization.yaml
    └── group-2
        └── [...]
```

It is critical for the `kustomization.yaml` files to be well-formed and consistent with the directory structure.
