package testutils

import (
	"context"
	"testing"
	"time"
)

// GetTimeoutContext returns a context that's already timed out
func GetTimeoutContext(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	t.Cleanup(func() { cancel() }) // Ensure cleanup
	time.Sleep(time.Millisecond)   // Force timeout
	return ctx
}

// GetCanceledContext returns an already canceled context
func GetCanceledContext(t *testing.T) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	return ctx
}

// GetContextWithTimeout returns a context with specified timeout
func GetContextWithTimeout(t *testing.T, duration time.Duration) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	t.Cleanup(func() { cancel() })
	return ctx
}
