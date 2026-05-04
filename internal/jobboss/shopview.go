package jobboss

import (
	"context"
	"net/url"
)

// ShopViewAuthenticateRequest models a ShopView authentication request payload.
type ShopViewAuthenticateRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// ShopViewAuthenticateResponse models a ShopView authentication response payload.
type ShopViewAuthenticateResponse struct {
	Authenticated bool   `json:"authenticated,omitempty"`
	Token         string `json:"token,omitempty"`
	Message       string `json:"message,omitempty"`
}

// ShopViewFilter models a ShopView filter payload.
type ShopViewFilter struct {
	Name    string   `json:"name,omitempty"`
	Type    string   `json:"type,omitempty"`
	Options []string `json:"options,omitempty"`
}

// ShopViewJob models a ShopView job payload.
type ShopViewJob struct {
	JobNumber  string `json:"jobNumber,omitempty"`
	PartNumber string `json:"partNumber,omitempty"`
	WorkCenter string `json:"workCenter,omitempty"`
	Status     string `json:"status,omitempty"`
	StartDate  string `json:"startDate,omitempty"`
	DueDate    string `json:"dueDate,omitempty"`
}

// ShopViewGridOption models a ShopView grid-option payload.
type ShopViewGridOption struct {
	OptionName  string `json:"optionName,omitempty"`
	OptionValue string `json:"optionValue,omitempty"`
}

// ShopViewKPIDefinition models a ShopView KPI definition payload.
type ShopViewKPIDefinition struct {
	Key         string `json:"key,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// ShopViewKPIValue models a ShopView KPI value payload.
type ShopViewKPIValue struct {
	Value int `json:"value,omitempty"`
}

// ShopViewResetResponse models a ShopView reset-grid-options response payload.
type ShopViewResetResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

// AuthenticateShopView authenticates a ShopView user.
func (c *Client) AuthenticateShopView(ctx context.Context, req ShopViewAuthenticateRequest) (SingleResponse[ShopViewAuthenticateResponse], error) {
	return PostJSON[ShopViewAuthenticateRequest, SingleResponse[ShopViewAuthenticateResponse]](ctx, c, "/api/v1/shopview/authenticate-user", req)
}

// GetShopViewFilters fetches ShopView filters.
func (c *Client) GetShopViewFilters(ctx context.Context, query url.Values) (ListResponse[ShopViewFilter], error) {
	return GetJSON[ListResponse[ShopViewFilter]](ctx, c, "/api/v1/shopview/filters", query)
}

// GetShopViewJobs fetches ShopView jobs.
func (c *Client) GetShopViewJobs(ctx context.Context, query url.Values) (ListResponse[ShopViewJob], error) {
	return GetJSON[ListResponse[ShopViewJob]](ctx, c, "/api/v1/shopview/get-jobs", query)
}

// CreateShopViewGridOption creates a ShopView grid option.
func (c *Client) CreateShopViewGridOption(ctx context.Context, option ShopViewGridOption) (SingleResponse[ShopViewGridOption], error) {
	return PostJSON[ShopViewGridOption, SingleResponse[ShopViewGridOption]](ctx, c, "/api/v1/shopview/grid-option", option)
}

// GetShopViewGridOptions fetches ShopView grid options.
func (c *Client) GetShopViewGridOptions(ctx context.Context, query url.Values) (ListResponse[ShopViewGridOption], error) {
	return GetJSON[ListResponse[ShopViewGridOption]](ctx, c, "/api/v1/shopview/grid-options", query)
}

// GetShopViewKPIDefinitions fetches ShopView KPI definitions.
func (c *Client) GetShopViewKPIDefinitions(ctx context.Context, query url.Values) (ListResponse[ShopViewKPIDefinition], error) {
	return GetJSON[ListResponse[ShopViewKPIDefinition]](ctx, c, "/api/v1/shopview/kpi-definitions", query)
}

// GetShopViewKPIJobsClosed fetches the ShopView jobs-closed KPI.
func (c *Client) GetShopViewKPIJobsClosed(ctx context.Context, query url.Values) (SingleResponse[ShopViewKPIValue], error) {
	return GetJSON[SingleResponse[ShopViewKPIValue]](ctx, c, "/api/v1/shopview/kpi/jobs-closed", query)
}

// GetShopViewKPIJobsInProgress fetches the ShopView jobs-in-progress KPI.
func (c *Client) GetShopViewKPIJobsInProgress(ctx context.Context, query url.Values) (SingleResponse[ShopViewKPIValue], error) {
	return GetJSON[SingleResponse[ShopViewKPIValue]](ctx, c, "/api/v1/shopview/kpi/jobs-in-progress", query)
}

// GetShopViewKPIJobsOnHold fetches the ShopView jobs-on-hold KPI.
func (c *Client) GetShopViewKPIJobsOnHold(ctx context.Context, query url.Values) (SingleResponse[ShopViewKPIValue], error) {
	return GetJSON[SingleResponse[ShopViewKPIValue]](ctx, c, "/api/v1/shopview/kpi/jobs-on-hold", query)
}

// GetShopViewKPIJobsPastDue fetches the ShopView jobs-past-due KPI.
func (c *Client) GetShopViewKPIJobsPastDue(ctx context.Context, query url.Values) (SingleResponse[ShopViewKPIValue], error) {
	return GetJSON[SingleResponse[ShopViewKPIValue]](ctx, c, "/api/v1/shopview/kpi/jobs-past-due", query)
}

// ResetShopViewGridOptions resets ShopView grid options.
func (c *Client) ResetShopViewGridOptions(ctx context.Context) (SingleResponse[ShopViewResetResponse], error) {
	return PostJSON[map[string]string, SingleResponse[ShopViewResetResponse]](ctx, c, "/api/v1/shopview/reset-grid-options", map[string]string{})
}
