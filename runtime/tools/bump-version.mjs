// Move the contract pin. Does not rewrite CHANGELOG history.
//   node runtime/tools/bump-version.mjs 5.5.0
import fs from 'node:fs';
import path from 'node:path';

const REPO = path.join(import.meta.dirname, '..', '..');
const next = process.argv[2];
if (!/^[0-9]+\.[0-9]+\.[0-9]+$/.test(next || '')) {
  console.error('usage: node runtime/tools/bump-version.mjs <major.minor.patch>');
  process.exit(1);
}

const versionsPath = path.join(REPO, 'registries', 'versions.json');
const versions = JSON.parse(fs.readFileSync(versionsPath, 'utf8'));
const from = versions.current;
if (from === next) {
  console.error(`already ${next}`);
  process.exit(1);
}

function read(rel) {
  return fs.readFileSync(path.join(REPO, rel), 'utf8');
}

function write(rel, text) {
  fs.writeFileSync(path.join(REPO, rel), text);
}

function replacePin(rel) {
  const text = read(rel);
  const nextText = text.split(`Orchestra ${from}`).join(`Orchestra ${next}`);
  if (nextText === text) {
    console.error(`${rel} has no Orchestra ${from}`);
    process.exit(1);
  }
  write(rel, nextText);
  console.log(`pin ${rel}`);
}

fs.writeFileSync(path.join(REPO, 'VERSION'), `${next}\n`);
versions.current = next;
versions.graph_version = next;
write('registries/versions.json', JSON.stringify(versions, null, 2) + '\n');
console.log('pin VERSION and versions.json');

for (const rel of ['registries/host-stack.json', 'registries/design-resource-graph.json']) {
  const text = read(rel);
  const needle = `"version": "${from}"`;
  if (!text.includes(needle)) {
    console.error(`${rel} has no ${needle}`);
    process.exit(1);
  }
  write(rel, text.replace(needle, `"version": "${next}"`));
  console.log(`pin ${rel}`);
}

for (const rel of [
  '.claude-plugin/plugin.json',
  '.cursor-plugin/plugin.json',
  '.codex-plugin/plugin.json',
  '.plugin/plugin.json',
]) {
  const text = read(rel);
  const needle = `"version": "${from}"`;
  if (!text.includes(needle)) {
    console.error(`${rel} has no ${needle}`);
    process.exit(1);
  }
  write(rel, text.replace(needle, `"version": "${next}"`));
  console.log(`pin ${rel}`);
}

for (const rel of [
  '.claude-plugin/marketplace.json',
  '.cursor-plugin/marketplace.json',
  '.agents/plugins/marketplace.json',
]) {
  const text = read(rel);
  const needle = `Version ${from}.`;
  if (!text.includes(needle)) {
    console.error(`${rel} has no ${needle}`);
    process.exit(1);
  }
  write(rel, text.replace(needle, `Version ${next}.`));
  console.log(`pin ${rel}`);
}

const live = [
  'AGENTS.md',
  'CLAUDE.md',
  '.cursorrules',
  'kit/antigravity/GEMINI.md',
  'kit/antigravity/MASTER-PROMPT.md',
  'kit/antigravity/ALWAYS-ON.md',
  'STACK.md',
  'WORKFLOW.md',
  'ARCHITECTURE.md',
  'kit/HOST_CAPABILITIES.md',
  'kit/ROLLBACK.md',
  'kit/bootstrap.ps1',
  'skills/ow-conductor/SKILL.md',
  'kit/codex/AGENTS.md',
  'kit/opencode/AGENTS.md',
  'kit/kimi/AGENTS.md',
  'README.md',
];
for (const rel of live) replacePin(rel);

const always = read('kit/antigravity/ALWAYS-ON.md');
const marker = `BEGIN ORCHESTRA ${from}`;
if (!always.includes(marker)) {
  console.error('ALWAYS-ON.md marker missing');
  process.exit(1);
}
write('kit/antigravity/ALWAYS-ON.md', always.replace(marker, `BEGIN ORCHESTRA ${next}`));

const readme = read('README.md');
const pinLine = `The pin in [\`VERSION\`](VERSION) is **${from}**.`;
if (!readme.includes(pinLine)) {
  console.error('README pin line missing');
  process.exit(1);
}
write('README.md', readme.replace(pinLine, `The pin in [\`VERSION\`](VERSION) is **${next}**.`));
console.log('pin README version line');

const skills = read('registries/skills.json');
const trigger = `Orchestra ${from} control plane.`;
if (!skills.includes(trigger)) {
  console.error('skills.json conductor trigger missing');
  process.exit(1);
}
write('registries/skills.json', skills.replace(trigger, `Orchestra ${next} control plane.`));
console.log('pin skills.json trigger');
console.log(`bumped ${from} -> ${next}`);
