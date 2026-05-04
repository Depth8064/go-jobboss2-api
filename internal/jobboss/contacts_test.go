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

func TestContactsGroup(t *testing.T) {
	phone := "555-2000"
	active := true

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/contacts":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(ListResponse[Contact]{
					Data: []Contact{{Object: "CUST", ContactCode: "ACME", Contact: "John", Email: "john@acme.com"}},
				})
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(SingleResponse[Contact]{
					Data: Contact{Object: "CUST", ContactCode: "ACME", Contact: "John", Email: "john@acme.com"},
				})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/api/v1/contacts/CUST/ACME/John":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(SingleResponse[Contact]{
					Data: Contact{Object: "CUST", ContactCode: "ACME", Contact: "John", Phone: "555-1000"},
				})
			case http.MethodPatch:
				_ = json.NewEncoder(w).Encode(SingleResponse[Contact]{
					Data: Contact{Object: "CUST", ContactCode: "ACME", Contact: "John", Phone: "555-2000", Active: true},
				})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
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

	listQ := url.Values{}
	listQ.Set("fields", "object,contactCode,contact,email")
	list, err := c.GetContacts(context.Background(), listQ)
	if err != nil || len(list.Data) != 1 {
		t.Fatalf("GetContacts failed: resp=%+v err=%v", list, err)
	}

	created, err := c.CreateContact(context.Background(), ContactCreate{Object: "CUST", ContactCode: "ACME", Contact: "John"})
	if err != nil || created.Data.Contact != "John" {
		t.Fatalf("CreateContact failed: resp=%+v err=%v", created, err)
	}

	one, err := c.GetContact(context.Background(), "CUST", "ACME", "John", nil)
	if err != nil || one.Data.Phone != "555-1000" {
		t.Fatalf("GetContact failed: resp=%+v err=%v", one, err)
	}

	upd, err := c.UpdateContact(context.Background(), "CUST", "ACME", "John", ContactUpdate{Phone: &phone, Active: &active})
	if err != nil || upd.StatusCode != http.StatusOK || upd.NoContent {
		t.Fatalf("UpdateContact failed: resp=%+v err=%v", upd, err)
	}

	_, err = c.GetContact(context.Background(), "", "ACME", "John", nil)
	if err == nil || !strings.Contains(err.Error(), "object is required") {
		t.Fatalf("expected object validation error, got %v", err)
	}

	_, err = c.GetContact(context.Background(), "CUST", "", "John", nil)
	if err == nil || !strings.Contains(err.Error(), "contact code is required") {
		t.Fatalf("expected contact code validation error, got %v", err)
	}

	_, err = c.GetContact(context.Background(), "CUST", "ACME", "", nil)
	if err == nil || !strings.Contains(err.Error(), "contact is required") {
		t.Fatalf("expected contact validation error, got %v", err)
	}

	_, err = c.UpdateContact(context.Background(), "", "ACME", "John", ContactUpdate{})
	if err == nil || !strings.Contains(err.Error(), "object is required") {
		t.Fatalf("expected object validation error, got %v", err)
	}

	_, err = c.UpdateContact(context.Background(), "CUST", "", "John", ContactUpdate{})
	if err == nil || !strings.Contains(err.Error(), "contact code is required") {
		t.Fatalf("expected contact code validation error, got %v", err)
	}

	_, err = c.UpdateContact(context.Background(), "CUST", "ACME", "", ContactUpdate{})
	if err == nil || !strings.Contains(err.Error(), "contact is required") {
		t.Fatalf("expected contact validation error, got %v", err)
	}
}
