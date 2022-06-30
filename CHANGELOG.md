# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project's packages adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

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

[Unreleased]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.21.0-gs1...HEAD
[1.21.0-gs1]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v1.21.0-gs0...v1.21.0-gs1
[1.21.0-gs0]: https://github.com/giantswarm/aws-cloud-controller-manager-app/compare/v0.0.0...v1.21.0-gs0
