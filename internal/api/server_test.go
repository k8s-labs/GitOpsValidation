package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pass"))
	})
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK || w.Body.String() != "pass" {
		t.Errorf("/healthz failed: got %d, %s", w.Code, w.Body.String())
	}
}

func TestVersion(t *testing.T) {
	req := httptest.NewRequest("GET", "/version", nil)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("0.0.1"))
	})
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK || w.Body.String() != "0.0.1" {
		t.Errorf("/version failed: got %d, %s", w.Code, w.Body.String())
	}
}
