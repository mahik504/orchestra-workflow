---
name: visual-golden-still
description: Representative desktop and mobile stills. PASS/FAIL only. Load only on the visual-golden-still route. Does not implement the product.
---

# Golden still visual evidence

Job-loaded. Skip unless the `visual-golden-still` route fired, or a `DESIGN.md` is contract-approved and stills are the unpaid gate.

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
