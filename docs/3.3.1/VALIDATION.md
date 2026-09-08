# 3.3.1 validation

Run from the public clone root:

```
cd runtime && go test ./...
node runtime/tools/hygiene.js
node runtime/tools/check-registry.js
node runtime/tools/check-host-stack.js
```

Expect:

- Classifier capability count **17**
- Contract approval does **not** unlock product `page.tsx`
- Stills approval does
- Named-reference survey accepts **exactly 3** and refuses 2 or 4
- No files under a product `Portfolio` tree
- No GSD/ECC clones in this repo
- No secrets

Do not claim HEALTHY without host evidence. Do not tag until the operator says so.
