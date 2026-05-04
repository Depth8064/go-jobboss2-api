package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ReasonCode models a reason code payload.
type ReasonCode struct {
	UniqueID    int    `json:"uniqueID,omitempty"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`
	Active      bool   `json:"active,omitempty"`
	LastModDate string `json:"lastModDate,omitempty"`
	LastModUser string `json:"lastModUser,omitempty"`
}

// ReportRequest models a report request payload.
type ReportRequest struct {
	ReportName   string                 `json:"reportName,omitempty"`
	OutputFormat string                 `json:"outputFormat,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
}

// ReportStatus models a report status payload.
type ReportStatus struct {
	RequestID   string `json:"requestId,omitempty"`
	Status      string `json:"status,omitempty"`
	DownloadURL string `json:"downloadUrl,omitempty"`
	Message     string `json:"message,omitempty"`
}

// Salesperson models a salesperson payload.
type Salesperson struct {
	UniqueID    int    `json:"uniqueID,omitempty"`
	SalesID     string `json:"salesID,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Active      bool   `json:"active,omitempty"`
	LastModDate string `json:"lastModDate,omitempty"`
	LastModUser string `json:"lastModUser,omitempty"`
}

// ShippingAddress models a shipping address payload.
type ShippingAddress struct {
	UniqueID     int    `json:"uniqueID,omitempty"`
	CustomerCode string `json:"customerCode,omitempty"`
	Location     string `json:"location,omitempty"`
	Name         string `json:"name,omitempty"`
	Address1     string `json:"address1,omitempty"`
	Address2     string `json:"address2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
	Country      string `json:"country,omitempty"`
	Active       bool   `json:"active,omitempty"`
	LastModDate  string `json:"lastModDate,omitempty"`
	LastModUser  string `json:"lastModUser,omitempty"`
}

// TaxCode models a tax code payload.
type TaxCode struct {
	UniqueID    int     `json:"uniqueID,omitempty"`
	TaxCode     string  `json:"taxCode,omitempty"`
	Description string  `json:"description,omitempty"`
	TaxRate     float64 `json:"taxRate,omitempty"`
	LastModDate string  `json:"lastModDate,omitempty"`
	LastModUser string  `json:"lastModUser,omitempty"`
}

// Term models a terms payload.
type Term struct {
	UniqueID    int    `json:"uniqueID,omitempty"`
	TermsCode   string `json:"termsCode,omitempty"`
	Description string `json:"description,omitempty"`
	DaysDue     int    `json:"daysDue,omitempty"`
	LastModDate string `json:"lastModDate,omitempty"`
	LastModUser string `json:"lastModUser,omitempty"`
}

// GetReasonCodes fetches reason codes.
func (c *Client) GetReasonCodes(ctx context.Context, query url.Values) (ListResponse[ReasonCode], error) {
	return GetJSON[ListResponse[ReasonCode]](ctx, c, "/api/v1/reason-codes", query)
}

// GetReasonCodeByID fetches a reason code by unique ID.
func (c *Client) GetReasonCodeByID(ctx context.Context, uniqueID int, query url.Values) (SingleResponse[ReasonCode], error) {
	if uniqueID <= 0 {
		return SingleResponse[ReasonCode]{}, fmt.Errorf("unique id must be greater than zero")
	}
	endpoint := "/api/v1/reason-codes/" + strconv.Itoa(uniqueID)
	return GetJSON[SingleResponse[ReasonCode]](ctx, c, endpoint, query)
}

// GetReleases fetches releases.
func (c *Client) GetReleases(ctx context.Context, query url.Values) (ListResponse[Release], error) {
	return GetJSON[ListResponse[Release]](ctx, c, "/api/v1/releases", query)
}

// CreateReport requests report generation.
func (c *Client) CreateReport(ctx context.Context, req ReportRequest) (SingleResponse[ReportStatus], error) {
	return PostJSON[ReportRequest, SingleResponse[ReportStatus]](ctx, c, "/api/v1/reports", req)
}

// GetReportStatus fetches report status by request ID.
func (c *Client) GetReportStatus(ctx context.Context, requestID string, query url.Values) (SingleResponse[ReportStatus], error) {
	requestID = strings.TrimSpace(requestID)
	if requestID == "" {
		return SingleResponse[ReportStatus]{}, fmt.Errorf("request id is required")
	}
	endpoint := "/api/v1/reports/" + url.PathEscape(requestID)
	return GetJSON[SingleResponse[ReportStatus]](ctx, c, endpoint, query)
}

// GetSalespersons fetches salespersons.
func (c *Client) GetSalespersons(ctx context.Context, query url.Values) (ListResponse[Salesperson], error) {
	return GetJSON[ListResponse[Salesperson]](ctx, c, "/api/v1/salespersons", query)
}

// UpdateSalesperson updates a salesperson by sales ID.
func (c *Client) UpdateSalesperson(ctx context.Context, salesID string, salesperson Salesperson) (UpdateResult, error) {
	salesID = strings.TrimSpace(salesID)
	if salesID == "" {
		return UpdateResult{}, fmt.Errorf("sales id is required")
	}
	endpoint := "/api/v1/salespersons/" + url.PathEscape(salesID)
	return PatchNoContentJSON[Salesperson](ctx, c, endpoint, salesperson)
}

// GetShippingAddresses fetches shipping addresses.
func (c *Client) GetShippingAddresses(ctx context.Context, query url.Values) (ListResponse[ShippingAddress], error) {
	return GetJSON[ListResponse[ShippingAddress]](ctx, c, "/api/v1/shipping-addresses", query)
}

// GetShippingAddress fetches a shipping address by customer/location.
func (c *Client) GetShippingAddress(ctx context.Context, customerCode, location string, query url.Values) (SingleResponse[ShippingAddress], error) {
	customerCode = strings.TrimSpace(customerCode)
	location = strings.TrimSpace(location)
	if customerCode == "" {
		return SingleResponse[ShippingAddress]{}, fmt.Errorf("customer code is required")
	}
	if location == "" {
		return SingleResponse[ShippingAddress]{}, fmt.Errorf("location is required")
	}
	endpoint := "/api/v1/shipping-addresses/" + url.PathEscape(customerCode) + "/" + url.PathEscape(location)
	return GetJSON[SingleResponse[ShippingAddress]](ctx, c, endpoint, query)
}

// CreateShippingAddress creates a shipping address.
func (c *Client) CreateShippingAddress(ctx context.Context, shippingAddress ShippingAddress) (SingleResponse[ShippingAddress], error) {
	return PostJSON[ShippingAddress, SingleResponse[ShippingAddress]](ctx, c, "/api/v1/shipping-addresses", shippingAddress)
}

// UpdateShippingAddress updates a shipping address by customer/location.
func (c *Client) UpdateShippingAddress(ctx context.Context, customerCode, location string, shippingAddress ShippingAddress) (UpdateResult, error) {
	customerCode = strings.TrimSpace(customerCode)
	location = strings.TrimSpace(location)
	if customerCode == "" {
		return UpdateResult{}, fmt.Errorf("customer code is required")
	}
	if location == "" {
		return UpdateResult{}, fmt.Errorf("location is required")
	}
	endpoint := "/api/v1/shipping-addresses/" + url.PathEscape(customerCode) + "/" + url.PathEscape(location)
	return PatchNoContentJSON[ShippingAddress](ctx, c, endpoint, shippingAddress)
}

// GetTaxCodes fetches tax codes.
func (c *Client) GetTaxCodes(ctx context.Context, query url.Values) (ListResponse[TaxCode], error) {
	return GetJSON[ListResponse[TaxCode]](ctx, c, "/api/v1/tax-codes", query)
}

// GetTaxCode fetches a tax code by code.
func (c *Client) GetTaxCode(ctx context.Context, taxCode string, query url.Values) (SingleResponse[TaxCode], error) {
	taxCode = strings.TrimSpace(taxCode)
	if taxCode == "" {
		return SingleResponse[TaxCode]{}, fmt.Errorf("tax code is required")
	}
	endpoint := "/api/v1/tax-codes/" + url.PathEscape(taxCode)
	return GetJSON[SingleResponse[TaxCode]](ctx, c, endpoint, query)
}

// GetTerms fetches terms.
func (c *Client) GetTerms(ctx context.Context, query url.Values) (ListResponse[Term], error) {
	return GetJSON[ListResponse[Term]](ctx, c, "/api/v1/terms", query)
}

// GetTerm fetches terms by terms code.
func (c *Client) GetTerm(ctx context.Context, termsCode string, query url.Values) (SingleResponse[Term], error) {
	termsCode = strings.TrimSpace(termsCode)
	if termsCode == "" {
		return SingleResponse[Term]{}, fmt.Errorf("terms code is required")
	}
	endpoint := "/api/v1/terms/" + url.PathEscape(termsCode)
	return GetJSON[SingleResponse[Term]](ctx, c, endpoint, query)
}
