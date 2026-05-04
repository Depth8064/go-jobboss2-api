package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// ProductCode models a product-code payload.
type ProductCode struct {
	UniqueID    int    `json:"uniqueID,omitempty"`
	ProductCode string `json:"productCode,omitempty"`
	Description string `json:"description,omitempty"`
	ProductType string `json:"productType,omitempty"`
	Active      bool   `json:"active,omitempty"`
	LastModDate string `json:"lastModDate,omitempty"`
	LastModUser string `json:"lastModUser,omitempty"`
}

// GetProductCodes fetches product codes.
func (c *Client) GetProductCodes(ctx context.Context, query url.Values) (ListResponse[ProductCode], error) {
	return GetJSON[ListResponse[ProductCode]](ctx, c, "/api/v1/product-codes", query)
}

// GetProductCodeByCode fetches a product code by code.
func (c *Client) GetProductCodeByCode(ctx context.Context, productCode string, query url.Values) (SingleResponse[ProductCode], error) {
	productCode = strings.TrimSpace(productCode)
	if productCode == "" {
		return SingleResponse[ProductCode]{}, fmt.Errorf("product code is required")
	}
	endpoint := "/api/v1/product-codes/" + url.PathEscape(productCode)
	return GetJSON[SingleResponse[ProductCode]](ctx, c, endpoint, query)
}
