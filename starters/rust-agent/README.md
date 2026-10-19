# MOTE Rust application starter

This starter embeds MOTE in a Rust application. The CLI loads an HTTP model
chain from `agent.yaml` and registers a Rust `uppercase` handler. The manifest
grants the handler permission to run. The application then asks the model to
write `report.txt` and checks both tool use and the file's exact contents.

The starter uses the repository's `mote-agent` path dependency (imported as
`mote`) and is not published to crates.io. You need a stable Rust toolchain.
It does not use `mote-bridge` or either language SDK.

## Run it

Run these commands from the repository root:

```bash
# First edit starters/rust-agent/agent.yaml for a reachable model endpoint and model.
cargo run --locked --manifest-path starters/rust-agent/Cargo.toml -- starters/rust-agent/agent.yaml
```

The YAML defaults to a local DMR-X-style endpoint and model `auto`. You supply
the provider or gateway. For an authenticated endpoint, set `auth_env` to the
name of an environment variable containing the credential; keep the credential
out of YAML. The endpoint must be reachable and return a supported completion.
See [CONFIGURATION](../../docs/CONFIGURATION.md) for manifest and model fields,
the [integration guide](../../docs/BUILDING_ON_MOTE.md) for application setup,
and the [quickstart](../../docs/QUICKSTART.md) for release and runtime details.

On success, the application prints `Verified report.txt: MOTE WORKS`. The file
is written under `starters/rust-agent/workspace/`; the workspace path is
relative to the YAML file. The checks require a call to `uppercase` and file
contents of exactly `MOTE WORKS`. A model's completion claim does not pass
those checks by itself.

## Adapt the starter

1. Copy this starter into your project and update its `mote-agent` path
   dependency. For a project outside this checkout, use a pinned Git revision
   of `https://github.com/danny-dis/MOTE` and package `mote-agent`. The
   crates.io package named `mote` is unrelated.
2. Replace `uppercase` with an operation your application owns. Validate its
   input and bound blocking work. The tool description guides the model; it
   does not validate input.
3. List only the tools the task needs in `capabilities`. Registering a handler
   does not grant permission to call it.
4. Replace the task and `run` acceptance checks with checks for your application.
5. Add the authentication, UI, persistence, scheduling, approvals, and workload
   isolation your application needs.

Host tools and the runtime run with your application's privileges. MOTE does
not provide an OS sandbox; isolate untrusted workloads yourself. Model
observations are sent to the configured model. Both observations and runtime
events may contain workspace data, so handle them as sensitive. The Rust API
returns its event list after the run, without live progress events.

## Deterministic checks

Run from the repository root:

```bash
cargo test --locked --manifest-path starters/rust-agent/Cargo.toml
cargo fmt --check --manifest-path starters/rust-agent/Cargo.toml
cargo clippy --locked --manifest-path starters/rust-agent/Cargo.toml --all-targets -- -D warnings
```

The tests use scripted `Model` implementations and need no API keys. They
check accepted output, rejected output, and a model that skips the required
callback. They do not test provider availability or live inference.
