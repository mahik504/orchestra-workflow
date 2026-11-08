// Writes docs/skills.md from the skill manifest and the routing table.
//   node runtime/tools/build-skills-guide.mjs
import fs from 'node:fs';
import path from 'node:path';

const REPO = path.join(import.meta.dirname, '..', '..');
const skills = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'skills.json'), 'utf8')).skills;
const routing = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'routing.json'), 'utf8'));

const routesBySkill = new Map();
for (const route of routing.routes) {
  for (const id of route.load) {
    if (!routesBySkill.has(id)) routesBySkill.set(id, []);
    routesBySkill.get(id).push(route);
  }
}

const lines = [
  '# Skills',
  '',
  'Generated from `registries/skills.json` and `registries/routing.json`. Do not edit by hand. Run `node runtime/tools/build-skills-guide.mjs`.',
  '',
  'A skill you name wins unless it is unsafe. A route loads the skills in its row. The skip cell is the reason to leave a skill unloaded.',
  '',
];

for (const skill of skills) {
  const routes = routesBySkill.get(skill.id) || [];
  lines.push(`## ${skill.id}`);
  lines.push('');
  lines.push(skill.triggers || skill.id);
  lines.push('');
  lines.push(`Owner: ${skill.owner}. License: ${skill.license}.`);
  lines.push('');
  if (routes.length === 0) {
    lines.push('No route loads this skill by itself. It runs when you name it, or when an engineering row names it.');
  } else {
    lines.push('| Route | Skip |');
    lines.push('| --- | --- |');
    for (const route of routes) {
      lines.push(`| ${route.job} | ${route.skip} |`);
    }
  }
  lines.push('');
}

fs.writeFileSync(path.join(REPO, 'docs', 'skills.md'), lines.join('\n'));
console.log(`skills ${skills.length}`);
