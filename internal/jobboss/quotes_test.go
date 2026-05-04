package jobboss

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestQuotesGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/quote-line-items":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(ListResponse[QuoteLineItem]{Data: []QuoteLineItem{{QuoteNumber: "Q-1", ItemNumber: 1}}})
		case "/api/v1/quotes":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(ListResponse[Quote]{Data: []Quote{{QuoteNumber: "Q-1", Status: "OPEN"}}})
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(SingleResponse[Quote]{Data: Quote{QuoteNumber: "Q-2", Status: "OPEN"}})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/api/v1/quotes/Q-1":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(SingleResponse[Quote]{Data: Quote{QuoteNumber: "Q-1", Description: "Bracket"}})
			case http.MethodPatch:
				_ = json.NewEncoder(w).Encode(SingleResponse[Quote]{Data: Quote{QuoteNumber: "Q-1", Status: "WON"}})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/api/v1/quotes/Q-1/quote-line-item/1":
			_ = json.NewEncoder(w).Encode(SingleResponse[QuoteLineItem]{Data: QuoteLineItem{QuoteNumber: "Q-1", ItemNumber: 1, Description: "Laser"}})
		case "/api/v1/quotes/Q-1/quote-line-items":
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[QuoteLineItem]{Data: QuoteLineItem{QuoteNumber: "Q-1", ItemNumber: 2, Description: "Form"}})
		case "/api/v1/quotes/Q-1/quote-line-items/1":
			if r.Method != http.MethodPatch {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[QuoteLineItem]{Data: QuoteLineItem{QuoteNumber: "Q-1", ItemNumber: 1, Description: "Laser Updated"}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
	q := url.Values{}
	q.Set("fields", "quoteNumber,status")

	lineItems, err := c.GetQuoteLineItems(context.Background(), q)
	if err != nil || len(lineItems.Data) != 1 || lineItems.Data[0].ItemNumber != 1 {
		t.Fatalf("GetQuoteLineItems failed: resp=%+v err=%v", lineItems, err)
	}
	quotes, err := c.GetQuotes(context.Background(), q)
	if err != nil || len(quotes.Data) != 1 || quotes.Data[0].QuoteNumber != "Q-1" {
		t.Fatalf("GetQuotes failed: resp=%+v err=%v", quotes, err)
	}
	quote, err := c.GetQuoteByNumber(context.Background(), "Q-1", nil)
	if err != nil || quote.Data.Description != "Bracket" {
		t.Fatalf("GetQuoteByNumber failed: resp=%+v err=%v", quote, err)
	}
	_, err = c.GetQuoteByNumber(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "quote number is required") {
		t.Fatalf("expected quote number validation error, got %v", err)
	}

	createdQuote, err := c.CreateQuote(context.Background(), Quote{QuoteNumber: "Q-2"})
	if err != nil || createdQuote.Data.QuoteNumber != "Q-2" {
		t.Fatalf("CreateQuote failed: resp=%+v err=%v", createdQuote, err)
	}
	updatedQuote, err := c.UpdateQuote(context.Background(), "Q-1", Quote{Status: "WON"})
	if err != nil || updatedQuote.Data.Status != "WON" {
		t.Fatalf("UpdateQuote failed: resp=%+v err=%v", updatedQuote, err)
	}
	_, err = c.UpdateQuote(context.Background(), "", Quote{})
	if err == nil || !strings.Contains(err.Error(), "quote number is required") {
		t.Fatalf("expected quote number validation error, got %v", err)
	}

	lineItem, err := c.GetQuoteLineItem(context.Background(), "Q-1", 1, nil)
	if err != nil || lineItem.Data.Description != "Laser" {
		t.Fatalf("GetQuoteLineItem failed: resp=%+v err=%v", lineItem, err)
	}
	_, err = c.GetQuoteLineItem(context.Background(), "", 1, nil)
	if err == nil || !strings.Contains(err.Error(), "quote number is required") {
		t.Fatalf("expected quote number validation error, got %v", err)
	}
	_, err = c.GetQuoteLineItem(context.Background(), "Q-1", 0, nil)
	if err == nil || !strings.Contains(err.Error(), "item number must be greater than zero") {
		t.Fatalf("expected item number validation error, got %v", err)
	}

	createdLine, err := c.CreateQuoteLineItem(context.Background(), "Q-1", QuoteLineItem{ItemNumber: 2})
	if err != nil || createdLine.Data.ItemNumber != 2 {
		t.Fatalf("CreateQuoteLineItem failed: resp=%+v err=%v", createdLine, err)
	}
	_, err = c.CreateQuoteLineItem(context.Background(), "", QuoteLineItem{})
	if err == nil || !strings.Contains(err.Error(), "quote number is required") {
		t.Fatalf("expected quote number validation error, got %v", err)
	}

	updatedLine, err := c.UpdateQuoteLineItem(context.Background(), "Q-1", 1, QuoteLineItem{Description: "Laser Updated"})
	if err != nil || updatedLine.Data.Description != "Laser Updated" {
		t.Fatalf("UpdateQuoteLineItem failed: resp=%+v err=%v", updatedLine, err)
	}
	_, err = c.UpdateQuoteLineItem(context.Background(), "", 1, QuoteLineItem{})
	if err == nil || !strings.Contains(err.Error(), "quote number is required") {
		t.Fatalf("expected quote number validation error, got %v", err)
	}
	_, err = c.UpdateQuoteLineItem(context.Background(), "Q-1", 0, QuoteLineItem{})
	if err == nil || !strings.Contains(err.Error(), "item number must be greater than zero") {
		t.Fatalf("expected item number validation error, got %v", err)
	}
}
