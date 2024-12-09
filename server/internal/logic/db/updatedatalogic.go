package db

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/repository"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatedataLogic struct {
	logx.Logger
	ctx                   context.Context
	svcCtx                *svc.ServiceContext
	applicationRepository *repository.ApplicationRepository
	ciRepository          *repository.CiRepository
	istioService          *svc.IstioService
}

func NewUpdatedataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatedataLogic {
	return &UpdatedataLogic{
		Logger:                logx.WithContext(ctx),
		ctx:                   ctx,
		svcCtx:                svcCtx,
		applicationRepository: repository.NewApplicationRepository(ctx, svcCtx),
		istioService:          svc.NewIstioService(ctx, svcCtx),
		ciRepository:          repository.NewCiRepository(ctx, svcCtx),
	}
}

func (l *UpdatedataLogic) Updatedata(req *types.UpdateDataReq) (resp *types.Response, err error) {
	_, _, tenant, _ := utils.GetPayload(l.ctx)
	switch req.Tablename {
	case "application":
		_, err = l.applicationRepository.Save(req.Body, tenant[0], false)
	case "application_faas":
		_, err = l.applicationRepository.SaveFaas(req.Body, false)
	case "git_repository":
		_, err = l.ciRepository.SaveGitRepository(req.Body)
	case "ci_trigger":
		_, err = l.ciRepository.SaveCiTrigger(req.Body)
	default:
		err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.UpdateData")).
			Table(req.Tablename).
			Where("id = ?", req.Id).
			Updates(req.Body).Error
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "更新成功",
	}
	return
}
