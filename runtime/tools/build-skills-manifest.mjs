// Writes registries/skills.json from skills/ and the rename map.
//   node runtime/tools/build-skills-manifest.mjs
import fs from 'node:fs';
import path from 'node:path';

const REPO = path.join(import.meta.dirname, '..', '..');
const HOST_LOCAL = new Set(['docx', 'pdf', 'pptx', 'xlsx', 'fullstack-mastery', 'creative-3d-web']);
const HOSTS = ['cursor', 'claude', 'antigravity', 'codex', 'jcode', 'agents'];

const RENAME = {
  'orchestra-conductor': 'ow-conductor',
  'prompt-forge': 'ow-forge',
  'orchestra-focus': 'ow-focus',
  'resource-scout': 'ow-resource-scout',
  'visual-forensics': 'ow-reference-forensics',
  'visual-golden-still': 'ow-golden-still',
  'fresh-context-review': 'ow-fresh-review',
  'final-showcase': 'ow-showcase',
  'backend-architecture': 'ow-backend-architecture',
  'r3f-shaders': 'ow-3d-engineering',
  'diagram-generator': 'ow-diagrams',
  'research-paper': 'ow-research-paper',
  'ml-experiment': 'ow-ml-experiment',
  'ship-safe': 'ow-ship-safe',
  'store-release': 'ow-store-release',
  'skill-security-check': 'ow-skill-security-check',
  'semgrep-adapter': 'ow-semgrep',
  'orchestra-career': 'ow-career',
  'orchestra-docs': 'ow-docs',
  'orchestra-marketing': 'ow-marketing',
  'orchestra-observability': 'ow-observability',
  'orchestra-presentation': 'ow-presentation',
  'orchestra-research': 'ow-research',
  'orchestra-ship': 'ow-ship',
  'orchestra-typescript': 'ow-typescript',
  'orchestra-vault': 'ow-vault',
  'orchestra-video': 'ow-video',
};

function descriptionOf(file) {
  const text = fs.readFileSync(file, 'utf8');
  const match = text.match(/^description:\s*(.+)$/m);
  return match ? match[1].trim() : '';
}

function licenseOf(dir) {
  const names = fs.readdirSync(dir).filter((name) => name.startsWith('LICENSE'));
  if (names.length === 0) return { license: 'MIT', license_file: '' };
  const file = names.sort()[0];
  const text = fs.readFileSync(path.join(dir, file), 'utf8').slice(0, 400);
  let license = 'SEE-LICENSE-FILE';
  if (/MIT License/i.test(text)) license = 'MIT';
  else if (/Apache/i.test(text)) license = 'Apache-2.0';
  return { license, license_file: file };
}

const skillsDir = path.join(REPO, 'skills');
const ids = fs.readdirSync(skillsDir).filter((name) => {
  if (HOST_LOCAL.has(name)) return false;
  const dir = path.join(skillsDir, name);
  return fs.statSync(dir).isDirectory() && fs.existsSync(path.join(dir, 'SKILL.md'));
}).sort();

const skills = ids.map((id) => {
  const dir = path.join(skillsDir, id);
  const skillFile = path.join(dir, 'SKILL.md');
  const description = descriptionOf(skillFile);
  const found = licenseOf(dir);
  const owner = RENAME[id] || Object.values(RENAME).includes(id) || id.startsWith('ow-') ? 'orchestra' : 'upstream';
  const license = found.license_file ? found.license : (owner === 'orchestra' ? 'MIT' : 'unspecified');
  return {
    id,
    rename_to: RENAME[id] && !Object.values(RENAME).includes(id) ? RENAME[id] : '',
    owner,
    tier: owner === 'orchestra' ? 'core' : 'vendored',
    source: owner === 'orchestra' ? 'original' : 'vendored',
    license,
    license_file: found.license_file,
    hosts: HOSTS,
    description_chars: description.length,
    triggers: description.slice(0, 300),
    context_chars: fs.statSync(skillFile).size,
  };
});

const aliases = { ...RENAME };

const stackPath = path.join(REPO, 'registries', 'host-stack.json');
const stack = JSON.parse(fs.readFileSync(stackPath, 'utf8'));
stack.skills = skills.map((row) => row.id);
fs.writeFileSync(stackPath, JSON.stringify(stack, null, 2) + '\n');
fs.writeFileSync(path.join(REPO, 'registries', 'skills.json'), JSON.stringify({ skills }, null, 2) + '\n');
fs.writeFileSync(path.join(REPO, 'registries', 'skill-aliases.json'), JSON.stringify({
  note: 'Old id to current id. Upstream skills are absent because their names do not change.',
  aliases,
}, null, 2) + '\n');
console.log(`skills ${skills.length}`);
console.log(`renames ${Object.keys(aliases).length}`);
