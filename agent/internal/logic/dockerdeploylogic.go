package logic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"os/user"
	"strings"

	"github.com/docker/cli/cli/config/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/errdefs"
	"github.com/docker/go-connections/nat"
	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DockerDeployLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDockerDeployLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DockerDeployLogic {
	return &DockerDeployLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// docker deploy
func (l *DockerDeployLogic) DockerDeploy(in *agent.DockerDeployRequest) (*agent.Response, error) {
	// stop
	if err := l.svcCtx.DockerClient.ContainerStop(l.ctx, in.ContainerName, container.StopOptions{}); err != nil {
		if errdefs.IsNotFound(err) {
			l.Logger.Infof("Container=%s not exists, ignore", in.ContainerName)
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else {
		l.Logger.Infof("Container=%s stopped successfully", in.ContainerName)
	}

	// remove
	if err := l.svcCtx.DockerClient.ContainerRemove(l.ctx, in.ContainerName, container.RemoveOptions{}); err != nil {
		if errdefs.IsNotFound(err) {
			l.Logger.Infof("Container=%s not exists, ignore", in.ContainerName)
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else {
		l.Logger.Infof("Container=%s removed successfully", in.ContainerName)
	}

	// start
	containerConfig := &container.Config{
		Image: in.Image,
	}
	hostConfig := &container.HostConfig{}
	if in.Network != "" { // --network
		hostConfig.NetworkMode = container.NetworkMode(in.Network)
	}
	if in.Ports != "" { // -p
		containerConfig.ExposedPorts = nat.PortSet{}
		hostConfig.PortBindings = nat.PortMap{}
		for _, portBinding := range strings.Split(in.Ports, ",") { // 8080:80,9090:90
			portKV := strings.Split(portBinding, ":")
			containerPort := portKV[1] + "/tcp"
			containerConfig.ExposedPorts[nat.Port(containerPort)] = struct{}{}
			hostConfig.PortBindings[nat.Port(containerPort)] = []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: portKV[0],
				},
			}
		}
	}
	if in.Volumes != "" { // -v
		hostConfig.Binds = strings.Split(in.Volumes, ",")
	}
	if in.Dns != "" { // --dns
		hostConfig.DNS = strings.Split(in.Dns, ",")
	}
	if in.WorkingDir != "" { // -w
		containerConfig.WorkingDir = in.WorkingDir
	}
	if len(in.Command) != 0 {
		containerConfig.Cmd = in.Command
	}

	// docker pull
	auth, err := getDockerCredential()
	if err != nil {
		l.Logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	l.Logger.Infof("Get docker auth from ~/.docker/config.json: %s", *auth)
	out, err := l.svcCtx.DockerClient.ImagePull(l.ctx, in.Image, image.PullOptions{RegistryAuth: *auth})
	if err != nil {
		l.Logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer out.Close()
	if _, err := io.Copy(io.Discard, out); err != nil {
		l.Logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	l.Logger.Infof("Image=%s pulled successfully", in.Image)
	// docker create
	resp, err := l.svcCtx.DockerClient.ContainerCreate(l.ctx, containerConfig, hostConfig, &network.NetworkingConfig{}, nil, in.ContainerName)
	if err != nil {
		l.Logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	// docker start
	if err := l.svcCtx.DockerClient.ContainerStart(l.ctx, resp.ID, container.StartOptions{}); err != nil {
		l.Logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	l.Logger.Infof("Container=%s started successfully", in.ContainerName)

	return &agent.Response{
		Code: uint32(codes.OK),
		Data: []byte(resp.ID),
	}, nil
}

func getDockerCredential() (*string, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	configFile, err := os.ReadFile(usr.HomeDir + "/.docker/config.json")
	if err != nil {
		return nil, errorx.NewDefaultError("Failed to read Docker config file: %v", err)
	}
	// parse docker config.json
	var config struct {
		Auths map[string]struct {
			Auth string `json:"auth"`
		} `json:"auths"`
	}
	if err := json.Unmarshal(configFile, &config); err != nil {
		return nil, errorx.NewDefaultError("Failed to parse Docker config file: %v", err)
	}
	authStr := ""
	for _, auth := range config.Auths {
		decodedAuth, _ := base64.StdEncoding.DecodeString(auth.Auth)
		parts := strings.SplitN(string(decodedAuth), ":", 2)
		if len(parts) != 2 {
			return nil, errorx.NewDefaultError("Invalid auth format: expected username:password")
		}
		authConfig := types.AuthConfig{
			Username: parts[0],
			Password: parts[1],
		}
		encoded, _ := json.Marshal(authConfig)
		authStr = base64.URLEncoding.EncodeToString(encoded)
		break
	}
	return &authStr, nil
}
