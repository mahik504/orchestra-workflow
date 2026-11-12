// The routing table is copied by hand into the host files.
// This check fails when a Load cell names a skill that is not in skills/,
// when Skip text is missing, or when the four copies drift.
// It does not ask a model anything. Document checks are not model trials.
const fs = require('fs');
const path = require('path');

const REPO = path.join(__dirname, '..', '..');
const FILES = [
  'AGENTS.md',
  'CLAUDE.md',
  'kit/antigravity/GEMINI.md',
  'skills/ow-conductor/SKILL.md',
  'skills/ow-conductor/SKILL.md',
  '.cursorrules',
  'kit/codex/AGENTS.md',
  'kit/opencode/AGENTS.md',
  'kit/kimi/AGENTS.md',
].filter((rel) => fs.existsSync(path.join(REPO, rel)));
const START = '<!-- routing:start -->';
const END = '<!-- routing:end -->';
const RESERVED = new Set(['not ready', 'host-provided', 'none']);

const routingDoc = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'routing.json'), 'utf8'));
const EXPECTED = routingDoc.routes.map((route) => ({
  job: route.job,
  load: route.load,
  skip: route.skip_phrases,
}));

function sectionOf(file, start, end) {
  const text = fs.readFileSync(path.join(REPO, file), 'utf8').replace(/\r\n/g, '\n');
  const a = text.indexOf(start);
  const b = text.indexOf(end);
  if (a < 0 || b < a) return null;
  return text.slice(a, b + end.length);
}

function section(file) {
  const text = fs.readFileSync(path.join(REPO, file), 'utf8').replace(/\r\n/g, '\n');
  const starts = text.split(START).length - 1;
  const ends = text.split(END).length - 1;
  if (starts !== 1 || ends !== 1) return { block: null, starts, ends };
  const a = text.indexOf(START);
  const b = text.indexOf(END);
  if (b < a) return { block: null, starts, ends };
  return { block: text.slice(a, b + END.length), starts, ends };
}

function rowsOf(block) {
  const rows = [];
  for (const line of block.split('\n')) {
    if (!line.startsWith('|') || line.includes('---') || line.includes('Job |')) continue;
    const cells = line.split('|').slice(1, -1).map((cell) => cell.trim());
    if (cells.length < 3) continue;
    rows.push({
      job: cells[0],
      load: cells[1].split(',').map((item) => item.trim()).filter(Boolean),
      skip: cells[2],
    });
  }
  return rows;
}

const problems = [];
const blocks = FILES.map((file) => {
  const found = section(file);
  if (found.starts !== 1 || found.ends !== 1 || !found.block) {
    problems.push(`${file} needs one routing start and one end (starts ${found.starts}, ends ${found.ends})`);
  }
  return found.block;
});

const tableFile = 'docs/ROUTING.md';
const table = section(tableFile);
if (!table.block) problems.push('docs/ROUTING.md needs one routing table');

if (blocks[0]) {
  for (let i = 1; i < blocks.length; i++) {
    if (blocks[i] && blocks[i] !== blocks[0]) problems.push(`${FILES[i]} routing block differs from AGENTS.md`);
  }
  if (blocks[0].includes('| Job |')) problems.push('always-on adapters must point at docs/ROUTING.md and not inline the table');

  const rows = rowsOf(table.block || '');
  const jobs = rows.map((row) => row.job);
  const seen = new Set();
  for (const job of jobs) {
    const key = job.toLowerCase();
    if (seen.has(key)) problems.push(`duplicate job "${job}"`);
    seen.add(key);
  }
  let skillPaths = 0;
  for (const row of rows) {
    for (const id of row.load) {
      if (RESERVED.has(id)) continue;
      if (!/^[a-z0-9-]+$/.test(id)) {
        problems.push(`bad Load token "${id}"`);
        continue;
      }
      skillPaths += 1;
      const skill = path.join(REPO, 'skills', id, 'SKILL.md');
      if (!fs.existsSync(skill)) problems.push(`Load names missing skill ${id}`);
    }
    for (const id of row.load) {
      if (!RESERVED.has(id) && !/^[a-z0-9-]+$/.test(id)) continue;
      if (![...RESERVED].includes(id) && !/^[a-z0-9-]+$/.test(id)) {
        problems.push(`Load token "${id}" is neither reserved nor a skill id`);
      }
    }
  }

  let assertions = 0;
  for (const expected of EXPECTED) {
    assertions += 1;
    const row = rows.find((item) => item.job.toLowerCase() === expected.job.toLowerCase());
    if (!row) {
      problems.push(`missing row "${expected.job}"`);
      continue;
    }
    if (row.load.join(', ') !== expected.load.join(', ')) {
      problems.push(`"${expected.job}" load is "${row.load.join(', ')}", expected "${expected.load.join(', ')}"`);
    }
    for (const phrase of expected.skip) {
      assertions += 1;
      if (!row.skip.toLowerCase().includes(phrase.toLowerCase())) {
        problems.push(`"${expected.job}" skip missing "${phrase}"`);
      }
    }
  }

  const loose = rows.filter((row) => row.job.toLowerCase() === 'app');
  if (loose.length) problems.push('a bare job named "app" would steal "Ship a deployed app"');
  const ship = rows.find((row) => row.job === 'Ship a deployed app');
  if (ship && rows.some((row) => row.job.toLowerCase() === 'app')) {
    problems.push('substring route "app" matches a real job');
  }

  const conductorRel = fs.existsSync(path.join(REPO, 'skills/ow-conductor/SKILL.md'))
    ? 'skills/ow-conductor/SKILL.md'
    : 'skills/ow-conductor/SKILL.md';
  const conductor = fs.readFileSync(path.join(REPO, conductorRel), 'utf8');
  const outside = conductor.replace(/<!-- routing:start -->[\s\S]*<!-- routing:end -->/, '');
  if (!outside.includes('interaction-components')) {
    problems.push('conductor body outside the routing block does not name interaction-components');
  }
  if (!/One control loads `none`/.test(outside)) {
    problems.push('conductor body does not say One control loads none');
  }

  const agents = fs.readFileSync(path.join(REPO, 'AGENTS.md'), 'utf8').replace(/\r\n/g, '\n');
  const engStart = '<!-- engineering:start -->';
  const engEnd = '<!-- engineering:end -->';
  const engA = agents.indexOf(engStart);
  const engB = agents.indexOf(engEnd);
  if (engA < 0 || engB < engA) {
    problems.push('AGENTS.md needs one engineering table between markers');
  } else {
    const block = agents.slice(engA, engB);
    const named = new Set();
    for (const match of block.matchAll(/`([a-z0-9-]+)`/g)) named.add(match[1]);
    if (named.size === 0) problems.push('engineering table names no skills');
    for (const id of named) {
      const skill = path.join(REPO, 'skills', id, 'SKILL.md');
      if (!fs.existsSync(skill)) problems.push(`engineering table names missing skill ${id}`);
    }
  }

  const RULES_START = '<!-- rules:start -->';
  const RULES_END = '<!-- rules:end -->';
  const ruleFiles = ['.cursorrules', 'CLAUDE.md', 'kit/antigravity/GEMINI.md'];
  const ruleBlocks = ruleFiles.map((file) => sectionOf(file, RULES_START, RULES_END));
  if (!ruleBlocks[0]) problems.push('.cursorrules needs one shared rules block');
  for (let i = 1; i < ruleBlocks.length; i++) {
    if (ruleBlocks[i] !== ruleBlocks[0]) problems.push(`${ruleFiles[i]} rules block differs from .cursorrules`);
  }
  const geminiPath = path.join(REPO, 'kit/antigravity/GEMINI.md');
  if (fs.existsSync(geminiPath) && Buffer.byteLength(fs.readFileSync(geminiPath)) > 24 * 1024) {
    problems.push('kit/antigravity/GEMINI.md exceeds 24KB');
  }
  const caps = {
    'AGENTS.md': 3 * 1024,
    'CLAUDE.md': 3 * 1024,
    '.cursorrules': 3 * 1024,
    'kit/antigravity/GEMINI.md': 3 * 1024,
    'kit/codex/AGENTS.md': 3 * 1024,
    'kit/opencode/AGENTS.md': 3 * 1024,
    'kit/kimi/AGENTS.md': 3 * 1024,
    'skills/ow-conductor/SKILL.md': 6 * 1024,
  };
  for (const [rel, cap] of Object.entries(caps)) {
    const file = path.join(REPO, rel);
    if (!fs.existsSync(file)) continue;
    const bytes = Buffer.byteLength(fs.readFileSync(file));
    if (bytes > cap) problems.push(`${rel} is ${bytes} bytes; cap is ${cap}`);
  }

  const cmdSrc = path.join(REPO, 'commands');
  const cmdDst = path.join(REPO, '.cursor', 'commands');
  if (fs.existsSync(cmdSrc) && fs.existsSync(cmdDst)) {
    const names = fs.readdirSync(cmdSrc).filter((name) => name.endsWith('.md')).sort();
    for (const name of names) {
      const a = fs.readFileSync(path.join(cmdSrc, name));
      const bPath = path.join(cmdDst, name);
      if (!fs.existsSync(bPath) || !a.equals(fs.readFileSync(bPath))) {
        problems.push(`.cursor/commands/${name} is not a copy of commands/${name}`);
      }
    }
  }

  console.log('routing copies:    ' + FILES.length);
  console.log('required routes:   ' + EXPECTED.length);
  console.log('skill-path checks: ' + skillPaths);
  console.log('load/skip asserts: ' + assertions);
  console.log('document check:    not a model trial');
  console.log('problems:          ' + problems.length);
  for (const problem of problems) console.log('  ' + problem);
  if (problems.length) process.exit(1);
} else {
  console.log('problems:          ' + problems.length);
  for (const problem of problems) console.log('  ' + problem);
  process.exit(1);
}
