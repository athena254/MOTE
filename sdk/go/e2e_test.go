package mote_test

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	mote "github.com/danny-dis/MOTE/sdk/go"
)

// TestE2EBridge exercises the real bridge. It is skipped unless MOTE_BRIDGE
// names an executable bridge binary.
func TestE2EBridge(t *testing.T) {
	bridge := os.Getenv("MOTE_BRIDGE")
	if bridge == "" {
		t.Skip("MOTE_BRIDGE not set")
	}

	var calls int64
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		_ = r.Body.Close()

		content := `{"action":"complete","summary":"HELLO"}`
		if atomic.AddInt64(&calls, 1) == 1 {
			content = `{"action":"custom","name":"uppercase","input":{"text":"hello"},"summary":"uppercase hello"}`
		}
		payload, err := json.Marshal(map[string]any{
			"choices": []any{
				map[string]any{
					"message": map[string]any{
						"role":    "assistant",
						"content": content,
					},
				},
			},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(payload)
	})

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	server := &http.Server{Handler: handler}
	go func() { _ = server.Serve(listener) }()
	defer func() { _ = server.Close() }()

	manifest := mote.Manifest{
		Name:         "go-e2e",
		Workspace:    ".",
		Capabilities: []string{"uppercase"},
		Models: []mote.ModelConfig{
			{
				Provider: "openai-compatible",
				Model:    "fixture",
				Endpoint: "http://" + listener.Addr().String() + "/v1/chat/completions",
			},
		},
	}
	tools := map[string]mote.Tool{
		"uppercase": {
			Description: "Returns uppercase text.",
			Handler: func(_ context.Context, input map[string]any) (string, error) {
				text, _ := input["text"].(string)
				return strings.ToUpper(text), nil
			},
		},
	}

	client := &mote.Client{Binary: bridge, Timeout: 15 * time.Second}
	result, err := client.Run(context.Background(), manifest, "Use uppercase with input hello, then finish.", tools)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if !result.Success() {
		t.Fatalf("expected success, got %+v", result)
	}

	found := false
	for _, event := range result.Events {
		raw, err := json.Marshal(event)
		if err != nil {
			continue
		}
		if strings.Contains(string(raw), "HELLO") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("no event contained HELLO: %v", result.Events)
	}
}
