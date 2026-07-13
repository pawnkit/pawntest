# Limitations

- Pawntest runs 32-bit AMX files with [`pawnkit/goamx`](https://github.com/pawnkit/goamx).
- 16-bit and 64-bit AMX files are not supported.
- Unmodeled server and plugin natives require mocks or custom Go natives.
- Automatic compiler installation supports compatible openmultiplayer releases only.
- Linux and Windows compiler downloads are not available on ARM.
- Coverage instruments the test AMX, not provider AMXs.
- Provider callback arguments are limited to eight cells or floats.

To test source compilation locally, set `PAWNTEST_PAWNCC` and run `go test`.
