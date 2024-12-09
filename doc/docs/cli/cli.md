# Lizardcd-cli 用法
<!-- ## 安装
::: code-group
```sh [window]
下载压缩包`lzcli-linux-amd64-<version>.tar.gz`并解压到任意位置
```
```sh [linux]
下载压缩包`lzcli-linux-amd64-<version>.tar.gz`并解压到任意位置
将 lzcli 添加到路径：
export PATH=$PWD:$PATH
```
```sh [macOS]
$ 
```
::: -->

使用以下命令可以查看lzcli支持的操作的完整列表:
![](/images/lzcli.jpg)
## lzcli命令行详解
### `lzcli config`
用于配置lzcli配置文件
#### 用法
```sh
lzcli config [flags]
```
#### 选项 {#options}

| 选项            | 说明                                       |
| --------------- | ------------------------------------------ |
| `-s` | 指定lizardcd-server地址 例：`http://192.168.200.100:5117`    |
| `-c` | 指定lzcli配置文件所在路径 不加默认写入~/.lizardcd-cli.yaml    |


### `lzcli login`
登录lizardcd-server服务器
#### 用法
```sh
lzcli login [flags]
```
#### 选项 

| 选项            | 说明                                       |
| --------------- | ------------------------------------------ |
| `-u` | 指定lizardcd-server的用户名     |
| `-p` | 指定lizardcd-server的密码     |
| `--save` | 保存用户名和密码到lzcli配置文件     |

### `lzcli agent`
获取lizardcd-agent列表
#### 用法
```sh
lzcli agent [command]
```
#### 子命令

| 子命令            | 语法 | 描述                                       |
| --------------- | --------------- | ------------------------------------------ |
| `list` | lzcli agent list | 获取agent列表     |

### `lzcli application`
用于管理应用
#### 用法
```sh
lzcli application [command] [flags]
```
#### 子命令 

| 子命令            | 语法 | 描述                                       |
| --------------- | --------------- | ------------------------------------------ |
| `deploy` | lzcli application deploy --name NAME --image IMAGE --cluster NAME -n NAMESPACE | 发布应用，只支持所有工作负载使用相同镜像    |
| `list` | lzcli application list --cluster NAME -n NAMESPACE | 获取应用列表     |
| `restart` | lzcli application restart --name NAME --cluster NAME -n NAMESPACE | 重启指定应用的所有工作负载     |


##### 选项

| 选项            | 说明                                       |
| --------------- | ------------------------------------------ |
| `--name`| 指定要发布的应用名称 |
| `--image`| 指定应用所使用的的镜像 |
| `--cluster` | 指定应用所在集群    |
| `-n/--namespace` | 指定应用的k8s命名空间，默认default空间     |

### `lzcli apply`
使用一个go模板部署工作负载
#### 用法
```sh
lzcli apply -f go-template.yaml -v KEY=VALUE --cluster NAME --namespace NAMESPACE 
```
#### 选项 

| 选项            | 说明                                       |
| --------------- | ------------------------------------------ |
| `--cluster` | 指定部署的应用所在集群    |
| `-f/--file` | 指定部署应用所使用的go-template模板     |
| `-n/--namespace` | 指定部署应用的k8s命名空间，默认default空间     |
| `-v/--vars` | go-template模板的变量，格式为：Appname=test,Namespace=default,Port=3000 |

### `lzcli deployment`
管理工作负载
#### 用法
```sh
lzcli deployment [command] [flags]
```
#### 子命令 

| 子命令            | 语法 | 描述                                       |
| --------------- | --------------- | ------------------------------------------ |
| `list` | lzcli deployment list --cluster NAME -n NAMESPACE | 获取指定命名空间下所有工作负载    |
| `restart` | lzcli deployment restart --name NAME --cluster NAME -n NAMESPACE | 重启指定命名空间下的工作负载     |
| `scale` | lzcli deployment scale --name NAME --replicas NUM --cluster NAME -n NAMESPACE | 设置工作负载副本数     |
| `set` | lzcli deployment set --name NAME --container CONTAINER --image IMAGE --cluster NAME -n NAMESPACE | 配置工作负载镜像 |
| `show` | lzcli deployment show TYPE --cluster NAME -n NAMESPACE | 获取指定工作负载下的资源信息，支持的资源为：pod、event、container |
##### 选项

| 选项            | 说明                                       |
| --------------- | ------------------------------------------ |
| `--name`| 指定工作负载名称 |
| `--replicas`| 指定副本数 |
| `--container`| 指定容器名称 |
| `--image` | 指定要更新的工作镜像 |
| `--cluster` | 指定工作负载所在的集群    |
| `-n/--namespace` | 指定工作负载所在的k8s命名空间，默认default空间     |

### `lzcli statefulset`
管理有状态副本集
#### 用法
```sh
lzcli statefulset [command] [flags]
```
#### 子命令 

| 子命令            | 语法 | 描述                                       |
| --------------- | --------------- | ------------------------------------------ |
| `list` |  lzcli statefulset list --cluster NAME -n NAMESPACE | 获取指定命名空间下所有有状态副本集    |
| `show` |  lzcli statefulset show TYPE NAME --cluster NAME -n NAMESPACE | 获取指定有状态副本集下的资源信息，支持的资源为：pod、event、container |
##### 选项

| 选项            | 说明                                       |
| --------------- | ------------------------------------------ |
| `--cluster`| 指定有状态副本集所在集群 |
| `-n/--namespace` | 指定有状态副本集所在的k8s命名空间，默认default空间     |

### `lzcli task`
查看任务信息
#### 用法
```sh
lzcli task [command] [flags]
```
#### 子命令 

| 子命令            | 语法 | 描述                                       |
| --------------- | --------------- | ------------------------------------------ |
| `list` |  lzcli task list  | 获取所有任务列表    |
| `show` |  lzcli task show --id TASK_ID | 获取指定任务详情 |
##### 选项

| 选项            | 说明                                       |
| --------------- | ------------------------------------------ |
| `--id`| 指定任务id |

### `lzcli helm`
管理helm
#### 用法
```sh
lzcli helm [command] [flags]
```
#### 子命令 

| 子命令            | 语法 | 描述                                       |
| --------------- | --------------- | ------------------------------------------ |
| `install` |  lzcli helm install [RELEASE_NAME] [REPO_NAME/CHART_NAME] --version VERSION   | 安装chart，--version必须指定    |
| `pull` |  lzcli helm pull [REPO_NAME/CHART_NAME] --version VERSION | 下载指定版本的charts |
| `release` | lzcli helm release list --cluster NAME --namespace NAMESPACE | 获取指定命名空间下已发布版本 |
| `repo` | lzcli helm repo [ add \| list \| remove \| update ] --name NAME --url URL | 添加、列出、删除、更新chart仓库 |
| `rollback` | lzcli helm rollback [RELEASE_NAME] [REVISION] | 回滚发布到上一个版本 |
| `search` | lzcli helm search [REPO_NAME/CHART_NAME] | 在已添加的chart仓库中搜索helm chart |
| `show` | lzcli helm show  [ values \| readme ] --chart NAME --repo REPO_NAME --version VERSION | 显示chart的README或values |
| `uninstall` | lzcli helm uninstall [RELEASE_NAME] | 卸载已发布版本 |
| `upgrade` | lzcli helm upgrade [RELEASE_NAME] [REPO_NAME/CHART_NAME] --file values.yaml --version VERSION | 发布升级到新版的chart |
##### 选项

| 选项            | 说明                                       |
| --------------- | ------------------------------------------ |
| `--version`| 指定chart版本号 |
| `--name`| 指定chart仓库名称 |
| `--url`| 指定chart仓库url地址 |
| `--chart`| 指定chart名称 |
| `--repo`| 已添加的chart仓库名称 |
| `--file/-f` | 指定新版chart的valus.yaml文件 |
| `--cluster`| 指定helm操作所在的集群 |
| `--namespace` | 指定helm操作所在的命名空间 | 