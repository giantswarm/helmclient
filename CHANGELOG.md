# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/giantswarm/helmclient/compare/v1.0.2...HEAD
[1.0.2]: https://github.com/giantswarm/helmclient/compare/v1.0.1..v1.0.2
[1.0.1]: https://github.com/giantswarm/helmclient/compare/v1.0.0..v1.0.1
[1.0.0]: https://github.com/giantswarm/helmclient/compare/v0.2.2..v1.0.0
[0.2.2]: https://github.com/giantswarm/helmclient/compare/v0.2.1..v0.2.2
[0.2.1]: https://github.com/giantswarm/helmclient/compare/v0.2.0..v0.2.1
[0.2.0]: https://github.com/giantswarm/helmclient/compare/v0.1.0..v0.2.0

[0.1.0]: https://github.com/giantswarm/helmclient/releases/tag/v0.1.0
