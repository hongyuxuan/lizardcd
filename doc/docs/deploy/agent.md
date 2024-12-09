# agent 部署
本文档将向您展示如何独立安装 lizardcd-agent 组件。 有两种选择：
1. k8s 集群：在现有 k8s 集群上使用 helm 方式安装
2. 本地机器：在本地机器上使用二进制文件方式安装

## 在 k8s 集群上使用 helm 安装
可 [在此](https://artifacthub.io/packages/helm/lizardcd/lizardcd) 查看最新的 Helm Charts 安装包。

添加仓库：
```shell
helm repo add lizardcd https://funnyfpf.github.io/lizardcd/
```
安装chart：
```shell
helm install lizardcd-agent lizardcd/lizardcd --set externalEtcd.hosts="{10.0.0.2:2379,10.0.0.3:2379,10.0.0.4:2379}" \
--set server.enabled=false \
--set agent.enabled=true \
--set agent.clusterName="k8s" \
--set ui.enabled=false \
--set etcd.enabled=false -n default
```
:::tip
请将上述 10.0.0.2:2379,10.0.0.3:2379,10.0.0.4:2379 换成实际 etcd 地址，agent.clusterName 可根据实际情况自定义。
:::

以上操作将在您的 Kubernetes 集群的 default 命名空间下安装 `lizardcd-agent`工作负载，以及相应的`service`、`configmap`、`pvc`。
## 卸载
```sh
helm uninstall lizardcd-agent
```
## 在本地机器上安装
本地机器支持 Linux 和 Windows 系统
### Linux
下载压缩包并解压
```sh
tar zxf lizardcd-agent-linux-amd64-<version>.tar.gz 
```

将 lizardcd-agent 添加到路径：
```sh
export PATH=$PWD/bin:$PATH
```
启动 lizardcd-agent 服务：
```sh
lizardcd-agent --etcd-host 10.50.89.17:2379 --service-key lizardcd-agent.*.k8s --kubeconfig ~/.kube/config --grpc-addr 0.0.0.0:5017
```
:::tip
上述 10.50.89.17:2379 请换成实际etcd地址。
:::

lizardcd-agent 还能以配置文件启动，启动方法：
```sh
lizardcd-agent -f <path>/lizardcd-agent.yaml
```
其中 `lizardcd-agent.yaml` 配置文件参考 [agent 配置](/docs/deploy/configure/agent)

### Windows
下载压缩包 `lizardcd-agent-windows-amd64-<version>.tar.gz ` 并解压到任意位置

启动 lizardcd-agent 服务：
```sh
.\lizardcd-agent.exe --etcd-host 10.50.89.17:2379 --service-key lizardcd-agent.*.k8s --kubeconfig <path>\config --grpc-addr 0.0.0.0:5017
```
:::tip
上述 10.50.89.17:2379 请换成实际 etcd 地址。
:::