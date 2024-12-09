# 链路跟踪
Lizardcd-server 和 lizardcd-agent 基于 [go-zero](https://go-zero.dev/) 框架编写，支持 `OpenTelemetry` 协议的链路跟踪功能，详见 [链路追踪](https://go-zero.dev/docs/tutorials/monitor/index#%E9%93%BE%E8%B7%AF%E8%BF%BD%E8%B8%AA)。

## 概念
提到链路跟踪，或者叫全链路监控，或者叫 APM(Application Performance Management)，具体含义和原理不赘述，开源方案有 skywalking、zipkin、elasticapm 等工具，商业产品有基调听云等等，但在云原生领域，也有一个 CNCF 已毕业项目 jaeger 同样发展迅速。

## 前提条件
部署 jaeger 套件，参考 [云原生链路跟踪工具Jaeger+OpenTelemetry Collector部署详解](https://hongyuxuan.github.io/blog/posts/2128614243/)

## 开启链路跟踪
链路追踪支持 jaeger、zipkin，Lizardcd 使用的是 jaeger，只需要在配置中添加几行配置就可开启，无需修改任何代码，示例如下：

- lizardcd-server.yaml
```yaml
Name: lizardServer
Host: 0.0.0.0
Port: 5117
Timeout: 60000
......
Telemetry:
  Name: lizardcd-server
  Endpoint: http://otel-collector:14278/api/traces
  Sampler: 1.0
  Batcher: jaeger
```

- lizardcd-agent.yaml
```yaml
Name: LizardAgent
ListenOn: 0.0.0.0:5017
Timeout: 60000
......
Telemetry:
  Name: lizardcd-agent
  Endpoint: http://otel-collector:14278/api/traces
  Sampler: 1.0
  Batcher: jaeger
```

上述配置假定 jaeger 和 Opentelemetry Collector 已经部署完毕，其中 Opentelemetry Collector 监听在 `otel-collector:14278`。