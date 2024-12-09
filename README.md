# Lizardcd - Cloud Native Continuous Delivery for Kubernetes
<p align="center">
<img align="center" width="150px" src="https://project-1255547500.cos.ap-beijing.myqcloud.com/lizardcd%2Flizardcd-logo.jpg">
</p>

Lizardcd is a lightweight cloud native continuous delivery project, which is a server-agent architecture and support multi-cluster of kubernetes. It's composed of lizardcd-server, lizardcd-agent, lizardcd-ui and lizardcd-cli.

Full documentation please reference the official doc site: https://hongyuxuan.github.io/lizardcd-doc/.

# What is Lizardcd
![](https://project-1255547500.cos.ap-beijing.myqcloud.com/lizardcd%2Flizardcd%E6%9E%B6%E6%9E%84%E5%9B%BE.png)
Lizardcd is a cloud native continuous delivery tool for kubernetes. It works in server-agent mode with grpc framework. The agent can run as a `deployment` in one or more kubernetes cluster, or run as an `executed binary package` out of kubernetes with a kubeconfig. The server can run everywhere, as an executed binary package, a docker or a deployment in kubernetes. Lizardcd needs a service registry center for service automatically discovery, we supoorts `etcd`, `consul` and `nacos`, default `etcd`.

# Features
- Cloud native, server-agent architecture, and agent can run in/out of kubernetes cluster.
- No need of kubeconfig and clusterrole, with In-Cluster mode, agent use a serviceaccount to communicate with kubernetes APIServer.
- Support multi-cluster of kubernetes.
- Agent register and discover automatically to service center.
- Create resources(deployments/statefulsets/service and so on) by go-templates.
- Application management, support deploy resources by application[v1.1.0+]. 
- Support `Canary Release` by `Istio`[v1.1.1+], you can configure an application with traffic management using a `DestinationRule` and `VirtualService` by `Istio`, based on `weight` or `http header`.
- Support `helm`, which includes `helm repo`, `helm search`, `helm install`, `helm get`, `helm show` and almost all of helm features.
- Supoort task management, which can be integrated with CI platforms.
- Deploy type supports `Kubernetes`, `Machinery` and `HTTP`, the `HTTP` can integrated with 3rd CD platform like `ansible tower/awx`, `dolphinscheduler` and so on.
- Prometheus Metrics and opentelemetry built-in, support monitoring and tracing.

# How lizardcd works
Lizardcd is composed of lizardcd-agent, lizardcd-server, lizardcd-ui and lizardcd-cli.

- lizardcd-agent: Agent use kubernetes `client-go` SDK to communited with kubernetes APIServer. It listens on a grpc port, and wait of requests from server. The agent will register itself to `etcd` or `consul` or `nacos` with a key of `lizardcd-agent.<namespace>.<cluster>` and a `ServicePrefix` when it starts, please do not change the format of the key. The agent is stateless and can be scaled out with many replicas.

  Attention: if the agent run in a kubernetes pod, it will defaultly register to etcd or consul or nacos with the address of the `pod IP`, which may not be connected from a server, so you can specify a environment variables of `POD_IP`, which points to your agent service `NodePort` or `LoadBalancer`.

- lizardcd-server: Server will `watch` the etcd or consul or nacos services and store all agent service connections in itself, when it received a RESTAPI requests with kubernetes cluster/namespace information, it will take a connection from agent services and communicated with this agent throw grpc.

- lizardcd-ui: A web ui for lizardcd-server, simply supported list/patch/delete/rollout-restart/scale workloads for devops.

- lizardcd-cli: A command-line tool, which has colorful table-format output.

# Why is Lizardcd
For many other CD tools for kubernetes, like argocd, kubevela, which are often agentless, must interactive with kubernetes with a `kubeconfig`. In many cases, you don't have a `kubeconfig`, you don't know the apiserver's address, may be you are using a dashboard of kubernetes, like kubesphere or rancher, and you are just a namespace's user of kubernetes. Lizardcd can run with a `In-Cluster` mode, it is deployed as a `deployment` and use a `serviceaccount` with a role of it's rolebinding. For example, if you want to use lizardcd to set a `image` of a workload, you just need a role with `apiGroups: apps` and `resources: deployments` and `verbs: update`.

Lizardcd supports multi-cluster of kubernetes, one agent for a cluster, if the agent has a cluster scope role of serviceaccount or kubeconfig. Also you can deploy many agents in a cluster and one for a namespace, if you only have a namespace scope role. 

All agents will automatically register to a service center when started, and will automatically deregister when stopped. The server will also automatically discover the agents from the register center.

# Installation
We support Binary package, Docker, Kubernetes(Helm charts), you can get binary packages from release pages, and docker images in [DockerHub](https://hub.docker.com/)

## Kubernetes
We support Helm Charts to install LizardCD to your kubernetes cluster, please refer to https://artifacthub.io/packages/helm/lizardcd/lizardcd

## Linux
Unpack the compressed pacakge:
```shell
tar zxf lizardcd-server-linux-amd64-<version>.tar.gz
tar zxf lizardcd-agent-linux-amd64-<version>.tar.gz
tar zxf lzcli-linux-amd64-<version>.tar.gz
```

Then you can start the agent like this:
```shell
./lizardcd-agent --consul-host 10.50.89.17:8500 --service-key lizardcd-agent.*.tektonk8s --kubeconfig ~/.kube/config --grpc-addr 0.0.0.0:5017
# or
./lizardcd-agent --etcd-host 10.50.89.17:2379 --service-key lizardcd-agent.*.tektonk8s --kubeconfig ~/.kube/config --grpc-addr 0.0.0.0:5017
```

You can start the server like this:
```shell
./lizardcd-server --consul-addr 10.50.89.17:8500 --http-addr=0.0.0.0:5117
# or 
./lizardcd-server --etcd-addr 10.50.89.17:2379 --http-addr=0.0.0.0:5117
```

Then you can use a cli to connect to the server:
```shell
./lzcli config -s http://<lizardcd-server>:5117
./lzcli agent list
```
the config file is defaultly saved to ~/.lizardcd-cli.yaml, you can use `-c` to specified the config file.

## Windows
Unpack the compressed pacakge, and you can start the agent like this:
```powershell
.\lizardcd-agent.exe --consul-host 10.50.89.17:8500 --service-key lizardcd-agent.*.tektonk8s --kubeconfig <path-to-your>\config --grpc-addr 0.0.0.0:5017
# or
.\lizardcd-agent.exe --etcd-host 10.50.89.17:2379 --service-key lizardcd-agent.*.tektonk8s --kubeconfig <path-to-your>\config --grpc-addr 0.0.0.0:5017
```

You can start the server like this:
```powershell
.\lizardcd-server.exe --consul-addr 10.50.89.17:8500 --http-addr=0.0.0.0:5117
# or
.\lizardcd-server.exe --etcd-addr 10.50.89.17:2379 --http-addr=0.0.0.0:5117
```

Then you can use a cli to connect to the server:
```powershell
.\lzcli.exe config -s http://<lizardcd-server>:5117
.\lzcli.exe agent list
```

# Build
The project is open-source and written in golang, you can download the source code and build it as you like. First you must have a golang environment, and go to the root directory and run commands below.

## lizardcd-agent
```shell
cd lizardcd/agent
export BINARY=lizardcd-agent
export GOOS=linux
export GOARCH=amd64
make
```

## lizardcd-server
```shell
cd lizardcd/server
export BINARY=lizardcd-server
export GOOS=linux
export GOARCH=amd64
make
```

## lizardcd-cli
```shell
cd lizardcd/cli
export BINARY=lzcli
export GOOS=linux
export GOARCH=amd64
make
```

## lizardcd-ui
```shell
cd lizardcd/ui
npm run build
```

# Supported versions
Lizardcd is tested in kubernetes-v1.19+, and untested below v1.19, please use lizardcd with caution when your kubernetes version is below v1.19