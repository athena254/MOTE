# MOTE Python SDK

Use the MOTE Python SDK to run MOTE and handle tool calls in your Python
application. It communicates with the `mote-bridge` executable over a
versioned stdio protocol. Install the SDK from source; it is not published to
PyPI and uses only the Python standard library at runtime.

You need Python 3.10 or newer (`requires-python = ">=3.10"`). Download a
prebuilt bridge from [GitHub Releases](https://github.com/danny-dis/MOTE/releases)
to use it without a Rust toolchain. Source and downloads are public and do not
require a GitHub account. Release binaries are unsigned and not notarized, so
your OS may warn about them. Verify the archive checksum, but remember that a
checksum does not establish the publisher's identity. The
[QUICKSTART](../../docs/QUICKSTART.md) lists platform assets and extraction steps.

## Quick start

From the repository root, build the bridge if you have not downloaded one,
then install the SDK. On Windows, append `.exe` to the executable path.

```bash
# Source checkout option; requires stable Rust.
cargo build --locked --release --bin mote-bridge

python -m pip install ./sdk/python
```

Point `MOTE_BRIDGE` at the downloaded or built `mote-bridge`, set
`MOTE_MODEL_ENDPOINT` to a reachable OpenAI-compatible chat-completions
endpoint, and set `MOTE_MODEL` to the model name. If authentication is required, set
`MODEL_API_KEY`; the manifest must refer to the variable by name with
`auth_env: "MODEL_API_KEY"`. You supply the provider, gateway, model, and
credentials. The SDK does not install or start a service.

From the repository root:

```bash
python sdk/python/example.py
```

The example defaults to a local DMR-X-style endpoint and model
`auto`; change the environment for your provider. It checks that the callback
ran and that the recorded observation is `HELLO`. You need a reachable model endpoint
to run it.

See [CONFIGURATION](../../docs/CONFIGURATION.md) for fields, defaults, model
providers, workspace boundaries, and timeouts. The
[integration guide](../../docs/BUILDING_ON_MOTE.md) and
[quickstart](../../docs/QUICKSTART.md) cover the rest of the setup.

## API

```python
from mote_sdk import Client, Tool

client = Client(binary="/absolute/path/to/mote-bridge", timeout=150)
result = client.run(
    manifest={
        "name": "demo",
        "capabilities": ["uppercase"],
        "models": [{
            "provider": "openai-compatible",
            "model": "your-model",
            "endpoint": "https://your-provider.example/v1/chat/completions",
            "auth_env": "MODEL_API_KEY",
        }],
    },
    task='Call uppercase with {"text":"hello"}, then complete.',
    tools={"uppercase": Tool(
        "Input: text string. Returns uppercase text.",
        lambda data: data["text"].upper(),
    )},
)
if not result.success:
    raise RuntimeError("Agent run failed")
print(result.events)
```

`Client(binary="mote-bridge", timeout=150, cwd=None, env=None)` starts one
trusted executable per `run`. Use `binary` for the executable path or an
executable/argument list for your own launcher; the SDK never passes it through
a shell. `timeout` must be a positive, finite number of seconds within Python's
thread timer limit. `env` overrides inherited environment variables, and `cwd`
sets the bridge's working directory.

`run(manifest, task, tools=None)` returns `RunResult(state, events, returncode)`.
`success` is true only for state `completed` and exit code `0`. Failed or
cancelled runtime work returns a result with success set to false. Setup,
protocol, transport, and deadline failures raise `SetupError`, `ProtocolError`,
`TransportError` or `TimeoutError`, all subclasses of `MoteError`.

`Tool(description, handler)` accepts a JSON object and must return a string.
Validate the fields your handler uses. Exceptions and non-string results
produce a generic tool failure; the failure message does not include exception
text or sensitive input. To allow a registered tool to run, also add its name
to `manifest["capabilities"]`. The SDK does not retry runs or callbacks. A
failed or timed-out run may already have changed files or other application state.

Each bridge stdout frame is limited to 1 MiB, including its newline. The SDK
discards bridge stderr through the OS without buffering it. You receive `events`
and model observations after the run; there is no live-event or progress-stream
API. They may contain sensitive data. On timeout, the watchdog kills and reaps
the direct bridge
child, but it cannot forcibly interrupt Python callbacks in the host process.
Set your own bounds or isolation for blocking callbacks and descendants.
The SDK and bridge do not provide an OS sandbox. Host code runs with your
application's privileges.

## Tests

From the repository root:

```bash
python -m pip install "./sdk/python[test]"
python -m pytest sdk/python/tests -q
```

Set `MOTE_BRIDGE` to a bridge executable before running the tests. The end-to-end
tests exercise a real subprocess and callback against a local scripted HTTP
model. They are deterministic; they do not test live inference.
