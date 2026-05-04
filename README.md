# go-jobboss2-api

[![CI](https://github.com/Depth8064/go-jobboss2-api/actions/workflows/ci.yml/badge.svg)](https://github.com/Depth8064/go-jobboss2-api/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Depth8064/go-jobboss2-api)](go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/Depth8064/go-jobboss2-api)](https://goreportcard.com/report/github.com/Depth8064/go-jobboss2-api)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/Depth8064/go-jobboss2-api.svg)](https://pkg.go.dev/github.com/Depth8064/go-jobboss2-api)
[![License](https://img.shields.io/github/license/Depth8064/go-jobboss2-api)](LICENSE)

Comprehensive, dependency-free Go module for JobBOSS2 API access.

This module provides a fully featured, typed client for JobBOSS2 endpoints, including authentication, query helpers, and resource-specific API methods.

## Requirements

- Go 1.26+

## CI

- GitHub Actions runs tests on push to `main` and on pull requests
- CI uses the latest Go 1.26 patch release (`1.26.x`)

## Features

- Full OpenAPI v1 group coverage with typed request and response models
- OAuth2 client-credentials token handling
- Automatic token caching and refresh
- Automatic token refresh-and-retry on 401 (max 3 total attempts by default)
- Automatic retry with exponential backoff + jitter for transient errors (5xx, 429, network)
- Retry-After header support for 429 rate-limit responses
- Resource-specific API methods for create/read/update/delete workflows
- Generic `GET`, `POST`, `PATCH`, `DELETE` helpers for direct endpoint access
- Generic JSON helpers for typed responses
- Explicit metadata-first PATCH helper for 204/empty-body update endpoints
- ECI-aligned defaults for API and auth hosts
- Query helpers for fields/filter/sort/take/skip
- UTC date/date-time formatting helpers
- Authentication API device endpoints (`/api/devices`, `/api/devices/{id}/activation`)
- No runtime dependency on OpenAPI parsing

## Install

```bash
go get github.com/Depth8064/go-jobboss2-api
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"net/url"

	jobboss "github.com/Depth8064/go-jobboss2-api"
)

func main() {
	client := jobboss.NewClient(jobboss.Config{
		BaseURL:      "https://api-jb2.integrations.ecimanufacturing.com",
		TokenURL:     "https://api-user.integrations.ecimanufacturing.com/oauth2/api-user/token",
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
	})

	query := url.Values{}
	query.Set("take", "25")
	query.Set("fields", "customerCode,customerName")

	resp, err := client.GetCustomers(context.Background(), query)
	if err != nil {
		panic(err)
	}

	fmt.Printf("customers: %d\n", len(resp.Data))
}
```

## Typed API Usage

In addition to generic HTTP helpers, the client provides typed resource methods:

```go
ctx := context.Background()

orders, err := client.GetOrders(ctx, nil)
if err != nil {
	panic(err)
}

fmt.Printf("orders: %d\n", len(orders.Data))
```

## Update Semantics (PATCH)

Many JobBOSS update endpoints return `204 No Content` after a successful `PATCH`.

- Resource `Update*` methods now return `UpdateResult` (status metadata), not hydrated entities.
- Use a follow-up `Get*` call when you need the latest entity state.
- `PatchJSON` now expects a JSON response body and returns an error for empty responses.
- Use `PatchNoContentJSON` for metadata-only PATCH success handling.

Example:

```go
result, err := client.UpdateOrder(ctx, "SO-100", jobboss.OrderUpdate{Status: &status})
if err != nil {
	panic(err)
}
if result.StatusCode == 204 {
	updated, err := client.GetOrderByNumber(ctx, "SO-100", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(updated.Data.OrderNumber)
}
```

## Retry and Resilience

The client automatically retries transient failures with exponential backoff and jitter.

**Retryable conditions:** network errors, HTTP 429 Too Many Requests, HTTP 5xx server errors.

**Configuration (all optional — sensible defaults apply):**

| Field | Default | Description |
|---|---|---|
| `MaxRetries` | `3` | Max retry attempts after the initial request. Set to `-1` to disable. |
| `RetryBackoffBase` | `500ms` | Starting delay for exponential backoff. |
| `RetryBackoffMax` | `30s` | Maximum delay cap. |

**Retry-After support:** When the server returns a `429` with a `Retry-After` header (integer seconds or HTTP-date), the client waits exactly that duration before retrying.

```go
client := jobboss.NewClient(jobboss.Config{
	// ...
	MaxRetries:       5,
	RetryBackoffBase: 250 * time.Millisecond,
	RetryBackoffMax:  15 * time.Second,
})
```

On final failure after all retries, the returned error is an `*APIError`. For `429` responses, `APIError.RetryAfter` holds the raw `Retry-After` header value.

## Defaults

If omitted, `NewClient` uses:

- `BaseURL`: `https://api-jb2.integrations.ecimanufacturing.com`
- `TokenURL`: `https://api-user.integrations.ecimanufacturing.com/oauth2/api-user/token`
- `Scope`: `openid`

## Query Helpers

```go
q := jobboss.NewQuery()
q = jobboss.ApplyPaging(q, 200, 0)
q = jobboss.ApplyFields(q, "customerCode", "customerName")
q = jobboss.ApplyFilter(q, "customerCode", jobboss.FilterEqual, "ACME")
q = jobboss.ApplySort(q, "+customerCode")
```

## Date/Time Helpers

JobBOSS expects UTC date formats like `yyyy-MM-dd` and `yyyy-MM-ddTHH:mm:ssZ`.

```go
import "time"

day := jobboss.FormatUTCDate(time.Now())
ts := jobboss.FormatUTCDateTime(time.Now())
```

## Development

- Run tests: `go test ./...`
- Run tests with coverage profile: `go test ./... -coverprofile=coverage.out -covermode=count`
- View coverage summary (bash): `go tool cover -func=coverage.out | tail -n 1`
- View coverage summary (PowerShell): `go tool cover -func=coverage.out | Select-Object -Last 1`

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution workflow and quality checks.

## Security

See [SECURITY.md](SECURITY.md) for vulnerability reporting instructions.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE).
