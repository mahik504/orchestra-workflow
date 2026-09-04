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
  },
  '3d-portfolio': {
    quality_bar: 'EXPERIMENTAL',
    risk_rank: 8,
    platform: 'web',
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
      'The brief names papers, arXiv, footnotes, citations, or math typesetting',
      'Reflow, measure, and vertical rhythm are explicit concerns',
    ],
    skip_conditions: [
      'The content is short marketing copy',
      'Interaction or dashboards dominate the surface',
      'The brief wants heavy motion; this route protects reading',
    ],
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
  },
  'reverse-engineering': {
    quality_bar: 'STANDARD',
    risk_rank: 2,
    platform: 'web',
    trigger_conditions: [
      'The human supplied a specific URL to learn from',
      'The goal is to extract tokens, type scale, or layout rules into a design note',
      'A reference needs to become a documented system before any build',
    ],
    skip_conditions: [
      'No concrete URL was named — this is not a browsing route',
      'The intent is to copy branding, assets, copy, or source rather than principles',
      'A direction already exists and only needs implementing',
    ],
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

graph.version = '3.1.0';

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
