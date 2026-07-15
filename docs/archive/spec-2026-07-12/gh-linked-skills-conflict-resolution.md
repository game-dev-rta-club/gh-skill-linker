---
title: Linked Skills manual conflict resolution
updated: 2026-07-15
status: archived
---

# Manual conflict resolution

> [!SUMMARY]
> Pull does not choose meaning automatically. The user edits the Git diff3
> markers written to ordinary files.

## Flow

```text
pull
  -> write text conflicts to files
  -> advance the manifest baseline to the remote head
  -> user resolves markers
  -> status reports push
  -> push
```

Markers use this format:

```text
 <<<<<<< gh-linked-skills:local
 local content
 ||||||| gh-linked-skills:base:<tree-sha>
 base content
 =======
 remote content
 >>>>>>> gh-linked-skills:remote:<tree-sha>
```

The user deletes marker lines and unwanted content, leaving only the final
accepted content. While any marker remains, `status` reports `conflict` and
pull/push are unavailable with `unresolved_conflict`. There is no dedicated
index or conflict-state file.

Binary conflicts, modify/delete conflicts, and file/directory collisions do not
become markers and do not change local content. A text conflict inside skill
frontmatter uses the same markers; push revalidates Codex compatibility after
resolution.

## Related

- [Overview](gh-linked-skills.md)
- [Functions](gh-linked-skills-functions.md)
