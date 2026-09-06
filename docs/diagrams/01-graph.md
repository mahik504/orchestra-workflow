# Graph

Capability routes live in `registries/design-resource-graph.json`. Thirteen capabilities. Every row has a trigger **and** a skip.

```mermaid
flowchart LR
  brief[Brief]
  match[Cheap trigger and skip]
  cap[One capability]
  load[Load that route only]
  brief --> match --> cap --> load
```

Available ≠ loaded. A miss → research → one named capability → wait for **update**.
