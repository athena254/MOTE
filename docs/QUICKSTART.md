# MOTE binary quick-start

Download your platform's archive from https://github.com/danny-dis/MOTE/releases/latest.
Source and release downloads are public and do not require a GitHub account.

- Windows x64: `x86_64-pc-windows-msvc.zip` (Windows 10/11).
- Linux x64: `x86_64-unknown-linux-gnu.tar.gz` (Ubuntu 22.04 or newer compatible system, OpenSSL 3 and CA certificates required).
- macOS Apple Silicon: `aarch64-apple-darwin.tar.gz` (macOS 14 or newer).
- Intel Macs and other architectures: build from source.

Before extracting, compare the archive's checksum with `SHA256SUMS`:

```powershell
# Windows PowerShell: compare this digest with the matching SHA256SUMS line.
Get-FileHash .\mote-v*-x86_64-pc-windows-msvc.zip -Algorithm SHA256
```

```sh
# Linux (run from the download directory; missing other platforms are ignored).
sha256sum --check --ignore-missing SHA256SUMS
# macOS: compare with the matching SHA256SUMS line.
shasum -a 256 mote-v*-aarch64-apple-darwin.tar.gz
```

Each archive includes `mote` and `mote-bridge` (`.exe` on Windows), this guide,
`agent.yaml`, an empty `workspace`, and licenses. You do not need Rust installed.
The binaries are unsigned and not notarized, so your OS may warn about them.
Checksums check integrity, not publisher identity. Do not disable system-wide
security controls to run a download.

## Command line

Open a terminal in the extracted directory and check the executable:

```powershell
.\mote.exe --version
.\mote.exe --help
```

```sh
./mote --version
./mote --help
```

Set your provider's OpenAI-compatible chat-completions endpoint and model name in
`agent.yaml`. The example points to a local server; MOTE does not bundle a model.
If the endpoint needs authentication, set the environment variable named by
`auth_env` (`MOTE_API_KEY`) before running. For a keyless local server, remove
`auth_env` from the manifest. Never put a key in YAML or commit it.
Place only the files you want the agent to read in `workspace`.

```powershell
.\mote.exe --jsonl agent.yaml "List the workspace files and summarize them"
```

```sh
./mote --jsonl agent.yaml "List the workspace files and summarize them"
```

`--jsonl` prints tool observations and lifecycle events after the run ends, not
as a live stream. Without it, the CLI prints the final state and event count
instead of a chat answer. Event output may contain sensitive data.

The example grants `list_dir` and `read_file` only. It cannot write files or run
shell commands. Permission checks do not provide an OS sandbox: use trusted
workspaces and isolate untrusted workloads externally. Successful installation
does not establish provider availability or output correctness.

## Python and TypeScript

Install the SDK from the source tag matching your release. Packages are not
published to PyPI/npm. See https://github.com/danny-dis/MOTE/blob/production-ready/docs/BUILDING_ON_MOTE.md.
Give your client the extracted bridge's absolute path:

```python
from mote_sdk import Client
client = Client(binary=r"C:\path\to\mote-bridge.exe")
```

```javascript
import { Client } from "@mote-agent/sdk";
const client = new Client({ binary: "/absolute/path/to/mote-bridge" });
```

`mote-bridge` exchanges JSON over stdin/stdout. It is neither an interactive console
nor an HTTP server; launch it through the SDK. Application callbacks have the host
user's privileges, and bridge timeouts cannot terminate every descendant process.

## Troubleshooting

### Unknown field `api_key_env`

Use `auth_env: MOTE_API_KEY`. The v0.13.0 release example contained this typo;
upgrade to v0.14.0 or fix the field. Unknown manifest fields are rejected.

### Missing auth environment variable

Set the variable named by `auth_env` in the terminal or service that launches MOTE.
The field names a variable, not a literal key. Remove `auth_env` for a keyless local model.

### Connection refused, all models failed, or HTTP 429/503

Check the endpoint, model name, credentials, quota, and provider availability.
MOTE does not start a model server. An OpenAI-compatible endpoint needs the full
chat-completions route, not just `/v1`. More retries will not fix authentication
or quota errors.

### Capability rejected

The trusted manifest must grant the action. Custom tools also need registered
host handlers. Do not grant shell access just to hide the error.

### Bridge hangs when double-clicked

The bridge has no GUI. SDKs communicate with `mote-bridge` over stdin/stdout.
For CLI use, launch `mote` from a terminal.

### SDK cannot find the executable

Pass the bridge's absolute path, ending in `mote-bridge.exe` on Windows.
SDK installation does not install or download the bridge.

### Import fails in Node

Build the source package before installing it into your app. Use the ESM import
and package name in the [TypeScript SDK](https://github.com/danny-dis/MOTE/blob/production-ready/sdk/typescript/README.md).

### Output exists after failure, or a completed result is wrong

Side effects can happen before failure. Inspect events and check output independently.
Do not blindly retry writes or external callbacks. Completion does not validate meaning.

### No live progress events

CLI JSONL and SDK runtime events arrive after the run. Host-tool callbacks are
interactive; runtime event streaming is not implemented.

### OS download warning

Releases are unsigned and not notarized. Verify the origin and checksum without
disabling system-wide protections. Build from source if your policy requires it.

### Linux shared-library error

Use the matching architecture and an Ubuntu 22.04-compatible system with OpenSSL 3
and CA certificates, or build on the target system.

For other problems, check that CLI, bridge, and SDK source versions match. Reproduce
in a small, non-sensitive workspace and follow the
[bug-report guidance](https://github.com/danny-dis/MOTE/blob/production-ready/CONTRIBUTING.md#support-and-bug-reports).
