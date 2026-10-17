# Contributing to MOTE

Keep changes focused on the execution engine. Read
[Architecture and design principles](docs/ARCHITECTURE.md) and [Security](SECURITY.md)
before proposing core features.

## Scope and discussion

For a feature request, open an issue describing the task it enables, why application
code cannot provide it, and the dependency and security costs. For a bug, include a
minimal reproduction and a redacted manifest. Keep credentials and private workspace
contents out of issues. Report security problems privately as described in `SECURITY.md`.

## Support and bug reports

Use [GitHub Issues](https://github.com/danny-dis/MOTE/issues) for usage questions and
reproducible bugs. Check the [quick-start troubleshooting](docs/QUICKSTART.md#troubleshooting)
and existing issues first. There is no promised support SLA.

Include the CLI/bridge/SDK versions, OS and architecture, source tag or commit,
redacted manifest, exact command, exit status, and expected and actual results.
Use a small, non-sensitive workspace for the reproduction. State separately whether
the runtime completed and whether its output was correct. Never attach credentials,
private workspace contents, or unredacted logs. Use the
[Security](SECURITY.md#reporting-a-vulnerability) process for vulnerabilities.

## Development setup

Clone the repository and run the checks from its root. The integration branch is
`production-ready`, not `main`. Source and downloads are public.

You need stable Rust and the native build tools for your target:

- Windows: MSVC Rust plus Visual Studio C++ build tools, or a complete compatible
  GNU/MinGW toolchain. Do not mix GNU and MSVC binaries/toolchains.
- Linux: a C toolchain, `pkg-config`, and OpenSSL development headers.
- macOS: Xcode command-line tools.

Repository release tooling uses Python 3.11, as does CI. The Python SDK's package
metadata lists its supported runtime. Use Node 22 to match TypeScript SDK CI.

```sh
cargo build --locked --release --bins
cargo test --locked --all-targets
cargo test --locked --release --all-targets
cargo fmt --check
cargo clippy --locked --all-targets -- -D warnings
cargo run --locked --example basic_agent
cargo run --locked --example custom_tool
cargo doc --locked --no-deps
python -m unittest discover -s scripts/tests -v
```

## SDK and starter checks

Build the release bridge first, then set `MOTE_BRIDGE` to its absolute path.
Without it, SDK end-to-end tests skip. In a POSIX shell:

```sh
export MOTE_BRIDGE="$PWD/target/release/mote-bridge"
```

In Windows PowerShell:

```powershell
$env:MOTE_BRIDGE = (Resolve-Path .\target\release\mote-bridge.exe).Path
```

Run these checks from the repository root:

```sh
python -m pip install "./sdk/python[test]"
python -m pytest sdk/python/tests -q
npm --prefix sdk/typescript ci
npm --prefix sdk/typescript test
npm --prefix sdk/typescript run pack:check
cargo test --locked --manifest-path starters/rust-agent/Cargo.toml
cargo fmt --check --manifest-path starters/rust-agent/Cargo.toml
cargo clippy --locked --manifest-path starters/rust-agent/Cargo.toml --all-targets -- -D warnings
```

SDK integration fixtures run real subprocesses and callbacks against a local scripted
HTTP model server. They do not test live-provider reliability. For optional real-model
testing, read `python scripts/live_smoke.py --help`, configure your endpoint, and
inspect the outputs independently. Keep credentials out of PRs.

CI also audits dependencies and verifies the distributable Cargo package:

```sh
cargo install cargo-audit --locked
cargo audit
cargo package --locked
```

Package verification expects a clean tree. Reserve `--allow-dirty` for local
pre-commit packaging checks; it is not evidence of a published release.

## Pull requests

- Work on a focused branch and target `production-ready`.
- Add a regression test before fixing a behavioral defect.
- Update public API examples, configuration docs, and changelog entries together.
- Explain changes to trust boundaries and side effects. Avoid speculative abstractions.
- Run relevant checks and report what passed, failed, or was not run.
- Check model-generated artifacts independently of the exit status.
- Preserve lockfiles, existing work, historical evidence, and user data.

## Releases and documentation

Maintainers choose a version and push its tag; CI then builds and publishes binary
archives. Follow [Releasing](docs/RELEASING.md). Do not move published tags or replace
released assets to hide mistakes. Document shipped behavior in current docs;
Git history retains obsolete proposals and reports. Use the root README as the
documentation index and the changelog for version-specific release notes.
