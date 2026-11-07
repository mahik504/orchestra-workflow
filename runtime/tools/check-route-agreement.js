// Routing-table skills must appear on a capability skills array,
// and every graph skill id must have skills/<id>/SKILL.md.
// ow-forge and ow-focus run before classify.
// systematic-debugging is selected from the bug type, not a new route.
const fs = require('fs');
const path = require('path');

const REPO = path.join(__dirname, '..', '..');
const START = '<!-- routing:start -->';
const END = '<!-- routing:end -->';
const RESERVED = new Set(['none', 'not ready', 'host-provided']);
const EXEMPT = new Set(['ow-forge', 'ow-focus', 'systematic-debugging']);

const problems = [];
const conductor = fs.readFileSync(path.join(REPO, 'skills', 'ow-conductor', 'SKILL.md'), 'utf8').replace(/\r\n/g, '\n');
const a = conductor.indexOf(START);
const b = conductor.indexOf(END);
if (a < 0 || b < a) problems.push('conductor routing block missing');
const block = a >= 0 ? conductor.slice(a, b) : '';

const routed = new Set();
for (const line of block.split('\n')) {
  if (!line.startsWith('|') || line.includes('---') || line.includes('Job |')) continue;
  const cells = line.split('|').slice(1, -1).map((cell) => cell.trim());
  if (cells.length < 2) continue;
  for (const token of cells[1].split(',').map((item) => item.trim()).filter(Boolean)) {
    if (RESERVED.has(token) || EXEMPT.has(token)) continue;
    routed.add(token);
  }
}

const graph = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'design-resource-graph.json'), 'utf8'));
const onGraph = new Set();
for (const [id, cap] of Object.entries(graph.capabilities || {})) {
  const skills = Array.isArray(cap.skills) ? cap.skills : [];
  if (id === 'interaction-components' && skills.length !== 0) {
    problems.push('interaction-components skills must stay empty');
  }
  if (skills.includes('ow-video') && skills.includes('gsap-core')) {
    problems.push(`${id} lists both ow-video and gsap-core`);
  }
  for (const skill of skills) {
    onGraph.add(skill);
    const file = path.join(REPO, 'skills', skill, 'SKILL.md');
    if (!fs.existsSync(file)) problems.push(`graph skill ${skill} has no skills/${skill}/SKILL.md`);
  }
}

for (const skill of routed) {
  if (!onGraph.has(skill)) problems.push(`routing skill ${skill} is absent from the graph`);
}

if (problems.length) {
  console.log(problems.join('\n'));
  console.log(`route agreement problems ${problems.length}`);
  process.exit(1);
}
console.log(`route agreement skills ${onGraph.size}, routed ${routed.size}, problems 0`);
