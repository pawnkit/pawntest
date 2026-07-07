# Limitations

Pawntest uses the pure-Go `github.com/pawnkit/goamx` runtime for AMX execution.
That runtime targets 32-bit AMX files produced by the supported Pawn compiler.
It executes the compiler's active opcode set, compact images, standard
floating-point helpers, Pawntest natives, and configured mocks.

Deliberate support boundaries:

- 16-bit and 64-bit AMX cell formats are rejected.
- The initial scenario pack models common player identity, connection,
  position, money, messaging, broadcast, and kick natives. Other game-mode and
  plugin natives require mocks or custom Go-native modules.
- The canonical C adapter is a development-only parity oracle. Production
  execution is entirely Go-native through `pawnkit/goamx`.
- Automatic compiler installation is limited to compatible assets published by
  openmultiplayer. At present, generic Linux/Windows archives are treated as
  x86 binaries and are not offered on ARM.

Compiled Pawn fixture acceptance is available when `pawncc` is installed. The
test suite includes a local opt-in integration test; set `PAWNTEST_PAWNCC` to a
compiler path to exercise source compilation during `go test`. CI runs this
corpus with the official OpenMP compiler on Linux.
