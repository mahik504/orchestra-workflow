---
name: ow-presentation
description: Slide story, hierarchy, type, speaker notes, and slide QA. Use with ow-docs for a deck. Host pptx builds the file. Not a purple template.
---

# Presentation

`ow-docs` owns the deliverable. This skill owns the story. Source methods: PPT Master and Anthropic's pptx skill. Neither tree is copied. The host `pptx` skill builds the file when that host has it.

## Before slides

1. Name the audience and the decision the deck asks them to make.
2. One idea per slide. A slide that needs a second idea is two slides.
3. Order: claim, evidence, consequence. Speaker notes carry the sentence the slide does not.

## On the slide

- One type family for titles, one for body. Hierarchy is size and weight, not a new color per box.
- Charts and diagrams use the same names as the source note. Do not invent a metric.
- No purple gradient template, no icon row that repeats the title, no "proud to announce."

## QA

Read the deck once as the audience. Every slide answers "why is this here." A slide that only decorates is cut. Speaker notes say what the mouth will add. The type sizes still work when the deck is projected, not only on a laptop.

## Preconditions

- The audience and the decision are named. `ow-docs` owns the file.

## Procedure

1. The task's paths and acceptance line win.
2. One idea per slide: claim, evidence, consequence. Speaker notes carry the spoken sentence.
3. Cut a slide that only decorates. Do not invent a metric.

## Output

The slide list and the speaker notes.

## Verify

Read the deck once as the audience. A slide that does not change a decision is cut.

## Skip

Skip for a report that is not a deck, and for a purple template used as the design.
