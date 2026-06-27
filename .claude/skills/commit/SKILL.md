---
name: commit
description: Write git commit messages for this repo in Conventional Commits style — a "type: description" title (feat, fix, chore, docs, etc.) with an optional body. Use whenever creating a commit here.
---

# Commit messages

Commits in this repo follow **Conventional Commits**. Every commit title is:

```
<type>: <description>
```

- **`<type>`** — the kind of change (see the list below).
- **`<description>`** — a brief, imperative summary of *what changed*
  (e.g. "add splash screen", not "added" / "adds"). Lowercase start, no trailing period.
- Keep the title under ~72 characters.

## Allowed types

| Type       | Use for…                                                        |
|------------|-----------------------------------------------------------------|
| `feat`     | a new user-facing feature or capability                         |
| `fix`      | a bug fix                                                       |
| `docs`     | documentation only (README, EXPLAIN.md, AGENTS.md, comments)    |
| `chore`    | tooling, deps, config, build, or other non-product housekeeping |
| `refactor` | code change that neither fixes a bug nor adds a feature         |
| `style`    | formatting only (gofmt, whitespace) — no behavior change        |
| `test`     | adding or adjusting tests                                       |
| `perf`     | a performance improvement                                       |

If a change spans several types, pick the one that best describes the *primary*
intent (a feature that also tweaks docs is still `feat`).

## Body (optional)

Add a body when the title isn't enough — to explain *why*, list notable changes,
or call out anything non-obvious. Leave a blank line after the title, then use
short bullet points. Skip the body for small, self-explanatory commits.

## Examples

Simple:

```
fix: prevent list height going negative on short terminals
```

With a body:

```
feat: add startup splash screen

- show the animated BAKE banner full-screen for 1.5s on launch
- any key skips it; Ctrl+C exits the app
- home screen now uses a minimal static header instead of the banner
```

Docs-only:

```
docs: document the brand palette in AGENTS.md
```
