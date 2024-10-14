package storage

import "errors"

var (
	ErrConfigNotFound    = errors.New("configuration not found")
	ErrConfigNotSupplied = errors.New("configuration file not provided")
)
