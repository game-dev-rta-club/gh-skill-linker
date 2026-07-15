---
title: Source reference
updated: 2026-07-15
status: implemented
---

# Source reference

Branches and tags are stored as full Git refs.

| Kind | Full ref | Pull | Push |
| --- | --- | --- | --- |
| branch | `refs/heads/<name>` | yes | yes, with permission |
| tag | `refs/tags/<name>` | no | no |

Install accepts exactly one of `--branch` and `--tag`. Direct commit-SHA
selection is not supported.

## Resolution

1. Read the exact ref through the GitHub Git refs API.
2. Peel an annotated tag through the Git tags API to its final object.
3. Reject a final object that is not a commit.
4. Read discovery results and snapshots from that same commit.

Stored SHAs:

- `refSHA`: branch or lightweight-tag commit, or annotated-tag object
- `commitSHA`: commit after peeling
- `treeSHA`: skill subtree

If a tag's `refSHA` changes, status reports `tag_moved` even when content is
identical. The extension never accepts that change automatically.

## Operation rules

- install: resolve the ref once and read every snapshot from the same commit
- status: inspect tag local, baseline, current, and ref identity
- pull: reject a tag with `fixed_source_ref` before reading remote content
- push: reject a tag with `source_ref_read_only` before checking permission
- re-pin: select another tag for the same source through `install`; update only
  when local content is clean

Re-pin completes retrieval and validation before replacing the skill, then
writes the manifest atomically. On failure, it restores the original skill.
`status` detects an inconsistency caused by process termination.

## Compatibility

On read, schema v1 `branch` becomes `refs/heads/<branch>` with
`refSHA=commitSHA`. Read-only commands do not rewrite the file. The next
successful manifest mutation saves schema v2.

Not supported:

- `latest` or semantic-version ranges
- automatic branch/tag selection or switching
- default-branch fallback
- commit-SHA selection
- tag creation, movement, deletion, or push
