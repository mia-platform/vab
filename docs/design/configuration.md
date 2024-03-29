# Configuration

The configuration file will include all the information needed to set up the clusters, from modules/add-ons to download,
to eventual customizations desired for the various contexts.

```yaml
apiVersion: vab.mia-platform.eu/v1alpha1
kind: ClustersConfiguration
name: my-clusters
spec:
  modules:  # type: Object
    ingress/traefik/base:
      version: 1.20.1
    cni/cilium/base:
      version: 1.20.1
  addons:   # type: Object
    monitoring/traefik:
      version: 1.20.1
  groups:   # type: Array[]
    - name: group-1
      clusters:
        - name: cluster-1
          context: context-1
          addons:
            monitoring/traefik:
              version: 1.20.100
        - name: cluster-2
          context: context-2
          modules:
            cni/cilium/base:
              disable: true
            cni/calico/base:
              version: 1.20.20
        - name: cluster-3
          context: context-3
```

In the sample configuration file above:

- A `ClusterConfiguration` named `my-clusters` is defined.
- The `modules` field is a dictionary that will include the modules to install by default on every cluster unless
  otherwise specified. In this case, the configuration will download the modules `ingress/traefik` and `cni/cilium`
  with version `1.20.1`.  
  The `version` of the core modules will follow the release schedule and version of Kubernetes for majors and minors,
  while patches will be released asynchronously.
- The `addons` field is a dictionary that will include the add-ons to install by default on every cluster unless
  otherwise specified. In this case, the configuration will download the add-on `monitoring/traefik`
  with version `1.20.1`.
- The `groups` field is an array that will list all the cluster groups to which the default configuration
  will be applied. Each group will contain a list of clusters with their customizations.
  In this case, we have a cluster group named `group-1` that will include:
  - A cluster named `cluster-1` that will use `context-1` for the connection and that overrides the add-on
    `ingress-monitoring` with a different version (`1.20.100`).  
    This directive will download the new version in the corresponding vendor folder.
  - A cluster named `cluster-2` that will use `context-2` for the connection, disables the `cni/cilium` module
    and installs the `cni/calico` module.  
    The latter directive will download the `cni/calico` module in the corresponding vendor folder.
  - A cluster named `cluster-3` that will use `context-3` for the connection, without any customization.  
    Therefore, `cluster-3` will be configured with all the modules and add-ons specified by default.

The `sync` command will be in charge of updating the vendors to the latest configuration and creating the appropriate
directory structure. According to the example above, `clusters/group-1` will include the following directories:

- **`all-clusters`:** containing patches of the modules (`ingress/traefik v1.20.1`, `cni/cilium v1.20.1`)
  and add-ons (`monitoring/traefik v1.20.1`) that will be applied to all the clusters;
- **`cluster-1`:** containing patches of the modules (`ingress/traefik v1.20.1`, `cni/cilium v1.20.1`)
  and add-ons (`monitoring/traefik v1.20.100`) that will be applied to Cluster 1;
- **`cluster-2`:** containing patches of the modules (`ingress/traefik v1.20.1`, `cni/calico v1.20.20`)
  and add-ons (`monitoring/traefik v1.20.1`) that will be applied to Cluster 2;
- **`cluster-3`:** containing patches of the modules (`ingress/traefik v1.20.1`, `cni/cilium v1.20.1`)
  and add-ons (`monitoring/traefik v1.20.1`) that will be applied to Cluster 3.

Assuming that the folder names will be consistent with those specified in the configuration,
there will be no need of referencing them in the configuration file.
