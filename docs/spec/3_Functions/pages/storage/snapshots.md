---
title: Snapshot specification
updated: 2026-07-15
status: implemented
---

# Snapshots

Define the smallest unit used to compare and transfer a skill.

## SNAP-001 Snapshot identity

Exact identity includes:

- the set of relative paths
- raw bytes
- executable state

YAML meaning, newline conversion, timestamps, and directories are not compared.
Any execute bit means executable.

## SNAP-002 Supported entries

Remote entries must be `100644` or `100755` blobs. Local entries must be regular
files. Reject symlinks, submodules, special modes or files, and truncated trees.

A regular `SKILL.md` is required directly under the root. Validate its
frontmatter during remote reads and before a push.

- name: ASCII lowercase kebab-case, 1 to 64 bytes
- description: non-empty after trimming and no more than 1,024 Unicode code
  points
- line endings: LF and CRLF supported
- other fields: preserved

Because the installed destination is derived from name, local name must match
the parent directory. Managed name remains fixed after synchronization. Status
and push reject a local mismatch; status and pull reject a source mismatch. The
source-path basename does not need to match.

Status validates local frontmatter for push eligibility while preserving pull
eligibility. Pull does not repeat that validation at start. Only a new install
rejects an LFS pointer.

## SNAP-003 Write mode and umask boundary

Write modes are `0644/0755`. Because files are not chmodded after creation, the
system umask applies. A umask that removes execute bits prevents reproduction of
`100755`. Go [`os.WriteFile`](https://pkg.go.dev/os#WriteFile) also defines its
creation mode before umask application.

Only the manifest receives an explicit `chmod 0644`.
