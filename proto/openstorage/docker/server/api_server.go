package server

import (
	"errors"
	"time"

	"github.com/libopenstorage/openstorage/proto/openstorage"
	"github.com/libopenstorage/openstorage/proto/openstorage/docker"
	"go.pedge.io/proto/rpclog"
	"golang.org/x/net/context"
)

type apiServer struct {
	protorpclog.Logger
	openstorageAPIClient openstorage.APIClient
}

func newAPIServer(openstorageAPIClient openstorage.APIClient) *apiServer {
	return &apiServer{
		protorpclog.NewLogger("openstorage.docker.API"),
		openstorageAPIClient,
	}
}

func (a *apiServer) VolumeCreate(ctx context.Context, request *openstorage_docker.VolumeCreateRequest) (response *openstorage_docker.VolumeCreateResponse, err error) {
	defer func(start time.Time) { a.Log(request, response, err, time.Since(start)) }(time.Now())
	openstorageVolumeCreateRequest, err := toOpenstorageVolumeCreateRequest(request)
	if err != nil {
		// note that docker expects the error as part of the actual response
		return toVolumeCreateResponse(err), nil
	}
	_, err = a.openstorageAPIClient.VolumeCreate(ctx, openstorageVolumeCreateRequest)
	return toVolumeCreateResponse(err), nil
}

// TODO(pedge)
func toOpenstorageVolumeCreateRequest(request *openstorage_docker.VolumeCreateRequest) (*openstorage.VolumeCreateRequest, error) {
	return nil, errors.New("not implemented")
}

func toVolumeCreateResponse(err error) *openstorage_docker.VolumeCreateResponse {
	volumeCreateResponse := &openstorage_docker.VolumeCreateResponse{}
	if err != nil {
		volumeCreateResponse.Err = err.Error()
	}
	return volumeCreateResponse
}
