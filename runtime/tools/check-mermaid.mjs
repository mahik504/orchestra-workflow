// Parse every mermaid fence in README.md and docs/**/*.md inside Chromium.
//   node runtime/tools/check-mermaid.mjs
import { chromium } from 'playwright';
import fs from 'node:fs';
import path from 'node:path';

const REPO = path.join(import.meta.dirname, '..', '..');
const mermaidPath = path.join(REPO, 'runtime', 'node_modules', 'mermaid', 'dist', 'mermaid.min.js');

function walk(dir, out) {
  if (!fs.existsSync(dir)) return;
  for (const name of fs.readdirSync(dir)) {
    const full = path.join(dir, name);
    if (fs.statSync(full).isDirectory()) walk(full, out);
    else if (name.endsWith('.md')) out.push(full);
  }
}

function fences(text) {
  const found = [];
  const re = /```mermaid\r?\n([\s\S]*?)```/g;
  let match;
  while ((match = re.exec(text))) found.push(match[1]);
  return found;
}

const files = [path.join(REPO, 'README.md')];
walk(path.join(REPO, 'docs'), files);
const blocks = [];
for (const file of files) {
  const text = fs.readFileSync(file, 'utf8');
  for (const source of fences(text)) {
    blocks.push({ file: path.relative(REPO, file), source });
  }
}

if (!fs.existsSync(mermaidPath)) {
  console.error('mermaid bundle missing. Run npm ci in runtime/.');
  process.exit(1);
}
if (blocks.length === 0) {
  console.error('no mermaid fences found');
  process.exit(1);
}

const browser = await chromium.launch();
const page = await browser.newPage();
await page.addScriptTag({ path: mermaidPath });
const problems = [];
for (const block of blocks) {
  const message = await page.evaluate(async (source) => {
    try {
      mermaid.initialize({ startOnLoad: false, securityLevel: 'strict' });
      await mermaid.parse(source);
      return '';
    } catch (err) {
      return err && err.message ? err.message : String(err);
    }
  }, block.source);
  if (message) problems.push(`${block.file}: ${message.split('\n')[0]}`);
}
await browser.close();

console.log(`diagrams: ${blocks.length}`);
console.log(`problems: ${problems.length}`);
for (const problem of problems) console.log(`  ${problem}`);
process.exit(problems.length ? 1 : 0);
