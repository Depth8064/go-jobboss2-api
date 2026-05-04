package jobboss

import (
	"context"
	"net/url"
)

// CustomerReturn models a customer return payload.
type CustomerReturn struct {
	UniqueID             int    `json:"uniqueID,omitempty"`
	CustomerReturnNumber string `json:"customerReturnNumber,omitempty"`
	CustomerCode         string `json:"customerCode,omitempty"`
	OrderNumber          string `json:"orderNumber,omitempty"`
	Status               string `json:"status,omitempty"`
	ReturnDate           string `json:"returnDate,omitempty"`
	ReasonCode           string `json:"reasonCode,omitempty"`
	LastModDate          string `json:"lastModDate,omitempty"`
	LastModUser          string `json:"lastModUser,omitempty"`
}

// CustomerReturnLineItem models a customer return line item payload.
type CustomerReturnLineItem struct {
	UniqueID             int     `json:"uniqueID,omitempty"`
	CustomerReturnNumber string  `json:"customerReturnNumber,omitempty"`
	ItemNumber           int     `json:"itemNumber,omitempty"`
	PartNumber           string  `json:"partNumber,omitempty"`
	Description          string  `json:"description,omitempty"`
	Quantity             float64 `json:"quantity,omitempty"`
	Disposition          string  `json:"disposition,omitempty"`
	LastModDate          string  `json:"lastModDate,omitempty"`
	LastModUser          string  `json:"lastModUser,omitempty"`
}

// CustomerReturnRelease models a customer return release payload.
type CustomerReturnRelease struct {
	UniqueID             int     `json:"uniqueID,omitempty"`
	CustomerReturnNumber string  `json:"customerReturnNumber,omitempty"`
	ItemNumber           int     `json:"itemNumber,omitempty"`
	ReleaseNumber        int     `json:"releaseNumber,omitempty"`
	Quantity             float64 `json:"quantity,omitempty"`
	DueDate              string  `json:"dueDate,omitempty"`
	LastModDate          string  `json:"lastModDate,omitempty"`
	LastModUser          string  `json:"lastModUser,omitempty"`
}

// GetCustomerReturns fetches customer returns.
func (c *Client) GetCustomerReturns(ctx context.Context, query url.Values) (ListResponse[CustomerReturn], error) {
	return GetJSON[ListResponse[CustomerReturn]](ctx, c, "/api/v1/customer-returns", query)
}

// GetCustomerReturnLineItems fetches customer return line items.
func (c *Client) GetCustomerReturnLineItems(ctx context.Context, query url.Values) (ListResponse[CustomerReturnLineItem], error) {
	return GetJSON[ListResponse[CustomerReturnLineItem]](ctx, c, "/api/v1/customer-return-line-items", query)
}

// GetCustomerReturnReleases fetches customer return releases.
func (c *Client) GetCustomerReturnReleases(ctx context.Context, query url.Values) (ListResponse[CustomerReturnRelease], error) {
	return GetJSON[ListResponse[CustomerReturnRelease]](ctx, c, "/api/v1/customer-return-releases", query)
}
