// Fail on secret-shaped strings, curl|sh, unpinned @latest, and a UTF-8 BOM
// in the files Orchestra ships as instructions.
import fs from 'node:fs';
import path from 'node:path';

const repo = path.join(import.meta.dirname, '..', '..');
const files = [
  'AGENTS.md',
  'CLAUDE.md',
  '.cursorrules',
  'kit/mcp-core.json',
  'mcp_config.example.json',
  'kit/antigravity/mcp_config.example.json',
  'templates/mcp/optional-servers.json',
  '.cursor/hooks.json',
  '.claude/settings.json',
];
for (const dir of ['.cursor/rules', '.claude/rules']) {
  const abs = path.join(repo, dir);
  if (!fs.existsSync(abs)) continue;
  for (const name of fs.readdirSync(abs)) {
    if (name.endsWith('.md') || name.endsWith('.mdc')) files.push(path.join(dir, name));
  }
}

const problems = [];
for (const rel of files) {
  const abs = path.join(repo, rel);
  if (!fs.existsSync(abs)) continue;
  const buf = fs.readFileSync(abs);
  if (buf.length >= 3 && buf[0] === 0xef && buf[1] === 0xbb && buf[2] === 0xbf) {
    problems.push(`${rel} has a UTF-8 BOM`);
  }
  const text = buf.toString('utf8');
  if (/sk-[A-Za-z0-9]{20,}/.test(text) || /ghp_[A-Za-z0-9]{20,}/.test(text)) {
    problems.push(`${rel} has a secret-shaped string`);
  }
  if (/curl\b[^|\n]*\|\s*(sh|bash)\b/i.test(text) && !rel.endsWith('check-config-lint.mjs')) {
    problems.push(`${rel} contains curl piped to a shell`);
  }
  if ((rel.endsWith('mcp-core.json') || rel.endsWith('mcp_config.example.json')) && /@latest\b/.test(text)) {
    problems.push(`${rel} uses @latest`);
  }
}

if (problems.length) {
  for (const row of problems) console.error(row);
  process.exit(1);
}
console.log('config lint: 0 problems');
