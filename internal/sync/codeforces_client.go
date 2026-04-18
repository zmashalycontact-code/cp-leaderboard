package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Interface (port)
// ---------------------------------------------------------------------------

// CodeforcesClient is the anti-corruption layer between our sync engine and
// the external Codeforces API. Defining it as an interface lets us:
//   - Inject a mock in unit tests (no real HTTP calls, no rate-limit worries).
//   - Swap the implementation (e.g. add circuit breaker wrapping) without
//     touching the worker pool logic.
//
// The method is deliberately batch-oriented: callers hand it a slice of
// handles and receive back a slice of results. The implementation decides
// how to map those to actual HTTP requests.
type CodeforcesClient interface {
	// FetchUsers retrieves Codeforces profile data for the given handles.
	// Implementations may issue one or more HTTP requests (e.g. to honour
	// the per-request handle cap). The returned slice preserves no guaranteed
	// ordering relative to the input handles.
	//
	// A non-nil error is returned on any network failure, non-200 status, or
	// when the Codeforces API returns status:"FAILED". Partial results are
	// NOT returned on error — the caller should treat the batch as failed and
	// schedule a retry.
	FetchUsers(ctx context.Context, handles []string) ([]cfUser, error)
}

// ---------------------------------------------------------------------------
// HTTP implementation
// ---------------------------------------------------------------------------

const (
	// cfBaseURL is the Codeforces REST API base.
	cfBaseURL = "https://codeforces.com/api"

	// cfHandlesPerRequest is the documented maximum number of handles that
	// /user.info accepts in a single call. We stay safely under it at 50.
	cfHandlesPerRequest = 50

	// defaultHTTPTimeout caps individual HTTP round-trips.  The Codeforces
	// API typically responds in <1 s; 10 s gives plenty of slack for transient
	// slowness without hanging the worker indefinitely.
	defaultHTTPTimeout = 10 * time.Second
)

// HTTPCodeforcesClient is the production implementation of CodeforcesClient.
// It uses a plain net/http.Client (injected so callers can supply custom
// transports, e.g. with proxy settings or retry middleware).
type HTTPCodeforcesClient struct {
	httpClient *http.Client
	baseURL    string // overridable in tests to point at a mock server
}

// NewHTTPCodeforcesClient constructs a production-ready CodeforcesClient.
// Pass nil for httpClient to use a sensible default.
func NewHTTPCodeforcesClient(httpClient *http.Client) *HTTPCodeforcesClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultHTTPTimeout}
	}
	return &HTTPCodeforcesClient{
		httpClient: httpClient,
		baseURL:    cfBaseURL,
	}
}

// FetchUsers implements CodeforcesClient.
//
// If len(handles) > cfHandlesPerRequest the slice is automatically split into
// batches and each batch is fetched sequentially within this call.  The worker
// pool already rate-limits how often FetchUsers is invoked, so sequential
// batching here is intentional — we don't want nested concurrency making the
// rate-limit maths harder to reason about.
func (c *HTTPCodeforcesClient) FetchUsers(ctx context.Context, handles []string) ([]cfUser, error) {
	if len(handles) == 0 {
		return nil, nil
	}

	var all []cfUser

	for i := 0; i < len(handles); i += cfHandlesPerRequest {
		end := i + cfHandlesPerRequest
		if end > len(handles) {
			end = len(handles)
		}
		batch := handles[i:end]

		users, err := c.fetchBatch(ctx, batch)
		if err != nil {
			return nil, err // fail fast; caller will retry the whole job
		}
		all = append(all, users...)
	}

	return all, nil
}

// fetchBatch performs a single HTTP GET for up to cfHandlesPerRequest handles.
func (c *HTTPCodeforcesClient) fetchBatch(ctx context.Context, handles []string) ([]cfUser, error) {
	// Build the URL.  Codeforces expects handles joined by semicolons.
	// We use url.QueryEscape on the whole semicolon-joined string so that
	// special characters in handles (unlikely but possible) are escaped.
	endpoint := fmt.Sprintf(
		"%s/user.info?handles=%s",
		c.baseURL,
		url.QueryEscape(strings.Join(handles, ";")),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("codeforces: build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("codeforces: http GET: %w", err)
	}
	defer resp.Body.Close()

	// Read at most 1 MiB to guard against absurdly large or malicious bodies.
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("codeforces: read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("codeforces: unexpected status %d: %s", resp.StatusCode, body)
	}

	var envelope cfResponse
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("codeforces: decode response: %w", err)
	}

	if envelope.Status != "OK" {
		return nil, fmt.Errorf("codeforces: API error: %s", envelope.Comment)
	}

	return envelope.Result, nil
}
