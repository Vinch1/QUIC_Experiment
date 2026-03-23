# 项目计划

## Phase 0：脚手架与环境

目标：

- 固化目录结构
- 建立基础入口
- 编写环境文档
- 预留 Docker / Proto / Benchmark 目录

交付物：

- `README.md`
- `docs/setup.md`
- `docs/architecture.md`
- `cmd/node`
- `cmd/client`
- `cmd/bench`

## Phase 1：最小可运行 Raft

目标：

- 3 节点
- 选主
- 心跳
- 日志复制
- 内存状态机

验收标准：

- 单机可启动 3 节点
- 可完成 leader 选举
- client proposal 能复制到多数派

## Phase 2：双传输实现

目标：

- 抽象统一 RPC 传输接口
- 实现 TCP 版本
- 实现 QUIC 版本

验收标准：

- 相同 workload 下可仅切换配置完成 TCP / QUIC 对比
- 两个实现暴露一致的请求语义

## Phase 3：实验框架

目标：

- 构建可重复实验脚本
- 增加网络扰动注入
- 记录 benchmark 结果

核心指标：

- proposal latency
- commit latency
- throughput
- election duration
- reconnect latency
- CPU / memory

## Phase 4：结果分析

建议场景：

1. 稳定局域网基线
2. 50ms / 100ms RTT
3. 抖动场景
4. 丢包场景
5. leader 故障切换
6. 并发 proposal 场景

## Phase 5：增强项

- snapshot
- WAL
- membership change
- Prometheus + Grafana
- AWS 跨地域验证

## 不建议一开始做的内容

- 一次性做完整生产特性
- 先上云再调本地
- 先优化 QUIC 多 stream 再建立 TCP 基线

