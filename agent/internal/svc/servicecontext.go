package svc

import (
	"os"

	"github.com/hongyuxuan/lizardcd/agent/internal/config"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
)

type ServiceContext struct {
	Config        config.Config
	Clientset     *kubernetes.Clientset
	Dynamicclient dynamic.Interface
	Request_k8s   *utils.HttpClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	clientset, dynamicclient, k8sclient := createKubernetes(c)

	if c.Log.Level == "debug" {
		k8sclient.EnableDebug(true)
	}
	return &ServiceContext{
		Config:        c,
		Clientset:     clientset,
		Dynamicclient: dynamicclient,
		Request_k8s:   k8sclient,
	}
}

func createKubernetes(c config.Config) (*kubernetes.Clientset, dynamic.Interface, *utils.HttpClient) {
	var conf *rest.Config
	var err error

	if c.Kubeconfig != "" {
		logx.Infof("Using kubeconfig=%s", c.Kubeconfig)
		conf, err = clientcmd.BuildConfigFromFlags("", c.Kubeconfig)
	} else {
		logx.Info("Using in cluster config")
		conf, err = rest.InClusterConfig()
		conf.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(1000, 1000) // setting a big ratelimiter for client-side throttling, default 5
	}
	if err != nil {
		logx.Error(err)
		os.Exit(0)
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

	var k8sclient *utils.HttpClient = utils.NewHttpClient(otel.Tracer("imroc/req"))
	k8sclient.EnableInsecureSkipVerify().SetBaseURL(conf.Host)
	logx.Infof("Init k8s client %s success", k8sclient.BaseURL)

	return clientset, dynamicclient, k8sclient
}
