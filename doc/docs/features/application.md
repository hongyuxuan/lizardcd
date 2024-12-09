# 应用管理 <Badge type="warning" text="v1.2.1" />
应用部署是 Lizardcd 的核心功能，本章节讲述如何创建/编辑三种部署方式的应用、每种部署方式的应用场景；发布应用，查看应用的执行结果；发布失败后的处理方式等。

## 容器部署
Lizardcd 的原生部署方式。支持 Kubernetes `deployment` 和 `statefulset` 两种工作负载。

与 ArgoCD 的 GitOps 模式不同，ArgoCD 需要在 Git 上维护一份工作负载的 YAML 配置，并且通过 Git 的 `Webhook` 触发应用部署。

Lizardcd 是直接调用 `client-go` 中工作负载的 `Patch` 接口，参考如下 `DeploymentInterface` 接口定义：
```go
type DeploymentInterface interface {
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Deployment, err error)
}
```

其中 `PatchType` 有如下几种类型：
```go
const (
	JSONPatchType           PatchType = "application/json-patch+json"
	MergePatchType          PatchType = "application/merge-patch+json"
	StrategicMergePatchType PatchType = "application/strategic-merge-patch+json"
	ApplyPatchType          PatchType = "application/apply-patch+yaml"
)
```

Lizardcd 使用的是其中的 `StrategicMergePatchType`。关于更多 `PatchType` 的用法说明详见 [Update API Objects in Place Using kubectl patch](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/)。

关于 Lizardcd 容器部署的使用方法详见 [容器部署](/docs/features/application/k8s)

## 虚拟机/物理机部署
Lizardcd 的扩展部署模式。Lizardcd 设计之初是没有虚拟机/物理机部署方式的，但在公司内部的使用过程中，有越来越多的用户提出了这个需求。本着用户就是上帝的原则，Lizardcd 采纳了这个提议。

实体机的应用部署本质上是远程自动化任务执行，在这一领域，目前使用最多的方案是 `Ansible`，以及其图形化平台 `Ansible Tower`（企业版） 或 `AWX`（Tower 的开源版）。其它的诸如 `SaltStack` 和 `Puppet` 等已经没什么人用了。除此之外还有一些云厂商提供的平台，比如腾讯的 `蓝鲸PAAS平台`。

无论哪种产品，其实现方式上无外乎 `有 agent` 版和 `无 agent` 版。`Ansible` 是 `无 agent` 版，而 Lizardcd 采用的是 `有 agent` 版。Lizardcd 的 agent 除了能够纳管 Kubernetes 之外，还能部署于实体机上做任务执行，其本质是通过 `exec.Command` 系统调用执行 `shell` 脚本。

关于 Lizardcd 虚拟机部署的使用方法详见 [虚拟机部署](/docs/features/application/vm)

## HTTP部署
Lizardcd 的增强部署模式。对于大多数第三方部署平台（这里的第三方是相对于 Lizardcd 来说），例如 `ArgoCD`、`Ansible Tower/AWX`、`Apache Dolphinscheduler`等，只要提供了 HTTP 接口，Lizardcd 支持通过简单配置的方式，将应用部署请求转发到这些第三方平台，并通过配置的 `健康检查` 接口从第三方平台获取部署状态和部署结果。

HTTP 部署模式是 Lizardcd 的终极部署方法，理论上可以对接一切提供接口的第三方平台。对于企业内部已现存多种部署平台，想要集中管理；或者某些应用无法用 Lizardcd 直接部署的用户来说，是一个非常好的选择。

关于 Lizardcd HTTP部署的使用方法详见 [HTTP部署](/docs/features/application/http)