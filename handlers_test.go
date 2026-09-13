package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	app := &App{}
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	app.healthHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestMasterKeyMiddlewareRejectsInvalidKey(t *testing.T) {
	app := &App{MasterKey: "expected-master-key"}
	called := false
	handler := app.masterKeyAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequest(http.MethodPost, "/admin/keys", nil)
	req.Header.Set("Authorization", "Bearer wrong-key")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}
	if called {
		t.Fatal("protected handler should not be called with an invalid master key")
	}
}

func TestHashAPIKeyIsDeterministic(t *testing.T) {
	first := hashAPIKey("tm_key_example")
	second := hashAPIKey("tm_key_example")

	if first != second {
		t.Fatal("expected hashAPIKey to be deterministic")
	}
	if len(first) != 64 {
		t.Fatalf("expected SHA-256 hex length 64, got %d", len(first))
	}
}
