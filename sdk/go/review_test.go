package mote_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	mote "github.com/danny-dis/MOTE/sdk/go"
)

func TestPanicIsRedactedOnWire(t *testing.T) {
	tools := map[string]mote.Tool{"x": {Description: "x", Handler: func(context.Context, map[string]any) (string, error) { panic("secret panic input") }}}
	result, err := runMode(t, "success", tools, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	reply, ok := result.Events[0]["reply"].(map[string]any)
	if !ok || reply["error"] != "tool failed" || len(reply) != 3 {
		t.Fatalf("panic leaked or invalid reply: %v", reply)
	}
}

func TestUncooperativeHandlerRequiresHostRelease(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	entered, release, handlerDone := make(chan struct{}), make(chan struct{}), make(chan struct{})
	defer close(release)
	tools := map[string]mote.Tool{"x": {Description: "x", Handler: func(handlerCtx context.Context, _ map[string]any) (string, error) {
		defer close(handlerDone)
		close(entered)
		<-release
		return "", handlerCtx.Err()
	}}}
	outcome := make(chan error, 1)
	go func() {
		client := &mote.Client{Binary: fixtureBin, Timeout: 10 * time.Second, Env: []string{"TEST_BRIDGE_MODE=success"}}
		_, err := client.Run(ctx, fixtureManifest("x"), "go", tools)
		outcome <- err
	}()
	select {
	case <-entered:
	case err := <-outcome:
		t.Fatalf("run ended before callback: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("callback not reached")
	}
	cancel()
	select {
	case <-handlerDone:
		t.Fatal("handler stopped without host release")
	case err := <-outcome:
		t.Fatalf("Run returned before synchronous callback finished: %v", err)
	case <-time.After(50 * time.Millisecond):
	}
	// Release the callback without leaking it; the API can only return afterwards.
	release <- struct{}{}
	select {
	case err := <-outcome:
		if !errors.Is(err, mote.ErrTimeout) || !errors.Is(err, context.Canceled) {
			t.Fatalf("expected cancellation, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after host release")
	}
}

func TestBlockedStartupWriteReturnsTimeout(t *testing.T) {
	client := &mote.Client{Binary: fixtureBin, Timeout: 300 * time.Millisecond, Env: []string{"TEST_BRIDGE_MODE=noread"}}
	_, err := client.Run(context.Background(), fixtureManifest(), strings.Repeat("x", 512*1024), nil)
	if !errors.Is(err, mote.ErrTimeout) {
		t.Fatalf("error = %v, want ErrTimeout", err)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("missing deadline cause: %v", err)
	}
}

func TestCancellationDuringCallbackReturnsTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := &mote.Client{Binary: fixtureBin, Timeout: 5 * time.Second, Env: []string{"TEST_BRIDGE_MODE=success"}}
	tools := map[string]mote.Tool{"x": {Description: "x", Handler: func(context.Context, map[string]any) (string, error) { cancel(); return "", ctx.Err() }}}
	_, err := client.Run(ctx, fixtureManifest("x"), "go", tools)
	if !errors.Is(err, mote.ErrTimeout) || !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want timeout wrapping cancellation", err)
	}
}
