package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// remoteExecutor forwards runs to a standalone execd service over HTTP. The
// service does the sandboxing; this side just adds the shared-secret header and
// enforces a client-side deadline a hair beyond the run's own timeout.
type remoteExecutor struct {
	url    string
	token  string
	limits Limits
	client *http.Client
}

func newRemote(cfg Config) *remoteExecutor {
	return &remoteExecutor{
		url:    strings.TrimRight(cfg.SandboxURL, "/") + "/run",
		token:  cfg.SandboxToken,
		limits: cfg.Limits,
		client: &http.Client{},
	}
}

func (e *remoteExecutor) Kind() string  { return "remote" }
func (e *remoteExecutor) Enabled() bool { return true }

func (e *remoteExecutor) Run(ctx context.Context, req Request) (Result, error) {
	req.TimeoutMs = clampTimeout(req.TimeoutMs, e.limits)
	body, err := json.Marshal(req)
	if err != nil {
		return Result{}, err
	}
	// Give the network round-trip a couple of seconds beyond the run budget.
	rctx, cancel := context.WithTimeout(ctx, time.Duration(req.TimeoutMs)*time.Millisecond+3*time.Second)
	defer cancel()

	hreq, err := http.NewRequestWithContext(rctx, http.MethodPost, e.url, bytes.NewReader(body))
	if err != nil {
		return Result{}, err
	}
	hreq.Header.Set("Content-Type", "application/json")
	if e.token != "" {
		hreq.Header.Set("Authorization", "Bearer "+e.token)
	}
	resp, err := e.client.Do(hreq)
	if err != nil {
		return Result{}, fmt.Errorf("executor service unreachable: %w", err)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(resp.Body, int64(e.limits.MaxOutputBytes)*2+4096))
	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("executor service returned %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	var res Result
	if err := json.Unmarshal(data, &res); err != nil {
		return Result{}, fmt.Errorf("bad response from executor service: %w", err)
	}
	return res, nil
}
