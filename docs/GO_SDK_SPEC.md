# Official Go SDK specification — MOTE v0.14.1

## Status and scope

Accepted release scope for the next MOTE release, **v0.14.1**. Adding this
specification does not publish that release or change the bridge protocol.
The Go SDK is an official, repository-maintained integration alongside Python
and TypeScript, not a separate agent runtime.

- Source module: `sdk/go`, imported as `github.com/danny-dis/MOTE/sdk/go`.
- Minimum Go version: 1.22; runtime dependencies: Go standard library only.
- Transport: [bridge protocol v1](BRIDGE_PROTOCOL.md).
- Distribution: source checkout/local module replacement. No module-proxy
  publication or standalone SDK version tag is promised by this release.
- Runtime: use the matching MOTE release's `mote-bridge` executable. Go SDK source
  is first included in v0.14.1; protocol-v1 interoperability with v0.14.0 does
  not retroactively add SDK source to that release.

## Public API contract

`Client.Run(ctx, manifest, task, tools)` starts a fresh trusted executable per
run, without a command shell. `Client.Binary` defaults to `mote-bridge`;
`Args` supplies an argument array, `Cwd` sets the child working directory, and
`Env` overrides the inherited environment. `Timeout == 0` selects 150 seconds;
a negative timeout is a setup error. The earlier caller-context deadline applies.
A nil context uses `context.Background()`.

`Manifest`, `ModelConfig`, `ShellConfig` and `DecisionConfig` serialize the
runtime configuration; optional pointer fields preserve explicit zero values.
The Rust bridge remains the authority for manifest validation and defaults.
Credentials are referenced by environment-variable name (`AuthEnv`), not sent
as tool metadata.

`Tool` contains a description and a handler accepting a context and JSON object,
returning `(string, error)`. A tool must be both registered and named in
`Manifest.Capabilities` before its handler is invoked. Missing descriptions or
handlers are setup errors. Handler errors and panics produce a generic tool
failure without exposing their text or inputs to the model. Applications must
validate domain fields themselves. Neither runs nor callbacks are retried.

`RunResult` exposes `State`, `Events` and `ReturnCode`. `Success()` is true only
for state `completed` with exit code 0. Failed/cancelled runtime work returns a
result with exit code 1 rather than a successful result. Setup, malformed-protocol,
transport and deadline/cancellation errors return no result and a `*MoteError`.
Callers may inspect `Code` or use `errors.Is` with `ErrSetup`, `ErrProtocol`,
`ErrTransport` and `ErrTimeout`. Context cancellation/deadline errors retain their
context cause for `errors.Is`.

## Protocol and lifecycle requirements

- LF-terminated UTF-8 JSON objects; at most 1,048,576 bytes per incoming or
  outgoing frame, including the LF. Malformed, oversized, truncated, unexpected
  or blank frames fail closed.
- Sequential tool-call IDs start at 1. Tool input must be an object. Only one
  terminal result is accepted; output after it is rejected. Terminal state and
  child exit status must agree.
- Stdout buffering is bounded and applies backpressure. Stderr is diagnostic
  only, drained without exposing it as an API result.
- A timeout/cancellation stops and reaps the direct bridge process. Blocked pipe
  writes must not bypass the deadline or be misreported as ordinary transport
  failures when the context has expired.
- Host handlers execute synchronously and receive the run context. They must
  cooperate with cancellation for `Run` to return promptly. Go cannot forcibly
  stop arbitrary host code; the subprocess watchdog is not a hard deadline for
  an uncooperative callback. Applications must isolate such callbacks.
- Killing the direct bridge does not guarantee termination of all descendants.
  Neither capability grants nor the SDK are an OS sandbox. Failed/timed-out runs
  may already have side effects; events may contain sensitive workspace data.

## v0.14.1 acceptance checks

The repository CI gate must run the following on Linux, Windows and macOS, using
Go 1.22 for the minimum-version contract and a real built bridge for E2E coverage:

```sh
test -z "$(gofmt -l sdk/go)"
go -C sdk/go build ./...
go -C sdk/go vet ./...
go -C sdk/go test ./...
```

Tests must cover callbacks and grant denial; generic error/panic redaction;
invalid/trailing/oversized protocol frames; timeout/cancellation and blocked
writes; state/exit agreement; and a real bridge subprocess with a local scripted
model endpoint. The extracted-release bridge must also pass the Go suite in the
release workflow, alongside the Python and TypeScript suites. These tests do not
establish live-provider availability or model-output correctness.

The full existing Rust, SDK, starter, documentation, packaging and dependency
audit gates remain required. Keep `Cargo.toml` and published tags unchanged until
the maintainer performs the normal [release procedure](RELEASING.md).

See the [Go SDK guide](../sdk/go/README.md) for installation and examples.
