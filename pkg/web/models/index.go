package models

import "github.com/idpzero/idpzero/pkg/configuration"

type UrlInfo struct {
	Description string
	Url         string
}

type IndexModel struct {
	Urls    []UrlInfo
	Clients []configuration.ClientConfig
}
