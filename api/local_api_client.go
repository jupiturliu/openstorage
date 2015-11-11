package api

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// TODO(pedge): I will write a protoc plugin soon that auto-generates this

type localAPIClient struct {
	delegate APIServer
}

func newLocalAPIClient(delegate APIServer) *localAPIClient {
	return &localAPIClient{delegate}
}

func (a *localAPIClient) VolumeCreate(ctx context.Context, request *VolumeCreateRequest, _ ...grpc.CallOption) (*VolumeCreateResponse, error) {
	return a.delegate.VolumeCreate(ctx, request)
}
