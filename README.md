# gh-linked-skills

[![CI](https://github.com/game-dev-rta-club/gh-linked-skills/actions/workflows/ci.yml/badge.svg)](https://github.com/game-dev-rta-club/gh-linked-skills/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A GitHub CLI extension for installing and synchronizing project-local Agent
Skills from GitHub repositories.

It records each managed skill's source repository, path, branch or tag, and
last synchronized revision. Branch-backed skills can pull and push changes;
tag-backed skills remain fixed, read-only snapshots.

## Status

This project is pre-1.0 and supports macOS and Linux. Windows is not currently
supported.

The project is maintained by volunteers. Response times, releases, fixes, and
long-term maintenance are not guaranteed.

## Requirements

- macOS or Linux
- [GitHub CLI](https://cli.github.com/)
- system Git
- GitHub authentication from `gh auth login` or a supported token environment
  variable for remote operations

## Install

```sh
gh extension install game-dev-rta-club/gh-linked-skills
gh linked-skills --help
```

Upgrade an existing installation with:

```sh
gh extension upgrade game-dev-rta-club/gh-linked-skills
```

## Verify a release

Use v0.5.3 or later. Earlier releases remain available only as historical
records. Releases from v0.5.3 onward are immutable and include SHA-256
checksums and signed GitHub build provenance.

Download a release into an empty directory:

```sh
version=v0.5.3
gh release download "$version" --repo game-dev-rta-club/gh-linked-skills
```

Verify the checksums with `sha256sum -c SHA256SUMS` on Linux or
`shasum -a 256 -c SHA256SUMS` on macOS. Then verify that every asset came from
the release workflow, tag, and a GitHub-hosted runner:

```sh
for asset in SHA256SUMS gh-linked-skills_"$version"_*; do
  gh attestation verify "$asset" \
    --repo game-dev-rta-club/gh-linked-skills \
    --signer-workflow game-dev-rta-club/gh-linked-skills/.github/workflows/release.yml \
    --source-ref "refs/tags/$version" \
    --deny-self-hosted-runners
done
```

## Quick start

Run the extension inside an existing Git project. This example installs the
versioned Game Dev RTA Club skill set:

```sh
gh linked-skills install game-dev-rta-club/agent-skills --all --tag v1.0.0
git add .agents/skills .gh-linked-skills.json
git commit -m "chore: install agent skills"
```

Managed skills are stored at `.agents/skills/<name>`. Commit both the installed
skills and `.gh-linked-skills.json` so collaborators use the same source.

## Commands

```text
gh linked-skills install <owner>/<repository> [<skill-or-path>] (--branch <branch> | --tag <tag>)
gh linked-skills publish <owner>/<repository> <skill> --branch <branch>
gh linked-skills status [--json]
gh linked-skills pull <skill>
gh linked-skills push <skill>
gh linked-skills uninstall <skill> [--force]
```

Install without a skill name lists available skills. Add `--all` to install
every discovered skill. Run `gh linked-skills <command> --help` for complete
arguments and examples.

## Synchronization model

- Branch sources are writable and support `status`, `pull`, and `push`.
- Tag sources are fixed snapshots and cannot be pulled or pushed.
- Local changes are not silently discarded.
- Conflicting pulls write Git-style conflict markers for manual resolution.
- Push is rejected when the remote branch changed after the last synchronization.

Read the [safety model](docs/spec/3_Functions/pages/architecture/safety-model.md)
before automating write operations. The [documentation index](docs/README.md)
links to the complete command and implementation reference.

## Development

Requirements: Go 1.26.5 or later, GitHub CLI, and system Git.

```sh
go mod verify
go test ./...
go test -race ./...
go vet ./...
go build ./...
```

## Support and maintenance

Use [GitHub Issues](https://github.com/game-dev-rta-club/gh-linked-skills/issues)
for reproducible bugs and proposed improvements. General contact is available
through the [Game Dev RTA Club Google Group](https://groups.google.com/g/game-dev-rta-club).

There is no support SLA. Report vulnerabilities privately as described in
[SECURITY.md](SECURITY.md).

## Contributing

Contributions and forks are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for
the development and pull-request process.

## License

[MIT](LICENSE) © 2026 Game Dev RTA Club.
