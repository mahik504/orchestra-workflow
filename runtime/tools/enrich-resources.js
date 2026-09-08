// Backfill trigger/skip on every catalog row and upsert 3.3.0 free-resource rows.
//
//   node runtime/tools/enrich-resources.js          # dry run
//   node runtime/tools/enrich-resources.js --write
const fs = require('fs');
const path = require('path');

const ROOT = path.join(__dirname, '..', '..');
const REG = path.join(ROOT, 'registries');
const WRITE = process.argv.includes('--write');

function load(file) {
  return JSON.parse(fs.readFileSync(path.join(REG, file), 'utf8'));
}

function save(file, data) {
  fs.writeFileSync(path.join(REG, file), JSON.stringify(data, null, 2) + '\n', 'utf8');
}

const TODAY = '2026-09-06';

function row(partial) {
  return {
    last_verified: TODAY,
    ...partial,
  };
}

const UPSERTS = {
  'magic-ui-oss': {
    name: 'Magic UI OSS',
    canonical_url: 'https://github.com/magicuidesign/magicui',
    source_type: 'github_repository',
    source_repository: 'https://github.com/magicuidesign/magicui.git',
    documentation_url: 'https://magicui.design/docs/mcp',
    license: 'MIT',
    category: ['components', 'oss', 'motion'],
    representation: 'reference',
    routing_tags: ['magic-ui', 'interaction-components', 'free-registry'],
    acquisition_method: 'npm',
    runtime_method: 'project_scoped_install',
    status: 'CURATED_OPTIONAL',
    trigger_conditions: [
      'named Magic UI OSS component',
      'interaction-components route named Magic UI',
    ],
    avoid_conditions: [
      'buying Magic UI Pro this pass',
      'always-on Magic MCP',
      'claiming Pro registry access',
    ],
    skip_conditions: [
      'another motion kit already chosen',
      'the job is a full page — route the page first',
    ],
    policy_verdict: 'FREE OSS + FREE MCP; PRO NOT CLAIMED',
    rationale: 'Official free MCP is npx @magicuidesign/mcp@latest (or npx @magicuidesign/cli@latest install cursor). Do not buy Pro. Do not claim Pro blocks.',
    auth: 'NONE',
    free_capability: 'OSS components + free MCP catalog',
    paid_capability: 'Magic UI Pro blocks/templates',
    fallback: 'shadcn-ui',
  },
  'twenty-first-dev': {
    name: '21st.dev',
    canonical_url: 'https://21st.dev',
    source_type: 'component_gallery',
    documentation_url: 'https://docs.21st.dev/mcp',
    license: 'site-terms',
    category: ['components', 'mcp', 'optional-free'],
    representation: 'mcp',
    routing_tags: ['21st', 'component-search', 'optional-free-mcp'],
    acquisition_method: 'host_mcp_connect',
    runtime_method: 'mcp_connection',
    status: 'CURATED_OPTIONAL',
    trigger_conditions: [
      'human named 21st.dev',
      'component discovery / React UI catalog search',
      'one named echo from the free catalog',
    ],
    avoid_conditions: [
      'always-on 21st MCP',
      'buying Builder or Builder + AI this pass',
      'repeatedly exceeding two free installs per day',
    ],
    skip_conditions: [
      'MCP not Connected — use public catalog URLs instead',
      'operator HUD',
      'free daily install quota would be burned on a throwaway',
    ],
    policy_verdict: 'OPTIONAL_FREE_MCP',
    rationale: 'Search and publishing are free. Installs capped at two a day. 21st AI uses credits. Login in host. Never always-on. Do not buy Builder this pass.',
    auth: 'AUTH_REQUIRED',
    free_capability: 'catalog search; up to 2 component installs per day',
    paid_capability: 'Builder, Builder + AI, extra installs, AI generate credits',
    fallback: 'shadcn-ui',
  },
  tailkit: {
    name: 'Tailkit public/free resources',
    canonical_url: 'https://tailkit.com',
    source_type: 'component_gallery',
    documentation_url: 'https://tailkit.com',
    license: 'site-terms',
    category: ['components', 'reference'],
    representation: 'reference',
    routing_tags: ['tailkit', 'tailwind-kit', 'interaction-components', 'free-public'],
    acquisition_method: 'web_fetch',
    runtime_method: 'on_demand_research',
    status: 'REFERENCE',
    trigger_conditions: [
      'human named Tailkit',
      'interaction-components search for a public Tailkit pattern',
    ],
    avoid_conditions: [
      'buying Tailkit Unlimited this pass',
      'treating the Tailkit MCP as free',
      'mirroring the paid library',
    ],
    skip_conditions: [
      'another free kit already covers the control',
      'the job needs the paid MCP — that row is AUTH_REQUIRED',
    ],
    policy_verdict: 'FREE PUBLIC REFERENCE; MCP IS PAID',
    rationale: 'Public/free Tailkit pages are searchable references. The MCP is not free. Do not invent credentials.',
    auth: 'NONE',
    free_capability: 'public/free component pages as URL references',
    paid_capability: 'Tailkit MCP and paid kit',
    fallback: 'hyperui',
  },
  'aceternity-ui': {
    name: 'Aceternity UI free components',
    canonical_url: 'https://ui.aceternity.com/',
    source_type: 'component_gallery',
    documentation_url: 'https://ui.aceternity.com/',
    license: 'site-terms',
    category: ['components', 'reference'],
    representation: 'reference',
    routing_tags: [
      'aceternity',
      'interaction-components',
      'buttons',
      'backgrounds',
      'text',
      'cards',
      'forms',
      'navigation',
      'effects',
      '3d-components',
      'hero',
    ],
    acquisition_method: 'web_fetch',
    runtime_method: 'on_demand_research',
    status: 'CURATED_OPTIONAL',
    trigger_conditions: [
      'one named Aceternity echo on the card',
      'free category search: buttons, backgrounds, text, cards, forms, navigation, effects, 3D, hero',
      'official shadcn install command on a public component page',
    ],
    avoid_conditions: [
      'always-on Aceternity MCP',
      'buying All-Access this pass',
      'claiming Pro blocks/templates',
    ],
    skip_conditions: [
      'operator HUD',
      'college due-now GUI',
      'the card already named another kit',
    ],
    policy_verdict: 'FREE COMPONENTS ONLY',
    rationale: 'Use public components and official shadcn install commands. Do not claim All-Access.',
    auth: 'NONE',
    free_capability: 'public component catalog (buttons, backgrounds, text, cards, forms, navigation, effects, 3D, hero)',
    paid_capability: 'All-Access blocks/templates',
    fallback: 'magic-ui-oss',
  },
  'refero-design': {
    skip_conditions: [
      'MCP is AUTH_REQUIRED / PRO and the job can use public Refero pages instead',
      'backend-only, papers, operator HUD',
    ],
    policy_verdict: 'SKILL FREE; MCP AUTH_REQUIRED / PRO',
    rationale: 'Global research skill (referodesign/refero_skill). MCP at https://api.refero.design/mcp needs Pro. Public styles, DESIGN.md references, and public resource pages are free. Never copy branding.',
    auth: 'NONE',
    free_capability: 'skill methodology + public Refero pages, styles, design prompts',
    paid_capability: 'Refero MCP / Pro screen search',
    last_verified: TODAY,
  },
  getlayers: {
    last_verified: TODAY,
    auth: 'PURCHASE_PENDING',
    free_capability: 'none until purchase',
    paid_capability: 'Full Stack lifetime MCP and library',
  },
  jcode: {
    representation: 'executor',
    skip_conditions: [
      'Cursor, Antigravity, or Claude Code is already conducting',
      'no stdout evidence that jcode ran',
    ],
    last_verified: TODAY,
  },
};

const NEW_ROWS = [
  row({
    id: 'magic-ui-mcp',
    name: 'Magic UI free MCP',
    canonical_url: 'https://magicui.design/docs/mcp',
    source_type: 'mcp_server',
    source_repository: 'https://github.com/magicuidesign/mcp.git',
    documentation_url: 'https://magicui.design/docs/mcp',
    license: 'MIT',
    category: ['mcp', 'components', 'oss'],
    representation: 'mcp',
    routing_tags: ['magic-ui-mcp', 'free-mcp', 'interaction-components'],
    acquisition_method: 'host_mcp_connect',
    runtime_method: 'mcp_connection',
    status: 'CURATED_OPTIONAL',
    trigger_conditions: [
      'named Magic UI component install via MCP',
      'host already has magicuidesign-mcp connected',
    ],
    avoid_conditions: [
      'always-on Magic MCP',
      'Pro registry',
    ],
    skip_conditions: [
      'MCP not Connected — use magic-ui-oss URLs',
      'another motion kit already chosen',
    ],
    policy_verdict: 'FREE_MCP',
    rationale: 'Official install: npx @magicuidesign/cli@latest install cursor, or mcpServers.magicuidesign-mcp = npx -y @magicuidesign/mcp@latest. Pro is separate and unpaid.',
    auth: 'HOST_CONNECT',
    free_capability: 'free Magic UI MCP component catalog',
    paid_capability: 'Magic UI Pro',
    fallback: 'magic-ui-oss',
  }),
  row({
    id: 'refero-mcp',
    name: 'Refero MCP',
    canonical_url: 'https://api.refero.design/mcp',
    source_type: 'mcp_server',
    documentation_url: 'https://refero.design/',
    license: 'site-terms',
    category: ['mcp', 'visual_research'],
    representation: 'mcp',
    routing_tags: ['refero-mcp', 'pro'],
    acquisition_method: 'host_mcp_connect',
    runtime_method: 'mcp_connection',
    status: 'CURATED_OPTIONAL',
    trigger_conditions: [
      'human has Refero Pro and named Refero MCP',
      'operator said Connect Refero after Pro',
    ],
    avoid_conditions: [
      'pretending the MCP is free',
      'always-on Refero MCP',
    ],
    skip_conditions: [
      'no Pro token — use refero-public instead',
      'backend-only or papers',
    ],
    policy_verdict: 'AUTH_REQUIRED / PRO',
    rationale: 'Refero MCP requires Pro. Bearer token stays in the host. Do not invent credentials.',
    auth: 'AUTH_REQUIRED',
    free_capability: 'none on the MCP itself',
    paid_capability: 'Pro MCP screen search',
    fallback: 'refero-public',
  }),
  row({
    id: 'refero-public',
    name: 'Refero public resources',
    canonical_url: 'https://refero.design/',
    source_type: 'web_reference',
    documentation_url: 'https://refero.design/',
    license: 'site-terms',
    category: ['visual_research', 'reference'],
    representation: 'research_source',
    routing_tags: ['refero-public', 'design-md', 'public-styles'],
    acquisition_method: 'web_fetch',
    runtime_method: 'on_demand_research',
    status: 'REFERENCE',
    trigger_conditions: [
      'SaaS or product UI research without Pro',
      'public Refero styles or DESIGN.md references',
    ],
    avoid_conditions: [
      'cloning a screenshot as the product',
      'bulk scraping refero.design',
    ],
    skip_conditions: [
      'backend-only',
      'operator HUD',
    ],
    policy_verdict: 'RESOURCE / FREE',
    rationale: 'Public pages, styles, and design prompts. Not the Pro MCP.',
    auth: 'NONE',
    free_capability: 'public styles, DESIGN.md references, public resource pages',
    paid_capability: 'Pro MCP (see refero-mcp)',
    fallback: 'godly-design',
  }),
  row({
    id: 'tailkit-mcp',
    name: 'Tailkit MCP',
    canonical_url: 'https://tailkit.com',
    source_type: 'mcp_server',
    documentation_url: 'https://tailkit.com',
    license: 'site-terms',
    category: ['mcp', 'paid'],
    representation: 'mcp',
    routing_tags: ['tailkit-mcp', 'paid-mcp'],
    acquisition_method: 'host_mcp_connect',
    runtime_method: 'mcp_connection',
    status: 'CURATED_OPTIONAL',
    trigger_conditions: [
      'human has a paid Tailkit license and named the MCP',
    ],
    avoid_conditions: [
      'labelling Tailkit MCP as free because public pages exist',
      'inventing credentials',
    ],
    skip_conditions: [
      'no paid Tailkit license — use the tailkit public reference row',
      'purchase not authorized this pass',
    ],
    policy_verdict: 'AUTH_REQUIRED / PAID_ACCESS',
    rationale: 'MCP is paid. Free/public components live on the tailkit resource row.',
    auth: 'AUTH_REQUIRED',
    free_capability: 'none on the MCP',
    paid_capability: 'licensed Tailkit MCP',
    fallback: 'tailkit',
  }),
  row({
    id: 'lapa',
    name: 'Lapa Ninja',
    canonical_url: 'https://www.lapa.ninja/',
    source_type: 'web_reference',
    documentation_url: 'https://www.lapa.ninja/',
    license: 'site-terms',
    category: ['visual_research'],
    representation: 'research_source',
    routing_tags: ['lapa', 'landing-gallery', 'design-research'],
    acquisition_method: 'web_fetch',
    runtime_method: 'on_demand_research',
    status: 'REFERENCE',
    trigger_conditions: [
      'landing-page visual research',
      'named Lapa Ninja',
    ],
    avoid_conditions: [
      'bulk scrape',
      'fake MCP',
    ],
    skip_conditions: [
      'backend-only',
      'the job already has an approved DESIGN.md',
    ],
    policy_verdict: 'RESEARCH / NO MCP',
    rationale: 'Public gallery. Tell the conductor to find a relevant reference. Do not bulk scrape.',
    auth: 'NONE',
    free_capability: 'public landing-page gallery',
    paid_capability: 'none claimed',
    fallback: 'land-book',
  }),
  row({
    id: 'saasframe',
    name: 'SaaSFrame',
    canonical_url: 'https://www.saasframe.io/',
    source_type: 'web_reference',
    documentation_url: 'https://www.saasframe.io/',
    license: 'site-terms',
    category: ['visual_research', 'saas'],
    representation: 'research_source',
    routing_tags: ['saasframe', 'saas-gallery', 'design-research'],
    acquisition_method: 'web_fetch',
    runtime_method: 'on_demand_research',
    status: 'REFERENCE',
    trigger_conditions: [
      'SaaS product UI research',
      'named SaaSFrame',
    ],
    avoid_conditions: [
      'bulk scrape',
      'fake MCP',
    ],
    skip_conditions: [
      'the surface is not a SaaS product',
      'backend-only',
    ],
    policy_verdict: 'RESEARCH / NO MCP',
    rationale: 'Public SaaS screens for research. Do not claim MCP access.',
    auth: 'NONE',
    free_capability: 'public SaaS UI references',
    paid_capability: 'paid SaaSFrame plans if any — do not buy this pass',
    fallback: 'refero-public',
  }),
  row({
    id: 'three-js',
    name: 'Three.js',
    canonical_url: 'https://github.com/mrdoob/three.js',
    source_type: 'github_repository',
    source_repository: 'https://github.com/mrdoob/three.js.git',
    documentation_url: 'https://threejs.org/docs/',
    license: 'MIT',
    category: ['3d', 'webgl'],
    representation: 'dependency',
    routing_tags: ['threejs', 'webgl', '3d', 'canvas', 'shader'],
    acquisition_method: 'npm',
    runtime_method: 'project_scoped_install',
    status: 'ACTIVE',
    trigger_conditions: [
      '3D website, WebGL, interactive 3D, spatial experience, shader, canvas, 3D interaction',
      'portfolio with a 3D scene',
    ],
    avoid_conditions: [
      'global npm install',
      'decorative 3D that a flat hero would replace',
    ],
    skip_conditions: [
      'pure 2D layout',
      'the human refused a low-end fallback on EXPERIMENTAL',
    ],
    policy_verdict: 'PROJECT-SCOPED 3D CORE',
    rationale: 'OSS. Install in the app repo, never globally. Pair with R3F/drei when the stack is React.',
    npm_package: 'three',
    auth: 'NONE',
    free_capability: 'full OSS renderer',
    paid_capability: 'none',
    fallback: 'r3f',
  }),
  row({
    id: 'react-postprocessing',
    name: 'react-postprocessing',
    canonical_url: 'https://github.com/pmndrs/react-postprocessing',
    source_type: 'github_repository',
    source_repository: 'https://github.com/pmndrs/react-postprocessing.git',
    documentation_url: 'https://github.com/pmndrs/react-postprocessing',
    license: 'MIT',
    category: ['3d', 'webgl', 'cinematic'],
    representation: 'dependency',
    routing_tags: ['postprocessing', '3d', 'cinematic'],
    acquisition_method: 'npm',
    runtime_method: 'project_scoped_install',
    status: 'CURATED_OPTIONAL',
    trigger_conditions: [
      'cinematic 3D look named on the card',
      'bloom, SSAO, or similar post stacks',
    ],
    avoid_conditions: [
      'every 3D scene',
      'global npm',
    ],
    skip_conditions: [
      'simple R3F scene without cinematic grade',
      'performance budget forbids the extra pass',
    ],
    policy_verdict: '3D + CINEMATIC ONLY',
    rationale: 'Smallest sufficient stack: skip unless the card asks for cinematic grade.',
    npm_package: '@react-three/postprocessing',
    auth: 'NONE',
    free_capability: 'OSS postprocessing for R3F',
    paid_capability: 'none',
    fallback: 'drei',
  }),
  row({
    id: 'r3f-scroll-rig',
    name: 'r3f-scroll-rig',
    canonical_url: 'https://github.com/14islands/r3f-scroll-rig',
    source_type: 'github_repository',
    source_repository: 'https://github.com/14islands/r3f-scroll-rig.git',
    documentation_url: 'https://github.com/14islands/r3f-scroll-rig',
    license: 'MIT',
    category: ['3d', 'scroll'],
    representation: 'dependency',
    routing_tags: ['scroll-rig', 'scroll-driven-3d', 'r3f'],
    acquisition_method: 'npm',
    runtime_method: 'project_scoped_install',
    status: 'CURATED_OPTIONAL',
    trigger_conditions: [
      'scroll-driven 3D named on the card',
    ],
    avoid_conditions: [
      'every 3D portfolio',
      'global npm',
    ],
    skip_conditions: [
      'no scroll-linked scene',
      'Lenis + CSS already covers the motion',
    ],
    policy_verdict: 'SCROLL-DRIVEN 3D ONLY',
    rationale: 'Do not install on a static 3D hero.',
    npm_package: '@14islands/r3f-scroll-rig',
    auth: 'NONE',
    free_capability: 'OSS scroll-synced R3F',
    paid_capability: 'none',
    fallback: 'r3f',
  }),
  row({
    id: 'awesome-chatgpt-prompts',
    name: 'awesome-chatgpt-prompts',
    canonical_url: 'https://github.com/f/awesome-chatgpt-prompts',
    source_type: 'github_repository',
    source_repository: 'https://github.com/f/awesome-chatgpt-prompts.git',
    documentation_url: 'https://github.com/f/awesome-chatgpt-prompts',
    license: 'CC0-1.0',
    category: ['prompts', 'reference'],
    representation: 'prompt_library',
    routing_tags: ['prompts', 'ui-prompts', 'coding-prompts', 'agent-prompts'],
    acquisition_method: 'web_fetch',
    runtime_method: 'on_demand_research',
    status: 'REFERENCE',
    trigger_conditions: [
      'need a targeted public prompt from the UI, frontend, architecture, coding, agents, or design categories',
    ],
    avoid_conditions: [
      'dumping the whole list into context',
      'always-on prompt pack',
    ],
    skip_conditions: [
      'the job already has an Orchestra protocol for this step',
      'unrelated to those categories',
    ],
    policy_verdict: 'SEARCHABLE REFERENCE',
    rationale: 'Do not dump every prompt into active context. Search one category, use one prompt.',
    auth: 'NONE',
    free_capability: 'public prompt catalog',
    paid_capability: 'none',
    fallback: 'awesome-prompts',
  }),
  row({
    id: 'awesome-prompts',
    name: 'awesome-prompts',
    canonical_url: 'https://github.com/ai-boost/awesome-prompts',
    source_type: 'github_repository',
    source_repository: 'https://github.com/ai-boost/awesome-prompts.git',
    documentation_url: 'https://github.com/ai-boost/awesome-prompts',
    license: 'see-upstream',
    category: ['prompts', 'reference'],
    representation: 'prompt_library',
    routing_tags: ['prompts', 'design-prompts', 'architecture-prompts'],
    acquisition_method: 'web_fetch',
    runtime_method: 'on_demand_research',
    status: 'REFERENCE',
    trigger_conditions: [
      'second public prompt catalog when awesome-chatgpt-prompts has no fit',
    ],
    avoid_conditions: [
      'dumping the whole list into context',
    ],
    skip_conditions: [
      'already used one targeted prompt this turn',
      'Orchestra protocol already covers the step',
    ],
    policy_verdict: 'SEARCHABLE REFERENCE',
    rationale: 'Companion catalog. Targeted search only.',
    auth: 'NONE',
    free_capability: 'public prompt catalog',
    paid_capability: 'none',
    fallback: 'awesome-chatgpt-prompts',
  }),
  row({
    id: 'mermaid',
    name: 'Mermaid',
    canonical_url: 'https://github.com/mermaid-js/mermaid',
    source_type: 'github_repository',
    source_repository: 'https://github.com/mermaid-js/mermaid.git',
    documentation_url: 'https://mermaid.js.org/',
    license: 'MIT',
    category: ['diagrams', 'docs'],
    representation: 'dependency',
    routing_tags: ['mermaid', 'architecture-diagram', 'flowchart'],
    acquisition_method: 'npm',
    runtime_method: 'reference_only',
    status: 'ACTIVE',
    trigger_conditions: [
      'architecture diagram in README or docs',
      'diagram-generator produced Mermaid source',
    ],
    avoid_conditions: [
      'invented architecture',
      'diagrams as decoration with no repo grounding',
    ],
    skip_conditions: [
      'no structure in the repo to ground',
      'the human asked for a raster screenshot only',
    ],
    policy_verdict: 'GROUNDED DIAGRAM SOURCE',
    rationale: 'diagram-generator writes Mermaid; pretty-mermaid may render. Source of truth is the repo.',
    auth: 'NONE',
    free_capability: 'OSS diagram language',
    paid_capability: 'none',
    fallback: 'pretty-mermaid',
  }),
  row({
    id: 'devops-skills',
    name: 'devops-skills (reference only)',
    canonical_url: 'https://github.com/arjunprabhulal/devops-skills',
    source_type: 'github_repository',
    source_repository: 'https://github.com/arjunprabhulal/devops-skills.git',
    documentation_url: 'https://github.com/arjunprabhulal/devops-skills',
    license: 'see-upstream',
    category: ['devops', 'quarantined'],
    representation: 'quarantined_reference',
    routing_tags: ['devops-reference', 'ci-cd', 'docker', 'observability'],
    acquisition_method: 'none',
    runtime_method: 'reference_only',
    status: 'REFERENCE',
    trigger_conditions: [
      'human named a specific DevOps theme (CI/CD, Docker, deploy, observability, reliability) and existing Orchestra skills do not cover it',
    ],
    avoid_conditions: [
      'installing all 88 skills',
      'npx skills add --all',
      'second conductor',
    ],
    skip_conditions: [
      'ordinary coding turn',
      'an Orchestra skill already covers the job (ship-safe, strix, semgrep-adapter)',
    ],
    policy_verdict: 'QUARANTINED REFERENCE — EXTRACT THEMES ONLY',
    rationale: 'Do not vendor the dump. Use as a map: CI/CD, Docker, deployment, observability, security, reliability, infrastructure, performance.',
    auth: 'NONE',
    free_capability: 'public skill markdown as reference',
    paid_capability: 'none',
    fallback: 'ship-safe',
  }),
];

const rows = load('resources.json');
const byId = new Map(rows.map((r) => [r.id, r]));
let skipFilled = 0;
let triggerFilled = 0;
let upserted = 0;
let added = 0;

for (const r of rows) {
  if (!Array.isArray(r.trigger_conditions) || r.trigger_conditions.length === 0) {
    if (r.status === 'REJECTED') {
      r.trigger_conditions = ['human asked to dump-install this refused resource'];
    } else {
      r.trigger_conditions = ['the chosen route or the human named ' + r.id];
    }
    triggerFilled += 1;
  }
  if (!Array.isArray(r.skip_conditions) || r.skip_conditions.length === 0) {
    if (r.status === 'REJECTED') {
      r.skip_conditions = ['never load into runtime; quarantined or refused'];
    } else if (Array.isArray(r.avoid_conditions) && r.avoid_conditions.length > 0) {
      r.skip_conditions = [...r.avoid_conditions];
    } else {
      r.skip_conditions = ['the chosen route does not name this resource'];
    }
    skipFilled += 1;
  }
}

for (const [id, patch] of Object.entries(UPSERTS)) {
  const cur = byId.get(id);
  if (!cur) {
    console.error('upsert target missing:', id);
    process.exit(1);
  }
  Object.assign(cur, patch);
  upserted += 1;
}

for (const n of NEW_ROWS) {
  if (byId.has(n.id)) {
    Object.assign(byId.get(n.id), n);
    upserted += 1;
  } else {
    rows.push(n);
    byId.set(n.id, n);
    added += 1;
  }
}

const schema = load('schemas/resources.schema.json');
const item = schema.items;
item.properties.representation.enum = Array.from(new Set([
  ...item.properties.representation.enum,
  'plugin',
  'executor',
  'quarantined_reference',
  'project_tool',
  'resource',
]));
item.properties.auth = {
  type: 'string',
  enum: ['NONE', 'HOST_CONNECT', 'AUTH_REQUIRED', 'PURCHASE_PENDING'],
  description: 'Authentication boundary. Unauthorized is not HEALTHY.',
};
item.properties.free_capability = {
  type: 'string',
  description: 'What is legally usable without purchase',
};
item.properties.paid_capability = {
  type: 'string',
  description: 'What remains paid or Pro-gated',
};
item.properties.trigger_conditions.minItems = 1;
item.properties.skip_conditions.minItems = 1;
const req = new Set(item.required);
req.add('trigger_conditions');
req.add('skip_conditions');
item.required = [...req];

const graph = load('design-resource-graph.json');
graph.version = '3.3.1';
const d = graph.domains;
d.visual_research = Array.from(new Set([...(d.visual_research || []), 'lapa', 'saasframe', 'refero-public']));
d.interaction_components = Array.from(new Set([
  ...(d.interaction_components || []),
  'aceternity-ui',
  'tailkit',
  'twenty-first-dev',
  'magic-ui-mcp',
]));
d.webgl = Array.from(new Set([...(d.webgl || []), 'three-js', 'react-postprocessing', 'r3f-scroll-rig']));
d.prompt_references = ['awesome-chatgpt-prompts', 'awesome-prompts'];
d.devops_reference = ['devops-skills'];
d.diagrams = Array.from(new Set([...(d.diagrams || []), 'mermaid']));
d.saas_research = Array.from(new Set([...(d.saas_research || []), 'saasframe', 'refero-public']));

graph.capabilities['interaction-components'] = {
  name: 'Interaction component search',
  description: 'Search registered free component kits for one control. Inspect existing project components before generating another. Synthesize into the approved design system — not a collage.',
  primary_archetype: 'interaction_design',
  trigger_tags: [
    'interaction-components',
    'magnetic',
    'shimmer',
    'gradient-button',
    'ripple',
    'particle-button',
    'liquid-button',
    'glass-button',
    'command-button',
    'hold-button',
    'react-bits',
    'kokonut',
    'cult-ui',
    'hyperui',
    'magic-ui',
    'aceternity',
    'motion-primitives',
  ],
  discovery: ['interaction_research', 'interaction_components'],
  synthesis: ['design_synthesis'],
  implementation: ['interaction_components'],
  optional_extensions: ['motion_optional'],
  qa: ['playwright', 'impeccable'],
  anti_patterns: [
    'Generating a new button without inspecting existing project components',
    'Collage of every kit on one surface',
    'One permanent skill per button library',
  ],
  quality_bar: 'PREMIUM',
  risk_rank: 5,
  platform: 'web',
  trigger_conditions: [
    'The unit of work is a control (button, card, drawer, input), not a whole product',
    'The brief names a kit (shadcn, React Bits, Magic UI, Aceternity, Kokonut, Cult, HyperUI, Tailkit free, 21st free, Motion Primitives)',
    'The brief names a control kind: primary, secondary, destructive, loading, success, disabled, icon, magnetic, shimmer, gradient, ripple, particle, liquid, glass, command, hold, stateful',
    'Existing project components must be inspected before generating another',
  ],
  skip_conditions: [
    'The job is a full page, product, or app — route to premium-website, saas-dashboard, mobile-app, or 3d-portfolio first',
    'No product surface exists yet to attach the control to',
    '3D/WebGL is the point — route to 3d-portfolio',
    'The platform is native mobile — route to mobile-app',
  ],
};

const re = graph.capabilities['reverse-engineering'];
if (re) {
  const tags = new Set(re.trigger_tags || []);
  ['tint', 'recolor', 'keep-the-structure', 'screenshot-to-code', 'design-adaptation'].forEach((t) => tags.add(t));
  re.trigger_tags = [...tags];
  re.trigger_conditions = Array.from(new Set([
    ...(re.trigger_conditions || []),
    'A named reference, screenshot, or URL should keep structure while tokens change (tint)',
  ]));
  re.skip_conditions = [
    'No concrete URL was named — this is not a browsing route',
    'The intent is to copy branding, assets, copy, or source rather than principles',
    'A direction already exists and only needs implementing',
  ];
}

const oss = {
  version: '3.3.1',
  description: 'Searchable metadata for redistributable OSS. Not the Brain. Proprietary sites stay URL-only. No giant clones.',
  generated: TODAY,
  entries: [
    { id: 'shadcn-ui', source_url: 'https://github.com/shadcn-ui/ui', license: 'MIT', acquisition: 'metadata-only', local_copy: false, date: TODAY, capability: 'component foundation', categories: ['components'] },
    { id: 'react-bits', source_url: 'https://github.com/DavidHDev/react-bits', license: 'see-upstream', acquisition: 'metadata-only', local_copy: false, date: TODAY, capability: 'motion components', categories: ['motion', 'components'] },
    { id: 'magic-ui-oss', source_url: 'https://github.com/magicuidesign/magicui', license: 'MIT', acquisition: 'metadata-only', local_copy: false, date: TODAY, capability: 'animated UI OSS', categories: ['components'] },
    { id: 'kokonut-ui', source_url: 'https://github.com/kokonut-labs/kokonutui', license: 'see-upstream', acquisition: 'metadata-only', local_copy: false, date: TODAY, capability: 'interaction controls', categories: ['components'] },
    { id: 'motion-primitives', source_url: 'https://github.com/ibelick/motion-primitives', license: 'see-upstream', acquisition: 'metadata-only', local_copy: false, date: TODAY, capability: 'micro-interactions', categories: ['motion'] },
    { id: 'cult-ui', source_url: 'https://github.com/nolly-studio/cult-ui', license: 'see-upstream', acquisition: 'metadata-only', local_copy: false, date: TODAY, capability: 'animated blocks', categories: ['components'] },
    { id: 'three-js', source_url: 'https://github.com/mrdoob/three.js', license: 'MIT', acquisition: 'metadata-only', local_copy: false, date: TODAY, capability: 'WebGL renderer', categories: ['3d'] },
    { id: 'r3f', source_url: 'https://github.com/pmndrs/react-three-fiber', license: 'MIT', acquisition: 'metadata-only', local_copy: false, date: TODAY, capability: 'React Three Fiber', categories: ['3d'] },
    { id: 'drei', source_url: 'https://github.com/pmndrs/drei', license: 'MIT', acquisition: 'metadata-only', local_copy: false, date: TODAY, capability: 'R3F helpers', categories: ['3d'] },
    { id: 'react-postprocessing', source_url: 'https://github.com/pmndrs/react-postprocessing', license: 'MIT', acquisition: 'metadata-only', local_copy: false, date: TODAY, capability: 'cinematic 3D', categories: ['3d'] },
    { id: 'r3f-scroll-rig', source_url: 'https://github.com/14islands/r3f-scroll-rig', license: 'MIT', acquisition: 'metadata-only', local_copy: false, date: TODAY, capability: 'scroll-driven 3D', categories: ['3d'] },
    { id: 'gsap', source_url: 'https://github.com/greensock/gsap', license: 'see-upstream', acquisition: 'metadata-only', local_copy: false, date: TODAY, capability: 'timeline motion', categories: ['motion'] },
    { id: 'scrapling', source_url: 'https://github.com/D4Vinci/Scrapling', license: 'see-upstream', acquisition: 'metadata-only', local_copy: false, date: TODAY, capability: 'public web acquisition', categories: ['research'] },
    { id: 'mermaid', source_url: 'https://github.com/mermaid-js/mermaid', license: 'MIT', acquisition: 'metadata-only', local_copy: false, date: TODAY, capability: 'diagram source', categories: ['docs'] },
    { id: 'screenshot-to-code', source_url: 'https://github.com/abi/screenshot-to-code', license: 'see-upstream', acquisition: 'metadata-only', local_copy: false, date: TODAY, capability: 'reverse reference', categories: ['design'] },
  ],
};

const problems = [];
for (const r of rows) {
  if (!r.trigger_conditions?.length) problems.push(r.id + ': no trigger');
  if (!r.skip_conditions?.length) problems.push(r.id + ': no skip');
}
const ids = new Set(rows.map((r) => r.id));
for (const [dom, list] of Object.entries(graph.domains)) {
  for (const id of list) {
    if (!ids.has(id)) problems.push('dangling domain ' + dom + ':' + id);
  }
}

console.log('rows', rows.length);
console.log('skipFilled', skipFilled);
console.log('triggerFilled', triggerFilled);
console.log('upserted', upserted);
console.log('added', added);
console.log('capabilities', Object.keys(graph.capabilities).length);
console.log('problems', problems.length);
problems.forEach((p) => console.log(' ', p));
if (problems.length) process.exit(1);

if (!WRITE) {
  console.log('\ndry run. pass --write to apply.');
  process.exit(0);
}

save('resources.json', rows);
save('schemas/resources.schema.json', schema);
save('design-resource-graph.json', graph);
save('oss-index.json', oss);
console.log('wrote resources, schema, graph, oss-index');
