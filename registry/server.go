package registry

import "sync"

const (
	ServerPort  = "3000"
	ServicesURL = "http://localhost:" + ServerPort + "/services"
)

type registry struct {
	registrations []Registration
	mutex         sync.Mutex
}

func (r *registry) add(registration Registration) error {
	r.mutex.Lock()

	r.mutex.Unlock()
	return nil
}
