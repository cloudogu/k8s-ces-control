# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- [#7] Add API-endpoints for start, stop & restart dogus

## [v0.1.1] - 2023-08-31
### Added
- [#7] Add "k8s-etcd" as a dependency to the helm-chart

## [v0.1.0] - 2023-08-14
### Added
- [#5] Initialize a first version for the `k8s-ces-control`. In contrast to the prior poc status k8s-ces-control does not use TLS or service account verification, because the current Admin-Dogu does not support this.
- [#5] Add Helm chart release process to project