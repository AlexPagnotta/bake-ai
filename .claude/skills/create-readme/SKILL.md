---
name: create-readme
description: Generate the project's README.md. Use when asked to create, write, or regenerate the README.
---

# Create README

Generate a `README.md` for this project in **GitHub-flavored Markdown (GFM)**.

## Step 1 — Review the project first

Before writing anything, understand the project:

- Read the manifest / entry points (e.g. `go.mod`, `package.json`, `cmd/`, `main.*`)
  to learn what it is, what it does, and how it's run.
- Determine the **project type** — website, CLI tool, desktop/mobile app, library —
  because it decides which sections apply (see below).
- Look for an existing logo or screenshot in the repo (e.g. `assets/`, `docs/`,
  `.github/`, image files like `*.png`/`*.jpg`/`*.svg`) to feature near the top.
- Note the install method, available commands/entry points, and the dev setup
  steps (build, run, test).

Only write the README once you actually understand the project — don't guess.

## Step 2 — Structure

Build the README in this order. Include a section only when it applies to the
project type; skip the ones that don't.

1. **Title** — the project name, centered on the page.
2. **One-line description** — a short tagline directly below the title, also centered.
3. **Image** — a logo or an app screenshot, centered below the description. Use one
   found in the repo; if none exists, leave a clearly-marked placeholder rather than
   inventing a path.
4. **About the project** — a quick, plain-language explanation of what the project
   is and the problem it solves. Keep it non-technical: no implementation details,
   no internal architecture.
5. **Installation** — *only if the project is something a user installs* (a CLI
   tool, desktop/mobile app, etc. — not a hosted website). List the concrete steps
   to install it.
6. **Usage** — *only for CLI tools or similar*. List the available commands with a
   short explanation of each. If there are many, cover the main ones and point to
   `--help` for the rest.
7. **Setup** — how to set the project up for **local development**: clone, install
   dependencies, build, run, and test.

## Markdown conventions

- Use GFM throughout — headings (`##`), fenced code blocks with language hints
  (```` ```sh ````), tables, and lists.
- Centering: use a small HTML block, since Markdown can't center on its own, e.g.

  ```html
  <h1 align="center">Project Name</h1>
  <p align="center">A one-line description of what it does.</p>
  <p align="center">
    <img src="assets/cover.jpg" alt="Project Name" width="600" />
  </p>
  ```

- Keep commands copy-pasteable and accurate to what's actually in the repo.
- Match the project's real names, paths, and commands — verify against the code,
  don't assume.

## Notes

- If a `README.md` already exists, treat this as a rewrite: keep any still-accurate,
  project-specific content (badges, license, links) rather than discarding it.
- Keep the tone clear and friendly; favor short sections over walls of text.
