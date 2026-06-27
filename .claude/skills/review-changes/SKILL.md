---
name: review-changes
description: Review the code changes from a just-completed implementation using an independent sub-agent, surfacing bugs and problems grouped into High / Medium / Low priority. Use after finishing a feature or change, or whenever the user asks to review the work.
---

# Review changes

Run an **independent review** of the code that was just implemented and report the
findings back, grouped by priority. The review MUST be done by a separate sub-agent
so it looks at the diff with fresh eyes — do not review your own work inline.

## How to run it

1. Figure out the scope of "the change":
   - Prefer the uncommitted diff: `git status` + `git diff` (and `git diff --staged`).
   - If the work is already committed on a branch, diff against the base branch
     (e.g. `git diff main...HEAD`).
2. Spawn **one independent sub-agent** (Agent tool, `subagent_type: "general-purpose"`,
   or `Explore` if you only need read-only investigation) with the review prompt below.
   The sub-agent reads the diff and the surrounding code itself.
3. Relay the sub-agent's findings to the user using the **Output format** below.
   Do not fix anything yet — this skill only reports. Wait for the user to decide.

## Sub-agent prompt (template)

> You are an independent code reviewer. Review ONLY the following change (do not
> review unrelated existing code). Read the diff and the surrounding files for
> context.
>
> Scope: <the diff / changed files / branch range>.
>
> Find real problems and bugs. For each finding give: a one-line title, the
> `file:line`, why it's a problem, and a concrete suggested fix. Group findings
> into exactly three priority sections (High / Medium / Low) using the definitions
> below. If a section has no findings, say "None". Do not invent issues to fill a
> section — an empty section is a good outcome. Be concrete and cite locations.

## Priority definitions

- **High** — must be fixed. Real bugs, crashes, data loss, security issues, race
  conditions, broken/incorrect behavior, or anything that will likely cause
  problems in production.
- **Medium** — should be fixed. Could *potentially* introduce problems (missing
  edge cases, weak error handling, fragile assumptions), or code that is just
  pretty bad — unclean, not following best practices or project conventions.
- **Low** — minor / nice-to-have. Cleanups and nits: a leftover debug log, dead
  code, a typo, naming, formatting, a stray comment — things that don't affect
  behavior.

## Output format

Report back exactly like this (omit nothing; use "None" for empty sections):

```
## Review

### 🔴 High
1. <title> — `file:line`
   Problem: <why it matters>
   Fix: <concrete suggestion>

### 🟡 Medium
1. <title> — `file:line`
   ...

### 🟢 Low
1. <title> — `file:line`
   ...
```

End with a one-line summary (e.g. "3 high, 1 medium, 2 low") and ask whether to
fix the High items now.

## Notes

- Keep the review scoped to the change. Pre-existing issues outside the diff are
  out of scope unless the change makes them materially worse.
- Project-specific checks worth flagging: presentation code leaking into the engine
  packages, anything reading/storing/logging secrets, and `EXPLAIN.md` left out of
  sync with code changes (see `AGENTS.md`).
