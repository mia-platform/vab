# vab

`vab` is a cli for managing the installation of day 2 operation tools on multiple clusters for easier management and
updates.

## Building

`vab` provides various make command to handle various tasks that you may need during development, but you need at
least these dependencies installed on your machine:

- make
- bash
- docker with buildkit build engine available to use
- golang, for the exact version see the [.go-version](/.go-version) file in the repository

Once you have all the correct dependencies installed and the code pulled you can simply build the project with:

```bash
make build
```

This command will build the cli for your actual OS and architecture and will put the binary inside the folder
`bin/<os>/<arch>/`.

Or run the tests with:

```bash
make test
```

If you don’t plan to build a docker image but only to contribute to the code, we provide a devcontainer configuration
that will setup the correct dependencies and predownload the tools used for linting. Also if you use VSCode it will
setup three extensions that we recommend.

### Linting

For linting your files make provide the following command:

```bash
make lint
```

This command will run `go mod tidy` for cleaning up the `go.mod` and `go.sum` files and will stop if it senses that
the files are changed and not already commited or added to the git staging area, this check is done forcing the user
to not forgetting this steps and for breaking the ci/cd on GitHub if those files are not in the correct shape.  
Additionally the command will download and use the [`golangci-lint`][golangci-lint] cli for running various linters
on the code, the configuration used can be seen [here](/tools/.golangci.yml).

### Building Docker Image

If you need to use a docker image locally you can build it with:

```bash
make build-image
```

The command will first build the appropriate binary for your architecture and then build the correct docker image for
your platform based on Linux Alpine.

### Building Multiarch Docker Image

If you need to try and build a multiarch docker image locally:

```bash
make build-image-multiarch
```

For building this image you need to have installed `docker` and the `buildx` extension for emulating multiple
architecture on your pc. This command for now is created for using it in the ci.

[golangci-lint]: https://golangci-lint.run (Fast linters Runner for Go)
