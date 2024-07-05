# Lizardcd介绍

## 什么是Lizardcd
Lizardcd 是一个轻量级的云原生应用部署和持续交付平台，在最新的版本（v1.3.0）中，我们增加了对实体机应用部署和基于 HTTP 接口对接第三方平台（如 ArgoCD、Ansible/AWX）进行应用部署的功能。

Lizardcd 以 server-agent 架构工作运行，基于gRPC通信，详细介绍参考[部署架构](/docs/architecture/design)。Lizardcd 的每个组件均可部署于容器或实体机上，但我们推荐将其以容器方式部署于 Kubernetes。我们提供了 Docker Image 和 二进制安装包。对于 Kubernetes 我们还提供了 Helm Charts 的安装包，便于用户一键快速部署。

Lizardcd 的 agent 用于对接 Kubernetes 的 APIServer，支持以 `deployment` 形式运行于 Kubernetes 集群内部 （这种方式下 agent 使用 `serviceaccount` 与 APIServer 通信），也支持以 `二进制` 形式运行于 Kubernetes 集群之外（此时需要通过 `kubeconfig` 与APIServer 通信）。

Lizardcd 的 server 与 agent 通过服务注册中心通信，我们支持 `etcd`、`consul`、`nacos` 三种注册中心，其中 `etcd` 是默认选项。

## Lizardcd 的功能特性
- 云原生应用部署，同时支持实体机以及通过HTTP接口对接第三方部署/自动化平台，理论上可支持任意提供HTTP接口的平台/系统，如 ArgoCD、Ansible/AWX
- 无需 `kubeconfig`、`ClusterRole`，也无需知道 APIServer 地址、证书，agent 支持 `InCluster` 模式部署，特别适合于非容器集群管理员的普通用户进行应用部署
- 支持多集群管理，通过统一的 `ui` 和 `cli` 集中查看、管理、操作多个集群的工作负载，并进行应用部署
- agent 可以在启动时自动地注册到服务注册中心，只要网络端口可达。之后 server 会自动发现已注册并且网络端口可达的 agent，就此完成纳管过程
- 以应用为中心，提供完善的应用管理。无论应用部署于容器或是实体机，或是在第三方平台创建的一个部署对象，lizardcd 均支持对其进行持续交付
- 对于部署了服务网格（如 `Istio`）的集群，Lizardcd 支持将其接入，并结合应用管理实现灰度发布、流量控制等操作（例如基于header或权重）
- 支持完整的 `Helm Cli` 功能，包括 Repo 添加、删除，Release 安装、卸载、更新、历史记录，README 和 values 文件查看、编辑等
- 内置 `Prometheus Metrics` 指标监控，以及基于 `Opentelemetry` 和 `Jaeger` 的链路跟踪功能

其它更多功能特性详见`功能介绍`