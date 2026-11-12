// Measure what a host would load at each disclosure level.
//   node runtime/tools/measure-context.mjs
import fs from 'node:fs';
import path from 'node:path';

const REPO = path.join(import.meta.dirname, '..', '..');

function chars(file) {
  if (!fs.existsSync(file)) return 0;
  return fs.readFileSync(file, 'utf8').length;
}

function descriptionChars() {
  let total = 0;
  let ow = 0;
  let over = [];
  const dir = path.join(REPO, 'skills');
  for (const name of fs.readdirSync(dir)) {
    const file = path.join(dir, name, 'SKILL.md');
    if (!fs.existsSync(file)) continue;
    const text = fs.readFileSync(file, 'utf8');
    const match = text.match(/^description:\s*(.+)$/m);
    const n = match ? match[1].trim().length : 0;
    total += n;
    if (name.startsWith('ow-')) {
      ow += n;
      if (n > 300) over.push(`${name} ${n}`);
    }
  }
  return { total, ow, over };
}

function routeBytes() {
  const rows = [];
  for (const route of routing.routes) {
    if (!['Award-style site', 'Full-stack app', 'Mobile app'].includes(route.job)) continue;
    let bytes = 0;
    for (const id of route.load) {
      const file = path.join(REPO, 'skills', id, 'SKILL.md');
      if (fs.existsSync(file)) bytes += fs.statSync(file).size;
    }
    rows.push({ job: route.job, skills: route.load.length, bytes });
  }
  return rows;
}

const agents = chars(path.join(REPO, 'AGENTS.md'));
const routing = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'routing.json'), 'utf8'));
const routingChars = JSON.stringify(routing.routes.map((route) => ({
  job: route.job,
  load: route.load,
  skip: route.skip,
}))).length;
const desc = descriptionChars();
const levels = {
  L0: 'task and acceptance only',
  L1: 'route row from registries/routing.json',
  L2: 'selected SKILL.md',
  L3: 'references/ beside that skill',
  L4: 'tool output from the run',
};

const lines = [
  '# Context budget',
  '',
  'Measured by `node runtime/tools/measure-context.mjs` on this tree.',
  '',
  `| What | Characters |`,
  `| --- | ---: |`,
  `| AGENTS.md | ${agents} |`,
  `| Routing rows as JSON | ${routingChars} |`,
  `| Skill descriptions, all canonical skills | ${desc.total} |`,
  `| Skill descriptions, ow-* only | ${desc.ow} |`,
  `| ow-* descriptions over 300 | ${desc.over.length} |`,
  '',
  'The 4.2.0 tree had 78 canonical skills. This tree\'s count is in `registries/skills.json`. The slot cap of 80 is gone. The budget is the description length: 300 characters for an `ow-*` skill.',
  '',
  '## Disclosure',
  '',
  '| Level | What the worker sees |',
  '| --- | --- |',
];
for (const [id, text] of Object.entries(levels)) {
  lines.push(`| ${id} | ${text} |`);
}
lines.push('');
lines.push('A STANDARD job that names the file stays at L0. The conductor fast path does not print a brief or a receipt. A visual job starts at L1 and opens L2 only for the skills the route loaded.');
lines.push('');
if (desc.over.length) {
  lines.push('Over budget:');
  for (const row of desc.over) lines.push(`- ${row}`);
  lines.push('');
}

const heavy = routeBytes();
lines.push('## Heavy routes');
lines.push('');
lines.push('| Route | Skills | SKILL.md bytes |');
lines.push('| --- | ---: | ---: |');
for (const row of heavy) {
  lines.push(`| ${row.job} | ${row.skills} | ${row.bytes} |`);
}
lines.push('');
lines.push('These sizes are a measurement. They are not a cap.');
lines.push('');

const out = path.join(REPO, 'docs', 'CONTEXT-BUDGET.md');
const check = process.argv.includes('--check');
console.log(`AGENTS.md ${agents}`);
console.log(`descriptions ${desc.total}`);
console.log(`ow over 300: ${desc.over.length}`);
for (const row of heavy) console.log(`${row.job} ${row.skills} skills ${row.bytes} bytes`);
if (check) {
  let failed = false;
  if (agents > 3072) {
    console.log(`AGENTS.md is ${agents} characters; cap is 3072`);
    failed = true;
  }
  const existing = fs.existsSync(out) ? fs.readFileSync(out, 'utf8') : '';
  const baseline = Number((existing.match(/Skill descriptions, all canonical skills \| (\d+)/) || [])[1] || desc.total);
  if (desc.total > baseline) {
    console.log(`descriptions ${desc.total} exceed the recorded baseline ${baseline}`);
    failed = true;
  }
  if (desc.over.length) failed = true;
  if (failed) process.exit(1);
  console.log('check ok');
} else {
  fs.writeFileSync(out, lines.join('\n'));
  console.log(`wrote ${out}`);
  if (desc.over.length) process.exit(1);
}
