# Build an application on MOTE

MOTE runs the execution loop. Your application handles authentication, UI,
persistence, scheduling, business tools, approvals, and output validation.
You do not need an HTTP worker service.

## Choose an integration

- [Python SDK](../sdk/python/README.md): Python 3.10+ with a standard-library runtime.
  Install from source and register Python functions as tools.
- [TypeScript / Node SDK](../sdk/typescript/README.md): Node 20+ with no runtime npm
  dependencies. Build and install the source package, then register JavaScript callbacks.
- [Go SDK](../sdk/go/README.md): Go 1.22+ with a standard-library-only runtime.
  Install from source and register Go host tool handlers.
- [Rust starter](../starters/rust-agent/README.md): embed the library, register a Rust
  handler, and validate its output file. Rust does not need the bridge.
- [CLI](QUICKSTART.md): run YAML-configured tasks without embedding a library.

Use the language guides for installation, runnable examples, API details, errors,
and tests. SDKs are not published to npm/PyPI and do not install or download the
bridge. The Rust package is `mote-agent`, imported as `mote`. Use a local path or
pinned Git revision for this source, not the unrelated `mote` crate.

## Get the bridge

Python, TypeScript and Go users can [download](https://github.com/danny-dis/MOTE/releases/latest)
the CLI/bridge archive without installing Rust. Follow the [quick start](QUICKSTART.md),
verify the checksum, and pass the absolute path to `mote-bridge` (`mote-bridge.exe`
on Windows). Downloads are public, unsigned, and not notarized.

To build both executables yourself, use stable Rust and its native build tools:

```sh
cargo build --locked --release --bins
```

The binaries are in `target/release`. Each SDK run uses one bridge subprocess,
launched without a command shell. Executable paths may contain spaces. All SDKs
support built-in capabilities without custom callbacks. Choose a reachable model
and use `auth_env` to name its credentials environment variable. See
[Configuration](CONFIGURATION.md).

## Application responsibilities

- Register each tool handler and grant its name in the manifest. Models cannot
  invent executable tools or grant permissions.
- Validate callback input fields yourself. Callbacks accept JSON objects and return
  strings; descriptions do not enforce JSON Schema.
- Check results before accepting outputs or publishing changes. Failed and timed-out
  runs may already have written files or invoked callbacks. Do not blindly retry them.
- Bound or isolate blocking callbacks and untrusted workloads. Callbacks retain
  host-user privileges; neither the SDK nor the bridge is a sandbox.
- Protect logs and observations that may contain workspace data. Callbacks are
  interactive, while runtime event lists arrive after completion.
- Match SDK source to the runtime release. SDK package versions and the
  [bridge protocol](BRIDGE_PROTOCOL.md) are independent of the crate version.

Read [Security](../SECURITY.md) for deployment boundaries and
[Contributing](../CONTRIBUTING.md#sdk-and-starter-checks) for verification commands.
CI runs real subprocesses and callbacks against scripted local model fixtures.
Check live endpoint availability and model-output correctness separately.
