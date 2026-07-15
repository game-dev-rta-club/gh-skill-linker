---
title: Install reference
updated: 2026-07-15
status: implemented
---

# Install

Place a remote snapshot for the first time and record its synchronization point
in the manifest.

## INST-001 Source resolution

Accept an `OWNER/REPO` repository, optional selector, and source ref. The
repository and exactly one of `--branch` or `--tag` are required. Local sources
and repository-omitted installation are not supported.

- no selector: display discovery results
- name: select one matching discovery result
- `namespace/name`: select with the namespace
- exact path: bypass discovery
- `--all`: select every discovery result

An exact path identifies an Agent Skills directory or a path ending in
`/SKILL.md`. A file path is normalized to its parent directory. A repository-
root `SKILL.md` is excluded.

Resolve the source ref to a commit SHA once and build candidates from the
recursive tree at that SHA. Results are sorted by display name. Hidden
directories are excluded.

[Install a skill](../../../2_HowToUse/pages/install-skill.md) defines the
recognized paths. When one simple name matches multiple paths, require a
namespace.

## INST-002 Snapshot install

Validate name and description from the direct `SKILL.md` and require the name
to match the source directory. The name determines
`.agents/skills/<name>`. Description length is counted in Unicode code points.

Stage the entire directory as raw bytes, paths, and executable bits. Rename it
into place only when the target is absent. Do not add metadata to content.

## INST-003 Install collision and idempotency

Reject:

- a same-name entry with a different source or destination
- a matching entry with a different local snapshot
- switching between branch and tag
- the same source registered under another name
- an occupied destination
- an unmanaged destination

Only a rerun where entry, remote, and local are identical is a no-op.

If the manifest write fails, remove the target. Include a removal failure in
the returned error. There is no lock.

## INST-004 Tag re-pin

A tag-backed skill can be reinstalled from another tag when repository, source
path, and destination are unchanged.

1. Resolve the new tag, then retrieve and validate its snapshot.
2. Confirm that local content exactly matches the old baseline.
3. Replace local content with the new snapshot.
4. Update the manifest.

Reject local changes and do not merge. A changed `refSHA` under the same tag
name requires `--accept-moved-tag`. Reject that flag when selecting a different
tag.

Even when another tag points to the same tree, update the manifest. `--all`
does not re-pin existing skills.

## INST-005 Bulk install

`--all` uses one discovery commit as the baseline for every snapshot.

Before mutation, retrieve and validate every document, name, source,
destination, local snapshot, LFS state, and path-containment rule. A name
collision prevents every install.

After validation, install one skill at a time in display-name order. Each skill
uses the INST-003 directory and manifest transaction. A later failure does not
roll back earlier successful skills. The result contains successes and the
error.
