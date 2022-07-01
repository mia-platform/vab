# vab Modules

A `vab` module is a valid `kustomize` bundle of files that can be used as is or customized via additional resources
and patches applied to it.  
To allow `vab` to use it the cli is expecting a well defined structure of its files. Every module must be contain all
its files inside the same folder, in that folder you can have one or more subfolder that represents
flavours of the module. You can use the flavours for appling different configuration for specific cloud vendor or for
alternative installation of the module.  
Sharing files between modules is forbidden for avoiding problems of circular dependencies, the only sharing permitted
are between flavours of the module and for this reason a valid installation can have only one flavour of a module.

With the previous rules in mind we envisioned the following folder structure inside the repository:

```txt
./modules
|   ├── module-1
|   |   ├── flavour-1
|   |   |..
|   |   └── flavour-n
|   |   
|   └── module-2
|       ├── flavour-1
|       |..
|       └── flavour-n
└── README.md
```

## Versioning

The cli will pull modules that are versioned via git tags inside the repository, but because inside a repository you
can have multiple modules the tags must be in the form of `module-<module-name>-<semver-version>`.  
The cli will match the `module-name` with the folder name contained inside the `modules` one. and will pull all
the files contained in it, so you will have all the different flavour contained in it to ease the cross dependencies
between them.
