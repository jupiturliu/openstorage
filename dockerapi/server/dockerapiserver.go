package dockerapiserver

import (
	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/dockerapi"
)

func NewVolumeAPIServer(volumeAPIClient api.VolumeAPIClient) dockerapi.VolumeAPIServer {
	return newVolumeAPIServer(volumeAPIClient)
}
