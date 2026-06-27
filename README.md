# bake-ai

A personalized, project-aware AI assistant for the terminal — for research and
personal projects, not coding. **No provider lock-in, no subscription, no living
inside the ChatGPT/Gemini apps.**

`bake` is a thin, provider-agnostic wrapper around [goose](https://goose-docs.ai/)
that gives each of your projects its own prompt, docs, and context — so you never
start a new chat from zero.

## How it works

- **The tool is public; your data is private.** `bake` (this repo) holds only code
  and templates. Your projects live in a separate **workspace** (default `~/bake`),
  which you can make your own private git repo.
- **A project = a folder = a goose recipe**, with its own `.goosehints` (always-loaded
  context), `docs/` (curated reference), and `vault/` (notes).
- **Secrets stay in goose.** The OpenRouter API key lives in your OS keyring via
  goose — `bake` never reads, stores, or logs it.

## Install

Requires [Go](https://go.dev/) and [goose](https://formulae.brew.sh/formula/block-goose-cli):

```sh
brew install block-goose-cli
goose configure          # Configure Providers → OpenRouter → paste key
go install github.com/alexpagnotta/bake-ai/cmd/bake@latest
```

Make sure `$(go env GOPATH)/bin` is on your `PATH`.

## Usage

```sh
bake init                # set up workspace + non-secret config
bake new coffee          # scaffold a project
# edit ~/bake/projects/coffee/.goosehints to give it context
bake chat coffee         # start a session that already knows your project
bake list                # list projects
```

## Status

MVP (Phase 1): project scaffolding + per-project context + chat. See `PLAN.md`.
Auto-updating vault and a Charm TUI are planned — see `V2.md`.

## License

MIT
