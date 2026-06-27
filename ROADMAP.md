# bake-ai — Roadmap

Post-MVP additions, in build order. Builds on the MVP (`plans/initial.md`).

## 1. Auto-updating vault
After each chat, distill new facts/decisions into the project's `vault/` so you
never start cold.
- `bake distill <project>`, auto-run by `bake chat` on exit (wrapper-on-exit hook —
  deterministic; goose has no native session-end hook).
- Cheap model reads the session transcript + `vault/INDEX.md`, then **reconciles
  into existing topic notes** (never one file per chat) and updates `INDEX.md`.
- Optional auto-commit of the vault (see §4).
- Main risk: distill quality — keep the prompt strict; vault stays git-tracked so
  you can review/revert.

## 2. Chat conversation resume
Resume an old conversation for a project instead of always starting fresh.
- goose already persists sessions (`goose session --resume`, `--name`,
  `--session-id`, `--history`, `--fork`) — bake just surfaces them per project.
- `bake chat <project> --resume` (most recent) or `--resume <name>` (specific).
- TUI: list a project's past sessions (date + preview), pick one to resume;
  `--fork` to branch from it.
- Name sessions per project and scope listing to the project's working dir.

## 3. Project authoring
Manage a project's pieces — on creation, manually, or LLM-assisted. Logic in
`internal/` so CLI + TUI share it.

### On create
- `bake new` asks for an **optional description**. Store it in `.goosehints`
  ("About this project") — the always-loaded context the assistant actually uses —
  and mirror a one-line version into `recipe.yaml` `description:` (the picker
  subtitle / goose's load label). *Recommendation: `.goosehints` is the source of
  truth; recipe `description:` is just the short label.*

### Manual
- `bake skill add|list|remove <project> <name>`
- `bake prompt edit <project>` (re-renders recipe instructions)
- `bake doc add <project> <path|url>` (copy a file or fetch + convert a URL)
- `bake context edit <project>` (edit `.goosehints`)

### Assisted (LLM)
Natural-language edits: "create a skill to do X", "update the instructions to do Y".
- `bake assist <project> "<request>"` — drives goose with a dedicated **authoring
  recipe** pointed at the project folder.
- Guardrails: cwd = the project; authoring prompt restricted to "only modify files
  under this project"; limit enabled tools to file editing (no shell); **show a
  `git diff` and confirm before keeping** — revert with `git checkout` if rejected.
- Auto-commit each accepted change (§4) so every assisted edit is reviewable and
  revertible. git is the real safety net here.

## 4. Git sync
Workspace as its own private repo; let bake commit its own changes.
- Manual: `bake sync <project>` — stage + commit, optional `--push`.
- Auto: opt-in `auto_commit` in `config.toml` — authoring (§3) and vault distill
  (§1) commit themselves; push stays manual.
- Safe: scope to the touched project path, clear messages, never auto-push, no-op
  if not a repo; `.gitignore` guards secrets.
- Open question: should `bake init` auto-`git init` the workspace? Keep the tool
  repo and workspace repo separate (recommended).

## 5. Bake-owned chat via ACP (optional)
Render the chat in bake's own Bubble Tea + Glamour UI instead of handing off to
goose.
- Drive `goose acp` (stdio) as an ACP client; goose streams structured events,
  bake renders them and keeps doing the provider/MCP/tool work.
- Synergy: bake holds the transcript in memory → vault distill (§1) needs no
  session-file lookup.
- Cost: streaming, tool-call rendering, interrupts, ACP drift. The handoff stays
  the default.
- Rejected: scraping `goose run` output (fragile); full goose-agnostic rebuild
  (only if bake outgrows goose).

## 6. Project-less quick questions
A first-class way to use bake for one-off questions without picking (or creating)
a project — but still backed by a configurable recipe and skills, so it behaves
"like a project without being one."

- **The "scratch" pseudo-project.** A single, always-present context that lives
  outside `projects/` (e.g. `<workspace>/scratch/` or `<workspace>/quick/`) with
  the same shape as a project — `recipe.yaml`, `skills/`, `.goosehints` — just not
  listed as a project. `bake init` scaffolds it from a dedicated template
  (`internal/templates/scratch/...`) using a general-purpose recipe + default model.
- **Entry points:**
  - CLI: `bake ask "<question>"` — runs goose against the scratch recipe and
    answers (one-shot by default; `--chat`/`-i` to drop into an interactive
    session instead).
  - TUI: a pinned **"⚡ Quick question"** row at the top of the picker (above the
    project list), always available, no project required.
  - Bare `bake ask` with no args → interactive scratch session.
- **Configurable, like a project.** Reuse the §3 authoring commands against the
  scratch context so skills/recipe/context are editable the same way:
  `bake skill add scratch <name>`, `bake context edit scratch`,
  `bake prompt edit scratch`, `bake assist scratch "<request>"`. *Recommendation:
  make the engine functions take a project-or-scratch handle so CLI + TUI + these
  commands share one code path; scratch is just a `Project` with a fixed name and
  a path outside `projects/`.*
- **Engine changes:** `workspace.Scratch(c)` returns the scratch handle (creating
  it on demand if missing); `workspace.List` keeps excluding it so it never shows
  as a normal project; `gooserun` gains a one-shot `Ask` path alongside
  `LaunchChat`.
- **Optional vault:** scratch can have its own `vault/` so distill (§1) still
  applies, giving you a running memory of ad-hoc questions — opt-in, since one-offs
  are often throwaway.
- Open question: one scratch context, or named scratch presets
  (`bake ask --as research "<q>"`) for a few reusable "modes" without full
  projects? Start with one; add presets only if needed.

## Backlog / later
- Web view — open an answer in the browser.
- Homebrew tap for `brew install` (packaging is already brew-ready).
