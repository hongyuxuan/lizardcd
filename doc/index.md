---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: "Lizardcd - Cloud Native Continuous Delivery Platform"
  # text: "A lightweight cloud native continuous delivery project"
  tagline: 轻量级的云原生持续交付平台，支持Kubernetes部署、实体机部署，以及通过HTTP接口对接第三方平台进行应用部署
  actions:
    - theme: brand
      text: 快速开始
      link: /docs/quickstart/deployment
    - theme: alt
      text: 文档
      link: /docs/introduce

features:
  - title: 应用部署
    details: 原生支持以 Deployments 和 StatefulSets 两种工作负载方式进行应用部署，同时也支持在实体机上（如虚拟机、物理机）进行应用部署
  - title: Helm部署
    details: 支持全部 Helm Cli 功能，包括 Repo 添加、删除，Release 安装、卸载、更新、历史记录，README 和 values 文件查看、编辑等
  - title: 多集群管理
    details: 通过 agent 管理多个 Kubernetes 集群，支持将应用部署于这些 Kubernetes 集群上（通过工作负载和 Helm）
  - title: 对接服务网格
    details: Lizardcd 目前支持对接 Istio 实现集成部分网格功能，如灰度发布、流量控制等。未来将接入更多网格项目
  - title: 应用管理
    details: Lizardcd 以应用为中心，以任务形式执行部署动作。支持应用的发布、重启、扩缩容、删除、重装等
  - title: API接口
    details: Lizardcd-server提供 Restful API 支持对接上游 CI 平台，以实现 DevOps 流水线
---

