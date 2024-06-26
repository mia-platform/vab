# vab

<center>

[![Build Status][github-actions-svg]][github-actions]
[![Go Report Card][go-report-card]][go-report-card-link]
[![GoDoc][godoc-svg]][godoc-link]

</center>

`vab` is a cli for managing the installation of day 2 operation tools on multiple kubernetes clusters for easier
management and updates.

`vab` is the acronym for Vehicle Assembly Building that is designed to assemble large pre-manufactured
space vehicle components.

## To Start Using `vab`

Read the documentation [here](./docs/10_overview.md).

## To Start Developing `vab`

To start developing the CLI you must have this requirements:

- golang 1.22
- make

Once you have pulled the code locally, you can build the code with make:

```sh
make build
```

`make` will download all the dependencies needed and will build the binary for your current system that you can find
in the `bin` folder.

To build the docker image locally run:

```sh
make docker-build
```

## Testing `vab`

To run the tests use the command:

```sh
make test
```

Or add the `DEBUG_TEST` flag to run the test with debug mode enabled:

```sh
make test DEBUG_TEST=1
```

Before sending a PR be sure that all the linter pass with success:

```sh
make lint
```

[github-actions]: https://github.com/mia-platform/vab/actions
[github-actions-svg]: https://github.com/mia-platform/vab/workflows/Continuous%20Integration%20Pipeline/badge.svg
[godoc-svg]: https://godoc.org/github.com/mia-platform/vab?status.svg
[godoc-link]: https://godoc.org/github.com/mia-platform/vab
[go-report-card]: https://goreportcard.com/badge/github.com/mia-platform/vab
[go-report-card-link]: https://goreportcard.com/report/github.com/mia-platform/vab
