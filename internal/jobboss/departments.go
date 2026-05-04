package jobboss

import (
	"context"
	"net/url"
)

// Department models a department payload.
type Department struct {
	DepartmentNumber   string `json:"departmentNumber,omitempty"`
	Description        string `json:"description,omitempty"`
	GLCode             string `json:"GLCode,omitempty"`
	IsActive           bool   `json:"isActive,omitempty"`
	DefaultWorkCenter  string `json:"defaultWorkCenter,omitempty"`
	Shift1EmployeeCode string `json:"shift1DefaultEmployeeCode,omitempty"`
	Shift2EmployeeCode string `json:"shift2DefaultEmployeeCode,omitempty"`
	Shift3EmployeeCode string `json:"shift3DefaultEmployeeCode,omitempty"`
	LastModDate        string `json:"lastModDate,omitempty"`
	LastModUser        string `json:"lastModUser,omitempty"`
}

// GetDepartments fetches departments from /api/v1/departments.
func (c *Client) GetDepartments(ctx context.Context, query url.Values) (ListResponse[Department], error) {
	return GetJSON[ListResponse[Department]](ctx, c, "/api/v1/departments", query)
}
