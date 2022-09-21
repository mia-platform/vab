# vab Add-ons

A `vab` add-on is a valid `kustomize Component` that can create new resources or patches to apply to its target
module(s).  
To allow `vab` to use it, the CLI expects a well-defined structure of its files. Every add-on must contain all
its files inside the same folder.  
The sharing of files between add-ons is strictly forbidden, as they must be self-contained and not dependent on other
add-on. You can use an add-on to add resources to a module using `crd`s from another module
(like monitoring via the Prometheus operator's `crd`s), or to add further functionalities with
additional resources (like adding additional frontends for visualizing
data from a add-on, or Grafana to a monitoring stack).  
Forbidding dependency between add-ons is required to increase compatibility between add-ons as much as possible.

With the previous rules in mind, we envisioned the following folder structure inside the repository:

```txt
./add-ons
|   ├── category-1
|   |   ├── add-on-1
|   |   |..
|   |   └── add-on-n
|   └── category-2
|   |   ├── add-on-1
|   |   |..
|   |   └── add-on-n
└── README.md
```

## Versioning

The CLI will pull add-ons that are versioned via git tags inside the repository, but because inside a repository you
can have multiple add-ons the tags must be in the form of
`add-on-<add-on-category>-<add-on-name>-<`[`semver`][semver]`-version>`.  
The CLI will match the `add-on-category/add-on-name` path inside the `add-ons` folder; and will pull all
the files contained in it.

[semver]: https://semver.org/spec/v2.0.0.html "semantic versioning v2.0.0 site"
