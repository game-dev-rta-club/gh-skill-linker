---
title: Pull reference
updated: 2026-07-15
status: implemented
---

# Pull

Reject a tag source with `fixed_source_ref` before reading local or remote
content.

Apply the current source to local content, using a three-way merge when needed.

## PULL-001 Selection and prerequisites

Select exactly one skill by name or project-relative destination. Zero or
multiple matches are errors.

Every local file must be tracked. The manifest's tracked state is ignored.
Markers, unsafe paths, an invalid source, a source name different from the
managed name, or a remote failure leave content unchanged.

## PULL-002 No-op and baseline-only update

When local equals current, do not change content. Matching SHAs are a no-op. If
only the SHA advanced, update the baseline and return `Changed=true`.

## PULL-003 Clean and merged apply

When local equals base, replace it with current. Otherwise, perform a three-way
merge.

Stage under the same parent and rename the target to a backup. Reread the target
before activation; if it differs from the starting snapshot, roll back and
return `workspace changed during pull`.

Update the baseline after activation. A manifest failure restores the original
directory. A failed rollback returns the backup path and leaves the transaction
for recovery.

A failure to remove only the backup is still an error, but new content and the
baseline remain valid. The baseline also advances when markers are produced.

Return text-conflict paths, sorted and project-relative, as `ConflictPaths`.
The CLI lists the files and exits with `1`.
