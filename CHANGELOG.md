# Changelog

Notable changes are documented in GitHub Releases. The project follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and semantic
versioning after `v1.0.0`.

## 1.2.9 - 2026-08-02

### Fixed

- Authenticate GitHub release metadata requests when a workflow token is
  available, avoiding unauthenticated API rate limits.

## 1.2.8 - 2026-08-02

### Fixed

- Treat Windows access-denied lock opens as concurrent cache contention.

## 1.2.7 - 2026-08-02

### Changed

- Use pawn-project 0.34.2.

## 1.2.6 - 2026-08-02

### Changed

- Use pawn-project 0.34.1 and pawnkit-core 0.5.0.

## 1.2.5 - 2026-07-29

### Added

- Added the test adapter used by `pawn check`.

## 1.2.4 - 2026-07-29

### Fixed

- Updated the plugin host to reject incomplete worker frames.

## 1.2.3 - 2026-07-29

### Fixed

- Updated the runtime and plugin host to fix AMX switch-table execution.

## 1.2.2 - 2026-07-25

### Changed

- Require differential fixtures to name their reference engine and runtime
  tier.

## 1.2.1 - 2026-07-25

### Fixed

- Updated the support record to match the runtime metadata now emitted.

## 1.2.0 - 2026-07-25

### Added

- Report runtime-fidelity metadata in JSON, TAP, JUnit, and plain output.
- Declare and validate the repository support policy in CI.

## 1.1.4 - 2026-07-23

### Changed

- Updated to the current Pawn project release.

## 1.1.3 - 2026-07-23

### Fixed

- Updated runtime, plugin-host, and project dependencies.

## 1.1.2 - 2026-07-21

### Added

- Added source file and line metadata to JSON test discovery.

## 1.1.1 - 2026-07-20

### Fixed

- Included pawntest headers in release archives for editor tooling.

## 1.1.0 - 2026-07-20

### Added

- Project manifest configuration and instruction profiles.
- Opt-in legacy plugin natives through `pawn-plugin-host`.
- JSON test discovery for editor integrations.

## 1.0.1 - 2026-07-13

### Added

- Runnable examples for common pawntest workflows.

### Fixed

- Test filtering now applies consistently.

## 1.0.0 - 2026-07-13

### Added

- Pawn test runner, scenario models, mocks, fixtures, reports, coverage, and snapshots.
- Versioned Pawn native-provider ABI.
- Compiler discovery, verified installation, caching, diagnostics, and doctor checks.

### Security

- Release checksums, SBOMs, and build provenance.
