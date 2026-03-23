# 环境搭建

本文档用于搭建 `quic-raft` 的本地实验环境，目标是先跑通最小 3 节点实验，再逐步引入 QUIC、故障注入和性能观测。

## 1. 当前环境检查结果

本次初始化时，本机环境检测结果如下：

- 操作系统：macOS `arm64`
- Go：已安装，版本为 `go1.26.1`
- Docker：已安装
- protoc：已安装

可重复执行：

```bash
make doctor
```

## 2. 推荐依赖

可选：

- `tc` / `netem` 或 `toxiproxy`
- `hey`、`wrk` 或自定义压测脚本
- Prometheus / Grafana

## 3. macOS 安装建议

如果你使用 Homebrew，可以按下面的顺序准备环境：

```bash
brew install go
brew install protobuf
brew install make
brew install --cask docker
```

安装后验证：

```bash
go version
docker --version
docker compose version
protoc --version
make --version
```

## 4. 项目初始化

进入项目目录：

```bash
cd /Users/leo/Code/Local/quic_experiment
```

执行基础检查：

```bash
make doctor
make fmt
make test
```

构建入口程序：

```bash
make build
```

运行一个本地节点：

```bash
make run-node
```

在另一个终端写入测试数据：

```bash
make run-client
```

读取测试数据：

```bash
go run ./cmd/client --addr 127.0.0.1:9001 --command get --key demo
```

查看节点状态：

```bash
go run ./cmd/client --addr 127.0.0.1:9001 --command status
```

## 5. 目录用途

- `cmd/node`：Raft 节点服务入口
- `cmd/client`：客户端占位入口
- `cmd/bench`：实验驱动入口
- `internal/raft`：Raft 核心状态与配置
- `internal/transport`：TCP / QUIC 传输抽象
- `api/proto`：Raft RPC 协议
- `deploy/docker-compose.yml`：本地多节点编排预留

## 6. 建议的下一步安装顺序

第一优先级：

1. Prometheus
2. Grafana
3. `toxiproxy`

## 7. 环境完成标准

环境搭建完成后，应满足以下最小标准：

- `make doctor` 无关键缺失
- `make test` 成功
- `make build` 成功
- `make run-node` 能启动单节点
- 后续补全 `docker compose up` 时，能够同时启动 3 个节点
