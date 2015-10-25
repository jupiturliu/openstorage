package fuse

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/libopenstorage/openstorage/proto/openstorage"
	"github.com/libopenstorage/openstorage/volume"
)

// Provider provides fuse.FS and fuse.MountOptions, given an *openstorage.VolumeSpec.
type Provider interface {
	GetFS(volumeSpec *openstorage.VolumeSpec) (fs.FS, error)
	GetMountOptions(volumeSpec *openstorage.VolumeSpec) ([]fuse.MountOption, error)
}

// NewVolumeDriver creates a new volume.VolumeDriver for fuse.
func NewVolumeDriver(name string, baseDirPath string, provider Provider) volume.VolumeDriver {
	return newVolumeDriver(name, baseDirPath, provider)
}
