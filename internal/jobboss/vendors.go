package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// Vendor models a vendor payload.
type Vendor struct {
	UniqueID    int    `json:"uniqueID,omitempty"`
	VendorCode  string `json:"vendorCode,omitempty"`
	VendorName  string `json:"vendorName,omitempty"`
	Address1    string `json:"address1,omitempty"`
	Address2    string `json:"address2,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	PostalCode  string `json:"postalCode,omitempty"`
	Country     string `json:"country,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty"`
	Status      string `json:"status,omitempty"`
	LastModDate string `json:"lastModDate,omitempty"`
	LastModUser string `json:"lastModUser,omitempty"`
}

// VendorReturn models a vendor return payload.
type VendorReturn struct {
	UniqueID           int    `json:"uniqueID,omitempty"`
	VendorReturnNumber string `json:"vendorReturnNumber,omitempty"`
	VendorCode         string `json:"vendorCode,omitempty"`
	Status             string `json:"status,omitempty"`
	ReturnDate         string `json:"returnDate,omitempty"`
	ReasonCode         string `json:"reasonCode,omitempty"`
	LastModDate        string `json:"lastModDate,omitempty"`
	LastModUser        string `json:"lastModUser,omitempty"`
}

// VendorReturnLineItem models a vendor return line item payload.
type VendorReturnLineItem struct {
	UniqueID           int     `json:"uniqueID,omitempty"`
	VendorReturnNumber string  `json:"vendorReturnNumber,omitempty"`
	ItemNumber         int     `json:"itemNumber,omitempty"`
	PartNumber         string  `json:"partNumber,omitempty"`
	Description        string  `json:"description,omitempty"`
	Quantity           float64 `json:"quantity,omitempty"`
	LastModDate        string  `json:"lastModDate,omitempty"`
	LastModUser        string  `json:"lastModUser,omitempty"`
}

// VendorReturnRelease models a vendor return release payload.
type VendorReturnRelease struct {
	UniqueID           int     `json:"uniqueID,omitempty"`
	VendorReturnNumber string  `json:"vendorReturnNumber,omitempty"`
	ItemNumber         int     `json:"itemNumber,omitempty"`
	ReleaseNumber      int     `json:"releaseNumber,omitempty"`
	Quantity           float64 `json:"quantity,omitempty"`
	DueDate            string  `json:"dueDate,omitempty"`
	LastModDate        string  `json:"lastModDate,omitempty"`
	LastModUser        string  `json:"lastModUser,omitempty"`
}

// GetVendors fetches vendors.
func (c *Client) GetVendors(ctx context.Context, query url.Values) (ListResponse[Vendor], error) {
	return GetJSON[ListResponse[Vendor]](ctx, c, "/api/v1/vendors", query)
}

// GetVendorByCode fetches a vendor by vendor code.
func (c *Client) GetVendorByCode(ctx context.Context, vendorCode string, query url.Values) (SingleResponse[Vendor], error) {
	vendorCode = strings.TrimSpace(vendorCode)
	if vendorCode == "" {
		return SingleResponse[Vendor]{}, fmt.Errorf("vendor code is required")
	}
	endpoint := "/api/v1/vendors/" + url.PathEscape(vendorCode)
	return GetJSON[SingleResponse[Vendor]](ctx, c, endpoint, query)
}

// CreateVendor creates a vendor.
func (c *Client) CreateVendor(ctx context.Context, vendor Vendor) (SingleResponse[Vendor], error) {
	return PostJSON[Vendor, SingleResponse[Vendor]](ctx, c, "/api/v1/vendors", vendor)
}

// UpdateVendor updates a vendor by code.
func (c *Client) UpdateVendor(ctx context.Context, vendorCode string, vendor Vendor) (SingleResponse[Vendor], error) {
	vendorCode = strings.TrimSpace(vendorCode)
	if vendorCode == "" {
		return SingleResponse[Vendor]{}, fmt.Errorf("vendor code is required")
	}
	endpoint := "/api/v1/vendors/" + url.PathEscape(vendorCode)
	return PatchJSON[Vendor, SingleResponse[Vendor]](ctx, c, endpoint, vendor)
}

// GetVendorReturns fetches vendor returns.
func (c *Client) GetVendorReturns(ctx context.Context, query url.Values) (ListResponse[VendorReturn], error) {
	return GetJSON[ListResponse[VendorReturn]](ctx, c, "/api/v1/vendor-returns", query)
}

// GetVendorReturnLineItems fetches vendor return line items.
func (c *Client) GetVendorReturnLineItems(ctx context.Context, query url.Values) (ListResponse[VendorReturnLineItem], error) {
	return GetJSON[ListResponse[VendorReturnLineItem]](ctx, c, "/api/v1/vendor-returns-line-items", query)
}

// GetVendorReturnReleases fetches vendor return releases.
func (c *Client) GetVendorReturnReleases(ctx context.Context, query url.Values) (ListResponse[VendorReturnRelease], error) {
	return GetJSON[ListResponse[VendorReturnRelease]](ctx, c, "/api/v1/vendor-returns-releases", query)
}
