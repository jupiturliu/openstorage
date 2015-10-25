// +build linux freebsd

package serve

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/sockets"
)

func newUnixServeHelper(
	volumeDriverName string,
	group string,
) (*serveHelper, error) {
	start := make(chan struct{})
	path, err := fullSocketAddress(volumeDriverName)
	if err != nil {
		return nil, err
	}
	listener, err := sockets.NewUnixSocket(path, group, start)
	if err != nil {
		return nil, err
	}
	return &serveHelper{
		listener,
		volumeDriverName,
		path,
		start,
	}, nil
}

func fullSocketAddress(address string) (string, error) {
	dir := filepath.Join(UnixSockDir, address)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, address+".sock"), nil
}
