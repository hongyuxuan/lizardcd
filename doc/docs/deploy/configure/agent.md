# agent 配置

## 配置文件
::: details agent 的完整配置文件 `lizardcd-agent.yaml`
```yaml
Name: LizardAgent
ListenOn: 0.0.0.0:5017
Timeout: 60000
Log:
  Encoding: plain
  Level: info 
Prometheus:
  Host: 0.0.0.0
  Port: 15017
  Path: /metrics
Telemetry:
  Name: lizardcd-agent
  Endpoint: http://otel-collector:14278/api/traces
  Sampler: 1.0
  Batcher: jaeger
Kubeconfig: ~/.kube/config
Etcd:
  Hosts:
    - 10.50.89.17:2379
  Key: lizardcd-agent.default.k8s
Consul:
  Host: 10.50.89.17:8500
  Key: lizardcd-agent.default.k8s
  TTL: 60
  Meta: # custom your metadata here
    Protocol: grpc
    Service: lizardcd-agent
    Namespace: 'default'
    Cluster: k8s
Nacos:
  Host: 10.100.67.41:8848
  Key: lizardcd-agent.default.k8s
  NamespaceId: public
  Group: default
  Username: test
  Password: test
  Meta:
    Protocol: grpc
    Service: lizardcd-agent
    Namespace: "*"
    Cluster: tektonk8s
ServicePrefix: my-lizardcd-
KubernetesSecretPrefix: default-token
```
:::

详细说明：

<b>基础配置</b>
| 配置项 | 说明 | 可选值 | 默认值 |
| --- | ------ | --- | --- |
| Name | agent 的名称 | 任意 | LizardAgent |
| ListenOn | 监听地址 | | 0.0.0.0:5017 |
| Timeout | 访问 agent 的超时时间，单位ms | | 60000 |
| Log.Encoding | 日志格式 | plain,json | plain |
| Log.Level | 日志级别，受限于 go-zero 的功能，没有 warn 级别 | debug,info,error | info |
| Prometheus.Host | agent 的 metrics 监听地址 | | 0.0.0.0 |
| Prometheus.Port | agent 的 metrics 监听端口 | | 15017 |
| Prometheus.Path | agent 的 metrics 路径 | | /metrics |
| Telemetry.Name | Opentelemetry 中的服务名称 | | lizardcd-agent |
| Telemetry.Endpoint | Opentelemetry 服务端地址，需要单独部署，详见 [链路跟踪](/docs/operate/trace) | | http://otel-collector:14278<br>/api/traces |
| Telemetry.Sampler | 采样版本，参考 [Sampling](https://opentelemetry.opendocs.io/docs/specs/otel/trace/sdk#sampling) | | 1.0 |
| Telemetry.Batcher | Opentelemetry 导出器，参考 [Exporters](https://opentelemetry.io/docs/languages/go/exporters/) | jaeger,zipkin,file | jaeger |\
| Kubeconfig | 二进制部署时需要指定 kubeconfig 文件地址 | | ~/.kube/config |
| ServicePrefix| 注册的服务名前缀，server 需配置相同前缀，并根据此前缀匹配 agent 的服务 | | my-lizardcd- |
| KubernetesSecretPrefix | `InCluster` 部署时所需的 token 名称前缀 | | default-token |

<b>Etcd 配置</b>
:::tip
服务注册类型 etcd、consul、nacas 只能选择其中一种
:::

| 配置项 | 说明 | 可选值 | 默认值 |
| --- | ------ | --- | --- |
| Etcd.Hosts | 数组，etcd 地址。Helm 安装时如使用内置 etcd，则自动填充；否则需要指定外部 etcd 地址 | [\<ip\>:\<port\>,\<ip\>:\<port\>] | |
| Etcd.Key | agent 注册到 etcd 的 Key 值。对于 Kubernetes 纳管，写法是 lizardcd-agent.\<namespace\>.\<cluster\>。其中 namespace 是 agent 所具有权限的 ns，可以写 * 表示对所有 ns 都有权限；cluster 名字自定义；对于实体机纳管，写法是 lizardcd-agent.\<system_name\>.\<ip\>，system_name 自定义，ip 是本机 IP | | lizardcd-agent.default.k8s |

<b>Consul 配置</b>
| 配置项 | 说明 | 可选值 | 默认值 |
| --- | ------ | --- | --- |
| Consul.Host | consul 地址。Helm 安装时如使用内置 consul，则自动填充；否则需要指定外部 consul 地址 | \<ip:port\>,\<ip:port\>,... | |
| Consul.Key | agent 注册到 consul 的 Key 值，同 Etcd.Key | | lizardcd-agent.default.k8s |
| Consul.TTL | agent 注册到 consul 的 服务生存时间，超过后会自动重连，单位秒 | | 60 |
| Consul.Meta | agent 注册到 consul 的 meta 信息， map[string]string 格式 | | |

<b>Nacos 配置</b>
| 配置项 | 说明 | 可选值 | 默认值 |
| --- | ------ | --- | --- |
| Nacos.Host | Nacos 地址。Helm 安装时如使用内置 Nacos，则自动填充；否则需要指定外部 Nacos 地址 | \<ip:port\> | |
| Nacos.Key | agent 注册到 Nacos 的 Key 值，同 Etcd.Key | | lizardcd-agent.default.k8s |
| Nacos.NamespaceId | Nacos 的 namespace ID | | public |
| Nacos.Group | Nacos 的 Group name | | default |
| Nacos.Username | Nacos 登录用户名 | | |
| Nacos.Password | Nacos 登录密码 | | |
| Nacos.Meta | agent 注册到 Nacos 的 meta 信息， map[string]string 格式 | | |

:::tip
服务注册 Key 的格式必须严格遵循如下规则：
- Kubernetes 纳管：lizardcd-agent.\<namespace\>.\<cluster\>
- 实体机纳管：lizardcd-agent.\<system_name\>.\<ip\>

对于 Kubernetes，已注册的 \<namespace\> 和 \<cluster\> 会记录在 server 端，用于租户设置和多集群管理时进行集群选择。

对于实体机，已注册的 \<system_name\> 和 \<ip\> 也会记录在 server 端，用于租户设置（设置租户仅能看到选择的 system_name 下的 IP，便于做系统隔离），\<ip\> 用于应用配置时选择要部署的目标服务器。
:::

## 启动参数
启动参数说明可用 `--help` 查看
```shell
# ./lizardcd-agent --help
usage: lizardcd-agent [<flags>]

Flags:
  -h, --help                   Show context-sensitive help (also try --help-long and --help-man).
  -f, --config=""              config file
      --log.level=""           Log level.
      --consul-host=""         Consul hosts.
      --etcd-host=""           Etcd hosts.
      --nacos-host=""          Nacos hosts.
      --nacos-namespace-id=""  Nacos namespaceId.
      --nacos-username=""      Nacos username.
      --nacos-password=""      Nacos password.
      --nacos-group=""         Nacos group.
      --service-key=""         Service key for registry. Format must be: lizardcd-agent.<namespace>.<cluster>
      --service-prefix=""      Prefix of service key for registry. Can be empty
      --kubeconfig=""          Kubeconfig file, must be specified when agent is out-of-k8s deployed
      --grpc-addr=""           Grpc listen address.
      --metrics-addr=""        Prometheus metrics listen address.
      --version                Show application version.
```