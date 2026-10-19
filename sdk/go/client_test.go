package mote_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	mote "github.com/danny-dis/MOTE/sdk/go"
)

func TestSuccessCallbackAndEvents(t *testing.T) {
	var got string
	tools := map[string]mote.Tool{
		"x": {
			Description: "x",
			Handler: func(_ context.Context, input map[string]any) (string, error) {
				got, _ = input["value"].(string)
				return strings.ToUpper(got), nil
			},
		},
	}

	result, err := runMode(t, "success", tools, 5*time.Second)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if !result.Success() {
		t.Fatalf("expected success, got state=%q code=%d", result.State, result.ReturnCode)
	}
	if got != "hello" {
		t.Fatalf("handler input = %q, want %q", got, "hello")
	}
	if len(result.Events) != 1 {
		t.Fatalf("events = %v, want exactly one event", result.Events)
	}
	if result.Events[0]["type"] != "hello" {
		t.Fatalf("event = %v, want type=hello", result.Events[0])
	}
}

func TestProtocolErrors(t *testing.T) {
	cases := []struct {
		name string
		mode string
	}{
		{"malformed", "malformed"},
		{"wrongid", "wrongid"},
		{"oversized", "oversized"},
		{"invalidutf8", "invalidutf8"},
		{"duplicate", "duplicate"},
		{"mismatch", "mismatch"},
		{"missingresult", "missingresult"},
		{"blank", "blank"},
		{"truncated", "truncated"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := runMode(t, tc.mode, upperTool(), 5*time.Second)
			if err == nil {
				t.Fatalf("expected error, got result %+v", result)
			}
			if !errors.Is(err, mote.ErrProtocol) {
				t.Fatalf("error = %v, want ErrProtocol", err)
			}
		})
	}
}

func TestCallbackErrorRedacted(t *testing.T) {
	secret := "top-secret-input"
	tools := map[string]mote.Tool{
		"x": {
			Description: "x",
			Handler: func(_ context.Context, _ map[string]any) (string, error) {
				return "", errors.New(secret)
			},
		},
	}

	result, err := runMode(t, "success", tools, 5*time.Second)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if !result.Success() {
		t.Fatalf("expected success after redacted callback error, got %+v", result)
	}
	for _, event := range result.Events {
		if strings.Contains(fmt.Sprint(event), secret) {
			t.Fatalf("callback error leaked into events: %v", event)
		}
	}
}

func TestTimeoutRejects(t *testing.T) {
	started := time.Now()
	result, err := runMode(t, "hang", upperTool(), 300*time.Millisecond)
	if err == nil {
		t.Fatalf("expected timeout, got result %+v", result)
	}
	if !errors.Is(err, mote.ErrTimeout) {
		t.Fatalf("error = %v, want ErrTimeout", err)
	}
	if elapsed := time.Since(started); elapsed > 3*time.Second {
		t.Fatalf("timeout not enforced promptly: %v", elapsed)
	}
}
