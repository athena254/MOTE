---
name: Bug report
about: Report a reproducible runtime, SDK, packaging, or documentation problem
title: "[Bug] "
labels: ''
assignees: ''
---

## What happened?
What did you expect, and what happened instead?

## Reproduction
Share the smallest command or code example that reproduces the problem. Include
the manifest with credentials and private data removed.

## Environment
- MOTE version / commit:
- OS and architecture:
- CLI, Rust embedding, Python SDK, or TypeScript SDK:
- Python/Node version, if relevant:
- Provider adapter (no credentials):

## Evidence
Include the exit status and relevant errors after removing sensitive data. If
the run completed but produced the wrong output, show the mismatch if it is safe
to share.

## Checks
- [ ] I removed credentials, private data, and sensitive logs.
- [ ] I checked the troubleshooting advice in docs/QUICKSTART.md and searched existing issues.
- [ ] This report does not disclose a security vulnerability; I will follow SECURITY.md for private reporting.
