---
title: Automation specification
updated: 2026-07-15
status: implemented
---

# Automation

## AUTO-001 CI

Pull requests and pushes to main run tests, the race detector, vet,
vulnerability scanning, and build on Ubuntu and macOS. Workflow permission is
`contents: read`. Release tags also run tests and vet.

## AUTO-002 Scheduled live E2E

The workflow runs manually and every Monday at 03:17 UTC. Only the canonical
repository may create its owned branch.

Using the built binary, discovery and `install --all` create two fixtures and a
manifest. The workflow checks status, push, remote pull, and conflict behavior.
After recording evidence, it verifies deletion of the owned branch.
