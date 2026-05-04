package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// OperationCode models an operation code payload.
type OperationCode struct {
	UniqueID              int     `json:"uniqueID,omitempty"`
	OperationCode         string  `json:"operationCode,omitempty"`
	OperationNumber       int     `json:"operationNumber,omitempty"`
	Description           string  `json:"description,omitempty"`
	Active                bool    `json:"active,omitempty"`
	DefaultCycleTime      float64 `json:"defaultCycleTime,omitempty"`
	DefaultCycleUnit      string  `json:"defaultCycleUnit,omitempty"`
	DefaultSetupTime      float64 `json:"defaultSetupTime,omitempty"`
	DefaultTimeUnit       string  `json:"defaultTimeUnit,omitempty"`
	MachinesRunByOperator int     `json:"machinesRunByOperator,omitempty"`
	MachinesJobRunOn      int     `json:"machinesJobRunOn,omitempty"`
	PercentEfficiency     float64 `json:"percentEfficiency,omitempty"`
	ScrapPercent          float64 `json:"scrapPercent,omitempty"`
	TeamSize              int     `json:"teamSize,omitempty"`
	LastModDate           string  `json:"lastModDate,omitempty"`
	LastModUser           string  `json:"lastModUser,omitempty"`
}

// GetOperationCodes fetches operation codes.
func (c *Client) GetOperationCodes(ctx context.Context, query url.Values) (ListResponse[OperationCode], error) {
	return GetJSON[ListResponse[OperationCode]](ctx, c, "/api/v1/operation-codes", query)
}

// GetOperationCodeByCode fetches an operation code by operation code.
func (c *Client) GetOperationCodeByCode(ctx context.Context, operationCode string, query url.Values) (SingleResponse[OperationCode], error) {
	operationCode = strings.TrimSpace(operationCode)
	if operationCode == "" {
		return SingleResponse[OperationCode]{}, fmt.Errorf("operation code is required")
	}
	endpoint := "/api/v1/operation-codes/" + url.PathEscape(operationCode)
	return GetJSON[SingleResponse[OperationCode]](ctx, c, endpoint, query)
}
