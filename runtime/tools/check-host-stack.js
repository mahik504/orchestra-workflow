// Clone-safe host-stack integrity.
//   node runtime/tools/check-host-stack.js
const fs = require('fs');
const path = require('path');

const REPO = path.join(__dirname, '..', '..');
const stackPath = path.join(REPO, 'registries', 'host-stack.json');
const skillsDir = path.join(REPO, 'skills');
const problems = [];

if (!fs.existsSync(stackPath)) {
  console.error('missing registries/host-stack.json');
  process.exit(1);
}

const stack = JSON.parse(fs.readFileSync(stackPath, 'utf8'));

if (stack.version !== '3.1.0') {
  problems.push(`host-stack version is "${stack.version}", expected 3.1.0`);
}

const listed = Array.isArray(stack.skills) ? stack.skills : [];
const dirs = fs.readdirSync(skillsDir).filter((n) => {
  const p = path.join(skillsDir, n);
  return fs.statSync(p).isDirectory() && fs.existsSync(path.join(p, 'SKILL.md'));
}).sort();

const listedSorted = [...listed].sort();
if (listedSorted.length !== dirs.length) {
  problems.push(`skill count mismatch: host-stack has ${listedSorted.length}, skills/ has ${dirs.length}`);
}
for (const id of listedSorted) {
  if (!dirs.includes(id)) problems.push(`host-stack skill '${id}' missing under skills/`);
}
for (const id of dirs) {
  if (!listedSorted.includes(id)) problems.push(`skills/${id} not listed in host-stack.json`);
}

const LEAK = /mahik504(?!\/orchestra-workflow)|AstroVerse|AirLens|YUMIT|ODYSS|EvoNex|XplainAI|orchestra-brain|career\.md|C:\\\\Users|C:\\\\projects|gehlot\.mahisingh/i;
const txt = fs.readFileSync(stackPath, 'utf8');
const m = txt.match(LEAK);
if (m) problems.push(`personal marker '${m[0]}' in host-stack.json`);

if (!stack.resources || stack.resources.see !== 'registries/resources.json') {
  problems.push('host-stack.resources.see must be registries/resources.json');
}

console.log(`host-stack skills: ${listed.length}`);
console.log(`skills/ dirs:      ${dirs.length}`);
console.log(`problems:          ${problems.length}`);
for (const p of problems) console.log(`  ${p}`);
process.exit(problems.length ? 1 : 0);
