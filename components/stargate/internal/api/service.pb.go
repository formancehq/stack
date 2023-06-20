// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: internal/api/service.proto

package api

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

type Values struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values []string `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *Values) Reset() {
	*x = Values{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Values) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Values) ProtoMessage() {}

func (x *Values) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Values.ProtoReflect.Descriptor instead.
func (*Values) Descriptor() ([]byte, []int) {
	return file_internal_api_service_proto_rawDescGZIP(), []int{0}
}

func (x *Values) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

type StargateServerMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CorrelationId string `protobuf:"bytes,1,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	// Types that are assignable to Event:
	//
	//	*StargateServerMessage_ApiCall
	//	*StargateServerMessage_Ping_
	Event isStargateServerMessage_Event `protobuf_oneof:"event"`
}

func (x *StargateServerMessage) Reset() {
	*x = StargateServerMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StargateServerMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StargateServerMessage) ProtoMessage() {}

func (x *StargateServerMessage) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StargateServerMessage.ProtoReflect.Descriptor instead.
func (*StargateServerMessage) Descriptor() ([]byte, []int) {
	return file_internal_api_service_proto_rawDescGZIP(), []int{1}
}

func (x *StargateServerMessage) GetCorrelationId() string {
	if x != nil {
		return x.CorrelationId
	}
	return ""
}

func (m *StargateServerMessage) GetEvent() isStargateServerMessage_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (x *StargateServerMessage) GetApiCall() *StargateServerMessage_APICall {
	if x, ok := x.GetEvent().(*StargateServerMessage_ApiCall); ok {
		return x.ApiCall
	}
	return nil
}

func (x *StargateServerMessage) GetPing() *StargateServerMessage_Ping {
	if x, ok := x.GetEvent().(*StargateServerMessage_Ping_); ok {
		return x.Ping
	}
	return nil
}

type isStargateServerMessage_Event interface {
	isStargateServerMessage_Event()
}

type StargateServerMessage_ApiCall struct {
	ApiCall *StargateServerMessage_APICall `protobuf:"bytes,101,opt,name=api_call,json=apiCall,proto3,oneof"`
}

type StargateServerMessage_Ping_ struct {
	Ping *StargateServerMessage_Ping `protobuf:"bytes,102,opt,name=ping,proto3,oneof"`
}

func (*StargateServerMessage_ApiCall) isStargateServerMessage_Event() {}

func (*StargateServerMessage_Ping_) isStargateServerMessage_Event() {}

type StargateClientMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CorrelationId string `protobuf:"bytes,1,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	// Types that are assignable to Event:
	//
	//	*StargateClientMessage_ApiCallResponse
	//	*StargateClientMessage_Pong_
	Event isStargateClientMessage_Event `protobuf_oneof:"event"`
}

func (x *StargateClientMessage) Reset() {
	*x = StargateClientMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StargateClientMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StargateClientMessage) ProtoMessage() {}

func (x *StargateClientMessage) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StargateClientMessage.ProtoReflect.Descriptor instead.
func (*StargateClientMessage) Descriptor() ([]byte, []int) {
	return file_internal_api_service_proto_rawDescGZIP(), []int{2}
}

func (x *StargateClientMessage) GetCorrelationId() string {
	if x != nil {
		return x.CorrelationId
	}
	return ""
}

func (m *StargateClientMessage) GetEvent() isStargateClientMessage_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (x *StargateClientMessage) GetApiCallResponse() *StargateClientMessage_APICallResponse {
	if x, ok := x.GetEvent().(*StargateClientMessage_ApiCallResponse); ok {
		return x.ApiCallResponse
	}
	return nil
}

func (x *StargateClientMessage) GetPong() *StargateClientMessage_Pong {
	if x, ok := x.GetEvent().(*StargateClientMessage_Pong_); ok {
		return x.Pong
	}
	return nil
}

type isStargateClientMessage_Event interface {
	isStargateClientMessage_Event()
}

type StargateClientMessage_ApiCallResponse struct {
	ApiCallResponse *StargateClientMessage_APICallResponse `protobuf:"bytes,101,opt,name=api_call_response,json=apiCallResponse,proto3,oneof"`
}

type StargateClientMessage_Pong_ struct {
	Pong *StargateClientMessage_Pong `protobuf:"bytes,102,opt,name=pong,proto3,oneof"`
}

func (*StargateClientMessage_ApiCallResponse) isStargateClientMessage_Event() {}

func (*StargateClientMessage_Pong_) isStargateClientMessage_Event() {}

type StargateServerMessage_APICall struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Method  string             `protobuf:"bytes,1,opt,name=method,proto3" json:"method,omitempty"`
	Path    string             `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	Query   map[string]*Values `protobuf:"bytes,3,rep,name=query,proto3" json:"query,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Body    []byte             `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
	Headers map[string]*Values `protobuf:"bytes,5,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *StargateServerMessage_APICall) Reset() {
	*x = StargateServerMessage_APICall{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StargateServerMessage_APICall) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StargateServerMessage_APICall) ProtoMessage() {}

func (x *StargateServerMessage_APICall) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StargateServerMessage_APICall.ProtoReflect.Descriptor instead.
func (*StargateServerMessage_APICall) Descriptor() ([]byte, []int) {
	return file_internal_api_service_proto_rawDescGZIP(), []int{1, 0}
}

func (x *StargateServerMessage_APICall) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *StargateServerMessage_APICall) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *StargateServerMessage_APICall) GetQuery() map[string]*Values {
	if x != nil {
		return x.Query
	}
	return nil
}

func (x *StargateServerMessage_APICall) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *StargateServerMessage_APICall) GetHeaders() map[string]*Values {
	if x != nil {
		return x.Headers
	}
	return nil
}

type StargateServerMessage_Ping struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StargateServerMessage_Ping) Reset() {
	*x = StargateServerMessage_Ping{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StargateServerMessage_Ping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StargateServerMessage_Ping) ProtoMessage() {}

func (x *StargateServerMessage_Ping) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StargateServerMessage_Ping.ProtoReflect.Descriptor instead.
func (*StargateServerMessage_Ping) Descriptor() ([]byte, []int) {
	return file_internal_api_service_proto_rawDescGZIP(), []int{1, 1}
}

type StargateClientMessage_APICallResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode int32              `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	Body       []byte             `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	Headers    map[string]*Values `protobuf:"bytes,3,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *StargateClientMessage_APICallResponse) Reset() {
	*x = StargateClientMessage_APICallResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StargateClientMessage_APICallResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StargateClientMessage_APICallResponse) ProtoMessage() {}

func (x *StargateClientMessage_APICallResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StargateClientMessage_APICallResponse.ProtoReflect.Descriptor instead.
func (*StargateClientMessage_APICallResponse) Descriptor() ([]byte, []int) {
	return file_internal_api_service_proto_rawDescGZIP(), []int{2, 0}
}

func (x *StargateClientMessage_APICallResponse) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *StargateClientMessage_APICallResponse) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *StargateClientMessage_APICallResponse) GetHeaders() map[string]*Values {
	if x != nil {
		return x.Headers
	}
	return nil
}

type StargateClientMessage_Pong struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StargateClientMessage_Pong) Reset() {
	*x = StargateClientMessage_Pong{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StargateClientMessage_Pong) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StargateClientMessage_Pong) ProtoMessage() {}

func (x *StargateClientMessage_Pong) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StargateClientMessage_Pong.ProtoReflect.Descriptor instead.
func (*StargateClientMessage_Pong) Descriptor() ([]byte, []int) {
	return file_internal_api_service_proto_rawDescGZIP(), []int{2, 1}
}

var File_internal_api_service_proto protoreflect.FileDescriptor

var file_internal_api_service_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x66, 0x6f,
	0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x2e,
	0x61, 0x70, 0x69, 0x22, 0x20, 0x0a, 0x06, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x16, 0x0a,
	0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x9f, 0x05, 0x0a, 0x15, 0x53, 0x74, 0x61, 0x72, 0x67, 0x61,
	0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x51, 0x0a, 0x08, 0x61, 0x70, 0x69, 0x5f, 0x63, 0x61,
	0x6c, 0x6c, 0x18, 0x65, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x66, 0x6f, 0x72, 0x6d, 0x61,
	0x6e, 0x63, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x53, 0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x41, 0x50, 0x49, 0x43, 0x61, 0x6c, 0x6c, 0x48, 0x00,
	0x52, 0x07, 0x61, 0x70, 0x69, 0x43, 0x61, 0x6c, 0x6c, 0x12, 0x47, 0x0a, 0x04, 0x70, 0x69, 0x6e,
	0x67, 0x18, 0x66, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e,
	0x63, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x53, 0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x48, 0x00, 0x52, 0x04, 0x70, 0x69,
	0x6e, 0x67, 0x1a, 0xb1, 0x03, 0x0a, 0x07, 0x41, 0x50, 0x49, 0x43, 0x61, 0x6c, 0x6c, 0x12, 0x16,
	0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x55, 0x0a, 0x05, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x3f, 0x2e, 0x66, 0x6f, 0x72, 0x6d,
	0x61, 0x6e, 0x63, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x41, 0x50, 0x49, 0x43, 0x61, 0x6c, 0x6c, 0x2e,
	0x51, 0x75, 0x65, 0x72, 0x79, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72,
	0x79, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x5b, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73,
	0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x41, 0x2e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63,
	0x65, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53,
	0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2e, 0x41, 0x50, 0x49, 0x43, 0x61, 0x6c, 0x6c, 0x2e, 0x48, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x73, 0x1a, 0x57, 0x0a, 0x0a, 0x51, 0x75, 0x65, 0x72, 0x79, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x33, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1d, 0x2e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x73, 0x74, 0x61,
	0x72, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x59, 0x0a, 0x0c, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x33, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x66,
	0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x06, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x42, 0x07,
	0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x8d, 0x04, 0x0a, 0x15, 0x53, 0x74, 0x61, 0x72,
	0x67, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x6a, 0x0a, 0x11, 0x61, 0x70, 0x69, 0x5f,
	0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x65, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x3c, 0x2e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x73,
	0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53, 0x74, 0x61, 0x72,
	0x67, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x41, 0x50, 0x49, 0x43, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x48, 0x00, 0x52, 0x0f, 0x61, 0x70, 0x69, 0x43, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x04, 0x70, 0x6f, 0x6e, 0x67, 0x18, 0x66, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x31, 0x2e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x73, 0x74,
	0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x67,
	0x61, 0x74, 0x65, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x50, 0x6f, 0x6e, 0x67, 0x48, 0x00, 0x52, 0x04, 0x70, 0x6f, 0x6e, 0x67, 0x1a, 0x86, 0x02,
	0x0a, 0x0f, 0x41, 0x50, 0x49, 0x43, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f,
	0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x63, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x49, 0x2e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e,
	0x63, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x53, 0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x41, 0x50, 0x49, 0x43, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x1a, 0x59, 0x0a, 0x0c, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x33, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x66,
	0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x06, 0x0a, 0x04, 0x50, 0x6f, 0x6e, 0x67, 0x42, 0x07,
	0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x32, 0x7d, 0x0a, 0x0f, 0x53, 0x74, 0x61, 0x72, 0x67,
	0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x6a, 0x0a, 0x08, 0x53, 0x74,
	0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x12, 0x2c, 0x2e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63,
	0x65, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53,
	0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x1a, 0x2c, 0x2e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x2e,
	0x73, 0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53, 0x74, 0x61,
	0x72, 0x67, 0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x28, 0x01, 0x30, 0x01, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x68, 0x71, 0x2f,
	0x73, 0x74, 0x61, 0x63, 0x6b, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73,
	0x2f, 0x73, 0x74, 0x61, 0x72, 0x67, 0x61, 0x74, 0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_api_service_proto_rawDescOnce sync.Once
	file_internal_api_service_proto_rawDescData = file_internal_api_service_proto_rawDesc
)

func file_internal_api_service_proto_rawDescGZIP() []byte {
	file_internal_api_service_proto_rawDescOnce.Do(func() {
		file_internal_api_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_api_service_proto_rawDescData)
	})
	return file_internal_api_service_proto_rawDescData
}

var file_internal_api_service_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_internal_api_service_proto_goTypes = []interface{}{
	(*Values)(nil),                        // 0: formance.stargate.api.Values
	(*StargateServerMessage)(nil),         // 1: formance.stargate.api.StargateServerMessage
	(*StargateClientMessage)(nil),         // 2: formance.stargate.api.StargateClientMessage
	(*StargateServerMessage_APICall)(nil), // 3: formance.stargate.api.StargateServerMessage.APICall
	(*StargateServerMessage_Ping)(nil),    // 4: formance.stargate.api.StargateServerMessage.Ping
	nil,                                   // 5: formance.stargate.api.StargateServerMessage.APICall.QueryEntry
	nil,                                   // 6: formance.stargate.api.StargateServerMessage.APICall.HeadersEntry
	(*StargateClientMessage_APICallResponse)(nil), // 7: formance.stargate.api.StargateClientMessage.APICallResponse
	(*StargateClientMessage_Pong)(nil),            // 8: formance.stargate.api.StargateClientMessage.Pong
	nil,                                           // 9: formance.stargate.api.StargateClientMessage.APICallResponse.HeadersEntry
}
var file_internal_api_service_proto_depIdxs = []int32{
	3,  // 0: formance.stargate.api.StargateServerMessage.api_call:type_name -> formance.stargate.api.StargateServerMessage.APICall
	4,  // 1: formance.stargate.api.StargateServerMessage.ping:type_name -> formance.stargate.api.StargateServerMessage.Ping
	7,  // 2: formance.stargate.api.StargateClientMessage.api_call_response:type_name -> formance.stargate.api.StargateClientMessage.APICallResponse
	8,  // 3: formance.stargate.api.StargateClientMessage.pong:type_name -> formance.stargate.api.StargateClientMessage.Pong
	5,  // 4: formance.stargate.api.StargateServerMessage.APICall.query:type_name -> formance.stargate.api.StargateServerMessage.APICall.QueryEntry
	6,  // 5: formance.stargate.api.StargateServerMessage.APICall.headers:type_name -> formance.stargate.api.StargateServerMessage.APICall.HeadersEntry
	0,  // 6: formance.stargate.api.StargateServerMessage.APICall.QueryEntry.value:type_name -> formance.stargate.api.Values
	0,  // 7: formance.stargate.api.StargateServerMessage.APICall.HeadersEntry.value:type_name -> formance.stargate.api.Values
	9,  // 8: formance.stargate.api.StargateClientMessage.APICallResponse.headers:type_name -> formance.stargate.api.StargateClientMessage.APICallResponse.HeadersEntry
	0,  // 9: formance.stargate.api.StargateClientMessage.APICallResponse.HeadersEntry.value:type_name -> formance.stargate.api.Values
	2,  // 10: formance.stargate.api.StargateService.Stargate:input_type -> formance.stargate.api.StargateClientMessage
	1,  // 11: formance.stargate.api.StargateService.Stargate:output_type -> formance.stargate.api.StargateServerMessage
	11, // [11:12] is the sub-list for method output_type
	10, // [10:11] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_internal_api_service_proto_init() }
func file_internal_api_service_proto_init() {
	if File_internal_api_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_api_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Values); i {
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
		file_internal_api_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StargateServerMessage); i {
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
		file_internal_api_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StargateClientMessage); i {
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
		file_internal_api_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StargateServerMessage_APICall); i {
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
		file_internal_api_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StargateServerMessage_Ping); i {
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
		file_internal_api_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StargateClientMessage_APICallResponse); i {
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
		file_internal_api_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StargateClientMessage_Pong); i {
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
	file_internal_api_service_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*StargateServerMessage_ApiCall)(nil),
		(*StargateServerMessage_Ping_)(nil),
	}
	file_internal_api_service_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*StargateClientMessage_ApiCallResponse)(nil),
		(*StargateClientMessage_Pong_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_api_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_api_service_proto_goTypes,
		DependencyIndexes: file_internal_api_service_proto_depIdxs,
		MessageInfos:      file_internal_api_service_proto_msgTypes,
	}.Build()
	File_internal_api_service_proto = out.File
	file_internal_api_service_proto_rawDesc = nil
	file_internal_api_service_proto_goTypes = nil
	file_internal_api_service_proto_depIdxs = nil
}
