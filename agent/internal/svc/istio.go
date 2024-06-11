package svc

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	istiometav1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ktypes "k8s.io/apimachinery/pkg/types"
)

type IstioService struct {
	logx.Logger
	ctx    context.Context
	svcCtx *ServiceContext
}

func GetIstioService(ctx context.Context, svcCtx *ServiceContext) *IstioService {
	return &IstioService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (istio *IstioService) CreateDestinationRule(namespace string, destinationrule *v1beta1.DestinationRule) (*v1beta1.DestinationRule, error) {
	return istio.svcCtx.Istioclient.NetworkingV1beta1().DestinationRules(namespace).Create(istio.ctx, destinationrule, metav1.CreateOptions{})
}

func (istio *IstioService) ListDestinationRule(namespace string, labelSelector string) (res []*v1beta1.DestinationRule, err error) {
	var conti string
	for {
		var listRes *v1beta1.DestinationRuleList
		if listRes, err = istio.svcCtx.Istioclient.NetworkingV1beta1().DestinationRules(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
			Continue:      conti,
			Limit:         10,
		}); err != nil {
			istio.Logger.Error(err)
			return
		}
		res = append(res, processDestinationItems(listRes)...)
		if listRes.Continue == "" {
			break
		}
		conti = listRes.Continue
	}
	return
}

func (istio *IstioService) PatchDestinationRule(namespace, name string, data []byte) (*v1beta1.DestinationRule, error) {
	istio.Logger.Infof("Patch destinationrule data: %s", string(data))
	return istio.svcCtx.Istioclient.NetworkingV1beta1().DestinationRules(namespace).Patch(istio.ctx, name, ktypes.MergePatchType, data, metav1.PatchOptions{})
}

func (istio *IstioService) GetDestinationRule(namespace, name string) (res *v1beta1.DestinationRule, err error) {
	res, err = istio.svcCtx.Istioclient.NetworkingV1beta1().DestinationRules(namespace).Get(istio.ctx, name, metav1.GetOptions{})
	res.APIVersion = "networking.istio.io/v1beta1"
	res.Kind = "DestinationRule"
	res.UID = ""
	res.Generation = 0
	res.SelfLink = ""
	res.ManagedFields = nil
	return
}

func (istio *IstioService) DeleteDestinationRule(namespace, name string) error {
	return istio.svcCtx.Istioclient.NetworkingV1beta1().DestinationRules(namespace).Delete(istio.ctx, name, metav1.DeleteOptions{})
}

func (istio *IstioService) CreateVirtualService(namespace string, virtualservice *v1beta1.VirtualService) (*v1beta1.VirtualService, error) {
	return istio.svcCtx.Istioclient.NetworkingV1beta1().VirtualServices(namespace).Create(istio.ctx, virtualservice, metav1.CreateOptions{})
}

func (istio *IstioService) ListVirtualService(namespace string, labelSelector string) (res []*v1beta1.VirtualService, err error) {
	var conti string
	for {
		var listRes *v1beta1.VirtualServiceList
		if listRes, err = istio.svcCtx.Istioclient.NetworkingV1beta1().VirtualServices(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
			Continue:      conti,
			Limit:         10,
		}); err != nil {
			istio.Logger.Error(err)
			return
		}
		res = append(res, processVirtualItems(listRes)...)
		if listRes.Continue == "" {
			break
		}
		conti = listRes.Continue
	}
	return
}

func (istio *IstioService) PatchVirtualService(namespace, name string, data []byte) (*v1beta1.VirtualService, error) {
	istio.Logger.Infof("Patch virtualservice data: %s", string(data))
	return istio.svcCtx.Istioclient.NetworkingV1beta1().VirtualServices(namespace).Patch(istio.ctx, name, ktypes.MergePatchType, data, metav1.PatchOptions{})
}

func (istio *IstioService) GetVirtualService(namespace, name string) (res *v1beta1.VirtualService, err error) {
	res, err = istio.svcCtx.Istioclient.NetworkingV1beta1().VirtualServices(namespace).Get(istio.ctx, name, metav1.GetOptions{})
	res.APIVersion = "networking.istio.io/v1beta1"
	res.Kind = "VirtualService"
	res.UID = ""
	res.Generation = 0
	res.SelfLink = ""
	res.ManagedFields = nil
	return
}

func (istio *IstioService) DeleteVirtualService(namespace, name string) error {
	return istio.svcCtx.Istioclient.NetworkingV1beta1().VirtualServices(namespace).Delete(istio.ctx, name, metav1.DeleteOptions{})
}

func processDestinationItems(listRes *v1beta1.DestinationRuleList) []*v1beta1.DestinationRule {
	for i := range listRes.Items {
		listRes.Items[i].ManagedFields = nil
		listRes.Items[i].Spec = istiometav1.DestinationRule{}
		delete(listRes.Items[i].Annotations, "kubectl.kubernetes.io/last-applied-configuration")
	}
	return listRes.Items
}

func processVirtualItems(listRes *v1beta1.VirtualServiceList) []*v1beta1.VirtualService {
	for i := range listRes.Items {
		listRes.Items[i].ManagedFields = nil
		listRes.Items[i].Spec = istiometav1.VirtualService{}
		delete(listRes.Items[i].Annotations, "kubectl.kubernetes.io/last-applied-configuration")
	}
	return listRes.Items
}
