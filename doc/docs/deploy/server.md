# server 部署
本文档将向您展示如何独立安装 lizardcd-server 组件。 有两种选择：
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
helm install lizardcd-server lizardcd/lizardcd --set externalEtcd.hosts="{10.0.0.2:2379,10.0.0.3:2379,10.0.0.4:2379}" \
--set server.enabled=true \
--set agent.enabled=false \
--set ui.enabled=false \
--set etcd.enabled=false -n default
```
:::warning
请将上述 `10.0.0.2:2379,10.0.0.3:2379,10.0.0.4:2379` 换成实际 etcd 地址
:::

以上操作将在您的 Kubernetes 集群的 default 命名空间下安装 `lizardcd-server`工作负载，以及相应的`service`、`configmap`、`pvc`。
## 卸载
```sh
helm uninstall lizardcd-server
```
## 在本地机器上安装
本地机器支持 Linux 和 Windows 系统

### Linux
<b>初始化 sqlite 数据库</b>

下载 migrate 压缩包并解压
```sh
tar zxf migrate-linux-amd64-<version>.tar.gz 
 
```
将 migrate 添加到路径：
```sh
export PATH=$PWD/bin:$PATH
```
执行以下命令初始化 sqlite 数据库

```sh
migrate -d <path>/lizardcd.db 
```

看到输出下图所示结果，则表示初始化成功

![](/images/migrate.jpg)

:::tip
建议保存初始化的 admin 密码以备后续使用 ui 登录使用。
:::

<b>安装 lizardcd-server</b>

下载 lizardcd-server 安装包并解压
```sh
tar zxf lizardcd-server-linux-amd64-<version>.tar.gz
```
将 lizardcd-server 添加到路径：
```sh
export PATH=$PWD/bin:$PATH
```
启动 lizardcd-server 服务：
```sh
lizardcd-server --etcd-addr 10.50.89.17:2379 --http-addr=0.0.0.0:5117 --db=<path>/lizardcd.db
```
:::tip
上述 10.50.89.17:2379 请换成实际 etcd 地址。--db 指定上一步初始化数据库生成的 lizardcd.db 文件。
:::

lizardcd-server 还能以配置文件启动，启动方法：
```sh
lizardcd-server -f <path>/lizardcd-server.yaml
```
其中 `lizardcd-server.yaml` 配置文件参考 [server 配置](/docs/deploy/configure/server)

### Windows

<b>初始化 sqlite 数据库</b>

下载 `migrate-windows-amd64-<version>.tar.gz` 压缩包并解压到任意位置。

执行以下命令初始化 sqlite 数据库

```sh
.\migrate.exe -d <path>\lizardcd.db 
```

<b>安装 lizardcd-server</b>

下载压缩包 `lizardcd-server-windows-amd64-<version>.tar.gz` 并解压到任意位置

启动 lizardcd-server 服务：
```sh
.\lizardcd-server.exe --etcd-addr 10.50.89.17:2379 --http-addr=0.0.0.0:5117 --db <path>\lizardcd.db
```
:::tip
上述 10.50.89.17:2379 请换成实际 etcd 地址。--db 指定上一步初始化数据库生成的 lizardcd.db 文件。
:::