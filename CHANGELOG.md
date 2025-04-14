# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- [#60] Set sensible resource requests and limits

## [v1.3.1] - 2025-03-06
### Fixed
- [#45] Fixed sporadically failing test
- [#57] Fix missing RBAC role to list blueprints

## [v1.3.0] - 2025-01-31
### Changed
- [#55] Use Dogu-CR & DoguRestart-CR to administer dogus instead of deployments
- [#55] Update go version to 1.23.2

## [v1.2.1] - 2024-12-19
### Fixed
- [#53] Fix critical CVEs by updating grpc_health_probe to 0.4.36

## [v1.2.0] - 2024-12-10
### Added
- [#51] NetworkPolicy to block all ingress traffic
  - Dependent Dogus and Components must bring their own NetworkPolicy to access k8s-ces-control
- [#51] Ingress Network Policy for Loki-gateway so that k8s-ces-control can access it

## [v1.1.0] - 2024-11-25
### Fixed
- [#48] Fix not applying development values to helm values

### Changed
- [#47] drastically reduce RBAC permissions

## [v1.0.0] - 2024-10-30
### Changed
- [#43] Make imagePullSecrets configurable via helm values and use `ces-container-registries` as default.
- [#41] Use dogu v2 api
- [#41] Update go version to 1.23.2

## [v0.11.1] - 2024-10-07
### Fixed
- [#39] Fix helm dependency version constraint for loki and allow loki versions >=2.9.1-0 and <4.0.0-0

## [v0.11.0] - 2024-09-20
### Changed
- [#35] - Use k8s-registry-lib to read and write log level

## [v0.10.0] - 2024-09-18
### Changed
- [#37] - Relicense to AGPL-3.0-only

## [v0.9.0] - 2024-07-15

### Changed
- [#33] - Use configmap based dogu registry instead of etcd to retrieve dogu.jsons

## [v0.8.1] - 2024-05-24
### Fixed
- [#31] After setting a new log level for a dogu, the dogu restarts even when it is stopped - dogu should stay stopped

## [v0.8.0] - 2024-05-23
### Changed
- Update Makefiles to 9.0.5

### Added
- [#27] API GetDoguList reports current LogLevel
- [#25] API to set log level for a dogu
  - when the provided log level differs from the log level set in the config, a dogu restart is triggered to apply the change

### Fixed
- A bug where the healthcheck never returned the correct container state (#29) 

## [v0.7.0] - 2024-05-02
### Added
- [#23] Add new query API for Dogu-Logs
  - The GRPC-API-definition is now located at https://github.com/cloudogu/ces-control-api


## [v0.6.0] - 2024-03-27
### Added
- GetBlueprintId endpoint added (#21)
    - retrieves the blueprint id of the currently installed blueprint, if applicable.

## [v0.5.0] - 2023-12-11
### Added
- [#17] Provide logs from loki
- [#18] Patch-templates for mirroring into airgapped environments
### Changed
- [#18] Extract yaml wallpaper into helm templates folder

## [v0.4.0] - 2023-11-14
### Added
- [#15] Add first version of debug mode for dogus without data collection and log rotation.

## [v0.3.0] - 2023-09-15
### Changed
- [#13] Move component-dependencies to helm-annotations

## [v0.2.0] - 2023-09-05
### Added
- [#9] Add API-endpoints for start, stop & restart dogus
- [#11] Add API-endpoints for dogu-health

## [v0.1.1] - 2023-08-31
### Added
- [#7] Add "k8s-etcd" as a dependency to the helm-chart

## [v0.1.0] - 2023-08-14
### Added
- [#5] Initialize a first version for the `k8s-ces-control`. In contrast to the prior poc status k8s-ces-control does not use TLS or service account verification, because the current Admin-Dogu does not support this.
- [#5] Add Helm chart release process to project