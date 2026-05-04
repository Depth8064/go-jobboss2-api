package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// JobMaterial models a job material payload.
type JobMaterial struct {
	UniqueID     int     `json:"uniqueID,omitempty"`
	JobNumber    string  `json:"jobNumber,omitempty"`
	OrderNumber  string  `json:"orderNumber,omitempty"`
	PartNumber   string  `json:"partNumber,omitempty"`
	Description  string  `json:"description,omitempty"`
	ProductCode  string  `json:"productCode,omitempty"`
	VendorCode   string  `json:"vendorCode,omitempty"`
	PONumber     string  `json:"PONumber,omitempty"`
	POItemNumber int     `json:"POItemNumber,omitempty"`
	GLCode       string  `json:"GLCode,omitempty"`
	PricingUnit  string  `json:"pricingUnit,omitempty"`
	StockingCost float64 `json:"stockingCost,omitempty"`
	ResalePrice  float64 `json:"resalePrice,omitempty"`
	StepNumber   int     `json:"stepNumber,omitempty"`
	LastModDate  string  `json:"lastModDate,omitempty"`
	LastModUser  string  `json:"lastModUser,omitempty"`
}

// JobRequirement models a job requirement payload.
type JobRequirement struct {
	UniqueID              int     `json:"uniqueID,omitempty"`
	JobNumber             string  `json:"jobNumber,omitempty"`
	OrderNumber           string  `json:"orderNumber,omitempty"`
	PartNumber            string  `json:"partNumber,omitempty"`
	PartDescription       string  `json:"partDescription,omitempty"`
	ProductCode           string  `json:"productCode,omitempty"`
	VendorCode            string  `json:"vendorCode,omitempty"`
	PONumber              string  `json:"PONumber,omitempty"`
	POItemNumber          int     `json:"POItemNumber,omitempty"`
	GLCode                string  `json:"GLCode,omitempty"`
	PricingUnit           string  `json:"pricingUnit,omitempty"`
	PurchaseUnit          string  `json:"purchaseUnit,omitempty"`
	PurchaseQuantity      float64 `json:"purchaseQuantity,omitempty"`
	QuantityToBuy         float64 `json:"quantityToBuy,omitempty"`
	Price                 float64 `json:"price,omitempty"`
	Cost                  float64 `json:"cost,omitempty"`
	StepNumber            int     `json:"stepNumber,omitempty"`
	CertificationRequired bool    `json:"certificationRequired,omitempty"`
	OutsideService        bool    `json:"outsideService,omitempty"`
	LastModDate           string  `json:"lastModDate,omitempty"`
	LastModUser           string  `json:"lastModUser,omitempty"`
}

// Material models a material payload.
type Material struct {
	UniqueID         int     `json:"uniqueID,omitempty"`
	PartNumber       string  `json:"partNumber,omitempty"`
	SubPartNumber    string  `json:"subPartNumber,omitempty"`
	Description      string  `json:"description,omitempty"`
	Vendor           string  `json:"vendor,omitempty"`
	MaterialDetailID string  `json:"materialDetailID,omitempty"`
	ItemNumber       int     `json:"itemNumber,omitempty"`
	StepNumber       int     `json:"stepNumber,omitempty"`
	Unit             string  `json:"unit,omitempty"`
	Quantity         float64 `json:"quantity,omitempty"`
	UnitCost         float64 `json:"unitCost,omitempty"`
	UnitPrice        float64 `json:"unitPrice,omitempty"`
	TotalCost        float64 `json:"totalCost,omitempty"`
	TotalPrice       float64 `json:"totalPrice,omitempty"`
	LastModDate      string  `json:"lastModDate,omitempty"`
	LastModUser      string  `json:"lastModUser,omitempty"`
}

// GetJobMaterials fetches job materials.
func (c *Client) GetJobMaterials(ctx context.Context, query url.Values) (ListResponse[JobMaterial], error) {
	return GetJSON[ListResponse[JobMaterial]](ctx, c, "/api/v1/job-materials", query)
}

// GetJobMaterialByID fetches a job material by unique ID.
func (c *Client) GetJobMaterialByID(ctx context.Context, uniqueID string, query url.Values) (SingleResponse[JobMaterial], error) {
	uniqueID = strings.TrimSpace(uniqueID)
	if uniqueID == "" {
		return SingleResponse[JobMaterial]{}, fmt.Errorf("unique id is required")
	}
	endpoint := "/api/v1/job-materials/" + url.PathEscape(uniqueID)
	return GetJSON[SingleResponse[JobMaterial]](ctx, c, endpoint, query)
}

// GetJobRequirements fetches job requirements.
func (c *Client) GetJobRequirements(ctx context.Context, query url.Values) (ListResponse[JobRequirement], error) {
	return GetJSON[ListResponse[JobRequirement]](ctx, c, "/api/v1/job-requirements", query)
}

// GetJobRequirementByID fetches a job requirement by unique ID.
func (c *Client) GetJobRequirementByID(ctx context.Context, uniqueID string, query url.Values) (SingleResponse[JobRequirement], error) {
	uniqueID = strings.TrimSpace(uniqueID)
	if uniqueID == "" {
		return SingleResponse[JobRequirement]{}, fmt.Errorf("unique id is required")
	}
	endpoint := "/api/v1/job-requirements/" + url.PathEscape(uniqueID)
	return GetJSON[SingleResponse[JobRequirement]](ctx, c, endpoint, query)
}

// GetMaterials fetches materials.
func (c *Client) GetMaterials(ctx context.Context, query url.Values) (ListResponse[Material], error) {
	return GetJSON[ListResponse[Material]](ctx, c, "/api/v1/materials", query)
}

// GetEstimateMaterial fetches a material by estimate part number and sub part number.
func (c *Client) GetEstimateMaterial(ctx context.Context, partNumber, subPartNumber string, query url.Values) (SingleResponse[Material], error) {
	partNumber = strings.TrimSpace(partNumber)
	if partNumber == "" {
		return SingleResponse[Material]{}, fmt.Errorf("part number is required")
	}
	subPartNumber = strings.TrimSpace(subPartNumber)
	if subPartNumber == "" {
		return SingleResponse[Material]{}, fmt.Errorf("sub part number is required")
	}
	endpoint := "/api/v1/estimates/" + url.PathEscape(partNumber) + "/materials/" + url.PathEscape(subPartNumber)
	return GetJSON[SingleResponse[Material]](ctx, c, endpoint, query)
}
