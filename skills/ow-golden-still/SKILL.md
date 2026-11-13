---
name: ow-golden-still
description: Representative desktop and mobile stills. PASS/FAIL only. Load only on the ow-golden-still route. Does not implement the product.
---

# Golden still visual evidence

Job-loaded. Skip unless the `ow-golden-still` route fired, or a `DESIGN.md` is contract-approved and stills are the unpaid gate.

A `DESIGN.md` is not visual evidence.

## Skip

- The lab was waived with “skip the lab”
- STANDARD backend, paper, or security with no UI
- Stills already human-approved and the job is implementation
- Full product build with no stills gate named — classify the product route first

## Produce

Under `.orchestra/design-lab/golden-stills/`:

- One representative **desktop** still
- One representative **mobile** still
- Optional throwaway renderer in that folder only

Verdict is **PASS** or **FAIL**. No product `src/` / `app/` writes. Missing files or a silent skip refuse `ApproveStills`.

## Preconditions

- The task names the screen and the acceptance line. A `DESIGN.md` is not the still.

## Procedure

1. The task's paths and acceptance line win.
2. Capture one desktop still and one mobile still under `.orchestra/design-lab/golden-stills/`.
3. Mark PASS or FAIL. Do not write product `src/` from this skill.

## Output

The two still paths and the PASS or FAIL line.

## Verify

Open both files. A missing file is FAIL.

## Skip

Skip when the lab was waived, the job has no UI, or the stills are already approved and the job is implementation.
