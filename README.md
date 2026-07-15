# gh-linked-skills — Versioned Agent Skills for GitHub

[![CI](https://github.com/game-dev-rta-club/gh-linked-skills/actions/workflows/ci.yml/badge.svg)](https://github.com/game-dev-rta-club/gh-linked-skills/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/game-dev-rta-club/gh-linked-skills)](https://github.com/game-dev-rta-club/gh-linked-skills/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Install, pin, pull, and publish project-local Agent Skills without hiding their
source or silently overwriting local work.

`gh-linked-skills` is a GitHub CLI extension for teams that want skill files in
the project they affect. It records the source repository, path, branch or tag,
and last synchronized revision so collaborators can review and reproduce the
same setup.

## Quick start

Requirements: macOS or Linux, an authenticated
[GitHub CLI](https://cli.github.com/), system Git, and a Git project.

```sh
gh extension install game-dev-rta-club/gh-linked-skills
gh linked-skills install game-dev-rta-club/agent-skills --all --tag v1.0.0
gh linked-skills status
```

The install prints the destination of each skill. A successful status check
looks like this:

```text
SKILL                PATH                                STATE  PULL                           PUSH
rubber-duck-caller   .agents/skills/rubber-duck-caller   clean  ineligible (fixed_source_ref)  ineligible (source_ref_read_only)
rubber-duck-partner  .agents/skills/rubber-duck-partner  clean  ineligible (fixed_source_ref)  ineligible (source_ref_read_only)
```

This is expected for a tag-backed install: the files are clean, while pull and
push are disabled to preserve the fixed snapshot.

Commit the installed files and their source record together:

```sh
git add .agents/skills .gh-linked-skills.json
git commit -m "chore: install agent skills"
```

```text
your-project/
├── .agents/skills/<skill>/...   # files the agent reads
└── .gh-linked-skills.json       # source ref and synchronized revision
```

The extension does not commit the parent project for you.

## Choose a tag or branch

| Source | Best for | Behavior |
| --- | --- | --- |
| `--tag <tag>` | Teams consuming a reviewed release | Fixed, read-only snapshot. `pull` and `push` are disabled. |
| `--branch <branch>` | Skill authors and collaborators | Tracks a writable branch. Supports `status`, `pull`, and `push`. |

Use a tag when you only need to consume a skill. Use a branch when you intend
to exchange changes with the source repository.

```sh
# Install every discovered skill at a fixed release.
gh linked-skills install OWNER/REPO --all --tag TAG

# List available skills on a branch, then install one by name or path.
gh linked-skills install OWNER/REPO --branch BRANCH
gh linked-skills install OWNER/REPO SKILL --branch BRANCH
```

## What it protects

- Local changes are not silently discarded.
- A push stops if the source skill changed since the last synchronization.
- Text conflicts are left as Git-style conflict markers for manual resolution.
- Tag-backed skills stay fixed and cannot push to their source.
- The manifest records where every managed skill came from.

Review any Agent Skill before installing it. Skills are instructions to an
agent and should be treated as trusted code.

## Common workflows

| Goal | Command |
| --- | --- |
| Discover skills | `gh linked-skills install OWNER/REPO --branch BRANCH` |
| Install one skill | `gh linked-skills install OWNER/REPO SKILL --branch BRANCH` |
| Install a fixed release | `gh linked-skills install OWNER/REPO SKILL --tag TAG` |
| Check local and source changes | `gh linked-skills status` |
| Bring source changes into the project | `gh linked-skills pull SKILL` |
| Send a local skill change to its source branch | `gh linked-skills push SKILL` |
| Publish a new local skill | `gh linked-skills publish OWNER/REPO SKILL --branch BRANCH` |
| Stop managing a skill | `gh linked-skills uninstall SKILL` |

Run `gh linked-skills <command> --help` for complete arguments and examples.
For installation, branch collaboration, tag upgrades, conflicts, and removal,
see the [user guide](docs/user-guide.md).

## Security and release integrity

Use v0.5.3 or later. Releases from v0.5.3 onward are immutable and include
SHA-256 checksums and signed GitHub build provenance. Earlier releases remain
available only as historical records.

See [Verifying release artifacts](docs/release-verification.md) for the exact
checksum and attestation commands. Read the
[safety model](docs/spec/3_Functions/pages/architecture/safety-model.md) before
automating write operations. Report vulnerabilities privately as described in
[SECURITY.md](SECURITY.md).

## Requirements and upgrades

- macOS or Linux; Windows is not currently supported
- [GitHub CLI](https://cli.github.com/)
- system Git
- `gh auth login` or a supported token environment variable for remote
  operations

Upgrade an existing installation with:

```sh
gh extension upgrade game-dev-rta-club/gh-linked-skills
```

## Documentation

Start at the [documentation index](docs/README.md) for links organized by goal.
The detailed design specifications are retained separately for contributors.

## Project status

This project is pre-1.0 and maintained by volunteers. Response times, releases,
fixes, and long-term maintenance are not guaranteed. There is no support SLA.

Use [GitHub Issues](https://github.com/game-dev-rta-club/gh-linked-skills/issues)
for reproducible bugs and proposed improvements. General contact is available
through the [Game Dev RTA Club Google Group](https://groups.google.com/g/game-dev-rta-club).

## Development

Requirements: Go 1.26.5 or later, GitHub CLI, and system Git.

```sh
go mod verify
go test ./...
go test -race ./...
go vet ./...
go build ./...
```

## Contributing

Contributions and forks are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for
the development and pull-request process.

## License

[MIT](LICENSE) © 2026 Game Dev RTA Club.
