---
name: cut-release
description: Cut a new versioned release of bake. Use whenever releasing a new version.
---

# Cut a release

Releases use **Semantic Versioning** via annotated git tags (`vMAJOR.MINOR.PATCH`,
e.g. `v1.2.0` — the `v` prefix is mandatory for Go modules). GoReleaser builds the
cross-platform binaries, generates the GitHub Release notes from the commits since
the last tag, and updates the Homebrew tap formula.

## 1. Pick the next version

Look at the commits since the latest tag (`git log $(git describe --tags --abbrev=0)..HEAD`)
and bump based on the highest-impact change. The Conventional Commit types (see the
`commit` skill) map straight onto SemVer:

| Change in the range                     | Bump        | Example          |
|-----------------------------------------|-------------|------------------|
| A breaking change (`!` or `BREAKING CHANGE`) | **major** | `1.4.2` → `2.0.0` |
| Any `feat:`                             | **minor**   | `1.4.2` → `1.5.0` |
| Only `fix:` / `perf:` (no `feat:`)      | **patch**   | `1.4.2` → `1.4.3` |
| Only `docs:` / `chore:` / `refactor:` … | no release needed (skip) |        |

While pre-1.0 the rules are looser, but bake starts at **1.0.0** so the table
above applies from the first release on.

## 2. Tag and push

```sh
git tag -a v1.2.0 -m "v1.2.0"
git push origin v1.2.0
```

## 3. Publish

```sh
export GITHUB_TOKEN=...   # PAT with write access to bake-ai AND homebrew-tap
make release             # goreleaser release --clean
```

GoReleaser auto-generates the GitHub Release notes from the commits in the range,
grouped into Features / Bug Fixes / Performance / Documentation. Nothing to write
by hand.

## First release only (v1.0.0)

The history before v1.0.0 predates Conventional Commits, so auto-generated notes
would read poorly. For **v1.0.0 only**, hand-write the notes and pass them in:

```sh
# NOTES-v1.0.0.md — short "initial public release" summary (what bake is + init/new/list/chat)
goreleaser release --clean --release-notes=NOTES-v1.0.0.md
```

From v1.0.1 onward, use the normal `make release` flow — every commit after v1.0.0
follows the convention, so auto-generation works cleanly.

## Notes

- `make release` assumes the tag is already pushed — tag first (step 2), then
  publish (step 3).
- Dry-run a build without publishing: `make snapshot`.
- Verify the config before a release: `goreleaser check`.
