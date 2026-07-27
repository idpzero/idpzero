package scim

import "github.com/idpzero/idpzero/pkg/configuration"

// ToSCIMUser maps an idpzero configured user onto a SCIM 2.0 User resource.
//
// The idpzero Subject is used as the SCIM externalId so provisioning is
// idempotent: the target is queried by externalId to decide between create and
// replace. Claim fields are optional pointers on the source, so every access is
// nil-guarded.
func ToSCIMUser(u *configuration.User) User {
	c := u.Claims

	user := User{
		Schemas:    []string{UserSchema},
		ExternalID: u.Subject,
		UserName:   deref(c.PreferredUsername, u.Subject),
		Active:     true,
	}

	// Display name prefers the login display, then the name claim.
	if u.LoginDisplay != "" {
		user.DisplayName = u.LoginDisplay
	} else if c.Name != nil {
		user.DisplayName = *c.Name
	}

	name := Name{
		Formatted:  deref(c.Name, ""),
		GivenName:  deref(c.GivenName, ""),
		MiddleName: deref(c.MiddleName, ""),
		FamilyName: deref(c.FamilyName, ""),
	}
	if name != (Name{}) {
		user.Name = &name
	}

	if c.Email != nil && *c.Email != "" {
		user.Emails = []Email{{Value: *c.Email, Type: "work", Primary: true}}
	}

	if c.Phone != nil && *c.Phone != "" {
		user.PhoneNumbers = []PhoneNumber{{Value: *c.Phone, Type: "work", Primary: true}}
	}

	return user
}

// deref returns the pointed-to string, or fallback when the pointer is nil or empty.
func deref(v *string, fallback string) string {
	if v == nil || *v == "" {
		return fallback
	}
	return *v
}
