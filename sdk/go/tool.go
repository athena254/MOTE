package mote

import "context"

// Tool is a host callback registered for a run. Registration is not permission:
// the tool name must also appear in Manifest.Capabilities to be granted.
type Tool struct {
	Description string
	Handler     func(ctx context.Context, input map[string]any) (string, error)
}
