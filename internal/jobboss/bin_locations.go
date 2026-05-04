package jobboss

import (
	"context"
	"net/url"
)

// BinLocation models a bin location payload.
type BinLocation struct {
	UniqueID             int     `json:"uniqueID,omitempty"`
	BinLocation          string  `json:"binLocation,omitempty"`
	PartNumber           string  `json:"partNumber,omitempty"`
	LotNumber            string  `json:"lotNumber,omitempty"`
	VendorCode           string  `json:"vendorCode,omitempty"`
	ReceiverNumber       string  `json:"receiverNumber,omitempty"`
	DeliveryTicketNumber string  `json:"deliveryTicketNumber,omitempty"`
	QuantityOnHand       float64 `json:"quantityOnHand,omitempty"`
	Cost                 float64 `json:"cost,omitempty"`
	POItemNumber         int     `json:"POItemNumber,omitempty"`
	DatePosted           string  `json:"datePosted,omitempty"`
	LastModDate          string  `json:"lastModDate,omitempty"`
	LastModUser          string  `json:"lastModUser,omitempty"`
}

// GetBinLocations fetches bin locations from /api/v1/bin-locations.
func (c *Client) GetBinLocations(ctx context.Context, query url.Values) (ListResponse[BinLocation], error) {
	return GetJSON[ListResponse[BinLocation]](ctx, c, "/api/v1/bin-locations", query)
}
