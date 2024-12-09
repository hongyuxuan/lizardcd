package svc

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/golang-module/carbon"
	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
)

func WritePodData(cluster, namespace, workloadType, workloadName string) (data [][]string) {
	var res *types.PodRes
	if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/kubernetes/cluster/%s/namespace/%s/%s/%s/pods", cluster, namespace, workloadType, workloadName)).SetResult(&res).Do(context.Background()).Err; err != nil {
		common.PrintFatal("failed to get pods of cluster=%s, namespace=%s, %s=%s: %v\n", cluster, namespace, workloadType, workloadName, err)
	}

	for _, d := range res.Data {
		var row []string
		var stateReason string
		var message string
		if d.Status.ContainerStatuses != nil {
			state := d.Status.ContainerStatuses[0].State
			if state.Running != nil {
				stateReason = fmt.Sprintf("Created %s", carbon.FromStdTime(state.Running.StartedAt.Time).DiffForHumans())
			} else if state.Waiting != nil {
				stateReason = state.Waiting.Reason
				message = state.Waiting.Message
			}
			row = []string{d.Name, d.Status.PodIP, fmt.Sprintf("%s(%s)", d.Spec.NodeName, d.Status.HostIP), stateReason, message}
		} else {
			stateReason = d.Status.Conditions[0].Reason
			message = d.Status.Conditions[0].Message
			row = []string{d.Name, "", "", stateReason, message}
		}
		data = append(data, row)
	}
	return
}

func GetPod(cluster, namespace, workloadType, workloadName, podName string) corev1.Pod {
	var res *types.PodRes
	if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/kubernetes/cluster/%s/namespace/%s/%s/%s/pods", cluster, namespace, workloadType, workloadName)).SetResult(&res).Do(context.Background()).Err; err != nil {
		common.PrintFatal("failed to get statefulset pods of cluster=%s, namespace=%s, %s=%s: %v\n", cluster, namespace, workloadType, workloadName, err)
	}
	p, ok := lo.Find(res.Data, func(item corev1.Pod) bool {
		return item.Name == podName
	})
	if !ok {
		common.PrintFatal("cannot find pod=%s in cluster=%s namespace=%s %s=%s\n", podName, cluster, namespace, workloadType, workloadName)
	}
	return p
}

func WriteContainerData(containerStatus []corev1.ContainerStatus, containers []corev1.Container, podCondition []corev1.PodCondition, initContainer bool) (data [][]string) {
	var row []string
	var stateReason string
	if containerStatus != nil {
		for _, c := range containerStatus {
			stateReason = c.Image
			var status string
			if c.State.Running != nil {
				status = "running"
			} else if c.State.Waiting != nil {
				stateReason = c.State.Waiting.Reason
				status = "waiting"
			} else if c.State.Terminated != nil {
				stateReason = c.State.Terminated.Reason
				status = "terminated"
			}
			containerName := c.Name
			if initContainer {
				containerName += "(initContainer)"
			}
			row = []string{containerName, status, stateReason, strconv.Itoa(int(c.RestartCount))}
			data = append(data, row)
		}
	} else {
		for _, c := range containers {
			containerName := c.Name
			if initContainer {
				containerName += "(initContainer)"
			}
			row = []string{containerName, "waiting", podCondition[0].Reason, ""}
			data = append(data, row)
		}
	}
	return
}

func GetEvents(cluster, namespace, podName string) (data [][]string) {
	var res *types.EventRes
	if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/kubernetes/cluster/%s/namespace/%s/Pod/%s/events", cluster, namespace, podName)).SetResult(&res).Do(context.Background()).Err; err != nil {
		common.PrintFatal("failed to get events of cluster=%s, namespace=%s, pod=%s: %v\n", cluster, namespace, podName, err)
	}
	for _, d := range res.Data {
		row := []string{d.Type, d.Reason, carbon.FromStdTime(d.EventTime.Time).Format("Y-m-d H:i:s"), d.Source.Component, d.Message}
		data = append(data, row)
	}
	return
}

func SortData(data [][]string, sortBy, order string) {
	sort.Slice(data, func(i, j int) bool {
		if sortBy == "lastUpdateTime" {
			if order == "desc" {
				return data[i][2] > data[j][2] // data[i][1] is lastUpdateTime
			} else {
				return data[i][2] < data[j][2]
			}
		} else {
			if order == "desc" {
				return data[i][0] > data[j][0]
			} else {
				return data[i][0] < data[j][0]
			}
		}
	})
}
