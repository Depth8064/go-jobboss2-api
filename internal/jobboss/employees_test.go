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

func TestEmployeesGroup(t *testing.T) {
	name := "John Smith"
	shift := 2

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/employee-trainings":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(ListResponse[EmployeeTraining]{
				Data: []EmployeeTraining{{EmployeeCode: 101, EmployeeName: "John", TrainingCode: "SAFE", Completed: true}},
			})
		case "/api/v1/employees":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(ListResponse[Employee]{
				Data: []Employee{{EmployeeCode: 101, EmployeeName: "John", DepartmentNumber: "MACH", Rate1: 20}},
			})
		case "/api/v1/employees/101":
			if r.Method != http.MethodPatch {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[Employee]{
				Data: Employee{EmployeeCode: 101, EmployeeName: "John Smith", ShiftID: 2},
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
	q.Set("fields", "employeeCode,employeeName,trainingCode")
	trainings, err := c.GetEmployeeTrainings(context.Background(), q)
	if err != nil || len(trainings.Data) != 1 || trainings.Data[0].TrainingCode != "SAFE" {
		t.Fatalf("GetEmployeeTrainings failed: resp=%+v err=%v", trainings, err)
	}

	employees, err := c.GetEmployees(context.Background(), q)
	if err != nil || len(employees.Data) != 1 || employees.Data[0].EmployeeCode != 101 {
		t.Fatalf("GetEmployees failed: resp=%+v err=%v", employees, err)
	}

	updated, err := c.UpdateEmployee(context.Background(), "101", EmployeeUpdate{EmployeeName: &name, ShiftID: &shift})
	if err != nil || updated.Data.ShiftID != 2 || updated.Data.EmployeeName != "John Smith" {
		t.Fatalf("UpdateEmployee failed: resp=%+v err=%v", updated, err)
	}

	_, err = c.UpdateEmployee(context.Background(), "", EmployeeUpdate{})
	if err == nil || !strings.Contains(err.Error(), "employee code is required") {
		t.Fatalf("expected employee code validation error, got %v", err)
	}
}
