---
name: orchestra-docs
description: College/product slides, reports, and papers. Use when the user needs a PPTX, PDF report, README, or academic paper. Not for app UI. Prefer markdown then export. Do not fake citations.
---

# Docs, slides, papers

Do not install Anthropic `frontend-design`. Document skills only.

## What to use

| Output | How |
| --- | --- |
| Slides | Write content in markdown first. Build `.pptx` with `python-pptx` or `pptxgenjs`. Few slides, one idea each, real hierarchy. Optional later: Anthropic [pptx skill](https://github.com/anthropics/skills/tree/main/skills/pptx) inside Claude/Antigravity (proprietary — do not copy the whole anthropics/skills repo). |
| Course report / README | Markdown in the repo. Honest, short. |
| Research paper | LaTeX + `.bib`. Real citations only. Conference template from the venue. Thin MIT reference: [research-paper-writing](https://github.com/nousresearch/hermes-agent/blob/main/skills/research/research-paper-writing/SKILL.md) as **method**, not a 12-agent install. |

## Taste

Same as `Preferences.md`: no purple AI-SaaS slide templates, no fake NeurIPS claims, no invented references.

## Never

- Dump `anthropics/skills` or academic 12-agent packs into Cursor.
- Generate a bibliography you did not read.
