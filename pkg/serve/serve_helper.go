package serve

import (
	"net"
	"os"
)

type serveHelper struct {
	l     net.Listener
	addr  string
	spec  string
	start chan struct{}
}

func (s *serveHelper) Listener() net.Listener {
	return s.l
}

func (s *serveHelper) Address() string {
	return s.addr
}

func (s *serveHelper) OnStart() {
	close(s.start)
}

func (s *serveHelper) OnShutdownInitiated() {
	_ = os.Remove(s.spec)
}
