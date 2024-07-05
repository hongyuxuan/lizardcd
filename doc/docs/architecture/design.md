# 部署架构
本章节介绍Lizardcd的各个组件以及架构设计。

## 组件介绍
Lizardcd由以下4个组件构成：
- lizardcd-agent：agent 使用 Kubernetes `client-go` SDK 与 APIServer 通信。agent 启动时监听在 gRPC 端口，同时将自己注册到`服务注册中心`。Lizardcd 目前支持三种注册中心： `etcd`、`consul` 或 `nacas`（使用 Helm 部署时，可选择是否同时部署 etcd/consul/nacos，或者使用一个外部的etcd/consul/nacos）。agent 的启动配置文件需要指定自己注册时使用的 `Key`，`Key` 的取值必须满足一定格式，详见 [agent 配置说明](/docs/deploy/configure/agent)。agent 是无状态的，可以任意横向扩展。
- lizardcd-server：server 启动后将持续监听`服务注册中心`，并根据 agent 注册的地址与 agent 的 gRPC 端口建立长连接，同时将该条连接保存于 server 自身的内存中。server 通过 Restful API 与 lizardcd-ui 或 lizardcd-cli 通信，接收图形化或客户端指令，对于部分指令通过 gRPC 下发给 agent 执行。
- lizardcd-ui：Lizardcd 的图形化界面，主要的平台入口。
- lizardcd-cli：Lizardcd 的客户端工具，目前支持 Windows 和 Linux。注意 cli 仅支持部分 Lizardcd 功能，详见 [客户端用法](/docs/cli/cli)

## 架构设计
Lizardcd 的后端服务（包括 server 和 agent）和客户端（cli）全部由Golang编写，基于 [go-zero](https://go-zero.dev/) 框架。前端 ui 由 Javascript 和 Vue3 编写。

### 总体架构
Lizardcd的总体架构如下：

![](https://project-1255547500.cos.ap-beijing.myqcloud.com/lizardcd%2Flizardcd%E6%9E%B6%E6%9E%84%E5%9B%BE.png)

### agent设计
Lizardcd 的 agent 支持以 `deployment` 形式部署于 Kubernetes 中，此时 agent 通过 `serviceaccount` 与 APIServer 通讯，无需 `kubeconfig`，也不用知道 APIServer 的地址和证书。agent所具有的权限取决于 `serviceaccount` 的 RBAC 权限。详见 [RBAC权限设置](/docs/deploy/configure/rbac)。

agent 还支持以`二进制`形式运行于 Windows 或 Linux 系统上，此时 agent 需要通过 `kubeconfig` 与 APIServer 通信。agent 所具有的权限取决于 `kubeconfig` 中设置的 user 的 RBAC权限。

每一个 agent 的纳管范围（通常指namespace）取决于其所使用的 `serviceaccount` 或 `kubeconfig` 的权限范围，假设其所使用的 `kubeconfig` 对集群所有 namespace 都有权限，那么 agent 的纳管范围即该集群的所有 namespace。而通常一个 agent 最多纳管一个集群（`kubeconfig` 中的 `current-context` 指向的 APIServer）。

### server设计
server 的纳管范围是所有已连通的 agent 纳管范围的集合，而我们通常说一个 Lizardcd 能纳管多少个集群，实际上是说所有 agent 纳管了多少集群。

server 设计为有状态的，与 server 通信需要先通过 `login` 接口获取一个 `JWT Token`，默认的 `JWT Token` 过期时间是 1 天。

server 与 agent 之间通过 gRPC 协议通信，符合大多数云原生工具的组件间通信习惯（如Prometheus）。server 会将与每一个 agent 建立的 gRPC 连接保存在内存中，其结构如下：
```go
type RpcAgent struct {
	Client        lizardagent.LizardAgent
	ServiceSource string
	Cli           zrpc.Client
	Count         int
}
```
其中 `lizardagent.LizardAgent` 即 agent 的 `gRPC client` 接口类型。`ServiceSource` 标识该 agent 注册于哪一种注册中心（取值：etcd/consul/nacos）。`zrpc.Client` 是 server 端的 Conn 对象。`Count` 标识当前与 server 建立连接的 agent 有多少个实例（agent 可以横向扩容成多实例），每当一个 agent 的新实例添加到`服务注册中心`时，server 会捕捉到该 agent 注册的 服务名 Key，如果 Key 已存在，则 Count ++。

当 server 收到上游指令时，其快速从内存中取出一条对应的连接，即 `RpcAgent` 对象，并通过对象中的 `Client` 对 agent 发起 gRPC 调用。这种设计的一个弊端是，当 agent 数量不断增大时，server 的内存会持续增长，因此 Lizardcd 也许并不适合海量规模的集群管理，但对于中小规模的集群，其具有相当可观的性能表现。