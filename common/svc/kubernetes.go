package svc

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/constant"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	triggerclient "github.com/tektoncd/triggers/pkg/client/clientset/versioned"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/logx"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/api/apps/v1"
	autov2 "k8s.io/api/autoscaling/v2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ktypes "k8s.io/apimachinery/pkg/types"
	uyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
)

type K8sService struct {
	logx.Logger
	ctx           context.Context
	Clientset     *kubernetes.Clientset
	dynamicclient dynamic.Interface
	TektonService *TektonService
	HelmService   *HelmService
}

func NewK8sService(ctx context.Context, Clientset *kubernetes.Clientset, dynamicclient dynamic.Interface, tektonset *tektonclient.Clientset, triggerset *triggerclient.Clientset) *K8sService {
	return &K8sService{
		Logger:        logx.WithContext(ctx),
		ctx:           ctx,
		Clientset:     Clientset,
		dynamicclient: dynamicclient,
		TektonService: NewTektonService(ctx, tektonset, triggerset, dynamicclient),
		HelmService:   NewHelmService(ctx),
	}
}

func CreateKubernetes(kubeconfigPath, kubeconfig string) (*kubernetes.Clientset, *rest.Config, dynamic.Interface, *versionedclient.Clientset, *tektonclient.Clientset, *triggerclient.Clientset, string) {
	var conf *rest.Config
	var token string
	var err error
	var tektonClient *tektonclient.Clientset
	var triggerClient *triggerclient.Clientset
	if kubeconfig != "" {
		config, err := clientcmd.NewClientConfigFromBytes([]byte(kubeconfig))
		if err != nil {
			logx.Error(err)
			os.Exit(0)
		}
		conf, err = config.ClientConfig()
		if err != nil {
			logx.Error(err)
			os.Exit(0)
		}

	} else if kubeconfigPath != "" {
		logx.Infof("Using kubeconfig=%s", kubeconfigPath)
		conf, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			logx.Error(err)
			os.Exit(0)
		}
	} else {
		logx.Info("Using in cluster config")
		conf, err = rest.InClusterConfig()
		if conf != nil {
			conf.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(1000, 1000) // setting a big ratelimiter for client-side throttling, default 5
		}
		if err != nil {
			logx.Errorf("Cannot build connection to kubernetes with kubeconfig or in-cluster, err: %v. Will start as an vm agent.", err)
			return nil, nil, nil, nil, nil, nil, ""
		}
	}
	Clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		logx.Error(err)
		os.Exit(0)
	}
	dynamicclient, err := dynamic.NewForConfig(conf)
	if err != nil {
		logx.Error(err)
		os.Exit(0)
	}
	if conf != nil {
		token = conf.BearerToken
		tektonClient, _ = tektonclient.NewForConfig(conf)
		triggerClient, _ = triggerclient.NewForConfig(conf)
	}
	logx.Infof("Connected to k8s: %s success", conf.Host)

	// create istio client if posible
	istioclient, err := versionedclient.NewForConfig(conf)
	if err != nil {
		logx.Errorf("Failed to create istio client: %s", err)
	}
	return Clientset, conf, dynamicclient, istioclient, tektonClient, triggerClient, token
}

func (k8s *K8sService) IsValid() bool {
	return k8s.Clientset != nil
}

func (k8s *K8sService) PatchWorkload(namespace, workloadType, workloadName, containerName, imageName string) (res []byte, err error) {
	k8s.Logger.Infof("Patch namespace=%s %s[%v] container=%v image=%v", namespace, workloadType, workloadName, containerName, imageName)
	var data string
	switch workloadType {
	case constant.K8S_RESOURCE_TYPE_DEPLOYMENTS:
		var r *v1.Deployment
		if r, err = k8s.Clientset.AppsV1().Deployments(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{}); err != nil {
			return nil, fmt.Errorf("error in GetDeployment: %v", err)
		}
		if data = getPatchWorkloadData(false, r.Spec.Template.Spec.Containers, containerName, imageName); data == "" {
			if data = getPatchWorkloadData(true, r.Spec.Template.Spec.InitContainers, containerName, imageName); data == "" {
				return nil, errorx.NewDefaultError("Deployment[%s] cannot find container[name=%s]", workloadName, containerName)
			}
		}
		if r, err = k8s.Clientset.AppsV1().Deployments(namespace).Patch(k8s.ctx, workloadName, ktypes.StrategicMergePatchType, []byte(data), metav1.PatchOptions{}); err != nil {
			return nil, fmt.Errorf("error in PatchDeployment: %v", err)
		}
		res, _ = json.Marshal(r)
	case constant.K8S_RESOURCE_TYPE_STATEFULSETS:
		var r *v1.StatefulSet
		if r, err = k8s.Clientset.AppsV1().StatefulSets(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{}); err != nil {
			return nil, fmt.Errorf("error in GetStatefulSet: %v", err)
		}
		if data = getPatchWorkloadData(false, r.Spec.Template.Spec.Containers, containerName, imageName); data == "" {
			if data = getPatchWorkloadData(true, r.Spec.Template.Spec.InitContainers, containerName, imageName); data == "" {
				return nil, errorx.NewDefaultError("StatefulSet[%s] cannot find container[name=%s]", workloadName, containerName)
			}
		}
		if r, err = k8s.Clientset.AppsV1().StatefulSets(namespace).Patch(k8s.ctx, workloadName, ktypes.StrategicMergePatchType, []byte(data), metav1.PatchOptions{}); err != nil {
			return nil, fmt.Errorf("error in PatchStatefulSet: %v", err)
		}
		res, _ = json.Marshal(r)
	case constant.K8S_RESOURCE_TYPE_JOB:
		var r *batchv1.Job
		if r, err = k8s.PatchJob(namespace, workloadName, containerName, imageName); err != nil {
			return
		}
		res, _ = json.Marshal(r)
	default:
		return nil, fmt.Errorf("unknown workload type")
	}
	return
}

func (k8s *K8sService) PatchJob(namespace, workloadName, containerName, imageName string) (res *batchv1.Job, err error) {
	if res, err = k8s.Clientset.BatchV1().Jobs(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{}); err != nil {
		return nil, fmt.Errorf("error in GetJob: %v", err)
	}
	// jobs image is immutable, must delete job first and recreate
	var seconds int64 = 0
	if err = k8s.Clientset.BatchV1().Jobs(namespace).Delete(k8s.ctx, workloadName, metav1.DeleteOptions{
		GracePeriodSeconds: &seconds,
	}); err != nil {
		return nil, fmt.Errorf("error in DeleteJob: %v", err)
	}
	// delete pod
	var pods []corev1.Pod
	if pods, _, err = k8s.ListPods(namespace, res.Spec.Selector.MatchLabels, "", "", 0); err != nil {
		return nil, fmt.Errorf("error in ListPods: %v", err)
	}
	for _, pod := range pods {
		if err = k8s.Clientset.CoreV1().Pods(namespace).Delete(k8s.ctx, pod.Name, metav1.DeleteOptions{
			GracePeriodSeconds: &seconds,
		}); err != nil {
			return nil, fmt.Errorf("error in DeletePod: %v", err)
		}
	}
	// recreate job
	for i, container := range res.Spec.Template.Spec.Containers {
		if container.Name == containerName {
			container.Image = imageName
			res.Spec.Template.Spec.Containers[i] = container
			break
		}
	}
	for i, container := range res.Spec.Template.Spec.InitContainers {
		if container.Name == containerName {
			container.Image = imageName
			res.Spec.Template.Spec.Containers[i] = container
			break
		}
	}
	res.ObjectMeta.Annotations = nil
	res.ObjectMeta.Labels = nil
	res.ObjectMeta.ResourceVersion = ""
	res.Spec.Selector = nil
	res.Spec.Template.ObjectMeta = metav1.ObjectMeta{}
	// wait 10s for deletion of jobs
	time.Sleep(10 * time.Second)
	if res, err = k8s.Clientset.BatchV1().Jobs(namespace).Create(k8s.ctx, res, metav1.CreateOptions{}); err != nil {
		return nil, fmt.Errorf("error in CreateJob: %v", err)
	}
	return
}

func (k8s *K8sService) ScaleWorkload(namespace, workloadType, workloadName string, replicas uint32) (err error) {
	data := fmt.Sprintf(`{ "spec": { "replicas": %d } }`, replicas)
	switch workloadType {
	case constant.K8S_RESOURCE_TYPE_DEPLOYMENTS:
		if _, err = k8s.Clientset.AppsV1().Deployments(namespace).Patch(k8s.ctx, workloadName, ktypes.StrategicMergePatchType, []byte(data), metav1.PatchOptions{}); err != nil {
			return fmt.Errorf("error in ScaleDeployment: %v", err)
		}
	case constant.K8S_RESOURCE_TYPE_STATEFULSETS:
		if _, err = k8s.Clientset.AppsV1().StatefulSets(namespace).Patch(k8s.ctx, workloadName, ktypes.StrategicMergePatchType, []byte(data), metav1.PatchOptions{}); err != nil {
			return fmt.Errorf("error in ScaleStatefulSet: %v", err)
		}
	default:
		return fmt.Errorf("unknown workload type")
	}
	k8s.Logger.Infof("Scale %s[%s] replicas to %d success", workloadType, workloadName, replicas)
	return
}

func (k8s *K8sService) DeleteResource(namespace, resourceType, resourceName string, force bool) (err error) {
	switch resourceType {
	case constant.K8S_RESOURCE_TYPE_PODS:
		var period int64 = 0
		var propagationPolicy = metav1.DeletePropagationBackground
		deleteOptions := metav1.DeleteOptions{}
		if force {
			deleteOptions.GracePeriodSeconds = &period
			deleteOptions.PropagationPolicy = &propagationPolicy
		}
		return k8s.Clientset.CoreV1().Pods(namespace).Delete(k8s.ctx, resourceName, deleteOptions)
	case constant.K8S_RESOURCE_TYPE_SERVICE:
		return k8s.Clientset.CoreV1().Services(namespace).Delete(k8s.ctx, resourceName, metav1.DeleteOptions{})
	case constant.K8S_RESOURCE_TYPE_PVC:
		return k8s.Clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(k8s.ctx, resourceName, metav1.DeleteOptions{})
	case constant.K8S_RESOURCE_TYPE_JOB:
		return k8s.Clientset.BatchV1().Jobs(namespace).Delete(k8s.ctx, resourceName, metav1.DeleteOptions{})
	case constant.K8S_RESOURCE_TYPE_CRONJOB:
		return k8s.Clientset.BatchV1().CronJobs(namespace).Delete(k8s.ctx, resourceName, metav1.DeleteOptions{})
	case constant.K8S_RESOURCE_TYPE_INGRESSES:
		return k8s.Clientset.NetworkingV1().Ingresses(namespace).Delete(k8s.ctx, resourceName, metav1.DeleteOptions{})
	case constant.K8S_RESOURCE_TYPE_DEPLOYMENTS:
		return k8s.Clientset.AppsV1().Deployments(namespace).Delete(k8s.ctx, resourceName, metav1.DeleteOptions{})
	case constant.K8S_RESOURCE_TYPE_STATEFULSETS:
		return k8s.Clientset.AppsV1().StatefulSets(namespace).Delete(k8s.ctx, resourceName, metav1.DeleteOptions{})
	case constant.K8S_RESOURCE_TYPE_CONFIGMAPS:
		return k8s.Clientset.CoreV1().ConfigMaps(namespace).Delete(k8s.ctx, resourceName, metav1.DeleteOptions{})
	case constant.K8S_RESOURCE_TYPE_SECRETS:
		return k8s.Clientset.CoreV1().Secrets(namespace).Delete(k8s.ctx, resourceName, metav1.DeleteOptions{})
	case constant.K8S_RESOURCE_TYPE_SERVICEACCOUNTS:
		return k8s.Clientset.CoreV1().ServiceAccounts(namespace).Delete(k8s.ctx, resourceName, metav1.DeleteOptions{})
	default:
		return errorx.NewDefaultError("Unknown resource type")
	}
}

func (k8s *K8sService) GetReplicas(namespace, workloadType string, workloads []string) (workloadReplicas []commontypes.WorkloadReplica, err error) {
	for _, workload := range workloads {
		var replicas int32
		switch workloadType {
		case constant.K8S_RESOURCE_TYPE_DEPLOYMENTS:
			r, err := k8s.Clientset.AppsV1().Deployments(namespace).Get(k8s.ctx, workload, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			replicas = *r.Spec.Replicas
		case constant.K8S_RESOURCE_TYPE_STATEFULSETS:
			r, err := k8s.Clientset.AppsV1().StatefulSets(namespace).Get(k8s.ctx, workload, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			replicas = *r.Spec.Replicas
		default:
			return nil, errorx.NewDefaultError("Unknown resource type")
		}
		workloadReplicas = append(workloadReplicas, commontypes.WorkloadReplica{
			Name:     workload,
			Replicas: replicas,
		})
	}
	if workloadReplicas == nil {
		workloadReplicas = []commontypes.WorkloadReplica{}
	}
	return
}

func (k8s *K8sService) GetImages(namespace, workloadType string, initContainer bool, workloads []string) (workloadImages []commontypes.WorkloadImage, err error) {
	for _, workload := range workloads {
		var containers []corev1.Container
		var initContainers []corev1.Container
		switch workloadType {
		case constant.K8S_RESOURCE_TYPE_DEPLOYMENTS:
			r, err := k8s.Clientset.AppsV1().Deployments(namespace).Get(k8s.ctx, workload, metav1.GetOptions{})
			if errors.IsNotFound(err) {
				continue
			} else if err != nil {
				return nil, err
			}
			containers = r.Spec.Template.Spec.Containers
			initContainers = r.Spec.Template.Spec.InitContainers
		case constant.K8S_RESOURCE_TYPE_STATEFULSETS:
			r, err := k8s.Clientset.AppsV1().StatefulSets(namespace).Get(k8s.ctx, workload, metav1.GetOptions{})
			if errors.IsNotFound(err) {
				continue
			} else if err != nil {
				return nil, err
			}
			containers = r.Spec.Template.Spec.Containers
			initContainers = r.Spec.Template.Spec.InitContainers
		case constant.K8S_RESOURCE_TYPE_JOB:
			r, err := k8s.Clientset.BatchV1().Jobs(namespace).Get(k8s.ctx, workload, metav1.GetOptions{})
			if errors.IsNotFound(err) {
				continue
			} else if err != nil {
				return nil, err
			}
			containers = r.Spec.Template.Spec.Containers
			initContainers = r.Spec.Template.Spec.InitContainers
		case constant.K8S_RESOURCE_TYPE_CRONJOB:
			r, err := k8s.Clientset.BatchV1().CronJobs(namespace).Get(k8s.ctx, workload, metav1.GetOptions{})
			if errors.IsNotFound(err) {
				continue
			} else if err != nil {
				return nil, err
			}
			containers = r.Spec.JobTemplate.Spec.Template.Spec.Containers
			initContainers = r.Spec.JobTemplate.Spec.Template.Spec.InitContainers
		default:
			return nil, errorx.NewDefaultError("Unknown resource type")
		}
		var image string
		if initContainer {
			if len(initContainers) == 0 {
				return nil, errorx.NewDefaultError("There is no initContainers in workload \"%s\"", workload)
			}
			image = initContainers[0].Image
		} else {
			image = containers[0].Image
		}
		workloadImages = append(workloadImages, commontypes.WorkloadImage{
			Image: image,
			Name:  workload,
		})
	}
	if workloadImages == nil {
		workloadImages = []commontypes.WorkloadImage{}
	}
	return
}

func (k8s *K8sService) GetResourceStatus(namespace, resourceType, resourceName string) (res *commontypes.WorkloadStatus, err error) {
	var labels map[string]string
	var revision string
	ready := true
	switch resourceType {
	case constant.K8S_RESOURCE_TYPE_DEPLOYMENTS:
		r, err := k8s.Clientset.AppsV1().Deployments(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, errorx.NewDefaultError(err.Error())
		}
		labels = r.Spec.Selector.MatchLabels
	case constant.K8S_RESOURCE_TYPE_STATEFULSETS:
		r, err := k8s.Clientset.AppsV1().Deployments(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, errorx.NewDefaultError(err.Error())
		}
		labels = r.Spec.Selector.MatchLabels
	case constant.K8S_RESOURCE_TYPE_JOB:
		r, err := k8s.Clientset.BatchV1().Jobs(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, errorx.NewDefaultError(err.Error())
		}
		labels = r.Spec.Selector.MatchLabels
		for k, v := range r.GetAnnotations() {
			if strings.Contains(k, "kubernetes.io/revision") {
				revision = v
				break
			}
		}
	case constant.K8S_RESOURCE_TYPE_CONFIGMAPS:
		r, err := k8s.Clientset.CoreV1().ConfigMaps(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, errorx.NewDefaultError(err.Error())
		}
		return &commontypes.WorkloadStatus{
			Name:   r.GetName(),
			Status: constant.STATUS_SUCCESS,
			Ready:  true,
		}, nil
	case constant.K8S_RESOURCE_TYPE_SERVICE:
		r, err := k8s.Clientset.CoreV1().Services(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, errorx.NewDefaultError(err.Error())
		}
		return &commontypes.WorkloadStatus{
			Name:   r.GetName(),
			Status: constant.STATUS_SUCCESS,
			Ready:  true,
		}, nil
	case constant.K8S_RESOURCE_TYPE_INGRESSES:
		r, err := k8s.Clientset.NetworkingV1().Ingresses(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, errorx.NewDefaultError(err.Error())
		}
		return &commontypes.WorkloadStatus{
			Name:   r.GetName(),
			Status: constant.STATUS_SUCCESS,
			Ready:  true,
		}, nil
	case constant.K8S_RESOURCE_TYPE_SECRETS:
		r, err := k8s.Clientset.CoreV1().Secrets(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, errorx.NewDefaultError(err.Error())
		}
		return &commontypes.WorkloadStatus{
			Name:   r.GetName(),
			Status: constant.STATUS_SUCCESS,
			Ready:  true,
		}, nil
	case constant.K8S_RESOURCE_TYPE_PVC:
		r, err := k8s.Clientset.CoreV1().PersistentVolumeClaims(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, errorx.NewDefaultError(err.Error())
		}
		if r.Status.Phase != corev1.ClaimBound {
			ready = false
		}
		return &commontypes.WorkloadStatus{
			Name:   r.GetName(),
			Status: string(r.Status.Phase),
			Ready:  ready,
		}, nil
	}
	Pods, _, err := k8s.ListPods(namespace, labels, "", "", 0)
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	var pods []commontypes.PodStatus

	for _, v := range Pods {
		readyStatus := "False"
		if resourceType == constant.K8S_RESOURCE_TYPE_JOB {
			if v.Status.Phase == "Succeeded" {
				readyStatus = "True"
			}
		} else {
			for _, c := range v.Status.Conditions {
				if c.Type == "Ready" {
					readyStatus = string(c.Status)
					break
				}
			}
		}
		image := make(map[string]string)
		for _, container := range v.Spec.Containers {
			image[container.Name] = container.Image
		}
		for _, container := range v.Spec.InitContainers {
			image[container.Name] = container.Image
		}
		podstatus := commontypes.PodStatus{
			PodName: v.Name,
			Ready:   readyStatus,
			Image:   image,
		}
		pods = append(pods, podstatus)
		if readyStatus == "False" {
			ready = false
		}
	}
	workLoadStatus := commontypes.WorkloadStatus{
		Name:     resourceName,
		Revision: revision,
		Pods:     pods,
		Ready:    ready,
	}
	return &workLoadStatus, nil
}

func (k8s *K8sService) ListPods(namespace string, labels map[string]string, fieldSelector, conti string, limit int64) (res []corev1.Pod, continu string, err error) {
	if limit == 0 {
		limit = 500
	}
	delete(labels, "pod-template-hash")
	var labelSelector []string
	for k, v := range labels {
		labelSelector = append(labelSelector, k+"="+v)
	}
	var podRes *corev1.PodList
	if podRes, err = k8s.Clientset.CoreV1().Pods(namespace).List(k8s.ctx, metav1.ListOptions{
		LabelSelector: strings.Join(labelSelector, ","),
		FieldSelector: fieldSelector,
		Limit:         limit,
		Continue:      conti,
	}); err != nil {
		k8s.Logger.Error(err)
		return
	}
	for i := range podRes.Items {
		podRes.Items[i].ManagedFields = nil
		podRes.Items[i].ObjectMeta.OwnerReferences = nil
	}
	// k8s.Logger.Debugf("Get pods with namespace: %s, labelSelector: %s, items=%+v", namespace, strings.Join(labelSelector, ","), res.Items)
	return podRes.Items, podRes.Continue, nil
}

func (k8s *K8sService) GetEvents(namespace, objectKind, objectName string) ([]corev1.Event, error) {
	res, err := k8s.Clientset.CoreV1().Events(namespace).List(k8s.ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=%s", objectName, objectKind),
	})
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}
	return res.Items, nil
}

func (k8s *K8sService) GetResourceYAML(namespace, resourceType, resourceNames string) (yamlContent string, err error) {
	var yamlItems []string
	for _, resourceName := range strings.Split(resourceNames, ",") {
		var unstructuredMap map[string]interface{}
		switch resourceType {
		case constant.K8S_RESOURCE_TYPE_PODS:
			res, err := k8s.Clientset.CoreV1().Pods(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				k8s.Logger.Error(err)
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.K8S_RESOURCE_TYPE_SERVICE:
			res, err := k8s.Clientset.CoreV1().Services(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				k8s.Logger.Error(err)
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.K8S_RESOURCE_TYPE_PVC:
			res, err := k8s.Clientset.CoreV1().PersistentVolumeClaims(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				k8s.Logger.Error(err)
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.K8S_RESOURCE_TYPE_DEPLOYMENTS:
			res, err := k8s.Clientset.AppsV1().Deployments(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				k8s.Logger.Error(err)
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.K8S_RESOURCE_TYPE_STATEFULSETS:
			res, err := k8s.Clientset.AppsV1().StatefulSets(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				k8s.Logger.Error(err)
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.K8S_RESOURCE_TYPE_JOB:
			res, err := k8s.Clientset.BatchV1().Jobs(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				k8s.Logger.Error(err)
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.K8S_RESOURCE_TYPE_CRONJOB:
			res, err := k8s.Clientset.BatchV1().CronJobs(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				k8s.Logger.Error(err)
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.K8S_RESOURCE_TYPE_INGRESSES:
			res, err := k8s.Clientset.NetworkingV1().Ingresses(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				k8s.Logger.Error(err)
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.K8S_RESOURCE_TYPE_CONFIGMAPS:
			res, err := k8s.Clientset.CoreV1().ConfigMaps(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				k8s.Logger.Error(err)
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.K8S_RESOURCE_TYPE_SECRETS:
			res, err := k8s.Clientset.CoreV1().Secrets(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				k8s.Logger.Error(err)
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		case constant.K8S_RESOURCE_TYPE_SERVICEACCOUNTS:
			res, err := k8s.Clientset.CoreV1().ServiceAccounts(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				k8s.Logger.Error(err)
				return "", err
			}
			res.ObjectMeta = processObjectMeta(res.ObjectMeta)
			unstructuredMap, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(res)
		default:
			return "", errorx.NewDefaultError("Unknown resource type")
		}
		delete(unstructuredMap, "status")
		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
		unstructuredObj.SetKind(getKindVersion(resourceType)[0])
		unstructuredObj.SetAPIVersion(getKindVersion(resourceType)[1])
		annotations := unstructuredObj.GetAnnotations()
		delete(annotations, "kubectl.kubernetes.io/last-applied-configuration")
		unstructuredObj.SetAnnotations(annotations)
		yamlPrinter := &printers.YAMLPrinter{}
		var buf bytes.Buffer
		if err = yamlPrinter.PrintObj(unstructuredObj, &buf); err != nil {
			k8s.Logger.Error(err)
			return "", err
		}
		yamlItems = append(yamlItems, buf.String())
		k8s.Logger.Debugf("\n%s", buf.String())
	}
	return strings.Join(yamlItems, "---\n"), nil
}

func (k8s *K8sService) GetConfigMap(namespace string, configmapName string) (map[string]string, error) {
	res, err := k8s.Clientset.CoreV1().ConfigMaps(namespace).Get(k8s.ctx, configmapName, metav1.GetOptions{})
	if err != nil {
		k8s.Logger.Error(err)
		return nil, errorx.NewDefaultError(err.Error())
	}
	return res.Data, nil
}

func (k8s *K8sService) PatchConfig(namespace, resourceType, resourceName, key string, value string) (res []byte, err error) {
	data := fmt.Sprintf(`{ "data": { "%s": %q } }`, key, value)
	switch resourceType {
	case constant.K8S_RESOURCE_TYPE_CONFIGMAPS:
		r, err := k8s.Clientset.CoreV1().ConfigMaps(namespace).Patch(k8s.ctx, resourceName, ktypes.StrategicMergePatchType, []byte(data), metav1.PatchOptions{})
		if err != nil {
			k8s.Logger.Error(err)
			return nil, errorx.NewDefaultError(err.Error())
		}
		res, _ = json.Marshal(r)
	case constant.K8S_RESOURCE_TYPE_SECRETS:
		r, err := k8s.Clientset.CoreV1().Secrets(namespace).Patch(k8s.ctx, resourceName, ktypes.StrategicMergePatchType, []byte(data), metav1.PatchOptions{})
		if err != nil {
			k8s.Logger.Error(err)
			return nil, errorx.NewDefaultError(err.Error())
		}
		res, _ = json.Marshal(r)
	default:
		return nil, errorx.NewDefaultError("Unknown resource type")
	}
	return res, nil
}

func (k8s *K8sService) UpdateFromYaml(namespace, applyYaml, kind string) (err error) {
	k8s.Logger.Debugf("\n%s", applyYaml)
	d := uyaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(applyYaml), 4096)
	var unstructureObj *unstructured.Unstructured
	for {
		unstructureObj, err = utils.GetUnstructured(d)
		if err == io.EOF {
			break
		}
		if err != nil {
			k8s.Logger.Error(err)
			return
		}
		if namespace != unstructureObj.GetNamespace() {
			return errorx.NewDefaultError("Namespace must be %s", namespace)
		}
		if kind != "" && kind != unstructureObj.GetKind() {
			return errorx.NewDefaultError("Kind must be %s", kind)
		}
		var gvr schema.GroupVersionResource
		if gvr, err = k8s.gtGVR(unstructureObj.GroupVersionKind()); err != nil {
			return
		}
		foundObj, getErr := k8s.dynamicclient.Resource(gvr).Namespace(namespace).Get(context.Background(), unstructureObj.GetName(), metav1.GetOptions{})
		if getErr != nil { // not found, create resource
			if _, err = k8s.dynamicclient.Resource(gvr).Namespace(namespace).Create(context.Background(), unstructureObj, metav1.CreateOptions{}); err != nil {
				return
			}
			name := unstructureObj.GetName()
			if name == "" {
				name = unstructureObj.GetGenerateName()
			}
			k8s.Logger.Infof("Create resource[%s] success", name)
		} else { // found & update resource
			unstructureObj.SetResourceVersion(foundObj.GetResourceVersion())
			if _, err = k8s.dynamicclient.Resource(gvr).Namespace(namespace).Update(context.Background(), unstructureObj, metav1.UpdateOptions{}); err != nil {
				k8s.Logger.Errorf("unable to update resource[%s]: %+v", unstructureObj.GetName(), err)
				return errorx.NewDefaultError("unable to update resource[%s]: %v", unstructureObj.GetName(), err)
			}
			k8s.Logger.Infof("Update resource[%s] success", unstructureObj.GetName())
		}
	}
	return nil
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
			err := k8s.dynamicclient.Resource(gvr).Namespace(namespace).Delete(context.Background(), unstructureObj.GetName(), metav1.DeleteOptions{})
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
func (k8s *K8sService) RolloutWorkload(namespace, workloadType, workloadName string) (res []byte, err error) {
	data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().Format("2006-01-02T15:04:05Z"))
	switch workloadType {
	case constant.K8S_RESOURCE_TYPE_DEPLOYMENTS:
		r, err := k8s.Clientset.AppsV1().Deployments(namespace).Patch(k8s.ctx, workloadName, ktypes.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
		if err != nil {
			return nil, err
		}
		res, _ = json.Marshal(r)
	case constant.K8S_RESOURCE_TYPE_STATEFULSETS:
		r, err := k8s.Clientset.AppsV1().StatefulSets(namespace).Patch(k8s.ctx, workloadName, ktypes.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
		if err != nil {
			return nil, err
		}
		res, _ = json.Marshal(r)
	case constant.K8S_RESOURCE_TYPE_JOB:
		// jobs image is immutable, must delete job first and recreate
		var r *batchv1.Job
		if r, err = k8s.Clientset.BatchV1().Jobs(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{}); err != nil {
			return nil, fmt.Errorf("error in GetJob: %w", err)
		}
		var seconds int64 = 0
		if err = k8s.Clientset.BatchV1().Jobs(namespace).Delete(k8s.ctx, workloadName, metav1.DeleteOptions{GracePeriodSeconds: &seconds}); err != nil {
			return nil, fmt.Errorf("error in DeleteJob: %w", err)
		}
		time.Sleep(1 * time.Second)
		r.ObjectMeta.Annotations = nil
		r.ObjectMeta.Labels = nil
		r.ObjectMeta.ResourceVersion = ""
		r.Spec.Selector = nil
		r.Spec.Template.ObjectMeta = metav1.ObjectMeta{}
		if r, err = k8s.Clientset.BatchV1().Jobs(namespace).Create(k8s.ctx, r, metav1.CreateOptions{}); err != nil {
			return nil, fmt.Errorf("error in CreateJob: %w", err)
		}
		res, _ = json.Marshal(r)
	default:
		return nil, errorx.NewDefaultError("Unknown resource type")
	}
	return
}

func (k8s *K8sService) GetNamespaces(labelSelector string) (res *corev1.NamespaceList, err error) {
	if res, err = k8s.Clientset.CoreV1().Namespaces().List(k8s.ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	}); err != nil {
		k8s.Logger.Error(err)
	}
	return
}

func (k8s *K8sService) ListResource(namespace, resourceType, labelSelector, fieldSelector, conti string, limit int64) (res []byte, err error) {
	if limit == 0 {
		limit = 500
	}
	switch resourceType {
	case constant.K8S_RESOURCE_TYPE_DEPLOYMENTS:
		var r *v1.DeploymentList
		if r, err = k8s.Clientset.AppsV1().Deployments(namespace).List(k8s.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListDeployment: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
			r.Items[i].Spec = v1.DeploymentSpec{}
			delete(r.Items[i].Annotations, "kubectl.kubernetes.io/last-applied-configuration")
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	case constant.K8S_RESOURCE_TYPE_STATEFULSETS:
		var r *v1.StatefulSetList
		if r, err = k8s.Clientset.AppsV1().StatefulSets(namespace).List(k8s.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListStatefulSet: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
			r.Items[i].Spec = v1.StatefulSetSpec{}
			delete(r.Items[i].Annotations, "kubectl.kubernetes.io/last-applied-configuration")
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	case constant.K8S_RESOURCE_TYPE_JOB:
		var r *batchv1.JobList
		if r, err = k8s.Clientset.BatchV1().Jobs(namespace).List(k8s.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListJob: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
			r.Items[i].Spec = batchv1.JobSpec{}
			delete(r.Items[i].Annotations, "kubectl.kubernetes.io/last-applied-configuration")
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	case constant.K8S_RESOURCE_TYPE_CRONJOB:
		var r *batchv1.CronJobList
		if r, err = k8s.Clientset.BatchV1().CronJobs(namespace).List(k8s.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListCronJob: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
			delete(r.Items[i].Annotations, "kubectl.kubernetes.io/last-applied-configuration")
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	case constant.K8S_RESOURCE_TYPE_PVC:
		var r *corev1.PersistentVolumeClaimList
		if r, err = k8s.Clientset.CoreV1().PersistentVolumeClaims(namespace).List(k8s.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListPVC: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	case constant.K8S_RESOURCE_TYPE_SERVICE:
		var r *corev1.ServiceList
		if r, err = k8s.Clientset.CoreV1().Services(namespace).List(k8s.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListService: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	case constant.K8S_RESOURCE_TYPE_INGRESSES:
		var r *networkingv1.IngressList
		if r, err = k8s.Clientset.NetworkingV1().Ingresses(namespace).List(k8s.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListIngress: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	case constant.K8S_RESOURCE_TYPE_PODS:
		labels := make(map[string]string)
		if labelSelector != "" {
			for _, ls := range strings.Split(labelSelector, ",") {
				labelParts := strings.Split(ls, "=")
				labels[labelParts[0]] = labelParts[1]
			}
		}
		r, conti, err := k8s.ListPods(namespace, labels, fieldSelector, conti, limit)
		if err != nil {
			return nil, fmt.Errorf("error in ListPod: %w", err)
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r,
			Continue: conti,
		})
	case constant.K8S_RESOURCE_TYPE_CONFIGMAPS:
		var r *corev1.ConfigMapList
		if r, err = k8s.Clientset.CoreV1().ConfigMaps(namespace).List(k8s.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListConfigMaps: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	case constant.K8S_RESOURCE_TYPE_SECRETS:
		var r *corev1.SecretList
		if r, err = k8s.Clientset.CoreV1().Secrets(namespace).List(k8s.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListSecrets: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	case constant.K8S_RESOURCE_TYPE_SERVICEACCOUNTS:
		var r *corev1.ServiceAccountList
		if r, err = k8s.Clientset.CoreV1().ServiceAccounts(namespace).List(k8s.ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
			Continue:      conti,
			Limit:         limit,
		}); err != nil {
			return nil, fmt.Errorf("error in ListServiceAccounts: %w", err)
		}
		for i := range r.Items {
			r.Items[i].ManagedFields = nil
		}
		res, _ = json.Marshal(commontypes.ListResponseData{
			Results:  r.Items,
			Continue: r.Continue,
		})
	default:
		return nil, errorx.NewDefaultError("Unknown resource type")
	}
	return
}

func (k8s *K8sService) GetResource(namespace, resourceType, resourceName string) (res []byte, labels map[string]string, err error) {
	switch resourceType {
	case constant.K8S_RESOURCE_TYPE_DEPLOYMENTS:
		r, err := k8s.Clientset.AppsV1().Deployments(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("error in GetDeployment: %w", err)
		}
		r.ManagedFields = nil
		labels = r.Spec.Selector.MatchLabels
		res, _ = json.Marshal(r)
	case constant.K8S_RESOURCE_TYPE_STATEFULSETS:
		r, err := k8s.Clientset.AppsV1().StatefulSets(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("error in GetStatefulSet: %w", err)
		}
		r.ManagedFields = nil
		labels = r.Spec.Selector.MatchLabels
		res, _ = json.Marshal(r)
	case constant.K8S_RESOURCE_TYPE_JOB:
		r, err := k8s.Clientset.BatchV1().Jobs(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("error in GetJob: %w", err)
		}
		r.ManagedFields = nil
		labels = r.Spec.Selector.MatchLabels
		res, _ = json.Marshal(r)
	case constant.K8S_RESOURCE_TYPE_CRONJOB:
		r, err := k8s.Clientset.BatchV1().CronJobs(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("error in GetCronJob: %w", err)
		}
		r.ManagedFields = nil
		labels = r.GetLabels()
		res, _ = json.Marshal(r)
	case constant.K8S_RESOURCE_TYPE_PODS:
		r, err := k8s.Clientset.CoreV1().Pods(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("error in GetPod: %w", err)
		}
		r.ManagedFields = nil
		labels = r.GetLabels()
		res, _ = json.Marshal(r)
	case constant.K8S_RESOURCE_TYPE_PVC:
		r, err := k8s.Clientset.CoreV1().PersistentVolumeClaims(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("error in GetPVC: %w", err)
		}
		r.ManagedFields = nil
		labels = r.GetLabels()
		res, _ = json.Marshal(r)
	case constant.K8S_RESOURCE_TYPE_SERVICE:
		r, err := k8s.Clientset.CoreV1().Services(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("error in GetService: %w", err)
		}
		r.ManagedFields = nil
		labels = r.Spec.Selector
		res, _ = json.Marshal(r)
	case constant.K8S_RESOURCE_TYPE_INGRESSES:
		r, err := k8s.Clientset.NetworkingV1().Ingresses(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("error in GetIngress: %w", err)
		}
		r.ManagedFields = nil
		labels = r.GetLabels()
		res, _ = json.Marshal(r)
	case constant.K8S_RESOURCE_TYPE_CONFIGMAPS:
		r, err := k8s.Clientset.CoreV1().ConfigMaps(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("error in GetConfigMap: %w", err)
		}
		r.ManagedFields = nil
		labels = r.GetLabels()
		res, _ = json.Marshal(r)
	case constant.K8S_RESOURCE_TYPE_SECRETS:
		r, err := k8s.Clientset.CoreV1().Secrets(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("error in GetSecret: %w", err)
		}
		r.ManagedFields = nil
		labels = r.GetLabels()
		res, _ = json.Marshal(r)
	case constant.K8S_RESOURCE_TYPE_SERVICEACCOUNTS:
		r, err := k8s.Clientset.CoreV1().ServiceAccounts(namespace).Get(k8s.ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("error in GetServiceAccount: %w", err)
		}
		r.ManagedFields = nil
		labels = r.GetLabels()
		res, _ = json.Marshal(r)
	default:
		return nil, nil, errorx.NewDefaultError("Unknown resource type")
	}
	return
}

func (k8s *K8sService) GetWorkloadQuota(namespace, workloadType, workloadName string) (*corev1.ResourceRequirements, error) {
	switch workloadType {
	case constant.K8S_RESOURCE_TYPE_DEPLOYMENTS:
		res, err := k8s.Clientset.AppsV1().Deployments(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return &res.Spec.Template.Spec.Containers[0].Resources, nil
	case constant.K8S_RESOURCE_TYPE_STATEFULSETS:
		res, err := k8s.Clientset.AppsV1().StatefulSets(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return &res.Spec.Template.Spec.Containers[0].Resources, nil
	default:
		return nil, errorx.NewDefaultError("Unknown resource type")
	}
}

func (k8s *K8sService) GetPodLog(namespace, podname, containerName string, lines int64, follow, timestamps bool, stream agent.LizardAgent_GetPodLogFollowServer, sendMsg func([]byte)) (logs string, err error) {
	podLogOption := corev1.PodLogOptions{
		Container:  containerName,
		Follow:     follow,
		Timestamps: timestamps,
	}
	if lines != 0 {
		podLogOption.TailLines = &lines
	}
	req := k8s.Clientset.CoreV1().Pods(namespace).GetLogs(podname, &podLogOption)
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
			if stream != nil {
				if err = stream.Send(&agent.YamlResponse{
					Code: uint32(codes.OK),
					Data: string(b),
				}); err != nil {
					k8s.Logger.Error(err)
					return
				}
			} else {
				sendMsg(b)
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
	res, err := k8s.Clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return k8s.Clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).Create(k8s.ctx, hpa, metav1.CreateOptions{})
	} else {
		hpa.ResourceVersion = res.ResourceVersion
		return k8s.Clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).Update(k8s.ctx, hpa, metav1.UpdateOptions{})
	}
}

func (k8s *K8sService) GetPodHorizonAutoscaler(namespace, workloadName string) (*autov2.HorizontalPodAutoscaler, error) {
	res, err := k8s.Clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).Get(k8s.ctx, workloadName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, err
	}
	res.ManagedFields = nil
	return res, nil
}

func (k8s *K8sService) gtGVR(gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {
	gr, err := restmapper.GetAPIGroupResources(k8s.Clientset.Discovery())
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

func processObjectMeta(metadata metav1.ObjectMeta) metav1.ObjectMeta {
	metadata.ManagedFields = nil
	metadata.UID = ""
	metadata.ResourceVersion = ""
	metadata.Generation = 0
	metadata.OwnerReferences = nil
	return metadata
}

func getKindVersion(resourceType string) []string {
	resourceMap := map[string][]string{
		constant.K8S_RESOURCE_TYPE_DEPLOYMENTS:     {"Deployment", "apps/v1"},
		constant.K8S_RESOURCE_TYPE_STATEFULSETS:    {"StatefulSet", "apps/v1"},
		constant.K8S_RESOURCE_TYPE_PODS:            {"Pod", "v1"},
		constant.K8S_RESOURCE_TYPE_SERVICE:         {"Service", "v1"},
		constant.K8S_RESOURCE_TYPE_PVC:             {"PersistentVolumeClaim", "v1"},
		constant.K8S_RESOURCE_TYPE_CRONJOB:         {"CronJob", "batch/v1"},
		constant.K8S_RESOURCE_TYPE_JOB:             {"Job", "batch/v1"},
		constant.K8S_RESOURCE_TYPE_INGRESSES:       {"Ingress", "networking.k8s.io/v1"},
		constant.K8S_RESOURCE_TYPE_CONFIGMAPS:      {"ConfigMap", "v1"},
		constant.K8S_RESOURCE_TYPE_SECRETS:         {"Secret", "v1"},
		constant.K8S_RESOURCE_TYPE_SERVICEACCOUNTS: {"ServiceAccount", "v1"},
		constant.TEKTON_CRD_TYPE_TASKS:             {"Task", "tekton.dev/v1"},
		constant.TEKTON_CRD_TYPE_TASKRUNS:          {"TaskRun", "tekton.dev/v1"},
		constant.TEKTON_CRD_TYPE_PIPELINES:         {"Pipeline", "tekton.dev/v1"},
		constant.TEKTON_CRD_TYPE_PIPELINERUNS:      {"PipelineRun", "tekton.dev/v1"},
		constant.TEKTON_CRD_TYPE_TRIGGERBINDINGS:   {"TriggerBinding", "triggers.tekton.dev/v1beta1"},
		constant.TEKTON_CRD_TYPE_TRIGGERTEMPLATES:  {"TriggerTemplates", "triggers.tekton.dev/v1beta1"},
		constant.TEKTON_CRD_TYPE_EVENTLISTENERS:    {"EventListener", "triggers.tekton.dev/v1beta1"},
		constant.TEKTON_CRD_TYPE_APPROVALTASKS:     {"ApprovalTask", "tekton.automatiko.io/v1beta1"},
	}
	if v, ok := resourceMap[resourceType]; ok {
		return v
	} else {
		return []string{"unknown", "unknown"}
	}
}

func getPatchWorkloadData(init bool, containers []corev1.Container, containerName, imageName string) string {
	for _, container := range containers {
		if container.Name == containerName {
			if init {
				return fmt.Sprintf(`{ "spec": { "template": { "spec": { "initContainers": [ { "name": "%s", "image": "%s" } ] } } } }`, containerName, imageName)
			}
			return fmt.Sprintf(`{ "spec": { "template": { "spec": { "containers": [ { "name": "%s", "image": "%s" } ] } } } }`, containerName, imageName)
		}
	}
	return ""
}
