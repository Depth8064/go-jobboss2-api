package jobboss

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing")

func TestCustomersSection(t *testing.T) {
	active := true
	name := "Acme Industries"
	phone := "555-0100"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/customers":
			switch r.Method {
			case http.MethodGet:
				if r.URL.Query().Get("fields") == "" {
					t.Fatalf("expected fields query")
				}
				_ = json.NewEncoder(w).Encode(ListResponse[Customer]{
					Data: []Customer{{CustomerCode: "ACME", CustomerName: "Acme", Active: true}},
				})
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(SingleResponse[Customer]{
					Data: Customer{CustomerCode: "ACME", CustomerName: "Acme", Active: true},
				})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/api/v1/customers/ACME":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(SingleResponse[Customer]{
					Data: Customer{CustomerCode: "ACME", CustomerName: "Acme", Active: true},
				})
			case http.MethodPatch:
				_ = json.NewEncoder(w).Encode(SingleResponse[Customer]{
					Data: Customer{CustomerCode: "ACME", CustomerName: "Acme Industries", Active: true, Phone: "555-0100"},
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

	q := url.Values{}
	q.Set("fields", "customerCode,customerName")
	listResp, err := c.GetCustomers(context.Background(), q)
	if err != nil {
		t.Fatalf("GetCustomers failed: %v", err)
	}
	if len(listResp.Data) != 1 || listResp.Data[0].CustomerCode != "ACME" {
		t.Fatalf("unexpected customer list response: %+v", listResp)
	}

	singleResp, err := c.GetCustomerByCode(context.Background(), "ACME", nil)
	if err != nil {
		t.Fatalf("GetCustomerByCode failed: %v", err)
	}
	if singleResp.Data.CustomerCode != "ACME" {
		t.Fatalf("unexpected customer response: %+v", singleResp)
	}

	_, err = c.GetCustomerByCode(context.Background(), "   ", nil)
	if err == nil || !strings.Contains(err.Error(), "customer code is required") {
		t.Fatalf("expected empty customer code error, got %v", err)
	}

	createResp, err := c.CreateCustomer(context.Background(), CustomerCreate{
		CustomerCode: "ACME",
		CustomerName: "Acme",
		Active:       &active,
	})
	if err != nil {
		t.Fatalf("CreateCustomer failed: %v", err)
	}
	if createResp.Data.CustomerCode != "ACME" {
		t.Fatalf("unexpected created customer: %+v", createResp)
	}

	updateResp, err := c.UpdateCustomer(context.Background(), "ACME", CustomerUpdate{
		CustomerName: &name,
		Phone:        &phone,
	})
	if err != nil {
		t.Fatalf("UpdateCustomer failed: %v", err)
	}
	if updateResp.Data.CustomerName != "Acme Industries" || updateResp.Data.Phone != "555-0100" {
		t.Fatalf("unexpected updated customer: %+v", updateResp)
	}

	_, err = c.UpdateCustomer(context.Background(), "", CustomerUpdate{})
	if err == nil || !strings.Contains(err.Error(), "customer code is required") {
		t.Fatalf("expected empty customer code update error, got %v", err)
	}
}


