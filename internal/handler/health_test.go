package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

type stubDatabase struct {
	err error
}

func (database stubDatabase) Ping(context.Context) error {
	return database.err
}

func TestHealth(t *testing.T) {
	t.Parallel()

	e := echo.New()
	e.GET("/health", Health(stubDatabase{}))

	recorder := httptest.NewRecorder()
	e.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got, want := recorder.Body.String(), "{\"status\":\"ok\"}\n"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestHealthWhenDatabaseIsUnavailable(t *testing.T) {
	t.Parallel()

	e := echo.New()
	e.GET("/health", Health(stubDatabase{err: errors.New("connection refused")}))

	recorder := httptest.NewRecorder()
	e.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health", nil))

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusServiceUnavailable)
	}
	if got, want := recorder.Body.String(), "{\"message\":\"database unavailable\",\"code\":\"service_unavailable\"}\n"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}
