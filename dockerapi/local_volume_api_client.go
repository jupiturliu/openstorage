package dockerapi

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type localVolumeAPIClient struct {
	delegate VolumeAPIServer
}

func newLocalVolumeAPIClient(delegate VolumeAPIServer) *localVolumeAPIClient {
	return &localVolumeAPIClient{delegate}
}

func (a *localVolumeAPIClient) VolumeCreate(ctx context.Context, request *NameOptsRequest, opts ...grpc.CallOption) (*ErrResponse, error) {
	return a.delegate.VolumeCreate(ctx, request)
}

func (a *localVolumeAPIClient) VolumeRemove(ctx context.Context, request *NameRequest, opts ...grpc.CallOption) (*ErrResponse, error) {
	return a.delegate.VolumeRemove(ctx, request)
}

func (a *localVolumeAPIClient) VolumePath(ctx context.Context, request *NameRequest, opts ...grpc.CallOption) (*MountpointErrResponse, error) {
	return a.delegate.VolumePath(ctx, request)
}

func (a *localVolumeAPIClient) VolumeMount(ctx context.Context, request *NameRequest, opts ...grpc.CallOption) (*MountpointErrResponse, error) {
	return a.delegate.VolumeMount(ctx, request)
}

func (a *localVolumeAPIClient) VolumeUnmount(ctx context.Context, request *NameRequest, opts ...grpc.CallOption) (*ErrResponse, error) {
	return a.delegate.VolumeUnmount(ctx, request)
}
