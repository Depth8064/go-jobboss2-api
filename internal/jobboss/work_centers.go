package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// WorkCenter models a work center payload.
type WorkCenter struct {
	UniqueID         int     `json:"uniqueID,omitempty"`
	WorkCenter       string  `json:"workCenter,omitempty"`
	Description      string  `json:"description,omitempty"`
	DepartmentNumber string  `json:"departmentNumber,omitempty"`
	Capacity         float64 `json:"capacity,omitempty"`
	Active           bool    `json:"active,omitempty"`
	LastModDate      string  `json:"lastModDate,omitempty"`
	LastModUser      string  `json:"lastModUser,omitempty"`
}

// WorkCenterMaintenance models a work-center-maintenance payload.
type WorkCenterMaintenance struct {
	UniqueID          int    `json:"uniqueID,omitempty"`
	WorkCenter        string `json:"workCenter,omitempty"`
	Description       string `json:"description,omitempty"`
	MaintenanceType   string `json:"maintenanceType,omitempty"`
	ScheduleFrequency string `json:"scheduleFrequency,omitempty"`
	LastPerformedDate string `json:"lastPerformedDate,omitempty"`
	NextDueDate       string `json:"nextDueDate,omitempty"`
	Status            string `json:"status,omitempty"`
	LastModDate       string `json:"lastModDate,omitempty"`
	LastModUser       string `json:"lastModUser,omitempty"`
}

// GetWorkCenters fetches work centers.
func (c *Client) GetWorkCenters(ctx context.Context, query url.Values) (ListResponse[WorkCenter], error) {
	return GetJSON[ListResponse[WorkCenter]](ctx, c, "/api/v1/work-centers", query)
}

// GetWorkCenterByCode fetches a work center by code.
func (c *Client) GetWorkCenterByCode(ctx context.Context, workCenter string, query url.Values) (SingleResponse[WorkCenter], error) {
	workCenter = strings.TrimSpace(workCenter)
	if workCenter == "" {
		return SingleResponse[WorkCenter]{}, fmt.Errorf("work center is required")
	}
	endpoint := "/api/v1/work-centers/" + url.PathEscape(workCenter)
	return GetJSON[SingleResponse[WorkCenter]](ctx, c, endpoint, query)
}

// CreateWorkCenter creates a work center.
func (c *Client) CreateWorkCenter(ctx context.Context, workCenter WorkCenter) (SingleResponse[WorkCenter], error) {
	return PostJSON[WorkCenter, SingleResponse[WorkCenter]](ctx, c, "/api/v1/work-centers", workCenter)
}

// UpdateWorkCenter updates a work center by code.
func (c *Client) UpdateWorkCenter(ctx context.Context, workCenter string, payload WorkCenter) (UpdateResult, error) {
	workCenter = strings.TrimSpace(workCenter)
	if workCenter == "" {
		return UpdateResult{}, fmt.Errorf("work center is required")
	}
	endpoint := "/api/v1/work-centers/" + url.PathEscape(workCenter)
	return PatchNoContentJSON[WorkCenter](ctx, c, endpoint, payload)
}

// GetWorkCenterMaintenance fetches work-center-maintenance records.
func (c *Client) GetWorkCenterMaintenance(ctx context.Context, query url.Values) (ListResponse[WorkCenterMaintenance], error) {
	return GetJSON[ListResponse[WorkCenterMaintenance]](ctx, c, "/api/v1/work-center-maintenance", query)
}
