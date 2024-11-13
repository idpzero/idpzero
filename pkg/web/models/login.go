package models

type UserLoginModel struct {
	Error         string
	AuthRequestID string
	Users         []UserOption
}

type UserOption struct {
	ID          string
	DisplayName string
}
