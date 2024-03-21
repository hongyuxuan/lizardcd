// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.19.4
// source: agent.proto

package agent

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PatchWorkloadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace    string `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	WorkloadName string `protobuf:"bytes,2,opt,name=workload_name,json=workloadName,proto3" json:"workload_name,omitempty"`
	Container    string `protobuf:"bytes,3,opt,name=container,proto3" json:"container,omitempty"`
	Image        string `protobuf:"bytes,4,opt,name=image,proto3" json:"image,omitempty"`
}

func (x *PatchWorkloadRequest) Reset() {
	*x = PatchWorkloadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PatchWorkloadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PatchWorkloadRequest) ProtoMessage() {}

func (x *PatchWorkloadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PatchWorkloadRequest.ProtoReflect.Descriptor instead.
func (*PatchWorkloadRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{0}
}

func (x *PatchWorkloadRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *PatchWorkloadRequest) GetWorkloadName() string {
	if x != nil {
		return x.WorkloadName
	}
	return ""
}

func (x *PatchWorkloadRequest) GetContainer() string {
	if x != nil {
		return x.Container
	}
	return ""
}

func (x *PatchWorkloadRequest) GetImage() string {
	if x != nil {
		return x.Image
	}
	return ""
}

type ListWorkloadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace     string `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	LabelSelector string `protobuf:"bytes,2,opt,name=label_selector,json=labelSelector,proto3" json:"label_selector,omitempty"`
}

func (x *ListWorkloadRequest) Reset() {
	*x = ListWorkloadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListWorkloadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListWorkloadRequest) ProtoMessage() {}

func (x *ListWorkloadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListWorkloadRequest.ProtoReflect.Descriptor instead.
func (*ListWorkloadRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{1}
}

func (x *ListWorkloadRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *ListWorkloadRequest) GetLabelSelector() string {
	if x != nil {
		return x.LabelSelector
	}
	return ""
}

type GetWorkloadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace    string `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	WorkloadName string `protobuf:"bytes,2,opt,name=workload_name,json=workloadName,proto3" json:"workload_name,omitempty"`
}

func (x *GetWorkloadRequest) Reset() {
	*x = GetWorkloadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetWorkloadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetWorkloadRequest) ProtoMessage() {}

func (x *GetWorkloadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetWorkloadRequest.ProtoReflect.Descriptor instead.
func (*GetWorkloadRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{2}
}

func (x *GetWorkloadRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *GetWorkloadRequest) GetWorkloadName() string {
	if x != nil {
		return x.WorkloadName
	}
	return ""
}

type GetPodEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace string `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	PodName   string `protobuf:"bytes,2,opt,name=pod_name,json=podName,proto3" json:"pod_name,omitempty"`
}

func (x *GetPodEventRequest) Reset() {
	*x = GetPodEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPodEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPodEventRequest) ProtoMessage() {}

func (x *GetPodEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPodEventRequest.ProtoReflect.Descriptor instead.
func (*GetPodEventRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{3}
}

func (x *GetPodEventRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *GetPodEventRequest) GetPodName() string {
	if x != nil {
		return x.PodName
	}
	return ""
}

type YamlRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace string `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Ymlstring string `protobuf:"bytes,2,opt,name=ymlstring,proto3" json:"ymlstring,omitempty"`
}

func (x *YamlRequest) Reset() {
	*x = YamlRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *YamlRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*YamlRequest) ProtoMessage() {}

func (x *YamlRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use YamlRequest.ProtoReflect.Descriptor instead.
func (*YamlRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{4}
}

func (x *YamlRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *YamlRequest) GetYmlstring() string {
	if x != nil {
		return x.Ymlstring
	}
	return ""
}

type RolloutWorkloadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace    string `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	WorkloadName string `protobuf:"bytes,2,opt,name=workload_name,json=workloadName,proto3" json:"workload_name,omitempty"`
}

func (x *RolloutWorkloadRequest) Reset() {
	*x = RolloutWorkloadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RolloutWorkloadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RolloutWorkloadRequest) ProtoMessage() {}

func (x *RolloutWorkloadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RolloutWorkloadRequest.ProtoReflect.Descriptor instead.
func (*RolloutWorkloadRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{5}
}

func (x *RolloutWorkloadRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *RolloutWorkloadRequest) GetWorkloadName() string {
	if x != nil {
		return x.WorkloadName
	}
	return ""
}

type GetYamlRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace    string `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	ResourceType string `protobuf:"bytes,2,opt,name=resource_type,json=resourceType,proto3" json:"resource_type,omitempty"`
	ResourceName string `protobuf:"bytes,3,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
}

func (x *GetYamlRequest) Reset() {
	*x = GetYamlRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetYamlRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetYamlRequest) ProtoMessage() {}

func (x *GetYamlRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetYamlRequest.ProtoReflect.Descriptor instead.
func (*GetYamlRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{6}
}

func (x *GetYamlRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *GetYamlRequest) GetResourceType() string {
	if x != nil {
		return x.ResourceType
	}
	return ""
}

func (x *GetYamlRequest) GetResourceName() string {
	if x != nil {
		return x.ResourceName
	}
	return ""
}

type ScaleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace    string `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	WorkloadName string `protobuf:"bytes,2,opt,name=workload_name,json=workloadName,proto3" json:"workload_name,omitempty"`
	Replicas     uint32 `protobuf:"varint,3,opt,name=replicas,proto3" json:"replicas,omitempty"`
}

func (x *ScaleRequest) Reset() {
	*x = ScaleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScaleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScaleRequest) ProtoMessage() {}

func (x *ScaleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScaleRequest.ProtoReflect.Descriptor instead.
func (*ScaleRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{7}
}

func (x *ScaleRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *ScaleRequest) GetWorkloadName() string {
	if x != nil {
		return x.WorkloadName
	}
	return ""
}

func (x *ScaleRequest) GetReplicas() uint32 {
	if x != nil {
		return x.Replicas
	}
	return 0
}

type LabelSelector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LabelSelector string `protobuf:"bytes,1,opt,name=label_selector,json=labelSelector,proto3" json:"label_selector,omitempty"`
}

func (x *LabelSelector) Reset() {
	*x = LabelSelector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LabelSelector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LabelSelector) ProtoMessage() {}

func (x *LabelSelector) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LabelSelector.ProtoReflect.Descriptor instead.
func (*LabelSelector) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{8}
}

func (x *LabelSelector) GetLabelSelector() string {
	if x != nil {
		return x.LabelSelector
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{9}
}

func (x *Response) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Response) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Response) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type YamlResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Data string `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *YamlResponse) Reset() {
	*x = YamlResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_agent_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *YamlResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*YamlResponse) ProtoMessage() {}

func (x *YamlResponse) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use YamlResponse.ProtoReflect.Descriptor instead.
func (*YamlResponse) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{10}
}

func (x *YamlResponse) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *YamlResponse) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

var File_agent_proto protoreflect.FileDescriptor

var file_agent_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x22, 0x8d, 0x01, 0x0a, 0x14, 0x50, 0x61, 0x74, 0x63, 0x68, 0x57, 0x6f,
	0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a,
	0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x77,
	0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x12, 0x14,
	0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69,
	0x6d, 0x61, 0x67, 0x65, 0x22, 0x5a, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x6b,
	0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6e,
	0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x6c, 0x61, 0x62,
	0x65, 0x6c, 0x5f, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x22, 0x57, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x77, 0x6f, 0x72,
	0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x4d, 0x0a, 0x12, 0x47, 0x65, 0x74,
	0x50, 0x6f, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x19, 0x0a,
	0x08, 0x70, 0x6f, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x70, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x49, 0x0a, 0x0b, 0x59, 0x61, 0x6d, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x79, 0x6d, 0x6c, 0x73, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x79, 0x6d, 0x6c, 0x73, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x22, 0x5b, 0x0a, 0x16, 0x52, 0x6f, 0x6c, 0x6c, 0x6f, 0x75, 0x74, 0x57, 0x6f,
	0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a,
	0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x77,
	0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x4e, 0x61, 0x6d, 0x65,
	0x22, 0x78, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x59, 0x61, 0x6d, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65,
	0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x6d, 0x0a, 0x0c, 0x53, 0x63,
	0x61, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61,
	0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e,
	0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x77, 0x6f, 0x72, 0x6b,
	0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x08, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x22, 0x36, 0x0a, 0x0d, 0x4c, 0x61, 0x62,
	0x65, 0x6c, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x25, 0x0a, 0x0e, 0x6c, 0x61,
	0x62, 0x65, 0x6c, 0x5f, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x22, 0x4c, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x36, 0x0a, 0x0c, 0x59, 0x61, 0x6d, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0x9d, 0x07, 0x0a, 0x0b, 0x4c, 0x69, 0x7a, 0x61,
	0x72, 0x64, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x12, 0x3f, 0x0a, 0x0f, 0x70, 0x61, 0x74, 0x63, 0x68,
	0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1b, 0x2e, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x2e, 0x50, 0x61, 0x74, 0x63, 0x68, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x10, 0x70, 0x61, 0x74, 0x63,
	0x68, 0x53, 0x74, 0x61, 0x74, 0x65, 0x66, 0x75, 0x6c, 0x73, 0x65, 0x74, 0x12, 0x1b, 0x2e, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x2e, 0x50, 0x61, 0x74, 0x63, 0x68, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f,
	0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x0e, 0x6c, 0x69,
	0x73, 0x74, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1a, 0x2e, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x0f, 0x6c, 0x69, 0x73,
	0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x66, 0x75, 0x6c, 0x73, 0x65, 0x74, 0x12, 0x1a, 0x2e, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x10, 0x67, 0x65, 0x74,
	0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x50, 0x6f, 0x64, 0x12, 0x19, 0x2e,
	0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x11, 0x67, 0x65, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x66, 0x75, 0x6c, 0x73, 0x65, 0x74, 0x50, 0x6f, 0x64, 0x12, 0x19,
	0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f,
	0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x0b, 0x67, 0x65,
	0x74, 0x50, 0x6f, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x19, 0x2e, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x59,
	0x61, 0x6d, 0x6c, 0x12, 0x12, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x59, 0x61, 0x6d, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x09, 0x61, 0x70, 0x70, 0x6c,
	0x79, 0x59, 0x61, 0x6d, 0x6c, 0x12, 0x12, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x59, 0x61,
	0x6d, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x07, 0x67, 0x65,
	0x74, 0x79, 0x61, 0x6d, 0x6c, 0x12, 0x15, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x47, 0x65,
	0x74, 0x59, 0x61, 0x6d, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x2e, 0x59, 0x61, 0x6d, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x43, 0x0a, 0x11, 0x72, 0x6f, 0x6c, 0x6c, 0x6f, 0x75, 0x74, 0x44, 0x65, 0x70, 0x6c,
	0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1d, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x52,
	0x6f, 0x6c, 0x6c, 0x6f, 0x75, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x12, 0x72, 0x6f, 0x6c, 0x6c, 0x6f, 0x75,
	0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x66, 0x75, 0x6c, 0x73, 0x65, 0x74, 0x12, 0x1d, 0x2e, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x2e, 0x52, 0x6f, 0x6c, 0x6c, 0x6f, 0x75, 0x74, 0x57, 0x6f, 0x72, 0x6b,
	0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x0f,
	0x73, 0x63, 0x61, 0x6c, 0x65, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12,
	0x13, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x53, 0x63, 0x61, 0x6c, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x10, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x66, 0x75, 0x6c, 0x73, 0x65, 0x74, 0x12, 0x13, 0x2e, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2e, 0x53, 0x63, 0x61, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f,
	0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x36, 0x0a, 0x0d, 0x67, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73,
	0x12, 0x14, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x53, 0x65,
	0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x1a, 0x0f, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_agent_proto_rawDescOnce sync.Once
	file_agent_proto_rawDescData = file_agent_proto_rawDesc
)

func file_agent_proto_rawDescGZIP() []byte {
	file_agent_proto_rawDescOnce.Do(func() {
		file_agent_proto_rawDescData = protoimpl.X.CompressGZIP(file_agent_proto_rawDescData)
	})
	return file_agent_proto_rawDescData
}

var file_agent_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_agent_proto_goTypes = []interface{}{
	(*PatchWorkloadRequest)(nil),   // 0: agent.PatchWorkloadRequest
	(*ListWorkloadRequest)(nil),    // 1: agent.ListWorkloadRequest
	(*GetWorkloadRequest)(nil),     // 2: agent.GetWorkloadRequest
	(*GetPodEventRequest)(nil),     // 3: agent.GetPodEventRequest
	(*YamlRequest)(nil),            // 4: agent.YamlRequest
	(*RolloutWorkloadRequest)(nil), // 5: agent.RolloutWorkloadRequest
	(*GetYamlRequest)(nil),         // 6: agent.GetYamlRequest
	(*ScaleRequest)(nil),           // 7: agent.ScaleRequest
	(*LabelSelector)(nil),          // 8: agent.LabelSelector
	(*Response)(nil),               // 9: agent.Response
	(*YamlResponse)(nil),           // 10: agent.YamlResponse
}
var file_agent_proto_depIdxs = []int32{
	0,  // 0: agent.LizardAgent.patchDeployment:input_type -> agent.PatchWorkloadRequest
	0,  // 1: agent.LizardAgent.patchStatefulset:input_type -> agent.PatchWorkloadRequest
	1,  // 2: agent.LizardAgent.listDeployment:input_type -> agent.ListWorkloadRequest
	1,  // 3: agent.LizardAgent.listStatefulset:input_type -> agent.ListWorkloadRequest
	2,  // 4: agent.LizardAgent.getDeploymentPod:input_type -> agent.GetWorkloadRequest
	2,  // 5: agent.LizardAgent.getStatefulsetPod:input_type -> agent.GetWorkloadRequest
	3,  // 6: agent.LizardAgent.getPodEvent:input_type -> agent.GetPodEventRequest
	4,  // 7: agent.LizardAgent.deleteYaml:input_type -> agent.YamlRequest
	4,  // 8: agent.LizardAgent.applyYaml:input_type -> agent.YamlRequest
	6,  // 9: agent.LizardAgent.getyaml:input_type -> agent.GetYamlRequest
	5,  // 10: agent.LizardAgent.rolloutDeployment:input_type -> agent.RolloutWorkloadRequest
	5,  // 11: agent.LizardAgent.rolloutStatefulset:input_type -> agent.RolloutWorkloadRequest
	7,  // 12: agent.LizardAgent.scaleDeployment:input_type -> agent.ScaleRequest
	7,  // 13: agent.LizardAgent.scaleStatefulset:input_type -> agent.ScaleRequest
	8,  // 14: agent.LizardAgent.getNamespaces:input_type -> agent.LabelSelector
	9,  // 15: agent.LizardAgent.patchDeployment:output_type -> agent.Response
	9,  // 16: agent.LizardAgent.patchStatefulset:output_type -> agent.Response
	9,  // 17: agent.LizardAgent.listDeployment:output_type -> agent.Response
	9,  // 18: agent.LizardAgent.listStatefulset:output_type -> agent.Response
	9,  // 19: agent.LizardAgent.getDeploymentPod:output_type -> agent.Response
	9,  // 20: agent.LizardAgent.getStatefulsetPod:output_type -> agent.Response
	9,  // 21: agent.LizardAgent.getPodEvent:output_type -> agent.Response
	9,  // 22: agent.LizardAgent.deleteYaml:output_type -> agent.Response
	9,  // 23: agent.LizardAgent.applyYaml:output_type -> agent.Response
	10, // 24: agent.LizardAgent.getyaml:output_type -> agent.YamlResponse
	9,  // 25: agent.LizardAgent.rolloutDeployment:output_type -> agent.Response
	9,  // 26: agent.LizardAgent.rolloutStatefulset:output_type -> agent.Response
	9,  // 27: agent.LizardAgent.scaleDeployment:output_type -> agent.Response
	9,  // 28: agent.LizardAgent.scaleStatefulset:output_type -> agent.Response
	9,  // 29: agent.LizardAgent.getNamespaces:output_type -> agent.Response
	15, // [15:30] is the sub-list for method output_type
	0,  // [0:15] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_agent_proto_init() }
func file_agent_proto_init() {
	if File_agent_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_agent_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PatchWorkloadRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_agent_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListWorkloadRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_agent_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetWorkloadRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_agent_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPodEventRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_agent_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*YamlRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_agent_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RolloutWorkloadRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_agent_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetYamlRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_agent_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScaleRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_agent_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LabelSelector); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_agent_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_agent_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*YamlResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_agent_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_agent_proto_goTypes,
		DependencyIndexes: file_agent_proto_depIdxs,
		MessageInfos:      file_agent_proto_msgTypes,
	}.Build()
	File_agent_proto = out.File
	file_agent_proto_rawDesc = nil
	file_agent_proto_goTypes = nil
	file_agent_proto_depIdxs = nil
}