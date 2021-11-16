// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package operations

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// OperationsClient is the client API for Operations service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OperationsClient interface {
	Get(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Value, error)
	Put(ctx context.Context, in *KeyValue, opts ...grpc.CallOption) (*Ack, error)
	Append(ctx context.Context, in *KeyValue, opts ...grpc.CallOption) (*Ack, error)
	Del(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Ack, error)
	GetInternal(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Value, error)
	PutInternal(ctx context.Context, in *KeyValue, opts ...grpc.CallOption) (*Ack, error)
	AppendInternal(ctx context.Context, in *KeyValue, opts ...grpc.CallOption) (*Ack, error)
	DelInternal(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Ack, error)
	Migration(ctx context.Context, in *KeyCost, opts ...grpc.CallOption) (*Outcome, error)
	Ping(ctx context.Context, in *PingMessage, opts ...grpc.CallOption) (*Ack, error)
	Join(ctx context.Context, in *JoinMessage, opts ...grpc.CallOption) (*JoinResponse, error)
	RequestJoin(ctx context.Context, in *RequestJoinMessage, opts ...grpc.CallOption) (*JoinMessage, error)
	LeaveCluster(ctx context.Context, in *RequestJoinMessage, opts ...grpc.CallOption) (*Ack, error)
}

type operationsClient struct {
	cc grpc.ClientConnInterface
}

func NewOperationsClient(cc grpc.ClientConnInterface) OperationsClient {
	return &operationsClient{cc}
}

func (c *operationsClient) Get(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Value, error) {
	out := new(Value)
	err := c.cc.Invoke(ctx, "/operations.Operations/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operationsClient) Put(ctx context.Context, in *KeyValue, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/operations.Operations/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operationsClient) Append(ctx context.Context, in *KeyValue, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/operations.Operations/Append", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operationsClient) Del(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/operations.Operations/Del", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operationsClient) GetInternal(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Value, error) {
	out := new(Value)
	err := c.cc.Invoke(ctx, "/operations.Operations/GetInternal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operationsClient) PutInternal(ctx context.Context, in *KeyValue, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/operations.Operations/PutInternal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operationsClient) AppendInternal(ctx context.Context, in *KeyValue, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/operations.Operations/AppendInternal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operationsClient) DelInternal(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/operations.Operations/DelInternal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operationsClient) Migration(ctx context.Context, in *KeyCost, opts ...grpc.CallOption) (*Outcome, error) {
	out := new(Outcome)
	err := c.cc.Invoke(ctx, "/operations.Operations/Migration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operationsClient) Ping(ctx context.Context, in *PingMessage, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/operations.Operations/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operationsClient) Join(ctx context.Context, in *JoinMessage, opts ...grpc.CallOption) (*JoinResponse, error) {
	out := new(JoinResponse)
	err := c.cc.Invoke(ctx, "/operations.Operations/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operationsClient) RequestJoin(ctx context.Context, in *RequestJoinMessage, opts ...grpc.CallOption) (*JoinMessage, error) {
	out := new(JoinMessage)
	err := c.cc.Invoke(ctx, "/operations.Operations/RequestJoin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operationsClient) LeaveCluster(ctx context.Context, in *RequestJoinMessage, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/operations.Operations/LeaveCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OperationsServer is the server API for Operations service.
// All implementations must embed UnimplementedOperationsServer
// for forward compatibility
type OperationsServer interface {
	Get(context.Context, *Key) (*Value, error)
	Put(context.Context, *KeyValue) (*Ack, error)
	Append(context.Context, *KeyValue) (*Ack, error)
	Del(context.Context, *Key) (*Ack, error)
	GetInternal(context.Context, *Key) (*Value, error)
	PutInternal(context.Context, *KeyValue) (*Ack, error)
	AppendInternal(context.Context, *KeyValue) (*Ack, error)
	DelInternal(context.Context, *Key) (*Ack, error)
	Migration(context.Context, *KeyCost) (*Outcome, error)
	Ping(context.Context, *PingMessage) (*Ack, error)
	Join(context.Context, *JoinMessage) (*JoinResponse, error)
	RequestJoin(context.Context, *RequestJoinMessage) (*JoinMessage, error)
	LeaveCluster(context.Context, *RequestJoinMessage) (*Ack, error)
	mustEmbedUnimplementedOperationsServer()
}

// UnimplementedOperationsServer must be embedded to have forward compatible implementations.
type UnimplementedOperationsServer struct {
}

func (UnimplementedOperationsServer) Get(context.Context, *Key) (*Value, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedOperationsServer) Put(context.Context, *KeyValue) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedOperationsServer) Append(context.Context, *KeyValue) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Append not implemented")
}
func (UnimplementedOperationsServer) Del(context.Context, *Key) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Del not implemented")
}
func (UnimplementedOperationsServer) GetInternal(context.Context, *Key) (*Value, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInternal not implemented")
}
func (UnimplementedOperationsServer) PutInternal(context.Context, *KeyValue) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutInternal not implemented")
}
func (UnimplementedOperationsServer) AppendInternal(context.Context, *KeyValue) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendInternal not implemented")
}
func (UnimplementedOperationsServer) DelInternal(context.Context, *Key) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelInternal not implemented")
}
func (UnimplementedOperationsServer) Migration(context.Context, *KeyCost) (*Outcome, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Migration not implemented")
}
func (UnimplementedOperationsServer) Ping(context.Context, *PingMessage) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedOperationsServer) Join(context.Context, *JoinMessage) (*JoinResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedOperationsServer) RequestJoin(context.Context, *RequestJoinMessage) (*JoinMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestJoin not implemented")
}
func (UnimplementedOperationsServer) LeaveCluster(context.Context, *RequestJoinMessage) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveCluster not implemented")
}
func (UnimplementedOperationsServer) mustEmbedUnimplementedOperationsServer() {}

// UnsafeOperationsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OperationsServer will
// result in compilation errors.
type UnsafeOperationsServer interface {
	mustEmbedUnimplementedOperationsServer()
}

func RegisterOperationsServer(s grpc.ServiceRegistrar, srv OperationsServer) {
	s.RegisterService(&Operations_ServiceDesc, srv)
}

func _Operations_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationsServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/operations.Operations/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationsServer).Get(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _Operations_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationsServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/operations.Operations/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationsServer).Put(ctx, req.(*KeyValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _Operations_Append_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationsServer).Append(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/operations.Operations/Append",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationsServer).Append(ctx, req.(*KeyValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _Operations_Del_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationsServer).Del(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/operations.Operations/Del",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationsServer).Del(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _Operations_GetInternal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationsServer).GetInternal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/operations.Operations/GetInternal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationsServer).GetInternal(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _Operations_PutInternal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationsServer).PutInternal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/operations.Operations/PutInternal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationsServer).PutInternal(ctx, req.(*KeyValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _Operations_AppendInternal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationsServer).AppendInternal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/operations.Operations/AppendInternal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationsServer).AppendInternal(ctx, req.(*KeyValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _Operations_DelInternal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationsServer).DelInternal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/operations.Operations/DelInternal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationsServer).DelInternal(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _Operations_Migration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyCost)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationsServer).Migration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/operations.Operations/Migration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationsServer).Migration(ctx, req.(*KeyCost))
	}
	return interceptor(ctx, in, info, handler)
}

func _Operations_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationsServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/operations.Operations/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationsServer).Ping(ctx, req.(*PingMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Operations_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationsServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/operations.Operations/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationsServer).Join(ctx, req.(*JoinMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Operations_RequestJoin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestJoinMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationsServer).RequestJoin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/operations.Operations/RequestJoin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationsServer).RequestJoin(ctx, req.(*RequestJoinMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Operations_LeaveCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestJoinMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationsServer).LeaveCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/operations.Operations/LeaveCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationsServer).LeaveCluster(ctx, req.(*RequestJoinMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// Operations_ServiceDesc is the grpc.ServiceDesc for Operations service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Operations_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "operations.Operations",
	HandlerType: (*OperationsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Operations_Get_Handler,
		},
		{
			MethodName: "Put",
			Handler:    _Operations_Put_Handler,
		},
		{
			MethodName: "Append",
			Handler:    _Operations_Append_Handler,
		},
		{
			MethodName: "Del",
			Handler:    _Operations_Del_Handler,
		},
		{
			MethodName: "GetInternal",
			Handler:    _Operations_GetInternal_Handler,
		},
		{
			MethodName: "PutInternal",
			Handler:    _Operations_PutInternal_Handler,
		},
		{
			MethodName: "AppendInternal",
			Handler:    _Operations_AppendInternal_Handler,
		},
		{
			MethodName: "DelInternal",
			Handler:    _Operations_DelInternal_Handler,
		},
		{
			MethodName: "Migration",
			Handler:    _Operations_Migration_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Operations_Ping_Handler,
		},
		{
			MethodName: "Join",
			Handler:    _Operations_Join_Handler,
		},
		{
			MethodName: "RequestJoin",
			Handler:    _Operations_RequestJoin_Handler,
		},
		{
			MethodName: "LeaveCluster",
			Handler:    _Operations_LeaveCluster_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "operations/operations.proto",
}
