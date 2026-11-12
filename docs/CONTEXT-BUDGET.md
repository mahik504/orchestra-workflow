# Context budget

Measured by `node runtime/tools/measure-context.mjs` on this tree.

| What | Characters |
| --- | ---: |
| AGENTS.md | 2047 |
| Routing rows as JSON | 4380 |
| Skill descriptions, all canonical skills | 25614 |
| Skill descriptions, ow-* only | 7303 |
| ow-* descriptions over 300 | 0 |

The 4.2.0 tree had 78 canonical skills. This tree's count is in `registries/skills.json`. The slot cap of 80 is gone. The budget is the description length: 300 characters for an `ow-*` skill.

## Disclosure

| Level | What the worker sees |
| --- | --- |
| L0 | task and acceptance only |
| L1 | route row from registries/routing.json |
| L2 | selected SKILL.md |
| L3 | references/ beside that skill |
| L4 | tool output from the run |

A STANDARD job that names the file stays at L0. The conductor fast path does not print a brief or a receipt. A visual job starts at L1 and opens L2 only for the skills the route loaded.

## Heavy routes

| Route | Skills | SKILL.md bytes |
| --- | ---: | ---: |
| Award-style site | 7 | 56621 |
| Full-stack app | 4 | 18698 |
| Mobile app | 5 | 39661 |

These sizes are a measurement. They are not a cap.
