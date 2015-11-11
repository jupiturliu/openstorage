package dockerapi

// NewLocalVolumeAPIClient constructs an VolumeAPIClient that directly calls the given VolumeAPIServer.
func NewLocalVolumeAPIClient(volumeAPIServer VolumeAPIServer) VolumeAPIClient {
	return newLocalVolumeAPIClient(volumeAPIServer)
}
