package models

import (
	"fmt"
	"sort"

	"github.com/idpzero/idpzero/pkg/configuration"
)

// ClaimRow is a single displayable claim (name/value) for a user.
type ClaimRow struct {
	Name  string
	Value string
}

// UserView is a user prepared for display, with its claims flattened into an
// ordered list of name/value rows (including any custom claims).
type UserView struct {
	Subject      string
	LoginDisplay string
	Claims       []ClaimRow
}

type UsersModel struct {
	Users []UserView
}

// NewUsersModel builds the users view model from the configured users, flattening
// each user's claims into ordered rows so they can be rendered without the UI
// needing to know about every individual claim field.
func NewUsersModel(users []*configuration.User) UsersModel {
	model := UsersModel{Users: make([]UserView, 0, len(users))}

	for _, u := range users {
		if u == nil {
			continue
		}

		model.Users = append(model.Users, UserView{
			Subject:      u.Subject,
			LoginDisplay: u.LoginDisplay,
			Claims:       flattenClaims(u.Claims),
		})
	}

	return model
}

// flattenClaims turns the strongly typed (and mostly optional) claims into an
// ordered list of name/value rows, appending custom claims sorted by name. Only
// claims that are actually set are included.
func flattenClaims(c configuration.UserClaims) []ClaimRow {
	rows := make([]ClaimRow, 0)

	appendStr := func(name string, v *string) {
		if v != nil {
			rows = append(rows, ClaimRow{Name: name, Value: *v})
		}
	}
	appendBool := func(name string, v *bool) {
		if v != nil {
			rows = append(rows, ClaimRow{Name: name, Value: fmt.Sprintf("%t", *v)})
		}
	}

	appendStr("preferred_username", c.PreferredUsername)
	appendStr("name", c.Name)
	appendStr("given_name", c.GivenName)
	appendStr("middle_name", c.MiddleName)
	appendStr("family_name", c.FamilyName)
	appendStr("nickname", c.Nickname)
	appendStr("email", c.Email)
	appendBool("email_verified", c.EmailVerified)
	appendStr("phone", c.Phone)
	appendBool("phone_verified", c.PhoneVerified)
	if c.UpdatedAt != nil {
		rows = append(rows, ClaimRow{Name: "updated_at", Value: c.UpdatedAt.Format("2006-01-02 15:04:05 MST")})
	}

	// custom claims, sorted by name for stable rendering
	customNames := make([]string, 0, len(c.Custom))
	for name := range c.Custom {
		customNames = append(customNames, name)
	}
	sort.Strings(customNames)
	for _, name := range customNames {
		rows = append(rows, ClaimRow{Name: name, Value: fmt.Sprintf("%v", c.Custom[name])})
	}

	return rows
}
