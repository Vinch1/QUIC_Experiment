# quic-raft

一个面向思路验证的实验项目：使用同一套 Raft 核心，对比 `TCP` 与 `QUIC` 作为节点间传输层时的行为差异。

## 当前目标

- 建立最小可运行的 Raft 实验骨架
- 保持传输层可替换，确保 TCP / QUIC 对比公平
- 优先支持本地 3 节点实验，再扩展到弱网与故障注入
- 先完成环境搭建、脚手架和实验计划，再逐步补充协议与算法实现

## 当前阶段

当前仓库提供：

- Go 项目脚手架
- `node` / `client` / `bench` 三个入口
- Raft、传输层、存储层、状态机的最小接口骨架
- 基于 HTTP 的最小控制面，支持 `put` / `get` / `status`
- 本地环境检查脚本
- 项目计划与环境搭建文档
- 预留的 Docker Compose / Proto / netem 目录

当前仓库尚未完成：

- 完整 Raft 选主与日志复制
- TCP / QUIC 真正网络实现
- Protobuf 代码生成
- Prometheus 指标采集
- 自动化 benchmark 执行器

## 目录结构

```text
.
├── api/proto
├── cmd
│   ├── bench
│   ├── client
│   └── node
├── deploy
│   ├── docker-compose.yml
│   └── netem
├── docs
├── internal
│   ├── cluster
│   ├── metrics
│   ├── raft
│   ├── statemachine
│   ├── storage
│   └── transport
├── scripts
├── Dockerfile
├── Makefile
└── go.mod
```

## 快速开始

1. 检查本地依赖

```bash
make doctor
```

2. 运行节点脚手架

```bash
make run-node
```

3. 写入一个键值

```bash
make run-client
```

4. 读取一个键值

```bash
go run ./cmd/client --addr 127.0.0.1:9001 --command get --key demo
```

5. 查看节点状态

```bash
go run ./cmd/client --addr 127.0.0.1:9001 --command status
```

## 推荐开发顺序

1. 完成 `api/proto/raft.proto`
2. 将控制面与节点间传输解耦
3. 补齐 `internal/transport/tcp` 和 `internal/transport/quic`
4. 在 `internal/raft` 中实现最小选主和日志复制
5. 在 `cmd/bench` 中补充基线测试场景

## 文档

- 环境搭建：`/Users/leo/Code/Local/quic_experiment/docs/setup.md`
- 架构说明：`/Users/leo/Code/Local/quic_experiment/docs/architecture.md`
- 实验计划：`/Users/leo/Code/Local/quic_experiment/docs/experiment-plan.md`
