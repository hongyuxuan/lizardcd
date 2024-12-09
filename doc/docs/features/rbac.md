# 用户与权限管理 <Badge type="warning" text="v1.2.1" />
本章节介绍 Lizardcd 的用户、租户、角色和权限体系设计。

用户使用密码登录平台，用户属于租户（一个租户下可以拥有多个用户，但目前一个用户只能属于一个租户），用户具有一种角色，角色具有对应的权限。详细介绍如下。

## 租户
租户的设计目的之一，是用来对 Kubernetes 的 namespace 进行可见性区分，目的之二是用于对租户创建的各类资源（如镜像仓库、应用、任务、YAML模板）进行可见性区分。

agent 启动时注册的 Key 标识了它所纳管的 namespace，或者当 agent 用于实体机部署时，Key 标识了它所属于的系统（详见 [agent配置](/docs/deploy/configure/agent)），server 会保管所有的 namespace（和系统）列表。当添加租户时，可以为租户设置它所能看见的 namespace（或系统，多选）。之后，该租户仅能看见已设置的 namespace（和系统） 下的资源，其它 namespace（和系统） 下的资源将不可见。

只有 `admin` 角色才能管理租户。Lizardcd 默认内置了一个 `admin` 租户。添加及设置租户的方法如下：

1. 使用 `admin` 用户（属于租户 admin）登录 Lizardcd， 点击顶部导航栏`配置`按钮，并将页面切换到`租户管理`（只有 `admin` 角色才有这个页面），可以看到默认具有一个 `admin` 租户；

![](/images/tenant-list.png)

2. 点击`新建租户`，输入租户名，并设置命名空间，提交保存即可。

![](/images/new-tenant.png)

假设创建了一个租户1，之后租户1下的用户A登录系统后，将只能看见租户1已设置的 namespace（和系统） 下的资源。同时，用户A所创建的所有资源，从可见性上将属于租户1。当另一个用户B（也属于租户1）登录后，因为它和A同属一个租户，因此用户B能够看到用户A创建的所有资源。当用户C（属于租户2）登录后，将无法看见用户A创建的资源，因为它们属于不同的租户。

## 用户
用户是用于登录平台的，Lizardcd 默认安装好后内置了 `admin` 用户，密码显示在初始化任务 `lizardcd-lizardcd-initjob` 的 Pod 日志中。

1. 现在使用 `admin` 用户登录平台，点击顶部导航栏`配置`按钮，并将页面切换到`用户管理`，可以看到默认有一个 `admin` 用户，该用户属于租户 `admin`，用户角色（权限）也是 `admin`。

![](/images/user-list.png)

2. 点击`新建用户`，输入用户名，选择一个所属租户，以及选择一个角色权限，提交保存即可。

![](/images/new-user.png)

3. 创建用户成功后，该用户的密码将即时显示在页面弹窗中，且仅显示这一次，务必将密码妥善保存。

![](/images/password.png)

4. 使用新用户登录平台，登录成功后 server 会颁发给该用户一个 `JWT Token`，该 `JWT Token` 可用于接口调用，且其所属租户、角色权限和它对应的用户相同。`JWT Token` 的默认过期时间是 1 天，可在 server 的启动配置文件修改，详见 [server配置](/docs/deploy/configure/server)

## 角色权限
Lizardcd 总共有三种角色权限：
- admin：超级管理员权限，可以干任何事，当需要设置多个管理员时使用
- readonly：只读权限，可以访问页面上所有通过 HTTP `GET` 接口获取的数据，例如查看应用列表、查看任务列表、查看 Kubernetes 资源等
- readwrite：读写权限，页面上所有 HTTP 接口数据均可访问，包括 `POST` `PUT` `PATCH`，如创建应用、执行任务、重启 Kubernetes 资源等

::: tip
目前一个用户只能属于一个租户，未来根据使用情况再决定是否优化。
:::