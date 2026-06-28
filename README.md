<h1 align="center">BAKE</h1>

<p align="center">A personalized, project-aware AI assistant for your terminal.</p>

<p align="center">
  <img src="assets/cover.jpg" alt="BAKE" width="640" />
</p>

## About

`bake` is a thin, provider-agnostic wrapper around
[goose](https://block.github.io/goose/) that gives each of your projects its own
prompt, docs, and accumulated context — so you never start a chat from zero.

It's built for research and personal projects (not coding): **no provider
lock-in, no subscription, and no living inside the ChatGPT or Gemini apps.**

A few ideas keep it simple:

- **The tool is public; your data is private.** This repo holds only code and
  templates. Your projects live in a separate **workspace** (default `~/bake`),
  which you can make your own private git repo.
- **A project is just a folder** — a goose recipe with its own always-loaded
  context (`.goosehints`), curated reference material (`docs/`), reusable
  playbooks (`skills/`), and notes (`vault/`).
- **Secrets stay in goose.** Your OpenRouter API key lives in your OS keyring via
  goose — `bake` never reads, stores, or logs it.

### What a project looks like

```
coffee/
  recipe.yaml     # how the assistant is wired (model, role, tools)
  .goosehints     # always-loaded context
  docs/           # curated reference, read on demand
  skills/         # reusable task playbooks, read on demand
  vault/          # accumulated notes & decisions (INDEX.md)
```

`recipe.yaml` and `.goosehints` are loaded on **every** message, while `docs/`,
`skills/`, and `vault/` are read **on demand** — the assistant opens them only
when relevant, so they don't cost tokens every turn.

## Installation

`bake` needs [Go](https://go.dev/) and
[goose](https://block.github.io/goose/docs/getting-started/installation).

```sh
# 1. Install goose and point it at a provider (OpenRouter recommended)
brew install block-goose-cli
goose configure          # Configure Providers → OpenRouter → paste your key

# 2. Install bake
go install github.com/alexpagnotta/bake-ai/cmd/bake@latest
```

Make sure `$(go env GOPATH)/bin` is on your `PATH`, then run `bake init`.

## Usage

Run `bake` with no arguments to open the interactive home screen, or use the
subcommands directly:

```sh
bake init                # set up the workspace + non-secret config
bake new coffee          # scaffold a new project (interactive form if no name)
bake list                # list projects in your workspace
bake chat coffee         # start a session that already knows your project
```

| Command       | What it does                                                            |
| ------------- | ----------------------------------------------------------------------- |
| `bake init`   | Creates the workspace and config, and checks for goose + OpenRouter.    |
| `bake new`    | Scaffolds a project from the templates. Omit the name for a guided form.|
| `bake list`   | Prints every project found in the workspace.                            |
| `bake chat`   | Hands the terminal to goose, scoped to the project's context.           |

A typical first run:

```sh
bake init
bake new coffee
# add context in ~/bake/projects/coffee/.goosehints
bake chat coffee
```

Run `bake --help` (or `bake <command> --help`) for all flags.

## Development

Clone the repo and build from source:

```sh
git clone https://github.com/alexpagnotta/bake-ai.git
cd bake-ai

go build ./...           # build everything
go run ./cmd/bake        # run the CLI locally
go test ./...            # run the tests
```

Releases are cut with [GoReleaser](https://goreleaser.com/) (see
`.goreleaser.yaml`).

## License

[MIT](LICENSE)
