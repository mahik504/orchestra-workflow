// The routing table is copied by hand into the host files.
// This check fails when a Load cell names a skill that is not in skills/,
// or when the four copies drift. It does not ask a model anything.
const fs = require('fs');
const path = require('path');

const REPO = path.join(__dirname, '..', '..');
const FILES = [
  'AGENTS.md',
  'CLAUDE.md',
  'kit/antigravity/GEMINI.md',
  'skills/orchestra-conductor/SKILL.md',
];
const START = '<!-- routing:start -->';
const END = '<!-- routing:end -->';
const NOT_SKILLS = new Set(['not ready', 'host-provided', 'none']);

// Document check only. These phrases must hit a row. They are not a model trial.
const PROMPTS = [
  { says: 'award-style site', wants: 'impeccable' },
  { says: 'bug', wants: 'systematic-debugging' },
  { says: 'paper', wants: 'research-paper' },
  { says: 'apk', wants: 'not ready' },
  { says: 'saas', wants: 'backend-architecture' },
];

function section(file) {
  const text = fs.readFileSync(path.join(REPO, file), 'utf8').replace(/\r\n/g, '\n');
  const a = text.indexOf(START);
  const b = text.indexOf(END);
  if (a < 0 || b < 0 || b < a) return null;
  return text.slice(a, b + END.length);
}

const problems = [];
const blocks = FILES.map((file) => {
  const block = section(file);
  if (!block) problems.push(`${file} missing routing block`);
  return block;
});

if (blocks[0]) {
  for (let i = 1; i < blocks.length; i++) {
    if (blocks[i] && blocks[i] !== blocks[0]) problems.push(`${FILES[i]} routing block differs from AGENTS.md`);
  }

  const rows = [];
  for (const line of blocks[0].split('\n')) {
    if (!line.startsWith('|') || line.includes('---') || line.includes('Job |')) continue;
    const cells = line.split('|').slice(1, -1).map((cell) => cell.trim());
    if (cells.length < 2) continue;
    rows.push({ job: cells[0].toLowerCase(), load: cells[1] });
    for (const raw of cells[1].split(',')) {
      const id = raw.trim();
      if (!id || NOT_SKILLS.has(id)) continue;
      if (!/^[a-z0-9-]+$/.test(id)) {
        problems.push(`bad Load token "${id}"`);
        continue;
      }
      const skill = path.join(REPO, 'skills', id, 'SKILL.md');
      if (!fs.existsSync(skill)) problems.push(`Load names missing skill ${id}`);
    }
  }

  for (const prompt of PROMPTS) {
    const row = rows.find((item) => item.job.includes(prompt.says));
    if (!row) {
      problems.push(`document check missed a row for "${prompt.says}"`);
      continue;
    }
    if (!row.load.includes(prompt.wants)) {
      problems.push(`document check: "${prompt.says}" does not load ${prompt.wants}`);
    }
  }
}

console.log('routing copies:    ' + FILES.length);
console.log('document prompts:  ' + PROMPTS.length + ' (table check, not a model trial)');
console.log('problems:          ' + problems.length);
for (const problem of problems) console.log('  ' + problem);
process.exit(problems.length ? 1 : 0);
