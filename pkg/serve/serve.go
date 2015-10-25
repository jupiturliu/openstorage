package serve

import "net"

const (
	// TCPSpecDir is the directory where plugin specifications are based.
	TCPSpecDir = "/etc/docker/plugins"
	// UnixSockDir is the directory where unix sockets are based.
	UnixSockDir = "/run/docker/plugins"
)

// ServeHelper provides helpers for serving docker volume plugins.
//
// TODO(pedge): name?
type ServeHelper interface {
	Listener() net.Listener
	Address() string
	OnStart()
	OnShutdownInitiated()
}

// NewTCPServeHelper constructs a new ServeHelper for TCP.
func NewTCPServeHelper(
	volumeDriverName string,
	address string,
) (ServeHelper, error) {
	return nil, nil
}

// NewUnixServeHelper constructs a new ServeHelper for Unix sockets.
func NewUnixServeHelper(
	volumeDriverName string,
	group string,
) (ServeHelper, error) {
	return nil, nil
}
