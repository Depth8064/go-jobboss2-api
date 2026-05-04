package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// GLCode models a GL code payload.
type GLCode struct {
	UniqueID            int    `json:"uniqueID,omitempty"`
	GLAccountNumber     string `json:"GLAccountNumber,omitempty"`
	AccountingID        string `json:"accountingID,omitempty"`
	AccountingCode      string `json:"accountingCode,omitempty"`
	Description         string `json:"description,omitempty"`
	CurrencyCode        string `json:"currencyCode,omitempty"`
	Division            string `json:"division,omitempty"`
	VendorCode          string `json:"vendorCode,omitempty"`
	AccumulationAccount bool   `json:"accumulationAccount,omitempty"`
	Active              bool   `json:"active,omitempty"`
	LastModDate         string `json:"lastModDate,omitempty"`
	LastModUser         string `json:"lastModUser,omitempty"`
}

// GetGLCodes fetches GL codes.
func (c *Client) GetGLCodes(ctx context.Context, query url.Values) (ListResponse[GLCode], error) {
	return GetJSON[ListResponse[GLCode]](ctx, c, "/api/v1/gl-codes", query)
}

// GetGLCodeByAccountNumber fetches a GL code by account number.
func (c *Client) GetGLCodeByAccountNumber(ctx context.Context, glAccountNumber string, query url.Values) (SingleResponse[GLCode], error) {
	glAccountNumber = strings.TrimSpace(glAccountNumber)
	if glAccountNumber == "" {
		return SingleResponse[GLCode]{}, fmt.Errorf("gl account number is required")
	}
	endpoint := "/api/v1/gl-codes/" + url.PathEscape(glAccountNumber)
	return GetJSON[SingleResponse[GLCode]](ctx, c, endpoint, query)
}
