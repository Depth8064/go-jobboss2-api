package jobboss

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// DeviceCredentialType values from Authentication API.
type DeviceCredentialType string

const (
	DeviceCredentialNone           DeviceCredentialType = "None"
	DeviceCredentialToken          DeviceCredentialType = "Token"
	DeviceCredentialSecret         DeviceCredentialType = "Secret"
	DeviceCredentialTokenAndSecret DeviceCredentialType = "TokenAndSecret"
)

type DeviceCreateModel struct {
	ClientID       string                `json:"clientId"`
	CredentialType *DeviceCredentialType `json:"credentialType,omitempty"`
	Description    string                `json:"description,omitempty"`
	Origin         string                `json:"origin,omitempty"`
	DeviceType     string                `json:"deviceType,omitempty"`
}

type DeviceCreateResponse struct {
	Code        string `json:"code,omitempty"`
	ID          string `json:"id,omitempty"`
	CodeExpires string `json:"codeExpires,omitempty"`
}

type DeviceActivateRequestModel struct {
	Code string `json:"code"`
}

type DeviceCredentialModel struct {
	Secret        string `json:"secret,omitempty"`
	TokenEndpoint string `json:"tokenEndpoint,omitempty"`
	TokenIssuer   string `json:"tokenIssuer,omitempty"`
	TokenType     string `json:"tokenType,omitempty"`
	AccessToken   string `json:"accessToken,omitempty"`
	IDToken       string `json:"idToken,omitempty"`
	ClientID      string `json:"clientId,omitempty"`
	ExpiresAt     *int64 `json:"expiresAt,omitempty"`
	RefreshToken  string `json:"refreshToken,omitempty"`
}

type DeviceActivatedResponse struct {
	ID          string                 `json:"id,omitempty"`
	IsActivated bool                   `json:"isActivated"`
	Properties  map[string]string      `json:"properties,omitempty"`
	Credential  *DeviceCredentialModel `json:"credential,omitempty"`
}

// CreateDevice calls POST /api/devices on the Authentication API.
func (c *Client) CreateDevice(ctx context.Context, request DeviceCreateModel) (*DeviceCreateResponse, error) {
	resp, err := PostJSON[DeviceCreateModel, DeviceCreateResponse](ctx, c, "/api/devices", request)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// TryActivateDevice calls POST /api/devices/{id}/activation.
// Returns activated=false with nil error when the API returns 204 (not yet activated).
func (c *Client) TryActivateDevice(ctx context.Context, id string, request DeviceActivateRequestModel) (*DeviceActivatedResponse, bool, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, false, fmt.Errorf("device id is required")
	}

	body, err := c.Post(ctx, "/api/devices/"+url.PathEscape(id)+"/activation", request)
	if err != nil {
		return nil, false, err
	}

	if len(body) == 0 {
		return nil, false, nil
	}

	var resp DeviceActivatedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, false, fmt.Errorf("decode activation response: %w", err)
	}
	return &resp, resp.IsActivated, nil
}
