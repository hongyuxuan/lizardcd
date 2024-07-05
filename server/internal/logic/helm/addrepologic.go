package helm

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRepoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddRepoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRepoLogic {
	return &AddRepoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddRepoLogic) AddRepo(req *types.AddRepoReq) (resp *types.Response, err error) {
	_, _, tenant, _ := utils.GetPayload(l.ctx)
	var repo *commontypes.HelmRepositories
	if err = l.svcCtx.Sqlite.Where("name = ?", req.Name).First(&repo).Error; err == nil { // err is nil means find one
		err = errorx.NewDefaultError("Repository \"%s\" has exists", req.Name)
		l.Logger.Error(err)
		return
	}
	if err = l.svcCtx.Sqlite.Create(&commontypes.HelmRepositories{
		Name:     req.Name,
		URL:      req.Url,
		Username: req.Username,
		Password: req.Password,
		Tenant:   tenant,
	}).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	l.Logger.Infof("Add repository \"%s\" to database success", req.Name)
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "添加成功",
	}
	return resp, nil
}
