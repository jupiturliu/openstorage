// TODO(pedge): package name

package server

import (
	"github.com/libopenstorage/openstorage/proto/openstorage"
	"github.com/libopenstorage/openstorage/proto/openstorage/docker"
)

func NewAPIServer(openstorageAPIClient openstorage.APIClient) openstorage_docker.APIServer {
	return newAPIServer(openstorageAPIClient)
}
