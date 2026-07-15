---
title: Filesystem specification
updated: 2026-07-15
status: implemented
---

# Filesystem

Restrict writes to the project and define symlink and path-collision
boundaries.

## PATH-001 Filesystem containment

Convert root and target to absolute, clean paths. Reject a path outside the
root, a symlink, or an intermediate non-directory. Do not inspect components
after the first missing component.

Reject remote paths that are absolute, contain backslashes, `.` or `..`, or are
not canonical.

## PATH-002 Existing-install containment gap

New install, status, pull, and push operations check containment.

Only a rerun of an existing managed install reads the destination first. If a
parent is replaced with a symlink, this read may reach outside the project. It
does not write outside the project.

## PATH-003 Case and Unicode boundary

The extension does not detect case-folding or Unicode-normalization collisions.
On affected filesystems, they may produce collisions, order-dependent behavior,
or errors. Map iteration order is undefined.
