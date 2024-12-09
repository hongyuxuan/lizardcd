package svc

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/apps/v1"
	autov2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ktypes "k8s.io/apimachinery/pkg/types"
	uyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/restmapper"
)

type K8sService struct {
	logx.Logger
	ctx    context.Context
	svcCtx *ServiceContext
}

func GetK8sService(ctx context.Context, svcCtx *ServiceContext) *K8sService {
	return &K8sService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (k8s *K8sService) PatchDeployment(namespace, workloadName, containerName, imageName string) (res *v1.Deployment, err error) {
	if res, err = k8s.svcCtx.Clientset.AppsV1().Deployments(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{}); err != nil {
		return
	}
	var data string
	for _, container := range res.Spec.Template.Spec.Containers {
		if container.Name == containerName {
			data = fmt.Sprintf(`{ "spec": { "template": { "spec": { "containers": [ { "name": "%s", "image": "%s" } ] } } } }`, containerName, imageName)
			break
		}
	}
	for _, container := range res.Spec.Template.Spec.InitContainers {
		if container.Name == containerName {
			data = fmt.Sprintf(`{ "spec": { "template": { "spec": { "initContainers": [ { "name": "%s", "image": "%s" } ] } } } }`, containerName, imageName)
			break
		}
	}
	if data == "" {
		return nil, errorx.NewDefaultError("Deployment[%s] cannot find container[name=%s]", workloadName, containerName)
	}
	return k8s.svcCtx.Clientset.AppsV1().Deployments(namespace).Patch(k8s.ctx, workloadName, ktypes.StrategicMergePatchType, []byte(data), metav1.PatchOptions{})
}

func (k8s *K8sService) PatchStatefulset(namespace, workloadName, containerName, imageName string) (res *v1.StatefulSet, err error) {
	if res, err = k8s.svcCtx.Clientset.AppsV1().StatefulSets(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{}); err != nil {
		return
	}
	var data string
	for _, container := range res.Spec.Template.Spec.Containers {
		if container.Name == containerName {
			data = fmt.Sprintf(`{ "spec": { "template": { "spec": { "containers": [ { "name": "%s", "image": "%s" } ] } } } }`, containerName, imageName)
			break
		}
	}
	for _, container := range res.Spec.Template.Spec.InitContainers {
		if container.Name == containerName {
			data = fmt.Sprintf(`{ "spec": { "template": { "spec": { "initContainers": [ { "name": "%s", "image": "%s" } ] } } } }`, containerName, imageName)
			break
		}
	}
	if data == "" {
		return nil, errorx.NewDefaultError("Deployment[%s] cannot find container[name=%s]", workloadName, containerName)
	}
	return k8s.svcCtx.Clientset.AppsV1().StatefulSets(namespace).Patch(k8s.ctx, workloadName, ktypes.StrategicMergePatchType, []byte(data), metav1.PatchOptions{})
}

func (k8s *K8sService) ScaleDeployment(namespace string, workloadName string, replicas uint32) error {
	data := fmt.Sprintf(`{ "spec": { "replicas": %d } }`, replicas)
	_, err := k8s.svcCtx.Clientset.AppsV1().Deployments(namespace).Patch(k8s.ctx, workloadName, ktypes.StrategicMergePatchType, []byte(data), metav1.PatchOptions{})
	if err != nil {
		return errorx.NewDefaultError(err.Error())
	}
	k8s.Logger.Infof("Scale deployment[%s] replicas to %d success", workloadName, replicas)
	return nil
}

func (k8s *K8sService) DeleteDeployment(namespace, workloadName string) (err error) {
	return k8s.svcCtx.Clientset.AppsV1().Deployments(namespace).Delete(k8s.ctx, workloadName, metav1.DeleteOptions{})
}

func (k8s *K8sService) DeleteStatefulset(namespace, workloadName string) (err error) {
	return k8s.svcCtx.Clientset.AppsV1().StatefulSets(namespace).Delete(k8s.ctx, workloadName, metav1.DeleteOptions{})
}

func (k8s *K8sService) GetDeploymentReplicas(namespace string, workloads []string) (workloadReplicas []map[string]interface{}, err error) {
	for _, workload := range workloads {
		res, err := k8s.svcCtx.Clientset.AppsV1().Deployments(namespace).Get(k8s.ctx, workload, metav1.GetOptions{})
		if err != nil {
			return nil, errorx.NewDefaultError(err.Error())
		}
		workloadReplicas = append(workloadReplicas, map[string]interface{}{
			"name":     workload,
			"replicas": int(*res.Spec.Replicas),
		})
	}
	return
}

func (k8s *K8sService) GetDeploymentImages(namespace string, initContainer bool, workloads []string) (workloadImages []map[string]interface{}, err error) {
	for _, workload := range workloads {
		res, err := k8s.svcCtx.Clientset.AppsV1().Deployments(namespace).Get(k8s.ctx, workload, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			continue
		} else if err != nil {
			return nil, errorx.NewDefaultError(err.Error())
		}
		var image string
		if initContainer {
			if len(res.Spec.Template.Spec.InitContainers) == 0 {
				return nil, errorx.NewDefaultError("There is no initContainers in workload \"%s\"", workload)
			}
			image = res.Spec.Template.Spec.InitContainers[0].Image
		} else {
			image = res.Spec.Template.Spec.Containers[0].Image
		}
		workloadImages = append(workloadImages, map[string]interface{}{
			"name":  workload,
			"image": image,
		})
	}
	return
}

func (k8s *K8sService) ScaleStatefulset(namespace string, workloadName string, replicas uint32) error {
	data := fmt.Sprintf(`{ "spec": { "replicas": %d } }`, replicas)
	_, err := k8s.svcCtx.Clientset.AppsV1().StatefulSets(namespace).Patch(k8s.ctx, workloadName, ktypes.StrategicMergePatchType, []byte(data), metav1.PatchOptions{})
	if err != nil {
		return errorx.NewDefaultError(err.Error())
	}
	k8s.Logger.Infof("Scale statefulset[%s] replicas to %d success", workloadName, replicas)
	return nil
}

func (k8s *K8sService) GetDeploymentPodInfo(namespace, workloadName string) ([]corev1.Pod, error) {
	res, err := k8s.svcCtx.Clientset.AppsV1().Deployments(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{})
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	labels := res.Spec.Template.ObjectMeta.Labels
	return k8s.GetPodStatus(namespace, labels)
}

func (k8s *K8sService) GetStatefulsetPodInfo(namespace, workloadName string) ([]corev1.Pod, error) {
	res, err := k8s.svcCtx.Clientset.AppsV1().StatefulSets(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{})
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	labels := res.Spec.Template.ObjectMeta.Labels
	return k8s.GetPodStatus(namespace, labels)
}

func (k8s *K8sService) GetDeploymentPodStatus(namespace, workloadName string) (*commontypes.WorkloadStatus, error) {
	res, err := k8s.svcCtx.Clientset.AppsV1().Deployments(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{})
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	labels := res.Spec.Template.ObjectMeta.Labels
	Pods, err := k8s.GetPodStatus(namespace, labels)
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	var pods []commontypes.PodStatus
	for _, v := range Pods {
		var readyStatus string
		for _, c := range v.Status.Conditions {
			if c.Type == "Ready" {
				readyStatus = string(c.Status)
				break
			}
		}
		podstatus := commontypes.PodStatus{
			PodName: v.Name,
			Ready:   readyStatus,
		}
		pods = append(pods, podstatus)
	}
	var revision string
	for k, v := range res.ObjectMeta.Annotations {
		if strings.Contains(k, "kubernetes.io/revision") {
			revision = v
			break
		}
	}
	workLoadStatus := commontypes.WorkloadStatus{
		Name:     workloadName,
		Revision: revision,
		Pods:     pods,
	}
	return &workLoadStatus, nil
}

func (k8s *K8sService) GetStatefulsetPodStatus(namespace, workloadName string) (*commontypes.WorkloadStatus, error) {
	res, err := k8s.svcCtx.Clientset.AppsV1().StatefulSets(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{})
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	labels := res.Spec.Template.ObjectMeta.Labels
	Pods, err := k8s.GetPodStatus(namespace, labels)
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	var pods []commontypes.PodStatus
	for _, v := range Pods {
		var readyStatus string
		for _, c := range v.Status.Conditions {
			if c.Type == "Ready" {
				readyStatus = string(c.Status)
				break
			}
		}
		podstatus := commontypes.PodStatus{
			PodName: v.Name,
			Ready:   readyStatus,
		}
		pods = append(pods, podstatus)
	}
	workLoadStatus := commontypes.WorkloadStatus{
		Name: workloadName,
		Pods: pods,
	}
	return &workLoadStatus, nil
}

func (k8s *K8sService) GetPodStatus(namespace string, labels map[string]string) ([]corev1.Pod, error) {
	delete(labels, "pod-template-hash")
	var labelSelector []string
	for k, v := range labels {
		labelSelector = append(labelSelector, k+"="+v)
	}
	res, err := k8s.svcCtx.Clientset.CoreV1().Pods(namespace).List(k8s.ctx, metav1.ListOptions{
		LabelSelector: strings.Join(labelSelector, ","),
	})
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	for i := range res.Items {
		res.Items[i].ManagedFields = nil
	}
	// k8s.Logger.Debugf("Get pods with namespace: %s, labelSelector: %s, items=%+v", namespace, strings.Join(labelSelector, ","), res.Items)
	return res.Items, nil
}

func (k8s *K8sService) GetEvents(namespace, objectKind, objectName string) ([]corev1.Event, error) {
	res, err := k8s.svcCtx.Clientset.CoreV1().Events(namespace).List(k8s.ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=%s", objectName, objectKind),
	})
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	return res.Items, nil
}

func (k8s *K8sService) GetCoreV1ResourceYAML(namespace, resourceType, resourceNames string) (string, error) {
	secret, err := k8s.GetSecret(namespace)
	if err != nil {
		return "", err
	}
	if secret == nil {
		errMsg := fmt.Sprintf("Cannot find secret with prefix \"%s\" in namespace \"%s\"", k8s.svcCtx.Config.KubernetesSecretPrefix, namespace)
		k8s.Logger.Error(errMsg)
		return "", errorx.NewDefaultError(errMsg)
	}
	token := string(secret.Data["token"])
	k8s.Logger.Debugf("Get default token=%s", token)
	var yamlItems []string
	for _, resourceName := range strings.Split(resourceNames, ",") {
		res, err := k8s.svcCtx.Request_k8s.R().SetBearerAuthToken(token).Get(fmt.Sprintf("/api/v1/namespaces/%s/%s/%s", namespace, resourceType, resourceName))
		if err != nil {
			k8s.Logger.Error(err)
			return "", errorx.NewDefaultError(err.Error())
		}
		v := make(map[string]interface{})
		res.UnmarshalJson(&v)
		delete(v, "status")
		metadata := v["metadata"].(map[string]interface{})
		delete(metadata, "managedFields")
		delete(metadata, "uid")
		if resourceType != "services" {
			delete(metadata, "resourceVersion")
		}
		delete(metadata, "selfLink")
		v["metadata"] = metadata
		manifest, _ := yaml.Marshal(v)
		yamlItems = append(yamlItems, string(manifest))
	}
	return strings.Join(yamlItems, "---\n"), nil
}

func (k8s *K8sService) GetAppsV1ResourceYAML(namespace, resourceType, resourceNames string) (string, error) {
	secret, err := k8s.GetSecret(namespace)
	if err != nil {
		k8s.Logger.Error(err)
		return "", errorx.NewDefaultError(err.Error())
	}
	if secret == nil {
		errMsg := fmt.Sprintf("Cannot find secret with prefix \"%s\" in namespace \"%s\"", k8s.svcCtx.Config.KubernetesSecretPrefix, namespace)
		k8s.Logger.Error(errMsg)
		return "", errorx.NewDefaultError(errMsg)
	}
	token := string(secret.Data["token"])
	k8s.Logger.Debugf("Get default token=%s", token)
	var yamlItems []string
	for _, resourceName := range strings.Split(resourceNames, ",") {
		res, err := k8s.svcCtx.Request_k8s.R().SetBearerAuthToken(token).Get(fmt.Sprintf("/apis/apps/v1/namespaces/%s/%s/%s", namespace, resourceType, resourceName))
		if err != nil {
			panic(err)
		}
		v := make(map[string]interface{})
		res.UnmarshalJson(&v)
		delete(v, "status")
		metadata := v["metadata"].(map[string]interface{})
		delete(metadata, "managedFields")
		delete(metadata, "uid")
		delete(metadata, "resourceVersion")
		delete(metadata, "generation")
		delete(metadata, "selfLink")
		v["metadata"] = metadata
		manifest, _ := yaml.Marshal(v)
		yamlItems = append(yamlItems, string(manifest))
	}
	return strings.Join(yamlItems, "---\n"), nil
}

func (k8s *K8sService) GetIngressYAML(namespace, resourceNames string) (string, error) {
	secret, err := k8s.GetSecret(namespace)
	if err != nil {
		return "", err
	}
	if secret == nil {
		errMsg := fmt.Sprintf("Cannot find secret with prefix \"%s\" in namespace \"%s\"", k8s.svcCtx.Config.KubernetesSecretPrefix, namespace)
		k8s.Logger.Error(errMsg)
		return "", errorx.NewDefaultError(errMsg)
	}
	token := string(secret.Data["token"])
	k8s.Logger.Debugf("Get default token=%s", token)
	var yamlItems []string
	for _, resourceName := range strings.Split(resourceNames, ",") {
		res, err := k8s.svcCtx.Request_k8s.R().SetBearerAuthToken(token).Get(fmt.Sprintf("/apis/networking.k8s.io/v1/namespaces/%s/ingresses/%s", namespace, resourceName))
		if err != nil {
			k8s.Logger.Error(err)
			return "", errorx.NewDefaultError(err.Error())
		}
		v := make(map[string]interface{})
		res.UnmarshalJson(&v)
		delete(v, "status")
		metadata := v["metadata"].(map[string]interface{})
		delete(metadata, "managedFields")
		delete(metadata, "uid")
		delete(metadata, "resourceVersion")
		delete(metadata, "selfLink")
		v["metadata"] = metadata
		manifest, _ := yaml.Marshal(v)
		yamlItems = append(yamlItems, string(manifest))
	}
	return strings.Join(yamlItems, "---\n"), nil
}

func (k8s *K8sService) GetConfigMap(namespace string, configmapName string) (map[string]string, error) {
	res, err := k8s.svcCtx.Clientset.CoreV1().ConfigMaps(namespace).Get(k8s.ctx, configmapName, metav1.GetOptions{})
	if err != nil {
		k8s.Logger.Error(err)
		return nil, errorx.NewDefaultError(err.Error())
	}
	return res.Data, nil
}

func (k8s *K8sService) PatchConfigMap(namespace, configmapName, key string, value string) (*corev1.ConfigMap, error) {
	data := fmt.Sprintf(`{ "data": { "%s": "%s" } }`, key, value)
	res, err := k8s.svcCtx.Clientset.CoreV1().ConfigMaps(namespace).Patch(k8s.ctx, configmapName, ktypes.StrategicMergePatchType, []byte(data), metav1.PatchOptions{})
	if err != nil {
		k8s.Logger.Error(err)
		return nil, errorx.NewDefaultError(err.Error())
	}
	return res, nil
}

func (k8s *K8sService) GetSecret(namespace string) (*corev1.Secret, error) {
	res, err := k8s.svcCtx.Clientset.CoreV1().Secrets(namespace).List(k8s.ctx, metav1.ListOptions{})
	if err != nil {
		k8s.Logger.Error(err)
		return nil, errorx.NewDefaultError(err.Error())
	}
	for _, secret := range res.Items {
		if strings.HasPrefix(secret.Name, k8s.svcCtx.Config.KubernetesSecretPrefix) {
			return &secret, nil
		}
	}
	return nil, nil
}

func (k8s *K8sService) UpdateFromYaml(namespace, applyYaml, kind string, taskResult chan map[string]interface{}) {
	d := uyaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(applyYaml), 4096)
	for {
		unstructureObj, err := utils.GetUnstructured(d)
		if err == io.EOF {
			break
		}
		if err != nil {
			k8s.Logger.Error(err)
			taskResult <- map[string]interface{}{"success": false, "message": err.Error()}
			return
		}
		if namespace != unstructureObj.GetNamespace() {
			taskResult <- map[string]interface{}{"success": false, "message": fmt.Sprintf("Namespace must be %s", namespace)}
			return
		}
		if kind != "" && kind != unstructureObj.GetKind() {
			taskResult <- map[string]interface{}{"success": false, "message": fmt.Sprintf("Kind must be %s", kind)}
			return
		}
		gvr, err := k8s.gtGVR(unstructureObj.GroupVersionKind())
		if err != nil {
			taskResult <- map[string]interface{}{"success": false, "message": err.Error()}
			return
		}
		_, getErr := k8s.svcCtx.Dynamicclient.Resource(gvr).Namespace(namespace).Get(context.Background(), unstructureObj.GetName(), metav1.GetOptions{})
		if getErr != nil {
			_, createErr := k8s.svcCtx.Dynamicclient.Resource(gvr).Namespace(namespace).Create(context.Background(), unstructureObj, metav1.CreateOptions{})
			if createErr != nil {
				taskResult <- map[string]interface{}{"success": false, "message": createErr.Error()}
				return
			}
			k8s.Logger.Infof("Create resource[%s] success", unstructureObj.GetName())
			taskResult <- map[string]interface{}{
				"success":      true,
				"message":      "success",
				"workloadType": unstructureObj.GetKind(),
				"workloadName": unstructureObj.GetName(),
			}
			return
		}
		if namespace == unstructureObj.GetNamespace() {
			_, err = k8s.svcCtx.Dynamicclient.Resource(gvr).Namespace(namespace).Update(context.Background(), unstructureObj, metav1.UpdateOptions{})
			if err != nil {
				k8s.Logger.Errorf("unable to update resource[%s]: %+v", unstructureObj.GetName(), err)
				taskResult <- map[string]interface{}{
					"success": false, "message": fmt.Sprintf("Unable to update resource[%s]: %s", unstructureObj.GetName(), err.Error()),
				}
				return
			}
			k8s.Logger.Infof("Update resource[%s] success", unstructureObj.GetName())
			taskResult <- map[string]interface{}{
				"success":      true,
				"message":      "success",
				"workloadType": unstructureObj.GetKind(),
				"workloadName": unstructureObj.GetName(),
			}
		} else {
			_, err = k8s.svcCtx.Dynamicclient.Resource(gvr).Update(context.Background(), unstructureObj, metav1.UpdateOptions{})
			if err != nil {
				k8s.Logger.Errorf("Ns is nil unable to update resource: %v", err)
				taskResult <- map[string]interface{}{"success": false, "message": fmt.Sprintf("Ns is nil unable to update resource: %s", err.Error())}
				return
			}
			taskResult <- map[string]interface{}{
				"success":      true,
				"message":      "success",
				"workloadType": unstructureObj.GetKind(),
				"workloadName": unstructureObj.GetName(),
			}
		}
	}
}

func (k8s *K8sService) DeleteFromYaml(namespace string, applyYaml string, taskResult chan map[string]interface{}) {
	d := uyaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(applyYaml), 4096)
	for {
		unstructureObj, err := utils.GetUnstructured(d)
		if err == io.EOF {
			break
		}
		if err != nil {
			taskResult <- map[string]interface{}{"success": false, "message": err.Error()}
			return
		}
		gvr, err := k8s.gtGVR(unstructureObj.GroupVersionKind())
		if err != nil {
			taskResult <- map[string]interface{}{"success": false, "message": err.Error()}
			return
		}

		if namespace == unstructureObj.GetNamespace() {
			err := k8s.svcCtx.Dynamicclient.Resource(gvr).Namespace(namespace).Delete(context.Background(), unstructureObj.GetName(), metav1.DeleteOptions{})
			if err != nil {
				k8s.Logger.Errorf("unable to delete resource[%s]: %v", unstructureObj.GetName(), err)
				taskResult <- map[string]interface{}{
					"success": false, "message": fmt.Sprintf("Unable to delete resource[%s]: %s", unstructureObj.GetName(), err.Error()),
				}
				return
			}
			k8s.Logger.Infof("Delete resource[%s] success", unstructureObj.GetName())
			taskResult <- map[string]interface{}{"success": true, "message": "success"}
		}
	}
}

// curl -v -XPATCH  -H "Content-Type: application/strategic-merge-patch+json" -H "User-Agent: kubectl/v1.23.5 (linux/amd64) kubernetes/c285e78" -H "Accept: application/json, */*" 'https://10.21.131.253:6443/apis/apps/v1/namespaces/default/deployments/tools?fieldManager=kubectl-rollout'
// request body: {"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"2023-04-17T08:32:32Z"}}}}}
func (k8s *K8sService) RolloutDeployment(namespace, workloadName string) (res *v1.Deployment, err error) {
	data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().Format("2006-01-02T15:04:05Z"))
	res, err = k8s.svcCtx.Clientset.AppsV1().Deployments(namespace).Patch(k8s.ctx, workloadName, ktypes.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
	if err != nil {
		k8s.Logger.Error(err)
	}
	return
}

func (k8s *K8sService) RolloutStatefulset(namespace, workloadName string) (res *v1.StatefulSet, err error) {
	data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().Format("2006-01-02T15:04:05Z"))
	res, err = k8s.svcCtx.Clientset.AppsV1().StatefulSets(namespace).Patch(k8s.ctx, workloadName, ktypes.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
	if err != nil {
		k8s.Logger.Error(err)
	}
	return
}

func (k8s *K8sService) gtGVR(gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {
	gr, err := restmapper.GetAPIGroupResources(k8s.svcCtx.Clientset.Discovery())
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	mapper := restmapper.NewDiscoveryRESTMapper(gr)

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	return mapping.Resource, nil
}

// func (k8s *K8sService) getUnstructured(d *uyaml.YAMLOrJSONDecoder) (unstructureObj *unstructured.Unstructured, err error) {
// 	var rawObj runtime.RawExtension
// 	err = d.Decode(&rawObj)
// 	if err == io.EOF {
// 		return
// 	}
// 	if err != nil {
// 		err = errorx.NewDefaultError("Decode is err: %v", err.Error())
// 		return
// 	}
// 	obj, _, err := syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
// 	if err != nil {
// 		err = errorx.NewDefaultError("Rawobj is err: %v", err.Error())
// 		return
// 	}
// 	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
// 	if err != nil {
// 		err = errorx.NewDefaultError("Tounstructured is err %v", err.Error())
// 		return
// 	}
// 	unstructureObj = &unstructured.Unstructured{Object: unstructuredMap}
// 	return
// }

func (k8s *K8sService) GetNamespaces(labelSelector string) (res *corev1.NamespaceList, err error) {
	if res, err = k8s.svcCtx.Clientset.CoreV1().Namespaces().List(k8s.ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	}); err != nil {
		k8s.Logger.Error(err)
	}
	return
}

func (k8s *K8sService) ListDeployment(namespace, labelSelector string) (res []v1.Deployment, err error) {
	var conti string
	for {
		var listRes *v1.DeploymentList
		if listRes, err = k8s.svcCtx.Clientset.AppsV1().Deployments(namespace).List(k8s.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			Continue:      conti,
			Limit:         10,
		}); err != nil {
			k8s.Logger.Error(err)
			return
		}
		res = append(res, processDeploymentItems(listRes)...)
		if listRes.Continue == "" {
			break
		}
		conti = listRes.Continue
	}
	return
}

func (k8s *K8sService) ListStatefulset(namespace, labelSelector string) (res []v1.StatefulSet, err error) {
	var conti string
	for {
		var listRes *v1.StatefulSetList
		if listRes, err = k8s.svcCtx.Clientset.AppsV1().StatefulSets(namespace).List(k8s.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			Continue:      conti,
			Limit:         10,
		}); err != nil {
			k8s.Logger.Error(err)
			return
		}
		res = append(res, processStatefulsetItems(listRes)...)
		if listRes.Continue == "" {
			break
		}
		conti = listRes.Continue
	}
	return
}

func (k8s *K8sService) GetDeployment(namespace, workloadName string) (*v1.Deployment, error) {
	res, err := k8s.svcCtx.Clientset.AppsV1().Deployments(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	res.ManagedFields = nil
	return res, nil
}

func (k8s *K8sService) GetStatefulset(namespace, workloadName string) (*v1.StatefulSet, error) {
	res, err := k8s.svcCtx.Clientset.AppsV1().StatefulSets(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	res.ManagedFields = nil
	return res, nil
}

func (k8s *K8sService) GetDeploymentResourceQuota(namespace, workloadName string) (*corev1.ResourceRequirements, error) {
	res, err := k8s.svcCtx.Clientset.AppsV1().Deployments(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &res.Spec.Template.Spec.Containers[0].Resources, nil
}

func (k8s *K8sService) GetStatefulsetResourceQuota(namespace, workloadName string) (*corev1.ResourceRequirements, error) {
	res, err := k8s.svcCtx.Clientset.AppsV1().StatefulSets(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &res.Spec.Template.Spec.Containers[0].Resources, nil
}

func (k8s *K8sService) GetPodLog(namespace, podname, containerName string, lines int64, follow bool, stream agent.LizardAgent_GetPodLogFollowServer) (logs string, err error) {
	req := k8s.svcCtx.Clientset.CoreV1().Pods(namespace).GetLogs(podname, &corev1.PodLogOptions{
		TailLines: &lines,
		Container: containerName,
		Follow:    follow,
	})
	var ioLogs io.ReadCloser
	if ioLogs, err = req.Stream(k8s.ctx); err != nil {
		k8s.Logger.Error(err)
		return
	}
	defer ioLogs.Close()
	if follow {
		r := bufio.NewReader(ioLogs)
		for {
			var b []byte
			if b, err = r.ReadBytes('\n'); err != nil {
				k8s.Logger.Error(err)
				return
			}
			if err = stream.Send(&agent.YamlResponse{
				Code: uint32(codes.OK),
				Data: string(b),
			}); err != nil {
				k8s.Logger.Error(err)
				return
			}
		}
	} else {
		buf := new(bytes.Buffer)
		if _, err = io.Copy(buf, ioLogs); err != nil {
			k8s.Logger.Error(err)
			return
		}
		logs = buf.String()
	}
	return
}

func (k8s *K8sService) SetPodHorizonAutoscaler(namespace, workloadName string, max, min, cpu, memory *int32) (*autov2.HorizontalPodAutoscaler, error) {
	var metrics []autov2.MetricSpec
	if cpu != nil {
		metrics = append(metrics, autov2.MetricSpec{
			Type: "Resource",
			Resource: &autov2.ResourceMetricSource{
				Name: "cpu",
				Target: autov2.MetricTarget{
					Type:               "Utilization",
					AverageUtilization: cpu,
				},
			},
		})
	}
	if memory != nil {
		metrics = append(metrics, autov2.MetricSpec{
			Type: "Resource",
			Resource: &autov2.ResourceMetricSource{
				Name: "memory",
				Target: autov2.MetricTarget{
					Type:               "Utilization",
					AverageUtilization: memory,
				},
			},
		})
	}
	hpa := &autov2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name: workloadName,
		},
		Spec: autov2.HorizontalPodAutoscalerSpec{
			MinReplicas: min,
			MaxReplicas: *max,
			Metrics:     metrics,
			ScaleTargetRef: autov2.CrossVersionObjectReference{
				Kind:       "Deployment",
				APIVersion: "apps/v1",
				Name:       workloadName,
			},
		},
	}
	res, err := k8s.svcCtx.Clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return k8s.svcCtx.Clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).Create(k8s.ctx, hpa, metav1.CreateOptions{})
	} else {
		hpa.ResourceVersion = res.ResourceVersion
		return k8s.svcCtx.Clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).Update(k8s.ctx, hpa, metav1.UpdateOptions{})
	}
}

func (k8s *K8sService) GetPodHorizonAutoscaler(namespace, workloadName string) (*autov2.HorizontalPodAutoscaler, error) {
	res, err := k8s.svcCtx.Clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, err
	}
	res.ManagedFields = nil
	return res, nil
}

func processDeploymentItems(listRes *v1.DeploymentList) []v1.Deployment {
	for i := range listRes.Items {
		listRes.Items[i].ManagedFields = nil
		listRes.Items[i].Spec = v1.DeploymentSpec{}
		delete(listRes.Items[i].Annotations, "kubectl.kubernetes.io/last-applied-configuration")
	}
	return listRes.Items
}

func processStatefulsetItems(listRes *v1.StatefulSetList) []v1.StatefulSet {
	for i := range listRes.Items {
		listRes.Items[i].ManagedFields = nil
		listRes.Items[i].Spec = v1.StatefulSetSpec{}
		delete(listRes.Items[i].Annotations, "kubectl.kubernetes.io/last-applied-configuration")
	}
	return listRes.Items
}
