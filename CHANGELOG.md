# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

- update go version to 1.22.4

## [v0.9.0] - 2023-02-08

### Added

- init command: initialize a new folder with base folders and config files
- sync command: download modules data and configura base kustomize files
- build command: print on stdout the generated files from all modules, addons and customizations
- apply command: apply generated files to targets clusters
- validate command: validate the configuration file

[Unreleased]: https://github.com/mia-platform/jpl/compare/v0.9.0...HEAD
[v0.9.0]: https://github.com/mia-platform/vab/releases/tag/v0.1.0
