package svc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hongyuxuan/lizardcd/common/constant"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	triggerbeta1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1beta1"
	triggerclient "github.com/tektoncd/triggers/pkg/client/clientset/versioned"
	"github.com/zeromicro/go-zero/core/logx"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ktypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/dynamic"
)

type TektonService struct {
	logx.Logger
	ctx           context.Context
	clientset     *tektonclient.Clientset
	triggerset    *triggerclient.Clientset
	dynamicclient dynamic.Interface
}

func NewTektonService(ctx context.Context, clientset *tektonclient.Clientset, triggerset *triggerclient.Clientset, dynamicclient dynamic.Interface) *TektonService {
	return &TektonService{
		Logger:        logx.WithContext(ctx),
		ctx:           ctx,
		clientset:     clientset,
		triggerset:    triggerset,
		dynamicclient: dynamicclient,
	}
}

func (t *TektonService) IsValid() bool {
	return t.clientset != nil && t.triggerset != nil
}

func (t *TektonService) ListResource(namespace, resourceType, labelSelector, fieldSelector, conti string, limit int64) (res []byte, err error) {
	if limit == 0 {
		limit = 500
	}
	switch resourceType {
	case constant.TEKTON_CRD_TYPE_TASKS:
		var r *tektonv1.TaskList
		if r, err = t.clientset.TektonV1().Tasks(namespace).List(t.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListTasks: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
			delete(r.Items[i].Annotations, "kubectl.kubernetes.io/last-applied-configuration")
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	case constant.TEKTON_CRD_TYPE_PIPELINES:
		var r *tektonv1.PipelineList
		if r, err = t.clientset.TektonV1().Pipelines(namespace).List(t.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListPipelines: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
			delete(r.Items[i].Annotations, "kubectl.kubernetes.io/last-applied-configuration")
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	case constant.TEKTON_CRD_TYPE_PIPELINERUNS:
		var r *tektonv1.PipelineRunList
		if r, err = t.clientset.TektonV1().PipelineRuns(namespace).List(t.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListPipelineRuns: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
			delete(r.Items[i].Annotations, "kubectl.kubernetes.io/last-applied-configuration")
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	case constant.TEKTON_CRD_TYPE_TASKRUNS:
		var r *tektonv1.TaskRunList
		if r, err = t.clientset.TektonV1().TaskRuns(namespace).List(t.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListTaskRuns: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
			delete(r.Items[i].Annotations, "kubectl.kubernetes.io/last-applied-configuration")
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	case constant.TEKTON_CRD_TYPE_TRIGGERBINDINGS:
		var r *triggerbeta1.TriggerBindingList
		if r, err = t.triggerset.TriggersV1beta1().TriggerBindings(namespace).List(t.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListTriggerBindings: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
			delete(r.Items[i].Annotations, "kubectl.kubernetes.io/last-applied-configuration")
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	case constant.TEKTON_CRD_TYPE_TRIGGERTEMPLATES:
		var r *triggerbeta1.TriggerTemplateList
		if r, err = t.triggerset.TriggersV1beta1().TriggerTemplates(namespace).List(t.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListTriggerTemplates: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
			delete(r.Items[i].Annotations, "kubectl.kubernetes.io/last-applied-configuration")
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	case constant.TEKTON_CRD_TYPE_EVENTLISTENERS:
		var r *triggerbeta1.EventListenerList
		if r, err = t.triggerset.TriggersV1beta1().EventListeners(namespace).List(t.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListEventListeners: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
			delete(r.Items[i].Annotations, "kubectl.kubernetes.io/last-applied-configuration")
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	case constant.TEKTON_CRD_TYPE_APPROVALTASKS:
		gvr := schema.GroupVersionResource{
			Group:    "tekton.automatiko.io",
			Version:  "v1beta1",
			Resource: "approvaltasks",
		}
		var r *unstructured.UnstructuredList
		if r, err = t.dynamicclient.Resource(gvr).Namespace(namespace).List(t.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListApprovalTasks: %w", err)
		}
		for i := range r.Items {
			r.Items[i].SetManagedFields(nil)
			annotations := r.Items[i].GetAnnotations()
			delete(annotations, "kubectl.kubernetes.io/last-applied-configuration")
			r.Items[i].SetAnnotations(annotations)
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.GetContinue(),
		})
	default:
		return nil, errorx.NewDefaultError("Unknown resource type")
	}
	return
}

func (t *TektonService) GetResource(namespace, resourceType, resourceName string) (res []byte, err error) {
	switch resourceType {
	case constant.TEKTON_CRD_TYPE_TASKS:
		r, err := t.clientset.TektonV1().Tasks(namespace).Get(t.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("error in GetTasks: %w", err)
		}
		r.ManagedFields = nil
		res, _ = json.Marshal(r)
	case constant.TEKTON_CRD_TYPE_PIPELINES:
		r, err := t.clientset.TektonV1().Pipelines(namespace).Get(t.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("error in GetPipelines: %w", err)
		}
		r.ManagedFields = nil
		res, _ = json.Marshal(r)
	case constant.TEKTON_CRD_TYPE_PIPELINERUNS:
		r, err := t.clientset.TektonV1().PipelineRuns(namespace).Get(t.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("error in GetPipelineRuns: %w", err)
		}
		r.ManagedFields = nil
		res, _ = json.Marshal(r)
	case constant.TEKTON_CRD_TYPE_TASKRUNS:
		r, err := t.clientset.TektonV1().TaskRuns(namespace).Get(t.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("error in GetTaskRuns: %w", err)
		}
		r.ManagedFields = nil
		res, _ = json.Marshal(r)
	case constant.TEKTON_CRD_TYPE_TRIGGERBINDINGS:
		r, err := t.triggerset.TriggersV1beta1().TriggerBindings(namespace).Get(t.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("error in GetTriggerBindings: %w", err)
		}
		r.ManagedFields = nil
		res, _ = json.Marshal(r)
	case constant.TEKTON_CRD_TYPE_TRIGGERTEMPLATES:
		r, err := t.triggerset.TriggersV1beta1().TriggerTemplates(namespace).Get(t.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("error in GetPVC: %w", err)
		}
		r.ManagedFields = nil
		res, _ = json.Marshal(r)
	case constant.TEKTON_CRD_TYPE_EVENTLISTENERS:
		r, err := t.triggerset.TriggersV1beta1().EventListeners(namespace).Get(t.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("error in GetIngress: %w", err)
		}
		r.ManagedFields = nil
		res, _ = json.Marshal(r)
	case constant.TEKTON_CRD_TYPE_APPROVALTASKS:
		gvr := schema.GroupVersionResource{
			Group:    "tekton.automatiko.io",
			Version:  "v1beta1",
			Resource: "approvaltasks",
		}
		var r *unstructured.Unstructured
		if r, err = t.dynamicclient.Resource(gvr).Namespace(namespace).Get(t.ctx, resourceName, metav1.GetOptions{}); err != nil {
			return nil, fmt.Errorf("error in GetApprovalTasks: %w", err)
		}
		r.SetManagedFields(nil)
		res, _ = json.Marshal(r)
	default:
		return nil, errorx.NewDefaultError("Unknown resource type")
	}
	return
}

func (t *TektonService) GetResourceYAML(namespace, resourceType, resourceNames, withStatus string) (yamlContent string, err error) {
	var yamlItems []string
	for _, resourceName := range strings.Split(resourceNames, ",") {
		var unstructuredMap map[string]interface{}
		switch resourceType {
		case constant.TEKTON_CRD_TYPE_TASKS:
			res, err := t.clientset.TektonV1().Tasks(namespace).Get(t.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.TEKTON_CRD_TYPE_PIPELINES:
			res, err := t.clientset.TektonV1().Pipelines(namespace).Get(t.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.TEKTON_CRD_TYPE_PIPELINERUNS:
			res, err := t.clientset.TektonV1().PipelineRuns(namespace).Get(t.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.TEKTON_CRD_TYPE_TASKRUNS:
			res, err := t.clientset.TektonV1().TaskRuns(namespace).Get(t.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.TEKTON_CRD_TYPE_TRIGGERBINDINGS:
			res, err := t.triggerset.TriggersV1beta1().TriggerBindings(namespace).Get(t.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.TEKTON_CRD_TYPE_TRIGGERTEMPLATES:
			res, err := t.triggerset.TriggersV1beta1().TriggerTemplates(namespace).Get(t.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.TEKTON_CRD_TYPE_EVENTLISTENERS:
			res, err := t.triggerset.TriggersV1beta1().EventListeners(namespace).Get(t.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.TEKTON_CRD_TYPE_APPROVALTASKS:
			gvr := schema.GroupVersionResource{
				Group:    "tekton.automatiko.io",
				Version:  "v1beta1",
				Resource: "approvaltasks",
			}
			res, err := t.dynamicclient.Resource(gvr).Namespace(namespace).Get(t.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				return "", err
			}
			res.SetManagedFields(nil)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		default:
			return "", errorx.NewDefaultError("Unknown resource type")
		}
		if withStatus == "" {
			delete(unstructuredMap, "status")
		} else if withStatus == "true" {
			if _, ok := unstructuredMap["status"]; ok {
				unstructuredMap = unstructuredMap["status"].(map[string]interface{})
			} else {
				unstructuredMap = map[string]interface{}{}
			}
		}
		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
		unstructuredObj.SetKind(getKindVersion(resourceType)[0])
		unstructuredObj.SetAPIVersion(getKindVersion(resourceType)[1])
		annotations := unstructuredObj.GetAnnotations()
		delete(annotations, "kubectl.kubernetes.io/last-applied-configuration")
		unstructuredObj.SetAnnotations(annotations)
		yamlPrinter := &printers.YAMLPrinter{}
		var buf bytes.Buffer
		if err = yamlPrinter.PrintObj(unstructuredObj, &buf); err != nil {
			return "", err
		}
		yamlItems = append(yamlItems, buf.String())
		t.Logger.Debugf("\n%s", buf.String())
	}
	return strings.Join(yamlItems, "---\n"), nil
}

func (t *TektonService) DeleteResource(namespace, resourceType, resourceName string) (err error) {
	switch resourceType {
	case constant.TEKTON_CRD_TYPE_TASKS:
		return t.clientset.TektonV1().Tasks(namespace).Delete(t.ctx, resourceName, metav1.DeleteOptions{})
	case constant.TEKTON_CRD_TYPE_PIPELINES:
		return t.clientset.TektonV1().Pipelines(namespace).Delete(t.ctx, resourceName, metav1.DeleteOptions{})
	case constant.TEKTON_CRD_TYPE_PIPELINERUNS:
		return t.clientset.TektonV1().PipelineRuns(namespace).Delete(t.ctx, resourceName, metav1.DeleteOptions{})
	case constant.TEKTON_CRD_TYPE_TASKRUNS:
		return t.clientset.TektonV1().TaskRuns(namespace).Delete(t.ctx, resourceName, metav1.DeleteOptions{})
	case constant.TEKTON_CRD_TYPE_TRIGGERBINDINGS:
		return t.triggerset.TriggersV1beta1().TriggerBindings(namespace).Delete(t.ctx, resourceName, metav1.DeleteOptions{})
	case constant.TEKTON_CRD_TYPE_TRIGGERTEMPLATES:
		return t.triggerset.TriggersV1beta1().TriggerTemplates(namespace).Delete(t.ctx, resourceName, metav1.DeleteOptions{})
	case constant.TEKTON_CRD_TYPE_EVENTLISTENERS:
		return t.triggerset.TriggersV1beta1().EventListeners(namespace).Delete(t.ctx, resourceName, metav1.DeleteOptions{})
	case constant.TEKTON_CRD_TYPE_APPROVALTASKS:
		gvr := schema.GroupVersionResource{
			Group:    "tekton.automatiko.io",
			Version:  "v1beta1",
			Resource: "approvaltasks",
		}
		return t.dynamicclient.Resource(gvr).Namespace(namespace).Delete(t.ctx, resourceName, metav1.DeleteOptions{})
	default:
		return errorx.NewDefaultError("Unknown resource type")
	}
}

func (t *TektonService) PatchResource(namespace, resourceType, resourceName string, patchOptions []byte) (res []byte, err error) {
	switch resourceType {
	case constant.TEKTON_CRD_TYPE_PIPELINERUNS:
		var r *tektonv1.PipelineRun
		if r, err = t.clientset.TektonV1().PipelineRuns(namespace).Patch(t.ctx, resourceName, ktypes.JSONPatchType, patchOptions, metav1.PatchOptions{}); err != nil {
			return nil, fmt.Errorf("error in PatchPipelineRuns: %w", err)
		}
		res, _ = json.Marshal(r)
	case constant.TEKTON_CRD_TYPE_TASKRUNS:
		var r *tektonv1.TaskRun
		if r, err = t.clientset.TektonV1().TaskRuns(namespace).Patch(t.ctx, resourceName, ktypes.JSONPatchType, patchOptions, metav1.PatchOptions{}); err != nil {
			return nil, fmt.Errorf("error in PatchPipelineRuns: %w", err)
		}
		res, _ = json.Marshal(r)
	default:
		return nil, errorx.NewDefaultError("Unknown resource type")
	}
	return
}
