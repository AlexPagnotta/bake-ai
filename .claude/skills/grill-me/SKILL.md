---
name: grill-me
description: Interrogate the user about a plan or design, one question at a time, walking the design tree and resolving decisions until a shared understanding is reached. Use when planning a feature or making design decisions and the user wants to be grilled / pushed to think it through.
---

# Grill me

Interview the user relentlessly about every aspect of this plan until you reach a
shared understanding. Walk down each branch of the design tree, resolving
dependencies between decisions one-by-one. For each question, provide your
recommended answer.

Ask the questions **one at a time**, waiting for feedback on each question before
continuing. Asking multiple questions at once is bewildering.

If a question can be answered by exploring the codebase, explore the codebase
instead of asking.

## How to run it

- Start from the highest-leverage, most upstream decision — the one other choices
  depend on — and work downward. Don't jump around the tree randomly.
- One question per turn. After each answer, integrate it and let it reshape the
  remaining branches before asking the next.
- For every question, state **your recommended answer** and a one-line why, so the
  user can just confirm or push back rather than starting from a blank page.
- Before asking, check whether the codebase already answers it (search/read the
  relevant files). Only ask about things that genuinely require the user's intent,
  taste, or external knowledge.
- Surface dependencies explicitly: "this depends on the previous answer; given you
  chose X, the options now are…".
- Keep going until the design is pinned down — no significant open questions, no
  hidden assumptions. Then summarize the agreed plan back in a short recap so the
  shared understanding is explicit.

## Notes

- Prefer concrete, decision-shaped questions over vague ones. "Should sessions
  persist across restarts, or start fresh each time?" beats "How should sessions
  work?".
- It's fine to be blunt and challenge weak answers — the point is to stress-test
  the plan, not to rubber-stamp it.
- If a decision is genuinely 50/50 or low-stakes, say so, pick the recommended
  default, and move on rather than belaboring it.
