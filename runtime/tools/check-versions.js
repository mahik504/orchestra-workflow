// Current contract version lives in registries/versions.json.
//   node runtime/tools/check-versions.js
const fs = require('fs');
const path = require('path');

const REPO = path.join(__dirname, '..', '..');
const manifest = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'versions.json'), 'utf8'));
const problems = [];
const current = manifest.current;

if (!/^[0-9]+\.[0-9]+\.[0-9]+$/.test(current || '')) {
  problems.push('versions.json current is not major.minor.patch');
}

const versionFile = fs.readFileSync(path.join(REPO, 'VERSION'), 'utf8').trim();
if (versionFile !== current) problems.push(`VERSION is "${versionFile}", expected ${current}`);

const stack = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'host-stack.json'), 'utf8'));
if (stack.version !== current) problems.push(`host-stack version is "${stack.version}", expected ${current}`);

const graph = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'design-resource-graph.json'), 'utf8'));
if (graph.version !== manifest.graph_version) {
  problems.push(`graph.version is "${graph.version}", expected ${manifest.graph_version}`);
}

for (const rel of manifest.live_files || []) {
  const file = path.join(REPO, rel);
  if (rel === 'VERSION') continue;
  if (!fs.existsSync(file)) {
    problems.push(`live file missing: ${rel}`);
    continue;
  }
  const text = fs.readFileSync(file, 'utf8');
  if (!text.includes(current)) problems.push(`${rel} does not name ${current}`);
}

const pinned = String(current);
for (const rel of manifest.readers_must_not_hardcode_current || []) {
  const text = fs.readFileSync(path.join(REPO, rel), 'utf8');
  if (text.includes(pinned)) problems.push(`${rel} hardcodes ${pinned}; read registries/versions.json`);
}

console.log(`current:  ${current}`);
console.log(`graph:    ${manifest.graph_version}`);
console.log(`problems: ${problems.length}`);
for (const problem of problems) console.log(`  ${problem}`);
process.exit(problems.length ? 1 : 0);
