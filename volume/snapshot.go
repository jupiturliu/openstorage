package volume

import (
	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/proto/openstorage"
)

type SnapshotNotSupported struct {
}

func (s *SnapshotNotSupported) Snapshot(volumeID string, readonly bool, locator *openstorage.VolumeLocator) (string, error) {
	return api.BadVolumeID, ErrNotSupported
}
