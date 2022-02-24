package registry

type Registration struct {
	ServiceName      ServiceName
	ServiceURL       string
	RequiredServices []ServiceName
	ServiceUpdateURL string
}

type ServiceName string

const (
	LogServiceName  = ServiceName("Log Service")
	UserServiceName = ServiceName("User Service")
)

type patchEntry struct {
	Name ServiceName
	URL  string
}

type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
