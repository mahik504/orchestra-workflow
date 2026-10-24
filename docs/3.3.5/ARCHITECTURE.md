# Orchestra 3.3.5 architecture

One conductor. Skills are capabilities. The registry is the resource list. The private vault is memory and is not part of this public repo.

```mermaid
flowchart TB
  prompt["plain request"] --> forge["prompt-forge"]
  forge --> conductor["conductor"]
  conductor --> skills["skills for this route"]
  conductor --> lanes["host lanes"]
  lanes --> verify["two checks then the running app"]
```

The design bar comes from project kind. Hiring work gets the full lab. Other showable work gets one approval. College work skips the lab unless asked.

A claim that this version beats a raw host waits for the 3H measurement.
