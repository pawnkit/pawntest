# Examples

Each directory is a standalone pawntest project.

| Example | Covers |
|---|---|
| [`basic`](basic) | Testing stock functions and assertions |
| [`gamemode`](gamemode) | Player lifecycle, commands, and scenario state |
| [`mocks`](mocks) | Mocking plugin natives and checking calls |
| [`provider`](provider) | Supplying plugin-native implementations in Pawn |
| [`roleplay`](roleplay) | Complete gamemode project using the full workflow |

From an example directory, run:

```sh
pawntest
```

Pawntest reads the local `pawntest.toml`. `pawncc` must be on `PATH` or passed
with `--pawncc`.
