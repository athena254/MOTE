# MOTE TypeScript / Node.js SDK

Use the MOTE Node SDK to run MOTE and handle tool calls in your Node
application. It communicates with the `mote-bridge` executable over a
versioned stdio protocol. The package is marked private, is not published to
npm, and has no runtime npm dependencies.

You need Node.js 20 or newer (`engines.node = ">=20"`). Download a prebuilt
bridge from [GitHub Releases](https://github.com/danny-dis/MOTE/releases) to
use it without a Rust toolchain. Source and downloads are public and do not
require a GitHub account. Release binaries are unsigned and not notarized, so
your OS may warn about them. Verify archive checksums. The
[QUICKSTART](../../docs/QUICKSTART.md) lists release assets and extraction steps.

## Quick start

From the repository root, build the bridge if you have not downloaded one,
then build the SDK. Use your downloaded `mote-bridge` path when configuring
the client; append `.exe` on Windows.

```bash
# Source checkout option; requires stable Rust.
cargo build --locked --release --bin mote-bridge

npm --prefix sdk/typescript ci
npm --prefix sdk/typescript run build
```

Install the built local package from your application's directory:

```bash
npm install /absolute/path/to/MOTE/sdk/typescript
```

Set `MOTE_BRIDGE` to the downloaded or built executable,
`MOTE_MODEL_ENDPOINT` to a reachable OpenAI-compatible chat-completions
endpoint, and `MOTE_MODEL` to the model name. For authentication, set
`MODEL_API_KEY` and use `auth_env: "MODEL_API_KEY"` in the manifest. Keep the
credential out of YAML and source. You supply the provider, gateway, or model
service; the SDK does not install or start one.

From the repository root, after building the SDK:

```bash
node sdk/typescript/example.mjs
```

The example defaults to a local DMR-X-style endpoint and model `auto`; change
the environment for your provider. It needs a reachable model endpoint and
checks that the callback ran and that the recorded observation is `HELLO`.

See [CONFIGURATION](../../docs/CONFIGURATION.md) for manifest fields, model
authentication, and limits. The [integration guide](../../docs/BUILDING_ON_MOTE.md)
and [quickstart](../../docs/QUICKSTART.md) cover the rest of the setup.

## API

```typescript
import { Client } from "@mote-agent/sdk";

const result = await new Client({
  binary: "/absolute/path/to/mote-bridge",
  timeoutMs: 150_000,
}).run({
  manifest: {
    name: "demo",
    capabilities: ["uppercase"],
    models: [{
      provider: "openai-compatible",
      model: "your-model",
      endpoint: "https://your-provider.example/v1/chat/completions",
      auth_env: "MODEL_API_KEY",
    }],
  },
  task: 'Call uppercase with {"text":"hello"}, then complete.',
  tools: {
    uppercase: {
      description: "Input: text string. Returns uppercase text.",
      handler: async (input) => {
        if (typeof input.text !== "string") throw new Error("Expected text");
        return input.text.toUpperCase();
      },
    },
  },
});
if (!result.success) throw new Error("Agent run failed");
console.log(result.events);
```

`new Client({ binary, args, timeoutMs, cwd, env })` starts one trusted
executable per run without a shell. Use `args` for your own launcher.
`timeoutMs` must be positive, finite, and within the Node timer limit; it
defaults to 150,000 ms. `env` overrides inherited environment variables, and
`cwd` sets the child process's working directory.

`run({ manifest, task, tools })` resolves to a `RunResult` containing `state`,
`events`, `returncode` and `success`. Failed or cancelled runtime work returns
`success: false`. Setup, protocol, transport, and deadline failures reject
with `MoteError`. Inspect its `code` (`SETUP`, `PROTOCOL`, `TRANSPORT` or
`TIMEOUT`) to identify the failure. A completed run still needs your
application's acceptance checks.

Tool handlers receive JSON objects and must return strings or promises of
strings. Validate the fields each handler uses. Exceptions and non-string
results produce a generic tool error; the error message does not include
exception text or sensitive input. To allow a registered tool to run, also
grant its name in the manifest's capabilities. The SDK does not retry calls
or runs. A failed run may already have changed files or other application state.

Frames are limited to 1 MiB, including LF. The SDK retains at most 64 KiB of
stderr. You receive events and observations in the final result; there is no
live-event or progress-stream API. Avoid logging sensitive events. On timeout,
the SDK kills and reaps the direct bridge child even if a handler promise
never resolves. It cannot stop blocking JavaScript work or forcibly cancel
arbitrary descendants. The SDK and bridge do not provide an OS sandbox;
callbacks run with your application's privileges.

## Tests and package check

From the repository root, with `MOTE_BRIDGE` set to a bridge executable:

```bash
npm --prefix sdk/typescript test
npm --prefix sdk/typescript run pack:check
```

The tests exercise real subprocesses and callbacks against a local scripted
HTTP model; they do not test live inference. `pack:check` installs the built
package in a temporary consumer and checks its public import and TypeScript
declarations. It does not publish to npm.
