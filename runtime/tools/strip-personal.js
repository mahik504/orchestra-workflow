// Maintenance script for the public resource registry.
//   1. Removes the operator's personal catalog rows (identity, products, career).
//   2. Fills honest provenance on rows that have no redistributable software licence.
//   3. Reports anything still missing, and any residual personal leak.
//
// "site-terms" means: no licence grant, governed by the publisher's terms, and we
// extract principles only. It is not a claim that the content is reusable.
//
// Usage: node runtime/tools/strip-personal.js [--write]
const fs = require('fs');
const path = require('path');

const REG = path.join(__dirname, '..', '..', 'registries', 'resources.json');

const DROP = new Set([
  // operator identity
  'mahik504-profile', 'mahik504-readme', 'aryan-dani-profile',
  // operator products
  'airlens', 'astroverse', 'lyra', 'evonex', 'xplain-ai', 'dbms-tsh', 'pbl',
  'smart-event-access', 'penfight', 'evocentric', 'planner-app',
  // operator infrastructure
  'orchestra-brain', 'orchestra-workflow', 'local-omniroute',
  // operator career
  'future-ready-talent', 'gcp-student-credits', 'hackathons-master-sheet',
  'leetsync', 'ai-job-search',
]);

// Licence defaults by source_type for rows we consult but never redistribute.
const BY_SOURCE_TYPE = {
  web_reference: 'site-terms',
  digital_archive: 'site-terms',
  component_gallery: 'site-terms',
  asset_archive: 'site-terms',
  standard_specification: 'vendor-documentation',
  specification: 'site-terms',
  mcp_server: 'vendor-service',
  local_runtime: 'not-applicable',
  agent_skill: 'MIT',
};
// Explicit overrides where the source_type default would be wrong.
const BY_ID = { 'motion-dev': 'MIT' };

const PROVENANCE = ['canonical_url', 'acquisition_method', 'license', 'status', 'rationale'];
const LEAK = /mahik504|AstroVerse|AirLens|YUMIT|ODYSS|penfight|EvoNex|XplainAI|C:\\\\(projects|Users)|orchestra-brain/i;

const rows = JSON.parse(fs.readFileSync(REG, 'utf8'));
const removed = rows.filter((r) => DROP.has(r.id));
const kept = rows.filter((r) => !DROP.has(r.id));

let filled = 0;
for (const r of kept) {
  if (r.license && String(r.license).trim()) continue;
  const lic = BY_ID[r.id] || BY_SOURCE_TYPE[r.source_type];
  if (lic) { r.license = lic; filled += 1; }
}

const gaps = kept
  .map((r) => ({ id: r.id, missing: PROVENANCE.filter((f) => !r[f] || String(r[f]).trim() === '') }))
  .filter((g) => g.missing.length > 0);
const leaks = kept.filter((r) => LEAK.test(JSON.stringify(r)));
const notFound = [...DROP].filter((id) => !rows.some((r) => r.id === id));

console.log(`rows in:              ${rows.length}`);
console.log(`personal removed:     ${removed.length}`);
console.log(`rows out:             ${kept.length}`);
console.log(`licences filled:      ${filled}`);
console.log(`provenance gaps left: ${gaps.length}`);
console.log(`personal leaks left:  ${leaks.length}`);
if (notFound.length) console.log(`drop-ids not found:   ${notFound.join(', ')}`);
for (const g of gaps) console.log(`  GAP ${g.id}: ${g.missing.join(', ')}`);
for (const l of leaks) console.log(`  LEAK ${l.id}`);

if (process.argv.includes('--write')) {
  fs.writeFileSync(REG, JSON.stringify(kept, null, 2) + '\n', 'utf8');
  console.log(`\nwrote ${REG}`);
} else {
  console.log('\ndry run. pass --write to apply.');
}
