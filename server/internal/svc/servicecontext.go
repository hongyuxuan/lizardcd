package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	capi "github.com/hashicorp/consul/api"
	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/config"
	"github.com/hongyuxuan/lizardcd/server/internal/middleware"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	nacosconstant "github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
)

type ServiceContext struct {
	Config         config.Config
	AgentList      map[string]*types.RpcAgent
	K8sList        map[string]*types.K8sConnection
	EtcdClient     *clientv3.Client
	ConsulClient   *capi.Client
	NacosClient    naming_client.INamingClient
	Database       *gorm.DB
	Version        string
	Validateuser   rest.Middleware
	Hub            *Hub
	Cron           *cron.Cron
	CronIdMap      map[string]types.CronData
	LeaderElection *LeaderElection
}

type NsCluster struct {
	Namespace string `json:"namespace"`
	Cluster   string `json:"cluster"`
}

func NewServiceContext(c config.Config) *ServiceContext {
	svcCtx := &ServiceContext{
		Config:       c,
		AgentList:    make(map[string]*types.RpcAgent),
		Validateuser: middleware.NewValidateuserMiddleware().Handle,
		Cron:         cron.New(),
		CronIdMap:    make(map[string]types.CronData),
		K8sList:      make(map[string]*types.K8sConnection),
	}
	if c.Database.Type == "sqlite" {
		svcCtx.Database = utils.NewSQLite(c.Database.DBfile, c.Log.Level)
	} else if c.Database.Type == "tidb" {
		svcCtx.Database = utils.NewTidb(c.Database.Host, c.Database.Username, c.Database.Password, c.Database.Database, c.Log.Level, c.Database.Port)
	}
	if c.Etcd.Address != "" {
		etcdHosts := strings.Split(c.Etcd.Address, ",")
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   etcdHosts,
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			logx.Errorf("Failed to connect to etcd: %v", err)
			os.Exit(0)
		}
		logx.Infof("Connect to etcd host=%s success", c.Etcd.Address)
		svcCtx.EtcdClient = client
		le, err := NewLeaderElection(client, c)
		if err != nil {
			logx.Errorf("Failed to elect leader by etcd: %v", err)
			os.Exit(0)
		}
		svcCtx.LeaderElection = le
		le.Start() // start election
	}
	if c.Consul.Address != "" {
		svcCtx.ConsulClient = utils.CreateConsul(c.Consul.Address)
	}
	if c.Nacos.Address != "" {
		fields := strings.Split(c.Nacos.Address, ":")
		host := fields[0]
		port, _ := strconv.ParseUint(fields[1], 10, 64)
		client, err := clients.NewNamingClient(
			vo.NacosClientParam{
				ServerConfigs: []nacosconstant.ServerConfig{
					{
						IpAddr: host,
						Port:   port,
					},
				},
				ClientConfig: &nacosconstant.ClientConfig{
					NamespaceId:         c.Nacos.NamespaceId,
					NotLoadCacheAtStart: true,
					LogDir:              "./logs",
					Username:            c.Nacos.Username,
					Password:            c.Nacos.Password,
				},
			},
		)
		if err != nil {
			logx.Errorf("Failed to connect to nacos: %v", err)
			os.Exit(0)
		}
		logx.Infof("Connect to nacos host=%s success", c.Nacos.Address)
		svcCtx.NacosClient = client
	}
	// start websocket server
	hub := NewHub()
	go hub.Run()
	svcCtx.Hub = hub
	return svcCtx
}

func (s *ServiceContext) GetAgent(cluster, namespace string) (agent lizardagent.LizardAgent, k8sService *commonsvc.K8sService, err error) {
	for k, v := range s.AgentList {
		re, _ := regexp.Compile(k)
		if re.MatchString(fmt.Sprintf("%slizardcd-agent.%s.%s", s.Config.ServicePrefix, namespace, cluster)) {
			return v.Client, nil, nil
		}
	}
	for k, v := range s.K8sList {
		re, _ := regexp.Compile(k)
		if re.MatchString(fmt.Sprintf("%slizardcd-agent.%s.%s", s.Config.ServicePrefix, namespace, cluster)) {
			return nil, v.K8sService, nil
		}
	}
	return nil, nil, errorx.NewDefaultError("Cannot find lizardcd-agent of cluster=%s namespace=%s, maybe the server cannot communicated with the agent", cluster, namespace)
}

func (s *ServiceContext) ListServices(role string, namespaces []string) (services []*commontypes.ServiceMeta) {
	for k, v := range s.AgentList {
		meta, _ := utils.GetServiceMata(s.Config.ServicePrefix, k)
		if meta != nil {
			meta.ServiceSource = v.ServiceSource
			meta.ServiceName = k
			meta.ServiceType = "agent"
			meta.Labels = v.Labels
			services = append(services, meta)
		}
	}
	for k, v := range s.K8sList {
		meta, _ := utils.GetServiceMata(s.Config.ServicePrefix, k)
		if meta != nil {
			meta.Protocol = "https"
			meta.ServiceName = k
			meta.ServiceType = "kubeconfig"
			meta.Labels = v.Labels
			services = append(services, meta)
		}
	}
	return
}

func (s *ServiceContext) GetNamespaces(serviceName string) (res []string) {
	if _, ok := s.AgentList[serviceName]; ok {
		if rpcResponse, err := s.AgentList[serviceName].Client.GetNamespaces(context.Background(), &lizardagent.LabelSelector{LabelSelector: ""}); err != nil {
			return
		} else {
			var r []corev1.Namespace
			json.Unmarshal(rpcResponse.Data, &r)
			for _, ns := range r {
				res = append(res, ns.Name)
			}
		}
	}
	if v, ok := s.K8sList[serviceName]; ok {
		if !v.K8sService.IsValid() {
			logx.Errorf("Cannot connect to k8s: %s", v.RestConfig.Host)
			return
		}
		r, err := v.K8sService.GetNamespaces("")
		if err != nil {
			logx.Error(err)
			return
		}
		for _, ns := range r.Items {
			res = append(res, ns.Name)
		}
	}
	res = lo.Uniq(res)
	return
}

func (s *ServiceContext) GetTargetAgent(ip string) (agent lizardagent.LizardAgent, err error) {
	for k, v := range s.AgentList {
		re, _ := regexp.Compile(fmt.Sprintf("%slizardcd-agent_vm\\.(.+?)\\.%s", s.Config.ServicePrefix, ip))
		if re.MatchString(k) {
			return v.Client, nil
		}
	}
	return nil, errorx.NewDefaultError("Cannot find lizardcd-agent of ip=%s, maybe the server cannot communicated with the agent", ip)
}

func (s *ServiceContext) SetVersion(version string) {
	s.Version = version
}

func (s *ServiceContext) GetHelmSettings(tenant []string) (wait bool, timeout int64, err error) {
	var settings []commontypes.Settings
	if err = s.Database.Where("tenant IN ?", tenant).Find(&settings).Error; err != nil {
		return
	}
	for _, s := range settings {
		if s.SettingKey == "helm_wait" {
			if s.SettingValue == "true" {
				wait = true
			} else {
				wait = false
			}
		}
		if s.SettingKey == "helm_timeout" {
			timeout, _ = strconv.ParseInt(s.SettingValue, 10, 64)
		}
	}
	return
}

func (s *ServiceContext) GetJwtToken(user commontypes.User, expireTime *int64) (accessToken string, now int64, err error) {
	var tenants []commontypes.Tenant
	s.Database.Model(&commontypes.Tenant{}).Where("tenant_name IN ?", strings.Split(user.Tenant, ",")).Find(&tenants)
	var namespaces []string
	for _, t := range tenants {
		var nscluster []NsCluster
		if err = json.Unmarshal([]byte(t.Namespaces), &nscluster); err != nil {
			logx.Error(err)
			return
		}
		namespaces = append(namespaces, lo.Map(nscluster, func(item NsCluster, _ int) string {
			return item.Namespace
		})...)
	}
	payloads := map[string]interface{}{
		"userid":    user.Userid,
		"username":  user.Username,
		"email":     user.Email,
		"role":      user.Role,
		"tenant":    user.Tenant,
		"namespace": strings.Join(namespaces, ","),
	}
	if expireTime == nil {
		expireTime = &s.Config.Auth.AccessExpire
	}
	now = time.Now().Unix()
	if accessToken, err = s.generateToken(now, payloads, *expireTime); err != nil {
		return
	}
	return
}

// 手动注册agent，不通过etcd
func (s *ServiceContext) RegisterAgent(serviceKey, endpoint, proxy, kubeconfig string, labels []string) (err error) {
	if _, ok := s.AgentList[serviceKey]; ok {
		return
	}
	if _, ok := s.K8sList[serviceKey]; ok {
		return
	}
	if endpoint != "" { // agent模式
		if proxy != "" { // 通过代理访问grpc
			var proxyURL *url.URL
			if proxyURL, err = url.Parse(proxy); err != nil {
				return
			}
			proxyClient := &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				},
			}
			cli, err := zrpc.NewClient(zrpc.RpcClientConf{
				Timeout: s.Config.Rpc.Timeout,
				Target:  endpoint,
			}, zrpc.WithDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{ // add keepalive option
				Time:                time.Duration(s.Config.Rpc.KeepaliveTime) * time.Second,
				Timeout:             time.Second,
				PermitWithoutStream: true,
			})), zrpc.WithDialOption(grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
				return proxyClient.Transport.(*http.Transport).DialContext(ctx, "tcp", addr)
			})))
			if err != nil {
				return err
			}
			s.AgentList[serviceKey] = &types.RpcAgent{
				Client:        lizardagent.NewLizardAgent(cli),
				ServiceSource: "manual",
				Cli:           cli,
				Labels:        labels,
				Count:         1,
			}
		} else { // 直连grpc
			cli, err := zrpc.NewClient(zrpc.RpcClientConf{
				Timeout: s.Config.Rpc.Timeout,
				Target:  endpoint,
			}, zrpc.WithDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{ // add keepalive option
				Time:                time.Duration(s.Config.Rpc.KeepaliveTime) * time.Second,
				Timeout:             time.Second,
				PermitWithoutStream: true,
			})))
			if err != nil {
				return err
			}
			s.AgentList[serviceKey] = &types.RpcAgent{
				Client:        lizardagent.NewLizardAgent(cli),
				ServiceSource: "manual",
				Cli:           cli,
				Labels:        labels,
				Count:         1,
			}
		}
		logx.Infof("A new lizardcd-agent: %s registered into lizardcd-server manually", serviceKey)
	} else { // 通过kubeconfig直连K8S
		clientset, restConf, dynamicclient, _, tektonClient, triggerClient, _ := commonsvc.CreateKubernetes("", kubeconfig)
		k8sService := commonsvc.NewK8sService(context.Background(), clientset, dynamicclient, tektonClient, triggerClient)
		k8sService.HelmService.SetKubeconfig(kubeconfig)
		s.K8sList[serviceKey] = &types.K8sConnection{
			K8sService: k8sService,
			RestConfig: restConf,
			Labels:     labels,
		}
		return
	}
	return
}

func (s *ServiceContext) generateToken(iat int64, payloads map[string]interface{}, seconds int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["payloads"] = payloads
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(s.Config.Auth.AccessSecret))
}
