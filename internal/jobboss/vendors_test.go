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

func TestVendorsGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/vendors":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(ListResponse[Vendor]{Data: []Vendor{{VendorCode: "V1", VendorName: "Acme Supply"}}})
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(SingleResponse[Vendor]{Data: Vendor{VendorCode: "V2", VendorName: "New Vendor"}})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/api/v1/vendors/V1":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(SingleResponse[Vendor]{Data: Vendor{VendorCode: "V1", Status: "ACTIVE"}})
			case http.MethodPatch:
				_ = json.NewEncoder(w).Encode(SingleResponse[Vendor]{Data: Vendor{VendorCode: "V1", Status: "INACTIVE"}})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/api/v1/vendor-returns":
			_ = json.NewEncoder(w).Encode(ListResponse[VendorReturn]{Data: []VendorReturn{{VendorReturnNumber: "VR-1", VendorCode: "V1"}}})
		case "/api/v1/vendor-returns-line-items":
			_ = json.NewEncoder(w).Encode(ListResponse[VendorReturnLineItem]{Data: []VendorReturnLineItem{{VendorReturnNumber: "VR-1", ItemNumber: 1}}})
		case "/api/v1/vendor-returns-releases":
			_ = json.NewEncoder(w).Encode(ListResponse[VendorReturnRelease]{Data: []VendorReturnRelease{{VendorReturnNumber: "VR-1", ReleaseNumber: 1}}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
	q := url.Values{}
	q.Set("fields", "vendorCode,vendorName")

	vendors, err := c.GetVendors(context.Background(), q)
	if err != nil || len(vendors.Data) != 1 || vendors.Data[0].VendorCode != "V1" {
		t.Fatalf("GetVendors failed: resp=%+v err=%v", vendors, err)
	}
	vendor, err := c.GetVendorByCode(context.Background(), "V1", nil)
	if err != nil || vendor.Data.Status != "ACTIVE" {
		t.Fatalf("GetVendorByCode failed: resp=%+v err=%v", vendor, err)
	}
	_, err = c.GetVendorByCode(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "vendor code is required") {
		t.Fatalf("expected vendor code validation error, got %v", err)
	}

	created, err := c.CreateVendor(context.Background(), Vendor{VendorCode: "V2"})
	if err != nil || created.Data.VendorCode != "V2" {
		t.Fatalf("CreateVendor failed: resp=%+v err=%v", created, err)
	}
	updated, err := c.UpdateVendor(context.Background(), "V1", Vendor{Status: "INACTIVE"})
	if err != nil || updated.Data.Status != "INACTIVE" {
		t.Fatalf("UpdateVendor failed: resp=%+v err=%v", updated, err)
	}
	_, err = c.UpdateVendor(context.Background(), "", Vendor{})
	if err == nil || !strings.Contains(err.Error(), "vendor code is required") {
		t.Fatalf("expected vendor code validation error, got %v", err)
	}

	returns, err := c.GetVendorReturns(context.Background(), nil)
	if err != nil || len(returns.Data) != 1 || returns.Data[0].VendorReturnNumber != "VR-1" {
		t.Fatalf("GetVendorReturns failed: resp=%+v err=%v", returns, err)
	}
	items, err := c.GetVendorReturnLineItems(context.Background(), nil)
	if err != nil || len(items.Data) != 1 || items.Data[0].ItemNumber != 1 {
		t.Fatalf("GetVendorReturnLineItems failed: resp=%+v err=%v", items, err)
	}
	releases, err := c.GetVendorReturnReleases(context.Background(), nil)
	if err != nil || len(releases.Data) != 1 || releases.Data[0].ReleaseNumber != 1 {
		t.Fatalf("GetVendorReturnReleases failed: resp=%+v err=%v", releases, err)
	}
}
