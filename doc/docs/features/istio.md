# Istio 管理 <Badge type="warning" text="v1.2.1" />
本章节将向您展示如何在 lizardcd-ui 管理您的 Istio 资源。

## 前提条件
参考 [平台设置](/docs/features/platform) 开启 `对接Istio`

## Istio 资源
Lizardcd 目前支持 Istio 的 `DestinationRule` 和 `VirtualService` 两种资源类型的管理。

当您按照 [灰度发布](/docs/features/application/k8s.html#灰度发布) 章节为应用开启了灰度发布策略、并创建了应用以后，Lizardcd 会自动在工作负载所在的集群创建 `DestinationRule` 和 `VirtualService` 两个 Istio 资源。打开 `Istio 管理` 页面可以查看这两种资源，如图所示：

![](/images/destinationrule.png)

集群和命名空间选择好以后，会列出 Lizarcd 自动创建的资源。点击右边编辑按钮可以查看 YAML 配置，支持编辑。

您还可以直接点击 `创建目标规则` 或 `创建虚拟服务` 直接创建资源。

更多 Istio 功能接入敬请期待。