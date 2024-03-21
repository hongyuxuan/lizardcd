// Code generated by goctl. DO NOT EDIT!
// Source: agent.proto

package server

import (
	"context"

	"github.com/hongyuxuan/lizardcd/agent/internal/logic"
	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
)

type LizardAgentServer struct {
	svcCtx *svc.ServiceContext
	agent.UnimplementedLizardAgentServer
}

func NewLizardAgentServer(svcCtx *svc.ServiceContext) *LizardAgentServer {
	return &LizardAgentServer{
		svcCtx: svcCtx,
	}
}

func (s *LizardAgentServer) PatchDeployment(ctx context.Context, in *agent.PatchWorkloadRequest) (*agent.Response, error) {
	l := logic.NewPatchDeploymentLogic(ctx, s.svcCtx)
	return l.PatchDeployment(in)
}

func (s *LizardAgentServer) PatchStatefulset(ctx context.Context, in *agent.PatchWorkloadRequest) (*agent.Response, error) {
	l := logic.NewPatchStatefulsetLogic(ctx, s.svcCtx)
	return l.PatchStatefulset(in)
}

func (s *LizardAgentServer) ListDeployment(ctx context.Context, in *agent.ListWorkloadRequest) (*agent.Response, error) {
	l := logic.NewListDeploymentLogic(ctx, s.svcCtx)
	return l.ListDeployment(in)
}

func (s *LizardAgentServer) ListStatefulset(ctx context.Context, in *agent.ListWorkloadRequest) (*agent.Response, error) {
	l := logic.NewListStatefulsetLogic(ctx, s.svcCtx)
	return l.ListStatefulset(in)
}

func (s *LizardAgentServer) GetDeploymentPod(ctx context.Context, in *agent.GetWorkloadRequest) (*agent.Response, error) {
	l := logic.NewGetDeploymentPodLogic(ctx, s.svcCtx)
	return l.GetDeploymentPod(in)
}

func (s *LizardAgentServer) GetStatefulsetPod(ctx context.Context, in *agent.GetWorkloadRequest) (*agent.Response, error) {
	l := logic.NewGetStatefulsetPodLogic(ctx, s.svcCtx)
	return l.GetStatefulsetPod(in)
}

func (s *LizardAgentServer) GetPodEvent(ctx context.Context, in *agent.GetPodEventRequest) (*agent.Response, error) {
	l := logic.NewGetPodEventLogic(ctx, s.svcCtx)
	return l.GetPodEvent(in)
}

func (s *LizardAgentServer) DeleteYaml(ctx context.Context, in *agent.YamlRequest) (*agent.Response, error) {
	l := logic.NewDeleteYamlLogic(ctx, s.svcCtx)
	return l.DeleteYaml(in)
}

func (s *LizardAgentServer) ApplyYaml(ctx context.Context, in *agent.YamlRequest) (*agent.Response, error) {
	l := logic.NewApplyYamlLogic(ctx, s.svcCtx)
	return l.ApplyYaml(in)
}

func (s *LizardAgentServer) Getyaml(ctx context.Context, in *agent.GetYamlRequest) (*agent.YamlResponse, error) {
	l := logic.NewGetyamlLogic(ctx, s.svcCtx)
	return l.Getyaml(in)
}

func (s *LizardAgentServer) RolloutDeployment(ctx context.Context, in *agent.RolloutWorkloadRequest) (*agent.Response, error) {
	l := logic.NewRolloutDeploymentLogic(ctx, s.svcCtx)
	return l.RolloutDeployment(in)
}

func (s *LizardAgentServer) RolloutStatefulset(ctx context.Context, in *agent.RolloutWorkloadRequest) (*agent.Response, error) {
	l := logic.NewRolloutStatefulsetLogic(ctx, s.svcCtx)
	return l.RolloutStatefulset(in)
}

func (s *LizardAgentServer) ScaleDeployment(ctx context.Context, in *agent.ScaleRequest) (*agent.Response, error) {
	l := logic.NewScaleDeploymentLogic(ctx, s.svcCtx)
	return l.ScaleDeployment(in)
}

func (s *LizardAgentServer) ScaleStatefulset(ctx context.Context, in *agent.ScaleRequest) (*agent.Response, error) {
	l := logic.NewScaleStatefulsetLogic(ctx, s.svcCtx)
	return l.ScaleStatefulset(in)
}

func (s *LizardAgentServer) GetNamespaces(ctx context.Context, in *agent.LabelSelector) (*agent.Response, error) {
	l := logic.NewGetNamespacesLogic(ctx, s.svcCtx)
	return l.GetNamespaces(in)
}