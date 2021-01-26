# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

- Update controller-runtime version

## [4.2.0] - 2021-01-26

- Pass SkipCRDs to Helm client.

## [4.1.0] - 2020-12-14

### Changed

- Updated Helm to v3.4.2.

## [4.0.0] - 2020-12-11

### Changed

- Get `k8sClient`, `restClient` and `restConfig` from generator.

### Removed

- Delete `MergeValue` function.

## [3.0.1] - 2020-10-29

### Added

- Validate the cache to prevent pulling stale cache objects.

### Fixed

- Add replace for moby v20.10.0-beta1 to fix build issue on darwin.

## [3.0.0] - 2020-10-27

### Changed

- Updated k8sclient to v5.
- Prepare module v3.

## [2.1.4] - 2020-10-01

### Added

- Added release name as a label into the event count metrics.

## [2.1.3] - 2020-09-24

### Security

- Updated Helm to v3.3.4.
- Updated Kubernetes dependencies to v1.18.9.

## [2.1.2] - 2020-09-22

### Added

- Added event count metrics for delete, install, rollback and update of Helm releases.

### Changed

- Fix structs merging error.


## [2.1.1] - 2020-09-21

### Security

- Updated Helm to v3.3.3.

## [2.1.0] - 2020-08-17

### Changed

- Updated Helm to v3.3.0.

## [2.0.0] - 2020-08-10

### Changed

- Updated Kubernetes dependencies to v1.18.5.
- Updated Helm to v3.2.4.
- Disable OpenAPI validation as some charts we need to deploy will contain
validation errors.

## [1.0.6] - 2020-08-05

### Added

- Add rollback support.

## [1.0.5] - 2020-07-17

### Added

- Add timeouts for installing or upgrading helm releases.

## [1.0.4] - 2020-07-13

### Changed

- Upgrade k8sclient to v3.1.1.

## [1.0.3] 2020-06-15

### Changed

- Add new yaml parsing errors.

## [1.0.2] 2020-05-29

- Add new validation error for both install and upgrade cases.

## [1.0.1] 2020-05-26

### Added

- Added manifest validation error.

### Security

- Updated Helm to [v3.1.3](https://github.com/helm/helm/releases/tag/v3.1.3) for security fix.

## [1.0.0] 2020-05-18

### Changed

- Updated to support Helm 3; To keep using Helm 2, please use version 0.2.X.

## [0.2.2] 2020-04-09

### Changed

- Add Helm revision number to GetReleaseHistory.

## [0.2.1] 2020-04-08

### Changed

- Handle 503 errors when pulling chart tarballs fails.
- Make HTTP client timeout for pulling chart tarballs configurable.

## [0.2.0] 2020-03-25

### Changed

- Switch from dep to go modules.
- Use architect orb.

## [0.1.0] 2020-03-19

### Added

- First release.

[Unreleased]: https://github.com/giantswarm/helmclient/compare/v4.2.0...HEAD
[4.2.0]: https://github.com/giantswarm/helmclient/compare/v4.1.0...v4.2.0
[4.1.0]: https://github.com/giantswarm/helmclient/compare/v4.0.0...v4.1.0
[4.0.0]: https://github.com/giantswarm/helmclient/compare/v3.0.1...v4.0.0
[3.0.1]: https://github.com/giantswarm/helmclient/compare/v3.0.0...v3.0.1
[3.0.0]: https://github.com/giantswarm/helmclient/compare/v2.1.4...v3.0.0
[2.1.4]: https://github.com/giantswarm/helmclient/compare/v2.1.3...v2.1.4
[2.1.3]: https://github.com/giantswarm/helmclient/compare/v2.1.2...v2.1.3
[2.1.2]: https://github.com/giantswarm/helmclient/compare/v2.1.1...v2.1.2
[2.1.1]: https://github.com/giantswarm/helmclient/compare/v2.1.0...v2.1.1
[2.1.0]: https://github.com/giantswarm/helmclient/compare/v2.0.0...v2.1.0
[2.0.0]: https://github.com/giantswarm/helmclient/compare/v1.0.6...v2.0.0
[1.0.6]: https://github.com/giantswarm/helmclient/compare/v1.0.5...v1.0.6
[1.0.5]: https://github.com/giantswarm/helmclient/compare/v1.0.4...v1.0.5
[1.0.4]: https://github.com/giantswarm/helmclient/compare/v1.0.3...v1.0.4
[1.0.3]: https://github.com/giantswarm/helmclient/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/giantswarm/helmclient/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/giantswarm/helmclient/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/giantswarm/helmclient/compare/v0.2.2...v1.0.0
[0.2.2]: https://github.com/giantswarm/helmclient/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/giantswarm/helmclient/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/helmclient/compare/v0.1.0...v0.2.0

[0.1.0]: https://github.com/giantswarm/helmclient/releases/tag/v0.1.0
