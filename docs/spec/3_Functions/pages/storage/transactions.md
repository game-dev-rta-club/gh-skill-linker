---
title: Transaction specification
updated: 2026-07-15
status: implemented
---

# Transactions

Roll back local mutations where possible. Do not roll back a remote push. The
extension does not recover automatically from process interruption.

## TXN-001 Interruption and locking boundary

There is no process lock, journal, startup recovery, or signal handler. SIGINT
and similar interruptions may leave temporary directories behind.

If pull stops between its two renames, the target may be absent while only the
backup remains. Recovery is manual.

| Operation | When the manifest update fails |
| --- | --- |
| install | Remove the target |
| install `--all` | Remove the failed target and preserve successful skills |
| pull | Restore the original |
| push | Do not restore the remote |
| publish | Do not restore the remote; rerun the same command to retry registration |
| uninstall | Restore the moved original |
| manifest | Rename within the same directory after fsync |

When rollback fails, include the backup path in the error and leave the
transaction for manual recovery.
