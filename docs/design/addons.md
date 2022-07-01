# vab Add-ons

A `vab` add-on is a valid `kustomize` bundle of files that can be used as is or customized via additional resources
and patches applied to it.  
To allow `vab` to use it the cli is expecting a well defined structure of its files. Every add-on must contain all
its files inside the same folder.  
Sharing of files between add-on are strictly forbidden and they must be self contained and not dependent from any
module. An add-on for this reason is stricly additive for a target module, you can use add-on for adding resources
to a module that will be deployed from another one (like monitoring via the prometheus-operator `crd`s),
or for adding additional functionality with additional resources (like adding additional frontends for visualizing
data from a module, like adding grafana to a monitoring stack).  
Forbidding dependency between add-on is a required condition for being able to increase at the maximum possibile
compatibily between different add-ons.

With the previous rules in mind we envisioned the following folder structure inside the repository:

```txt
./add-ons
|   ├── add-on-1
|   ├── add-on-2
|   |..
|   └── add-on-n
└── README.md
```

## Versioning

The cli will pull add-ons that are versioned via git tags inside the repository, but because inside a repository you
can have multiple add-ons the tags must be in the form of `add-on-<add-on-name>-<semver-version>`.  
The cli will match the `add-on-name` with the folder name contained inside the `add-ons` one. and will pull all
the files contained in it.
