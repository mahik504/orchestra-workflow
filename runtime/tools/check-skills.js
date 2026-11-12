// registries/skills.json matches skills/ and host-stack.json.
//   node runtime/tools/check-skills.js
const fs = require('fs');
const path = require('path');

const REPO = path.join(__dirname, '..', '..');

function hostLocalFromGitignore() {
  const text = fs.readFileSync(path.join(REPO, '.gitignore'), 'utf8');
  const names = new Set();
  for (const line of text.split(/\r?\n/)) {
    const match = line.match(/^skills\/([A-Za-z0-9._-]+)\/$/);
    if (match) names.add(match[1]);
  }
  return names;
}

const HOST_LOCAL = hostLocalFromGitignore();
const problems = [];

const catalog = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'skills.json'), 'utf8'));
const aliases = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'skill-aliases.json'), 'utf8'));
const stack = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'host-stack.json'), 'utf8'));
const rows = catalog.skills || [];
const ids = rows.map((row) => row.id);
const dirs = fs.readdirSync(path.join(REPO, 'skills')).filter((name) => {
  if (HOST_LOCAL.has(name)) return false;
  const dir = path.join(REPO, 'skills', name);
  return fs.statSync(dir).isDirectory() && fs.existsSync(path.join(dir, 'SKILL.md'));
}).sort();

if (stack.skills_manifest !== 'registries/skills.json') {
  problems.push('host-stack skills_manifest must be registries/skills.json');
}
const listed = [...stack.skills].sort();
if (listed.join() !== [...ids].sort().join()) {
  problems.push('host-stack skills and skills.json ids differ');
}
if (dirs.join() !== [...ids].sort().join()) {
  problems.push(`skills/ has ${dirs.length} canonical dirs, skills.json has ${ids.length}`);
}

const seen = new Set();
for (const row of rows) {
  if (seen.has(row.id)) problems.push(`duplicate skill ${row.id}`);
  seen.add(row.id);
  if (!row.license) problems.push(`${row.id} missing license`);
  if (row.owner === 'orchestra' && row.license !== 'MIT') {
    problems.push(`${row.id} is orchestra-owned and its license is ${row.license}`);
  }
  if (row.tier === 'vendored' && (!row.license || row.license === 'unspecified')) {
    problems.push(`${row.id} vendored license is unspecified`);
  }
  if (row.id.startsWith('ow-') && row.description_chars > 300) {
    problems.push(`${row.id} description is ${row.description_chars} characters; budget is 300`);
  }
  const skill = path.join(REPO, 'skills', row.id, 'SKILL.md');
  if (!fs.existsSync(skill)) problems.push(`missing ${skill}`);
  else {
    const raw = fs.readFileSync(skill);
    if (raw[0] === 0xef && raw[1] === 0xbb && raw[2] === 0xbf) problems.push(`${row.id} SKILL.md starts with a UTF-8 BOM`);
    const text = raw.toString('utf8');
    if (!text.startsWith('---\n') && !text.startsWith('---\r\n')) problems.push(`${row.id} SKILL.md has no frontmatter`);
    if (row.id.startsWith('ow-') && !/^description:/m.test(text)) problems.push(`${row.id} SKILL.md has no description`);
    const front = text.match(/^---\r?\n([\s\S]*?)\r?\n---/);
    const stub = front && /stub:\s*true/.test(front[1]);
    if (row.id.startsWith('ow-') && !stub && raw.length > 6 * 1024) {
      problems.push(`${row.id} body is ${raw.length} bytes; cap is 6144`);
    }
    if (row.owner === 'orchestra' && !stub) {
      for (const heading of ['Preconditions', 'Procedure', 'Output', 'Verify', 'Skip']) {
        if (!text.includes(`## ${heading}`)) problems.push(`${row.id} missing ## ${heading} (or set stub: true with a reason)`);
      }
    }
  }
  if (row.rename_to && aliases.aliases[row.id] !== row.rename_to) {
    problems.push(`alias for ${row.id} does not match rename_to`);
  }
}

for (const [from, to] of Object.entries(aliases.aliases || {})) {
  const row = rows.find((item) => item.id === from || item.id === to);
  if (!row) problems.push(`alias ${from} -> ${to} matches no skill row`);
  const dest = path.join(REPO, 'skills', row && row.id === to ? to : from, 'SKILL.md');
  if (row && !fs.existsSync(dest) && !fs.existsSync(path.join(REPO, 'skills', to, 'SKILL.md'))) {
    problems.push(`alias ${from} -> ${to} has no skill directory`);
  }
}

console.log(`skills:   ${rows.length}`);
console.log(`aliases:  ${Object.keys(aliases.aliases || {}).length}`);
console.log(`problems: ${problems.length}`);
for (const problem of problems) console.log(`  ${problem}`);
process.exit(problems.length ? 1 : 0);
