# 平台设置 <Badge type="warning" text="v1.0.0" />
本章节介绍 Lizardcd 的平台设置，部分设置仅有 `admin` 权限的用户可以使用。

## 设置

### 开启对接 Istio
[Istio](https://istio.io/) 是一个云原生的服务网格解决方案，提供微服务调用、流量治理、网关、服务跟踪监控等组件功能。

Lizardcd 通过 [Istio client-go SDK](https://pkg.go.dev/istio.io/client-go/pkg/apis/networking/v1beta1) 对接 `Istio`，该功能要求 agent 所纳管的 Kubernetes 集群已安装好 `Istio`，具备 `istio.io` Group 下的相关 apiVersion，同时 agent 启动时使用的 `serviceaccount` 或 `kubeconfig` 必须具备对 `istio.io` Group 下资源的相关操作权限。

:::tip
如何查看集群是否已经安装 `Istio`，可执行如下命令：
```shell
# kubectl api-versions | grep istio
extensions.istio.io/v1alpha1
networking.istio.io/v1alpha3
networking.istio.io/v1beta1
security.istio.io/v1
security.istio.io/v1beta1
telemetry.istio.io/v1alpha1
```
如能显示 *.istio.io 即表示已安装
:::

:::tip
一个对 `istio.io` Group 资源具有 * 权限的 `role` 参考配置如下：
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  annotations:
    kubesphere.io/creator: system
  creationTimestamp: "2024-04-29T08:31:32Z"
  name: operator
  namespace: default
rules:
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  - apps
  - extensions
  - batch
  - monitoring.coreos.com
  - autoscaling
  - app.k8s.io
  - config.istio.io
  - events.k8s.io
  - snapshot.storage.k8s.io
  - monitoring.coreos.com
  - extensions.istio.io
  - networking.istio.io
  - networking.istio.io
  - security.istio.io
  - security.istio.io
  - telemetry.istio.io
  resources:
  - '*'
  - tasks
  verbs:
  - '*'
```
:::

对接 `Istio` 功能默认是关闭的，管理员可在设置页面开启对接 `Istio`。当启用后，左边导航栏将会出现 `Istio管理` 菜单，详细使用方法参见 [Istio管理](/docs/features/istio)

### 开启对接 Tekton
[Tekton](https://tekton.dev/) 是一个云原生的持续集成项目。Lizardcd 预留了对接 `Tekton` 的功能，预计在v1.4.0版本上线。

### 开启 Helm 包管理
[Helm](https://helm.sh/zh/docs/) 是一个 Kubernetes 的包管理器，类似于 CentOS 的 `yum` 或 Ubuntu 的 `apt`。Helm 源码由 `Helm pkg` 和 `Helm cli` 两部分构成，参考：[https://github.com/helm/helm](https://github.com/helm/helm)。`pkg` 是 Helm 的核心功能，`cli` 调用 `pkg` 实现了命令行工具。

Lizardcd 通过调用 `Helm pkg` 重写了 `Helm cli`，将全部 `Helm cli` 的功能集成入 `lizardcd-ui` 和 `lizardcd-cli`。

`Helm包管理` 功能默认是关闭的，管理员可在设置页面开启 `Helm包管理`。当启用后，左边导航栏将会出现 `Helm管理` 菜单，详细使用方法参见 [Helm管理](/docs/features/helm)

:::tip
设置对不同租户是独立的，亦即某租户开启了一项设置，和其它租户无关，只有同属一个租户的用户共享相同的设置。
:::

## 租户管理
详见 [用户与权限管理](/docs/features/rbac)

:::tip
该设置仅 `admin` 角色可见
:::

## 用户管理
详见 [用户与权限管理](/docs/features/rbac)

:::tip
该设置仅 `admin` 角色可见
:::

## 镜像仓库管理
镜像仓库是应用发布时用来选择镜像 tag（或者实体机部署时选择安装包）的，需要提前在本页面设置好。目前 Lizardcd 共支持四种类型的镜像仓库：

### Artifactory
杰蛙公司的产品 [Jfrog Artifactory](https://jfrog.com/)，经测试，目前仅支持企业版（因只有企业版才提供 API 接口）。创建该类型仓库时填写的内容：

- 仓库地址：Artifactory 的域名地址，以 http 或 https 开头，不要带任何路径
- 仓库账号，登录 Artifactory 的用户名
- 仓库密码：该登录账户的 `API Key`

### Harbor
[Harbor](https://goharbor.io/) 是一个 Docker 镜像仓库管理平台。创建该类型仓库时填写的内容：

- 仓库地址：Harbor 的域名地址，以 http 或 https 开头，不要带任何路径
- 仓库账号，登录 Harbor 的用户名
- 仓库密码：登录 Harbor 的密码

### DockerHub
Docker 官方镜像仓库。创建该类型仓库时填写的内容：

- 仓库地址：固定地址：https://hub.docker.com
- 仓库账号，登录 Docker 的用户名
- 仓库密码：Personal Access Tokens，关于如何获取详见 [Create and manage access tokens](https://docs.docker.com/security/for-developers/access-tokens/)

:::warning
自 2024.6 月起国内网络无法直接拉取 DockerHub 的镜像，可能需要施展一些魔法。
:::

### S3
兼容S3协议的对象存储，如 `Minio`、`腾讯云COS`、`阿里云OSS` 等。对于实体机部署，可将应用制品（安装包）存放于对象存储，Lizardcd 支持从对象存储下载制品。创建该类型仓库时填写的内容：

- 仓库地址：对象存储地址，以 http 或 https 开头
- 仓库账号，对象存储的 `AccessKey`
- 仓库密码：对象存储的 `SecretKey`

:::tip
镜像仓库管理对不同租户是独立的，亦即某租户添加的仓库，其它租户将看不见，只有同属一个租户的用户能够看见。
:::