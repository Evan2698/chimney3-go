# chimney3-go (refactored)

This repository was refactored to follow a conventional Go layout. The executable entrypoint now lives in `cmd/chimney`.

Quick start:

1. Build the binary:

```bash
go build -o bin/chimney ./cmd/chimney
```

2. Run with the config next to the executable (or pass path in code):

```bash
# run as server
./bin/chimney
```



Other packages and source files have not been reorganized beyond this entrypoint move.

If you'd like, I can continue reorganizing packages into `internal/` and add tests.
