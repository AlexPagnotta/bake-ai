# bake-ai — MVP plan (V1)

## Context

`bake-ai` is a **public, open-source CLI tool** that gives anyone a personalized,
non-coding AI assistant for research and personal projects (coffee company, new
house, etc.) with three hard constraints: **no provider lock-in, no subscription,
not living inside the Gemini/ChatGPT apps**.

Each project carries its own prompt, skills, docs, and context, so you don't start
a new chat from zero. The **auto-updating vault, the session hook, and the
goose-agnostic layer are deferred to `V2.md`** — this plan is the smallest thing
that proves the idea and is publishable.

The MVP is split into two phases so each step is controllable, and so the two
*independent* unknowns (the goose pipeline vs. building a Charm TUI) are debugged
one at a time:
- **Phase 1 — Engine + thin CLI:** the core logic (`CreateProject`, `ListProjects`,
  `LaunchChat`) behind the thinnest possible command skin. Proves goose +
  OpenRouter + per-project context + tool/data separation.
- **Phase 2 — TUI polish:** a *second* front-end (Charm) over the **same** core
  functions — purely additive, nothing rewritten.

**Architectural rule that makes the split free (not overhead):** all logic lives in
core functions with no presentation; the Phase 1 CLI and the Phase 2 TUI are both
thin skins calling them. The plain commands survive as the scriptable/headless
interface — the TUI doesn't replace them. So Phase 1 should keep output minimal
(don't polish plain-text formatting you'd later drop).

### Verdict on the stack (OpenRouter + goose-cli + thin CLI)
- **OpenRouter** — keep it. Pay-per-token, one key, swap models freely.
- **goose-cli** — good foundation: general-purpose, provider-agnostic (OpenRouter
  first-class), ships Recipes/`.goosehints`/MCP that map onto the "project" concept.
- **Thin wrapper** — its job is to scaffold/switch projects and launch goose with
  the right recipe + working dir. (The post-chat distill that *also* justifies a
  wrapper lives in V2.)

## Public tool vs. private data (the key design call)

Separate the **tool** (published) from the **data** (private), like `git` vs. your
repos, or Obsidian vs. your vault.

- **Tool repo (public, on GitHub):** CLI code + templates only. No user data, ever.
- **Workspace (private, on your machine):** your projects. Default `~/bake`,
  configurable, XDG-aware. You can `git init` the workspace as **your own private
  repo** to back up/sync — completely separate from the public tool repo.
- **Secrets: bake handles none.** goose already stores the OpenRouter key in the
  **OS keyring** (macOS Keychain / Windows Credential Manager / Linux Secret
  Service) by default, with precedence env var → `config.yaml` (non-secrets) →
  keyring → `secrets.yaml` fallback. bake never reads, stores, or logs the key —
  it delegates to goose. `bake init` only **verifies** goose has an OpenRouter
  provider configured and, if not, points the user to `goose configure`.
- **bake's own `~/.config/bake/config.toml` holds non-secret settings only**
  (workspace path, default model name). Never committed.
- *V2/goose-agnostic note:* if bake ever calls a provider directly, use the
  keyring ladder — OS keyring (`zalando/go-keyring`) > env var > config file
  `0600` > never a `--flag` (leaks to shell history/`ps`); prompt without echo
  (`x/term.ReadPassword`), never log the key.

```
# Public tool repo (published)            # Private workspace (~/bake, your repo)
bake-ai/                                   ~/bake/
  cmd/bake/        # main entry point         config.toml      # optional, gitignore secrets
  internal/        # CLI/TUI logic            projects/
  templates/       # scaffolding (embedded)     coffee/
  go.mod  go.sum                                  recipe.yaml      # goose recipe (model + role/instructions)
  README.md  LICENSE  .gitignore                  .goosehints      # static context + persona, loaded each session
                                                  skills/          # project-specific snippets
                                                  docs/            # curated static reference
                                                  vault/           # notes (HAND-MAINTAINED in MVP)
                                                    INDEX.md
                                               house/ …
```

## Architecture
Thin **Go CLI/TUI** built on the **Charm** stack, shelling out to `goose`. Go for
a single static binary (trivial distribution) + a polished animated terminal UI.
A project = a folder = a goose recipe, living in the private workspace.

TUI stack (introduced in Phase 2):
- **Bubble Tea** — app runtime / state model (Elm-style), animations.
- **Lip Gloss** — styling and layout.
- **Bubbles** — prebuilt components (lists, spinners, viewports, inputs).
- **Huh** — interactive forms for `bake new` (name, model, scaffold options).
- **Glamour** — render the assistant's markdown answers in the terminal (pairs
  with the future "open answer in a web page" idea).

## Installation & distribution (callable from anywhere)
- Go compiles to a **single static binary, no runtime** — that binary on PATH is
  what makes `bake` callable from any directory.
- **Now (dev):** `go build -o bake ./cmd/bake` then move it onto PATH, or
  `go install ./cmd/bake` (drops `bake` in `$GOBIN`/`$GOPATH/bin`). `go run` for
  quick iteration.
- **Users today:** `go install github.com/alexpagnotta/bake-ai/cmd/bake@latest`,
  or download a prebuilt binary from GitHub Releases.
- **Future (brew):** a Homebrew tap formula just ships the prebuilt binary (or
  `go build` from source) — far simpler than any runtime-based tool. Automate the
  release binaries with GoReleaser.

## Docs/context loading — three tiers (reach for RAG last)
No embeddings needed:
1. **Stuff it in** — small docs/vault loaded whole via `.goosehints` / recipe.
2. **Agentic file reading** — let goose's developer extension `grep`/open only the
   files it needs on demand. Retrieval without embeddings; the sweet spot for a
   curated `docs/` folder. **Defer RAG indefinitely at this tier.**
3. **RAG (embeddings)** — only when the corpus is too big to fit *and* too big to
   navigate file-by-file. Build last, if ever.
Long-context models (e.g. Gemini via OpenRouter) just raise the tier-1 ceiling —
but you pay for every stuffed token each turn, so tier 2 usually wins even there.

## Shared risks / call-outs
- **Two TUIs can't share the terminal.** goose runs its own interactive UI. So
  `bake chat` must **hand the terminal to goose** for the chat itself, keeping
  Bubble Tea for everything *around* it (picker, forms, list, status). A fully
  custom chat screen (bake drives goose headlessly + renders with Glamour) is a
  later option that reimplements the chat loop — don't commit the MVP to it.
- **`.goosehints` loads entirely every request** → watch token cost as context
  grows; agentic file reading is the mitigation.
- **Tool/data leakage** → ship a `.gitignore` and keep the workspace path *outside*
  the tool repo by default; never scaffold projects inside the published repo.
- **Over-engineering** → build one project end-to-end before generalizing.

---

# Phase 1 — Engine + thin CLI

> **Status (2026-06-27): implemented & verified.** Repo scaffolded (Go + Cobra,
> module `github.com/alexpagnotta/bake-ai`); `internal/` engine (`config`,
> `workspace`, `gooserun`, `templates`); `init`/`new`/`list`/`chat` working;
> templates embedded; recipe validates; `.goosehints` context load proven
> (espresso-blend test); README + MIT LICENSE + `.goreleaser.yaml` added. Remaining
> user-side: add `~/go/bin` to PATH; optional GitHub remote + push.

**Goal:** prove the whole pipeline with the core engine behind a minimal command
skin. De-risk goose + OpenRouter + workspace separation before investing a single
line in UI. Output is minimal plain text; commands take args (no interactive
screens yet).

### Scope
Core engine (`internal/`, presentation-free, reused unchanged by Phase 2):
- `CreateProject`, `ListProjects`, `LaunchChat`, plus config/workspace resolution.

Thin CLI skin over the engine:
- `bake init` — set workspace path + OpenRouter config (first run).
- `bake new <project>` — scaffold a project folder from embedded templates.
- `bake list` — list workspace projects (minimal output).
- `bake chat <project>` — launch goose with that recipe, working dir = the project
  folder (so `.goosehints` + recipe load); hand the terminal to goose for the chat.

Under the hood: `goose run --recipe <workspace>/projects/<p>/recipe.yaml` with the
working directory set to the project folder.

### Build steps
1. Install goose; `goose configure` → OpenRouter + API key. Confirm a plain
   `goose run` works against an OpenRouter model.
2. Scaffold the public tool repo: `go mod init`, `cmd/bake/main.go`, a CLI router
   (e.g. Cobra). `git init`. `go install ./cmd/bake` and confirm `bake --help` runs
   from an unrelated directory.
3. Implement `bake init` (workspace path + non-secret global config in
   `~/.config/bake`); verify goose has an OpenRouter provider, else point the user
   to `goose configure`. bake stores no secrets.
4. Write `templates/` (recipe.yaml, .goosehints, skills/, docs/, vault/INDEX.md),
   embedded in the binary via `go:embed`. (Role/persona lives in recipe.yaml
   `instructions` + `.goosehints`; no separate prompt.md — goose wouldn't auto-read
   it. A managed prompt file is a V2 authoring feature.)
5. Implement the core engine in `internal/` (`CreateProject`, `ListProjects`,
   `LaunchChat`); wire `new`, `list`, `chat` as thin commands calling it (goose
   launched with recipe + workdir). Keep output minimal.
6. Write README + LICENSE so it's publishable; set up GoReleaser for binaries.

### Verification
1. `bake init` → workspace created at `~/bake`, non-secret config written, no
   secret stored by bake; goose OpenRouter config detected (or user guided to it).
2. `bake new coffee` → project scaffolded *in the workspace*, not the tool repo.
3. Add a fact to `coffee/.goosehints` or `coffee/docs/` (e.g. "our espresso blend is
   a lighter roast").
4. `bake chat coffee`; ask "what's our espresso blend?" → answered from the static
   context without re-explaining. ✅ Proves per-project context works.
5. Tool repo contains **no** project data and pushes cleanly to GitHub; the
   workspace can be its own separate (private) repo.

---

# Phase 2 — TUI polish (Charm)

**Goal:** turn the working CLI into the cool, animated experience as a second
front-end calling the **same** `internal/` core functions from Phase 1 — no core
logic rewritten, and the plain commands still work headless for scripting.

### Scope
- **App shell:** running `bake` with no args opens a Bubble Tea home — a project
  picker (Bubbles `list`) with Lip Gloss styling and animations.
- **`bake new`:** interactive **Huh** form (name, model, which scaffold pieces).
- **`bake chat`:** Bubbles spinner/viewport + status chrome *around* the handoff to
  goose (per the shared-risks note, goose still owns the chat screen itself).
- **Answer rendering / previews:** use **Glamour** to render any markdown bake
  itself shows (project summaries, `INDEX.md` previews, help).
- Consistent theme/styling layer (Lip Gloss) reused across screens.

### Build steps
1. Add an app-shell entry: no-arg `bake` launches the Bubble Tea program; existing
   subcommands remain callable directly for scripting.
2. Build the project picker (Bubbles `list`) from the workspace; selecting a project
   routes to chat.
3. Replace `bake new`'s plain prompts with a Huh form.
4. Add spinner/viewport + Lip Gloss chrome around `chat`; render markdown previews
   with Glamour.
5. Extract a shared theme/styles package; polish animations and transitions.

### Verification
1. `bake` (no args) → animated picker lists workspace projects; arrow-keys + enter
   select one and start its chat.
2. `bake new` → guided Huh form scaffolds the project.
3. Markdown that bake renders (summaries/help/INDEX preview) shows styled via
   Glamour.
4. All commands still work non-interactively with args (scripting unbroken).
