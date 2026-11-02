# Contributing

Thanks for contributing.

## Development Setup

1. Install Go 1.26 or newer.
2. Clone the repository.
3. Run tests:

```bash
go test ./...
```

## Quality Gate

Before opening a pull request, run:

```bash
go test ./... -coverprofile=coverage.out -covermode=count
go tool cover -func=coverage.out | tail -n 1
```

CI requires 100% statement coverage across all packages, including the live-smoke
command. Its tests use local HTTP servers with the real API client; no tenant
credentials or external API calls are needed.

## Releases

Releases are driven by version tags, not automatic version bumps. After merging
changes to `main` and confirming CI passes, create and push a new tag:

```bash
git switch main
git pull --ff-only
git tag -a v0.3.1 -m "Release v0.3.1"
git push origin v0.3.1
```

Use the next appropriate version; the tag above is an example. Never move an
existing published tag.

The `release` Actions workflow validates the tagged commit using the same vet,
tests, and 100% coverage gate as CI, then publishes a GitHub Release with generated
notes. Tags must use `vMAJOR.MINOR.PATCH`, optionally followed by a prerelease suffix
such as `-rc.1`. Prereleases are marked as such and are not marked Latest.
This is a Go library: GitHub provides source archives and no binaries are uploaded.
The workflow uses the built-in `GITHUB_TOKEN`; no additional secret is required.

To publish a release for an existing tag that has no GitHub Release, select
**Actions > release > Run workflow** on `main` and enter the tag. The workflow tests
that tag's commit, not the current `main`. Older tags that fail the quality gate
cannot be released through this workflow; publish a new fixed version instead.
Running the workflow for an already released tag fails without modifying its
release.

## Pull Request Checklist

1. Keep changes focused and scoped.
2. Add or update tests for behavior changes.
3. Ensure tests pass locally.
4. Update README when public usage changes.
5. Use clear commit messages.
