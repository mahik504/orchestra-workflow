// Public-repo hygiene gate. Runs in CI and locally with the same rules.
//
//   node runtime/tools/hygiene.js
//
// Fails on: the operator's real paths, private-repo references, personal product
// catalog ids, credential-shaped strings, and a VERSION that is not the expected
// contract version.
//
// Only tracked files are scanned, so a local build cache cannot trip or hide a rule.
const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

const REPO = path.join(__dirname, '..', '..');
const EXPECTED_VERSION = JSON.parse(
  fs.readFileSync(path.join(REPO, 'registries', 'versions.json'), 'utf8')
).current;

// These files define or document the rules, so they contain the patterns by design.
const SELF = new Set([
  'runtime/tools/hygiene.js',
  'runtime/tools/check-registry.js',
  'runtime/tools/check-host-stack.js',
  'runtime/tools/strip-personal.js',
  'runtime/tools/genericize.js',
  '.github/workflows/hygiene.yml',
]);

const RULES = [
  // --- operator identity and machine ---
  { name: 'operator home directory', re: /C:\\Users\\Asus|C:\/Users\/Asus/i },
  { name: 'private vault path', re: /C:\\projects\\orchestra-brain|C:\/projects\/orchestra-brain/i },
  { name: 'private vault repo', re: /orchestra-brain/i },
  { name: 'operator email', re: /gehlot\.mahisingh/i },
  // The public repo's own clone URL is allowed; every other use of the handle is not.
  { name: 'operator github handle', re: /mahik504(?!\/(?:orchestra-workflow|limitless-claude))/i },

  // --- personal product catalog ---
  { name: 'personal product id', re: /\b(AstroVerse|AirLens|YUMIT|ODYSS|penfight|EvoNex|XplainAI|AstroLens|EvoMoE|TTB Agro|SIH26174)\b/i },
  { name: 'personal operator HUD id', re: /\bLYRA\b/ },

  // --- credentials ---
  { name: 'openai-style key', re: /\bsk-[A-Za-z0-9]{20,}/ },
  { name: 'github classic token', re: /\bghp_[A-Za-z0-9]{20,}/ },
  { name: 'github fine-grained token', re: /\bgithub_pat_[A-Za-z0-9_]{20,}/ },
  { name: 'google api key', re: /\bAIza[0-9A-Za-z_-]{30,}/ },
  { name: 'xai key', re: /\bxai-[A-Za-z0-9]{20,}/ },
  { name: 'bearer token', re: /bearer\s+[A-Za-z0-9._-]{24,}/i },
  { name: 'anthropic document license', re: /materials outside the Services/ },
];

const LOCAL_SKILL_PREFIXES = [
  'skills/docx/',
  'skills/pdf/',
  'skills/pptx/',
  'skills/xlsx/',
  'skills/fullstack-mastery/',
  'skills/creative-3d-web/',
];

const BINARY = new Set(['.png', '.jpg', '.jpeg', '.gif', '.webp', '.ico', '.pdf', '.zip', '.exe', '.woff', '.woff2', '.ttf']);

let entries;
try {
  entries = execSync('git ls-files -s', { cwd: REPO, encoding: 'utf8' }).split('\n').filter(Boolean);
} catch (e) {
  console.error('not a git repository, or git is unavailable');
  process.exit(2);
}

const failures = [];
const files = [];
for (const line of entries) {
  const tab = line.indexOf('\t');
  if (tab < 0) continue;
  const rel = line.slice(tab + 1).replace(/\\/g, '/');
  const mode = line.slice(0, tab).split(' ')[0];
  if (mode === '160000') {
    failures.push(`${rel}  gitlink (mode 160000): commit the files, not a nested repository`);
    continue;
  }
  files.push(rel);
}

for (const rel of files) {
  for (const prefix of LOCAL_SKILL_PREFIXES) {
    if (rel.startsWith(prefix)) failures.push(`${rel} is a host-local skill and must stay gitignored`);
  }
  if (SELF.has(rel)) continue;
  if (BINARY.has(path.extname(rel).toLowerCase())) continue;
  const abs = path.join(REPO, rel);
  if (!fs.existsSync(abs)) continue;
  const stat = fs.statSync(abs);
  if (stat.isDirectory()) {
    failures.push(`${rel}  tracked path is a directory`);
    continue;
  }
  let text;
  try {
    text = fs.readFileSync(abs, 'utf8');
  } catch (e) {
    failures.push(`${rel}  unreadable: ${e.code || e.message}`);
    continue;
  }
  if (text.includes('\u0000')) {
    failures.push(`${rel}  contains NUL bytes (UTF-16 or a binary mislabeled as text)`);
    continue;
  }
  const lines = text.split(/\r?\n/);
  lines.forEach((line, i) => {
    for (const rule of RULES) {
      const m = line.match(rule.re);
      if (m) failures.push(`${rel}:${i + 1}  ${rule.name}: ${m[0].slice(0, 60)}`);
    }
  });
}

// VERSION must match the contract version.
const versionPath = path.join(REPO, 'VERSION');
const version = fs.existsSync(versionPath) ? fs.readFileSync(versionPath, 'utf8').trim() : '(missing)';
if (version !== EXPECTED_VERSION) {
  failures.push(`VERSION  is "${version}", expected "${EXPECTED_VERSION}"`);
}

// The contract rule must be present in every contract file.
const CONTRACT_FILES = ['AGENTS.md', '.cursorrules', 'CLAUDE.md', 'ARCHITECTURE.md', 'WORKFLOW.md', 'skills/ow-conductor/SKILL.md'];
for (const rel of CONTRACT_FILES) {
  const abs = path.join(REPO, rel);
  if (!fs.existsSync(abs)) { failures.push(`${rel}  missing contract file`); continue; }
  if (!/ORCHESTRA = CONTROL PLANE/.test(fs.readFileSync(abs, 'utf8'))) {
    failures.push(`${rel}  missing the immutable control-plane rule`);
  }
}

console.log(`scanned files: ${files.length}`);
console.log(`VERSION:       ${version}`);
console.log(`failures:      ${failures.length}`);
for (const f of failures) console.log(`  ${f}`);
if (failures.length) process.exit(1);

try {
  execSync('node runtime/tools/check-host-stack.js', { cwd: REPO, stdio: 'inherit' });
  execSync('node runtime/tools/check-routing.js', { cwd: REPO, stdio: 'inherit' });
  execSync('node runtime/tools/check-route-agreement.js', { cwd: REPO, stdio: 'inherit' });
  execSync('node runtime/tools/check-41.js', { cwd: REPO, stdio: 'inherit' });
} catch (e) {
  process.exit(1);
}
process.exit(0);
