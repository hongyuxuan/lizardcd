# 工作负载 <Badge type="warning" text="v1.0.0" />
本章节将向您展示如何在 lizardcd-ui 管理您的工作负载。

## 工作负载列表
在列表上方可以过滤容器 `集群`、`命名空间`，并搜索工作负载名称 `metadata.name`。如图所示：

![](/images/workload-list.jpg)

列表展示了工作负载的名称、状态、更新时间。

:::tip
根据 Kubernetes 的 APIServer 返回值，只有 `Deployment` 才有更新时间 `status.conditions.lastUpdateTime`。而 `Statefulset` 没有更新时间，只有创建时间 `metadata.creationTimestamp`。
:::

对每一项工作负载可以执行如下操作
- 重启。类似于调用 `kubectl rollout restart` 命令。
- 设置镜像：为工作负载的 `第一个container` 重新设置镜像。
- 设置副本：设置 `replicas`。
- 编辑YAML：展示工作负载的 YAML 并且可以编辑、提交。
- 删除：请谨慎操作。

## 工作负载详情
点击一项工作负载的负载名称可以进入工作负载的详情页。

![](/images/workload-detail.jpg)

详情页左边展示了工作负载的基本信息。右边展示了工作负载的副本数、`第一个container` 的 Resource、标签、注解、事件。

工作负载的 Pod 信息，包括更新时间、所在 Node 节点、Pod IP、事件查看等。

展开 Pod，将展示 Pod 下的 `initContainers` 和 `containers` 信息，包括镜像、状态、重启次数。