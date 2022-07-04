# vab Modules

A `vab` module is a valid `kustomize` bundle of files, that can be used as-is or customized by applying
additional resources and patches.  
To allow `vab` to use it, the CLI expects a well-defined structure of its files. Every module must contain all
its files inside the same folder, in which you can have one or more sub-folders representing the various flavors
of the module. You can use the flavors to apply different configurations for specific cloud vendors, or as
alternative installations of the module.  
The sharing of files between modules is forbidden to avoid problems of circular dependencies. The only sharing permitted
is between flavors of the same module. For this reason, a valid installation can have only one flavor of a single module.

With the previous rules in mind, we envisioned the following folder structure inside the repository:

```txt
./modules
|   ├── module-1
|   |   ├── flavor-1
|   |   |..
|   |   └── flavor-n
|   └── module-2
|       ├── flavor-1
|       |..
|       └── flavor-n
└── README.md
```

## Versioning

The CLI will pull modules versioned via git tags inside the repository. However, since you can have multiple modules
inside a repository, the tags must be in the form of `module-<module-name>-<`[`semver`][semver]`-version>`.  
The CLI will match the `module-name` with the folder name included in the `modules` directory. and will pull all
the files contained in it, so you will have all the different flavors contained in it to ease the cross dependencies
between them.

[semver]: https://semver.org/spec/v2.0.0.html (semantic versioning v2.0.0 site)
