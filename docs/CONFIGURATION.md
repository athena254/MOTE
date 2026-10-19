# Configuration reference

The CLI reads a YAML manifest defined by [`src/config.rs`](../src/config.rs).
It rejects unknown fields. Your application must control the manifest it accepts;
arbitrary caller-supplied grants do not establish authorization.

## Minimal read-only manifest

```yaml
name: document-worker
workspace: ./workspace
capabilities: [list_dir, read_file]
max_iterations: 8
models:
  - provider: openai-compatible
    model: your-model-name
    endpoint: http://127.0.0.1:8080/v1/chat/completions
    auth_env: MOTE_API_KEY
```

Set the endpoint and model to your provider's values. `auth_env` names an environment
variable; it never contains the key itself. Omit `auth_env` when a local endpoint
needs no authentication. MOTE neither expands `${VARIABLE}` inside YAML nor starts
a model server. The example endpoint is not a hosted MOTE service.

The CLI resolves relative `workspace` paths from the manifest's location. SDK
manifests resolve them from the bridge process's working directory (`cwd`).
Use an absolute workspace to avoid that difference.

## Top-level fields

| Field | Default | Meaning |
|---|---|---|
| `name` | Required | Agent/run configuration name |
| `workspace` | `.` | Root for built-in file operations |
| `capabilities` | `[]` | Granted tool names; empty means no tools |
| `max_iterations` | `20` | Maximum execution-loop iterations |
| `max_per_tool` | `5` | Physical invocation budget per capability/tool name |
| `max_runtime_seconds` | `120` | Elapsed runtime budget; not hard OS isolation |
| `command_timeout_seconds` | `30` | Timeout for direct shell/Git commands |
| `max_output_bytes` | `65536` | Observation limit; permitted range `1` to `1048576` |
| `shell` | Empty allowlist, `unsafe_shell: false` | Shell/Git execution settings |
| `output_file` | None | Optional workspace-relative final report path |
| `models` | `[]` | Ordered HTTP model chain; CLI/bridge require at least one |
| `decision` | None | Optional typed decision provider configuration |

An embedded Rust application can supply its own `Model` instead of the HTTP model
chain. Optional fields still need valid combinations: runtime setup checks grants,
paths, and observation limits.

## Models and authentication

Each model entry needs `provider`, `model`, and `endpoint`. It also accepts:

- `auth_env`: the environment-variable name. If it is present but unset, inference fails.
- `timeout_seconds`: the HTTP timeout, defaulting to `30`.

Supported adapter names:

| `provider` | Request format |
|---|---|
| `openai` or `openai-compatible` | OpenAI-style chat completions; Bearer authentication |
| `dmr-x` or `dmrx` | Same chat-completions format, for a configured DMR-X gateway |
| `google` or `gemini` | Gemini generate-content format; `x-goog-api-key` |
| `anthropic` | Anthropic messages format; `x-api-key` |

Use the full endpoint route for the adapter. Gemini's generate-content URL includes
the selected model; changing the `model` label does not rewrite that URL. An adapter
cannot guarantee compatibility with every provider, model, or endpoint version.
Manage credentials, billing, quotas, and network access outside MOTE.

When inference fails, MOTE tries the next model in order. Empty responses get one
retry before fallback. Responses marked unfinished by the provider are rejected
so a truncated action is not executed. MOTE does not automatically retry tool side
effects or whole runs, which could duplicate writes or external operations.

## Capabilities

Built-in names are `list_dir`, `read_file`, `write_file`, `shell`, `git`, and optional
`decision`. `complete` ends a run. If it writes `output_file`, the write needs
`write_file` permission and consumes the normal write budget. A successful explicit
write to `output_file` can also complete the run. An existing report does not prove
that the run produced fresh output successfully.

The host must register custom tools and grant their names in `capabilities`.
YAML cannot import Python/JavaScript functions or load Rust plug-ins. See the
[application guide](BUILDING_ON_MOTE.md).

Built-in file paths must be relative, must not traverse `..`, and must resolve
inside the workspace when validated. These checks do not cover all filesystem
races, child processes, or trusted handlers accessing the host. Read [Security](../SECURITY.md).

## Shell and Git

```yaml
capabilities: [shell]
shell:
  allowed_programs: [whoami]
  unsafe_shell: false
command_timeout_seconds: 10
```

Shell access is opt-in and is absent from the read-only quick start. By default,
MOTE parses the executable and arguments and rejects shell operators. An allowlisted
interpreter still has substantial host access. `unsafe_shell: true` executes through
a shell; do not enable it for untrusted workloads.

The Git capability restricts commands by default and requires the shell
authorization checked during runtime setup. Keep grants narrow and trusted;
do not copy broad grants from historical examples.

## Optional decisions

The `decision` object requires `provider` and `endpoint`. It accepts `model`
(default `jev-latest`), `auth_env`, `fallback`, and `timeout_seconds` (default `30`).
This typed-choice adapter is not a general approval service. See
[`src/decision.rs`](../src/decision.rs) and [the example](../specs/decision-gate.yaml)
for the request contract. The host must choose a valid fail-closed fallback.
A failed decision provider must not cause MOTE to invent an allowed action.

## CLI and events

```sh
mote --version
mote --help
mote agent.yaml "Your task"
mote --jsonl agent.yaml "Your task"
```

Use `./mote` or `.\mote.exe` when the executable is not on `PATH`.

- Exit `0`: runtime completed.
- Exit `1`: runtime failed or was cancelled.
- Exit `2`: CLI/setup error, such as an invalid manifest.

Normal output gives the state and event count. `--jsonl` prints one lifecycle event
per line after the run. Event types include `RunStarted`, `ActionProposed`,
`ActionExecuted`, `ObservationReceived`, `RunCompleted`, and `Error`. Observations
may contain sensitive file contents or command output and may be sent to the model.

A zero exit code confirms execution completed. Check answers and generated files
independently against your business requirements.
