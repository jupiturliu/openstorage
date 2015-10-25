package openstoragedocker

// NewLocalAPIClient constructs an APIClient that directly calls the given APIServer.
func NewLocalAPIClient(apiServer APIServer) APIClient {
	return newLocalAPIClient(apiServer)
}
