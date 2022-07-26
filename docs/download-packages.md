# Download Packages

As stated in the documentation relative to the [configuration specification] a module or an add-on is
identified with a **name** and a **version**.  
With this two pieces of information we will have to reference unequivocally a bundle of one or more file to download.

In the [modules] and [addons] documentation is stated that we expect a certain folder structure inside the repository
that is containing the modules and/or add-ons and how we expect the tags name for referencing the correct versions.

In this documentation will outline how all of this is tied together and what at high level the download will do
for retriving the correct files to use.

## High Level Design

The implementation is based on avoiding external dependencies, so we will use a library that implement
natively the git commands like [go-git] this will lessen the burden of setting up the machine with additional
dependencies and for avoiding incopatibility with different version of the `git` cli.

The implementation flow will be based on these steps:

- cleaning old files if presents
- downloading the files matching the name and versions of modules and add-ons
- copying them in the correct locations

If the first and last steps are straightforward to implement, the second point is the main protagonist of this doc.

Downloding a module or an add-on can be seen essentially as a clone operation targeted to a specific tag.  
The remote url is set as the url of the mia-platform monorepo containing all the modules and add-ons, and the
various tags will be built using the name and the the version contained in the configuration file.

Once the correct tag and url are generated we can use them to create a temporary clone and then using it for copying
all the files contained inside the correct folders (add-ons or module, remembering that all the module with all the
flavors will be copied for mantaining cross dependencies between them).

For the first version only the mia-platform offical public repo will be supported via the https conncetion
and so we don’t have to support particular connection credentials; but the interaction with git must be
incapsulated in a dedicated module for easier modification in future for supporting different repos, connection
and credentials.

## Open Points for Future Enhancement

For this first implementation we will have the following open points of possibile improvments:

1. parallelization of the downloads, for now the download will be done sequentially
1. clone the target repository only once and done the different checkout of the tags without cloning multiple time
  the same repository
1. support different respositories, for now all the modules and add-on must be in the mia-platform offical monorepo
1. support of credentials and different protocols for the connections, without the need to have sensitive data
  written inside the configuration file
1. cache the download for avoiding downloading the same data over when run the command multiple times without changes

[configuration specification]: design/configuration.md "vab configuration specifications"
[modules]: design/modules.md "modules specification"
[addons]: design/addons.md "add-ons specification"
[go-git]: https://github.com/go-git/go-git "go-git repository"
