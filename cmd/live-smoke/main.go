package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	jobboss "github.com/Depth8064/go-jobboss2-api"
)

type checkFn func(ctx context.Context, c *jobboss.Client, q url.Values) (int, error)

type smokeCheck struct {
	name string
	fn   checkFn
}

type checkResult struct {
	name     string
	count    int
	duration time.Duration
	err      error
}

func main() {
	cfg, q, selected, timeout, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(2)
	}

	client := jobboss.NewClient(cfg)
	checks := availableChecks()
	fmt.Printf("Starting JobBOSS live smoke: %d checks, timeout=%s\n", len(selected), timeout)
	fmt.Printf("API URL: %s\n", client.BaseURL())
	fmt.Printf("Token URL: %s\n", client.TokenURL())

	authCtx, authCancel := context.WithTimeout(context.Background(), timeout)
	fmt.Print("Checking auth token      ... ")
	_, authErr := client.AccessToken(authCtx)
	authCancel()
	if authErr != nil {
		fmt.Printf("FAIL (%s)\n", timeout)
		fmt.Println("Auth preflight failed; skipping API checks.")
		fmt.Printf("Error: %v\n", authErr)
		os.Exit(1)
	}
	fmt.Println("OK")

	results := make([]checkResult, 0, len(selected))
	for _, checkName := range selected {
		fmt.Printf("Checking %-16s ... ", checkName)

		check, ok := checks[checkName]
		if !ok {
			fmt.Println("unknown check")
			results = append(results, checkResult{name: checkName, err: fmt.Errorf("unknown check")})
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		start := time.Now()
		count, runErr := check.fn(ctx, client, q)
		cancel()

		results = append(results, checkResult{
			name:     check.name,
			count:    count,
			duration: time.Since(start),
			err:      runErr,
		})

		if runErr != nil {
			fmt.Printf("FAIL (%s)\n", time.Since(start).Truncate(time.Millisecond))
			continue
		}
		fmt.Printf("OK (%s, records=%d)\n", time.Since(start).Truncate(time.Millisecond), count)
	}

	printSummary(results)
	if hasFailures(results) {
		os.Exit(1)
	}
}

func loadConfig() (jobboss.Config, url.Values, []string, time.Duration, error) {
	clientID := strings.TrimSpace(os.Getenv("JB2_CLIENT_ID"))
	clientSecret := strings.TrimSpace(os.Getenv("JB2_CLIENT_SECRET"))
	if clientID == "" || clientSecret == "" {
		return jobboss.Config{}, nil, nil, 0, errors.New("JB2_CLIENT_ID and JB2_CLIENT_SECRET are required")
	}

	timeoutSec, err := envInt("JB2_TIMEOUT_SEC", 20)
	if err != nil {
		return jobboss.Config{}, nil, nil, 0, fmt.Errorf("JB2_TIMEOUT_SEC: %w", err)
	}

	take, err := envInt("JB2_TAKE", 5)
	if err != nil {
		return jobboss.Config{}, nil, nil, 0, fmt.Errorf("JB2_TAKE: %w", err)
	}
	if take < 1 {
		take = 1
	}
	if take > jobboss.MaxTake {
		take = jobboss.MaxTake
	}

	selected := parseChecks(strings.TrimSpace(os.Getenv("JB2_CHECKS")))
	if len(selected) == 0 {
		selected = []string{"company", "currency-codes", "departments", "employees", "work-centers"}
	}

	cfg := jobboss.Config{
		BaseURL:      strings.TrimSpace(os.Getenv("JB2_BASE_URL")),
		TokenURL:     strings.TrimSpace(os.Getenv("JB2_TOKEN_URL")),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        strings.TrimSpace(os.Getenv("JB2_SCOPE")),
	}

	query := jobboss.NewQuery()
	query = jobboss.ApplyPaging(query, take, 0)

	return cfg, query, selected, time.Duration(timeoutSec) * time.Second, nil
}

func parseChecks(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	seen := map[string]bool{}
	for _, part := range parts {
		name := strings.ToLower(strings.TrimSpace(part))
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		out = append(out, name)
	}
	return out
}

func envInt(name string, fallback int) (int, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("must be integer")
	}
	return value, nil
}

func availableChecks() map[string]smokeCheck {
	return map[string]smokeCheck{
		"company": {
			name: "company",
			fn: func(ctx context.Context, c *jobboss.Client, _ url.Values) (int, error) {
				resp, err := c.GetCompany(ctx, nil)
				if err != nil {
					return 0, err
				}
				return len(resp.Data), nil
			},
		},
		"currency-codes": {
			name: "currency-codes",
			fn: func(ctx context.Context, c *jobboss.Client, q url.Values) (int, error) {
				resp, err := c.GetCurrencyCodes(ctx, q)
				if err != nil {
					return 0, err
				}
				return len(resp.Data), nil
			},
		},
		"departments": {
			name: "departments",
			fn: func(ctx context.Context, c *jobboss.Client, q url.Values) (int, error) {
				resp, err := c.GetDepartments(ctx, q)
				if err != nil {
					return 0, err
				}
				return len(resp.Data), nil
			},
		},
		"employees": {
			name: "employees",
			fn: func(ctx context.Context, c *jobboss.Client, q url.Values) (int, error) {
				resp, err := c.GetEmployees(ctx, q)
				if err != nil {
					return 0, err
				}
				return len(resp.Data), nil
			},
		},
		"work-centers": {
			name: "work-centers",
			fn: func(ctx context.Context, c *jobboss.Client, q url.Values) (int, error) {
				resp, err := c.GetWorkCenters(ctx, q)
				if err != nil {
					return 0, err
				}
				return len(resp.Data), nil
			},
		},
		"customers": {
			name: "customers",
			fn: func(ctx context.Context, c *jobboss.Client, q url.Values) (int, error) {
				resp, err := c.GetCustomers(ctx, q)
				if err != nil {
					return 0, err
				}
				return len(resp.Data), nil
			},
		},
		"vendors": {
			name: "vendors",
			fn: func(ctx context.Context, c *jobboss.Client, q url.Values) (int, error) {
				resp, err := c.GetVendors(ctx, q)
				if err != nil {
					return 0, err
				}
				return len(resp.Data), nil
			},
		},
		"tax-codes": {
			name: "tax-codes",
			fn: func(ctx context.Context, c *jobboss.Client, q url.Values) (int, error) {
				resp, err := c.GetTaxCodes(ctx, q)
				if err != nil {
					return 0, err
				}
				return len(resp.Data), nil
			},
		},
		"terms": {
			name: "terms",
			fn: func(ctx context.Context, c *jobboss.Client, q url.Values) (int, error) {
				resp, err := c.GetTerms(ctx, q)
				if err != nil {
					return 0, err
				}
				return len(resp.Data), nil
			},
		},
	}
}

func printSummary(results []checkResult) {
	fmt.Println("JobBOSS live smoke checks")
	fmt.Println("-------------------------")
	for _, result := range results {
		if result.err != nil {
			fmt.Printf("FAIL %-16s %8s  %s\n", result.name, result.duration.Truncate(time.Millisecond), result.err)
			continue
		}
		fmt.Printf("PASS %-16s %8s  records=%d\n", result.name, result.duration.Truncate(time.Millisecond), result.count)
	}
}

func hasFailures(results []checkResult) bool {
	for _, result := range results {
		if result.err != nil {
			return true
		}
	}
	return false
}
