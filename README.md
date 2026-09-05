# chimney3-go

This repository follows a conventional Go layout, with the executable entrypoint in `cmd/chimney`.

`chimney3-go` 是一个使用 Go 编写的网络代理服务，提供 SOCKS5、代理转发和 KCP 三种服务入口，并包含面向 Android VPN 集成的 `vpncore` 模块。

项目入口位于 `cmd/chimney`。程序从 JSON 文件读取配置，在启动时规范化并校验地址、服务类型、运行模式和加密方法，然后启动对应的服务运行器。

## 功能

- SOCKS5 服务
- 通用代理转发服务
- KCP 服务
- 用户名和密码认证配置
- 可配置的加密方法
- TCP、UDP 及 HTTP 相关监听地址
- Android `gomobile` 绑定支持
- 收到 `SIGINT` 或 `SIGTERM` 后优雅退出

## 快速开始

要求：Go `1.26` 或更高版本。

在项目根目录构建：

```bash
go build -o bin/chimney ./cmd/chimney
```

直接运行时，程序默认读取可执行文件目录下的 `configs/setting.json`：

```bash
./bin/chimney
```

也可以通过 `-config` 指定配置文件：

```bash
./bin/chimney -config /path/to/setting.json
```

运行测试：

```bash
go test ./...
```

## 配置

默认配置文件为 `configs/setting.json`。配置使用扁平 JSON 结构：

```json
{
  "listen": "[::1]:9025",
  "remote_listen": "127.0.0.1:9090",
  "username": "where",
  "password": "change-me",
  "method": "CHACHA-20",
  "which": "socks5",
  "udplisten": "0.0.0.0:15999",
  "httpurl": "0.0.0.0:4010",
  "mode": "server"
}
```

字段说明：

| 字段 | 说明 |
| --- | --- |
| `listen` | 本地服务监听地址，支持 `host:port` 或仅端口 |
| `remote_listen` | 远端服务地址 |
| `username` | 用户名 |
| `password` | 密码 |
| `method` | 加密方法，名称会被规范化为大写 |
| `which` | 服务类型：`socks5`、`proxy` 或 `kcp` |
| `udplisten` | UDP 监听地址 |
| `httpurl` | HTTP 相关监听地址 |
| `mode` | 运行模式：`server` 或 `client` |

`which` 和 `mode` 不区分大小写。地址字段可以使用 IPv4、IPv6 或仅端口形式；程序会在启动前拒绝无效地址、未知服务、未知模式和不支持的加密方法。

## 项目结构

```text
cmd/chimney/   命令行程序入口
configs/       默认 JSON 配置
all/           服务选择与运行分发
socks5/        SOCKS5 服务
proxy/         代理转发服务
kcpproxy/      KCP 服务
privacy/       加密方法实现
core/          网络连接与数据转发基础组件
vpncore/       Android VPN / gomobile 集成
utils/         通用工具
```

## Android 构建

安装 `gomobile` 后，可以将 `vpncore` 绑定为 Android AAR：

```bash
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init
gomobile bind -v -o bin/vpn.aar -target=android -androidapi 30 -ldflags '-w -s' ./vpncore
```

生成的 AAR 文件位于 `bin/vpn.aar`。Android 集成需要由调用方提供 VPN 文件描述符、TCP/UDP 代理地址以及 socket 保护函数，具体接口定义见 `vpncore` 包。

## 许可证

本项目使用仓库中的 [LICENSE](LICENSE) 文件所规定的许可证。

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
