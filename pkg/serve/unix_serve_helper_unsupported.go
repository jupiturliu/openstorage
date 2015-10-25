// +build !linux,!freebsd

package serve

import "errors"

var (
	errOnlySupportedOnLinuxAndFreeBSD = errors.New("unix socket creation is only supported on linux and freebsd")
)

func newUnixServeHelper(
	volumeDriverName string,
	group string,
) (*serveHelper, error) {
	return nil, errOnlySupportedOnLinuxAndFreeBSD
}
