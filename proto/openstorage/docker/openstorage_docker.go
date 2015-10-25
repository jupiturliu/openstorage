// TODO(pedge): I don't like the package structure here, and underscores in package names
// is not recommended - revisit this

package openstorage_docker

// NewLocalAPIClient constructs an APIClient that directly calls the given APIServer.
func NewLocalAPIClient(apiServer APIServer) APIClient {
	return newLocalAPIClient(apiServer)
}
