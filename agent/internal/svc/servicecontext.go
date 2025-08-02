package svc

import (
	"os"
	"strings"
	"time"

	"github.com/hongyuxuan/lizardcd/agent/internal/config"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"

	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	triggerclient "github.com/tektoncd/triggers/pkg/client/clientset/versioned"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	dockerclient "github.com/docker/docker/client"
)

type ServiceContext struct {
	Config        config.Config
	EtcdClient    *clientv3.Client
	Clientset     *kubernetes.Clientset
	Dynamicclient dynamic.Interface
	Token         string
	Istioclient   *versionedclient.Clientset
	TektonClient  *tektonclient.Clientset
	TriggerClient *triggerclient.Clientset
	DockerClient  *dockerclient.Client
	RestConfig    *rest.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	clientset, restConf, dynamicclient, istioclient, tektonClient, triggerClient, token := commonsvc.CreateKubernetes(c.Kubeconfig, "")

	var etcdClient *clientv3.Client
	if len(c.Etcd.Hosts) > 0 {
		var err error
		etcdClient, err = clientv3.New(clientv3.Config{
			Endpoints:   c.Etcd.Hosts,
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			logx.Errorf("Failed to connect to etcd: %v", err)
			os.Exit(0)
		}
		logx.Infof("Connect to etcd host=%s success", strings.Join(c.Etcd.Hosts, ","))
	}

	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithAPIVersionNegotiation())
	if err == nil {
		logx.Infof("Connect to docker daemon success")
	}

	return &ServiceContext{
		Config:        c,
		Clientset:     clientset,
		Dynamicclient: dynamicclient,
		Token:         token,
		EtcdClient:    etcdClient,
		Istioclient:   istioclient,
		TektonClient:  tektonClient,
		TriggerClient: triggerClient,
		DockerClient:  cli,
		RestConfig:    restConf,
	}
}
