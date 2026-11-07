// Writes docs/audit/FILE-INVENTORY.json from git ls-files.
//   node runtime/tools/inventory.mjs
import { execSync } from 'node:child_process';
import fs from 'node:fs';
import path from 'node:path';

const REPO = path.join(import.meta.dirname, '..', '..');
const outDir = path.join(REPO, 'docs', 'audit');

function groupOf(file) {
  if (file.startsWith('skills/')) return 'skills';
  if (file.startsWith('runtime/internal/')) return 'runtime';
  if (file.startsWith('runtime/cmd/')) return 'cli';
  if (file.startsWith('runtime/tools/')) return 'tooling';
  if (file.startsWith('runtime/benchmarks/')) return 'benchmarks';
  if (file.startsWith('runtime/')) return 'runtime';
  if (file.startsWith('registries/')) return 'registries';
  if (file.startsWith('protocols/')) return 'protocols';
  if (file.startsWith('docs/')) return 'docs';
  if (file.startsWith('kit/')) return 'installers';
  if (file.startsWith('templates/')) return 'templates';
  if (file.startsWith('.github/')) return 'ci';
  if (file.startsWith('workspace-template/')) return 'workspace-template';
  return 'root';
}

function statusOf(file) {
  if (file.includes('templates/shelf/')) return 'historical';
  if (file.startsWith('docs/diagrams/0') || file === 'docs/AGENT_PORTABILITY.md') return 'historical';
  if (file.includes('_test.go') || file.includes('.test.')) return 'test';
  return 'active';
}

function purposeOf(file) {
  const g = groupOf(file);
  if (g === 'skills') {
    const id = file.split('/')[1];
    if (file.endsWith('SKILL.md')) return `Skill instructions for ${id}`;
    if (file.includes('LICENSE')) return `License notice for ${id}`;
    return `Supporting file for skill ${id}`;
  }
  if (file.endsWith('_test.go')) return 'Go test';
  if (g === 'registries') return 'Registry or schema';
  if (g === 'protocols') return 'Protocol the conductor follows';
  if (g === 'ci') return 'GitHub Actions workflow';
  if (g === 'installers') return 'Host install or sync script';
  if (g === 'docs') return 'Documentation';
  if (g === 'benchmarks') return 'Benchmark fixture or record';
  if (g === 'cli') return 'orchestra CLI';
  if (g === 'tooling') return 'Node check or helper';
  if (g === 'runtime') return 'Go engine package';
  if (g === 'templates') return 'Template copied into a product repo';
  return 'Repository root contract or metadata';
}

const files = execSync('git ls-files', { cwd: REPO, encoding: 'utf8' })
  .split(/\r?\n/)
  .filter(Boolean);

const inventory = files.map((file) => {
  const abs = path.join(REPO, file);
  let bytes = 0;
  try { bytes = fs.statSync(abs).size; } catch { bytes = 0; }
  const tested = file.endsWith('_test.go') || files.some((other) => other === file.replace(/\.go$/, '_test.go'));
  return {
    path: file,
    group: groupOf(file),
    purpose: purposeOf(file),
    status: statusOf(file),
    bytes,
    tested: file.endsWith('.go') ? tested : groupOf(file) === 'tooling',
    documented: groupOf(file) === 'docs' || file === 'README.md' || file === 'AGENTS.md',
  };
});

const byGroup = {};
for (const row of inventory) {
  byGroup[row.group] = (byGroup[row.group] || 0) + 1;
}

fs.mkdirSync(outDir, { recursive: true });
const payload = {
  generated: new Date().toISOString().slice(0, 10),
  commit: execSync('git rev-parse --short HEAD', { cwd: REPO, encoding: 'utf8' }).trim(),
  tracked_files: inventory.length,
  by_group: byGroup,
  files: inventory,
};
fs.writeFileSync(path.join(outDir, 'FILE-INVENTORY.json'), JSON.stringify(payload, null, 2) + '\n');
console.log(`inventory ${inventory.length} files`);
console.log(JSON.stringify(byGroup));
