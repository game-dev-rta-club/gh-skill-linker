---
title: Development and testing policy
updated: 2026-07-15
status: implemented
---

# Development

Keep external services behind adapters so normal tests run entirely locally.

Use test-driven development for behavior changes:

1. Write a failing test.
2. Add the smallest implementation.
3. Simplify it.
4. Update user guidance and specifications.

| Layer | Test |
| --- | --- |
| domain | pure unit test |
| service | fake dependency |
| filesystem | `t.TempDir()` |
| Git | temporary repository and bare remote |
| GitHub API | test transport or server |
| CLI | injected dependency |

Normal tests do not connect to live GitHub.

```bash
go test ./...
go test -race ./...
go vet ./...
go build ./...
```

Live E2E requires an owned branch, canonical-repository guard, and verified
cleanup.
