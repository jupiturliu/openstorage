package serve

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/sockets"
)

func newTCPServeHelper(
	volumeDriverName string,
	address string,
) (*serveHelper, error) {
	start := make(chan struct{})
	listener, err := sockets.NewTCPSocket(address, nil, start)
	if err != nil {
		return nil, err
	}
	spec, err := writeSpec(volumeDriverName, listener.Addr().String())
	if err != nil {
		return nil, err
	}
	return &serveHelper{
		listener,
		address,
		spec,
		start,
	}, nil
}

func writeSpec(name string, address string) (string, error) {
	dir := filepath.Join(TCPSpecDir, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	spec := filepath.Join(dir, name+".spec")
	if strings.HasPrefix(address, "[::]:") {
		address = fmt.Sprintf("0.0.0.0:%s", strings.TrimPrefix(address, "[::]:"))
	}
	url := "tcp://" + address
	if err := ioutil.WriteFile(spec, []byte(url), 0644); err != nil {
		return "", err
	}
	return spec, nil
}
