package scim

import (
	"testing"

	"github.com/idpzero/idpzero/pkg/configuration"
)

func strptr(s string) *string { return &s }
func boolptr(b bool) *bool    { return &b }

func TestToSCIMUser_FullClaims(t *testing.T) {
	u := &configuration.User{
		Subject:      "john_verified_email",
		LoginDisplay: "John Jones",
		Claims: configuration.UserClaims{
			Name:          strptr("John Jones"),
			GivenName:     strptr("John"),
			FamilyName:    strptr("Jones"),
			Email:         strptr("john@example.com"),
			EmailVerified: boolptr(true),
		},
	}

	got := ToSCIMUser(u)

	if got.ExternalID != "john_verified_email" {
		t.Errorf("ExternalID = %q, want the subject", got.ExternalID)
	}
	// No preferred_username, so userName falls back to the subject.
	if got.UserName != "john_verified_email" {
		t.Errorf("UserName = %q, want subject fallback", got.UserName)
	}
	if got.DisplayName != "John Jones" {
		t.Errorf("DisplayName = %q, want login display", got.DisplayName)
	}
	if !got.Active {
		t.Error("Active = false, want true")
	}
	if got.Name == nil || got.Name.GivenName != "John" || got.Name.FamilyName != "Jones" {
		t.Errorf("Name = %+v, want given/family populated", got.Name)
	}
	if len(got.Emails) != 1 || got.Emails[0].Value != "john@example.com" || !got.Emails[0].Primary {
		t.Errorf("Emails = %+v, want one primary work email", got.Emails)
	}
	if len(got.Schemas) != 1 || got.Schemas[0] != UserSchema {
		t.Errorf("Schemas = %v, want [%s]", got.Schemas, UserSchema)
	}
}

func TestToSCIMUser_PreferredUsernameWins(t *testing.T) {
	u := &configuration.User{
		Subject: "abc123",
		Claims: configuration.UserClaims{
			PreferredUsername: strptr("jsmith"),
		},
	}

	got := ToSCIMUser(u)

	if got.UserName != "jsmith" {
		t.Errorf("UserName = %q, want preferred_username", got.UserName)
	}
	if got.ExternalID != "abc123" {
		t.Errorf("ExternalID = %q, want subject", got.ExternalID)
	}
}

func TestToSCIMUser_MinimalNilClaims(t *testing.T) {
	// A user with only a subject and no claims must not panic and must omit
	// the optional complex/multi-valued attributes.
	u := &configuration.User{Subject: "bare"}

	got := ToSCIMUser(u)

	if got.UserName != "bare" {
		t.Errorf("UserName = %q, want subject", got.UserName)
	}
	if got.Name != nil {
		t.Errorf("Name = %+v, want nil for empty claims", got.Name)
	}
	if got.Emails != nil {
		t.Errorf("Emails = %+v, want nil", got.Emails)
	}
	if got.PhoneNumbers != nil {
		t.Errorf("PhoneNumbers = %+v, want nil", got.PhoneNumbers)
	}
}
