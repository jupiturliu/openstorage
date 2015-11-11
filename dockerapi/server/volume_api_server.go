package dockerapiserver

import (
	"time"

	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/dockerapi"
	"go.pedge.io/pkg/map"
	"go.pedge.io/proto/rpclog"
	"golang.org/x/net/context"
)

type volumeAPIServer struct {
	protorpclog.Logger
	delegate api.VolumeAPIClient
}

func newVolumeAPIServer(delegate api.VolumeAPIClient) *volumeAPIServer {
	return &volumeAPIServer{
		protorpclog.NewLogger("openstorage.docker.api.VolumeAPI"),
		delegate,
	}
}

func (a *volumeAPIServer) VolumeCreate(ctx context.Context, request *dockerapi.NameOptsRequest) (response *dockerapi.ErrResponse, err error) {
	defer func(start time.Time) { a.Log(request, response, err, time.Since(start)) }(time.Now())
	return doNameOptsToErr(ctx, request, a.volumeCreate)
}

func (a *volumeAPIServer) volumeCreate(ctx context.Context, name string, opts pkgmap.StringStringMap) error {
	volumeCreateRequest, err := toVolumeCreateRequest(name, opts)
	if err != nil {
		return err
	}
	_, err = a.delegate.VolumeCreate(ctx, volumeCreateRequest)
	return err
}

func (a *volumeAPIServer) VolumeRemove(ctx context.Context, request *dockerapi.NameRequest) (response *dockerapi.ErrResponse, err error) {
	defer func(start time.Time) { a.Log(request, response, err, time.Since(start)) }(time.Now())
	return doNameToErr(ctx, request, a.volumeRemove)
}

func (a *volumeAPIServer) volumeRemove(ctx context.Context, name string) error {
	return nil
}

func (a *volumeAPIServer) VolumePath(ctx context.Context, request *dockerapi.NameRequest) (response *dockerapi.MountpointErrResponse, err error) {
	defer func(start time.Time) { a.Log(request, response, err, time.Since(start)) }(time.Now())
	return doNameToMountpointErr(ctx, request, a.volumePath)
}

func (a *volumeAPIServer) volumePath(ctx context.Context, name string) (string, error) {
	return "", nil
}

func (a *volumeAPIServer) VolumeMount(ctx context.Context, request *dockerapi.NameRequest) (response *dockerapi.MountpointErrResponse, err error) {
	defer func(start time.Time) { a.Log(request, response, err, time.Since(start)) }(time.Now())
	return doNameToMountpointErr(ctx, request, a.volumeMount)
}

func (a *volumeAPIServer) volumeMount(ctx context.Context, name string) (string, error) {
	return "", nil
}

func (a *volumeAPIServer) VolumeUnmount(ctx context.Context, request *dockerapi.NameRequest) (response *dockerapi.ErrResponse, err error) {
	defer func(start time.Time) { a.Log(request, response, err, time.Since(start)) }(time.Now())
	return doNameToErr(ctx, request, a.volumeUnmount)
}

func (a *volumeAPIServer) volumeUnmount(ctx context.Context, name string) error {
	return nil
}

func fromNameOptsRequest(request *dockerapi.NameOptsRequest) (string, pkgmap.StringStringMap) {
	return request.Name, pkgmap.StringStringMap(request.Opts).Copy()
}

func fromNameRequest(request *dockerapi.NameRequest) string {
	return request.Name
}

func toErrResponse(err error) (*dockerapi.ErrResponse, error) {
	response := &dockerapi.ErrResponse{}
	if err != nil {
		response.Err = err.Error()
	}
	return response, nil
}

func toMountpointErrResponse(mountpoint string, err error) (*dockerapi.MountpointErrResponse, error) {
	response := &dockerapi.MountpointErrResponse{
		Mountpoint: mountpoint,
	}
	if err != nil {
		response.Err = err.Error()
	}
	return response, nil
}

func doNameOptsToErr(ctx context.Context, request *dockerapi.NameOptsRequest, f func(context.Context, string, pkgmap.StringStringMap) error) (*dockerapi.ErrResponse, error) {
	name, opts := fromNameOptsRequest(request)
	return toErrResponse(f(ctx, name, opts))
}

func doNameToErr(ctx context.Context, request *dockerapi.NameRequest, f func(context.Context, string) error) (*dockerapi.ErrResponse, error) {
	return toErrResponse(f(ctx, fromNameRequest(request)))
}

func doNameToMountpointErr(ctx context.Context, request *dockerapi.NameRequest, f func(context.Context, string) (string, error)) (*dockerapi.MountpointErrResponse, error) {
	mountpoint, err := f(ctx, fromNameRequest(request))
	return toMountpointErrResponse(mountpoint, err)
}

func toVolumeCreateRequest(name string, opts pkgmap.StringStringMap) (*api.VolumeCreateRequest, error) {
	volumeSource := &api.VolumeSource{}
	parentVolumeID, err := opts.GetString("parent_volume_id")
	if err != nil {
		return nil, err
	}
	volumeSource.ParentVolumeId = parentVolumeID
	seedURI, err := opts.GetString("seed_uri")
	if err != nil {
		return nil, err
	}
	volumeSource.SeedUri = seedURI
	volumeSpec := &api.VolumeSpec{}
	ephemeral, err := opts.GetBool("ephemeral")
	if err != nil {
		return nil, err
	}
	volumeSpec.Ephemeral = ephemeral
	sizeBytes, err := opts.GetUint64("size_bytes")
	if err != nil {
		return nil, err
	}
	volumeSpec.SizeBytes = sizeBytes
	fsTypeObj, err := opts.GetString("fs_type")
	if err != nil {
		return nil, err
	}
	fsType, err := api.FSTypeSimpleValueOf(fsTypeObj)
	if err != nil {
		return nil, err
	}
	volumeSpec.FsType = fsType
	blockSize, err := opts.GetInt64("block_size")
	if err != nil {
		return nil, err
	}
	volumeSpec.BlockSize = blockSize
	haLevel, err := opts.GetInt32("ha_level")
	if err != nil {
		return nil, err
	}
	volumeSpec.HaLevel = haLevel
	cosObj, err := opts.GetString("cos")
	if err != nil {
		return nil, err
	}
	cos, err := api.COSSimpleValueOf(cosObj)
	if err != nil {
		return nil, err
	}
	volumeSpec.Cos = cos
	deduplicate, err := opts.GetBool("deduplicate")
	if err != nil {
		return nil, err
	}
	volumeSpec.Deduplicate = deduplicate
	snapshotIntervalMin, err := opts.GetUint32("snapshot_interval_min")
	if err != nil {
		return nil, err
	}
	volumeSpec.SnapshotIntervalMin = snapshotIntervalMin
	return &api.VolumeCreateRequest{
		// TODO(pedge): what to do with labels? one idea is to
		// have labels be any fields that do not map into something
		// specific for VolumeLocator, VolumeSource, VolumeSpec,
		// but we also could ignore them for docker api requests
		VolumeLocator: &api.VolumeLocator{
			Name: name,
		},
		VolumeSource: volumeSource,
		// TODO(pedge): labels?
		VolumeSpec: volumeSpec,
	}, nil
}
