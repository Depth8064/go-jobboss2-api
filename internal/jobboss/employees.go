package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// EmployeeTraining models an employee training payload.
type EmployeeTraining struct {
	UniqueID      int     `json:"uniqueID,omitempty"`
	EmployeeCode  int     `json:"employeeCode,omitempty"`
	EmployeeName  string  `json:"employeeName,omitempty"`
	TrainingCode  string  `json:"trainingCode,omitempty"`
	Description   string  `json:"description,omitempty"`
	Instructor    string  `json:"instructor,omitempty"`
	VendorCode    string  `json:"vendorCode,omitempty"`
	InvoiceNumber string  `json:"invoiceNumber,omitempty"`
	Location      string  `json:"location,omitempty"`
	StartDate     string  `json:"startDate,omitempty"`
	EndDate       string  `json:"endDate,omitempty"`
	Cost          float64 `json:"cost,omitempty"`
	Completed     bool    `json:"completed,omitempty"`
	LastModDate   string  `json:"lastModDate,omitempty"`
	LastModUser   string  `json:"lastModUser,omitempty"`
}

// Employee models an employee payload.
type Employee struct {
	UniqueID          int     `json:"uniqueID,omitempty"`
	EmployeeCode      int     `json:"employeeCode,omitempty"`
	EmployeeName      string  `json:"employeeName,omitempty"`
	EmployeeShortName string  `json:"employeeShortName,omitempty"`
	Active            bool    `json:"active,omitempty"`
	DepartmentNumber  string  `json:"departmentNumber,omitempty"`
	ShiftID           int     `json:"shiftID,omitempty"`
	Phone             string  `json:"phone,omitempty"`
	City              string  `json:"city,omitempty"`
	State             string  `json:"state,omitempty"`
	Country           string  `json:"country,omitempty"`
	Rate1             float64 `json:"rate1,omitempty"`
	LastModDate       string  `json:"lastModDate,omitempty"`
	LastModUser       string  `json:"lastModUser,omitempty"`
}

// EmployeeUpdate is a payload for updating employees.
type EmployeeUpdate struct {
	DepartmentNumber *string  `json:"departmentNumber,omitempty"`
	EmployeeName     *string  `json:"employeeName,omitempty"`
	ShiftID          *int     `json:"shiftID,omitempty"`
	Rate1            *float64 `json:"rate1,omitempty"`
	Phone            *string  `json:"phone,omitempty"`
}

// GetEmployeeTrainings fetches employee trainings.
func (c *Client) GetEmployeeTrainings(ctx context.Context, query url.Values) (ListResponse[EmployeeTraining], error) {
	return GetJSON[ListResponse[EmployeeTraining]](ctx, c, "/api/v1/employee-trainings", query)
}

// GetEmployees fetches employees.
func (c *Client) GetEmployees(ctx context.Context, query url.Values) (ListResponse[Employee], error) {
	return GetJSON[ListResponse[Employee]](ctx, c, "/api/v1/employees", query)
}

// UpdateEmployee updates an employee by employee code.
func (c *Client) UpdateEmployee(ctx context.Context, employeeCode string, payload EmployeeUpdate) (SingleResponse[Employee], error) {
	employeeCode = strings.TrimSpace(employeeCode)
	if employeeCode == "" {
		return SingleResponse[Employee]{}, fmt.Errorf("employee code is required")
	}
	endpoint := "/api/v1/employees/" + url.PathEscape(employeeCode)
	return PatchJSON[EmployeeUpdate, SingleResponse[Employee]](ctx, c, endpoint, payload)
}
