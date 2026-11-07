// Render registries/routing.json into host adapters.
//   node runtime/tools/render-adapters.mjs
import fs from 'node:fs';
import path from 'node:path';

const REPO = path.join(import.meta.dirname, '..', '..');
const START = '<!-- routing:start -->';
const END = '<!-- routing:end -->';
const ENG_START = '<!-- engineering:start -->';
const ENG_END = '<!-- engineering:end -->';

export function routingMarkdown(doc) {
  const lines = [
    START,
    '## Routing',
    '',
    doc.preamble,
    '',
    '| Job | Load | Skip |',
    '| --- | --- | --- |',
  ];
  for (const route of doc.routes) {
    lines.push(`| ${route.job} | ${route.load.join(', ')} | ${route.skip} |`);
  }
  lines.push(END);
  return lines.join('\n');
}

export function engineeringMarkdown(doc) {
  const lines = [ENG_START, '| Job | Skill |', '| --- | --- |'];
  for (const [job, skill] of doc.engineering) {
    lines.push(`| ${job} | ${skill} |`);
  }
  lines.push(ENG_END);
  return lines.join('\n');
}

function replaceBlock(text, start, end, block) {
  const a = text.indexOf(start);
  const b = text.indexOf(end);
  if (a < 0 || b < a) return null;
  return text.slice(0, a) + block + text.slice(b + end.length);
}

const version = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'versions.json'), 'utf8')).current;
const doc = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'routing.json'), 'utf8'));
const routing = routingMarkdown(doc);
const engineering = engineeringMarkdown(doc);

const copies = [
  'AGENTS.md',
  'CLAUDE.md',
  'kit/antigravity/GEMINI.md',
  'skills/ow-conductor/SKILL.md',
  '.cursorrules',
];

for (const rel of copies) {
  const file = path.join(REPO, rel);
  if (!fs.existsSync(file)) continue;
  const next = replaceBlock(fs.readFileSync(file, 'utf8').replace(/\r\n/g, '\n'), START, END, routing);
  if (!next) {
    console.error(`${rel} has no routing markers`);
    process.exit(1);
  }
  fs.writeFileSync(file, next);
  console.log(`rendered ${rel}`);
}

const agentsPath = path.join(REPO, 'AGENTS.md');
const agentsNext = replaceBlock(fs.readFileSync(agentsPath, 'utf8'), ENG_START, ENG_END, engineering);
if (!agentsNext) {
  console.error('AGENTS.md has no engineering markers');
  process.exit(1);
}
fs.writeFileSync(agentsPath, agentsNext);

function adapter(title, line) {
  return `# ${title} — Orchestra ${version}\n\nYou are running **Orchestra ${version}**. This file is an adapter. It does not invent a second loop.\n\n${routing}\n\n${line}\n\nSay **skip orchestra** to stand this down. Plain text. Not a slash command.\n`;
}

const fresh = [
  ['kit/codex/AGENTS.md', 'Codex adapter', 'Codex reads this file. The conductor is still Orchestra.'],
  ['kit/opencode/AGENTS.md', 'OpenCode adapter', 'OpenCode reads this file. The conductor is still Orchestra.'],
  ['kit/kimi/AGENTS.md', 'Kimi adapter', 'Kimi reads this file. The conductor is still Orchestra.'],
];
for (const [rel, title, line] of fresh) {
  const file = path.join(REPO, rel);
  fs.mkdirSync(path.dirname(file), { recursive: true });
  fs.writeFileSync(file, adapter(title, line));
  console.log(`wrote ${rel}`);
}
