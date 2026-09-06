// One-shot leftover scanner. Personal from-strings were removed 2026-09-06.
// Usage: node runtime/tools/genericize.js
const fs = require('fs');
const path = require('path');

const REG = path.join(__dirname, '..', '..', 'registries', 'resources.json');
const raw = fs.readFileSync(REG, 'utf8');
JSON.parse(raw);
const leftover = [...raw.matchAll(/\bLYRA\b|\bMAHI\b/g)].length;
console.log(`leftover operator tokens in resources.json: ${leftover}`);
process.exit(leftover ? 1 : 0);
