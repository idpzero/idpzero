package models

import "github.com/idpzero/idpzero/pkg/configuration"

type UsersModel struct {
	Users []*configuration.User
}
