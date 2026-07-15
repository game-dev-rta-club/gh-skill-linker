---
title: Linked Skills
updated: 2026-07-15
status: archived
tags:
  - gh-extension
  - linked-skills
---

# Linked Skills

> [!SUMMARY]
> A GitHub CLI extension that installs skills developed across projects without
> changing source bytes, then pulls and pushes them safely.

## Value

Treat a skill not only as an installed artifact, but as a shared source that
receives improvements from multiple projects. The basic workflow is
`pull -> edit -> push`; the command detects remote updates, missing permission,
and conflicts first.

## Product boundary

- distribute the extension through `gh extension`
- officially support GitHub.com, macOS, Linux, and project scope
- place managed skills at `.agents/skills/<name>`
- store source information in `.gh-linked-skills.json` at the project root
- never add tracking metadata to skill content
- use `gh`, GitHub APIs, and system Git for transport, authentication, and merge
- exclude PR creation, force push, user scope, GHES, and Windows from the MVP

Do not migrate existing skills or user data. No managed skills exist before the
unreleased prototype is replaced directly.

## User flow

```bash
gh linked-skills install OWNER/REPO --path skills/example --branch main
git add .agents/skills/example .gh-linked-skills.json
git commit
gh linked-skills status
gh linked-skills pull example
# edit the skill
gh linked-skills push example
```

See [Functions](gh-linked-skills-functions.md) for operations and states, and
[Manual conflict resolution](gh-linked-skills-conflict-resolution.md) for the
conflict workflow.

## Pages

| Page | Contents |
| --- | --- |
| [Functions](gh-linked-skills-functions.md) | command interface, states, rejection rules |
| [Implementation](gh-linked-skills-implementation.md) | components, data, transactions, tests |
| [Distribution and support](gh-linked-skills-distribution.md) | extension and workflow-skill installation |
| [Manual conflict resolution](gh-linked-skills-conflict-resolution.md) | marker-based resolution workflow |
