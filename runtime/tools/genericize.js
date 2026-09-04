// One-shot: replace operator-specific product references in the public registry
// with the generic capability they actually describe.
// Usage: node runtime/tools/genericize.js [--write]
const fs = require('fs');
const path = require('path');

const REG = path.join(__dirname, '..', '..', 'registries', 'resources.json');

// Longest first so the specific phrasing wins over the general one.
const SUBS = [
  ['MAHI may self-host (localhost:20128/v1) and point LYRA .env.local at it. Skip inside Cursor/Antigravity.',
    'Can be self-hosted (for example localhost:20128/v1) and pointed at from an application .env.local. Do not install it inside the IDE.'],
  ['Free tier aggregator. LYRA has Studio/Groq. Fallback only if both die.',
    'Free-tier aggregator. Fallback only when the primary providers are unavailable.'],
  ['Not LYRA crimson taste profile.', 'Not a bespoke operator-HUD taste profile.'],
  ['LYRA operator HUD', 'an operator HUD'],
  ['LYRA HUD surfaces', 'operator HUD surfaces'],
  ['LYRA nucleus', 'operator HUD core'],
  ['Not LYRA HUD.', 'Not an operator HUD.'],
  ['LYRA HUD', 'operator HUD'],
  ['MAHI', 'the operator'],
];

let raw = fs.readFileSync(REG, 'utf8');
const before = raw;
const report = [];
for (const [from, to] of SUBS) {
  const n = raw.split(from).length - 1;
  if (n > 0) { raw = raw.split(from).join(to); report.push(`${n}x  "${from.slice(0, 55)}" -> "${to.slice(0, 55)}"`); }
}

JSON.parse(raw); // fail loudly rather than write broken JSON

const leftover = [...raw.matchAll(/\bLYRA\b|\bMAHI\b/g)].length;
console.log(report.join('\n') || 'no substitutions matched');
console.log(`\nchanged: ${before !== raw}`);
console.log(`leftover LYRA/MAHI tokens: ${leftover}`);

if (process.argv.includes('--write')) {
  fs.writeFileSync(REG, raw, 'utf8');
  console.log(`wrote ${REG}`);
} else {
  console.log('dry run. pass --write to apply.');
}
