# Benchmark classes

Nine task classes. Each `c*.md` file is the prompt shape. It is not a score.

| Id | Class | Record |
| --- | --- | --- |
| c1-premium-landing | premium landing | three trials, both arms built, judge 2–0 with one tie, claim none |
| c2-immersive-3d | one 3D scene | not run |
| c3-fullstack-saas | SaaS slice | not run |
| c4-fastapi-ai | FastAPI route | not run |
| c5-backend-only | backend change | not run |
| c6-bugfix | bug fix | three trials on `composer-2.5-fast`, both arms tied, claim none |
| c7-college | college assignment | not run |
| c8-diagram | architecture diagram | not run |
| c9-readme | README and post | not run |

`compare.go` enables a class only when a paired record exists. None are enabled.
