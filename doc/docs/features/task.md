# 任务管理 <Badge type="warning" text="v1.2.1" />
本章节将向您展示如何在 lizardcd-ui 管理您的任务。

Lizardcd 的任务类型分成两种：`deploy` 和 `rollout`。在 [应用管理](/docs/features/application) 页面，对应用进行 `发布` 将触发 `deploy` 类型任务；对应用进行 `重启` 将触发 `rollout` 类型任务。
```go
const K8S_TASK_TYPE_ROLLOUT = "rollout"
const K8S_TASK_TYPE_DEPLOY = "deploy"
```

## 任务列表
页面导航点击 `任务管理`，即打开任务列表页面。任务列表如图所示：

![](/images/task/task-list.png)

列表展示了应用名称、任务类型、触发类型、执行结果、状态、所属租户、耗时、操作等。

点击应用名称将跳转到应用管理页面。

如果任务执行失败，鼠标移动到状态列会弹出错误信息。

有操作三个按钮：
- 查看详情：点击弹框展示任务详情，包括任务的标签信息。
- 重跑：点击将以完全相同的参数重新运行任务。
- 删除。

触发类型：
- 手动触发：在应用管理页面发布或重启应用，任务类型为 `手动触发`。
- 其它类型：调用 `RunTask` 接口时传递 `TriggerType` 自定义触发类型。

## 任务详情
点击一项任务最左边的 `>` 按钮展开查看任务详情。如图所示。

任务有三种部署类型，参见 [应用管理](/docs/features/application)：
- 容器部署：列出该应用的每一个工作负载所在集群、命名空间、负载类型（deployments/statefulsets）、容器名称、镜像、Pod 状态、输出信息
- 虚拟机部署：列出目标类型（vm）、目标地址（IP）、制品、状态（agent 部署返回的标准输出信息）、输出信息（agent 部署返回的标准错误信息）
- HTTP部署：列出目标类型（HTTP）、目标地址（HTTP BaseURL）、制品、输出信息（HTTP接口返回值提取的输出信息）