# AGENTS.md

Guidance for AI coding agents working in this repository. (Human-readable too —
think of it as the house rules.) Agents read this file automatically at the start
of a session; please follow it.

## What this project is

`bake` is a small Go CLI that wraps the external `goose` tool to give each project
its own AI assistant context. See `README.md` for the product overview and
`EXPLAIN.md` for a file-by-file explanation of the code.

## Golden rule: keep EXPLAIN.md in sync

`EXPLAIN.md` is a beginner-friendly map of the Go code, maintained for a user who
is **not** a Go developer. It must always match the code.

**Whenever you change the code, update `EXPLAIN.md` in the same change.** Specifically:

- **Added a Go file** → add a section for it under the right package, describing
  what it does and its key functions/types.
- **Deleted a Go file** → remove its section.
- **Renamed/moved a file or package** → update the path and any references to it
  (including the dependency diagram and the "typical flow" section).
- **Changed what a function/command does**, added/removed a command, flag, or
  exported function → update the relevant bullet so the description stays accurate.
- **Added or removed a third-party dependency** that matters to how the code works
  → update the libraries list in the "big picture" section.
- **Changed the project scaffold** (`internal/templates/project/...`) → update the
  template files list.

Keep the existing style of `EXPLAIN.md`: plain language, short bullets, explain
*why* not just *what*, and assume the reader doesn't know Go.

### How to verify you've kept it in sync

Before finishing, sanity-check that every `*.go` file has a corresponding mention:

```sh
find cmd internal -name '*.go' | sort   # every file here should appear in EXPLAIN.md
```

If you changed code but `EXPLAIN.md` has no diff, you almost certainly missed an
update.

## Design

The TUI is built from this palette and **nothing else** — no grays, no black.
Keep new UI on-brand with these only:

- **Pink** `#FE5283` (`brandPink`) — primary accent: titles, highlights, selected items.
- **Cyan** `#6AD5FD` (`brandCyan`) — secondary accent: borders, descriptions, gradients.
- **White** `#FFFFFF` (`brandWhite`) — primary text, only when an accent doesn't fit.
- **Lilac** `#EEADEE` (`brandLilac`) — muted / secondary text, only when an accent doesn't fit.

They're defined once as constants in `internal/tui/theme.go`; reuse those
constants rather than hardcoding hex elsewhere, and don't introduce new colors.
Aim for a fun, "poppy" feel (e.g. the animated pink→cyan gradient on the
home-screen `BAKE` banner).

## Always review after changing code

After completing any request that changes code, **automatically run the
`review-changes` skill** (an independent sub-agent review, grouped into
High/Medium/Low). Do this before considering the task done — the user should not
have to ask for it.

- Skip it only when the user explicitly opts out for that request (e.g. "no
  review", "skip the review", "don't review this").
- Also skip for changes that aren't code: docs-only edits, this file, skill files,
  config tweaks, or pure formatting.
- Also skip when the change is truly trivial, so a full review wouldn't make sense
  (e.g. a one-line tweak, a renamed variable, a typo fix, a comment change). Use
  judgment: if there's nothing meaningful for a reviewer to catch, don't run it.
- Run the review **before** committing, so any High findings can be addressed
  first. Report the findings, then continue.

## Other conventions

- The **engine** packages (`internal/config`, `internal/workspace`,
  `internal/gooserun`, `internal/templates`) must stay presentation-free — no
  printing/TUI code there, so the CLI and TUI keep sharing the same behavior.
  If you blur this boundary, call it out and update `EXPLAIN.md`'s layer table.
- **Never** read, store, or log secrets (the OpenRouter API key). goose owns that
  in the OS keyring; bake only does secret-free checks.
- Build/run: `go build ./cmd/bake` and `go vet ./...` before committing.
