# go-jobboss2-api

Minimal, dependency-free Go module for JobBOSS2 API access.

This module intentionally provides only a low-level authenticated HTTP client so higher-level services can be built by downstream code.

## Features

- OAuth2 client-credentials token handling
- Automatic token caching and refresh
- Basic `GET`, `POST`, `PATCH`, `DELETE` helpers
- Generic JSON helpers for typed responses
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
		TokenURL:     "https://your-auth-host/oauth2/token",
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
