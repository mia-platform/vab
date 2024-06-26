# vab CLI

`vab` is a Command Line Interface that simplify the management of the Magellano kubernetes distribution on one or
more remote clusters.

In future it will also allow to download custom modules and addons that follow the design documentations.

In essence vab can be viewed as an utility for downloading and setting up [kustomization] files with a structure
that allow sharing adding and patching resources that will be applied on different Kubernetes clusters.  
The structure that is created and mantained by the tools can alse be used with different tools that support `kustomize`
files like [Argo CD].

Applying the configuration files with `vab` will bring additional capabilities than using `kustomize` directly like
finer resource ordering, applied resoruce tracking for pruning resources not available in subsequent apply,
a little bit of validation for missing resource types in the target cluster and chaining deployments against
multiple cluster grouped together with a single command.

## Functionalities

The `vab` CLI functionalities can be summarized within its main subcommands:

- `apply`: apply all the manifests to one or more targeted cluster specified in the configuration file
- `build`: print all the manifests that the `apply` command would eventually apply to the cluster(s)
- `create`: create and empty configuration file and starting files structures in the target folder
- `sync`: donwload the modules and addons of the distribution locally and update the file structure if needed
- `validate`: validate the configuration file to check its validity or attention points

## Guides

Below, you can find additional documentaion for `vab`:

- [Setup](./20_setup.md)
- [Configuration Grammar](./30_grammar.md)
- Design Documentation
  - [Configuration](./design/10_configuration.md)
  - [Folder Structures](./design/20_folder-structure.md)
  - [Modules](./design/30_modules.md)
  - [Addons](./design/40_addons.md)
  - [Download Packages](./design/50_download-packages.md)

[kustomization]: https://kustomize.io "Kubernetes native configuration management"
[Argo CD]: https://argo-cd.readthedocs.io/en/stable/ "Declarative GitOps CD for Kubernetes"
