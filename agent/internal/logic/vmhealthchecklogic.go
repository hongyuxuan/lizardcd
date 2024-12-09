package logic

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/imroc/req/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type VmHealthCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVmHealthCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VmHealthCheckLogic {
	return &VmHealthCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VmHealthCheckLogic) VmHealthCheck(in *agent.VmHealthCheckRequest) (resp *agent.Response, err error) {
	if in.Type == "http" {
		var res *req.Response
		if in.Method == "get" {
			res, err = req.C().SetBaseURL(fmt.Sprintf("http://localhost:%s", in.Port)).R().Get(in.Uri)
		} else if in.Method == "post" {
			res, err = req.C().SetBaseURL(fmt.Sprintf("http://localhost:%s", in.Port)).R().Post(in.Uri)
		} else {
			return nil, errorx.NewDefaultError("Unsupported http method: %s", in.Port)
		}
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if res.IsError() {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Healthcheck return http code = %d", res.StatusCode))
		}
		return &agent.Response{
			Code: uint32(codes.OK),
			Data: res.Bytes(),
		}, nil
	} else if in.Type == "shell" {
		checkfile := fmt.Sprintf("/tmp/tmp-healthcheck-%d.sh", time.Now().UnixMicro())
		if err = os.WriteFile(checkfile, []byte(in.Shell), 0755); err != nil {
			l.Logger.Error(err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		// execute healthcheck
		cmd := exec.Command("/bin/bash", "-c", checkfile)
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		defer os.Remove(checkfile)
		resp = &agent.Response{
			Code: uint32(codes.OK),
			Data: out.Bytes(),
		}
		if err != nil {
			l.Logger.Error(err)
			return resp, status.Error(codes.Internal, err.Error())
		}
		l.Logger.Infof("Successfully execute healthcheck \"%s\"", in.Shell)
		return resp, nil
	} else if in.Type == "tcp" {
		conn, err := net.DialTimeout("tcp", "localhost:"+in.Port, 2*time.Second)
		if err != nil {
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				return nil, status.Error(codes.Internal, fmt.Sprintf("Port: %s connection timed out", in.Port))
			} else {
				return nil, status.Error(codes.Internal, fmt.Sprintf("Port: %s connection refused", in.Port))
			}
		} else {
			conn.Close()
			return &agent.Response{
				Code: uint32(codes.OK),
				Data: []byte(fmt.Sprintf("Successfully connect to port = %s", in.Port)),
			}, nil
		}
	} else {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Unsupported healthcheck type: %s", in.Type))
	}
}
