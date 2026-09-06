# Capability graph

Thirteen routes in `registries/design-resource-graph.json`. Every row has a trigger **and** a skip.

```mermaid
flowchart LR
  brief[Brief]
  match[Trigger and skip]
  one[One capability]
  load[Load that route only]
  miss[No match]
  research[Research then wait for update]
  brief --> match
  match -->|hit| one --> load
  match -->|miss| miss --> research
```

`interaction-components` searches registered kits for one control. It is not a skill per button library.
