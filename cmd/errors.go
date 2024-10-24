package cmd

import "errors"

var (
	ErrAlreadyInitialized = errors.New("configuration already initialized")
)
