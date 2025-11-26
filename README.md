# chimney3-go 

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


 server and client are the same program. 
```for server
  "mode": "server"


  for client:
  "mode": "client"
```
