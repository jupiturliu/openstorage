package api

import (
	"fmt"
	"strconv"
	"strings"
)

// NewLocalVolumeAPIClient constructs an VolumeAPIClient that directly calls the given VolumeAPIServer.
func NewLocalVolumeAPIClient(apiServer VolumeAPIServer) VolumeAPIClient {
	return newLocalVolumeAPIClient(apiServer)
}

func FSTypeSimpleValueOf(s string) (FSType, error) {
	fsTypeObj, ok := FSType_value[fmt.Sprintf("FS_TYPE_%s", strings.ToUpper(s))]
	if !ok {
		return FSType_FS_TYPE_NONE, fmt.Errorf("no openstorage.FSType for %s", s)
	}
	return FSType(fsTypeObj), nil
}

func (x FSType) SimpleString() string {
	s, ok := FSType_name[int32(x)]
	if !ok {
		return strconv.Itoa(int(x))
	}
	return strings.TrimPrefix(strings.ToLower(s), "fs_type_")
}

func COSSimpleValueOf(s string) (COS, error) {
	cosObj, ok := COS_value[fmt.Sprintf("COS_%s", strings.ToUpper(s))]
	if !ok {
		return COS_COS_NONE, fmt.Errorf("no openstorage.COS for %s", s)
	}
	return COS(cosObj), nil
}

func (x COS) SimpleString() string {
	s, ok := COS_name[int32(x)]
	if !ok {
		return strconv.Itoa(int(x))
	}
	return strings.TrimPrefix(strings.ToLower(s), "cos_")
}
