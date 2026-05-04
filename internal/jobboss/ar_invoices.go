package jobboss

import (
	"context"
	"net/url"
)

// ARInvoiceDetail models an AR invoice detail payload.
type ARInvoiceDetail struct {
	UniqueID             int     `json:"uniqueID,omitempty"`
	InvoiceNumber        string  `json:"invoiceNumber,omitempty"`
	DeliveryTicketNumber string  `json:"deliveryTicketNumber,omitempty"`
	PartNumber           string  `json:"partNumber,omitempty"`
	PartDescription      string  `json:"partDescription,omitempty"`
	JobNumber            string  `json:"jobNumber,omitempty"`
	PONumber             string  `json:"PONumber,omitempty"`
	ProductCode          string  `json:"productCode,omitempty"`
	WorkCode             string  `json:"workCode,omitempty"`
	Unit                 string  `json:"unit,omitempty"`
	UnitPrice            float64 `json:"unitPrice,omitempty"`
	UnitPriceForeign     float64 `json:"unitPriceForeign,omitempty"`
	QuantityOrdered      float64 `json:"quantityOrdered,omitempty"`
	QuantityShipped      float64 `json:"quantityShipped,omitempty"`
	LineTotal            float64 `json:"lineTotal,omitempty"`
	LineTotalForeign     float64 `json:"lineTotalForeign,omitempty"`
	TaxStatus            bool    `json:"taxStatus,omitempty"`
	LastModDate          string  `json:"lastModDate,omitempty"`
}

// ARInvoice models an AR invoice payload.
type ARInvoice struct {
	UniqueID               int     `json:"uniqueID,omitempty"`
	InvoiceNumber          string  `json:"invoiceNumber,omitempty"`
	CustomerCode           string  `json:"customerCode,omitempty"`
	CustomerDescription    string  `json:"customerDescription,omitempty"`
	CurrencyCode           string  `json:"currencyCode,omitempty"`
	PaymentStatus          string  `json:"paymentStatus,omitempty"`
	InvoiceDate            string  `json:"invoiceDate,omitempty"`
	NetDueDate             string  `json:"netDueDate,omitempty"`
	PaymentDate            string  `json:"paymentDate,omitempty"`
	InvoiceTotal           float64 `json:"invoiceTotal,omitempty"`
	InvoiceTotalForeign    float64 `json:"invoiceTotalForeign,omitempty"`
	OpenInvoiceAmount      float64 `json:"openInvoiceAmount,omitempty"`
	AmountPaidSoFar        float64 `json:"amountPaidSoFar,omitempty"`
	AmountPaidSoFarForeign float64 `json:"amountPaidSoFarForeign,omitempty"`
	SalesID                string  `json:"salesID,omitempty"`
	TaxCode                string  `json:"taxCode,omitempty"`
	TermsCode              string  `json:"termsCode,omitempty"`
	Territory              string  `json:"territory,omitempty"`
	Location               string  `json:"location,omitempty"`
	LastModDate            string  `json:"lastModDate,omitempty"`
}

// GetARInvoiceDetails fetches AR invoice details from /api/v1/ar-invoice-details.
func (c *Client) GetARInvoiceDetails(ctx context.Context, query url.Values) (ListResponse[ARInvoiceDetail], error) {
	return GetJSON[ListResponse[ARInvoiceDetail]](ctx, c, "/api/v1/ar-invoice-details", query)
}

// GetARInvoices fetches AR invoices from /api/v1/ar-invoices.
func (c *Client) GetARInvoices(ctx context.Context, query url.Values) (ListResponse[ARInvoice], error) {
	return GetJSON[ListResponse[ARInvoice]](ctx, c, "/api/v1/ar-invoices", query)
}
