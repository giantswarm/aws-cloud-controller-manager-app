# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project's packages adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.25.14-gs3] - 2024-05-09

### Changed

- Reduce minimum CPU and memory requests.

## [1.25.14-gs2] - 2024-01-22

### Changed

- Configure `gsoci.azurecr.io` as the default container image registry.

## [1.25.14-gs1] - 2024-01-16

### Changed

- Bump to upstream version 1.25.14.

## [1.24.1-gs10] - 2023-10-18

### Fixed

- Add required values for pss policies.

### Added

- Add toggle mechanism for `PSPs`.

## [1.24.1-gs9] - 2023-06-30

### Changed

- Adjusted minimum allowed CPU and memory

## [1.24.1-gs8] - 2023-06-12

### Changed

- Always install the VPA CR if `verticalPodAutoscaler.enabled` is true, no matter if the VPA CRD is present or not.

## [1.24.1-gs7] - 2023-05-11

### Fixed

- Quote environment variables that contain numeric values, because it's required by kubernetes.

## [1.24.1-gs6] - 2023-05-10

### Added

- Added two new values to set `KUBERNETES_SERVICE_HOST` and `KUBERNETES_SERVICE_PORT`. This is needed because sometimes we install the app before the CNI is installed, and the controller can't talk to the k8s api using the k8s `Service`.

## [1.24.1-gs5] - 2023-05-10

### Fixed

- Fix label selector in `NetworkPolicy`.

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

[Unreleased]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.25.14-gs3...HEAD
[1.25.14-gs3]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.25.14-gs2...v1.25.14-gs3
[1.25.14-gs2]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.25.14-gs1...v1.25.14-gs2
[1.25.14-gs1]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.24.1-gs10...v1.25.14-gs1
[1.24.1-gs10]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.24.1-gs9...v1.24.1-gs10
[1.24.1-gs9]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.24.1-gs8...v1.24.1-gs9
[1.24.1-gs8]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.24.1-gs7...v1.24.1-gs8
[1.24.1-gs7]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.24.1-gs6...v1.24.1-gs7
[1.24.1-gs6]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.24.1-gs5...v1.24.1-gs6
[1.24.1-gs5]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.24.1-gs4...v1.24.1-gs5
[1.24.1-gs4]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.24.1-gs3...v1.24.1-gs4
[1.24.1-gs3]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.24.1-gs2...v1.24.1-gs3
[1.24.1-gs2]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.24.1-gs1...v1.24.1-gs2
[1.24.1-gs1]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.23.2-gs2...v1.24.1-gs1
[1.23.2-gs2]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.23.2-gs1...v1.23.2-gs2
[1.23.2-gs1]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.22.4-gs1...v1.23.2-gs1
[1.22.4-gs1]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.21.0-gs1...v1.22.4-gs1
[1.21.0-gs1]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.21.0-gs0...v1.21.0-gs1
[1.21.0-gs0]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v0.0.0...v1.21.0-gs0
