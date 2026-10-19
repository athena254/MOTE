# Security and deployment scope

Use MOTE with trusted operators and workspaces, a restricted manifest, and external
isolation. It is not a sandbox for untrusted model output, code, repositories, or
users. Passing CI does not change those limits.

## What MOTE enforces

- A tool needs a capability grant in the manifest. The host must control which
  manifest it accepts; caller-supplied grants are not a trusted permission decision.
- Built-in file actions reject absolute paths, parent traversal, and symlink
  targets that resolve outside the workspace when validated.
- `max_output_bytes` (1 to 1,048,576; default 65,536) caps shell output and rejects
  oversized file and directory observations before they enter the next model
  prompt or event stream.
- The runtime limits per-tool calls, iterations, and elapsed run time. A failed
  write cannot count as a fresh report.
- The JSONL action event omits `write_file` content. Observations and other action
  arguments are still recorded.

## What MOTE does not enforce

- Shell/Git programs have the invoking account's OS permissions. An allowlisted
  interpreter such as PowerShell, Python, or Node can execute code, read files, or
  use the network outside the workspace. `unsafe_shell: false` is not an OS sandbox.
- Path validation and file opening are separate operations. Another process with
  write access can change a symlink between the check and the open.
- Executable lookup uses the process environment's search path. An untrusted
  `PATH` directory can substitute a program.
- Timeouts and cancellation make a best-effort attempt to stop the direct process.
  Descendants and blocked operations can outlive the run.
- Model requests contain task prompts and tool observations. Logs can contain file
  contents, command output, and other action arguments. The provider may be remote.
- MOTE has no multi-tenant access control, tool plug-in sandbox, installed scheduler,
  or guarantee of live provider uptime. Host-registered handlers have the host
  process's permissions and can block past a runtime deadline. Registering a
  handler does not authorize its use.

## Deployment checklist

1. Run each workload under a dedicated low-privilege OS account or in an externally
   isolated container/VM. Give it only the files and network destinations it needs.
   Keep secrets and unrelated repositories out of its workspace and environment.
2. Grant only the capabilities the task needs. File-only agents do not need shell
   or Git. When shell access is necessary, prefer narrowly scoped programs and an
   operator-controlled search path. Do not allow interpreters on an untrusted host
   without OS isolation.
3. Choose and pin the model endpoint, and review where observations go. Protect and
   rotate provider credentials outside MOTE; keep them out of manifests. Register
   custom Rust handlers only from trusted host code. YAML cannot load handlers itself.
4. Restrict who can edit the manifest, workspace, executable search path, and runtime
   environment. Prevent adversaries from changing the workspace during a run.
5. Restrict access to JSONL logs and reports. Check report freshness and content;
   sample and historical reports do not establish current health. Arrange external
   monitoring and backups.
6. Test against the actual model endpoint and deployment environment before enabling
   unattended runs. Monitor failures, rate limits, and network timeouts externally.

## Reporting a vulnerability

Keep exploit details, credentials, and sensitive logs out of public issues. Contact
the maintainer through an existing private channel. If you have no private contact,
open an issue asking for one without disclosing the vulnerability. GitHub private
vulnerability reporting may not be enabled or available; this policy does not rely on it.

Include the affected release/commit, operating system, a minimal reproduction,
security impact, and proposed mitigation. Remove real keys and private user data.
There is no paid support contract, bug bounty, or guaranteed response deadline.

## Updates and release trust

Use the latest published patch release and read its changelog. MOTE is pre-1.0;
older releases have no promised long-term-support or backport schedule. A fix is
available only once its release or commit has been published.

Release archives are unsigned, and macOS binaries are not notarized. Compare the
archive hash with `SHA256SUMS` from the same trusted GitHub release. A checksum
checks integrity, not publisher identity: repository or account compromise can
affect both files. Do not disable system-wide security protections to run a download.

SDK timeouts terminate the bridge, but may leave descendant processes or application
callbacks running. Bound or isolate host code separately. MOTE reads credentials
from the variable named by `auth_env`; never commit its value to YAML.
