# vab

`vab` is a cli for managing the installation of day 2 operation tools on multiple kubernetes clusters for easier
management and updates.

`vab` is the acronym for Vehicle Assembly Building that is designed to assemble large pre-manufactured
space vehicle components.

## Building

`vab` provides various make command to handle various tasks that you may need during development, but you need at
least these dependencies installed on your machine:

- make
- bash
- docker with buildkit build engine available to use
- golang, for the exact version see the [.go-version](/.go-version) file in the repository

Once you have all the correct dependencies installed and the code pulled you can build the project with:

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

This command will run `go mod tidy` for cleaning up the `go.mod` and `go.sum` files.  
Additionally the command will download and use the [`golangci-lint`][golangci-lint] cli for running various linters
on the code, the configuration used can be seen [here](.golangci.yml).

### Building Docker Image

If you need to use a docker image locally you can build it with:

```bash
make docker-build
```

The command will first build the appropriate binary and then build the correct docker image for
your platform based on Linux Alpine.

### Building Multiarch Docker Image

If you need to try and build a multiarch docker image locally, you have to run these commands:

```bash
make docker-setup-multiarch
make docker-build-multiarch
```

For building this image you need to have installed `docker` and the `buildx` extension for emulating multiple
architecture on your pc. This command for now is created for using it in the ci.

[golangci-lint]: https://golangci-lint.run (Fast linters Runner for Go)
