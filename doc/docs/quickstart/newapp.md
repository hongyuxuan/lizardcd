# 新建应用
本章节将快速建立一个简单的demo应用，并对该应用做发布、重启等操作。

## 新建工作负载
Lizardcd默认提供了简易的 `deployment` YAML 模板，点击页面的 `YAML模板` 菜单，该模板遵循 `go-template` 模板语法，如图所示：

![](/images/yaml-template.jpg)

可以打开内置好的模板 `application_template_base`，点击编辑按钮：

![](/images/yaml-template-edit.jpg)

通常情况下不需要修改这个模板，关闭此页面。如有定制需求，可编辑该模板，或者新建您自己的模板。

点击 `工作负载`-`部署` 导航，选择集群 `k8s`、命名空间 `default`，点击新建工作负载，从模板导入选择 `application_template_base`，该模板默认会部署一个 `nginx` 镜像，Appname 为 `lizardtest`，点击提交。

::: warning
您可能无法访问DockerHub，请将 `Image` 修改为能访问的镜像
:::

## 查看工作负载
等待一会，在当前页面看是否 `lizardtest` 工作负载已经就绪：

![](/images/workload-list.jpg)

如果工作负载未就绪，可点击 `lizardtest` 进入工作负载详情页面查看：

![](/images/workload-detail.jpg)

## 添加镜像仓库
点击页面右上角配置，并切换到 `镜像仓库管理`，点击按钮 `新建仓库`，Lizardcd 支持四种类型仓库：
- Jfrog Artifactory
- Harbor
- DockerHub（现在可能已经无法使用了）
- S3

假设选择DockerHub，填入仓库账号和密码（DockerHub请填写Personal Access Tokens），保存即可。

## 新建应用
确保工作负载已经就绪后，点击 `应用管理` 导航，点击 `新建应用` 按钮，新增一个应用：

![](/images/new-application.jpg)

其中，部署方式选择 `容器`，镜像仓库选择上一步新建的 `DockerHub` 仓库，其余按照页面提示填写即可。在工作负载处点击`+`按钮新增一个负载。如果前序步骤部署正确，在容器集群和命名空间应该能出现 agent 所在的 Kubernetes 集群和命名空间。工作负载名称填写 `deployment` 的名称 `lizardtest`，容器名称填写 `spec.template.spec.container.name`，提交保存即可。

![](/images/application-list.jpg)

## 发布版本
在当前页面，找到 `lizardtest` 应用，更多操作按钮选择 `发布`，发布策略保持选择 `所有工作负载使用相同镜像`。如果前面配置镜像仓库没有问题，选择制品这里会列出来所有的镜像 tags，选择一个 tag。工作负载可以选择 `是否启用`，这里保持启用，点击提交。

![](/images/release-app.jpg)

之后页面会自动跳转到 [任务管理](/docs/quickstart/task)。