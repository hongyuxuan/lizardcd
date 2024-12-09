# 容器部署 <Badge type="warning" text="v1.2.1" />
本章节展示如何创建一个容器部署方式的应用

## 新建应用
导航栏点击 `应用管理`，点击按钮 `新建应用`，按如下图所示填写：

![](/images/application/new-k8s-app.png)

<b>基础配置</b>
- 应用名称：任意名称
- 部署方式：容器
- 镜像仓库/制品库：选择一个 [镜像仓库管理](/docs/features/platform.html#镜像仓库管理) 章节添加的仓库
- 仓库/项目：Artifactory 填写仓库名；Harbor 填写项目名；DockerHub 填写 namespace；S3 填写桶（bucket）名
- 镜像名/制品路径：剩余的镜像路径。Lizardcd 会将上述三个字段拼成实际的镜像 URL 并从选择的镜像仓库中获取
- 设置标签：可以为应用设置标签，标签 key 支持重复

<b>容器部署配置</b>
- 开启灰度发布：参考 [开启对接Istio](/docs/features/platform.html#%E5%BC%80%E5%90%AF%E5%AF%B9%E6%8E%A5istio) 章节设置，开启对接 Istio；如不需要，保持关闭
- 工作负载：点击 `+` 添加一个工作负载
  - 容器集群：agent 注册时设置的 Key 会标识其所纳管的集群名称，参考 [agent配置](/docs/deploy/configure/agent)。选择一个集群。
  - 命名空间：agent 注册时设置的 Key 会标识其所纳管的命名空间，参考 [agent配置](/docs/deploy/configure/agent)。选择一个命名空间。
  - 工作负载类型：按需选择。
  - 工作负载名称：`deployment` 或 `statefulset` 的 `metadata.name`。在最新的 <Badge type="warning" text="v1.3.1" /> 版本中已经改为从对应集群自动读取并列出下拉框选择，不允许直接填写。
  - 容器名称：`spec.template.spec.containers[*].name`。在最新的 <Badge type="warning" text="v1.3.1" /> 版本中已经改为从对应集群自动读取并列出下拉框选择，不允许直接填写。
  - 是否启用：如不需要部署此工作负载则关闭。此功能的用处在于用户临时不想部署某个工作负载时，可以将其禁用。

提交保存即可。

:::tip
一个应用可以添加多个工作负载，每个工作负载可以属于不同集群/命名空间。由此可见，Lizardcd 支持将应用的同一个版本发布到不同的集群/命名空间，这对于多数据中心、异地灾备的应用部署将非常有用。
:::

## 灰度发布
Lizardcd 支持对接 `Istio`，借助 `Istio` 的流量管理功能实现灰度发布。

### 前提条件
- 需要使用灰度发布的集群必须已经安装好 `Istio`。
- Lizardcd 的设置中开启 `Istio` 管理，详见 [开启对接 Istio](/docs/features/platform.html#开启对接istio)。
- 在已装有 `Istio` 的集群中安装 `lizardcd-agent`，并且设置服务注册 Key 为 `lizardcd-agent.default.istiok8s`。
- lizardcd-ui 页面 agent 列表能成功显示上一步注册的 agent。

### 原理
Istio 的灰度发布参考 [配置请求路由](https://istio.io/latest/zh/docs/tasks/traffic-management/request-routing/)，基本步骤如下：
1. 创建两个版本的 `deployment`，并分别打上标签 `version: v1` 和 `version: v2`。
2. 创建一个service，将请求转发到第 1 步创建的两个 `deployment`。
3. 创建 `目标规则`，将两个 `deployment` 分成两个 `subset`。
4. 创建 `虚拟服务`，将流量根据 `基于HTTP头部字段` 或 `基于权重` 路由到两个 `deployment`。

### 步骤
本例将创建一个 `flaskapp` 以说明如何通过 Lizardcd 实现灰度发布。

#### 新建工作负载
1. 创建工作负载 `flaskapp-v1`。导航切换到 `工作负载`-`部署`，集群选择 `istiok8s`，命名空间选择 `default`，点击按钮 `新建工作负载`。从模板导入选择 `application_template_with_service`，编辑 YAML 文件（主要给 labels 添加 version）：
```yaml
kind: Deployment
apiVersion: apps/v1
metadata:
  name: {{ .Appname }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Appname }}
    version: v1
  annotations:
    lizardcd/creator: lizardcd
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ .Appname }}
      version: v1
  template:
    metadata:
      labels:
        app: {{ .Appname }}
        version: v1
    spec:
      volumes:
        - name: host-time
          hostPath:
            path: /etc/localtime
            type: ''
        - name: timezone
          hostPath:
            path: /etc/timezone
            type: ''
      containers:
        - name: {{ .Appname }}-container
          image: {{ .Image }}
          ports:
            - name: tcp
              containerPort: {{ .Port }}
              protocol: TCP
          resources:
            limits:
              cpu: 200m
              memory: 200Mi
            requests:
              cpu: 100m
              memory: 100Mi
          volumeMounts:
            - name: timezone
              readOnly: true
              mountPath: /etc/timezone
            - name: host-time
              readOnly: true
              mountPath: /etc/localtime
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: Always
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      serviceAccountName: default
      serviceAccount: default
      securityContext: {}
      schedulerName: default-scheduler
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 25%
      maxSurge: 25%
  revisionHistoryLimit: 10
  progressDeadlineSeconds: 600
---
kind: Service
apiVersion: v1
metadata:
  name: {{ .Appname }}
  namespace: {{ .Namespace }}
  annotations:
    lizardcd/creator: lizardcd
spec:
  ports:
    - name: tcp
      protocol: TCP
      port: {{ .Port }}
      targetPort: tcp
  selector:
    app: {{ .Appname }}
  type: ClusterIP
  sessionAffinity: None
```

模板变量修改：
- Appname：flaskapp-v1
- Image：dustise/flaskapp:v0.2.7
- Port：80
- Namespace：default

提交保存。

2. 创建工作负载 `flaskapp-v2`。从模板导入选择 `application_template_base`，编辑 YAML 文件（主要给 labels 添加 version）：
```yaml
kind: Deployment
apiVersion: apps/v1
metadata:
  name: {{ .Appname }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Appname }}
    version: v2
  annotations:
    lizardcd/creator: lizardcd
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ .Appname }}
      version: v2
  template:
    metadata:
      labels:
        app: {{ .Appname }}
        version: v2
    spec:
      volumes:
        - name: host-time
          hostPath:
            path: /etc/localtime
            type: ''
        - name: timezone
          hostPath:
            path: /etc/timezone
            type: ''
      containers:
        - name: {{ .Appname }}-container
          image: {{ .Image }}
          ports:
            - name: tcp
              containerPort: {{ .Port }}
              protocol: TCP
          resources:
            limits:
              cpu: 200m
              memory: 200Mi
            requests:
              cpu: 100m
              memory: 100Mi
          volumeMounts:
            - name: timezone
              readOnly: true
              mountPath: /etc/timezone
            - name: host-time
              readOnly: true
              mountPath: /etc/localtime
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: Always
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      serviceAccountName: default
      serviceAccount: default
      securityContext: {}
      schedulerName: default-scheduler
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 25%
      maxSurge: 25%
  revisionHistoryLimit: 10
  progressDeadlineSeconds: 600
```

模板变量修改：
- Appname：flaskapp-v2
- Image：dustise/flaskapp:v0.2.7
- Port：80
- Namespace：default

提交保存。

#### 新建应用
新建 `flaskapp` 应用，如下图所示。

- 开启灰度发布：开启

<b>灰度发布策略：基于权重</b>

新增两个工作负载

工作负载 1：

- 容器集群：istiok8s
- 命名空间：default
- 工作负载名称：flaskapp-v1
- 容器名称：flaskapp-container
- 版本号：v1
- 权重：10

工作负载 2：
- 容器集群：istiok8s
- 命名空间：default
- 工作负载名称：flaskapp-v2
- 容器名称：flaskapp-container
- 版本号：v2
- 权重：90

<b>灰度发布策略：基于头部字段</b>

新增两个工作负载

工作负载 1：

- 容器集群：istiok8s
- 命名空间：default
- 工作负载名称：flaskapp-v1
- 容器名称：flaskapp-container
- 版本号：v1
- 匹配头部字段：
  - 头部键：end-user
  - 匹配方式：exact
  - 头部值：kevin

工作负载 2：
- 容器集群：istiok8s
- 命名空间：default
- 工作负载名称：flaskapp-v2
- 容器名称：flaskapp-container
- 版本号：v2

提交保存。

#### 发布应用
找到 flaskapp 应用，点击 `更多操作`-`发布`，选择一个制品，提交。
![](/images/application/release-flaskapp.png)

#### 查看 Istio 资源
页面导航切换到 `Istio 管理`，在 `目标规则` 页面，集群选择 `istiok8s`，命名空间选择 `default`，可以看到已自动创建出一个目标规则 `flaskapp`。可查看 YAML 信息。

![](/images/application/destinationrule.png)

在 `虚拟服务` 页面，集群选择 `istiok8s`，命名空间选择 `default`，可以看到已自动创建出一个虚拟服务 `flaskapp`。可查看 YAML 信息。

![](/images/application/virtualservice.png)

#### 验证结果
在集群 `istiok8s` 中部署一个验证工作负载：
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "4"
  generation: 4
  name: sleep
  namespace: default
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: sleep
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: sleep
    spec:
      containers:
      - image: dustise/sleep:v0.9.7
        imagePullPolicy: IfNotPresent
        name: sleep
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      serviceAccount: lizardcd
      serviceAccountName: lizardcd
      terminationGracePeriodSeconds: 30
```

进入 Pod 后执行命令验证：

```shell
# kubectl exec -it sleep-5db69b44d4-55szj -- /bin/bash
```

验证基于权重，大约10%比例返回v1，90%比例返回v2
```shell
for((i=0; i<100; i++)); do curl http://flaskapp/env/version;echo ;done
```

验证基于头部字段，kevin全部返回v1，tom全部返回v2
```shell
for((i=0; i<10; i++)); do curl -H 'end-user:kevin' http://flaskapp/env/version;echo ;done
for((i=0; i<10; i++)); do curl -H 'end-user:tom' http://flaskapp/env/version;echo ;done
```