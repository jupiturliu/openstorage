package api

import (
	"time"

	"github.com/libopenstorage/openstorage/proto/openstorage"
)

// BadVolumeID invalid volume ID, usually accompanied by an error.
const BadVolumeID = ""

// VolumeStatus a health status.
type VolumeStatus string

const (
	// NotPresent This volume is not present.
	NotPresent = VolumeStatus("NotPresent")
	// Up status healthy
	Up = VolumeStatus("Up")
	// Down status failure.
	Down = VolumeStatus("Down")
	// Degraded status up but with degraded performance. In a RAID group, this may indicate a problem with one or more drives
	Degraded = VolumeStatus("Degraded")
)

// VolumeState is one of the below enumerations and reflects the state
// of a volume.
type VolumeState int

const (
	// VolumePending volume is transitioning to new state
	VolumePending VolumeState = 1 << iota
	// VolumeAvailable volume is ready to be assigned to a container
	VolumeAvailable
	// VolumeAttached is attached to container
	VolumeAttached
	// VolumeDetached is detached but associated with a container.
	VolumeDetached
	// VolumeDetaching is detach is in progress.
	VolumeDetaching
	// VolumeError is in Error State
	VolumeError
	// VolumeDeleted is deleted, it will remain in this state while resources are
	// asynchronously reclaimed.
	VolumeDeleted
)

// VolumeStateAny a filter that selects all volumes
const VolumeStateAny = VolumePending | VolumeAvailable | VolumeAttached | VolumeDetaching | VolumeDetached | VolumeError | VolumeDeleted

// MachineID is a node instance identifier for clustered systems.
type MachineID string

const MachineNone MachineID = ""

// Volume represents a live, created volume.
type Volume struct {
	// ID Self referential VolumeID
	ID string
	// Source
	Source *openstorage.VolumeSource
	// Readonly
	Readonly bool
	// Locator User specified locator
	Locator *openstorage.VolumeLocator
	// Ctime Volume creation time
	Ctime time.Time
	// Spec User specified VolumeSpec
	Spec *openstorage.VolumeSpec
	// Usage Volume usage
	Usage uint64
	// LastScan time when an integrity check for run
	LastScan time.Time
	// Format Filesystem type if any
	Format openstorage.FSType
	// Status see VolumeStatus
	Status VolumeStatus
	// State see VolumeState
	State VolumeState
	// AttachedOn - Node on which this volume is attached.
	AttachedOn MachineID
	// DevicePath
	DevicePath string
	// AttachPath
	AttachPath string
	// ReplicaSet Set of nodes no which this Volume is erasure coded - for clustered storage arrays
	ReplicaSet []MachineID
	// Error Last recorded error
	Error string
}

// Alerts
type Stats struct {
	// Reads completed successfully.
	Reads int64
	// ReadMs time spent in reads in ms.
	ReadMs int64
	// ReadBytes
	ReadBytes int64
	// Writes completed successfully.
	Writes int64
	// WriteBytes
	WriteBytes int64
	// WriteMs time spent in writes in ms.
	WriteMs int64
	// IOProgress I/Os curently in progress.
	IOProgress int64
	// IOMs time spent doing I/Os ms.
	IOMs int64
}

// Alerts
type Alerts struct {
}
