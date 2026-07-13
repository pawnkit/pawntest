# Go API

Use `pkg/pawntest` to embed the runner:

```go
runner := pawntest.NewRunner()
runner.PawnCC = "./tools/pawncc"
runner.Include = []string{"include"}

suite, err := runner.RunFileContext(ctx, "tests/example.test.pwn")
```

`NewRunner` and the zero value use a one-million-instruction limit. Set
`MaxInstructions` to a negative value to disable the limit. Context cancellation
stops discovery and compilation; AMX execution is bounded by the instruction
limit.
