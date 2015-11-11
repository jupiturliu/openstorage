package api

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// TODO(pedge): I will write a protoc plugin soon that auto-generates this

type localVolumeAPIClient struct {
	delegate VolumeAPIServer
}

func newLocalVolumeAPIClient(delegate VolumeAPIServer) *localVolumeAPIClient {
	return &localVolumeAPIClient{delegate}
}

func (a *localVolumeAPIClient) VolumeCreate(ctx context.Context, request *VolumeCreateRequest, opts ...grpc.CallOption) (*VolumeCreateResponse, error) {
	return a.delegate.VolumeCreate(ctx, request)
}
