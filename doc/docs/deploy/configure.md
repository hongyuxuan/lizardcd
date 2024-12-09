# 配置
Lizardcd 的 agent 和 server 基于 [go-zero](https://go-zero.dev/) 框架编写，有两种启动配置参数设置方式
- 使用 YAML 配置文件
- 使用命令行参数

对于同一个配置项，命令行参数优先级大于配置文件，命令行配置值会覆盖配置文件。

关于 agent 启动配置，参见 [agent 配置](/docs/deploy/configure/agent)

关于 server 启动配置，参见 [server 配置](/docs/deploy/configure/server)