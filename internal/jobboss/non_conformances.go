package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// NonConformance models a non-conformance payload.
type NonConformance struct {
	UniqueID         int     `json:"uniqueID,omitempty"`
	NCNumber         string  `json:"ncNumber,omitempty"`
	NCCode           string  `json:"ncCode,omitempty"`
	NCDescription    string  `json:"ncDescription,omitempty"`
	Status           string  `json:"status,omitempty"`
	OrderNumber      string  `json:"orderNumber,omitempty"`
	PartNumber       string  `json:"partNumber,omitempty"`
	PartDescription  string  `json:"partDescription,omitempty"`
	Quantity         float64 `json:"quantity,omitempty"`
	QuantityReturned float64 `json:"quantityReturned,omitempty"`
	ReasonCode       string  `json:"reasonCode,omitempty"`
	Disposition      string  `json:"disposition,omitempty"`
	CustomerCode     string  `json:"customerCode,omitempty"`
	VendorCode       string  `json:"vendorCode,omitempty"`
	EnterBy          string  `json:"enterBy,omitempty"`
	EnterDate        string  `json:"enterDate,omitempty"`
	LastModDate      string  `json:"lastModDate,omitempty"`
	LastModUser      string  `json:"lastModUser,omitempty"`
}

// GetNonConformances fetches non-conformance records.
func (c *Client) GetNonConformances(ctx context.Context, query url.Values) (ListResponse[NonConformance], error) {
	return GetJSON[ListResponse[NonConformance]](ctx, c, "/api/v1/non-conformances", query)
}

// GetNonConformanceByNumber fetches a non-conformance by NC number.
func (c *Client) GetNonConformanceByNumber(ctx context.Context, ncNumber string, query url.Values) (SingleResponse[NonConformance], error) {
	ncNumber = strings.TrimSpace(ncNumber)
	if ncNumber == "" {
		return SingleResponse[NonConformance]{}, fmt.Errorf("nc number is required")
	}
	endpoint := "/api/v1/non-conformances/" + url.PathEscape(ncNumber)
	return GetJSON[SingleResponse[NonConformance]](ctx, c, endpoint, query)
}
