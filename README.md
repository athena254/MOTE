# MOTE

A small execution engine for AI agents.

[![CI](https://github.com/danny-dis/MOTE/actions/workflows/ci.yml/badge.svg?branch=production-ready)](https://github.com/danny-dis/MOTE/actions/workflows/ci.yml)
[![Release](https://github.com/danny-dis/MOTE/actions/workflows/release.yml/badge.svg)](https://github.com/danny-dis/MOTE/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

MOTE (Minimal Orchestration & Task Execution) is a Rust library and command-line
program. It asks a model for the next action, checks capability grants, executes
the tool, and passes the result back. The loop stops when the task finishes or a
limit is reached. Your application handles users, approvals, storage, scheduling,
the interface, and output validation.

[Website](https://mote-runtime.vercel.app) ·
[Download](https://github.com/danny-dis/MOTE/releases/latest) ·
[Quick start](docs/QUICKSTART.md) · [Documentation](#documentation) · [Changelog](CHANGELOG.md)

> Source and downloads are public. You supply the model and credentials.
> Capability checks are not an OS sandbox. Use trusted workspaces and operators.

## Download and run

Choose an archive from [Releases](https://github.com/danny-dis/MOTE/releases/latest):

- Windows x64: ZIP containing `mote.exe` and `mote-bridge.exe`.
- Linux x64: tar.gz for an Ubuntu 22.04-compatible system with OpenSSL 3 and CA certificates.
- macOS Apple Silicon: tar.gz for macOS 14 or newer. Build from source for other architectures.

The binaries do not need Rust installed. Verify the checksum, extract the archive,
and follow the [quick start](docs/QUICKSTART.md) to configure `agent.yaml`. Choose a
reachable model endpoint and set the environment variable named by `auth_env`.
Remove `auth_env` for a keyless local model. Put a non-sensitive file in `workspace`, then run:

```sh
./mote --jsonl agent.yaml "List the workspace files and summarize them"
```

On Windows PowerShell, replace `./mote` with `.\mote.exe`.
The example allows only directory listing and file reads. JSONL events appear
after the run; without `--jsonl`, the CLI prints the final state and event count,
not a chat answer. Binaries are unsigned and not notarized. Checksums verify
integrity, not publisher identity.

## What is included

- The Rust library and CLI support YAML configuration, capability grants, model
  fallback, run/tool/observation limits, cancellation hooks, and post-run events.
- Built-in tools read and write files, list directories, run restricted shell/Git
  commands, and make optional typed decisions. Access is opt-in.
- Applications can register Rust handlers or Python/TypeScript/Go callbacks through
  the separate `mote-bridge` process. Registration does not grant permission.
- Thin Python/TypeScript/Go SDKs and a Rust starter have no third-party SDK runtime dependencies.
- Distribution checks include cross-platform CI and extracted-archive integration
  tests. Archives include dependency license notices; releases include SHA-256 checksums.

The Rust package is `mote-agent`, imported as `mote`. The unrelated `mote` package
on crates.io is not this project. Install the Python/TypeScript/Go SDKs from source,
not npm/PyPI, and match the source and SDKs to your installed release.

## What MOTE does not provide

MOTE has no bundled model, hosted inference, GUI, database/memory service,
scheduler, browser, MCP server, multi-user access control, dynamically loaded
plugins, or OS sandbox. `AGENT.md` does not configure the runtime; the CLI reads YAML.

Tool callbacks are interactive. Runtime events are buffered until the run ends.
Timeouts cannot guarantee that all descendant processes or blocking host callbacks
stop. A completed run or green CI does not establish model-output correctness or
provider availability. Deterministic tests use scripted model responses; check real
outputs in your deployment.

## Documentation

- [Quick start and troubleshooting](docs/QUICKSTART.md) ·
  [Configuration reference](docs/CONFIGURATION.md).
- [Choose an integration](docs/BUILDING_ON_MOTE.md) ·
  [Python SDK](sdk/python/README.md) · [TypeScript SDK](sdk/typescript/README.md) ·
  [Go SDK](sdk/go/README.md) · [Rust starter](starters/rust-agent/README.md).
- [Architecture and design principles](docs/ARCHITECTURE.md) ·
  [Bridge protocol](docs/BRIDGE_PROTOCOL.md) · [Security](SECURITY.md).
- [Official Go SDK specification for v0.14.1](docs/GO_SDK_SPEC.md).
- [Contributing and support](CONTRIBUTING.md) ·
  [Release procedure](docs/RELEASING.md) · [Changelog](CHANGELOG.md).

Source docs follow `production-ready`. For an installed release, read the matching
Git tag and bundled quick-start. SDK package versions and the bridge protocol
version are independent of the Rust crate version. Maintainers choose version
tags; ordinary commits do not publish binary releases.

## Community and license

[Report a bug or ask a question](https://github.com/danny-dis/MOTE/issues/new/choose) ·
[Security reporting](SECURITY.md#reporting-a-vulnerability)

MOTE uses the [MIT license](LICENSE). Binary archives include dependency license
texts and Rust standard-library notices. See [third-party notices](THIRD_PARTY_NOTICES.md).
