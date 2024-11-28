package server

import (
	"sync"

	"github.com/idpzero/idpzero/pkg/configuration"
)

type clients struct {
	lock    sync.Mutex
	clients map[string]configuration.ClientConfig
}

func newClients() *clients {
	return &clients{
		lock:    sync.Mutex{},
		clients: make(map[string]configuration.ClientConfig),
	}
}

func (u *clients) GetByID(id string) (configuration.ClientConfig, bool) {
	u.lock.Lock()
	defer u.lock.Unlock()

	client, ok := u.clients[id]
	return client, ok
}

func (u *clients) Update(clients []*configuration.ClientConfig) {
	u.lock.Lock()
	defer u.lock.Unlock()

	// clear existing
	u.clients = make(map[string]configuration.ClientConfig)
	for _, client := range clients {
		u.clients[client.ClientID] = *client
	}
}
