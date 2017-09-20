// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pipservice.proto

/*
Package pipservice is a generated protocol buffer package.

It is generated from these files:
	pipservice.proto

It has these top-level messages:
	Request
	Response
*/
package pipservice

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

type Response_Status int32

const (
	Response_OK               Response_Status = 0
	Response_NOTFOUND         Response_Status = 1
	Response_INSUFFICIENTDATA Response_Status = 2
	Response_NOTAPPLICABLE    Response_Status = 3
	Response_GENERALERROR     Response_Status = 4
	Response_SERVICEERROR     Response_Status = 5
	Response_FATALERROR       Response_Status = 6
)

var Response_Status_name = map[int32]string{
	0: "OK",
	1: "NOTFOUND",
	2: "INSUFFICIENTDATA",
	3: "NOTAPPLICABLE",
	4: "GENERALERROR",
	5: "SERVICEERROR",
	6: "FATALERROR",
}
var Response_Status_value = map[string]int32{
	"OK":               0,
	"NOTFOUND":         1,
	"INSUFFICIENTDATA": 2,
	"NOTAPPLICABLE":    3,
	"GENERALERROR":     4,
	"SERVICEERROR":     5,
	"FATALERROR":       6,
}

func (x Response_Status) String() string {
	return proto.EnumName(Response_Status_name, int32(x))
}
func (Response_Status) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

type Request struct {
	QueryURL string `protobuf:"bytes,1,opt,name=queryURL" json:"queryURL,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Request) GetQueryURL() string {
	if m != nil {
		return m.QueryURL
	}
	return ""
}

// use a a string for easy testing
type Response struct {
	Status     Response_Status `protobuf:"varint,1,opt,name=status,enum=pipservice.Response_Status" json:"status,omitempty"`
	Categories string          `protobuf:"bytes,2,opt,name=categories" json:"categories,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Response) GetStatus() Response_Status {
	if m != nil {
		return m.Status
	}
	return Response_OK
}

func (m *Response) GetCategories() string {
	if m != nil {
		return m.Categories
	}
	return ""
}

func init() {
	proto.RegisterType((*Request)(nil), "pipservice.Request")
	proto.RegisterType((*Response)(nil), "pipservice.Response")
	proto.RegisterEnum("pipservice.Response_Status", Response_Status_name, Response_Status_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for PIP service

type PIPClient interface {
	// Sends a query
	GetCategories(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type pIPClient struct {
	cc *grpc.ClientConn
}

func NewPIPClient(cc *grpc.ClientConn) PIPClient {
	return &pIPClient{cc}
}

func (c *pIPClient) GetCategories(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := grpc.Invoke(ctx, "/pipservice.PIP/GetCategories", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for PIP service

type PIPServer interface {
	// Sends a query
	GetCategories(context.Context, *Request) (*Response, error)
}

func RegisterPIPServer(s *grpc.Server, srv PIPServer) {
	s.RegisterService(&_PIP_serviceDesc, srv)
}

func _PIP_GetCategories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PIPServer).GetCategories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pipservice.PIP/GetCategories",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PIPServer).GetCategories(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _PIP_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pipservice.PIP",
	HandlerType: (*PIPServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCategories",
			Handler:    _PIP_GetCategories_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pipservice.proto",
}

func init() { proto.RegisterFile("pipservice.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 272 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0x9b, 0x56, 0x63, 0x1c, 0xda, 0xb2, 0x8e, 0x3d, 0x94, 0x0a, 0x22, 0x01, 0xc1, 0x53,
	0x0e, 0xed, 0xd5, 0xcb, 0x9a, 0x6e, 0xca, 0x62, 0xd8, 0x84, 0x4d, 0xe2, 0xbd, 0x96, 0x41, 0x72,
	0x31, 0x69, 0x76, 0x23, 0x88, 0xaf, 0xea, 0xc3, 0x88, 0xa9, 0xad, 0x15, 0x3c, 0xee, 0xc7, 0xb7,
	0x33, 0xf3, 0xff, 0xc0, 0xea, 0xb2, 0x36, 0xd4, 0xbc, 0x95, 0x1b, 0x0a, 0xea, 0xa6, 0xb2, 0x15,
	0xc2, 0x2f, 0xf1, 0x6f, 0xe1, 0x4c, 0xd3, 0xb6, 0x25, 0x63, 0x71, 0x06, 0xde, 0xb6, 0xa5, 0xe6,
	0xbd, 0xd0, 0xf1, 0xd4, 0xb9, 0x71, 0xee, 0xce, 0xf5, 0xe1, 0xed, 0x7f, 0x3a, 0xe0, 0x69, 0x32,
	0x75, 0xf5, 0x6a, 0x08, 0x17, 0xe0, 0x1a, 0xbb, 0xb6, 0xad, 0xe9, 0xb4, 0xf1, 0xfc, 0x2a, 0x38,
	0x5a, 0xb1, 0xb7, 0x82, 0xac, 0x53, 0xf4, 0x8f, 0x8a, 0xd7, 0x00, 0x9b, 0xb5, 0xa5, 0x97, 0xaa,
	0x29, 0xc9, 0x4c, 0xfb, 0xdd, 0xfc, 0x23, 0xe2, 0x7f, 0x80, 0xbb, 0xfb, 0x81, 0x2e, 0xf4, 0x93,
	0x47, 0xd6, 0xc3, 0x21, 0x78, 0x2a, 0xc9, 0xa3, 0xa4, 0x50, 0x4b, 0xe6, 0xe0, 0x04, 0x98, 0x54,
	0x59, 0x11, 0x45, 0x32, 0x94, 0x42, 0xe5, 0x4b, 0x9e, 0x73, 0xd6, 0xc7, 0x0b, 0x18, 0xa9, 0x24,
	0xe7, 0x69, 0x1a, 0xcb, 0x90, 0x3f, 0xc4, 0x82, 0x0d, 0x90, 0xc1, 0x70, 0x25, 0x94, 0xd0, 0x3c,
	0x16, 0x5a, 0x27, 0x9a, 0x9d, 0x7c, 0x93, 0x4c, 0xe8, 0x27, 0x19, 0x8a, 0x1d, 0x39, 0xc5, 0x31,
	0x40, 0xc4, 0xf3, 0xbd, 0xe1, 0xce, 0x43, 0x18, 0xa4, 0x32, 0xc5, 0x7b, 0x18, 0xad, 0xc8, 0x86,
	0x87, 0xa3, 0xf0, 0xf2, 0x6f, 0xb2, 0xae, 0xa7, 0xd9, 0xe4, 0xbf, 0xb8, 0x7e, 0xef, 0xd9, 0xed,
	0xda, 0x5d, 0x7c, 0x05, 0x00, 0x00, 0xff, 0xff, 0x1b, 0xc3, 0x94, 0xe9, 0x71, 0x01, 0x00, 0x00,
}
