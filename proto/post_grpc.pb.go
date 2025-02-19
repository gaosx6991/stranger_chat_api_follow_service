// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
// source: proto/post.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	PostService_GetUserPosts_FullMethodName  = "/proto.PostService/GetUserPosts"
	PostService_SearchPosts_FullMethodName   = "/proto.PostService/SearchPosts"
	PostService_GetHotPosts_FullMethodName   = "/proto.PostService/GetHotPosts"
	PostService_GetPostDetail_FullMethodName = "/proto.PostService/GetPostDetail"
)

// PostServiceClient is the client API for PostService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PostServiceClient interface {
	GetUserPosts(ctx context.Context, in *GetUserPostsRequest, opts ...grpc.CallOption) (*GetUserPostsResponse, error)
	SearchPosts(ctx context.Context, in *SearchPostsRequest, opts ...grpc.CallOption) (*SearchPostsResponse, error)
	GetHotPosts(ctx context.Context, in *GetHotPostsRequest, opts ...grpc.CallOption) (*GetHotPostsResponse, error)
	GetPostDetail(ctx context.Context, in *GetPostDetailRequest, opts ...grpc.CallOption) (*GetPostDetailResponse, error)
}

type postServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPostServiceClient(cc grpc.ClientConnInterface) PostServiceClient {
	return &postServiceClient{cc}
}

func (c *postServiceClient) GetUserPosts(ctx context.Context, in *GetUserPostsRequest, opts ...grpc.CallOption) (*GetUserPostsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserPostsResponse)
	err := c.cc.Invoke(ctx, PostService_GetUserPosts_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) SearchPosts(ctx context.Context, in *SearchPostsRequest, opts ...grpc.CallOption) (*SearchPostsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchPostsResponse)
	err := c.cc.Invoke(ctx, PostService_SearchPosts_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetHotPosts(ctx context.Context, in *GetHotPostsRequest, opts ...grpc.CallOption) (*GetHotPostsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetHotPostsResponse)
	err := c.cc.Invoke(ctx, PostService_GetHotPosts_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetPostDetail(ctx context.Context, in *GetPostDetailRequest, opts ...grpc.CallOption) (*GetPostDetailResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPostDetailResponse)
	err := c.cc.Invoke(ctx, PostService_GetPostDetail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PostServiceServer is the server API for PostService service.
// All implementations must embed UnimplementedPostServiceServer
// for forward compatibility.
type PostServiceServer interface {
	GetUserPosts(context.Context, *GetUserPostsRequest) (*GetUserPostsResponse, error)
	SearchPosts(context.Context, *SearchPostsRequest) (*SearchPostsResponse, error)
	GetHotPosts(context.Context, *GetHotPostsRequest) (*GetHotPostsResponse, error)
	GetPostDetail(context.Context, *GetPostDetailRequest) (*GetPostDetailResponse, error)
	mustEmbedUnimplementedPostServiceServer()
}

// UnimplementedPostServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPostServiceServer struct{}

func (UnimplementedPostServiceServer) GetUserPosts(context.Context, *GetUserPostsRequest) (*GetUserPostsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserPosts not implemented")
}
func (UnimplementedPostServiceServer) SearchPosts(context.Context, *SearchPostsRequest) (*SearchPostsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchPosts not implemented")
}
func (UnimplementedPostServiceServer) GetHotPosts(context.Context, *GetHotPostsRequest) (*GetHotPostsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHotPosts not implemented")
}
func (UnimplementedPostServiceServer) GetPostDetail(context.Context, *GetPostDetailRequest) (*GetPostDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPostDetail not implemented")
}
func (UnimplementedPostServiceServer) mustEmbedUnimplementedPostServiceServer() {}
func (UnimplementedPostServiceServer) testEmbeddedByValue()                     {}

// UnsafePostServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PostServiceServer will
// result in compilation errors.
type UnsafePostServiceServer interface {
	mustEmbedUnimplementedPostServiceServer()
}

func RegisterPostServiceServer(s grpc.ServiceRegistrar, srv PostServiceServer) {
	// If the following call pancis, it indicates UnimplementedPostServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PostService_ServiceDesc, srv)
}

func _PostService_GetUserPosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserPostsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetUserPosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_GetUserPosts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetUserPosts(ctx, req.(*GetUserPostsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_SearchPosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchPostsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).SearchPosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_SearchPosts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).SearchPosts(ctx, req.(*SearchPostsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetHotPosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetHotPostsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetHotPosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_GetHotPosts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetHotPosts(ctx, req.(*GetHotPostsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetPostDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPostDetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetPostDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_GetPostDetail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetPostDetail(ctx, req.(*GetPostDetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PostService_ServiceDesc is the grpc.ServiceDesc for PostService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PostService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.PostService",
	HandlerType: (*PostServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserPosts",
			Handler:    _PostService_GetUserPosts_Handler,
		},
		{
			MethodName: "SearchPosts",
			Handler:    _PostService_SearchPosts_Handler,
		},
		{
			MethodName: "GetHotPosts",
			Handler:    _PostService_GetHotPosts_Handler,
		},
		{
			MethodName: "GetPostDetail",
			Handler:    _PostService_GetPostDetail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/post.proto",
}
