// Code generated by protoc-gen-go. DO NOT EDIT.
// source: event.proto

package eventpushdown

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Data struct {
	Ret                  int64             `protobuf:"varint,1,opt,name=ret,proto3" json:"ret,omitempty"`
	Data                 map[string]string `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Data) Reset()         { *m = Data{} }
func (m *Data) String() string { return proto.CompactTextString(m) }
func (*Data) ProtoMessage()    {}
func (*Data) Descriptor() ([]byte, []int) {
	return fileDescriptor_event_047130b9cccf1bd7, []int{0}
}
func (m *Data) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Data.Unmarshal(m, b)
}
func (m *Data) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Data.Marshal(b, m, deterministic)
}
func (dst *Data) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Data.Merge(dst, src)
}
func (m *Data) XXX_Size() int {
	return xxx_messageInfo_Data.Size(m)
}
func (m *Data) XXX_DiscardUnknown() {
	xxx_messageInfo_Data.DiscardUnknown(m)
}

var xxx_messageInfo_Data proto.InternalMessageInfo

func (m *Data) GetRet() int64 {
	if m != nil {
		return m.Ret
	}
	return 0
}

func (m *Data) GetData() map[string]string {
	if m != nil {
		return m.Data
	}
	return nil
}

type PushDownRequest struct {
	Mid                  string   `protobuf:"bytes,1,opt,name=mid,proto3" json:"mid,omitempty"`
	Version              string   `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	PKey                 string   `protobuf:"bytes,3,opt,name=pKey,proto3" json:"pKey,omitempty"`
	DName                string   `protobuf:"bytes,4,opt,name=dName,proto3" json:"dName,omitempty"`
	Obj                  string   `protobuf:"bytes,5,opt,name=obj,proto3" json:"obj,omitempty"`
	Action               string   `protobuf:"bytes,6,opt,name=action,proto3" json:"action,omitempty"`
	Data                 *Data    `protobuf:"bytes,7,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PushDownRequest) Reset()         { *m = PushDownRequest{} }
func (m *PushDownRequest) String() string { return proto.CompactTextString(m) }
func (*PushDownRequest) ProtoMessage()    {}
func (*PushDownRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_event_047130b9cccf1bd7, []int{1}
}
func (m *PushDownRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushDownRequest.Unmarshal(m, b)
}
func (m *PushDownRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushDownRequest.Marshal(b, m, deterministic)
}
func (dst *PushDownRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushDownRequest.Merge(dst, src)
}
func (m *PushDownRequest) XXX_Size() int {
	return xxx_messageInfo_PushDownRequest.Size(m)
}
func (m *PushDownRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PushDownRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PushDownRequest proto.InternalMessageInfo

func (m *PushDownRequest) GetMid() string {
	if m != nil {
		return m.Mid
	}
	return ""
}

func (m *PushDownRequest) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *PushDownRequest) GetPKey() string {
	if m != nil {
		return m.PKey
	}
	return ""
}

func (m *PushDownRequest) GetDName() string {
	if m != nil {
		return m.DName
	}
	return ""
}

func (m *PushDownRequest) GetObj() string {
	if m != nil {
		return m.Obj
	}
	return ""
}

func (m *PushDownRequest) GetAction() string {
	if m != nil {
		return m.Action
	}
	return ""
}

func (m *PushDownRequest) GetData() *Data {
	if m != nil {
		return m.Data
	}
	return nil
}

type PushDownReply struct {
	Status               int64    `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Err                  string   `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PushDownReply) Reset()         { *m = PushDownReply{} }
func (m *PushDownReply) String() string { return proto.CompactTextString(m) }
func (*PushDownReply) ProtoMessage()    {}
func (*PushDownReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_event_047130b9cccf1bd7, []int{2}
}
func (m *PushDownReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushDownReply.Unmarshal(m, b)
}
func (m *PushDownReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushDownReply.Marshal(b, m, deterministic)
}
func (dst *PushDownReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushDownReply.Merge(dst, src)
}
func (m *PushDownReply) XXX_Size() int {
	return xxx_messageInfo_PushDownReply.Size(m)
}
func (m *PushDownReply) XXX_DiscardUnknown() {
	xxx_messageInfo_PushDownReply.DiscardUnknown(m)
}

var xxx_messageInfo_PushDownReply proto.InternalMessageInfo

func (m *PushDownReply) GetStatus() int64 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *PushDownReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

func init() {
	proto.RegisterType((*Data)(nil), "eventpushdown.Data")
	proto.RegisterMapType((map[string]string)(nil), "eventpushdown.Data.DataEntry")
	proto.RegisterType((*PushDownRequest)(nil), "eventpushdown.PushDownRequest")
	proto.RegisterType((*PushDownReply)(nil), "eventpushdown.PushDownReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EventPushDownClient is the client API for EventPushDown service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EventPushDownClient interface {
	PushDown(ctx context.Context, in *PushDownRequest, opts ...grpc.CallOption) (*PushDownReply, error)
}

type eventPushDownClient struct {
	cc *grpc.ClientConn
}

func NewEventPushDownClient(cc *grpc.ClientConn) EventPushDownClient {
	return &eventPushDownClient{cc}
}

func (c *eventPushDownClient) PushDown(ctx context.Context, in *PushDownRequest, opts ...grpc.CallOption) (*PushDownReply, error) {
	out := new(PushDownReply)
	err := c.cc.Invoke(ctx, "/eventpushdown.EventPushDown/PushDown", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventPushDownServer is the server API for EventPushDown service.
type EventPushDownServer interface {
	PushDown(context.Context, *PushDownRequest) (*PushDownReply, error)
}

func RegisterEventPushDownServer(s *grpc.Server, srv EventPushDownServer) {
	s.RegisterService(&_EventPushDown_serviceDesc, srv)
}

func _EventPushDown_PushDown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushDownRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventPushDownServer).PushDown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eventpushdown.EventPushDown/PushDown",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventPushDownServer).PushDown(ctx, req.(*PushDownRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _EventPushDown_serviceDesc = grpc.ServiceDesc{
	ServiceName: "eventpushdown.EventPushDown",
	HandlerType: (*EventPushDownServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PushDown",
			Handler:    _EventPushDown_PushDown_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "event.proto",
}

func init() { proto.RegisterFile("event.proto", fileDescriptor_event_047130b9cccf1bd7) }

var fileDescriptor_event_047130b9cccf1bd7 = []byte{
	// 309 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0x4f, 0x4b, 0xc3, 0x30,
	0x18, 0xc6, 0xcd, 0xda, 0x6d, 0xee, 0x1d, 0x43, 0x89, 0x22, 0x61, 0xa8, 0x8c, 0x5e, 0xdc, 0x69,
	0x60, 0x3d, 0xf8, 0xe7, 0xbc, 0x5d, 0x14, 0x44, 0x7a, 0xf5, 0x94, 0xda, 0xc0, 0x3a, 0xdb, 0xa6,
	0x26, 0x69, 0x47, 0xef, 0x7e, 0x2a, 0x3f, 0x9d, 0xbc, 0x69, 0x3a, 0x71, 0xe8, 0xa5, 0x3c, 0xcf,
	0xdb, 0x3e, 0xbf, 0x27, 0x6f, 0x0a, 0x63, 0x51, 0x8b, 0xc2, 0x2c, 0x4a, 0x25, 0x8d, 0xa4, 0x13,
	0x6b, 0xca, 0x4a, 0xaf, 0x13, 0xb9, 0x2d, 0x82, 0x4f, 0x02, 0xfe, 0x92, 0x1b, 0x4e, 0x8f, 0xc1,
	0x53, 0xc2, 0x30, 0x32, 0x23, 0x73, 0x2f, 0x42, 0x49, 0xaf, 0xc1, 0x4f, 0xb8, 0xe1, 0xac, 0x37,
	0xf3, 0xe6, 0xe3, 0xf0, 0x62, 0xf1, 0x2b, 0xb8, 0xc0, 0x90, 0x7d, 0xac, 0x0a, 0xa3, 0x9a, 0xc8,
	0x7e, 0x3a, 0xbd, 0x85, 0xd1, 0x6e, 0x84, 0xc4, 0x77, 0xd1, 0x58, 0xe2, 0x28, 0x42, 0x49, 0x4f,
	0xa1, 0x5f, 0xf3, 0xac, 0x12, 0xac, 0x67, 0x67, 0xad, 0x79, 0xe8, 0xdd, 0x91, 0xe0, 0x8b, 0xc0,
	0xd1, 0x4b, 0xa5, 0xd7, 0x4b, 0xb9, 0x2d, 0x22, 0xf1, 0x51, 0x09, 0x6d, 0x30, 0x9f, 0xa7, 0x49,
	0x97, 0xcf, 0xd3, 0x84, 0x32, 0x18, 0xd6, 0x42, 0xe9, 0x54, 0x16, 0x8e, 0xd0, 0x59, 0x4a, 0xc1,
	0x2f, 0x9f, 0x44, 0xc3, 0x3c, 0x3b, 0xb6, 0x1a, 0xdb, 0x92, 0x67, 0x9e, 0x0b, 0xe6, 0xb7, 0x6d,
	0xd6, 0x20, 0x55, 0xc6, 0x1b, 0xd6, 0x6f, 0xa9, 0x32, 0xde, 0xd0, 0x33, 0x18, 0xf0, 0x37, 0x83,
	0xd0, 0x81, 0x1d, 0x3a, 0x47, 0xaf, 0xdc, 0xfe, 0xc3, 0x19, 0x99, 0x8f, 0xc3, 0x93, 0x3f, 0xf6,
	0x6f, 0xb7, 0x0e, 0xee, 0x61, 0xf2, 0x73, 0xf6, 0x32, 0x6b, 0x90, 0xa8, 0x0d, 0x37, 0x95, 0x76,
	0xd7, 0xe9, 0x1c, 0x76, 0x0b, 0xa5, 0xdc, 0xd9, 0x51, 0x86, 0xaf, 0x30, 0x59, 0x21, 0xb6, 0xcb,
	0xd3, 0x47, 0x38, 0xdc, 0xe9, 0xcb, 0xbd, 0xca, 0xbd, 0x0b, 0x9a, 0x9e, 0xff, 0xfb, 0xbe, 0xcc,
	0x9a, 0xe0, 0x20, 0x1e, 0xd8, 0x3f, 0x7e, 0xf3, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x77, 0xcc, 0x7e,
	0xd2, 0x00, 0x02, 0x00, 0x00,
}