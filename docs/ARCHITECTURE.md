# Architecture

MOTE is a library with small command-line and cross-language entry points. Your
application handles users, storage, UI, scheduling, and deployment.

```text
Rust application ───────────────────────────────┐
                                               │
CLI → YAML manifest → model chain → runtime ←───┤
                                   │           │
Python / TypeScript → stdio bridge ─┘           │
                                   │           │
                          permitted tool calls │
                                   ↓           │
                     observation → next action ┘
```

## The execution loop

1. Load a trusted manifest and establish a workspace.
2. Ask the configured model for one structured action.
3. Check its capability grant and the applicable limits.
4. Execute the built-in action or application-owned handler.
5. Record the observation and include it in the next model prompt.
6. Stop on completion, cancellation, error, or an exhausted limit.

Models propose actions. They cannot register code or grant permissions. The
application must validate outputs. The [Rust starter](../starters/rust-agent/README.md)
checks both callback execution and the resulting file.

## Source map

| Module | Responsibility |
|---|---|
| [`config.rs`](../src/config.rs) | Typed manifest and defaults |
| [`action.rs`](../src/action.rs) | Action data types |
| [`model.rs`](../src/model.rs) | Provider adapters, action parsing, fallback chain |
| [`runtime.rs`](../src/runtime.rs) | Execution loop, tools, shell/Git execution, budgets and state |
| [`capability.rs`](../src/capability.rs) | Workspace file/directory operations and path validation |
| [`decision.rs`](../src/decision.rs) | Optional typed decision providers and fallback |
| [`event.rs`](../src/event.rs) | Lifecycle event representation |
| [`main.rs`](../src/main.rs) | YAML CLI and post-run JSONL output |
| [`mote-bridge.rs`](../src/bin/mote-bridge.rs) | One-task, bidirectional stdio bridge |

## Language boundaries

Rust applications embed the public library directly. `Model` supplies actions and
`CustomTool` registers application handlers. Rust applications do not need the bridge.

Python and TypeScript applications launch one bridge process per run. The bridge
requests custom-tool results from the host and uses the Rust runtime's checks.
SDKs validate frames, enforce child-process deadlines, and check terminal results
against the exit status. They contain neither an HTTP worker service nor a second
agent loop.

A custom tool needs both a host registration and a grant in the trusted manifest.
The bridge carries data; it does not dynamically load code.

## Boundaries that matter

- Workspace checks cover built-in file operations, not every host operation.
- Shell, Git, and custom handlers retain host-user privileges.
- Timeouts cannot safely interrupt arbitrary synchronous callbacks or guarantee
  that every descendant process exits.
- CLI events and SDK event lists arrive after the run. Bridge callbacks are interactive.
- Deterministic tests check execution mechanics. Check live model responses and
  business outcomes separately in the intended deployment.

Read [Security](../SECURITY.md) before adding capabilities or using untrusted content.

## Design principles

Keep the execution kernel small and stable. YAGNI (You Aren't Gonna Need It) means
deferring features without a concrete consumer. Architecture that current users
need still belongs in the design.

- Add the smallest extension contract a real application needs. Prefer narrow
  protocols to provider-specific knowledge in the kernel.
- Keep specialized integrations and deployment behavior in application code.
  Replacing a component should not require restructuring the execution loop.
- Check model proposals deterministically through named grants, workspace
  validation, observations, and budgets. These checks do not provide an OS sandbox.
- Build the security boundaries required by the threat model, protect against data
  loss, and preserve stable interfaces that other systems already consume.
- Use YAGNI to guide engineering and review. It is not a runtime capability.

Before adding a core feature, identify its consumer and the task it enables.
Explain why an application tool cannot supply it, what complexity it removes, and
its dependency, security, and maintenance costs. Defer it if there is no concrete consumer.
