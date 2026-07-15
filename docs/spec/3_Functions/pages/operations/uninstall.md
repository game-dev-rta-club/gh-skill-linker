---
title: Uninstall operation
updated: 2026-07-15
status: implemented
---

# Uninstall

Project-localなmanaged skillを削除する。Remoteは参照・変更しない。

## UNIN-001 Selection

`gh linked-skills uninstall SKILL [--force]`を受け取る。`SKILL`はmanifest key、またはproject-relative destination。未管理skillは拒否する。

## UNIN-002 Safety

Destinationはproject内のregular directoryで、symlink componentを含まないことを確認する。

通常はlocal tree SHAとmanifestの`treeSHA`が一致し、追加された空directoryがない場合だけ削除する。Fileの追加、削除、byte、実行bit、空directoryの追加はlocal変更として拒否する。`--force`はlocal一致確認だけを省略し、path安全性は省略しない。

Destinationが不存在なら、data lossがないためmanifest entryだけ削除する。

## UNIN-003 Transaction

Destinationを同じparent内のtemporary directoryへ移動してから、開始時entryとのoptimistic comparisonでmanifest entryを削除する。

Manifest更新失敗時はdestinationを元へ戻す。成功後にtemporary directoryを削除する。Cleanup失敗はerrorにするが、manifest entryと元destinationは削除済みのままとする。

GitHub API、GitHub認証、remote source状態は必要としない。
