# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project's packages adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.24.1-gs4] - 2023-05-10

### Changed

- Change the `nodeSelector` in the `DaemonSet` from `node-role.kubernetes.io/master: ""` to `node-role.kubernetes.io/control-plane: ""`.

## [1.24.1-gs3] - 2023-05-02

### Added

- Support for running behind a proxy.
- Support for using `cluster-apps-operator` specific `cluster.proxy` values.

## [1.24.1-gs2] - 2023-05-01

### Fixed

- Fix indentation for `requests` and `limits` on the daemonset manifest.

## [1.24.1-gs1] - 2022-09-13

### Changed

- Bump to upstream version 1.24.1.

## [1.23.2-gs2] - 2022-07-20

### Changed

- Add to default catalog.

## [1.23.2-gs1] - 2022-07-13

### Changed

- Bump to upstream version 1.23.2.

## [1.22.4-gs1] - 2022-07-01

### Changed

- Bump to upstream version 1.22.4.
- Switch to default dnsPolicy to avoid circular dependency with coredns.
- Add readinessProbe to make rollouts smoother.

## [1.21.0-gs1] - 2022-04-20

### Added

- Add `config.giantswarm.io/version` annotation to Chart.yaml.

### Changed

- Push to `giantswarm` catalog instead of `control-plane` one.

## [1.21.0-gs0] - 2022-04-11

### Added

- Initial release.

[Unreleased]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.24.1-gs4...HEAD
[1.24.1-gs4]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.24.1-gs3...v1.24.1-gs4
[1.24.1-gs3]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.24.1-gs2...v1.24.1-gs3
[1.24.1-gs2]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.24.1-gs1...v1.24.1-gs2
[1.24.1-gs1]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.23.2-gs2...v1.24.1-gs1
[1.23.2-gs2]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.23.2-gs1...v1.23.2-gs2
[1.23.2-gs1]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.22.4-gs1...v1.23.2-gs1
[1.22.4-gs1]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.21.0-gs1...v1.22.4-gs1
[1.21.0-gs1]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.21.0-gs0...v1.21.0-gs1
[1.21.0-gs0]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v0.0.0...v1.21.0-gs0
