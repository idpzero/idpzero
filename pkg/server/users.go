package server

import (
	"sync"

	"github.com/idpzero/idpzero/pkg/configuration"
)

type users struct {
	lock  sync.Mutex
	users map[string]configuration.User
}

func newUsers() *users {
	return &users{
		lock:  sync.Mutex{},
		users: make(map[string]configuration.User),
	}
}

func (u *users) GetByID(id string) (configuration.User, bool) {
	u.lock.Lock()
	defer u.lock.Unlock()

	user, ok := u.users[id]
	return user, ok
}

func (u *users) Update(users []configuration.User) {
	u.lock.Lock()
	defer u.lock.Unlock()

	// clear existing
	u.users = make(map[string]configuration.User)
	for _, user := range users {
		u.users[user.Subject] = user
	}
}
