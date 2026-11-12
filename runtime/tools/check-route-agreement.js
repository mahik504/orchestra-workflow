// Skills named by registries/routing.json must exist on the graph,
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

const routing = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'routing.json'), 'utf8'));
const routed = new Set();
for (const route of routing.routes || []) {
  for (const id of route.load || []) {
    if (RESERVED.has(id) || EXEMPT.has(id)) continue;
    routed.add(id);
  }
}
if (routed.size === 0) problems.push('routing.json produced zero skills');

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

let offGraph = 0;
for (const skill of routed) {
  const file = path.join(REPO, 'skills', skill, 'SKILL.md');
  if (!fs.existsSync(file)) problems.push(`routing skill ${skill} has no skills/${skill}/SKILL.md`);
  if (!onGraph.has(skill)) offGraph += 1;
}

if (problems.length) {
  console.log(problems.join('\n'));
  console.log(`route agreement skills ${onGraph.size}, routed ${routed.size}, off graph ${offGraph}, problems ${problems.length}`);
  process.exit(1);
}
console.log(`route agreement skills ${onGraph.size}, routed ${routed.size}, off graph ${offGraph}, problems 0`);
