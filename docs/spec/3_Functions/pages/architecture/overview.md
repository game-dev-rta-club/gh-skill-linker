---
title: Architecture overview
updated: 2026-07-16
status: implemented
---

# Architecture

The CLI calls operation services while keeping domain logic separate from
external adapters.

## TECH-001 Runtime and build dependencies

- Go `1.26.5`
- `github.com/cli/go-gh/v2 v2.13.0`
- `gopkg.in/yaml.v3 v3.0.1`
- runtime: `gh`, Git, and network access to GitHub.com
- optional status inventory: Codex CLI

Node.js, Python, and a shell interpreter are not required.

```text
main / cli
  -> install / publish / status / pull / push / uninstall
    -> source / skill / syncstate / merge / proposal
    -> manifest / workspace / skillinventory / githubapi / gitcli / command
```

Operations use adapters through interfaces. Tests replace those adapters with
fakes.

`status`, `pull`, and `push` read the manifest and compare source, baseline, and
local state. `install` registers the source snapshot as the initial baseline.
`push --pr` and `publish --pr` share the proposal service; GitHub and hidden PR
metadata hold proposal identity while the manifest remains unchanged.
Status merges managed health with local `gh skill` inventory and enabled Codex
plugin manifests. External providers are informational only.

Each synchronization command changes only the selected skill.

Fixed boundaries:

- host: GitHub.com
- workspace: current Git worktree
- destination: `.agents/skills/<name>`
- source: repository, path, and full Git ref
- parent-project commits: outside the application
