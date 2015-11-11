// TODO(pedge): package name

package server

import (
	"github.com/libopenstorage/openstorage/proto/openstorage"
	"github.com/libopenstorage/openstorage/proto/openstoragedocker"
)

func NewAPIServer(openstorageAPIClient openstorage.APIClient) openstoragedocker.APIServer {
	return newAPIServer(openstorageAPIClient)
}
