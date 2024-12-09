# server 配置

## 配置文件
::: details server 的完整配置文件 `lizardcd-server.yaml`
```yaml
Name: lizardServer
Host: 0.0.0.0
Port: 5117
Timeout: 60000
Log:
  Encoding: plain
  Level: info
Prometheus:
  Host: 0.0.0.0
  Port: 15117
  Path: /metrics
Telemetry:
  Name: lizardcd-server
  Endpoint: http://otel-collector:14278/api/traces
  Sampler: 1.0
  Batcher: jaeger
Auth:
  AccessSecret: wLnOk8keh/WO5u7lX8H1dB1/mcuHvnI/jfWCMXMPg9o=
  AccessExpire: 86400
Etcd:
  Address: 10.50.89.17:2379
Consul:
  Address: 10.50.89.17:8500
Nacos:
  Address: 10.100.67.41:8848
  NamespaceId: public
  Group: default
  Username: test
  Password: test
Sqlite: ./lizardcd.db
ServicePrefix: my-lizardcd-
Rpc:
  Timeout: 2000 # millisecond
  KeepaliveTime: 600 # seconds
  RetryInterval: 600 # seconds 
```
:::

详细说明：

<b>基础配置</b>
| 配置项 | 说明 | 可选值 | 默认值 |
| --- | ------ | --- | --- |
| Name | server 的名称 | 任意 | lizardServer |
| Host | 监听地址 | | 0.0.0.0 |
| Port | 监听端口 | | 5117 |
| Timeout | 访问 server 的超时时间，单位ms | | 60000 |
| Log.Encoding | 日志格式 | plain,json | plain |
| Log.Level | 日志级别，受限于 go-zero 的功能，没有 warn 级别 | debug,info,error | info |
| Prometheus.Host | server 的 metrics 监听地址 | | 0.0.0.0 |
| Prometheus.Port | server 的 metrics 监听端口 | | 15117 |
| Prometheus.Path | server 的 metrics 路径 | | /metrics |
| Telemetry.Name | Opentelemetry 中的服务名称 | | lizardcd-server |
| Telemetry.Endpoint | Opentelemetry 服务端地址，需要单独部署，详见 [链路跟踪](/docs/operate/trace) | | http://otel-collector:14278<br>/api/traces |
| Telemetry.Sampler | 采样版本，参考 [Sampling](https://opentelemetry.opendocs.io/docs/specs/otel/trace/sdk#sampling) | | 1.0 |
| Telemetry.Batcher | Opentelemetry 导出器，参考 [Exporters](https://opentelemetry.io/docs/languages/go/exporters/) | jaeger,zipkin,file | jaeger |\
| ServicePrefix| 注册的服务名前缀，server 需配置和 agent 相同前缀 | | my-lizardcd- |
| Auth.AccessSecret | JWT Token 的密钥 | | |
| Auth.AccessExpire | JWT Token 的过期时间，单位秒 | | 86400 |
| Rpc.Timeout | gRPC 调用超时时间，单位毫秒 | | 2000 |
| Rpc.KeepaliveTime | gRPC 长连接保持时间，单位秒 | | 600 |
| Rpc.RetryInterval | gRPC 长连接断掉后重试间隔，单位秒 | | 600 |

<b>Etcd 配置</b>
:::tip
服务注册类型 etcd、consul、nacas 可以多选。如多选，将分别从各个注册中心获取 agent
:::

| 配置项 | 说明 | 可选值 | 默认值 |
| --- | ------ | --- | --- |
| Etcd.Address | etcd 地址，多个地址用英文逗号隔开。Helm 安装时如使用内置 etcd，则自动填充；否则需要指定外部 etcd 地址 | \<ip\>:\<port\>,\<ip\>:\<port\>,... | 同 agent |

<b>Consul 配置</b>
| 配置项 | 说明 | 可选值 | 默认值 |
| --- | ------ | --- | --- |
| Consul.Address | consul 任意节点地址或 VIP 地址。Helm 安装时如使用内置 consul，则自动填充；否则需要指定外部 consul 地址 | \<ip:port\> | 同 agent |

<b>Nacos 配置</b>
| 配置项 | 说明 | 可选值 | 默认值 |
| --- | ------ | --- | --- |
| Nacos.Address | Nacos 地址。Helm 安装时如使用内置 Nacos，则自动填充；否则需要指定外部 Nacos 地址 | \<ip:port\> | |
| Nacos.NamespaceId | Nacos 的 namespace ID | | public |
| Nacos.Group | Nacos 的 Group name | | default |
| Nacos.Username | Nacos 登录用户名 | | |
| Nacos.Password | Nacos 登录密码 | | |

## 启动参数
启动参数说明可用 `--help` 查看
```shell
# ./lizardcd-server --help
usage: lizardcd-server [<flags>]

Flags:
  -h, --help                   Show context-sensitive help (also try --help-long and --help-man).
  -f, --config=""              config file
      --consul-addr=""         Consul address.
      --etcd-addr=""           Etcd address.
      --nacos-addr=""          Nacos address.
      --nacos-namespace-id=""  Nacos namespaceId.
      --nacos-username=""      Nacos username.
      --nacos-password=""      Nacos password.
      --nacos-group=""         Nacos group.
      --service-prefix=""      Prefix of service key for registry. Can be empty
      --log.level=""           Log level.
      --http-addr=""           HTTP listen address.
      --metrics-addr=""        Prometheus metrics listen address.
      --db=""                  SQLite database file.
      --access.secret=""       Jwt token accessSecret.
      --access.expire=ACCESS.EXPIRE
                               Jwt token expire time.
      --version                Show application version.
```

server 内置了 `SQLite` 作为 Lizardcd 的数据库，--db 参数指定 `SQLite` db文件地址。