# ui 部署
本文档将向您展示如何独立安装 lizardcd-ui 组件。 

## 在 k8s 集群上使用 helm 安装
可 [在此](https://artifacthub.io/packages/helm/lizardcd/lizardcd) 查看最新的 Helm Charts 安装包。

添加仓库：
```shell
helm repo add lizardcd https://funnyfpf.github.io/lizardcd/
```
安装chart：
```shell
helm install lizardcd-ui lizardcd-lizardcd \ 
--set ui.enabled=true \
--set server.enabled=false \
--set agent.enabled=false \
--set etcd.enabled=false \
--set ui.externalServer.server="192.168.200.168" \
--set ui.externalServer.port=5117 -n default
```
:::tip
请将上述`192.168.200.168`及`5117`换成实际server ip地址和端口
:::

以上操作将在您的 Kubernetes 集群的 default 命名空间下安装 `lizardcd-ui` 工作负载，以及相应的`service`、`configmap`。
:::warning
默认 charts 并不会创建 `Ingress`，您可以编辑 `values.yaml` ，将 `Ingress` 打开，并填入您的域名：
```yaml
ui:
  ingress:
    enabled: false
    path: /
    hostname: lizardcd-ui.local
```
:::
## 访问 Lizardcd
打开浏览器，输入您配置的 Ingress 域名或通过 nodePort 的方式访问 lizardcd-ui。

如果能看到如下页面，则说明安装成功。否则，请检查您的工作负载 `lizardcd-ui` 的日志是否有报错。

![](/images/login.jpg)

若您的 lizardcd-ui 连接的是 Kubernetes 上部署的 lizardcd-server，则在您的 Kubernetes 任务页面找到 `lizardcd-lizardcd-initjob`，查看该任务的日志输出，找到初始化的 `admin` 用户密码并记住此密码。

若连接的是虚机部署的 lizardcd-server，则在初始化数据库的时候查看任务日志输出，详见 [server 部署 linux 部分](/docs/deploy/server.html#linux)

回到 Lizardcd 页面，用 `admin` 用户密码登录，您应该能看到如下页面：

![](/images/quick-deploy.jpg)

## 卸载
```sh
helm uninstall lizardcd-ui
```
