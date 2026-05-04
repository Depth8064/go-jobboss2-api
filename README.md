# go-jobboss2-api

Minimal, dependency-free Go module for JobBOSS2 API access.

This module intentionally provides only a low-level authenticated HTTP client so higher-level services can be built by downstream code.

## Features

- OAuth2 client-credentials token handling
- Automatic token caching and refresh
- Automatic token refresh-and-retry on 401 (max 3 total attempts by default)
- Basic `GET`, `POST`, `PATCH`, `DELETE` helpers
- Generic JSON helpers for typed responses
- ECI-aligned defaults for API and auth hosts
- Query helpers for fields/filter/sort/take/skip
- UTC date/date-time formatting helpers
- Authentication API device endpoints (`/api/devices`, `/api/devices/{id}/activation`)
- No runtime dependency on OpenAPI parsing

## Install

```bash
go get go-jobboss2-api
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"net/url"

	jobboss "go-jobboss2-api"
)

type Customer struct {
	CustomerCode string `json:"customerCode"`
	CustomerName string `json:"customerName"`
}

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

	resp, err := jobboss.GetJSON[jobboss.ListResponse[Customer]](context.Background(), client, "/api/v1/customers", query)
	if err != nil {
		panic(err)
	}

	fmt.Printf("customers: %d\n", len(resp.Data))
}
```

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
day := jobboss.FormatUTCDate(time.Now())
ts := jobboss.FormatUTCDateTime(time.Now())
```
