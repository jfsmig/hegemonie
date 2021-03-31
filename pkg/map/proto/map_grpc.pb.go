// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

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

// MapClient is the client API for Map service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MapClient interface {
	Maps(ctx context.Context, in *ListMapsReq, opts ...grpc.CallOption) (Map_MapsClient, error)
	// Paginated query for the vertices of the graph
	Vertices(ctx context.Context, in *ListVerticesReq, opts ...grpc.CallOption) (Map_VerticesClient, error)
	// Paginated query for the edges of the graph
	Edges(ctx context.Context, in *ListEdgesReq, opts ...grpc.CallOption) (Map_EdgesClient, error)
	// Paginated query of the location occupied by a City
	Cities(ctx context.Context, in *ListCitiesReq, opts ...grpc.CallOption) (Map_CitiesClient, error)
	// Request a path computation on the map
	GetPath(ctx context.Context, in *PathRequest, opts ...grpc.CallOption) (Map_GetPathClient, error)
}

type mapClient struct {
	cc grpc.ClientConnInterface
}

func NewMapClient(cc grpc.ClientConnInterface) MapClient {
	return &mapClient{cc}
}

func (c *mapClient) Maps(ctx context.Context, in *ListMapsReq, opts ...grpc.CallOption) (Map_MapsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Map_ServiceDesc.Streams[0], "/hege.map.Map/Maps", opts...)
	if err != nil {
		return nil, err
	}
	x := &mapMapsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Map_MapsClient interface {
	Recv() (*MapName, error)
	grpc.ClientStream
}

type mapMapsClient struct {
	grpc.ClientStream
}

func (x *mapMapsClient) Recv() (*MapName, error) {
	m := new(MapName)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mapClient) Vertices(ctx context.Context, in *ListVerticesReq, opts ...grpc.CallOption) (Map_VerticesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Map_ServiceDesc.Streams[1], "/hege.map.Map/Vertices", opts...)
	if err != nil {
		return nil, err
	}
	x := &mapVerticesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Map_VerticesClient interface {
	Recv() (*Vertex, error)
	grpc.ClientStream
}

type mapVerticesClient struct {
	grpc.ClientStream
}

func (x *mapVerticesClient) Recv() (*Vertex, error) {
	m := new(Vertex)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mapClient) Edges(ctx context.Context, in *ListEdgesReq, opts ...grpc.CallOption) (Map_EdgesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Map_ServiceDesc.Streams[2], "/hege.map.Map/Edges", opts...)
	if err != nil {
		return nil, err
	}
	x := &mapEdgesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Map_EdgesClient interface {
	Recv() (*Edge, error)
	grpc.ClientStream
}

type mapEdgesClient struct {
	grpc.ClientStream
}

func (x *mapEdgesClient) Recv() (*Edge, error) {
	m := new(Edge)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mapClient) Cities(ctx context.Context, in *ListCitiesReq, opts ...grpc.CallOption) (Map_CitiesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Map_ServiceDesc.Streams[3], "/hege.map.Map/Cities", opts...)
	if err != nil {
		return nil, err
	}
	x := &mapCitiesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Map_CitiesClient interface {
	Recv() (*CityLocation, error)
	grpc.ClientStream
}

type mapCitiesClient struct {
	grpc.ClientStream
}

func (x *mapCitiesClient) Recv() (*CityLocation, error) {
	m := new(CityLocation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mapClient) GetPath(ctx context.Context, in *PathRequest, opts ...grpc.CallOption) (Map_GetPathClient, error) {
	stream, err := c.cc.NewStream(ctx, &Map_ServiceDesc.Streams[4], "/hege.map.Map/GetPath", opts...)
	if err != nil {
		return nil, err
	}
	x := &mapGetPathClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Map_GetPathClient interface {
	Recv() (*PathElement, error)
	grpc.ClientStream
}

type mapGetPathClient struct {
	grpc.ClientStream
}

func (x *mapGetPathClient) Recv() (*PathElement, error) {
	m := new(PathElement)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MapServer is the server API for Map service.
// All implementations must embed UnimplementedMapServer
// for forward compatibility
type MapServer interface {
	Maps(*ListMapsReq, Map_MapsServer) error
	// Paginated query for the vertices of the graph
	Vertices(*ListVerticesReq, Map_VerticesServer) error
	// Paginated query for the edges of the graph
	Edges(*ListEdgesReq, Map_EdgesServer) error
	// Paginated query of the location occupied by a City
	Cities(*ListCitiesReq, Map_CitiesServer) error
	// Request a path computation on the map
	GetPath(*PathRequest, Map_GetPathServer) error
	mustEmbedUnimplementedMapServer()
}

// UnimplementedMapServer must be embedded to have forward compatible implementations.
type UnimplementedMapServer struct {
}

func (UnimplementedMapServer) Maps(*ListMapsReq, Map_MapsServer) error {
	return status.Errorf(codes.Unimplemented, "method Maps not implemented")
}
func (UnimplementedMapServer) Vertices(*ListVerticesReq, Map_VerticesServer) error {
	return status.Errorf(codes.Unimplemented, "method Vertices not implemented")
}
func (UnimplementedMapServer) Edges(*ListEdgesReq, Map_EdgesServer) error {
	return status.Errorf(codes.Unimplemented, "method Edges not implemented")
}
func (UnimplementedMapServer) Cities(*ListCitiesReq, Map_CitiesServer) error {
	return status.Errorf(codes.Unimplemented, "method Cities not implemented")
}
func (UnimplementedMapServer) GetPath(*PathRequest, Map_GetPathServer) error {
	return status.Errorf(codes.Unimplemented, "method GetPath not implemented")
}
func (UnimplementedMapServer) mustEmbedUnimplementedMapServer() {}

// UnsafeMapServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MapServer will
// result in compilation errors.
type UnsafeMapServer interface {
	mustEmbedUnimplementedMapServer()
}

func RegisterMapServer(s grpc.ServiceRegistrar, srv MapServer) {
	s.RegisterService(&Map_ServiceDesc, srv)
}

func _Map_Maps_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListMapsReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MapServer).Maps(m, &mapMapsServer{stream})
}

type Map_MapsServer interface {
	Send(*MapName) error
	grpc.ServerStream
}

type mapMapsServer struct {
	grpc.ServerStream
}

func (x *mapMapsServer) Send(m *MapName) error {
	return x.ServerStream.SendMsg(m)
}

func _Map_Vertices_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListVerticesReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MapServer).Vertices(m, &mapVerticesServer{stream})
}

type Map_VerticesServer interface {
	Send(*Vertex) error
	grpc.ServerStream
}

type mapVerticesServer struct {
	grpc.ServerStream
}

func (x *mapVerticesServer) Send(m *Vertex) error {
	return x.ServerStream.SendMsg(m)
}

func _Map_Edges_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListEdgesReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MapServer).Edges(m, &mapEdgesServer{stream})
}

type Map_EdgesServer interface {
	Send(*Edge) error
	grpc.ServerStream
}

type mapEdgesServer struct {
	grpc.ServerStream
}

func (x *mapEdgesServer) Send(m *Edge) error {
	return x.ServerStream.SendMsg(m)
}

func _Map_Cities_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListCitiesReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MapServer).Cities(m, &mapCitiesServer{stream})
}

type Map_CitiesServer interface {
	Send(*CityLocation) error
	grpc.ServerStream
}

type mapCitiesServer struct {
	grpc.ServerStream
}

func (x *mapCitiesServer) Send(m *CityLocation) error {
	return x.ServerStream.SendMsg(m)
}

func _Map_GetPath_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PathRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MapServer).GetPath(m, &mapGetPathServer{stream})
}

type Map_GetPathServer interface {
	Send(*PathElement) error
	grpc.ServerStream
}

type mapGetPathServer struct {
	grpc.ServerStream
}

func (x *mapGetPathServer) Send(m *PathElement) error {
	return x.ServerStream.SendMsg(m)
}

// Map_ServiceDesc is the grpc.ServiceDesc for Map service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Map_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hege.map.Map",
	HandlerType: (*MapServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Maps",
			Handler:       _Map_Maps_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Vertices",
			Handler:       _Map_Vertices_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Edges",
			Handler:       _Map_Edges_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Cities",
			Handler:       _Map_Cities_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetPath",
			Handler:       _Map_GetPath_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "map.proto",
}