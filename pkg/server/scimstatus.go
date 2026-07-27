package server

import (
	"sync"

	"github.com/idpzero/idpzero/pkg/scim"
)

// scimStatus holds the result of the most recent SCIM sync run. It is in-memory
// only and resets on restart, mirroring the users/clients stores.
type scimStatus struct {
	lock sync.Mutex
	last *scim.Report
}

func newSCIMStatus() *scimStatus {
	return &scimStatus{}
}

func (s *scimStatus) Set(r *scim.Report) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.last = r
}

func (s *scimStatus) Get() *scim.Report {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.last
}
