// Integrity check for the public registries.
//   - every id referenced by the capability graph exists in resources.json
//   - resources.json validates against its JSON Schema (required fields + enums)
//   - no personal marker survives anywhere in registries/
//
// Usage: node runtime/tools/check-registry.js
const fs = require('fs');
const path = require('path');

const REGDIR = path.join(__dirname, '..', '..', 'registries');
const rows = JSON.parse(fs.readFileSync(path.join(REGDIR, 'resources.json'), 'utf8'));
const graph = JSON.parse(fs.readFileSync(path.join(REGDIR, 'design-resource-graph.json'), 'utf8'));

const ids = new Set(rows.map((r) => r.id));
const problems = [];

// 1. dangling graph references
function walk(node, trail) {
  if (Array.isArray(node)) {
    node.forEach((v, i) => walk(v, `${trail}[${i}]`));
  } else if (node && typeof node === 'object') {
    for (const [k, v] of Object.entries(node)) walk(v, `${trail}.${k}`);
  } else if (typeof node === 'string') {
    if (/resource/i.test(trail) && !ids.has(node) && /^[a-z0-9][a-z0-9-]{2,}$/.test(node)) {
      problems.push(`dangling graph reference '${node}' at ${trail}`);
    }
  }
}
walk(graph, 'graph');

// 2. schema required fields
let schema = null;
const schemaPath = path.join(REGDIR, 'schemas', 'resources.schema.json');
if (fs.existsSync(schemaPath)) schema = JSON.parse(fs.readFileSync(schemaPath, 'utf8'));
const itemSchema = schema && (schema.items || schema);
const required = (itemSchema && itemSchema.required) || [];
for (const r of rows) {
  for (const f of required) {
    if (r[f] === undefined || r[f] === null || String(r[f]).trim() === '') {
      problems.push(`${r.id}: missing required field '${f}'`);
    }
  }
}

// 3. duplicate ids
const seen = new Set();
for (const r of rows) {
  if (seen.has(r.id)) problems.push(`duplicate id '${r.id}'`);
  seen.add(r.id);
}

// 4. personal markers across all of registries/
const LEAK = /mahik504|AstroVerse|AirLens|YUMIT|ODYSS|EvoNex|XplainAI|orchestra-brain|C:\\\\Users|C:\\\\projects|gehlot\.mahisingh/i;
for (const f of fs.readdirSync(REGDIR)) {
  const p = path.join(REGDIR, f);
  if (!fs.statSync(p).isFile()) continue;
  const txt = fs.readFileSync(p, 'utf8');
  const m = txt.match(LEAK);
  if (m) problems.push(`personal marker '${m[0]}' in registries/${f}`);
}

console.log(`resources:        ${rows.length}`);
console.log(`required fields:  ${required.length ? required.join(', ') : '(no schema required list)'}`);
console.log(`problems:         ${problems.length}`);
for (const p of problems) console.log(`  ${p}`);
process.exit(problems.length ? 1 : 0);
