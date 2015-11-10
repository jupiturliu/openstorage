package buse

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pborman/uuid"

	"github.com/portworx/kvdb"

	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/cluster"
	"github.com/libopenstorage/openstorage/proto/openstorage"
	"github.com/libopenstorage/openstorage/volume"
)

const (
	Name          = "buse"
	Type          = api.Block | api.Graph
	BuseDBKey     = "OpenStorageBuseKey"
	BuseMountPath = "/var/lib/openstorage/buse/"
)

// Implements the open storage volume interface.
type driver struct {
	*volume.DefaultEnumerator
	buseDevices map[string]*buseDev
}

// Implements the Device interface.
type buseDev struct {
	file string
	f    *os.File
	nbd  *NBD
}

func (d *buseDev) ReadAt(b []byte, off int64) (n int, err error) {
	return d.f.ReadAt(b, off)
}

func (d *buseDev) WriteAt(b []byte, off int64) (n int, err error) {
	return d.f.WriteAt(b, off)
}

func copyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func Init(params volume.DriverParams) (volume.VolumeDriver, error) {
	inst := &driver{
		DefaultEnumerator: volume.NewDefaultEnumerator(Name, kvdb.Instance()),
	}

	inst.buseDevices = make(map[string]*buseDev)

	err := os.MkdirAll(BuseMountPath, 0744)
	if err != nil {
		return nil, err
	}

	volumeInfo, err := inst.DefaultEnumerator.Enumerate(
		&openstorage.VolumeLocator{},
		nil)
	if err == nil {
		for _, info := range volumeInfo {
			if info.Status == "" {
				info.Status = api.Up
				inst.UpdateVol(&info)
			}
		}
	} else {
		log.Println("Could not enumerate Volumes, ", err)
	}

	c, err := cluster.Inst()
	if err != nil {
		log.Println("BUSE initializing in clustered mode")
		c.AddEventListener(inst)
	} else {
		log.Println("BUSE initializing in single node mode")
	}

	log.Println("BUSE initialized and driver mounted at: ", BuseMountPath)
	return inst, nil
}

//
// These functions below implement the volume driver interface.
//

func (d *driver) String() string {
	return Name
}

func (d *driver) Type() api.DriverType {
	return Type
}

// Status diagnostic information
func (d *driver) Status() [][2]string {
	return [][2]string{}
}

func (d *driver) Create(locator *openstorage.VolumeLocator, source *openstorage.VolumeSource, spec *openstorage.VolumeSpec) (string, error) {
	volumeID := uuid.New()
	volumeID = strings.TrimSuffix(volumeID, "\n")

	if spec.SizeBytes == 0 {
		return api.BadVolumeID, fmt.Errorf("Volume size cannot be zero", "buse")
	}

	if spec.FsType == openstorage.FSType_FS_TYPE_NONE {
		return api.BadVolumeID, fmt.Errorf("Missing volume format", "buse")
	}

	// Create a file on the local buse path with this UUID.
	buseFile := path.Join(BuseMountPath, volumeID)
	f, err := os.Create(buseFile)
	if err != nil {
		log.Println(err)
		return api.BadVolumeID, err
	}

	err = f.Truncate(int64(spec.SizeBytes))
	if err != nil {
		log.Println(err)
		return api.BadVolumeID, err
	}

	bd := &buseDev{
		file: buseFile,
		f:    f}

	nbd := Create(bd, int64(spec.SizeBytes))
	bd.nbd = nbd

	log.Infof("Connecting to NBD...")
	dev, err := bd.nbd.Connect()
	if err != nil {
		log.Println(err)
		return api.BadVolumeID, err
	}

	log.Infof("FsTypeting %s with %v", dev, spec.FsType.SimpleString())
	// TODO(pedge): fs type simple string could result in "none"
	// TODO(pedge): make common function
	cmd := "/sbin/mkfs." + spec.FsType.SimpleString()
	o, err := exec.Command(cmd, dev).Output()
	if err != nil {
		log.Warnf("Failed to run command %v %v: %v", cmd, dev, o)
		return api.BadVolumeID, err
	}

	log.Infof("BUSE mapped NBD device %s (size=%v) to block file %s", dev, spec.SizeBytes, buseFile)

	v := &api.Volume{
		ID:         volumeID,
		Source:     source,
		Locator:    locator,
		Ctime:      time.Now(),
		Spec:       spec,
		LastScan:   time.Now(),
		Format:     spec.FsType,
		State:      api.VolumeAvailable,
		Status:     api.Up,
		DevicePath: dev,
	}

	d.buseDevices[dev] = bd

	err = d.CreateVol(v)
	if err != nil {
		return api.BadVolumeID, err
	}
	return v.ID, err
}

func (d *driver) Delete(volumeID string) error {
	v, err := d.GetVol(volumeID)
	if err != nil {
		log.Println(err)
		return err
	}

	bd, ok := d.buseDevices[v.DevicePath]
	if !ok {
		err = fmt.Errorf("Cannot locate a BUSE device for %s", v.DevicePath)
		log.Println(err)
		return err
	}

	// Clean up buse block file and close the NBD connection.
	os.Remove(bd.file)
	bd.f.Close()
	bd.nbd.Disconnect()

	log.Infof("BUSE deleted volume %v at NBD device %s", volumeID, v.DevicePath)

	err = d.DeleteVol(volumeID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (d *driver) Mount(volumeID string, mountpath string) error {
	v, err := d.GetVol(volumeID)
	if err != nil {
		return fmt.Errorf("Failed to locate volume %q", volumeID)
	}
	// TODO(pedge): fs type simple string could result in "none"
	err = syscall.Mount(v.DevicePath, mountpath, v.Spec.FsType.SimpleString(), 0, "")
	if err != nil {
		log.Errorf("Mounting %s on %s failed because of %v", v.DevicePath, mountpath, err)
		return fmt.Errorf("Failed to mount %v at %v: %v", v.DevicePath, mountpath, err)
	}

	log.Infof("BUSE mounted NBD device %s at %s", v.DevicePath, mountpath)

	v.AttachPath = mountpath
	err = d.UpdateVol(v)

	return nil
}

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
	return err
}

func (d *driver) Snapshot(volumeID string, readonly bool, locator *openstorage.VolumeLocator) (string, error) {
	volIDs := make([]string, 1)
	volIDs[0] = volumeID
	vols, err := d.Inspect(volIDs)
	if err != nil {
		return api.BadVolumeID, nil
	}

	source := &openstorage.VolumeSource{ParentVolumeId: volumeID}
	newVolumeID, err := d.Create(locator, source, vols[0].Spec)
	if err != nil {
		return api.BadVolumeID, nil
	}

	// BUSE does not support snapshots, so just copy the block files.
	err = copyFile(BuseMountPath+volumeID, BuseMountPath+newVolumeID)
	if err != nil {
		d.Delete(newVolumeID)
		return api.BadVolumeID, nil
	}

	return newVolumeID, nil
}

func (d *driver) Set(volumeID string, locator *openstorage.VolumeLocator, spec *openstorage.VolumeSpec) error {
	if spec != nil {
		return volume.ErrNotSupported
	}
	v, err := d.GetVol(volumeID)
	if err != nil {
		return err
	}
	if locator != nil {
		v.Locator = locator
	}
	err = d.UpdateVol(v)
	return err
}

func (d *driver) Attach(volumeID string) (string, error) {
	// Nothing to do on attach.
	return path.Join(BuseMountPath, volumeID), nil
}

func (d *driver) Detach(volumeID string) error {
	// Nothing to do on detach.
	return nil
}

func (d *driver) Stats(volumeID string) (api.Stats, error) {
	return api.Stats{}, volume.ErrNotSupported
}

func (d *driver) Alerts(volumeID string) (api.Alerts, error) {
	return api.Alerts{}, volume.ErrNotSupported
}

func (d *driver) Shutdown() {
	log.Printf("%s Shutting down", Name)
	syscall.Unmount(BuseMountPath, 0)
}

func (d *driver) ClusterInit(self *api.Node, db *cluster.Database) error {
	return nil
}

func (d *driver) Init(self *api.Node, db *cluster.Database) error {
	return nil
}

func (d *driver) Join(self *api.Node, db *cluster.Database) error {
	return nil
}

func (d *driver) Add(self *api.Node) error {
	return nil
}

func (d *driver) Remove(self *api.Node) error {
	return nil
}

func (d *driver) Update(self *api.Node) error {
	return nil
}

func (d *driver) Leave(self *api.Node) error {
	return nil
}

func init() {

	// Register ourselves as an openstorage volume driver.
	volume.Register(Name, Init)
}
