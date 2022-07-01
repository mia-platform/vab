# vab Add-ons

A `vab` add-on is a valid `kustomize` bundle of files that can be used as-is or customized via additional resources
and patches applied to it.  
To allow `vab` to use it, the CLI expects a well-defined structure of its files. Every add-on must contain all
its files inside the same folder.  
The sharing of files between add-ons is strictly forbidden, as they must be self-contained and not dependent on any
module. For this reason, an add-on is strictly additive to a target module. You can use an add-on to add resources
to a module that will be deployed from another one (like monitoring via the Prometheus operator's `crd`s),
or to add further functionalities with additional resources (like adding additional frontends for visualizing
data from a module, or Grafana to a monitoring stack).  
Forbidding dependency between add-ons is required to increase compatibility between add-ons as much as possible.

With the previous rules in mind, we envisioned the following folder structure inside the repository:

```txt
./add-ons
|   ├── add-on-1
|   ├── add-on-2
|   |..
|   └── add-on-n
└── README.md
```

## Versioning

The CLI will pull add-ons that are versioned via git tags inside the repository, but because inside a repository you
can have multiple add-ons the tags must be in the form of `add-on-<add-on-name>-<semver-version>`.  
The CLI will match the `add-on-name` with the folder name contained inside the `add-ons` one. and will pull all
the files contained in it.
