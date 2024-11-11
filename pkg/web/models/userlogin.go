package models

type UserLoginModel struct {
	Error         string
	AuthRequestID string
	Users         []OptionGroup
}

type OptionGroup struct {
	DisplayName string
	Options     []Option
}

type Option struct {
	ID          string
	DisplayName string
}
