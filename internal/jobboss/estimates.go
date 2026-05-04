package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// Estimate models an estimate payload.
type Estimate struct {
	UniqueID         int     `json:"uniqueID,omitempty"`
	PartNumber       string  `json:"partNumber,omitempty"`
	Revision         string  `json:"revision,omitempty"`
	Description      string  `json:"description,omitempty"`
	ProductCode      string  `json:"productCode,omitempty"`
	CustomerCode     string  `json:"customerCode,omitempty"`
	VendorCode1      string  `json:"vendorCode1,omitempty"`
	GLCode           string  `json:"GLCode,omitempty"`
	PurchasingGLCode string  `json:"purchasingGLCode,omitempty"`
	PricingUnit      string  `json:"pricingUnit,omitempty"`
	PurchaseFactor   float64 `json:"purchaseFactor,omitempty"`
	StockingCost     float64 `json:"stockingCost,omitempty"`
	LeadTime         int     `json:"leadTime,omitempty"`
	Active           bool    `json:"active,omitempty"`
	LastModDate      string  `json:"lastModDate,omitempty"`
	LastModUser      string  `json:"lastModUser,omitempty"`
}

// EstimateCreate is a minimal payload for creating estimates.
type EstimateCreate struct {
	PartNumber       string  `json:"partNumber,omitempty"`
	Revision         string  `json:"revision,omitempty"`
	Description      string  `json:"description,omitempty"`
	ProductCode      string  `json:"productCode,omitempty"`
	CustomerCode     string  `json:"customerCode,omitempty"`
	VendorCode1      string  `json:"vendorCode1,omitempty"`
	GLCode           string  `json:"GLCode,omitempty"`
	PurchasingGLCode string  `json:"purchasingGLCode,omitempty"`
	PricingUnit      string  `json:"pricingUnit,omitempty"`
	PurchaseFactor   float64 `json:"purchaseFactor,omitempty"`
	StockingCost     float64 `json:"stockingCost,omitempty"`
	LeadTime         int     `json:"leadTime,omitempty"`
	Active           bool    `json:"active,omitempty"`
}

// EstimateUpdate is a payload for updating estimates.
type EstimateUpdate struct {
	Revision    *string `json:"revision,omitempty"`
	Description *string `json:"description,omitempty"`
	ProductCode *string `json:"productCode,omitempty"`
	GLCode      *string `json:"GLCode,omitempty"`
}

// GetEstimates fetches estimates.
func (c *Client) GetEstimates(ctx context.Context, query url.Values) (ListResponse[Estimate], error) {
	return GetJSON[ListResponse[Estimate]](ctx, c, "/api/v1/estimates", query)
}

// GetEstimateByPartNumber fetches an estimate by part number.
func (c *Client) GetEstimateByPartNumber(ctx context.Context, partNumber string, query url.Values) (SingleResponse[Estimate], error) {
	partNumber = strings.TrimSpace(partNumber)
	if partNumber == "" {
		return SingleResponse[Estimate]{}, fmt.Errorf("part number is required")
	}
	endpoint := "/api/v1/estimates/" + url.PathEscape(partNumber)
	return GetJSON[SingleResponse[Estimate]](ctx, c, endpoint, query)
}

// CreateEstimate creates an estimate.
func (c *Client) CreateEstimate(ctx context.Context, payload EstimateCreate) (SingleResponse[Estimate], error) {
	return PostJSON[EstimateCreate, SingleResponse[Estimate]](ctx, c, "/api/v1/estimates", payload)
}

// UpdateEstimate updates an estimate by part number.
func (c *Client) UpdateEstimate(ctx context.Context, partNumber string, payload EstimateUpdate) (SingleResponse[Estimate], error) {
	partNumber = strings.TrimSpace(partNumber)
	if partNumber == "" {
		return SingleResponse[Estimate]{}, fmt.Errorf("part number is required")
	}
	endpoint := "/api/v1/estimates/" + url.PathEscape(partNumber)
	return PatchJSON[EstimateUpdate, SingleResponse[Estimate]](ctx, c, endpoint, payload)
}
