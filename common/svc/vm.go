package svc

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/imroc/req/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/ssh"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VmService struct {
	logx.Logger
	ctx context.Context
}

func NewVmService(ctx context.Context) *VmService {
	return &VmService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (s *VmService) createSSHConnection(sshHost, sshUser, sshPassword, sshPrivateKey string, sshPort int) (sshclient *ssh.Client, err error) {
	if sshPort == 0 {
		sshPort = 22
	}
	if sshUser == "" {
		sshUser = "root"
	}
	var authMethod ssh.AuthMethod
	if sshPassword != "" {
		authMethod = ssh.Password(sshPassword)
	} else {
		signer, err := ssh.ParsePrivateKey([]byte(sshPrivateKey))
		if err != nil {
			return nil, err
		}
		authMethod = ssh.PublicKeys(signer)
	}
	config := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshclient, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshHost, sshPort), config)
	return
}

func (s *VmService) SSHDeploy(req *commontypes.SSHDeployReq, target string) (res []byte, err error) {
	var sshclient *ssh.Client
	port, _ := strconv.Atoi(req.SSHPort)
	if sshclient, err = s.createSSHConnection(target, req.SSHUser, req.SSHPassword, req.SSHPrivateKey, port); err != nil {
		s.Logger.Error(err)
		return
	}
	defer sshclient.Close()

	var utf8Output []byte
	// run pre_command
	if req.PreCommand != "" {
		var command string
		if runtime.GOOS != "windows" {
			command = fmt.Sprintf("$SHELL -s <<'EOF'\n%s\nEOF", req.PreCommand)
		}
		if utf8Output, err = s.RunCommand(nil, sshclient, command); err != nil {
			err = errorx.NewDefaultError("Failed to run pre_command: %v, output: %s", err, string(utf8Output))
			s.Logger.Error(err)
			return
		}
		s.Logger.Infof("Successfully run pre_command, output: %s", string(utf8Output))
	}

	// download package
	headers := []string{}
	for k, v := range req.ArtifactHeader {
		headers = append(headers, fmt.Sprintf("--header \"%s:%s\"", k, v))
	}
	command := fmt.Sprintf("cd %s && curl -O %s %s", req.DeployPath, strings.Join(headers, " "), req.ArtifactUrl)
	s.Logger.Debug(command)
	if utf8Output, err = s.RunCommand(nil, sshclient, command); err != nil {
		err = errorx.NewDefaultError("Failed to download package %s: %v", req.ArtifactUrl, err)
		s.Logger.Error(err)
		return
	}
	s.Logger.Infof("Successfully download package [%s] to %s, output: %s", req.ArtifactUrl, req.DeployPath, string(utf8Output))

	// run start_command
	if runtime.GOOS != "windows" {
		command = fmt.Sprintf("$SHELL -s <<'EOF'\n%s\nEOF", req.StartCommand)
	}
	if utf8Output, err = s.RunCommand(nil, sshclient, command); err != nil {
		err = errorx.NewDefaultError("Failed to run start_command: %v, output: %s", err, string(utf8Output))
		s.Logger.Error(err)
		return
	}
	s.Logger.Infof("Successfully run start_command, output: %s", string(utf8Output))

	return utf8Output, nil
}

func (s *VmService) SSHCheck(req *commontypes.SSHDeployReq, target string, healthCheck commontypes.HealthCheck) (res []byte, err error) {
	switch healthCheck.Type {
	case "http":
		return s.HttpCheck(healthCheck.Method, target, healthCheck.Port, healthCheck.Uri)
	case "tcp":
		return s.TcpCheck(target, healthCheck.Port)
	case "shell":
		var sshclient *ssh.Client
		port, _ := strconv.Atoi(req.SSHPort)
		if sshclient, err = s.createSSHConnection(target, req.SSHUser, req.SSHPassword, req.SSHPrivateKey, port); err != nil {
			s.Logger.Error(err)
			return
		}
		var command string
		if runtime.GOOS != "windows" {
			command = fmt.Sprintf("$SHELL -s <<'EOF'\n%s\nEOF", healthCheck.Shell)
		}
		var utf8Output []byte
		if utf8Output, err = s.RunCommand(nil, sshclient, command); err != nil {
			err = errorx.NewDefaultError("Failed to run check_command: %v, output: %s", err, string(utf8Output))
			s.Logger.Error(err)
			return
		}
		s.Logger.Infof("Successfully execute healthcheck \"%s\"", healthCheck.Shell)

		sshclient.Close()
		return utf8Output, nil
	case "none":
		return []byte("No need to healthcheck"), nil
	default:
		return nil, errorx.NewDefaultError("Unsupported healthcheck type: %s", healthCheck.Type)
	}
}

func (s *VmService) RunCommand(cmd *exec.Cmd, sshclient *ssh.Client, command string) (utf8Output []byte, err error) {
	var out bytes.Buffer
	if cmd != nil {
		cmd.Stdout = &out
		cmd.Stderr = &out
		err = cmd.Run()
	} else {
		// 创建会话
		var session *ssh.Session
		if session, err = sshclient.NewSession(); err != nil {
			return nil, err
		}
		defer session.Close()

		session.Stdout = &out
		session.Stderr = &out
		s.Logger.Debug(command)
		err = session.Run(command)
	}
	utf8Output = out.Bytes()
	if runtime.GOOS == "windows" {
		utf8Output, _, _ = transform.Bytes(simplifiedchinese.GBK.NewDecoder(), out.Bytes())
	}
	return
}

func (s *VmService) HttpCheck(method, target, port, uri string) (res []byte, err error) {
	var r *req.Response
	if method == "get" {
		r, err = req.C().SetBaseURL(fmt.Sprintf("http://%s:%s", target, port)).R().Get(uri)
	} else if method == "post" {
		r, err = req.C().SetBaseURL(fmt.Sprintf("http://%s:%s", target, port)).R().Post(uri)
	} else {
		return nil, errorx.NewDefaultError("Unsupported http method: %s", port)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if r.IsError() {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Healthcheck return http code = %d", r.StatusCode))
	}
	return r.Bytes(), nil
}

func (s *VmService) TcpCheck(target, port string) (res []byte, err error) {
	conn, err := net.DialTimeout("tcp", target+":"+port, 2*time.Second)
	if err != nil {
		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Port: %s connection timed out", port))
		} else {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Port: %s connection refused", port))
		}
	} else {
		conn.Close()
		return []byte(fmt.Sprintf("Successfully connect to port = %s", port)), nil
	}
}
