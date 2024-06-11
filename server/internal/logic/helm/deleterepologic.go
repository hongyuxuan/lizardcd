package helm

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRepoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRepoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRepoLogic {
	return &DeleteRepoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRepoLogic) DeleteRepo(req *types.RepoReq) (resp *types.Response, err error) {
	if err = l.svcCtx.Sqlite.Where("name = ?", req.RepoName).Delete(&commontypes.HelmRepositories{}).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	l.Logger.Infof("Delete repository \"%s\" success", req.RepoName)
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "删除成功",
	}
	return resp, nil
}
