# Limitations

- Pawntest runs 32-bit AMX files with [`pawnkit/goamx`](https://github.com/pawnkit/goamx).
- 16-bit and 64-bit AMX files are not supported.
- Unmodeled server and plugin natives require mocks or custom Go natives.
- Automatic compiler installation supports compatible openmultiplayer releases only.
- Linux and Windows compiler downloads are not available on ARM.
- Coverage instruments the test AMX, not provider AMXs.
- Provider callback arguments are limited to eight cells or floats.
- Legacy plugins run only on Linux and require a separately installed
  `pawn-plugin-host` worker and loader.
- The legacy bridge supports one native with up to 32 cell arguments. It does
  not yet support strings, references, arrays, callbacks, or worker reuse.

To test source compilation locally, set `PAWNTEST_PAWNCC` and run `go test`.
