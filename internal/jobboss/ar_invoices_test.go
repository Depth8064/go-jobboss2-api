package jobboss

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestARInvoicesSection(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/ar-invoice-details":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			if r.URL.Query().Get("fields") == "" {
				t.Fatalf("expected fields query")
			}
			_ = json.NewEncoder(w).Encode(ListResponse[ARInvoiceDetail]{
				Data: []ARInvoiceDetail{{InvoiceNumber: "INV-100", PartNumber: "PN-1", LineTotal: 100.5}},
			})
		case "/api/v1/ar-invoices":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			if r.URL.Query().Get("sort") == "" {
				t.Fatalf("expected sort query")
			}
			_ = json.NewEncoder(w).Encode(ListResponse[ARInvoice]{
				Data: []ARInvoice{{InvoiceNumber: "INV-100", CustomerCode: "ACME", PaymentStatus: "UNPAID", InvoiceTotal: 100.5}},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{
		BaseURL:           ts.URL,
		TokenURL:          ts.URL + "/oauth2/api-user/token",
		ClientID:          "id",
		ClientSecret:      "secret",
		HTTPClient:        ts.Client(),
		MaxAuthRetryOn401: 0,
	})

	detailsQuery := url.Values{}
	detailsQuery.Set("fields", "invoiceNumber,partNumber,lineTotal")
	detailsResp, err := c.GetARInvoiceDetails(context.Background(), detailsQuery)
	if err != nil {
		t.Fatalf("GetARInvoiceDetails failed: %v", err)
	}
	if len(detailsResp.Data) != 1 || detailsResp.Data[0].InvoiceNumber != "INV-100" {
		t.Fatalf("unexpected AR invoice details response: %+v", detailsResp)
	}

	invoiceQuery := url.Values{}
	invoiceQuery.Set("sort", "-invoiceDate")
	invResp, err := c.GetARInvoices(context.Background(), invoiceQuery)
	if err != nil {
		t.Fatalf("GetARInvoices failed: %v", err)
	}
	if len(invResp.Data) != 1 || invResp.Data[0].InvoiceNumber != "INV-100" {
		t.Fatalf("unexpected AR invoices response: %+v", invResp)
	}
}
