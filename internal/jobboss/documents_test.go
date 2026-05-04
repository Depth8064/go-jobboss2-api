package jobboss

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestDocumentsGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/document-controls":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(ListResponse[DocumentControl]{
				Data: []DocumentControl{{DocumentNumber: "DOC-100", DocumentStatus: "ACTIVE", Revision: "A"}},
			})
		case "/api/v1/document-histories":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(ListResponse[DocumentHistory]{
				Data: []DocumentHistory{{DocumentNumber: "DOC-100", Revision: "A", UserID: "ADMIN"}},
			})
		case "/api/v1/document-review":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(ListResponse[DocumentReview]{
				Data: []DocumentReview{{DocumentNumber: "DOC-100", ReviewCode: "VEND", VendorCode: "V-1", Cost: 10}},
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

	q := url.Values{}
	q.Set("fields", "documentNumber,documentStatus,revision")
	controls, err := c.GetDocumentControls(context.Background(), q)
	if err != nil || len(controls.Data) != 1 || controls.Data[0].DocumentNumber != "DOC-100" {
		t.Fatalf("GetDocumentControls failed: resp=%+v err=%v", controls, err)
	}

	histories, err := c.GetDocumentHistories(context.Background(), q)
	if err != nil || len(histories.Data) != 1 || histories.Data[0].UserID != "ADMIN" {
		t.Fatalf("GetDocumentHistories failed: resp=%+v err=%v", histories, err)
	}

	reviews, err := c.GetDocumentReviews(context.Background(), q)
	if err != nil || len(reviews.Data) != 1 || reviews.Data[0].ReviewCode != "VEND" {
		t.Fatalf("GetDocumentReviews failed: resp=%+v err=%v", reviews, err)
	}
}
