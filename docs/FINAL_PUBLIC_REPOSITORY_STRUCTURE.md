# Final Public Repository Structure (Orchestra V3)

Target: `mahik504/orchestra-workflow`

```text
orchestra-workflow/
├── .github/
│   └── workflows/
│       └── ci.yml                (GitHub Actions CI for Go runtime & schemas)
├── docs/
│   └── (Public documentation & benchmark methodologies)
├── kit/
│   └── (Orchestration setup utilities)
├── protocols/
│   └── (Public workflow protocols)
├── registries/
│   └── (Public capability & adapter mappings)
├── runtime/                      (Core V3 Engine - Pure Go)
│   ├── cmd/
│   │   └── orchestra/
│   │       └── main.go           (Compiled orchestra CLI binary source)
│   ├── internal/
│   │   ├── adapters/             (Verification and external runners)
│   │   ├── classifier/           (Task classification engine)
│   │   ├── handoff/              (SHA256 versioned state protocol)
│   │   ├── kernel/               (Core execution context)
│   │   ├── memory/               (Ponytail management)
│   │   ├── planner/              (Execution manifesto builder)
│   │   ├── resources/            (Capability metadata loader)
│   │   ├── router/               (4-Stage Capability Router)
│   │   └── verify/               (Verification gates)
│   ├── schemas/                  (JSON Schemas for Task, Plan, Handoff, etc.)
│   ├── tests/                    (Benchmark protocols and results)
│   ├── .gitignore                (Excludes go-portable binary layer)
│   └── go.mod                    (Go module root)
├── scripts/
│   └── (Public installation & utility scripts)
├── skills/
│   └── (Publicly accessible skill definitions)
├── templates/
│   └── Preferences.template.md   (Clean user preferences template)
├── workspace-template/
│   └── (Clean `.orchestra/` init template)
├── .gitignore                    (Global git exclusions)
├── AGENTS.md                     (Supported agents documentation)
├── ARCHITECTURE.md               (V3 architectural diagram & explanation)
├── CHANGELOG.md                  (Version history)
├── CONTRIBUTING.md               (Contribution guide)
├── LICENSE                       (MIT)
├── README.md                     (Comprehensive 20-section Production Documentation)
├── STACK.md                      (Supported technology stack)
├── VERSION                       (Contains exact string: 3.0.0)
└── WORKFLOW.md                   (End-to-end V3 workflow process)
```

## Privacy Assurance
The above tree represents the exact final state committed in V3. 
**Verification Status:** Confirmed that `memory/`, `.obsidian/`, `Mahi.plan/`, all private project notes (`TTB Agro`, `Lyra`, etc.), and stale V2 onboarding instructions (`START.md`, `STITCH.md`, `CONNECT.md`) have been completely eradicated from the public repository.
