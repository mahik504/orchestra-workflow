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
    '# Route table',
    '',
    'Read this file when you classify. It is not always loaded.',
    '',
    START,
    '## Routing',
    '',
    doc.preamble,
    '',
    '| Job | Load | Skip |',
    '| --- | --- | --- |',
  ];
  for (const route of doc.routes) {
    const session = route.session ? ` <!-- ${route.session} -->` : '';
    lines.push(`| ${route.job} | ${route.load.join(', ')} | ${route.skip} |${session}`);
  }
  lines.push(END);
  return lines.join('\n');
}

export function routingPointer() {
  return [
    START,
    '## Routing',
    '',
    'The route table is `docs/ROUTING.md`. Read that file when you classify. Name the skills, MCPs, plugins, and libraries that row loads. Omit that line when nothing new is loaded. A missing tool stops. It does not install itself.',
    END,
  ].join('\n');
}

export function rulesMarkdown(doc) {
  const lines = ['<!-- rules:start -->', '## Shared rules', ''];
  for (const rule of doc.rules || []) lines.push(`- ${rule}`);
  lines.push('<!-- rules:end -->');
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
const routing = routingPointer();
const rules = rulesMarkdown(doc);
const engineering = engineeringMarkdown(doc);
fs.writeFileSync(path.join(REPO, 'docs', 'ROUTING.md'), routingMarkdown(doc) + '\n');
console.log('wrote docs/ROUTING.md');

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
  let withRules = next;
  if (next.includes('<!-- rules:start -->')) {
    withRules = replaceBlock(next, '<!-- rules:start -->', '<!-- rules:end -->', rules);
    if (!withRules) {
      console.error(`${rel} has a broken rules block`);
      process.exit(1);
    }
  }
  fs.writeFileSync(file, withRules);
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

const cmdSrc = path.join(REPO, 'commands');
const cmdDst = path.join(REPO, '.cursor', 'commands');
fs.mkdirSync(cmdDst, { recursive: true });
const cmdNames = new Set(fs.readdirSync(cmdSrc).filter((name) => name.endsWith('.md')));
for (const name of cmdNames) {
  fs.copyFileSync(path.join(cmdSrc, name), path.join(cmdDst, name));
}
for (const name of fs.readdirSync(cmdDst)) {
  if (name.endsWith('.md') && !cmdNames.has(name)) fs.rmSync(path.join(cmdDst, name));
}
console.log(`commands ${cmdNames.size}`);
