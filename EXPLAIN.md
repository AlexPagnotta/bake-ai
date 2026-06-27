# EXPLAIN.md — How the Go code works

This document walks through every Go file in `bake`, explaining what it does and
how the pieces fit together. It's meant as a learning map of the codebase.

## The big picture

`bake` is a small CLI written in Go. Its job is **not** to talk to an LLM itself —
it's a thin wrapper around the external [`goose`](https://goose-docs.ai/) tool.
`bake` manages *projects* (folders, each with their own prompt/docs/context) and,
when you start a chat, it hands the terminal over to `goose` pointed at that
project's `recipe.yaml`.

The code is split into two layers, which is the single most important idea to
grasp:

| Layer | Location | Responsibility |
| --- | --- | --- |
| **Commands / UI** | `cmd/bake/` (CLI) and `internal/tui/` (interactive screens) | Parse input, show output. **No core logic.** |
| **Engine** | `internal/{config,workspace,gooserun,templates}/` | All the real work: read config, create/list projects, launch goose. Presentation-free. |

Because the engine is presentation-free, the CLI (`bake new foo`) and the
interactive TUI (the "＋ New project" screen) both call the *same* functions and
behave identically.

Dependency direction (who imports whom):

```
cmd/bake (commands)  ─┬─► internal/tui ──┐
                      │                  ├─► internal/workspace ─► internal/templates
                      ├─► internal/gooserun ─┘                     internal/config
                      └─► internal/config
```

Third-party libraries: **cobra** (CLI commands/flags), **Bubble Tea + Bubbles +
Lipgloss + Huh + Glamour** (the Charm TUI stack), and **BurntSushi/toml** (config
file parsing).

---

## `cmd/bake/` — the command layer

This is `package main`, the executable. Each file defines one cobra command.

### `main.go`
The entry point. Tiny on purpose:
- Declares `version = "dev"`, which is overwritten at build time via `-ldflags`
  (see `.goreleaser.yaml`) so released binaries report a real version.
- `main()` just calls `Execute()` (defined in `root.go`).

### `root.go`
Defines `rootCmd` — the base `bake` command — and wires everything together.
- `Execute()` runs the root command and is what `main` calls.
- `init()` registers the four subcommands: `init`, `new`, `list`, `chat`.
- `runHome()` is what runs when you type `bake` with **no arguments**. This is the
  heart of the interactive mode: an infinite loop that
  1. shows the project picker (`tui.RunPicker`),
  2. reacts to the chosen `Action`:
     - `ActionQuit` → return (exit),
     - `ActionNew` → show the new-project form, then `workspace.Create`,
     - `ActionChat` → look up the project, check goose is installed, print a
       header, and `gooserun.LaunchChat`.
  3. loops back to the picker after a chat session ends.
- Note the careful error handling around `LaunchChat`: a non-zero exit (e.g. you
  pressed Ctrl+C to leave goose) is treated as *normal* and returns you to the
  picker; only genuine launch failures are printed.

### `init.go`
The `bake init` command — first-time setup.
- If config already exists, it just prints the current settings.
- Otherwise it builds a `Default()` config, applies optional `--workspace` and
  `--model` flags, creates the projects directory, and saves the config.
- Finally it does **secret-free health checks**: is `goose` on PATH? Does goose
  have an OpenRouter provider configured? It only *checks* — it never stores keys.
- The two flags are registered in this file's own `init()` function.

### `new.go`
The `bake new [project]` command — scaffolds a project.
- If you pass a name (`bake new coffee`), it stays headless/scriptable.
- If you don't, it pops the interactive Huh form (`tui.NewProjectForm`) to collect
  name + model.
- Either way it calls `workspace.Create` and prints next steps.

### `list.go`
The `bake list` command. Calls `workspace.List` and prints one project name per
line (or a friendly "no projects yet" hint). Pure plumbing.

### `chat.go`
The `bake chat <project>` command — the non-interactive way to start a session.
- Loads config, resolves the project with `workspace.Get`, checks goose is
  installed, prints the styled header, and calls `gooserun.LaunchChat`.
- This is essentially the `ActionChat` branch of `runHome`, exposed as a direct
  command.

---

## `internal/config/` — non-secret settings

### `config.go`
Owns `~/.config/bake/config.toml` (or `$XDG_CONFIG_HOME/bake`). **Secrets are
never stored here** — goose owns the OpenRouter key in the OS keyring.
- `Config` struct: just two fields, `WorkspacePath` and `DefaultModel`, with toml
  tags so they map to the file.
- `Default()` — baseline config on first run (`~/bake` workspace,
  `google/gemini-2.5-flash` model).
- `Dir()` / `FilePath()` — resolve where the config lives (honoring XDG).
- `Exists()` — does the config file exist yet?
- `Load()` — decode the TOML; returns a helpful "run `bake init` first" error if
  missing.
- `Save()` — write the file (dir `0755`, file `0600`).

---

## `internal/workspace/` — core project operations

### `workspace.go`
The engine's core. **Presentation-free** so CLI and TUI share it.
- `nameRe` / `ValidName()` — project names must be lowercase letters, digits, `-`
  or `_`. Used both to validate on create and to drive the form's live validation.
- `Project` struct — a project is just a `Name` + `Path`.
- `ProjectsDir()` — `<workspace>/projects`.
- `Create()` — validates the name, defaults the model, errors if the folder
  already exists, then renders the embedded template tree via
  `templates.RenderProject`.
- `List()` — reads the projects directory and returns the folders that contain a
  `recipe.yaml` (that file is the marker of "this is a real project"), sorted by
  name. A missing directory yields an empty list, not an error.
- `Get()` — resolve a single project by name, erroring if its `recipe.yaml` is
  missing.

---

## `internal/templates/` — the project scaffold

### `embed.go`
Turns the files under `templates/project/` into a new project on disk, using Go's
`embed` to bake the template tree into the binary (no external files needed at
runtime).
- `//go:embed all:project` — compiles the whole `project/` tree into `projectFS`.
- `ProjectData` — the values substituted into templates: `Name` and `Model`.
- `RenderProject()` — walks the embedded tree; recreates each directory and
  renders each file into the destination.
- `renderFile()` — reads a template, parses it with `text/template`, and writes
  the executed result (with `ProjectData` substituted in).
- `mapName()` — turns a template path into its final on-disk name: strips the
  `.tmpl` suffix, and special-cases `goosehints.tmpl` → `.goosehints` (Go's embed
  can't cleanly ship dotfiles, so they're stored without the leading dot).

### The template files (not Go, but rendered by the above)
- `project/recipe.yaml.tmpl` — the goose recipe (model, role, tools). `{{.Model}}`
  gets filled in.
- `project/goosehints.tmpl` → `.goosehints` — always-loaded context for the assistant.
- `project/docs/README.md.tmpl` — starter for curated reference docs.
- `project/skills/README.md.tmpl` — starter for reusable task playbooks.
- `project/vault/INDEX.md.tmpl` — notes index (also shown as a preview in the chat header).

---

## `internal/gooserun/` — launching & inspecting goose

### `gooserun.go`
The bridge to the external `goose` binary. **bake never reads, stores, or logs
secrets** — goose owns the API key.
- `LaunchChat()` — runs `goose run --recipe <project>/recipe.yaml --interactive`,
  with the working directory set to the project (so goose picks up `.goosehints`,
  `docs/`, `vault/`). It wires the child process's stdin/stdout/stderr to bake's,
  effectively handing the terminal to goose until the session ends. While goose
  runs, bake **ignores Ctrl+C (SIGINT)** — goose and bake share the terminal, so
  the interrupt reaches both; letting goose own it means pressing Ctrl+C ends the
  goose session and returns control to bake (back to the picker) instead of also
  killing bake. The original signal behavior is restored when the session ends.
- `Installed()` — is the `goose` binary on PATH? (`exec.LookPath`).
- `HasOpenRouter()` — a best-effort, **secret-free** scan of goose's
  `config.yaml` for the word "openrouter", just to give a helpful hint during
  `bake init`. Returns false if it can't tell. Finds the config via
  `gooseConfigPath()`, which honors `XDG_CONFIG_HOME` before `~/.config/goose`.

---

## `internal/tui/` — the Charm / Bubble Tea front-end

A thin presentation layer over the engine. No core logic lives here.

### `theme.go`
The package doc comment (explaining the layer separation) lives here, plus the
visual identity: a warm "toasted amber" `accent` palette and a set of shared
Lipgloss `Style` values (`titleStyle`, `panelStyle`, `labelStyle`, etc.) reused
across every screen.

### `picker.go`
The interactive project list — the home screen — built on the Bubbles `list`
component, following Bubble Tea's **Model/Update/View** architecture.
- `pickerItem` — one row; either a real project (`kindProject`) or the
  "＋ New project" entry (`kindNew`). Implements the `list.Item` interface
  (`Title`, `Description`, `FilterValue`).
- `pickerModel` — holds the `list.Model` and the `Result` to return.
- `newPickerModel()` — builds the item list (projects + the "new" row), styles the
  selected row with the accent color, and enables fuzzy filtering.
- `Init()` — nothing to do up front.
- `Update()` — the event loop: resizes the list on window-resize; on key presses
  handles `q`/`ctrl+c`/`esc` (quit), `n` (new), and `enter` (chat or new depending
  on the selected row). Crucially, while the user is *typing a filter* it lets the
  list own every keystroke so letters aren't hijacked as shortcuts.
- `View()` — renders the list with a margin.

### `tui.go`
The public surface the command layer calls into — the glue functions:
- `Action` / `Result` — the small types describing what the user chose
  (`ActionQuit` / `ActionChat` / `ActionNew`, plus an optional project name).
- `RunPicker()` — runs the picker program in the alt screen and returns the
  chosen `Result`.
- `NewProjectForm()` — a **Huh** form collecting project name (with live
  validation via `workspace.ValidName`) and model (a select with the workspace
  default plus a few presets). Returns `ok = false` if the user aborts.
- `PrintChatHeader()` — prints the styled banner shown right before goose takes
  over: project name, path, model, and a Glamour-rendered preview of
  `vault/INDEX.md`.
- `recipeModel()` — best-effort regex read of `goose_model:` out of the project's
  `recipe.yaml` (for the header).
- `renderVaultPreview()` — renders `vault/INDEX.md` through Glamour (markdown →
  styled terminal output), or returns `""` if it can't.

---

## A typical flow, end to end

1. You run `bake` (no args) → `main.main` → `Execute` → `root.runHome`.
2. `runHome` calls `config.Load`, then loops on `tui.RunPicker`.
3. You pick a project and press Enter → `Result{Action: ActionChat, Project: ...}`.
4. `runHome` calls `workspace.Get`, checks `gooserun.Installed`, prints
   `tui.PrintChatHeader`, then `gooserun.LaunchChat`.
5. goose runs interactively using the project's `recipe.yaml` and context files.
6. You Ctrl+C out of goose → control returns to `runHome`, which loops back to the
   picker.

Creating a project follows the same shape: picker → `tui.NewProjectForm` →
`workspace.Create` → `templates.RenderProject` writes the embedded scaffold to
`<workspace>/projects/<name>/`.
