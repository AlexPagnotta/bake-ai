# Project structure

A project is just a folder. `bake new <name>` scaffolds it from a set of
templates, and every piece has a clear job:

```
coffee/
  recipe.yaml     # how the assistant is wired (model, role, tools)
  .goosehints     # always-loaded context
  docs/           # curated reference, read on demand
  skills/         # reusable task playbooks, read on demand
  vault/          # accumulated notes & decisions (INDEX.md)
```

## What's loaded, and when

The most important thing to understand is *when* each piece reaches the model:

- **Loaded on every message** — `recipe.yaml` and `.goosehints`. They're always
  in context, so they shape every reply but cost tokens every turn. Keep them
  tight.
- **Read on demand** — `docs/`, `skills/`, and `vault/`. The assistant opens
  these only when a request makes them relevant, so they can be as large as you
  like without taxing routine chats.

Rule of thumb: **always-known and short → `.goosehints`; everything bigger or
occasional → `docs/`, `skills/`, or `vault/`.**

## The pieces

- **`recipe.yaml`** — the goose recipe: the model and provider, the assistant's
  role (its `instructions`), and which tools are enabled. This is *how the agent
  is wired*, not what it knows. Edit it to change the model, the persona, or the
  available tools.
- **`.goosehints`** — durable, high-value facts the assistant should *always*
  know: who you are, the project's core facts, your preferred tone, and pointers
  to the rest. Loaded in full on every message, so keep it lean.
- **`docs/`** — larger or occasional *reference* material: product lists,
  pricing, specs, pasted notes or PDFs. Too big or too situational to always
  load, so the assistant reads only what it needs.
- **`skills/`** — *how-to* playbooks for recurring tasks (e.g. "write tasting
  notes", "draft marketing copy"): the method and format you want applied
  consistently.
- **`vault/`** — the project's *memory over time*: decisions and topic notes,
  indexed by `INDEX.md`. Hand-maintained today; auto-updated after each chat on
  the roadmap.

## Which one do I use?

| You want to…                                   | Put it in…    |
| ---------------------------------------------- | ------------- |
| Change the model, persona, or enabled tools    | `recipe.yaml` |
| State a fact the assistant must *always* know  | `.goosehints` |
| Add big or situational reference material      | `docs/`       |
| Define how a repeatable task should be done    | `skills/`     |
| Record a decision or outcome to build on later | `vault/`      |
