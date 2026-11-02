package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	jobboss "github.com/Depth8064/go-jobboss2-api"
)

func setConfigEnv(t *testing.T) {
	t.Helper()
	for _, name := range []string{
		"JB2_BASE_URL", "JB2_TOKEN_URL", "JB2_SCOPE", "JB2_TIMEOUT_SEC",
		"JB2_TAKE", "JB2_CHECKS",
	} {
		t.Setenv(name, "")
	}
	t.Setenv("JB2_CLIENT_ID", "test-client")
	t.Setenv("JB2_CLIENT_SECRET", "test-secret")
}

func captureOutput(t *testing.T, fn func()) (string, string) {
	t.Helper()
	stdout, err := os.CreateTemp(t.TempDir(), "stdout")
	if err != nil {
		t.Fatal(err)
	}
	defer stdout.Close()
	stderr, err := os.CreateTemp(t.TempDir(), "stderr")
	if err != nil {
		t.Fatal(err)
	}
	defer stderr.Close()
	oldOut, oldErr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = stdout, stderr
	defer func() { os.Stdout, os.Stderr = oldOut, oldErr }()
	fn()
	read := func(f *os.File) string {
		t.Helper()
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			t.Fatal(err)
		}
		b, err := io.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	return read(stdout), read(stderr)
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		wantErr  string
		take     string
		timeout  time.Duration
		selected []string
	}{
		{
			name: "defaults", take: "5", timeout: 20 * time.Second,
			selected: []string{"company", "currency-codes", "departments", "employees", "work-centers"},
		},
		{name: "missing client ID", env: map[string]string{"JB2_CLIENT_ID": " "}, wantErr: "JB2_CLIENT_ID and JB2_CLIENT_SECRET are required"},
		{name: "missing client secret", env: map[string]string{"JB2_CLIENT_SECRET": ""}, wantErr: "JB2_CLIENT_ID and JB2_CLIENT_SECRET are required"},
		{name: "invalid timeout", env: map[string]string{"JB2_TIMEOUT_SEC": "invalid"}, wantErr: "JB2_TIMEOUT_SEC: must be integer"},
		{name: "invalid take", env: map[string]string{"JB2_TAKE": "1.5"}, wantErr: "JB2_TAKE: must be integer"},
		{
			name: "minimum take", env: map[string]string{"JB2_TAKE": "-10"},
			take: "1", timeout: 20 * time.Second,
			selected: []string{"company", "currency-codes", "departments", "employees", "work-centers"},
		},
		{
			name: "maximum take", env: map[string]string{"JB2_TAKE": "999999"},
			take: "200", timeout: 20 * time.Second,
			selected: []string{"company", "currency-codes", "departments", "employees", "work-centers"},
		},
		{
			name: "custom settings",
			env: map[string]string{
				"JB2_BASE_URL": " https://example.test/api ", "JB2_TOKEN_URL": " https://example.test/token ",
				"JB2_CLIENT_ID": " custom-id ", "JB2_CLIENT_SECRET": " custom-secret ",
				"JB2_SCOPE": " custom-scope ", "JB2_TIMEOUT_SEC": " 3 ", "JB2_TAKE": " 12 ",
				"JB2_CHECKS": " Customers, vendors, CUSTOMERS, , terms ",
			},
			take: "12", timeout: 3 * time.Second, selected: []string{"customers", "vendors", "terms"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setConfigEnv(t)
			for key, value := range tt.env {
				t.Setenv(key, value)
			}
			cfg, q, selected, timeout, err := loadConfig()
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("error = %v, want %q", err, tt.wantErr)
				}
				if cfg != (jobboss.Config{}) || q != nil || selected != nil || timeout != 0 {
					t.Fatalf("configuration error returned partial settings: %+v, %v, %v, %v", cfg, q, selected, timeout)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if want := (url.Values{"take": {tt.take}}); !reflect.DeepEqual(q, want) {
				t.Errorf("paging = %v, want %v", q, want)
			}
			if timeout != tt.timeout || !reflect.DeepEqual(selected, tt.selected) {
				t.Errorf("timeout/checks = %v/%v, want %v/%v", timeout, selected, tt.timeout, tt.selected)
			}
			if tt.name == "custom settings" {
				want := jobboss.Config{
					BaseURL: "https://example.test/api", TokenURL: "https://example.test/token",
					ClientID: "custom-id", ClientSecret: "custom-secret", Scope: "custom-scope",
				}
				if cfg != want {
					t.Errorf("config = %+v, want %+v", cfg, want)
				}
			} else if cfg.ClientID != "test-client" || cfg.ClientSecret != "test-secret" ||
				cfg.BaseURL != "" || cfg.TokenURL != "" || cfg.Scope != "" {
				t.Errorf("unexpected default config: %+v", cfg)
			}
		})
	}
}

func TestParseChecks(t *testing.T) {
	for _, tt := range []struct {
		raw  string
		want []string
	}{
		{"", nil},
		{" , , ", []string{}},
		{" COMPANY, customers,company,, Vendors ", []string{"company", "customers", "vendors"}},
		{"unknown", []string{"unknown"}},
	} {
		if got := parseChecks(tt.raw); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parseChecks(%q) = %v, want %v", tt.raw, got, tt.want)
		}
	}
}

func TestEnvInt(t *testing.T) {
	for _, tt := range []struct {
		raw     string
		want    int
		wantErr bool
	}{
		{"", 7, false}, {"  ", 7, false}, {" 12 ", 12, false},
		{"0", 0, false}, {"-2", -2, false}, {"1.5", 0, true},
	} {
		t.Setenv("JB2_TEST_INT", tt.raw)
		got, err := envInt("JB2_TEST_INT", 7)
		if got != tt.want || (err != nil) != tt.wantErr {
			t.Errorf("envInt(%q) = %d, %v; want %d, error=%v", tt.raw, got, err, tt.want, tt.wantErr)
		}
	}
}

func TestAvailableChecks(t *testing.T) {
	paths := map[string]string{
		"company": "/api/v1/company", "currency-codes": "/api/v1/currency-codes",
		"departments": "/api/v1/departments", "employees": "/api/v1/employees",
		"work-centers": "/api/v1/work-centers", "customers": "/api/v1/customers",
		"vendors": "/api/v1/vendors", "tax-codes": "/api/v1/tax-codes", "terms": "/api/v1/terms",
	}
	checks := availableChecks()
	if len(checks) != len(paths) {
		t.Fatalf("got %d checks, want %d", len(checks), len(paths))
	}
	for name, path := range paths {
		t.Run(name, func(t *testing.T) {
			check, ok := checks[name]
			if !ok || check.name != name {
				t.Fatalf("missing or misnamed check %q: %+v", name, check)
			}
			for _, status := range []int{http.StatusOK, http.StatusBadRequest} {
				requests := 0
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/token" {
						if r.Method != http.MethodPost {
							t.Errorf("token method = %s", r.Method)
						}
						if err := r.ParseForm(); err != nil {
							t.Error(err)
						}
						if r.Form.Get("grant_type") != "client_credentials" || r.Form.Get("client_id") != "id" ||
							r.Form.Get("client_secret") != "secret" || r.Form.Get("scope") != "openid" {
							t.Errorf("unexpected token form: %v", r.Form)
						}
						io.WriteString(w, `{"access_token":"test-token","expires_in":3600}`)
						return
					}
					requests++
					if r.URL.Path != path || r.Method != http.MethodGet {
						t.Errorf("request = %s %s, want GET %s", r.Method, r.URL.Path, path)
					}
					if r.Header.Get("Authorization") != "Bearer test-token" {
						t.Errorf("missing bearer token: %v", r.Header)
					}
					wantQuery := url.Values{"take": {"2"}, "skip": {"0"}}
					if name == "company" {
						wantQuery = url.Values{}
					}
					if !reflect.DeepEqual(r.URL.Query(), wantQuery) {
						t.Errorf("query = %v, want %v", r.URL.Query(), wantQuery)
					}
					w.WriteHeader(status)
					if status == http.StatusOK {
						io.WriteString(w, `{"data":[{},{}]}`)
					} else {
						io.WriteString(w, "invalid query")
					}
				}))
				client := jobboss.NewClient(jobboss.Config{
					BaseURL: server.URL, TokenURL: server.URL + "/token",
					ClientID: "id", ClientSecret: "secret", HTTPClient: server.Client(),
				})
				count, err := check.fn(context.Background(), client, url.Values{"take": {"2"}, "skip": {"0"}})
				server.Close()
				if requests != 1 {
					t.Errorf("API requests = %d, want 1", requests)
				}
				if status == http.StatusOK {
					if err != nil || count != 2 {
						t.Errorf("check = %d, %v; want 2 records", count, err)
					}
				} else {
					var apiErr *jobboss.APIError
					if count != 0 || !errors.As(err, &apiErr) || apiErr.StatusCode != status || apiErr.Body != "invalid query" {
						t.Errorf("check = %d, %v; want original API error", count, err)
					}
				}
			}
		})
	}
}

func TestMainSuccess(t *testing.T) {
	setConfigEnv(t)
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		if r.URL.Path == "/token" {
			io.WriteString(w, `{"access_token":"test-token","expires_in":3600}`)
			return
		}
		io.WriteString(w, `{"data":[{},{}]}`)
	}))
	defer server.Close()
	t.Setenv("JB2_BASE_URL", server.URL)
	t.Setenv("JB2_TOKEN_URL", server.URL+"/token")
	stdout, stderr := captureOutput(t, main)
	server.Close()
	if stderr != "" {
		t.Errorf("stderr = %q", stderr)
	}
	for _, want := range []string{
		"Starting JobBOSS live smoke: 5 checks, timeout=20s",
		"API URL: " + server.URL, "Token URL: " + server.URL + "/token",
		"Checking auth token      ... OK", "JobBOSS live smoke checks",
		"PASS company", "PASS currency-codes", "PASS departments", "PASS employees", "PASS work-centers",
		"records=2",
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("stdout missing %q:\n%s", want, stdout)
		}
	}
	wantPaths := []string{"/token", "/api/v1/company", "/api/v1/currency-codes", "/api/v1/departments", "/api/v1/employees", "/api/v1/work-centers"}
	if !reflect.DeepEqual(paths, wantPaths) {
		t.Errorf("requests = %v, want %v", paths, wantPaths)
	}
	if strings.Contains(stdout, "FAIL") {
		t.Errorf("unexpected failure:\n%s", stdout)
	}
}

func TestRunFailures(t *testing.T) {
	for _, tt := range []struct {
		name       string
		checks     string
		configFail bool
		authFail   bool
		apiFail    bool
		wantCode   int
		wantPaths  []string
		wantOutput []string
	}{
		{name: "configuration", configFail: true, wantCode: 2},
		{
			name: "authentication", authFail: true, wantCode: 1, wantPaths: []string{"/token"},
			wantOutput: []string{"FAIL (20s)", "Auth preflight failed; skipping API checks.", "Error: jobboss API error 400"},
		},
		{
			name: "partial failure", checks: "company,customers,terms", apiFail: true, wantCode: 1,
			wantPaths:  []string{"/token", "/api/v1/company", "/api/v1/customers", "/api/v1/terms"},
			wantOutput: []string{"PASS company", "FAIL customers", "PASS terms", "invalid customer query"},
		},
		{
			name: "unknown check", checks: "unknown,company", wantCode: 1,
			wantPaths:  []string{"/token", "/api/v1/company"},
			wantOutput: []string{"unknown check", "FAIL unknown", "PASS company"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			setConfigEnv(t)
			var paths []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				paths = append(paths, r.URL.Path)
				if r.URL.Path == "/token" {
					if tt.authFail {
						http.Error(w, "invalid credentials", http.StatusBadRequest)
						return
					}
					io.WriteString(w, `{"access_token":"test-token","expires_in":3600}`)
					return
				}
				if tt.apiFail && r.URL.Path == "/api/v1/customers" {
					http.Error(w, "invalid customer query", http.StatusBadRequest)
					return
				}
				io.WriteString(w, `{"data":[{}]}`)
			}))
			defer server.Close()
			t.Setenv("JB2_BASE_URL", server.URL)
			t.Setenv("JB2_TOKEN_URL", server.URL+"/token")
			t.Setenv("JB2_CHECKS", tt.checks)
			if tt.configFail {
				t.Setenv("JB2_CLIENT_ID", "")
			}
			var codes []int
			stdout, stderr := captureOutput(t, func() {
				run(func(code int) { codes = append(codes, code) })
			})
			server.Close()
			if !reflect.DeepEqual(codes, []int{tt.wantCode}) {
				t.Errorf("exit codes = %v, want [%d]", codes, tt.wantCode)
			}
			if !reflect.DeepEqual(paths, tt.wantPaths) {
				t.Errorf("requests = %v, want %v", paths, tt.wantPaths)
			}
			if tt.configFail {
				if stdout != "" || stderr != "config error: JB2_CLIENT_ID and JB2_CLIENT_SECRET are required\n" {
					t.Errorf("unexpected config output: stdout=%q stderr=%q", stdout, stderr)
				}
			} else if stderr != "" {
				t.Errorf("stderr = %q", stderr)
			}
			for _, want := range tt.wantOutput {
				if !strings.Contains(stdout, want) {
					t.Errorf("stdout missing %q:\n%s", want, stdout)
				}
			}
		})
	}
}

func TestSummaryAndFailures(t *testing.T) {
	results := []checkResult{
		{name: "company", count: 2, duration: 1234 * time.Microsecond},
		{name: "customers", duration: 2345 * time.Microsecond, err: errors.New("request failed")},
	}
	stdout, stderr := captureOutput(t, func() { printSummary(results) })
	want := "JobBOSS live smoke checks\n-------------------------\n" +
		"PASS company               1ms  records=2\n" +
		"FAIL customers             2ms  request failed\n"
	if stdout != want || stderr != "" {
		t.Errorf("summary = %q, stderr=%q; want %q", stdout, stderr, want)
	}
	if hasFailures(nil) || hasFailures(results[:1]) || !hasFailures(results) {
		t.Error("hasFailures did not distinguish empty, successful, and mixed results")
	}
}
