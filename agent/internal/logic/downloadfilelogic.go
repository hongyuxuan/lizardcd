package logic

import (
	"context"
	"io"
	"os"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDownloadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadFileLogic {
	return &DownloadFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DownloadFileLogic) DownloadFile(in *agent.FileRequest, stream agent.LizardAgent_DownloadFileServer) (err error) {
	file, err := os.Open(in.Filename)
	if err != nil {
		l.Logger.Error(err)
		return status.Error(codes.Internal, err.Error())
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			l.Logger.Error(err)
			return status.Error(codes.Internal, err.Error())
		}
		stream.Send(&agent.Response{Data: buffer[:n]})
	}
	return nil
}
