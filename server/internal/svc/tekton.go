package svc

import (
	"context"
	"encoding/json"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	tektontypes "github.com/hongyuxuan/tekton-sdk-go/types"
	"github.com/samber/lo"
	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"github.com/zeromicro/go-zero/core/logx"
)

type TektonService struct {
	logx.Logger
	ctx    context.Context
	svcCtx *ServiceContext
}

func NewTektonService(ctx context.Context, svcCtx *ServiceContext) *TektonService {
	return &TektonService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (t *TektonService) CancelPipelineRun(ag lizardagent.LizardAgent, ks *commonsvc.K8sService, namespace, labels string) {
	r := struct {
		Results []tektonv1.PipelineRun `json:"results"`
	}{Results: []tektonv1.PipelineRun{}}
	var data []byte
	var err error
	if ag != nil {
		rpcResponse, err := ag.ListTektonResource(context.Background(), &lizardagent.TektonListRequest{
			Namespace:     namespace,
			ResourceType:  "pipelineruns",
			LabelSelector: labels,
		})
		if err != nil {
			t.Logger.Errorf("Failed to list PipelineRuns when cancelPipelinerun: %v", err)
			return
		}
		data = rpcResponse.Data
	} else if ks != nil && ks.IsValid() {
		if data, err = ks.TektonService.ListResource(namespace, "pipelineruns", labels, "", "", 0); err != nil {
			t.Logger.Error(err)
			return
		}
	} else {
		t.Logger.Errorf("Cannot ListTektonResource of namespace=%s labels=%s", namespace, labels)
		return
	}
	json.Unmarshal(data, &r)
	if len(r.Results) == 0 {
		t.Logger.Infof("No PipelineRuns need to be cancelled because there's no PipelineRun with labels: %s", labels)
		return
	}
	runningInstances := lo.Filter(r.Results, func(item tektonv1.PipelineRun, _ int) bool {
		if len(item.Status.Conditions) == 0 {
			return false
		}
		return item.Status.Conditions[0].Status == "Unknown" && item.Status.Conditions[0].Reason == "Running"
	})
	for _, instance := range runningInstances {
		t.Logger.Infof("Find running PipelineRun: %s with labels: %s, will cancelled it", instance.GetName(), labels)
		b, _ := json.Marshal([]tektontypes.PatchOptions{
			{
				Op:    "replace",
				Path:  "/spec/status",
				Value: "Cancelled",
			},
		})
		if ag != nil {
			if _, err = ag.PatchTektonResource(context.Background(), &lizardagent.TektonPatchRequest{
				Namespace:    namespace,
				ResourceType: "pipelineruns",
				ResourceName: instance.GetName(),
				PatchOptions: b,
			}); err != nil {
				t.Logger.Errorf("Error in cancelled PipelineRun: %v", err)
			}
		} else if ks != nil && ks.IsValid() {
			if _, err = ks.TektonService.PatchResource(namespace, "pipelineruns", instance.GetName(), b); err != nil {
				t.Logger.Error(err)
			}
		} else {
			t.Logger.Errorf("Cannot PatchTektonResource of namespace=%s pipelineruns=%s", namespace, instance.GetName())
		}
	}
}
