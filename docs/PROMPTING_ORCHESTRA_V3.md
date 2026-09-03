# Prompting Orchestra V3

A strong prompt is the difference between generic AI output and premium engineering. Orchestra V3 relies on your prompt to classify the task and load the correct capabilities.

## The Prompt Formula
The ideal prompt contains 6 elements:
`Intent + Context + Requirements + Constraints + References + Acceptance Criteria`

## Examples

### ❌ Bad Prompt
> "Make a landing page for my ag-tech startup."

*Why it fails:* Orchestra doesn't know the tone, the tech stack, the constraints, or what defines success. It will fallback to generic templates.

### ✅ Good Prompt
> **[Intent]** Build a premium landing page for an ag-tech supply chain company.
> **[Context]** We deal with enterprise clients handling $10M+ in international trade.
> **[Requirements]** It needs a hero section, a global supply chain visualization, and a contact form.
> **[Constraints]** Use Vite + React + Tailwind. No generic 3-column cards.
> **[References]** Use a restrained editorial style, similar to high-end trading houses.
> **[Acceptance Criteria]** Lighthouse score > 90, fully responsive, zero CLS.

## Explicit Directives
You have override authority over Orchestra's routing. You can explicitly shape the execution by including these phrases in your prompt:

- **"Make this premium"** -> Activates `taste-design`, strict typography rules, and asymmetric layout engines.
- **"Use a restrained editorial style"** -> Shifts design focus away from SaaS dashboards toward corporate/editorial layouts.
- **"Use 3D only if justified"** -> Prompts Orchestra to evaluate the performance cost before loading R3F/Three.js.
- **"Research current references"** -> Forces Orchestra to search the web or read specific files before planning.
- **"Use my provided PDF as authoritative"** -> Locks the factual context to your uploaded document.
- **"Optimize for quality, not token minimization"** -> Instructs the router to load the full heavyweight design ecosystem.
- **"Show me alternatives before implementation"** -> Forces an explicit human approval gate.
