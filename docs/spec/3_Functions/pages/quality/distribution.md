---
title: Distribution specification
updated: 2026-07-15
status: implemented
---

# Distribution

## DIST-001 Release artifacts

After tests and vet pass for a `v*` tag, run
`cli/gh-extension-precompile@v2`.

- darwin-amd64
- darwin-arm64
- linux-amd64
- linux-arm64

Build with CGO disabled and `-trimpath -s -w`. Do not create a Windows artifact.
The release job runs on Ubuntu with `contents: write`.
