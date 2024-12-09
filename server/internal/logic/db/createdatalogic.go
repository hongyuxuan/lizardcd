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

type CreatedataLogic struct {
	logx.Logger
	ctx           context.Context
	svcCtx        *svc.ServiceContext
	appRepository *repository.ApplicationRepository
	ciRepository  *repository.CiRepository
}

func NewCreatedataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatedataLogic {
	return &CreatedataLogic{
		Logger:        logx.WithContext(ctx),
		ctx:           ctx,
		svcCtx:        svcCtx,
		appRepository: repository.NewApplicationRepository(ctx, svcCtx),
		ciRepository:  repository.NewCiRepository(ctx, svcCtx),
	}
}

func (l *CreatedataLogic) Createdata(req *types.CreateDataReq) (resp *types.Response, err error) {
	_, _, tenant, _ := utils.GetPayload(l.ctx)
	data := make(map[string]interface{})
	switch req.Tablename {
	case "application":
		data["id"], err = l.appRepository.Save(req.Body, tenant[0], true)
	case "application_faas":
		data["id"], err = l.appRepository.SaveFaas(req.Body, true)
	case "ci_trigger":
		data["id"], err = l.ciRepository.SaveCiTrigger(req.Body)
	case "git_repository":
		data["id"], err = l.ciRepository.SaveGitRepository(req.Body)
	case "tenant":
		if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.CreateTenant")).Table(req.Tablename).Create(&req.Body).Error; err != nil {
			return
		}
		utils.AddSettings(req.Body["tenant_name"].(string), l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.SaveSettings")))
	default:
		l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.CreateData")).Table(req.Tablename).Create(&req.Body)
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "新增成功",
		Data:    data,
	}
	return
}
