---
title: Status reference
updated: 2026-07-15
status: implemented
---

# Status

Calculate state from local, source, and manifest tree identity, then determine
pull and push eligibility.

## STAT-001 Managed inventory and states

Inspect only manifest `skills`. Base is manifest `treeSHA`, current is the skill
tree SHA at the source ref, and local is the Git tree SHA calculated from the
destination.

Within one status call, batch all refs and push permissions. Read each
project-wide Git inventory once. There is no persistent cache. Local evaluation
does not download the baseline snapshot.

For a tag source, pull reason is `fixed_source_ref` and push reason is
`source_ref_read_only`. When remote `refSHA` differs from the baseline, pull
reason is `tag_moved`. Do not check permission for a tag.

| Local | Remote | State |
| --- | --- | --- |
| same | same | `clean` |
| changed | same | `push` |
| same | changed | `pull` |
| changed | changed | `conflict` |

If markers exist, return `conflict` without reading remote content. When current
commit SHA equals the baseline, do not read the repository tree. Otherwise,
read a tree once per repository and commit, then compare the current tree SHA.
Do not compare semantic meaning.

## STAT-002 Eligibility and reason precedence

An invalid path, file, source, or marker ends evaluation early. Then validate
local frontmatter and Git inventory before calculating state from local/current
tree SHAs.

Validation failures and `source_unavailable` produce state `null`. Markers
produce state `conflict`. A Git or permission failure preserves any state that
was already calculated.

Even with invalid frontmatter or a managed-name mismatch, calculate pull state
and eligibility; only push receives `invalid_local_skill`. A changed current
skill with a mismatched name produces `source_unavailable`. Push reason order is
frontmatter, permission, `remote_changed`, then `ignored_files`. Never overwrite
an earlier non-eligible reason.

`eligible` means only that service preconditions pass. It does not include the
final result of GitHub branch protection or rulesets.
