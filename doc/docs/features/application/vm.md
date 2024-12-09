# 虚拟机部署 <Badge type="warning" text="v1.3.0" />
agent 以二进制启动时，如果未找到 `kubeconfig`，会自动转为虚拟机模式。本章节展示如何将应用部署于虚拟机（物理机同理）的指定路径下，并且启动，及健康检查。

:::tip
虚拟机模式的 agent 启动配置中，注册 Key 的格式必须满足 `lizardcd-agent.<system_name>.<ip>`，其中 \<system_name\> 自定义，\<ip\> 是本机 IP。
:::

## 新建应用
本节以部署一个 `node_exporter` 到 `10.50.89.44` 为例。其中安装包存放于对象存储（minio）。

### 前提条件
- 已有一个对象存储 `minio`，桶名为 `lizardcd`，安装包存放路径 `http://minio:9000/lizardcd/node_exporter/node_exporter-1.8.1.linux-amd64.tar.gz`。
- 平台设置里添加一个镜像仓库，类型为 `S3`，地址为 http://minio:9000。

### 步骤
导航栏点击 `应用管理`，点击按钮 `新建应用`，按如下图所示填写：

![](/images/application/new-node_exporter.png)

`镜像仓库/制品库` 选择 http://minio:9000，`仓库/项目` 填写 lizardcd，`制品名` 填写 node_exporter

![](/images/application/new-node_exporter-2.png)

`部署前命令/脚本` 通常可执行停止服务、清理旧版本等操作，运行于下载制品之前。

`启动命令/脚本` 运行于下载制品、解压之后（lizardcd 会将制品自动下载于 `部署路径`）。

### 健康检查配置

1. HTTP方式

![](/images/application/healthcheck-http.png)

:::warning
agent 会在虚拟机 以 `GET/POST` 访问 `http://localhost:9100/metrics` 并判断返回码是否为 `200OK`。因此服务必须启动在 `0.0.0.0`，否则将不通。
:::

2. Shell方式

![](/images/application/healthcheck-shell.png)

3. TCP方式

![](/images/application/healthcheck-tcp.png)

:::warning
和 HTTP 检查一样，agent 会探测虚拟机 `localhost` 的 `9100` TCP端口。因此服务必须启动在 `0.0.0.0`，否则将不通。
:::

## 发布应用
选中 `node_exporter`，`更多操作`-`发布`，选择一个版本提交。

![](/images/application/release-node_exporter.png)