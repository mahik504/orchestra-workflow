# Custom skill (one named gap)

Use only after research found **no** graph/catalog match and the human said **update**.

Do not `npx skills add` a pack. Do not make this a 31st global skill unless they ask to promote it.

```
NAME: (kebab-case, one job)
PROBLEM: (one sentence)
WHEN TO LOAD: (trigger_conditions)
WHEN NOT TO: (avoid_conditions)
FALLBACK: (what to use if this fails)
HARNESS: cursor | antigravity | claude | jcode | any
STEPS:
1.
2.
VERIFY: (command or Playwright path)
```

Register: Brain overlay `memory/added-resources.json` first. Public `registries/resources.json` only if it should ship with the method.
