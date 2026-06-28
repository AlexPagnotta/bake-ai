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

```sh
# Installs bake + goose in one go
brew install alexpagnotta/tap/bake

# Point goose at a provider (OpenRouter recommended), then set up the workspace
goose configure          # Configure Providers → OpenRouter → paste your key
bake init
```

<details>
<summary>Install with Go instead</summary>

Requires [Go](https://go.dev/) and
[goose](https://block.github.io/goose/docs/getting-started/installation):

```sh
brew install block-goose-cli
go install github.com/alexpagnotta/bake-ai/cmd/bake@latest
```

Make sure `$(go env GOPATH)/bin` is on your `PATH`, then run `goose configure`
and `bake init`.

</details>

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

make build               # build everything (go build ./...)
make dev                 # run the CLI locally (go run ./cmd/bake)
make test                # run the tests (go test ./...)
```

### Releasing

Releases use [Semantic Versioning](https://semver.org/) via git tags and are cut
with [GoReleaser](https://goreleaser.com/) (see `.goreleaser.yaml`), which builds
the binaries, generates the GitHub Release notes, and updates the Homebrew tap.

```sh
git tag -a v1.2.0 -m "v1.2.0"
git push origin v1.2.0
export GITHUB_TOKEN=...   # PAT with write access to bake-ai and homebrew-tap
make release             # goreleaser release --clean
```

Use `make snapshot` to test a release build locally without publishing.

## License

[MIT](LICENSE)
