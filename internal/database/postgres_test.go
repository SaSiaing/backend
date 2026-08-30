package database

import (
	"context"
	"errors"
	"testing"
)

func TestOpenRequiresDatabaseURL(t *testing.T) {
	t.Parallel()

	pool, err := Open(context.Background(), "  ")
	if pool != nil {
		pool.Close()
		t.Fatal("Open() returned a pool without DATABASE_URL")
	}
	if !errors.Is(err, ErrMissingURL) {
		t.Fatalf("Open() error = %v, want %v", err, ErrMissingURL)
	}
}
