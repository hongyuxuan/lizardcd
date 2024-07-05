package utils

import (
	"context"
	"os"

	"github.com/hongyuxuan/lizardcd/common/errorx"

	capi "github.com/hashicorp/consul/api"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConsulUtil struct {
	logx.Logger
	consulClient *capi.Client
}

func NewConsulUtil(ctx context.Context, consulClient *capi.Client) *ConsulUtil {
	return &ConsulUtil{
		Logger:       logx.WithContext(ctx),
		consulClient: consulClient,
	}
}

func CreateConsul(address string) *capi.Client {
	client, err := capi.NewClient(&capi.Config{
		Address: address,
	})
	if err != nil {
		logx.Errorf("Error creating consul client: %v", err)
		os.Exit(0)
	}
	logx.Infof("Connect to consul host=%s success", address)
	return client
}

func (c *ConsulUtil) GetKV(ctx context.Context, key string) ([]byte, error) {
	pair, _, err := c.consulClient.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		e := errorx.NewDefaultError("Consul cannot find key: %s", key)
		return nil, e
	}
	c.Logger.Debugf("Consul key[%s]=%s", pair.Key, string(pair.Value))
	return pair.Value, nil
}

func (c *ConsulUtil) ListServices() (services map[string][]string, err error) {
	services, _, err = c.consulClient.Catalog().Services(&capi.QueryOptions{})
	return
}

func (c *ConsulUtil) GetService(serviceName string) (service []*capi.CatalogService, err error) {
	service, _, err = c.consulClient.Catalog().Service(serviceName, "", &capi.QueryOptions{})
	return
}
