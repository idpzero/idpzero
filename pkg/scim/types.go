// Package scim implements a minimal SCIM 2.0 provisioning client. idpzero acts
// as the SCIM client and pushes its configured users to an external SCIM service
// (the service provider). Only the User resource is supported.
package scim

// SCIM schema URNs (RFC 7643 / 7644).
const (
	UserSchema  = "urn:ietf:params:scim:schemas:core:2.0:User"
	ListSchema  = "urn:ietf:params:scim:api:messages:2.0:ListResponse"
	ErrorSchema = "urn:ietf:params:scim:api:messages:2.0:Error"

	// ContentType is the media type SCIM uses for request and response bodies.
	ContentType = "application/scim+json"
)

// User is a SCIM 2.0 core User resource. Only the subset of attributes idpzero
// can populate from its configured users is modelled here.
type User struct {
	Schemas      []string      `json:"schemas"`
	ID           string        `json:"id,omitempty"`         // assigned by the service provider
	ExternalID   string        `json:"externalId,omitempty"` // idpzero's stable subject
	UserName     string        `json:"userName"`
	DisplayName  string        `json:"displayName,omitempty"`
	Name         *Name         `json:"name,omitempty"`
	Emails       []Email       `json:"emails,omitempty"`
	PhoneNumbers []PhoneNumber `json:"phoneNumbers,omitempty"`
	Active       bool          `json:"active"`
}

// Name is the SCIM complex "name" attribute.
type Name struct {
	Formatted  string `json:"formatted,omitempty"`
	GivenName  string `json:"givenName,omitempty"`
	MiddleName string `json:"middleName,omitempty"`
	FamilyName string `json:"familyName,omitempty"`
}

// Email is a single SCIM multi-valued "emails" entry.
type Email struct {
	Value   string `json:"value"`
	Type    string `json:"type,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

// PhoneNumber is a single SCIM multi-valued "phoneNumbers" entry.
type PhoneNumber struct {
	Value   string `json:"value"`
	Type    string `json:"type,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

// ListResponse is the SCIM envelope returned by query/list operations.
type ListResponse struct {
	Schemas      []string `json:"schemas"`
	TotalResults int      `json:"totalResults"`
	Resources    []User   `json:"Resources"`
}

// ErrorResponse is the SCIM error envelope returned on non-2xx responses.
type ErrorResponse struct {
	Schemas []string `json:"schemas"`
	Detail  string   `json:"detail"`
	Status  string   `json:"status"`
}
