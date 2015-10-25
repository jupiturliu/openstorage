package server

import (
	"fmt"
	"time"

	"github.com/libopenstorage/openstorage/proto/openstorage"
	"github.com/libopenstorage/openstorage/proto/openstorage/docker"
	"go.pedge.io/pkg/map"
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

func toOpenstorageVolumeCreateRequest(request *openstorage_docker.VolumeCreateRequest) (*openstorage.VolumeCreateRequest, error) {
	openstorageVolumeLocator, err := toOpenstorageVolumeLocator(request)
	if err != nil {
		return nil, err
	}
	openstorageVolumeSource, err := toOpenstorageVolumeSource(request)
	if err != nil {
		return nil, err
	}
	openstorageVolumeSpec, err := toOpenstorageVolumeSpec(request)
	if err != nil {
		return nil, err
	}
	return &openstorage.VolumeCreateRequest{
		VolumeLocator: openstorageVolumeLocator,
		VolumeSource:  openstorageVolumeSource,
		VolumeSpec:    openstorageVolumeSpec,
	}, nil
}

func toOpenstorageVolumeLocator(request *openstorage_docker.VolumeCreateRequest) (*openstorage.VolumeLocator, error) {
	// TODO(pedge): what to do with labels? one idea is to
	// have labels be any fields that do not map into something
	// specific for VolumeLocator, VolumeSource, VolumeSpec,
	// but we also could ignore them for docker api requests
	return &openstorage.VolumeLocator{
		Name: request.Name,
	}, nil
}

func toOpenstorageVolumeSource(request *openstorage_docker.VolumeCreateRequest) (*openstorage.VolumeSource, error) {
	openstorageVolumeSource := &openstorage.VolumeSource{}
	if len(request.Opts) == 0 {
		return openstorageVolumeSource, nil
	}
	opts := pkgmap.StringStringMap(request.Opts)
	parentVolumeID, err := opts.GetString("parent_volume_id")
	if err != nil {
		return nil, err
	}
	openstorageVolumeSource.ParentVolumeId = parentVolumeID
	seedURI, err := opts.GetString("seed_uri")
	if err != nil {
		return nil, err
	}
	openstorageVolumeSource.SeedUri = seedURI
	return openstorageVolumeSource, nil
}

func toOpenstorageVolumeSpec(request *openstorage_docker.VolumeCreateRequest) (*openstorage.VolumeSpec, error) {
	openstorageVolumeSpec := &openstorage.VolumeSpec{}
	if len(request.Opts) == 0 {
		return openstorageVolumeSpec, nil
	}
	opts := pkgmap.StringStringMap(request.Opts)
	ephemeral, err := opts.GetBool("ephemeral")
	if err != nil {
		return nil, err
	}
	openstorageVolumeSpec.Ephemeral = ephemeral
	sizeBytes, err := opts.GetUint64("size_bytes")
	if err != nil {
		return nil, err
	}
	openstorageVolumeSpec.SizeBytes = sizeBytes
	fsTypeObj, err := opts.GetString("fs_type")
	if err != nil {
		return nil, err
	}
	// TODO(pedge): should we look up the value based on that
	// actual string, or strings.TrimSuffix(strings.ToLower(fsTypeObk), "fs_type_")?
	fsType, ok := openstorage.FSType_value[fsTypeObj]
	if !ok {
		return nil, fmt.Errorf("No openstorage.FSType value %s", fsTypeObj)
	}
	openstorageVolumeSpec.FsType = openstorage.FSType(fsType)
	blockSize, err := opts.GetInt64("block_size")
	if err != nil {
		return nil, err
	}
	openstorageVolumeSpec.BlockSize = blockSize
	haLevel, err := opts.GetInt32("ha_level")
	if err != nil {
		return nil, err
	}
	openstorageVolumeSpec.HaLevel = haLevel
	cosObj, err := opts.GetString("cos")
	if err != nil {
		return nil, err
	}
	cos, ok := openstorage.COS_value[cosObj]
	if !ok {
		return nil, fmt.Errorf("No openstorage.COS value %s", cosObj)
	}
	openstorageVolumeSpec.Cos = openstorage.COS(cos)
	deduplicate, err := opts.GetBool("deduplicate")
	if err != nil {
		return nil, err
	}
	openstorageVolumeSpec.Deduplicate = deduplicate
	snapshotIntervalMin, err := opts.GetUint32("snapshot_interval_min")
	if err != nil {
		return nil, err
	}
	openstorageVolumeSpec.SnapshotIntervalMin = snapshotIntervalMin
	// TODO(pedge): labels?
	return openstorageVolumeSpec, nil
}

func toVolumeCreateResponse(err error) *openstorage_docker.VolumeCreateResponse {
	volumeCreateResponse := &openstorage_docker.VolumeCreateResponse{}
	if err != nil {
		volumeCreateResponse.Err = err.Error()
	}
	return volumeCreateResponse
}
