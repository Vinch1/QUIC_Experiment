# 架构说明

## 1. 设计目标

本项目不是直接追求“完整生产级 Raft”，而是优先回答以下实验问题：

- QUIC 是否适合作为 Raft 节点间传输层
- 在同一套 Raft 核心下，QUIC 与 TCP 的性能差异是什么
- QUIC 多 stream 是否能改善不同类型 Raft 消息的相互干扰

## 2. 架构原则

- 同一套 Raft 核心
- 同一套消息模型
- 同一套存储策略
- 同一套 benchmark 场景
- 唯一变量是传输层

## 3. 模块划分

### `internal/raft`

负责节点状态、角色切换、选主、日志复制和后续快照扩展。

### `internal/transport`

提供统一接口，分为：

- `tcp`
- `quic`

第一阶段只要求接口一致，第二阶段再实现真实网络行为。

### `internal/storage`

第一阶段使用内存存储，避免磁盘因素干扰网络实验结果；后续再补充 WAL 和 snapshot。

### `internal/statemachine`

先采用最简单的 KV 状态机，用于验证 proposal → replicate → apply 路径。

### `cmd/bench`

用于承载实验场景配置、批量执行和结果汇总。

## 4. 分阶段策略

### Phase 0

完成脚手架、文档、接口与目录布局。

### Phase 1

实现最小 Raft：

- leader election
- append entries
- basic commit path

### Phase 2

实现两套真实传输：

- TCP
- QUIC

### Phase 3

实现 benchmark：

- 基线延迟
- 吞吐
- 选主时间
- 故障恢复

### Phase 4

再加入增强项：

- snapshot
- persistent storage
- membership change
- QUIC stream-level optimization

## 5. 为什么不先做完整 Raft

因为当前目标是“思路验证”，不是直接做生产级复制状态机。若一开始加入快照、成员变更、云部署和 tracing，实验成本会显著上升，而且很难快速定位网络层收益是否真实存在。

