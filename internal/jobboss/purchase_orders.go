package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// PackingList models a packing list payload.
type PackingList struct {
	UniqueID       int    `json:"uniqueID,omitempty"`
	PackingListNo  string `json:"packingListNo,omitempty"`
	OrderNumber    string `json:"orderNumber,omitempty"`
	CustomerCode   string `json:"customerCode,omitempty"`
	ShipDate       string `json:"shipDate,omitempty"`
	TrackingNumber string `json:"trackingNumber,omitempty"`
	LastModDate    string `json:"lastModDate,omitempty"`
	LastModUser    string `json:"lastModUser,omitempty"`
}

// PackingListLineItem models a packing-list-line-item payload.
type PackingListLineItem struct {
	UniqueID        int     `json:"uniqueID,omitempty"`
	PackingListNo   string  `json:"packingListNo,omitempty"`
	OrderNumber     string  `json:"orderNumber,omitempty"`
	PartNumber      string  `json:"partNumber,omitempty"`
	ItemNumber      int     `json:"itemNumber,omitempty"`
	QuantityShipped float64 `json:"quantityShipped,omitempty"`
	LastModDate     string  `json:"lastModDate,omitempty"`
	LastModUser     string  `json:"lastModUser,omitempty"`
}

// PurchaseOrder models a purchase order payload.
type PurchaseOrder struct {
	UniqueID     int     `json:"uniqueID,omitempty"`
	PONumber     string  `json:"poNumber,omitempty"`
	VendorCode   string  `json:"vendorCode,omitempty"`
	OrderDate    string  `json:"orderDate,omitempty"`
	ExpectedDate string  `json:"expectedDate,omitempty"`
	Status       string  `json:"status,omitempty"`
	Buyer        string  `json:"buyer,omitempty"`
	TotalAmount  float64 `json:"totalAmount,omitempty"`
	LastModDate  string  `json:"lastModDate,omitempty"`
	LastModUser  string  `json:"lastModUser,omitempty"`
}

// PurchaseOrderLineItem models a purchase-order-line-item payload.
type PurchaseOrderLineItem struct {
	UniqueID        int     `json:"uniqueID,omitempty"`
	PONumber        string  `json:"poNumber,omitempty"`
	PartNumber      string  `json:"partNumber,omitempty"`
	ItemNumber      int     `json:"itemNumber,omitempty"`
	Description     string  `json:"description,omitempty"`
	QuantityOrdered float64 `json:"quantityOrdered,omitempty"`
	UnitCost        float64 `json:"unitCost,omitempty"`
	LastModDate     string  `json:"lastModDate,omitempty"`
	LastModUser     string  `json:"lastModUser,omitempty"`
}

// PurchaseOrderRelease models a purchase-order-release payload.
type PurchaseOrderRelease struct {
	UniqueID         int     `json:"uniqueID,omitempty"`
	PONumber         string  `json:"poNumber,omitempty"`
	PartNumber       string  `json:"partNumber,omitempty"`
	ItemNumber       int     `json:"itemNumber,omitempty"`
	ReleaseNumber    int     `json:"releaseNumber,omitempty"`
	ReleaseDate      string  `json:"releaseDate,omitempty"`
	QuantityReleased float64 `json:"quantityReleased,omitempty"`
	LastModDate      string  `json:"lastModDate,omitempty"`
	LastModUser      string  `json:"lastModUser,omitempty"`
}

// GetPackingLists fetches packing lists.
func (c *Client) GetPackingLists(ctx context.Context, query url.Values) (ListResponse[PackingList], error) {
	return GetJSON[ListResponse[PackingList]](ctx, c, "/api/v1/packing-lists", query)
}

// GetPackingListLineItems fetches packing list line items.
func (c *Client) GetPackingListLineItems(ctx context.Context, query url.Values) (ListResponse[PackingListLineItem], error) {
	return GetJSON[ListResponse[PackingListLineItem]](ctx, c, "/api/v1/packing-list-line-items", query)
}

// GetPurchaseOrders fetches purchase orders.
func (c *Client) GetPurchaseOrders(ctx context.Context, query url.Values) (ListResponse[PurchaseOrder], error) {
	return GetJSON[ListResponse[PurchaseOrder]](ctx, c, "/api/v1/purchase-orders", query)
}

// GetPurchaseOrderByNumber fetches a purchase order by number.
func (c *Client) GetPurchaseOrderByNumber(ctx context.Context, poNumber string, query url.Values) (SingleResponse[PurchaseOrder], error) {
	poNumber = strings.TrimSpace(poNumber)
	if poNumber == "" {
		return SingleResponse[PurchaseOrder]{}, fmt.Errorf("po number is required")
	}
	endpoint := "/api/v1/purchase-orders/" + url.PathEscape(poNumber)
	return GetJSON[SingleResponse[PurchaseOrder]](ctx, c, endpoint, query)
}

// GetPurchaseOrderLineItems fetches purchase order line items.
func (c *Client) GetPurchaseOrderLineItems(ctx context.Context, query url.Values) (ListResponse[PurchaseOrderLineItem], error) {
	return GetJSON[ListResponse[PurchaseOrderLineItem]](ctx, c, "/api/v1/purchase-order-line-items", query)
}

// GetPurchaseOrderLineItem fetches a purchase order line item by composite key.
func (c *Client) GetPurchaseOrderLineItem(ctx context.Context, poNumber, partNumber, itemNumber string, query url.Values) (SingleResponse[PurchaseOrderLineItem], error) {
	poNumber = strings.TrimSpace(poNumber)
	partNumber = strings.TrimSpace(partNumber)
	itemNumber = strings.TrimSpace(itemNumber)
	if poNumber == "" {
		return SingleResponse[PurchaseOrderLineItem]{}, fmt.Errorf("po number is required")
	}
	if partNumber == "" {
		return SingleResponse[PurchaseOrderLineItem]{}, fmt.Errorf("part number is required")
	}
	if itemNumber == "" {
		return SingleResponse[PurchaseOrderLineItem]{}, fmt.Errorf("item number is required")
	}
	endpoint := "/api/v1/purchase-order-line-items/" + url.PathEscape(poNumber) + "/" + url.PathEscape(partNumber) + "/" + url.PathEscape(itemNumber)
	return GetJSON[SingleResponse[PurchaseOrderLineItem]](ctx, c, endpoint, query)
}

// GetPurchaseOrderReleases fetches purchase order releases.
func (c *Client) GetPurchaseOrderReleases(ctx context.Context, query url.Values) (ListResponse[PurchaseOrderRelease], error) {
	return GetJSON[ListResponse[PurchaseOrderRelease]](ctx, c, "/api/v1/purchase-order-releases", query)
}
