# System figures

GitHub renders these from the implementation. The three README figures are SVGs in `assets/`.

## Control plane

```mermaid
flowchart TD
  human[Human]
  conductor[One conductor]
  graph[Capability graph]
  note[Routing note]
  hosts[Host adapters]
  brain[Private Brain]
  human --> conductor
  conductor --> graph
  conductor --> note
  conductor --> hosts
  conductor --> brain
```

## Eight stages

```mermaid
flowchart LR
  d[Discover]
  c[Classify]
  r[Research]
  s[Synthesize]
  design[Design]
  gate[Human gate]
  impl[Implement]
  v[Verify]
  mem[Remember]
  d --> c --> r --> s --> design --> gate --> impl --> v --> mem
```

## Mission

```mermaid
flowchart LR
  mission[Mission]
  capsule[Worker capsule]
  budget[Budget]
  run[Execution]
  evidence[Evidence]
  fail[Failure class]
  recover[Recovery or human gate]
  mission --> capsule --> budget --> run --> evidence --> fail --> recover
```

## What a task loads

```mermaid
flowchart TD
  brief[Brief]
  route[One capability]
  domain[Design domain]
  stack[Technology rules]
  note[Load MCP plugins libraries]
  brief --> route --> domain
  route --> stack --> note
```

## Design gate

```mermaid
flowchart TD
  cards[Directions or three translations]
  contract[DESIGN.md]
  stills[Golden stills]
  code[Implementation]
  cards --> contract --> stills --> code
```

## Scene

```mermaid
flowchart TD
  still[Local still]
  canvas[One canvas]
  r3f[React Three Fiber]
  three[Three.js renderer]
  dispose[Dispose on unmount]
  still --> canvas --> r3f --> three --> dispose
```

## Data store

```mermaid
flowchart TD
  need{Durable data?}
  none[No database]
  pg[Postgres]
  sb{Product names Supabase?}
  supa[Supabase]
  redis{Cache TTL or rate limit?}
  cache[Redis]
  need -->|no| none
  need -->|yes| pg
  pg --> sb
  sb -->|yes| supa
  sb -->|no| redis
  redis -->|yes| cache
```

## Hosts

```mermaid
flowchart LR
  repo[orchestra-workflow]
  sync[sync-ides]
  cursor[Cursor]
  anti[Antigravity]
  claude[Claude Code]
  jcode[jcode]
  repo --> sync
  sync --> cursor
  sync --> anti
  sync --> claude
  sync --> jcode
```

Limitless Claude is a separate repository. It routes model providers. It is not a fifth copy of this contract.
