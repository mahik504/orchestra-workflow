const fs = require('fs');
const path = require('path');
const os = require('os');

const pluginDir = path.join(os.homedir(), 'AppData', 'Local', 'Claude-3p', 'local-agent-mode-sessions', 'skills-plugin', '00000000-0000-4000-8000-000000000001', '8e2b40ae-18e9-4e7f-b0e9-f7ce6b8a75c8');
const skillsDir = path.join(pluginDir, 'skills');
const manifestPath = path.join(pluginDir, 'manifest.json');

if (!fs.existsSync(manifestPath)) {
  console.log('Claude Desktop manifest not found at ' + manifestPath);
  process.exit(0);
}

const manifest = JSON.parse(fs.readFileSync(manifestPath, 'utf8'));
const existingMap = new Map();
for (const s of manifest.skills) {
  existingMap.set(s.name, s);
}

// Read descriptions from orchestra skills.json
const REPO = path.join(__dirname, '..', '..');
const owSkillsPath = path.join(REPO, 'registries', 'skills.json');
let owSkillsData = [];
if (fs.existsSync(owSkillsPath)) {
  owSkillsData = JSON.parse(fs.readFileSync(owSkillsPath, 'utf8')).skills || [];
}
const descMap = new Map();
for (const s of owSkillsData) {
  descMap.set(s.id, s.triggers || '');
}

const skillDirs = fs.readdirSync(skillsDir).filter(f => fs.statSync(path.join(skillsDir, f)).isDirectory());

let addedCount = 0;
for (const dirName of skillDirs) {
  if (existingMap.has(dirName)) continue;

  let desc = descMap.get(dirName) || '';
  if (!desc) {
    const skillMdPath = path.join(skillsDir, dirName, 'SKILL.md');
    if (fs.existsSync(skillMdPath)) {
      const content = fs.readFileSync(skillMdPath, 'utf8');
      const lines = content.split(/\r?\n/);
      let inFrontmatter = false;
      let descLines = [];
      for (const line of lines) {
        if (line.trim() === '---') {
          if (!inFrontmatter) {
            inFrontmatter = true;
          } else {
            break;
          }
          continue;
        }
        if (inFrontmatter && line.startsWith('description:')) {
          descLines.push(line.replace('description:', '').trim());
        }
      }
      desc = descLines.join(' ').replace(/^['"]|['"]$/g, '');
    }
  }

  const entry = {
    skillId: dirName,
    name: dirName,
    description: desc || ('Orchestra Workflow skill: ' + dirName),
    creatorType: 'orchestra',
    updatedAt: null,
    enabled: true
  };
  manifest.skills.push(entry);
  existingMap.set(dirName, entry);
  addedCount++;
}

manifest.lastUpdated = Date.now();
fs.writeFileSync(manifestPath, JSON.stringify(manifest, null, 2), 'utf8');
console.log(`Successfully indexed ${manifest.skills.length} skills into manifest.json (+${addedCount} newly added)!`);
