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

func TestSalesSupportGroups(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/reason-codes":
			_ = json.NewEncoder(w).Encode(ListResponse[ReasonCode]{Data: []ReasonCode{{UniqueID: 1, Code: "RC1"}}})
		case "/api/v1/reason-codes/1":
			_ = json.NewEncoder(w).Encode(SingleResponse[ReasonCode]{Data: ReasonCode{UniqueID: 1, Description: "Damaged"}})
		case "/api/v1/releases":
			_ = json.NewEncoder(w).Encode(ListResponse[Release]{Data: []Release{{OrderNumber: "SO-1", ItemNumber: 1}}})
		case "/api/v1/reports":
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[ReportStatus]{Data: ReportStatus{RequestID: "RPT-1", Status: "QUEUED"}})
		case "/api/v1/reports/RPT-1":
			_ = json.NewEncoder(w).Encode(SingleResponse[ReportStatus]{Data: ReportStatus{RequestID: "RPT-1", Status: "COMPLETED"}})
		case "/api/v1/salespersons":
			_ = json.NewEncoder(w).Encode(ListResponse[Salesperson]{Data: []Salesperson{{SalesID: "S1", FirstName: "A"}}})
		case "/api/v1/salespersons/S1":
			if r.Method != http.MethodPatch {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[Salesperson]{Data: Salesperson{SalesID: "S1", Active: true}})
		case "/api/v1/shipping-addresses":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(ListResponse[ShippingAddress]{Data: []ShippingAddress{{CustomerCode: "C1", Location: "MAIN"}}})
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(SingleResponse[ShippingAddress]{Data: ShippingAddress{CustomerCode: "C1", Location: "ALT"}})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/api/v1/shipping-addresses/C1/MAIN":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(SingleResponse[ShippingAddress]{Data: ShippingAddress{CustomerCode: "C1", Location: "MAIN", City: "Chicago"}})
			case http.MethodPatch:
				_ = json.NewEncoder(w).Encode(SingleResponse[ShippingAddress]{Data: ShippingAddress{CustomerCode: "C1", Location: "MAIN", City: "Peoria"}})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/api/v1/tax-codes":
			_ = json.NewEncoder(w).Encode(ListResponse[TaxCode]{Data: []TaxCode{{TaxCode: "IL"}}})
		case "/api/v1/tax-codes/IL":
			_ = json.NewEncoder(w).Encode(SingleResponse[TaxCode]{Data: TaxCode{TaxCode: "IL", TaxRate: 7.25}})
		case "/api/v1/terms":
			_ = json.NewEncoder(w).Encode(ListResponse[Term]{Data: []Term{{TermsCode: "N30"}}})
		case "/api/v1/terms/N30":
			_ = json.NewEncoder(w).Encode(SingleResponse[Term]{Data: Term{TermsCode: "N30", DaysDue: 30}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
	q := url.Values{}
	q.Set("fields", "code,description")

	reasons, err := c.GetReasonCodes(context.Background(), q)
	if err != nil || len(reasons.Data) != 1 || reasons.Data[0].Code != "RC1" {
		t.Fatalf("GetReasonCodes failed: resp=%+v err=%v", reasons, err)
	}
	reason, err := c.GetReasonCodeByID(context.Background(), 1, nil)
	if err != nil || reason.Data.Description != "Damaged" {
		t.Fatalf("GetReasonCodeByID failed: resp=%+v err=%v", reason, err)
	}
	_, err = c.GetReasonCodeByID(context.Background(), 0, nil)
	if err == nil || !strings.Contains(err.Error(), "unique id must be greater than zero") {
		t.Fatalf("expected unique id validation error, got %v", err)
	}

	releases, err := c.GetReleases(context.Background(), nil)
	if err != nil || len(releases.Data) != 1 || releases.Data[0].ItemNumber != 1 {
		t.Fatalf("GetReleases failed: resp=%+v err=%v", releases, err)
	}

	report, err := c.CreateReport(context.Background(), ReportRequest{ReportName: "OpenOrders"})
	if err != nil || report.Data.RequestID != "RPT-1" {
		t.Fatalf("CreateReport failed: resp=%+v err=%v", report, err)
	}
	reportStatus, err := c.GetReportStatus(context.Background(), "RPT-1", nil)
	if err != nil || reportStatus.Data.Status != "COMPLETED" {
		t.Fatalf("GetReportStatus failed: resp=%+v err=%v", reportStatus, err)
	}
	_, err = c.GetReportStatus(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "request id is required") {
		t.Fatalf("expected request id validation error, got %v", err)
	}

	sales, err := c.GetSalespersons(context.Background(), nil)
	if err != nil || len(sales.Data) != 1 || sales.Data[0].SalesID != "S1" {
		t.Fatalf("GetSalespersons failed: resp=%+v err=%v", sales, err)
	}
	updatedSales, err := c.UpdateSalesperson(context.Background(), "S1", Salesperson{Active: true})
	if err != nil || updatedSales.StatusCode != http.StatusOK || updatedSales.NoContent {
		t.Fatalf("UpdateSalesperson failed: resp=%+v err=%v", updatedSales, err)
	}
	_, err = c.UpdateSalesperson(context.Background(), "", Salesperson{})
	if err == nil || !strings.Contains(err.Error(), "sales id is required") {
		t.Fatalf("expected sales id validation error, got %v", err)
	}

	shipList, err := c.GetShippingAddresses(context.Background(), nil)
	if err != nil || len(shipList.Data) != 1 || shipList.Data[0].Location != "MAIN" {
		t.Fatalf("GetShippingAddresses failed: resp=%+v err=%v", shipList, err)
	}
	shipOne, err := c.GetShippingAddress(context.Background(), "C1", "MAIN", nil)
	if err != nil || shipOne.Data.City != "Chicago" {
		t.Fatalf("GetShippingAddress failed: resp=%+v err=%v", shipOne, err)
	}
	_, err = c.GetShippingAddress(context.Background(), "", "MAIN", nil)
	if err == nil || !strings.Contains(err.Error(), "customer code is required") {
		t.Fatalf("expected customer code validation error, got %v", err)
	}
	_, err = c.GetShippingAddress(context.Background(), "C1", "", nil)
	if err == nil || !strings.Contains(err.Error(), "location is required") {
		t.Fatalf("expected location validation error, got %v", err)
	}

	createdShip, err := c.CreateShippingAddress(context.Background(), ShippingAddress{CustomerCode: "C1", Location: "ALT"})
	if err != nil || createdShip.Data.Location != "ALT" {
		t.Fatalf("CreateShippingAddress failed: resp=%+v err=%v", createdShip, err)
	}
	updatedShip, err := c.UpdateShippingAddress(context.Background(), "C1", "MAIN", ShippingAddress{City: "Peoria"})
	if err != nil || updatedShip.StatusCode != http.StatusOK || updatedShip.NoContent {
		t.Fatalf("UpdateShippingAddress failed: resp=%+v err=%v", updatedShip, err)
	}
	_, err = c.UpdateShippingAddress(context.Background(), "", "MAIN", ShippingAddress{})
	if err == nil || !strings.Contains(err.Error(), "customer code is required") {
		t.Fatalf("expected customer code validation error, got %v", err)
	}
	_, err = c.UpdateShippingAddress(context.Background(), "C1", "", ShippingAddress{})
	if err == nil || !strings.Contains(err.Error(), "location is required") {
		t.Fatalf("expected location validation error, got %v", err)
	}

	taxCodes, err := c.GetTaxCodes(context.Background(), nil)
	if err != nil || len(taxCodes.Data) != 1 || taxCodes.Data[0].TaxCode != "IL" {
		t.Fatalf("GetTaxCodes failed: resp=%+v err=%v", taxCodes, err)
	}
	taxCode, err := c.GetTaxCode(context.Background(), "IL", nil)
	if err != nil || taxCode.Data.TaxRate != 7.25 {
		t.Fatalf("GetTaxCode failed: resp=%+v err=%v", taxCode, err)
	}
	_, err = c.GetTaxCode(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "tax code is required") {
		t.Fatalf("expected tax code validation error, got %v", err)
	}

	terms, err := c.GetTerms(context.Background(), nil)
	if err != nil || len(terms.Data) != 1 || terms.Data[0].TermsCode != "N30" {
		t.Fatalf("GetTerms failed: resp=%+v err=%v", terms, err)
	}
	term, err := c.GetTerm(context.Background(), "N30", nil)
	if err != nil || term.Data.DaysDue != 30 {
		t.Fatalf("GetTerm failed: resp=%+v err=%v", term, err)
	}
	_, err = c.GetTerm(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "terms code is required") {
		t.Fatalf("expected terms code validation error, got %v", err)
	}
}
