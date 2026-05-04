package jobboss

import (
	"context"
	"net/url"
)

// CurrencyCode models a currency code payload.
type CurrencyCode struct {
	CurrencyCode              string  `json:"currencyCode,omitempty"`
	CurrencyDescription       string  `json:"currencyDescription,omitempty"`
	ExchangeRate              float64 `json:"exchangeRate,omitempty"`
	LastModDate               string  `json:"lastModDate,omitempty"`
	LastModUser               string  `json:"lastModUser,omitempty"`
	CurrencyCountryCode       string  `json:"currencyCountryCode,omitempty"`
	CurrencyFormatMask        string  `json:"currencyFormatMask,omitempty"`
	CurrencyFormatDescription string  `json:"currencyFormatDescription,omitempty"`
}

// GetCurrencyCodes fetches currency codes from /api/v1/currency-codes.
func (c *Client) GetCurrencyCodes(ctx context.Context, query url.Values) (ListResponse[CurrencyCode], error) {
	return GetJSON[ListResponse[CurrencyCode]](ctx, c, "/api/v1/currency-codes", query)
}
