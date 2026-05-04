package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// Customer models a JobBOSS customer payload.
type Customer struct {
	UniqueID       int    `json:"uniqueID,omitempty"`
	CustomerCode   string `json:"customerCode,omitempty"`
	CustomerName   string `json:"customerName,omitempty"`
	Active         bool   `json:"active,omitempty"`
	Phone          string `json:"phone,omitempty"`
	CurrencyCode   string `json:"currencyCode,omitempty"`
	SalesID        string `json:"salesID,omitempty"`
	TermsCode      string `json:"termsCode,omitempty"`
	BillingAddress string `json:"billingAddress1,omitempty"`
	BillingCity    string `json:"billingCity,omitempty"`
	BillingState   string `json:"billingState,omitempty"`
	BillingZIPCode string `json:"billingZIPCode,omitempty"`
	BillingCountry string `json:"billingCountry,omitempty"`
}

// CustomerCreate is a minimal create payload for /api/v1/customers.
type CustomerCreate struct {
	CustomerCode string `json:"customerCode"`
	CustomerName string `json:"customerName"`
	Active       *bool  `json:"active,omitempty"`
	Phone        string `json:"phone,omitempty"`
	CurrencyCode string `json:"currencyCode,omitempty"`
	SalesID      string `json:"salesID,omitempty"`
	TermsCode    string `json:"termsCode,omitempty"`
}

// CustomerUpdate is a partial update payload for /api/v1/customers/{customerCode}.
type CustomerUpdate struct {
	CustomerName *string `json:"customerName,omitempty"`
	Active       *bool   `json:"active,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	CurrencyCode *string `json:"currencyCode,omitempty"`
	SalesID      *string `json:"salesID,omitempty"`
	TermsCode    *string `json:"termsCode,omitempty"`
}

// GetCustomers fetches customers from /api/v1/customers.
func (c *Client) GetCustomers(ctx context.Context, query url.Values) (ListResponse[Customer], error) {
	return GetJSON[ListResponse[Customer]](ctx, c, "/api/v1/customers", query)
}

// GetCustomerByCode fetches a single customer by code.
func (c *Client) GetCustomerByCode(ctx context.Context, customerCode string, query url.Values) (SingleResponse[Customer], error) {
	customerCode = strings.TrimSpace(customerCode)
	if customerCode == "" {
		return SingleResponse[Customer]{}, fmt.Errorf("customer code is required")
	}
	endpoint := "/api/v1/customers/" + url.PathEscape(customerCode)
	return GetJSON[SingleResponse[Customer]](ctx, c, endpoint, query)
}

// CreateCustomer creates a customer using /api/v1/customers.
func (c *Client) CreateCustomer(ctx context.Context, payload CustomerCreate) (SingleResponse[Customer], error) {
	return PostJSON[CustomerCreate, SingleResponse[Customer]](ctx, c, "/api/v1/customers", payload)
}

// UpdateCustomer updates a customer using PATCH /api/v1/customers/{customerCode}.
func (c *Client) UpdateCustomer(ctx context.Context, customerCode string, payload CustomerUpdate) (UpdateResult, error) {
	customerCode = strings.TrimSpace(customerCode)
	if customerCode == "" {
		return UpdateResult{}, fmt.Errorf("customer code is required")
	}
	endpoint := "/api/v1/customers/" + url.PathEscape(customerCode)
	return PatchNoContentJSON[CustomerUpdate](ctx, c, endpoint, payload)
}
