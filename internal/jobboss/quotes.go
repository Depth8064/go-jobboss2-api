package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Quote models a quote payload.
type Quote struct {
	UniqueID        int     `json:"uniqueID,omitempty"`
	QuoteNumber     string  `json:"quoteNumber,omitempty"`
	CustomerCode    string  `json:"customerCode,omitempty"`
	Status          string  `json:"status,omitempty"`
	Description     string  `json:"description,omitempty"`
	EstimatedAmount float64 `json:"estimatedAmount,omitempty"`
	ExpireDate      string  `json:"expireDate,omitempty"`
	LastModDate     string  `json:"lastModDate,omitempty"`
	LastModUser     string  `json:"lastModUser,omitempty"`
}

// QuoteLineItem models a quote line item payload.
type QuoteLineItem struct {
	UniqueID    int     `json:"uniqueID,omitempty"`
	QuoteNumber string  `json:"quoteNumber,omitempty"`
	ItemNumber  int     `json:"itemNumber,omitempty"`
	PartNumber  string  `json:"partNumber,omitempty"`
	Description string  `json:"description,omitempty"`
	Quantity    float64 `json:"quantity,omitempty"`
	UnitPrice   float64 `json:"unitPrice,omitempty"`
	LastModDate string  `json:"lastModDate,omitempty"`
	LastModUser string  `json:"lastModUser,omitempty"`
}

// GetQuoteLineItems fetches quote line items.
func (c *Client) GetQuoteLineItems(ctx context.Context, query url.Values) (ListResponse[QuoteLineItem], error) {
	return GetJSON[ListResponse[QuoteLineItem]](ctx, c, "/api/v1/quote-line-items", query)
}

// GetQuotes fetches quotes.
func (c *Client) GetQuotes(ctx context.Context, query url.Values) (ListResponse[Quote], error) {
	return GetJSON[ListResponse[Quote]](ctx, c, "/api/v1/quotes", query)
}

// GetQuoteByNumber fetches a quote by number.
func (c *Client) GetQuoteByNumber(ctx context.Context, quoteNumber string, query url.Values) (SingleResponse[Quote], error) {
	quoteNumber = strings.TrimSpace(quoteNumber)
	if quoteNumber == "" {
		return SingleResponse[Quote]{}, fmt.Errorf("quote number is required")
	}
	endpoint := "/api/v1/quotes/" + url.PathEscape(quoteNumber)
	return GetJSON[SingleResponse[Quote]](ctx, c, endpoint, query)
}

// CreateQuote creates a quote.
func (c *Client) CreateQuote(ctx context.Context, quote Quote) (SingleResponse[Quote], error) {
	return PostJSON[Quote, SingleResponse[Quote]](ctx, c, "/api/v1/quotes", quote)
}

// UpdateQuote updates a quote by number.
func (c *Client) UpdateQuote(ctx context.Context, quoteNumber string, quote Quote) (UpdateResult, error) {
	quoteNumber = strings.TrimSpace(quoteNumber)
	if quoteNumber == "" {
		return UpdateResult{}, fmt.Errorf("quote number is required")
	}
	endpoint := "/api/v1/quotes/" + url.PathEscape(quoteNumber)
	return PatchNoContentJSON[Quote](ctx, c, endpoint, quote)
}

// GetQuoteLineItem fetches a quote line item by quote and item number.
func (c *Client) GetQuoteLineItem(ctx context.Context, quoteNumber string, itemNumber int, query url.Values) (SingleResponse[QuoteLineItem], error) {
	quoteNumber = strings.TrimSpace(quoteNumber)
	if quoteNumber == "" {
		return SingleResponse[QuoteLineItem]{}, fmt.Errorf("quote number is required")
	}
	if itemNumber <= 0 {
		return SingleResponse[QuoteLineItem]{}, fmt.Errorf("item number must be greater than zero")
	}
	endpoint := "/api/v1/quotes/" + url.PathEscape(quoteNumber) + "/quote-line-item/" + strconv.Itoa(itemNumber)
	return GetJSON[SingleResponse[QuoteLineItem]](ctx, c, endpoint, query)
}

// CreateQuoteLineItem creates a quote line item under a quote.
func (c *Client) CreateQuoteLineItem(ctx context.Context, quoteNumber string, lineItem QuoteLineItem) (SingleResponse[QuoteLineItem], error) {
	quoteNumber = strings.TrimSpace(quoteNumber)
	if quoteNumber == "" {
		return SingleResponse[QuoteLineItem]{}, fmt.Errorf("quote number is required")
	}
	endpoint := "/api/v1/quotes/" + url.PathEscape(quoteNumber) + "/quote-line-items"
	return PostJSON[QuoteLineItem, SingleResponse[QuoteLineItem]](ctx, c, endpoint, lineItem)
}

// UpdateQuoteLineItem updates a quote line item under a quote.
func (c *Client) UpdateQuoteLineItem(ctx context.Context, quoteNumber string, itemNumber int, lineItem QuoteLineItem) (UpdateResult, error) {
	quoteNumber = strings.TrimSpace(quoteNumber)
	if quoteNumber == "" {
		return UpdateResult{}, fmt.Errorf("quote number is required")
	}
	if itemNumber <= 0 {
		return UpdateResult{}, fmt.Errorf("item number must be greater than zero")
	}
	endpoint := "/api/v1/quotes/" + url.PathEscape(quoteNumber) + "/quote-line-items/" + strconv.Itoa(itemNumber)
	return PatchNoContentJSON[QuoteLineItem](ctx, c, endpoint, lineItem)
}
