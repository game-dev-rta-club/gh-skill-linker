---
title: Uninstall
updated: 2026-07-15
status: implemented
---

# Uninstall

Managed skillを現在のprojectから削除する。Source repositoryは変更しない。

```bash
gh linked-skills uninstall SKILL
```

`SKILL`はskill名、または`.agents/skills/<name>`。Localが最後の同期点と一致する場合、skill directoryとmanifest entryを削除する。

Local変更は誤削除を避けるため拒否する。変更を破棄する場合だけ明示する。

```bash
gh linked-skills uninstall SKILL --force
```

Skill directoryが既に無い場合は、残ったmanifest entryだけ削除する。GitHub認証とnetwork接続は使わない。

内部: [[docs/spec/3_Functions/pages/operations/uninstall|Uninstall operation]]
