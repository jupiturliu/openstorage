package vfs

import (
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/proto/openstorage"
	"github.com/libopenstorage/openstorage/volume"
	"github.com/pborman/uuid"
	"github.com/portworx/kvdb"
)

const (
	Name       = "vfs"
	Type       = api.File
	volumeBase = "/var/lib/osd/"
)

type driver struct {
	*volume.DefaultBlockDriver
	*volume.DefaultEnumerator
	*volume.SnapshotNotSupported
}

// Init Driver intialization.
func Init(params volume.DriverParams) (volume.VolumeDriver, error) {
	return &driver{
		DefaultEnumerator: volume.NewDefaultEnumerator(Name, kvdb.Instance())}, nil
}

func (d *driver) String() string {
	return Name
}

func (d *driver) Type() api.DriverType {
	return Type
}

func (d *driver) Create(locator *openstorage.VolumeLocator, source *openstorage.VolumeSource, spec *openstorage.VolumeSpec) (string, error) {

	volumeID := uuid.New()
	volumeID = strings.TrimSuffix(volumeID, "\n")

	// Create a directory on the Local machine with this UUID.
	err := os.MkdirAll(path.Join(volumeBase, volumeID), 0744)
	if err != nil {
		log.Println(err)
		return api.BadVolumeID, err
	}

	v := &api.Volume{
		ID:         volumeID,
		Locator:    locator,
		Ctime:      time.Now(),
		Spec:       spec,
		LastScan:   time.Now(),
		Format:     openstorage.FSType_FS_TYPE_VFS,
		State:      api.VolumeAvailable,
		Status:     api.Up,
		DevicePath: path.Join(volumeBase, volumeID),
	}

	err = d.CreateVol(v)
	if err != nil {
		return api.BadVolumeID, err
	}

	err = d.UpdateVol(v)
	return v.ID, err

}

func (d *driver) Delete(volumeID string) error {

	// Check if volume exists
	_, err := d.GetVol(volumeID)
	if err != nil {
		log.Println("Volume not found ", err)
		return err
	}

	// Delete the directory
	os.RemoveAll(path.Join(volumeBase, volumeID))

	err = d.DeleteVol(volumeID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

// Mount volume at specified path
// Errors ErrEnoEnt, ErrVolDetached may be returned.
func (d *driver) Mount(volumeID string, mountpath string) error {

	v, err := d.GetVol(volumeID)
	if err != nil {
		log.Println(err)
		return err
	}
	syscall.Unmount(mountpath, 0)
	// TODO(pedge): fs type simple string could result in "none"
	err = syscall.Mount(path.Join(volumeBase, volumeID), mountpath, v.Spec.FsType.SimpleString(), syscall.MS_BIND, "")
	if err != nil {
		log.Printf("Cannot mount %s at %s because %+v",
			path.Join(volumeBase, volumeID), mountpath, err)
		return err
	}

	v.AttachPath = mountpath
	err = d.UpdateVol(v)

	return nil
}

// Unmount volume at specified path
// Errors ErrEnoEnt, ErrVolDetached may be returned.
func (d *driver) Unmount(volumeID string, mountpath string) error {
	v, err := d.GetVol(volumeID)
	if err != nil {
		return err
	}
	if v.AttachPath == "" {
		return fmt.Errorf("Device %v not mounted", volumeID)
	}
	err = syscall.Unmount(v.AttachPath, 0)
	if err != nil {
		return err
	}

	v.AttachPath = ""
	err = d.UpdateVol(v)
	return nil
}

// Set update volume with specified parameters.
func (d *driver) Set(volumeID string, locator *openstorage.VolumeLocator, spec *openstorage.VolumeSpec) error {
	if spec != nil {
		return volume.ErrNotSupported
	}
	v, err := d.GetVol(volumeID)
	if err != nil {
		return err
	}
	if locator != nil {
		v.Locator = *locator
	}
	err = d.UpdateVol(v)
	return err
}

// Stats Not Supported.
func (d *driver) Stats(volumeID string) (api.Stats, error) {
	return api.Stats{}, volume.ErrNotSupported
}

// Alerts Not Supported.
func (d *driver) Alerts(volumeID string) (api.Alerts, error) {
	return api.Alerts{}, volume.ErrNotSupported
}

// Status returns a set of key-value pairs which give low
// level diagnostic status about this driver.
func (d *driver) Status() [][2]string {
	return [][2]string{}
}

// Shutdown and cleanup.
func (d *driver) Shutdown() {
	log.Debugf("%s Shutting down", Name)
}

func init() {
	volume.Register(Name, Init)
}
