package jobboss

import (
	"context"
	"net/url"
)

// Routing models a routing payload.
type Routing struct {
	UniqueID         int     `json:"uniqueID,omitempty"`
	OperationCode    string  `json:"operationCode,omitempty"`
	Description      string  `json:"description,omitempty"`
	DepartmentNumber string  `json:"departmentNumber,omitempty"`
	WorkCenter       string  `json:"workCenter,omitempty"`
	SetupTime        float64 `json:"setupTime,omitempty"`
	RunTime          float64 `json:"runTime,omitempty"`
	TimeUnit         string  `json:"timeUnit,omitempty"`
	LastModDate      string  `json:"lastModDate,omitempty"`
	LastModUser      string  `json:"lastModUser,omitempty"`
}

// ToolingMaintenance models a tooling-maintenance payload.
type ToolingMaintenance struct {
	UniqueID          int    `json:"uniqueID,omitempty"`
	ToolNumber        string `json:"toolNumber,omitempty"`
	Description       string `json:"description,omitempty"`
	MaintenanceType   string `json:"maintenanceType,omitempty"`
	Frequency         string `json:"frequency,omitempty"`
	LastPerformedDate string `json:"lastPerformedDate,omitempty"`
	DueDate           string `json:"dueDate,omitempty"`
	Status            string `json:"status,omitempty"`
	LastModDate       string `json:"lastModDate,omitempty"`
	LastModUser       string `json:"lastModUser,omitempty"`
}

// UserLabel models a user label payload.
type UserLabel struct {
	UniqueID     int    `json:"uniqueID,omitempty"`
	LabelName    string `json:"labelName,omitempty"`
	DisplayText  string `json:"displayText,omitempty"`
	Module       string `json:"module,omitempty"`
	LanguageCode string `json:"languageCode,omitempty"`
	LastModDate  string `json:"lastModDate,omitempty"`
	LastModUser  string `json:"lastModUser,omitempty"`
}

// UserTransaction models a user transaction payload.
type UserTransaction struct {
	UniqueID        int    `json:"uniqueID,omitempty"`
	UserCode        string `json:"userCode,omitempty"`
	TransactionType string `json:"transactionType,omitempty"`
	Description     string `json:"description,omitempty"`
	TransactionDate string `json:"transactionDate,omitempty"`
	Reference       string `json:"reference,omitempty"`
	LastModDate     string `json:"lastModDate,omitempty"`
	LastModUser     string `json:"lastModUser,omitempty"`
}

// GetRoutings fetches routings.
func (c *Client) GetRoutings(ctx context.Context, query url.Values) (ListResponse[Routing], error) {
	return GetJSON[ListResponse[Routing]](ctx, c, "/api/v1/routings", query)
}

// GetToolingMaintenance fetches tooling-maintenance records.
func (c *Client) GetToolingMaintenance(ctx context.Context, query url.Values) (ListResponse[ToolingMaintenance], error) {
	return GetJSON[ListResponse[ToolingMaintenance]](ctx, c, "/api/v1/tooling-maintenance", query)
}

// GetUserLabels fetches user labels.
func (c *Client) GetUserLabels(ctx context.Context, query url.Values) (ListResponse[UserLabel], error) {
	return GetJSON[ListResponse[UserLabel]](ctx, c, "/api/v1/user-labels", query)
}

// GetUserTransactions fetches user transactions.
func (c *Client) GetUserTransactions(ctx context.Context, query url.Values) (ListResponse[UserTransaction], error) {
	return GetJSON[ListResponse[UserTransaction]](ctx, c, "/api/v1/user-transactions", query)
}
