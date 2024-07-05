package vm

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeployLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeployLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeployLogic {
	return &DeployLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeployLogic) Deploy(req *types.VmDeployReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	result := make(map[string]string)
	for _, ip := range req.Targets {
		if ag, err = l.svcCtx.GetTargetAgent(ip); err != nil {
			return
		}
		var rpcResponse *agent.Response
		header, _ := json.Marshal(req.ArtifactHeader)
		if rpcResponse, err = ag.VmDeploy(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.VmDeploy"), &agent.VmDeployRequest{
			ArtifactUrl:    req.ArtifactUrl,
			ArtifactHeader: header,
			DeployPath:     req.DeployPath,
			DeployUser:     req.DeployUser,
			PreCommand:     req.PreCommand,
			StartCommand:   req.StartCommand,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		result[ip] = string(rpcResponse.Data)
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "任务提交成功",
		Data:    result,
	}
	return
}
