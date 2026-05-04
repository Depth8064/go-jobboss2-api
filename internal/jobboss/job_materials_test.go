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

func TestJobAndMaterialsGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/job-materials":
			_ = json.NewEncoder(w).Encode(ListResponse[JobMaterial]{Data: []JobMaterial{{UniqueID: 11, JobNumber: "J-1", PartNumber: "PN-1"}}})
		case "/api/v1/job-materials/11":
			_ = json.NewEncoder(w).Encode(SingleResponse[JobMaterial]{Data: JobMaterial{UniqueID: 11, JobNumber: "J-1", PartNumber: "PN-1"}})
		case "/api/v1/job-requirements":
			_ = json.NewEncoder(w).Encode(ListResponse[JobRequirement]{Data: []JobRequirement{{UniqueID: 22, JobNumber: "J-2", PartNumber: "PN-2"}}})
		case "/api/v1/job-requirements/22":
			_ = json.NewEncoder(w).Encode(SingleResponse[JobRequirement]{Data: JobRequirement{UniqueID: 22, JobNumber: "J-2", PartNumber: "PN-2"}})
		case "/api/v1/materials":
			_ = json.NewEncoder(w).Encode(ListResponse[Material]{Data: []Material{{UniqueID: 33, PartNumber: "P-EST", SubPartNumber: "SUB-1"}}})
		case "/api/v1/estimates/P-EST/materials/SUB-1":
			_ = json.NewEncoder(w).Encode(SingleResponse[Material]{Data: Material{UniqueID: 33, PartNumber: "P-EST", SubPartNumber: "SUB-1"}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
	q := url.Values{}
	q.Set("fields", "uniqueID,partNumber")

	jm, err := c.GetJobMaterials(context.Background(), q)
	if err != nil || len(jm.Data) != 1 || jm.Data[0].UniqueID != 11 {
		t.Fatalf("GetJobMaterials failed: resp=%+v err=%v", jm, err)
	}
	jmo, err := c.GetJobMaterialByID(context.Background(), "11", nil)
	if err != nil || jmo.Data.UniqueID != 11 {
		t.Fatalf("GetJobMaterialByID failed: resp=%+v err=%v", jmo, err)
	}
	jr, err := c.GetJobRequirements(context.Background(), q)
	if err != nil || len(jr.Data) != 1 || jr.Data[0].UniqueID != 22 {
		t.Fatalf("GetJobRequirements failed: resp=%+v err=%v", jr, err)
	}
	jro, err := c.GetJobRequirementByID(context.Background(), "22", nil)
	if err != nil || jro.Data.UniqueID != 22 {
		t.Fatalf("GetJobRequirementByID failed: resp=%+v err=%v", jro, err)
	}
	m, err := c.GetMaterials(context.Background(), q)
	if err != nil || len(m.Data) != 1 || m.Data[0].UniqueID != 33 {
		t.Fatalf("GetMaterials failed: resp=%+v err=%v", m, err)
	}
	mo, err := c.GetEstimateMaterial(context.Background(), "P-EST", "SUB-1", nil)
	if err != nil || mo.Data.UniqueID != 33 {
		t.Fatalf("GetEstimateMaterial failed: resp=%+v err=%v", mo, err)
	}

	_, err = c.GetJobMaterialByID(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "unique id is required") {
		t.Fatalf("expected unique id validation error, got %v", err)
	}
	_, err = c.GetJobRequirementByID(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "unique id is required") {
		t.Fatalf("expected unique id validation error, got %v", err)
	}
	_, err = c.GetEstimateMaterial(context.Background(), "", "SUB-1", nil)
	if err == nil || !strings.Contains(err.Error(), "part number is required") {
		t.Fatalf("expected part number validation error, got %v", err)
	}
	_, err = c.GetEstimateMaterial(context.Background(), "P-EST", "", nil)
	if err == nil || !strings.Contains(err.Error(), "sub part number is required") {
		t.Fatalf("expected sub part number validation error, got %v", err)
	}
}
