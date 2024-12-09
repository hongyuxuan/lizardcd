# 指标监控
Lizardcd-server 和 lizardcd-agent 基于 [go-zero](https://go-zero.dev/) 框架编写，支持 [Prometheus](http://prometheus.io) 格式的指标收集。详见 [指标监控](https://go-zero.dev/docs/tutorials/monitor/index#%E6%8C%87%E6%A0%87%E7%9B%91%E6%8E%A7)。

## 开启指标收集
只需要在配置中添加几行配置就可开启，无需修改任何代码，示例如下：
- lizardcd-server.yaml
```yaml
......
Prometheus:
  Host: 0.0.0.0
  Port: 15117
  Path: /metrics
```

- lizardcd-agent.yaml
```yaml
......
Prometheus:
  Host: 0.0.0.0
  Port: 15017
  Path: /metrics
```

## 访问指标监控
```sh
$ curl http://127.0.0.1:15117/metrics
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 7.3427e-05
go_gc_duration_seconds{quantile="0.25"} 7.3427e-05
go_gc_duration_seconds{quantile="0.5"} 7.3427e-05
go_gc_duration_seconds{quantile="0.75"} 7.3427e-05
go_gc_duration_seconds{quantile="1"} 7.3427e-05
go_gc_duration_seconds_sum 7.3427e-05
go_gc_duration_seconds_count 1
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 12
```

go-zero 默认的指标如下：

<b>RPC Server</b>
| 指标名 | Label | 说明 |
|--- | --- | --- |
| rpc_server_requests_duration_ms | method | Histogram，耗时统计单位为毫秒 |
| rpc_server_requests_code_total | method、code | Counter，错误码统计 |

<b>RPC Client</b>
| 指标名 | Label | 说明 |
|--- | --- | --- |
| rpc_client_requests_duration_ms | method | Histogram，耗时统计单位为毫秒 |
| rpc_client_requests_code_total | method、code | Counter，错误码统计 |

<b>HTTP Server</b>
| 指标名 | Label | 说明 |
|--- | --- | --- |
| http_server_requests_duration_ms | path | Histogram，耗时统计单位为毫秒 |
| http_server_requests_code_total | path、code | Counter，错误码统计 |