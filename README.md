# chimney3-go

This repository follows a conventional Go layout, with the executable entrypoint in `cmd/chimney`.

Quick start:

1. Build the binary:

```bash
go build -o bin/chimney ./cmd/chimney
```

2. Run with the default config next to the executable, or pass a custom config path:

```bash
./bin/chimney
./bin/chimney -config /path/to/setting.json
```

The same program can run as either a server or a client, selected by the `mode` field in the JSON config:

```json
{
  "mode": "server"
}
```

```json
{
  "mode": "client"
}
```

Supported service selectors:

- `which: "socks5"`
- `which: "proxy"`
- `which: "kcp"`

The startup layer validates unsupported service names and invalid runtime modes early, and treats nil configuration as an explicit error instead of crashing.
