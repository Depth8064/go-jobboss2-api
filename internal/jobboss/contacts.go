package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// Contact models a contact payload.
type Contact struct {
	UniqueID    int    `json:"uniqueID,omitempty"`
	Object      string `json:"object,omitempty"`
	ContactCode string `json:"contactCode,omitempty"`
	Contact     string `json:"contact,omitempty"`
	Title       string `json:"title,omitempty"`
	Phone       string `json:"phone,omitempty"`
	CellPhone   string `json:"cellPhone,omitempty"`
	Extension   string `json:"extension,omitempty"`
	Email       string `json:"email,omitempty"`
	MobileEmail string `json:"mobileEmail,omitempty"`
	Comments    string `json:"comments,omitempty"`
	Active      bool   `json:"active,omitempty"`
	IsOptedOut  bool   `json:"isOptedOut,omitempty"`
	LastModDate string `json:"lastModDate,omitempty"`
	RevisedDate string `json:"revisedDate,omitempty"`
}

// ContactCreate is a payload for creating contacts.
type ContactCreate struct {
	Object      string `json:"object,omitempty"`
	ContactCode string `json:"contactCode,omitempty"`
	Contact     string `json:"contact,omitempty"`
	Title       string `json:"title,omitempty"`
	Phone       string `json:"phone,omitempty"`
	CellPhone   string `json:"cellPhone,omitempty"`
	Extension   string `json:"extension,omitempty"`
	Email       string `json:"email,omitempty"`
	MobileEmail string `json:"mobileEmail,omitempty"`
	Comments    string `json:"comments,omitempty"`
	Active      bool   `json:"active,omitempty"`
	IsOptedOut  bool   `json:"isOptedOut,omitempty"`
}

// ContactUpdate is a payload for patching contacts.
type ContactUpdate struct {
	Title     *string `json:"title,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	CellPhone *string `json:"cellPhone,omitempty"`
	Extension *string `json:"extension,omitempty"`
	Email     *string `json:"email,omitempty"`
	Comments  *string `json:"comments,omitempty"`
	Active    *bool   `json:"active,omitempty"`
}

// GetContacts fetches contacts from /api/v1/contacts.
func (c *Client) GetContacts(ctx context.Context, query url.Values) (ListResponse[Contact], error) {
	return GetJSON[ListResponse[Contact]](ctx, c, "/api/v1/contacts", query)
}

// CreateContact creates a contact.
func (c *Client) CreateContact(ctx context.Context, payload ContactCreate) (SingleResponse[Contact], error) {
	return PostJSON[ContactCreate, SingleResponse[Contact]](ctx, c, "/api/v1/contacts", payload)
}

// GetContact fetches a contact by object, contact code, and contact.
func (c *Client) GetContact(ctx context.Context, object, contactCode, contact string, query url.Values) (SingleResponse[Contact], error) {
	object = strings.TrimSpace(object)
	if object == "" {
		return SingleResponse[Contact]{}, fmt.Errorf("object is required")
	}
	contactCode = strings.TrimSpace(contactCode)
	if contactCode == "" {
		return SingleResponse[Contact]{}, fmt.Errorf("contact code is required")
	}
	contact = strings.TrimSpace(contact)
	if contact == "" {
		return SingleResponse[Contact]{}, fmt.Errorf("contact is required")
	}
	endpoint := "/api/v1/contacts/" + url.PathEscape(object) + "/" + url.PathEscape(contactCode) + "/" + url.PathEscape(contact)
	return GetJSON[SingleResponse[Contact]](ctx, c, endpoint, query)
}

// UpdateContact updates a contact by object, contact code, and contact.
func (c *Client) UpdateContact(ctx context.Context, object, contactCode, contact string, payload ContactUpdate) (SingleResponse[Contact], error) {
	object = strings.TrimSpace(object)
	if object == "" {
		return SingleResponse[Contact]{}, fmt.Errorf("object is required")
	}
	contactCode = strings.TrimSpace(contactCode)
	if contactCode == "" {
		return SingleResponse[Contact]{}, fmt.Errorf("contact code is required")
	}
	contact = strings.TrimSpace(contact)
	if contact == "" {
		return SingleResponse[Contact]{}, fmt.Errorf("contact is required")
	}
	endpoint := "/api/v1/contacts/" + url.PathEscape(object) + "/" + url.PathEscape(contactCode) + "/" + url.PathEscape(contact)
	return PatchJSON[ContactUpdate, SingleResponse[Contact]](ctx, c, endpoint, payload)
}
