# HTTP部署 <Badge type="warning" text="v1.3.0" />
Lizardcd 支持在 ui 直接配置 HTTP 部署，无需编写任何一行代码，即可对目标接口发起 HTTP 调用完成部署工作。同时还根据配置的健康检查接口对目标发起调用。根据配置的 ` 检查结束判断` 从返回值（返回值必须是 JSON 格式）提取 `JSONPATH`，与配置的关键词进行正则表达式匹配，判断任务是否结束。如果应用还配置了 `检查成功判断`，则从返回值中提取 `JSONPATH`，判断任务是否成功；否则接口返回 `200OK` 即认为成功。
:::tip
目前仅支持 `http`，暂不支持 `https`。HTTP 部署不依赖 agent。
:::

## 新建应用
本节以调用 `AWX` 平台实现 HTTP 部署 node_exporter 为例。

### 前提条件
- 已部署 `AWX` 平台。
- `AWX` 平台已创建 `job_template`，id 为 1。
- `AWX` 平台已创建 `inventory`，id 为 1。
- 已有一个对象存储 `minio`，桶名为 `lizardcd`，安装包存放路径 `http://minio:9000/lizardcd/node_exporter/node_exporter-1.8.1.linux-amd64.tar.gz`。

### 步骤
导航栏点击 `应用管理`，点击按钮 `新建应用`，按如下图所示填写：
![](/images/application/new-http.png)

![](/images/application/new-http-2.png)

`Base URL` 填写 `AWX` 平台的地址。

`Http Header` 填写 `AWX` 的 `Bearer Token`。

`部署Method` 选择 POST。

`部署Path` 填写 `AWX` 的固定接口 `/api/v2/job_templates/1/launch/`，其中 1 是 `job_template_Id`。

`提交Body` 固定格式如下：
```json
{
  "inventory": "1",
  "extra_vars": "{\"deploy_dir\":\"/opt\",\"app_name\":\"node_exporter\",\"artifact_url\":\"{{artifact_url}}\",\"env\":\"dev\"}"
}
```
其中 `extra_vars` 根据您 job 所需的参数填写，但 `artifact_url` 是固定格式，取值为变量 `artifact_url`，该变量在执行时会被替换为实际的制品地址。

`部署成功判断` 根据 `AWX` 的 `launch` 固定返回接口格式，如调用成功返回值的 JSON 中 `job` 字段为 `job_id`；否则为 `null`，因此正则表达式写 `\d+`。

![](/images/application/new-http-3.png)

`健康检查Path` 为 `AWX` 的固定写法，其中变量 `response\$\.job` 将被替换为上一步执行完后的返回值 JSON 中的 job_id。

`检查结束判断` 和 `检查成功判断` 都将从返回值中提取 `JSONPATH`, 并和设置的关键词进行正则表达式匹配。

`提取输出信息` 将从返回值中提取 `JSONPATH`，并显示在任务执行结果中。

## 发布应用
同虚拟机部署。