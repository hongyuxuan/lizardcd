---
outline: deep
---

# 快速部署Lizardcd
本章节讲述如何将 Lizardcd 快速部署于一个 Kubernetes 集群中。

## 前提条件
- Kubernetes >= 1.19
- Helm >= 3.8

## 使用Helm安装
可[在此](https://artifacthub.io/packages/helm/lizardcd/lizardcd)查看最新的 Helm Charts 安装包。

添加仓库：
```shell
helm repo add lizardcd https://funnyfpf.github.io/lizardcd/
```

安装chart：
```shell
helm install my-lizardcd lizardcd/lizardcd --version 2.0.3 -n default
```

注意，默认 charts 并不会创建 `Ingress`，您可以编辑 `values.yaml` ，将 `Ingress` 打开，并填入您的域名：
```yaml
ingress:
    enabled: false
    path: /
    hostname: lizardcd-server.local
```

然后使用 `upgrade` 命令更新：
```shell
helm upgrade my-lizardcd -f values.yaml
```

以上操作将在您的 Kubernetes 集群的 default 命名空间下安装 `etcd`、`lizardcd-server`、`lizardcd-agent`、`lizardcd-ui` 四个工作负载，并且配置好了 `Service`、`Ingress`、`ServiceAccount`、`ConfigMap`等资源。

## 访问Lizardcd
打开浏览器，输入您配置的 Ingress 域名（您可能需要手动设置 hosts）

http://lizardcd-server.local

如果能看到如下页面，则说明安装成功。否则，请检查您的工作负载 `lizardcd-server` 的日志是否有报错。

![](/images/login.jpg)

在您的 Kubernetes 任务页面找到 `lizardcd-lizardcd-initjob`，查看该任务的日志输出，找到初始化的 `admin` 用户密码并记住此密码。

回到 Lizardcd 页面，用 `admin` 用户密码登录，您应该能看到如下页面：

![](/images/quick-deploy.jpg)