// Add routing governance to every capability row.
//
// Every row must carry BOTH a trigger and a skip. A capability that can only be
// entered and never declined is not a route, it is a default — and defaults are
// how a restaurant site, a school SaaS, and a 3D portfolio end up looking alike.
//
//   quality_bar       STANDARD | PREMIUM | EXPERIMENTAL  (Design Lab default)
//   risk_rank         1 = cheapest to be wrong about, 10 = most expensive.
//                     Used by the silence fallback: no answer -> lower rank wins.
//   trigger_conditions  plain-language reasons to enter this route
//   skip_conditions     plain-language reasons to decline it
//
// Usage: node runtime/tools/enrich-graph.js [--write]
const fs = require('fs');
const path = require('path');

const GRAPH = path.join(__dirname, '..', '..', 'registries', 'design-resource-graph.json');
const SCHEMA = path.join(__dirname, '..', '..', 'registries', 'schemas', 'design-resource-graph.schema.json');

const GOV = {
  'premium-website': {
    quality_bar: 'PREMIUM',
    risk_rank: 6,
    platform: 'web',
    trigger_conditions: [
      'A stranger, client, or recruiter is the audience',
      'The brief names a brand, studio, agency, restaurant, or launch page',
      'Marketing or conversion copy is the primary content',
      'The request asks for a distinctive look rather than a working form',
      'A showcase carries commerce alongside the story — a shop or checkout bolted onto a brand page',
    ],
    add_trigger_tags: ['portfolio', 'shop', 'storefront', 'ecommerce', 'checkout', 'marketing-site'],
    skip_conditions: [
      'The surface is behind a login and only staff will see it',
      'The request is a bugfix, refactor, or backend change',
      'The product is a dense data tool — route to saas-dashboard or b2b-portal',
      'The human said "skip the lab" and only wants the page wired up',
      'A 3D scene is the point rather than the decoration — route to 3d-portfolio',
    ],
    skills: ['impeccable', 'taste-design', 'ow-golden-still', 'gsap-core', 'gsap-scrolltrigger', 'gsap-timeline'],
  },
  '3d-portfolio': {
    quality_bar: 'EXPERIMENTAL',
    risk_rank: 8,
    platform: 'web',
    skills: ['ow-3d-engineering'],
    trigger_conditions: [
      'The brief names WebGL, Three.js, R3F, shaders, or a 3D scene',
      'Camera movement or spatial navigation is part of the concept',
      'The showcase itself is the product, not a page describing a product',
    ],
    skip_conditions: [
      'No low-end or reduced-motion fallback is acceptable to the human',
      '3D is decorative and a flat hero would carry the same message',
      'The target is a low-power device or an audience on poor connections',
      'The brief is a portfolio but never mentions depth, scene, or motion — route to premium-website',
      'Physics simulation is the point rather than presentation — route to physics-canvas',
      'Checkout or an orders console is the job rather than a spatial scene',
    ],
  },
  'operator-hud': {
    quality_bar: 'PREMIUM',
    risk_rank: 7,
    platform: 'web',
    trigger_conditions: [
      'A single operator monitors live state and acts on it',
      'The brief names telemetry, cockpit, mission control, or a command surface',
      'Information density matters more than whitespace',
      'Monospace figures and always-on status are expected',
      'The brief describes watching servers, jobs, or streams while they run',
    ],
    add_trigger_tags: ['monitoring', 'uptime', 'observability', 'status-board', 'incident'],
    skip_conditions: [
      'The audience is many casual users rather than one trained operator',
      'The data is reported after the fact — route to saas-dashboard',
      'The brief expects a component kit look; this route is bespoke by definition',
      'Administrative CRUD is the real job — route to b2b-portal',
    ],
  },
  'b2b-portal': {
    quality_bar: 'STANDARD',
    risk_rank: 4,
    platform: 'web',
    trigger_conditions: [
      'Authenticated users perform recurring transactional work',
      'The brief names admin, tenants, roles, billing, or approvals',
      'The brief names an orders console or admin area next to a shop',
      'Tables, forms, and workflow states are the core surface',
      'Accessibility and predictability outrank visual novelty',
    ],
    skip_conditions: [
      'The page is public marketing — route to premium-website',
      'The value is charts and trends rather than records — route to saas-dashboard',
      'The human explicitly asked for an award-winning visual treatment',
      'There is no authentication and no multi-user concept',
    ],
  },
  'academic-reader': {
    quality_bar: 'PREMIUM',
    risk_rank: 3,
    platform: 'web',
    trigger_conditions: [
      'Long-form prose is the product and reading comfort is the metric',
      'The brief names papers, arXiv, footnotes, or math typesetting',
      'Reflow, measure, and vertical rhythm are explicit concerns',
    ],
    skip_conditions: [
      'The content is short marketing copy',
      'Interaction or dashboards dominate the surface',
      'The brief wants heavy motion; this route protects reading',
      'The job is to write a paper or produce citations rather than a reading surface — route to ow-research-paper',
    ],
  },
  'ow-research-paper': {
    quality_bar: 'STANDARD',
    risk_rank: 2,
    platform: 'any',
    trigger_conditions: [
      'The job is producing a paper, thesis, or methods writeup',
      'Citations must be real; inventing sources is forbidden',
      'There is no reading UI to design — the deliverable is prose',
    ],
    add_trigger_tags: ['citations', 'thesis', 'manuscript', 'methods-paper', 'academic-writing', 'write-a-paper'],
    skip_conditions: [
      'The job is a reading viewer or arXiv app — route to academic-reader',
      'The brief is a marketing page about research',
      'The request is a UI or dashboard, not a document',
    ],
    skills: ['ow-research-paper', 'humanizer'],
  },
  'micro-interactions': {
    quality_bar: 'PREMIUM',
    risk_rank: 5,
    platform: 'web',
    trigger_conditions: [
      'The unit of work is a component, not a page',
      'The brief names feel, tactility, spring, press feedback, or haptics',
      'An existing UI is being refined rather than created',
    ],
    skip_conditions: [
      'No product surface exists yet to attach the interaction to',
      'The whole page or product needs a direction first',
      'Motion has no stated purpose — do not animate to look busy',
      'The platform is native mobile — route to mobile-app',
    ],
  },
  'physics-canvas': {
    quality_bar: 'EXPERIMENTAL',
    risk_rank: 7,
    platform: 'web',
    trigger_conditions: [
      'Collision, gravity, particles, or rigid bodies are simulated',
      'The brief names a game, playground, or interactive simulation',
      'Per-frame computation drives the visuals',
    ],
    skip_conditions: [
      'The motion is scripted presentation rather than simulated behaviour',
      'A CSS or timeline animation achieves the same result',
      'The scene is a spatial showcase without simulation — route to 3d-portfolio',
    ],
  },
  'saas-dashboard': {
    quality_bar: 'STANDARD',
    risk_rank: 4,
    platform: 'web',
    trigger_conditions: [
      'Metrics, KPIs, charts, or trends are the primary content',
      'The brief names reporting, analytics, or an overview screen',
      'Users filter and compare rather than transact',
      'A school, clinic, or internal management tool needs its numbers surfaced',
    ],
    add_trigger_tags: ['dashboard', 'charts', 'reporting', 'metrics'],
    skip_conditions: [
      'The screens are records and forms rather than charts — route to b2b-portal',
      'The surface is live operational control — route to operator-hud',
      'The page is public marketing — route to premium-website',
    ],
    skills: ['ow-backend-architecture', 'test-driven-development'],
  },
  'mobile-app': {
    quality_bar: 'PREMIUM',
    risk_rank: 5,
    platform: 'mobile',
    trigger_conditions: [
      'The target is iOS, Android, Expo, or React Native',
      'Touch, gesture, or device haptics are part of the design',
      'Platform conventions (HIG, Material 3) constrain the layout',
    ],
    skip_conditions: [
      'The target is a responsive website that merely works on phones',
      'No device capability is used and no store build is planned',
      'The brief is a desktop-first dense tool',
    ],
  },
  'security-audit': {
    quality_bar: 'STANDARD',
    risk_rank: 1,
    platform: 'any',
    trigger_conditions: [
      'The job is to find or fix vulnerabilities in your own application',
      'The brief names SAST, DAST, pentest, dependency audit, or attack surface',
      'Secrets, authorization, or input handling are under review',
    ],
    add_trigger_tags: ['security', 'injection', 'xss', 'csrf', 'owasp', 'dependency-audit', 'supply-chain'],
    skip_conditions: [
      'The target is not owned by the human — never scan someone else\'s system',
      'The request is a feature build; security review is a later gate, not this route',
      'The finding is a plain bug with no security consequence',
    ],
    skills: ['ow-ship-safe'],
  },
  'reverse-engineering': {
    quality_bar: 'STANDARD',
    risk_rank: 2,
    platform: 'web',
    trigger_conditions: [
      'The human supplied a specific URL to learn from',
      'The goal is to extract tokens, type scale, or layout rules into a design note',
      'A reference needs to become a documented system before any build',
      'A named reference, screenshot, or URL should keep structure while tokens change (tint)',
    ],
    add_trigger_tags: ['tint', 'recolor', 'keep-the-structure', 'screenshot-to-code', 'design-adaptation'],
    skip_conditions: [
      'No concrete URL was named — this is not a browsing route',
      'The intent is to copy branding, assets, copy, or source rather than principles',
      'A direction already exists and only needs implementing',
    ],
  },
  'interaction-components': {
    quality_bar: 'PREMIUM',
    risk_rank: 5,
    platform: 'web',
    trigger_conditions: [
      'The unit of work is a control (button, card, drawer, input), not a whole product',
      'The brief names a kit (shadcn, React Bits, Watermelon, Magic UI, Aceternity, Kokonut, Cult, HyperUI, Tailkit free, 21st free, Motion Primitives)',
      'The brief names a control kind: primary, secondary, destructive, loading, success, disabled, icon, magnetic, shimmer, gradient, ripple, particle, liquid, glass, command, hold, stateful',
      'Existing project components must be inspected before generating another',
      'Pick one source in this order: existing project components, then shadcn for app chrome, then React Bits for one interaction, then Watermelon for one production-style block when the reference row exists, then Magic UI or Aceternity for one marketing effect, then Motion Primitives for one motion on the current runtime, then GSAP only when the route is scroll',
    ],
    add_trigger_tags: ['watermelon'],
    skip_conditions: [
      'The job is a full page, product, or app — route to premium-website, saas-dashboard, mobile-app, or 3d-portfolio first',
      'No product surface exists yet to attach the control to',
      '3D/WebGL is the point — route to 3d-portfolio',
      'The platform is native mobile — route to mobile-app',
      'A sample that needs a second motion runtime',
    ],
    skills: [],
  },
  'ow-reference-forensics': {
    quality_bar: 'PREMIUM',
    risk_rank: 4,
    platform: 'web',
    trigger_conditions: [
      'The brief names ow-reference-forensics or a reference-benchmark file',
      'The job is to measure composition, interaction, and typography of a named live website',
      'The output is REFERENCE_BENCHMARK.md with mechanism vs identity and forbidden identity',
      'A visual recreation of a live product needs a translation matrix before Design Lab',
    ],
    skip_conditions: [
      'The goal is only tint, recolor, or keep-the-structure token swap — route to reverse-engineering',
      'No named site, screenshot, or public source was supplied',
      'The intent is to copy branding, assets, copy, or source',
      'A locked DESIGN.md and golden stills already exist and the job is implementation',
    ],
    skills: ['ow-reference-forensics'],
  },
  'ow-golden-still': {
    quality_bar: 'PREMIUM',
    risk_rank: 5,
    platform: 'web',
    trigger_conditions: [
      'The brief names golden-still or ow-golden-still',
      'The job is a representative desktop opening scene and mobile opening scene',
      'Visual evidence PASS/FAIL is owed before product implementation',
      'A DESIGN.md exists and stills are the next gate',
    ],
    skip_conditions: [
      'The lab was waived with skip the lab',
      'The work is STANDARD backend, paper, or security with no UI',
      'Stills already human-approved and the job is implementation',
      'The request is a full product build with no stills gate named — classify the product route first',
    ],
  },
  'ow-fresh-review': {
    quality_bar: 'STANDARD',
    risk_rank: 3,
    platform: 'web',
    trigger_conditions: [
      'The brief names ow-fresh-review or independent review of a completed phase',
      'The human said the work is ready for review or pre-release',
      'The question is what is wrong that the builder cannot see',
      'A major frontend milestone needs a clean-context critic',
    ],
    skip_conditions: [
      'The phase is not complete — this is not a planning or implementation route',
      'The request is to build or redesign from scratch',
      'The human asked the same chat that built the work to rubber-stamp it without a fresh context',
    ],
  },
  'ow-showcase': {
    quality_bar: 'STANDARD',
    risk_rank: 1,
    platform: 'web',
    trigger_conditions: [
      'The brief names ow-showcase or brag as a delivery step',
      'The product is production-ready and the verification ledger is green',
      'The human approved recording a walkthrough and showcase stills',
      'README or social showcase assets are owed after ship, not before',
    ],
    skip_conditions: [
      'Design Lab or golden stills are still unpaid',
      'Implementation is in progress',
      'The verification ledger is not green',
      'The human has not approved production readiness',
      'The request is to change product behaviour or visuals — that is not showcase',
    ],
    skills: ['ow-ship'],
  },
  'typescript-review': {
    quality_bar: 'STANDARD', risk_rank: 2, platform: 'any',
    skills: ['ow-typescript'],
    trigger_conditions: [
      'The job is a TypeScript review of types at a boundary',
      'The brief names Zod or a typed parser at the edge of the program',
    ],
    skip_conditions: ['The job is a visual redesign', 'The job is to adopt a new framework'],
  },
  'slides': {
    quality_bar: 'STANDARD', risk_rank: 2, platform: 'any',
    skills: ['ow-docs', 'ow-presentation'],
    trigger_conditions: ['The deliverable is a slide deck', 'The brief asks for speaker notes on slides'],
    skip_conditions: ['The job is a methods paper or a thesis', 'The deliverable is a rendered mp4'],
  },
  'recent-research': {
    quality_bar: 'STANDARD', risk_rank: 2, platform: 'any',
    skills: ['ow-research'],
    trigger_conditions: [
      'The question is what people said recently in public',
      'Public discussions from the last month need a dated source',
    ],
    skip_conditions: ['The job is a methods paper or a thesis', 'The request asks to scrape, post, or send messages'],
  },
  'marketing-copy': {
    quality_bar: 'STANDARD', risk_rank: 2, platform: 'any',
    skills: ['ow-marketing', 'humanizer'],
    trigger_conditions: ['The job is the page title or the on-page wording', 'The brief asks for a positioning line'],
    skip_conditions: ['The job is a visual redesign or a new design system', 'The request invents a metric'],
  },
  'rendered-video': {
    quality_bar: 'STANDARD', risk_rank: 5, platform: 'any',
    skills: ['ow-video'],
    trigger_conditions: ['The deliverable is a rendered mp4', 'The brief names HyperFrames for an mp4 file'],
    skip_conditions: ['The job is scroll animation', 'The deliverable is a slide deck'],
  },
  'deployed-errors': {
    quality_bar: 'STANDARD', risk_rank: 2, platform: 'any',
    skills: ['ow-ship-safe', 'ow-observability'],
    trigger_conditions: [
      'A deployed application is throwing a production error',
      'The brief names a production error in a shipped app',
    ],
    skip_conditions: ['The request is to scan a third-party URL', 'The job is a local bug with no deployment'],
  },
  'career-notes': {
    quality_bar: 'STANDARD', risk_rank: 2, platform: 'any',
    skills: ['ow-career'],
    trigger_conditions: [
      'The job is to tailor a CV from notes already on disk',
      'The brief asks for a role brief from notes the human already has',
    ],
    skip_conditions: ['The request is to scrape a job board', 'The request is to apply on the human behalf'],
  },
  'architecture-diagram': {
    quality_bar: 'STANDARD', risk_rank: 2, platform: 'any',
    skills: ['ow-diagrams', 'pretty-mermaid'],
    trigger_conditions: ['The deliverable is one architecture diagram'],
    skip_conditions: ['The request wants a second diagram system', 'The request fills a README with charts'],
  },
  'nextjs-speed': {
    quality_bar: 'STANDARD', risk_rank: 2, platform: 'any',
    skills: ['react-best-practices'],
    trigger_conditions: ['The job is Next.js speed or a request waterfall', 'The brief asks about a published package size'],
    skip_conditions: ['The job is a visual redesign', 'The request is a bulk install of Vercel skills'],
  },
  'spec-review': {
    quality_bar: 'STANDARD', risk_rank: 2, platform: 'any',
    skills: ['egm-draft', 'egm-review', 'egm-resolve'],
    trigger_conditions: ['The job is a spec or a design doc before the build'],
    skip_conditions: ['A STANDARD fix already names the file'],
  },
  'codebase-map': {
    quality_bar: 'STANDARD', risk_rank: 2, platform: 'any',
    skills: ['floorplan', 'ow-diagrams'],
    trigger_conditions: ['The job is a map of a codebase the operator has'],
    skip_conditions: ['The repository is not on disk'],
  },
  'isometric-explainer': {
    quality_bar: 'STANDARD', risk_rank: 2, platform: 'any',
    skills: ['anatomy'],
    trigger_conditions: ['The job is an isometric explainer'],
    skip_conditions: ['The job is a software architecture chart'],
  },
  'mcp-server': {
    quality_bar: 'STANDARD', risk_rank: 2, platform: 'any',
    skills: ['mcp-builder'],
    trigger_conditions: ['The job is an MCP server'],
    skip_conditions: ['The request installs every MCP'],
  },
  'running-app-check': {
    quality_bar: 'STANDARD', risk_rank: 2, platform: 'web',
    skills: ['webapp-testing'],
    trigger_conditions: ['The job is a full-stack app with a running page'],
    skip_conditions: ['There is no running page'],
  },
  'ml-experiment': {
    quality_bar: 'STANDARD', risk_rank: 2, platform: 'any',
    skills: ['ow-ml-experiment'],
    trigger_conditions: ['The job is one ML experiment'],
    skip_conditions: ['The request claims a benchmark win'],
  },
  'award-site': {
    quality_bar: 'PREMIUM', risk_rank: 6, platform: 'web',
    skills: ['build-awwwards-quality-sites', 'unslop'],
    trigger_conditions: ['The job is an award-style site or a polish pass on a screen that already exists'],
    skip_conditions: ['The request copies a reference site'],
  },
};

const graph = JSON.parse(fs.readFileSync(GRAPH, 'utf8'));
const problems = [];
let touched = 0;

for (const [id, cap] of Object.entries(graph.capabilities)) {
  const gov = GOV[id];
  if (!gov) { problems.push(`no governance authored for capability '${id}'`); continue; }
  const { add_trigger_tags: extraTags, ...fields } = gov;
  Object.assign(cap, fields);
  if (extraTags) {
    const tags = new Set(cap.trigger_tags || []);
    extraTags.forEach((t) => tags.add(t));
    cap.trigger_tags = [...tags];
  }
  touched += 1;
}
for (const id of Object.keys(GOV)) {
  if (!graph.capabilities[id]) problems.push(`authored governance for unknown capability '${id}'`);
}

// Every row must have both a trigger and a skip.
for (const [id, cap] of Object.entries(graph.capabilities)) {
  if (!Array.isArray(cap.trigger_conditions) || cap.trigger_conditions.length === 0) problems.push(`${id}: no trigger_conditions`);
  if (!Array.isArray(cap.skip_conditions) || cap.skip_conditions.length === 0) problems.push(`${id}: no skip_conditions`);
  if (!['STANDARD', 'PREMIUM', 'EXPERIMENTAL'].includes(cap.quality_bar)) problems.push(`${id}: bad quality_bar`);
  if (typeof cap.risk_rank !== 'number') problems.push(`${id}: no risk_rank`);
}

graph.version = '4.0.0';

// Schema must allow the new fields (additionalProperties is false).
const schema = JSON.parse(fs.readFileSync(SCHEMA, 'utf8'));
const capSchema = schema.properties.capabilities.additionalProperties;
capSchema.properties.quality_bar = { type: 'string', enum: ['STANDARD', 'PREMIUM', 'EXPERIMENTAL'] };
capSchema.properties.risk_rank = { type: 'integer', minimum: 1, maximum: 10 };
capSchema.properties.platform = { type: 'string' };
capSchema.properties.trigger_conditions = { type: 'array', minItems: 1, items: { type: 'string' } };
capSchema.properties.skip_conditions = { type: 'array', minItems: 1, items: { type: 'string' } };
const req = new Set(capSchema.required || []);
['quality_bar', 'risk_rank', 'trigger_conditions', 'skip_conditions'].forEach((r) => req.add(r));
capSchema.required = [...req];

console.log(`capabilities:      ${Object.keys(graph.capabilities).length}`);
console.log(`governance added:  ${touched}`);
console.log(`problems:          ${problems.length}`);
for (const p of problems) console.log(`  ${p}`);

if (problems.length) process.exit(1);

if (process.argv.includes('--write')) {
  fs.writeFileSync(GRAPH, JSON.stringify(graph, null, 2) + '\n', 'utf8');
  fs.writeFileSync(SCHEMA, JSON.stringify(schema, null, 2) + '\n', 'utf8');
  console.log('\nwrote graph + schema');
} else {
  console.log('\ndry run. pass --write to apply.');
}
