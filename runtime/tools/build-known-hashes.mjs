// Hash each skill tree at the named revisions.
//   node runtime/tools/build-known-hashes.mjs
import { execFileSync } from 'node:child_process';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { hashTree } from './ow-sync.mjs';

const REPO = path.join(import.meta.dirname, '..', '..');
const wanted = ['v4.0.0', 'v4.1.0', 'v4.2.0', 'v4.5.0', '60cf9e7'];
const present = [];
for (const rev of wanted) {
  try {
    execFileSync('git', ['rev-parse', '--verify', rev], { cwd: REPO, stdio: 'ignore' });
    present.push(rev);
  } catch {
    console.log(`skip ${rev}`);
  }
}

const hashes = {};
const tmp = fs.mkdtempSync(path.join(os.tmpdir(), 'ow-known-'));
for (const rev of present) {
  const extract = path.join(tmp, rev.replace(/[^\w.-]/g, '_'));
  fs.mkdirSync(extract, { recursive: true });
  const tarPath = path.join(extract, 'skills.tar');
  let tar;
  try {
    tar = execFileSync('git', ['archive', rev, 'skills'], { cwd: REPO, maxBuffer: 256 * 1024 * 1024 });
  } catch (err) {
    console.log(`no skills tree at ${rev}`);
    continue;
  }
  fs.writeFileSync(tarPath, tar);
  execFileSync('tar', ['-xf', tarPath, '-C', extract]);
  const skills = path.join(extract, 'skills');
  if (!fs.existsSync(skills)) continue;
  for (const name of fs.readdirSync(skills)) {
    const dir = path.join(skills, name);
    if (!fs.statSync(dir).isDirectory()) continue;
    if (!fs.existsSync(path.join(dir, 'SKILL.md'))) continue;
    const digest = hashTree(dir);
    hashes[name] = hashes[name] || [];
    if (!hashes[name].includes(digest)) hashes[name].push(digest);
  }
  console.log(`hashed ${rev}`);
}

const out = { generated_from: present, hashes };
fs.writeFileSync(path.join(REPO, 'registries', 'known-skill-hashes.json'), JSON.stringify(out, null, 2) + '\n');
console.log(`ids ${Object.keys(hashes).length}`);
