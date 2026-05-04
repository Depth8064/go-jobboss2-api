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

func TestPurchaseAndPackingGroups(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/packing-lists":
			_ = json.NewEncoder(w).Encode(ListResponse[PackingList]{Data: []PackingList{{PackingListNo: "PL-1", OrderNumber: "SO-1"}}})
		case "/api/v1/packing-list-line-items":
			_ = json.NewEncoder(w).Encode(ListResponse[PackingListLineItem]{Data: []PackingListLineItem{{PackingListNo: "PL-1", PartNumber: "P1", ItemNumber: 1}}})
		case "/api/v1/purchase-orders":
			_ = json.NewEncoder(w).Encode(ListResponse[PurchaseOrder]{Data: []PurchaseOrder{{PONumber: "PO-1", VendorCode: "V1"}}})
		case "/api/v1/purchase-orders/PO-1":
			_ = json.NewEncoder(w).Encode(SingleResponse[PurchaseOrder]{Data: PurchaseOrder{PONumber: "PO-1", Status: "OPEN"}})
		case "/api/v1/purchase-order-line-items":
			_ = json.NewEncoder(w).Encode(ListResponse[PurchaseOrderLineItem]{Data: []PurchaseOrderLineItem{{PONumber: "PO-1", PartNumber: "P1", ItemNumber: 1}}})
		case "/api/v1/purchase-order-line-items/PO-1/P1/1":
			_ = json.NewEncoder(w).Encode(SingleResponse[PurchaseOrderLineItem]{Data: PurchaseOrderLineItem{PONumber: "PO-1", PartNumber: "P1", ItemNumber: 1, Description: "Line"}})
		case "/api/v1/purchase-order-releases":
			_ = json.NewEncoder(w).Encode(ListResponse[PurchaseOrderRelease]{Data: []PurchaseOrderRelease{{PONumber: "PO-1", ReleaseNumber: 1}}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
	q := url.Values{}
	q.Set("fields", "poNumber,status")

	packing, err := c.GetPackingLists(context.Background(), q)
	if err != nil || len(packing.Data) != 1 || packing.Data[0].PackingListNo != "PL-1" {
		t.Fatalf("GetPackingLists failed: resp=%+v err=%v", packing, err)
	}
	packingItems, err := c.GetPackingListLineItems(context.Background(), nil)
	if err != nil || len(packingItems.Data) != 1 || packingItems.Data[0].PartNumber != "P1" {
		t.Fatalf("GetPackingListLineItems failed: resp=%+v err=%v", packingItems, err)
	}

	orders, err := c.GetPurchaseOrders(context.Background(), q)
	if err != nil || len(orders.Data) != 1 || orders.Data[0].PONumber != "PO-1" {
		t.Fatalf("GetPurchaseOrders failed: resp=%+v err=%v", orders, err)
	}
	order, err := c.GetPurchaseOrderByNumber(context.Background(), "PO-1", nil)
	if err != nil || order.Data.Status != "OPEN" {
		t.Fatalf("GetPurchaseOrderByNumber failed: resp=%+v err=%v", order, err)
	}
	_, err = c.GetPurchaseOrderByNumber(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "po number is required") {
		t.Fatalf("expected po number validation error, got %v", err)
	}

	lineItems, err := c.GetPurchaseOrderLineItems(context.Background(), nil)
	if err != nil || len(lineItems.Data) != 1 || lineItems.Data[0].ItemNumber != 1 {
		t.Fatalf("GetPurchaseOrderLineItems failed: resp=%+v err=%v", lineItems, err)
	}
	lineItem, err := c.GetPurchaseOrderLineItem(context.Background(), "PO-1", "P1", "1", nil)
	if err != nil || lineItem.Data.Description != "Line" {
		t.Fatalf("GetPurchaseOrderLineItem failed: resp=%+v err=%v", lineItem, err)
	}
	_, err = c.GetPurchaseOrderLineItem(context.Background(), "", "P1", "1", nil)
	if err == nil || !strings.Contains(err.Error(), "po number is required") {
		t.Fatalf("expected po number validation error, got %v", err)
	}
	_, err = c.GetPurchaseOrderLineItem(context.Background(), "PO-1", "", "1", nil)
	if err == nil || !strings.Contains(err.Error(), "part number is required") {
		t.Fatalf("expected part number validation error, got %v", err)
	}
	_, err = c.GetPurchaseOrderLineItem(context.Background(), "PO-1", "P1", "", nil)
	if err == nil || !strings.Contains(err.Error(), "item number is required") {
		t.Fatalf("expected item number validation error, got %v", err)
	}

	releases, err := c.GetPurchaseOrderReleases(context.Background(), nil)
	if err != nil || len(releases.Data) != 1 || releases.Data[0].ReleaseNumber != 1 {
		t.Fatalf("GetPurchaseOrderReleases failed: resp=%+v err=%v", releases, err)
	}
}
