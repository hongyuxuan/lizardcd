package svc

import (
	"os"
	"strings"
	"time"

	"github.com/hongyuxuan/lizardcd/agent/internal/config"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel"

	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
)

type ServiceContext struct {
	Config        config.Config
	EtcdClient    *clientv3.Client
	Clientset     *kubernetes.Clientset
	Dynamicclient dynamic.Interface
	Request_k8s   *utils.HttpClient
	Istioclient   *versionedclient.Clientset
}

func NewServiceContext(c config.Config) *ServiceContext {
	clientset, dynamicclient, k8sclient, istioclient := createKubernetes(c)
	if c.Log.Level == "debug" {
		k8sclient.EnableDebug(true)
	}

	var client *clientv3.Client
	if len(c.Etcd.Hosts) > 0 {
		var err error
		client, err = clientv3.New(clientv3.Config{
			Endpoints:   c.Etcd.Hosts,
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			logx.Errorf("Failed to connect to etcd: %v", err)
			os.Exit(0)
		}
		logx.Infof("Connect to etcd host=%s success", strings.Join(c.Etcd.Hosts, ","))
	}

	return &ServiceContext{
		Config:        c,
		Clientset:     clientset,
		Dynamicclient: dynamicclient,
		Request_k8s:   k8sclient,
		EtcdClient:    client,
		Istioclient:   istioclient,
	}
}

func createKubernetes(c config.Config) (*kubernetes.Clientset, dynamic.Interface, *utils.HttpClient, *versionedclient.Clientset) {
	var conf *rest.Config
	var err error
	if c.Kubeconfig != "" {
		logx.Infof("Using kubeconfig=%s", c.Kubeconfig)
		conf, err = clientcmd.BuildConfigFromFlags("", c.Kubeconfig)
	} else {
		logx.Info("Using in cluster config")
		conf, err = rest.InClusterConfig()
		if conf != nil {
			conf.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(1000, 1000) // setting a big ratelimiter for client-side throttling, default 5
		}
	}
	if err != nil {
		logx.Errorf("Cannot build connection to kubernetes with kubeconfig or in-cluster, err: %v. Will start as an vm agent.", err)
		return nil, nil, nil, nil
	}
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		logx.Error(err)
		os.Exit(0)
	}
	dynamicclient, err := dynamic.NewForConfig(conf)
	if err != nil {
		logx.Error(err)
		os.Exit(0)
	}

	// create k8s httpclient
	var k8sclient *utils.HttpClient = utils.NewHttpClient(otel.Tracer("imroc/req"))
	k8sclient.EnableInsecureSkipVerify().SetBaseURL(conf.Host)
	logx.Infof("Init k8s client %s success", k8sclient.BaseURL)

	// create istio client if posible
	istioclient, err := versionedclient.NewForConfig(conf)
	if err != nil {
		logx.Errorf("Failed to create istio client: %s", err)
	}

	return clientset, dynamicclient, k8sclient, istioclient
}
