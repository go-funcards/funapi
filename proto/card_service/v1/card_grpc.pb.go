// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: card_service/v1/card.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CardClient is the client API for Card service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CardClient interface {
	CreateCard(ctx context.Context, in *CreateCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UpdateCard(ctx context.Context, in *UpdateCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UpdateManyCards(ctx context.Context, in *UpdateManyCardsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteCard(ctx context.Context, in *DeleteCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetCards(ctx context.Context, in *CardsRequest, opts ...grpc.CallOption) (*CardsResponse, error)
}

type cardClient struct {
	cc grpc.ClientConnInterface
}

func NewCardClient(cc grpc.ClientConnInterface) CardClient {
	return &cardClient{cc}
}

func (c *cardClient) CreateCard(ctx context.Context, in *CreateCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.v1.Card/CreateCard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardClient) UpdateCard(ctx context.Context, in *UpdateCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.v1.Card/UpdateCard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardClient) UpdateManyCards(ctx context.Context, in *UpdateManyCardsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.v1.Card/UpdateManyCards", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardClient) DeleteCard(ctx context.Context, in *DeleteCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.v1.Card/DeleteCard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardClient) GetCards(ctx context.Context, in *CardsRequest, opts ...grpc.CallOption) (*CardsResponse, error) {
	out := new(CardsResponse)
	err := c.cc.Invoke(ctx, "/proto.v1.Card/GetCards", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CardServer is the server API for Card service.
// All implementations must embed UnimplementedCardServer
// for forward compatibility
type CardServer interface {
	CreateCard(context.Context, *CreateCardRequest) (*emptypb.Empty, error)
	UpdateCard(context.Context, *UpdateCardRequest) (*emptypb.Empty, error)
	UpdateManyCards(context.Context, *UpdateManyCardsRequest) (*emptypb.Empty, error)
	DeleteCard(context.Context, *DeleteCardRequest) (*emptypb.Empty, error)
	GetCards(context.Context, *CardsRequest) (*CardsResponse, error)
	mustEmbedUnimplementedCardServer()
}

// UnimplementedCardServer must be embedded to have forward compatible implementations.
type UnimplementedCardServer struct {
}

func (UnimplementedCardServer) CreateCard(context.Context, *CreateCardRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCard not implemented")
}
func (UnimplementedCardServer) UpdateCard(context.Context, *UpdateCardRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCard not implemented")
}
func (UnimplementedCardServer) UpdateManyCards(context.Context, *UpdateManyCardsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateManyCards not implemented")
}
func (UnimplementedCardServer) DeleteCard(context.Context, *DeleteCardRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCard not implemented")
}
func (UnimplementedCardServer) GetCards(context.Context, *CardsRequest) (*CardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCards not implemented")
}
func (UnimplementedCardServer) mustEmbedUnimplementedCardServer() {}

// UnsafeCardServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CardServer will
// result in compilation errors.
type UnsafeCardServer interface {
	mustEmbedUnimplementedCardServer()
}

func RegisterCardServer(s grpc.ServiceRegistrar, srv CardServer) {
	s.RegisterService(&Card_ServiceDesc, srv)
}

func _Card_CreateCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardServer).CreateCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.v1.Card/CreateCard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardServer).CreateCard(ctx, req.(*CreateCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Card_UpdateCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardServer).UpdateCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.v1.Card/UpdateCard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardServer).UpdateCard(ctx, req.(*UpdateCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Card_UpdateManyCards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateManyCardsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardServer).UpdateManyCards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.v1.Card/UpdateManyCards",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardServer).UpdateManyCards(ctx, req.(*UpdateManyCardsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Card_DeleteCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardServer).DeleteCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.v1.Card/DeleteCard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardServer).DeleteCard(ctx, req.(*DeleteCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Card_GetCards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CardsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardServer).GetCards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.v1.Card/GetCards",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardServer).GetCards(ctx, req.(*CardsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Card_ServiceDesc is the grpc.ServiceDesc for Card service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Card_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.v1.Card",
	HandlerType: (*CardServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCard",
			Handler:    _Card_CreateCard_Handler,
		},
		{
			MethodName: "UpdateCard",
			Handler:    _Card_UpdateCard_Handler,
		},
		{
			MethodName: "UpdateManyCards",
			Handler:    _Card_UpdateManyCards_Handler,
		},
		{
			MethodName: "DeleteCard",
			Handler:    _Card_DeleteCard_Handler,
		},
		{
			MethodName: "GetCards",
			Handler:    _Card_GetCards_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "card_service/v1/card.proto",
}
