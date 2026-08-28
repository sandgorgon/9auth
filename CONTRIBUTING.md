# Contributing

Thanks for considering a contribution to `9auth`. Like its sibling
`9p`, this project has one hard rule that shapes everything else:
**no dependencies beyond the Go standard library.** Please don't
send a PR that adds one, even a small or well-known one — `go build`
should never have to reach outside `GOROOT`.

## Getting started

```
git clone git@github.com:sandgorgon/9auth.git
cd 9auth
go build ./...
go test ./...
```

No other setup is required — no `go.sum`, no toolchain beyond `go`
itself.

## Branching (Gitflow)

This repo follows [Gitflow](https://nvie.com/posts/a-successful-git-branching-model/),
the same as `9p`.

- **`master`** — release history only. Every commit on it is a
  tagged release (`vX.Y.Z`). Never commit to it directly; it only
  moves via a `release/*` or `hotfix/*` merge.
- **`develop`** — the integration branch, and GitHub's default
  branch. Unreleased work accumulates here, mirrored by the
  `[Unreleased]` section of `CHANGELOG.md`. Never commit to it
  directly either — merge in via PR from a `feature/*` branch.

| Branch      | Cut from  | Merges into        | Naming             | Purpose                                  |
|-------------|-----------|---------------------|--------------------|-------------------------------------------|
| `feature/*` | `develop` | `develop`           | `feature/<name>`   | New work, in progress or small changes    |
| `release/*` | `develop` | `master` + `develop`| `release/<X.Y.Z>`  | Stabilize a release (version bump, final CHANGELOG edits, no new features) |
| `fix/*`     | `develop` | `develop`           | `fix/<name>`       | Bug fix that isn't urgent enough for a hotfix |
| `hotfix/*`  | `master`  | `master` + `develop`| `hotfix/<name>`    | Urgent production fix, released outside the normal release cycle |

Open PRs against `develop`, not `master`.

## Before you send a PR

```
go build ./...
go vet ./...
gofmt -l .          # should print nothing
go test -race ./...
```

## What to send

- **Bug fixes**: welcome, with a test that fails before the fix and
  passes after.
- **New consumers' needs**: if 9vcs, 9sh, or another 9-family
  program needs something this package doesn't yet expose (a new
  `Permission` level, a different config-dir layout, etc.), open an
  issue describing the concrete need before sending a PR — this
  package is a shared dependency, so a change here affects every
  consumer at once.
- **New third-party dependencies, build tools, or codegen steps**:
  won't be merged, regardless of how small or standard.

## Code conventions

- Match the existing style: `gofmt`-formatted, doc comments on every
  exported identifier, no comments that just restate what the code
  does.
- This package parses two kinds of local files
  (`authorized-peers`, `known-peers`) and PEM/DER key material — none
  of it is data received over a network from an untrusted party, but
  parsing errors should still be specific (`file:line: ...`) rather
  than opaque, since these are hand-edited files.

## Reporting issues

Open a GitHub issue with what you expected, what happened instead,
and which consumer (9vcs, 9sh, ...) you hit it from.
