package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/danny-dis/MOTE/sdk/go"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "example failed:", err)
		os.Exit(1)
	}
}

func run() error {
	bridge := envOr("MOTE_BRIDGE", "mote-bridge")
	model := envOr("MOTE_MODEL", "auto")
	endpoint := envOr("MOTE_MODEL_ENDPOINT", "http://127.0.0.1:47113/v1/chat/completions")

	modelConfig := mote.ModelConfig{
		Provider: "openai-compatible",
		Model:    model,
		Endpoint: endpoint,
	}
	if os.Getenv("MODEL_API_KEY") != "" {
		modelConfig.AuthEnv = "MODEL_API_KEY"
	}

	manifest := mote.Manifest{
		Name:         "go-example",
		Capabilities: []string{"uppercase"},
		Models:       []mote.ModelConfig{modelConfig},
	}

	var calls []map[string]any
	tools := map[string]mote.Tool{
		"uppercase": {
			Description: "Input: text string. Returns uppercase text.",
			Handler: func(ctx context.Context, input map[string]any) (string, error) {
				text, ok := input["text"].(string)
				if !ok {
					return "", errors.New("expected text")
				}
				calls = append(calls, input)
				return strings.ToUpper(text), nil
			},
		},
	}

	task := `Call uppercase with {"text":"hello"}, then complete with its returned text.`

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
	defer cancel()

	result, err := (&mote.Client{Binary: bridge}).Run(ctx, manifest, task, tools)
	if err != nil {
		return err
	}

	observed := false
	for _, event := range result.Events {
		if event["type"] != "ObservationReceived" {
			continue
		}
		data, ok := event["data"].(map[string]any)
		if !ok {
			continue
		}
		if data["text"] == "HELLO" {
			observed = true
			break
		}
	}

	if !result.Success() || len(calls) == 0 || !observed {
		return errors.New("independent callback/observation checks did not pass")
	}

	summary := map[string]any{
		"state":         result.State,
		"observation":   "HELLO",
		"callback_used": len(calls) > 0,
	}
	encoded, err := json.Marshal(summary)
	if err != nil {
		return err
	}
	fmt.Println(string(encoded))
	return nil
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
