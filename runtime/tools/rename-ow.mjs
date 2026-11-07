// One-shot rename of Orchestra-owned skill directories to ow-*.
// Aliases stay in registries/skill-aliases.json. The generator keeps the old keys.
import { execSync } from 'node:child_process';
import fs from 'node:fs';
import path from 'node:path';

const REPO = path.join(import.meta.dirname, '..', '..');
const aliases = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'skill-aliases.json'), 'utf8')).aliases;
const pairs = Object.entries(aliases).sort((a, b) => b[0].length - a[0].length);

for (const [from, to] of pairs) {
  const src = path.join(REPO, 'skills', from);
  const dst = path.join(REPO, 'skills', to);
  if (!fs.existsSync(src)) {
    console.log(`skip missing ${from}`);
    continue;
  }
  if (fs.existsSync(dst)) throw new Error(`destination exists: ${to}`);
  execSync(`git mv "skills/${from}" "skills/${to}"`, { cwd: REPO, stdio: 'inherit' });
  console.log(`${from} -> ${to}`);
}

const skip = new Set([
  'registries/skill-aliases.json',
  'runtime/tools/build-skills-manifest.mjs',
  'runtime/tools/rename-ow.mjs',
  'CHANGELOG.md',
]);

function walk(dir, out = []) {
  for (const name of fs.readdirSync(dir)) {
    if (name === '.git' || name === 'node_modules') continue;
    const abs = path.join(dir, name);
    const rel = path.relative(REPO, abs).replace(/\\/g, '/');
    if (rel.startsWith('docs/audit/')) continue;
    if (skip.has(rel)) continue;
    const st = fs.statSync(abs);
    if (st.isDirectory()) walk(abs, out);
    else if (/\.(md|go|js|mjs|json|yml|yaml|ps1|sh|txt|html)$/.test(name)) out.push(abs);
  }
  return out;
}

for (const file of walk(REPO)) {
  let text = fs.readFileSync(file, 'utf8');
  let next = text;
  for (const [from, to] of pairs) {
    const re = new RegExp(`(?<![A-Za-z0-9])${from}(?![A-Za-z0-9])`, 'g');
    next = next.replace(re, to);
  }
  if (next !== text) fs.writeFileSync(file, next);
}

console.log('rewrote references');
