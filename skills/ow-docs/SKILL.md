---
name: ow-docs
description: College/product slides, reports, and papers. Use when the user needs a PPTX, PDF report, README, or academic paper. Not for app UI. Prefer markdown then export. Do not fake citations.
---

# Docs, slides, papers

Do not install Anthropic `frontend-design`. Document skills only.

## What to use

| Output | How |
| --- | --- |
| Slides | Write content in markdown first. Build `.pptx` with `python-pptx` or `pptxgenjs`. Few slides, one idea each, real hierarchy. Optional later: Anthropic [pptx skill](https://github.com/anthropics/skills/tree/main/skills/pptx) inside Claude/Antigravity (proprietary — do not copy the whole anthropics/skills repo). |
| Course report / README | Markdown in the repo. Honest, short. MarkItDown (`microsoft/markitdown`) can convert a document you already have into markdown. It is a reference, not a skill in this repo. |
| Research paper | LaTeX + `.bib`. Real citations only. Conference template from the venue. Thin MIT reference: [ow-research-paper-writing](https://github.com/nousresearch/hermes-agent/blob/main/skills/research/ow-research-paper-writing/SKILL.md) as **method**, not a 12-agent install. |

## Taste

Same as `Preferences.md`: no purple AI-SaaS slide templates, no fake NeurIPS claims, no invented references.

## Never

- Dump `anthropics/skills` or academic 12-agent packs into Cursor.
- Generate a bibliography you did not read.

## Preconditions

- The task names the deliverable: slides, a report, a README, or a paper.

## Procedure

1. The task's paths and acceptance line win.
2. Write the content in markdown first. Export only after the structure is real.
3. Every citation is one you opened. A missing citation is a gap, not a guessed DOI.

## Output

The markdown source and, if asked, the exported file.

## Verify

Open each citation you added. A citation you did not open comes out.

## Skip

Skip for app UI. A deck's story is `ow-presentation`. A paper's method is `ow-research-paper`.
