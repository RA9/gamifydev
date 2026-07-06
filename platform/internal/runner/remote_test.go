package runner

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRemoteRoundTrip stands up a stub execd (a local executor behind HTTP) and
// checks the remote client marshals the request, carries the bearer token, and
// unmarshals the result — the whole main-app ⇄ execd wire path.
func TestRemoteRoundTrip(t *testing.T) {
	le := localExec(t) // skips if python3 is absent
	const token = "s3cret"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/run" || r.Header.Get("Authorization") != "Bearer "+token {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad", http.StatusBadRequest)
			return
		}
		res, _ := le.Run(r.Context(), req)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	}))
	defer srv.Close()

	re := newRemote(Config{SandboxURL: srv.URL, SandboxToken: token, Limits: DefaultLimits()})
	res, err := re.Run(context.Background(), Request{Code: "print(6 * 7)"})
	if err != nil {
		t.Fatalf("remote run: %v", err)
	}
	if strings.TrimSpace(res.Stdout) != "42" {
		t.Fatalf("stdout = %q, want 42", res.Stdout)
	}
}

func TestRemoteRejectsBadToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer right" {
			http.Error(w, "nope", http.StatusUnauthorized)
			return
		}
		_ = json.NewEncoder(w).Encode(Result{})
	}))
	defer srv.Close()

	re := newRemote(Config{SandboxURL: srv.URL, SandboxToken: "wrong", Limits: DefaultLimits()})
	if _, err := re.Run(context.Background(), Request{Code: "print(1)"}); err == nil {
		t.Fatal("expected an error on a rejected token")
	}
}
