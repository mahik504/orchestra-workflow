# Research review 5.5

Recorded 2026-11-12. None of these resolve rates, pass rates, or cost deltas are an Orchestra target. A paper’s score belongs to that paper’s model, harness, and task set.

Two layers:

- **VERIFIED.** The current arXiv HTML was fetched in this session and the quoted numbers were read from that fetch.
- **REPORTED.** [Papers](73faad4f-9ccd-4a2f-be4b-b583487ab68d) said they opened the page. This writing pass did not re-open it. The URL is included so the sentence can be checked. Do not promote a REPORTED number into a requirement.

Full PDFs were not read, except where the notes say otherwise. The SWE-Skills v1 PDF was not opened here.

## VERIFIED

### Evaluating AGENTS.md

- **Title.** Evaluating AGENTS.md: Are Repository-Level Context Files Helpful for Coding Agents?
- **Authors.** Thibaud Gloaguen, Niels Mündler-Sasahara, Mark Niklas Müller, Veselin Raychev, Martin Vechev. ETH Zurich and LogicStar.ai.
- **Source.** https://arxiv.org/abs/2602.11988 HTML asset `v3`. BibTeX year 2026. The “Submitted on” line was not in the converted HTML.
- **Question.** Do repository context files raise task success?
- **Method, as stated on the page.** SWE-bench tasks with LLM-generated context files, and CTXbench: 138 instances from 12 repositories that already had developer context files. The specialist also reported SWE-bench Lite at 300 tasks and four agent–model pairs sampled once. Those sample details were in the specialist’s reading of later sections. The introduction and results paragraphs opened here support the rates below.
- **Findings opened here.** Context files do not generally improve success and raise inference cost by over 20 percent on average (abstract). Developer-committed files outperform LLM-generated ones by a margin the introduction calls 7 percent. The results text says LLM-generated files change resolution by about −0.5 percent on SWE-bench and −2 percent on CTXbench, with the page giving those tests as p-values of 87 percent and 37 percent. Cost rises about 20 percent and 23 percent. Developer-provided files are about +2.4 percent versus no file (p = 21 percent) and beat generated files (p = 3.8 percent). Instructions are followed. Repository overviews do not help the agent reach the gold-patch files sooner, per the abstract.
- **Limitations the specialist reported from the paper, not re-copied from a limitations section in this writing pass.** Python only. One sample per instance. CTXbench issue text was model-assisted. Status of that limitation paragraph: REPORTED.
- **Implication.** Keep human lines that state a convention the README does not. Do not auto-write a standing overview as a performance feature. `docs/CONDUCTOR.md` at 19,467 bytes is the kind of file this paper says to test before trusting.

### SkillsBench

- **Title.** SkillsBench: Benchmarking How Well Agent Skills Work Across Diverse Tasks.
- **Authors.** Xiangyi Li and a long co-author list ending with Dawn Song, as printed on the HTML. BenchFlow and academic affiliations.
- **Source.** https://arxiv.org/abs/2602.12670 HTML asset `v4`. BibTeX year 2026.
- **Question.** How much do curated skills change pass rate under matched conditions?
- **Method opened here.** 87 tasks across 8 domains, kept from 400 submissions. 18 model–harness configurations. Matched no-skill versus curated-skill. Tasks with no measurable separation were rejected as low-signal. A third condition, self-generated skills, is reported separately on three dedicated harnesses.
- **Findings opened here.** Curated skills move the reported average from 33.9 percent to 50.5 percent (+16.6 points). Configuration-level gains range from +4.1 to +25.7 points. Focused skills of at most three modules beat larger bundles (abstract). Self-generated skills land below the no-skill baseline on the three dedicated configs: −8.1, −11.3, and −11.5 points, while curated skills on those same configs add +18.2 to +24.8 points. One audited “win” for a self-generated skill named the verifier’s checks. The page says the Claude Code self-generated runs used a different effort level than the baseline.
- **Limitations opened here.** The gain is on tasks kept because the conditions separated, so it is not a default lift. Extra context length is a confound the authors name. The set is not a sample of ordinary low-quality skills. Terminal containers only. Self-generated runs mix bad content, discovery failure, and creator/solver interference.
- **Implication.** A short procedure can help a hard task. A self-written skill pack can hurt. Orchestra’s `PromoteSkill` gate matches the second sentence. The +16.6 point average must not be copied onto this repo.

### SWE-Skills-Bench, current HTML

- **Title.** SWE-Skills-Bench: Do Agent Skills Actually Help in Real-World Software Engineering?
- **Authors on the HTML.** Tingxu Han, Yi Zhang, Wei Song, Ye Shang, Chunrong Fang, Zhenyu Chen, Youcheng Sun, Lijie Hu.
- **Source.** https://arxiv.org/abs/2603.15401 HTML asset `v2`. Code link printed on the page: https://github.com/ControllableModel/SWE-Skills-Bench. That GitHub repo was not opened.
- **Question.** Does one public skill change requirement pass ratio on a pinned repo?
- **Method opened here.** 98 skills, ten requirements each, four agent–model configurations, 3,920 paired runs. Requirements are reverse-engineered from the skill. Verifiers are deterministic. Each pair is one clean run. “Ineffective” means a sign test on ten deltas did not find a directional effect.
- **Findings opened here.** The results table pools all pairs at 76.59 percent without the skill and 76.93 percent with it (+0.34). The abstract rounds that to 76.6 percent and 76.9 percent. 92 of 98 skills are called ineffective. 2,711 of 3,920 pairs were unchanged. 40 of 98 skills change sign between two harnesses on the same model. `viper-configuration` is cited at about +17.7 points because unassisted runs skip the library. `asyncio-patterns` is cited at about −17.0 points because templates override the task, including file paths.
- **Limitations opened here.** Requirements are written from the skill, so they test the skill’s claimed behavior rather than a natural request stream. Verifiers can favor one implementation. One run per pair.
- **v1 PDF.** The specialist reported a v1 PDF at the same id, dated on that PDF 16 Mar 2026, with 49 skills and about +1.2 percent, and a different code URL (`GeniusHTX/SWE-Skills-Bench`). That PDF was not opened. Those v1 figures stay REPORTED. Cite the current HTML, not the v1 numbers, until the PDF is opened.
- **Implication.** Skill utility is configuration-dependent. A template that names paths can beat the user’s paths. Orchestra skill text should say the task’s paths and acceptance line win.

The two skill papers disagree with each other. SkillsBench’s average is large on a filtered, curated set. SWE-Skills-Bench’s average is near zero on requirements written from the skills. Both disagree with transferring either average onto Orchestra.

## REPORTED by the papers specialist, not re-opened here

Use these as a reading list. The numbers are the specialist’s, from pages they said they opened. If a number matters for a design change, open the page again before coding.

| Paper | URL | One reported finding | Why it is not a requirement |
| --- | --- | --- | --- |
| SWE-bench, Jimenez et al. | https://arxiv.org/abs/2310.06770 | Claude 2 resolved 1.96 percent of 2,294 tasks | Old model and scaffold. Python only |
| SWE-bench Verified | https://openai.com/index/introducing-swe-bench-verified/ | 500 human-filtered instances | A later OpenAI post (https://openai.com/index/why-we-no-longer-evaluate-swe-bench-verified/, reported not re-opened) says the set is contaminated |
| Terminal-Bench 2.0 | https://arxiv.org/abs/2601.11868 | 89 tasks. The specialist reported frontier scores under 65 percent | Their models, their harness |
| SWE-agent | https://arxiv.org/abs/2405.15793 | Agent-computer interface beat a raw shell on a subset | Interface was shaped on one model |
| Agentless | https://arxiv.org/abs/2407.01489 | A fixed localize-repair-validate pipeline, model not choosing the next action, reported at 32 percent on SWE-bench Lite | One paper’s comparison, not a mandate to delete tool loops |
| OpenHands | https://arxiv.org/abs/2407.16741 | Event stream plus a sandbox | Harness paper, not a skill pack |
| Reflexion | https://arxiv.org/abs/2303.11366 | Verbal feedback raised a HumanEval number the specialist quoted as 91 percent versus 80 percent | Self-evaluation can stick in a local minimum. Not a repo-repair result |
| ReAct | https://arxiv.org/abs/2210.03629 | Interleaving a thought and a tool call beat reason-only and act-only on their tasks | A thought can itself be wrong |
| MemGPT | https://arxiv.org/abs/2310.08560 | Paging between the window and external storage | Document QA, not a coding harness. Letta is the later product name. The product was not opened |
| Voyager | https://arxiv.org/abs/2305.16291 | A library of executable Minecraft skills | Do not copy the item-count multiplier into a repo agent |
| AutoGen | https://arxiv.org/abs/2308.08155 | Multi-agent conversation framework | Not evidence that a committee fixes repository bugs |
| MetaGPT | https://arxiv.org/abs/2308.00352 | Role SOPs. Their SoftwareDev score is on a self-written set | Not a repository-issue result |
| A Survey of Context Engineering | https://arxiv.org/abs/2507.13334 | Map of retrieval, processing, and management. They say they analyzed over 1,400 papers | A survey, not an experiment on this repo |
| Context Engineering for AI Agents in Open-Source Software | https://arxiv.org/abs/2510.21413 | Mining study of repositories that adopted an agent config file. They did not measure patch quality | Adoption is not effectiveness |
| OSWorld | https://arxiv.org/abs/2404.07972 | GUI computer tasks. The specialist reported humans over 72 percent and the best model in that eval at 12.24 percent | A text coding harness does not become a UI agent by adding a screenshot |

Aider’s SWE-bench figure was reported only as a citation inside the OpenHands table. It was not opened. It is not used.

## What was not read

Full PDFs. Every results cell of SkillsBench. The SWE-Skills v1 PDF. Leaderboards. Any re-run of these benchmarks. Product docs for Letta, Mem0, DSPy, CrewAI, LangGraph, and Mastra.
