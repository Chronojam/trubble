// Code generated by protoc-gen-go.
// source: auth.proto
// DO NOT EDIT!

package api

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

type NewJwtRequest struct {
	Email    string `protobuf:"bytes,1,opt,name=email" json:"email,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *NewJwtRequest) Reset()                    { *m = NewJwtRequest{} }
func (m *NewJwtRequest) String() string            { return proto.CompactTextString(m) }
func (*NewJwtRequest) ProtoMessage()               {}
func (*NewJwtRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *NewJwtRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *NewJwtRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type NewJwtResponse struct {
	Token string `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
}

func (m *NewJwtResponse) Reset()                    { *m = NewJwtResponse{} }
func (m *NewJwtResponse) String() string            { return proto.CompactTextString(m) }
func (*NewJwtResponse) ProtoMessage()               {}
func (*NewJwtResponse) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *NewJwtResponse) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type UpdatePasswordRequest struct {
	Old string `protobuf:"bytes,1,opt,name=old" json:"old,omitempty"`
	New string `protobuf:"bytes,2,opt,name=new" json:"new,omitempty"`
}

func (m *UpdatePasswordRequest) Reset()                    { *m = UpdatePasswordRequest{} }
func (m *UpdatePasswordRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdatePasswordRequest) ProtoMessage()               {}
func (*UpdatePasswordRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *UpdatePasswordRequest) GetOld() string {
	if m != nil {
		return m.Old
	}
	return ""
}

func (m *UpdatePasswordRequest) GetNew() string {
	if m != nil {
		return m.New
	}
	return ""
}

type UpdatePasswordResponse struct {
}

func (m *UpdatePasswordResponse) Reset()                    { *m = UpdatePasswordResponse{} }
func (m *UpdatePasswordResponse) String() string            { return proto.CompactTextString(m) }
func (*UpdatePasswordResponse) ProtoMessage()               {}
func (*UpdatePasswordResponse) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func init() {
	proto.RegisterType((*NewJwtRequest)(nil), "api.NewJwtRequest")
	proto.RegisterType((*NewJwtResponse)(nil), "api.NewJwtResponse")
	proto.RegisterType((*UpdatePasswordRequest)(nil), "api.UpdatePasswordRequest")
	proto.RegisterType((*UpdatePasswordResponse)(nil), "api.UpdatePasswordResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Auth service

type AuthClient interface {
	IssueNewJWT(ctx context.Context, in *NewJwtRequest, opts ...grpc.CallOption) (*NewJwtResponse, error)
	UpdatePassword(ctx context.Context, in *UpdatePasswordRequest, opts ...grpc.CallOption) (*UpdatePasswordResponse, error)
}

type authClient struct {
	cc *grpc.ClientConn
}

func NewAuthClient(cc *grpc.ClientConn) AuthClient {
	return &authClient{cc}
}

func (c *authClient) IssueNewJWT(ctx context.Context, in *NewJwtRequest, opts ...grpc.CallOption) (*NewJwtResponse, error) {
	out := new(NewJwtResponse)
	err := grpc.Invoke(ctx, "/api.Auth/IssueNewJWT", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) UpdatePassword(ctx context.Context, in *UpdatePasswordRequest, opts ...grpc.CallOption) (*UpdatePasswordResponse, error) {
	out := new(UpdatePasswordResponse)
	err := grpc.Invoke(ctx, "/api.Auth/UpdatePassword", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Auth service

type AuthServer interface {
	IssueNewJWT(context.Context, *NewJwtRequest) (*NewJwtResponse, error)
	UpdatePassword(context.Context, *UpdatePasswordRequest) (*UpdatePasswordResponse, error)
}

func RegisterAuthServer(s *grpc.Server, srv AuthServer) {
	s.RegisterService(&_Auth_serviceDesc, srv)
}

func _Auth_IssueNewJWT_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewJwtRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).IssueNewJWT(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Auth/IssueNewJWT",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).IssueNewJWT(ctx, req.(*NewJwtRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_UpdatePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).UpdatePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Auth/UpdatePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).UpdatePassword(ctx, req.(*UpdatePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Auth_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IssueNewJWT",
			Handler:    _Auth_IssueNewJWT_Handler,
		},
		{
			MethodName: "UpdatePassword",
			Handler:    _Auth_UpdatePassword_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}

func init() { proto.RegisterFile("auth.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 223 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x2c, 0x2d, 0xc9,
	0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4e, 0x2c, 0xc8, 0x54, 0x72, 0xe4, 0xe2, 0xf5,
	0x4b, 0x2d, 0xf7, 0x2a, 0x2f, 0x09, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0x12, 0xe1, 0x62,
	0x4d, 0xcd, 0x4d, 0xcc, 0xcc, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x82, 0x70, 0x84, 0xa4,
	0xb8, 0x38, 0x0a, 0x12, 0x8b, 0x8b, 0xcb, 0xf3, 0x8b, 0x52, 0x24, 0x98, 0xc0, 0x12, 0x70, 0xbe,
	0x92, 0x1a, 0x17, 0x1f, 0xcc, 0x88, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x90, 0x19, 0x25, 0xf9,
	0xd9, 0xa9, 0x79, 0x30, 0x33, 0xc0, 0x1c, 0x25, 0x6b, 0x2e, 0xd1, 0xd0, 0x82, 0x94, 0xc4, 0x92,
	0xd4, 0x00, 0xa8, 0x4e, 0x98, 0x95, 0x02, 0x5c, 0xcc, 0xf9, 0x39, 0x29, 0x50, 0xc5, 0x20, 0x26,
	0x48, 0x24, 0x2f, 0xb5, 0x1c, 0x6a, 0x13, 0x88, 0xa9, 0x24, 0xc1, 0x25, 0x86, 0xae, 0x19, 0x62,
	0x99, 0x51, 0x2f, 0x23, 0x17, 0x8b, 0x63, 0x69, 0x49, 0x86, 0x90, 0x05, 0x17, 0xb7, 0x67, 0x71,
	0x71, 0x69, 0x2a, 0xc8, 0x31, 0xe1, 0x21, 0x42, 0x42, 0x7a, 0x89, 0x05, 0x99, 0x7a, 0x28, 0x9e,
	0x93, 0x12, 0x46, 0x11, 0x83, 0x18, 0xa0, 0xc4, 0x20, 0xe4, 0xcd, 0xc5, 0x87, 0x6a, 0xb8, 0x90,
	0x14, 0x58, 0x21, 0x56, 0xe7, 0x4a, 0x49, 0x63, 0x95, 0x83, 0x19, 0x96, 0xc4, 0x06, 0x0e, 0x5d,
	0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf9, 0x56, 0x09, 0x1f, 0x6b, 0x01, 0x00, 0x00,
}
