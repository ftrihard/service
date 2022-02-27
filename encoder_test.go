package service1_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ftrihard/service1"
)

func TestEncoder(t *testing.T) {
	cfg := service1.CreateConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := service1.New(ctx, next, cfg, "demo-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	if req.Header.Get("User-Id") == "" {
		t.Errorf("JWT cannot be empty!")
	}
}
