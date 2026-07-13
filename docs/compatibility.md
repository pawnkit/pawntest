# Compatibility

Pawntest uses semantic versioning after `v1.0.0`. Before `v1.0.0`, minor
releases may change the Go API, configuration schema, Pawn macros, provider ABI,
or report schemas. Changes are listed in GitHub release notes.

Stable compatibility surfaces:

- CLI flags and exit codes
- Configuration fields
- Pawn test macros and natives
- Provider ABI version
- JSON, JUnit, TAP, and coverage schemas
- `pkg/pawntest` public API

Provider source is compiled against the installed pawntest include. Precompiled
provider AMXs must use the same provider ABI version.
