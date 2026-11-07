// Checks a design-contract directory against the pack for this task.
// A bug pack requires nothing. A visual pack requires the headings in the templates.
//
//   node runtime/tools/check-design-contract.js <dir> <profile>
const fs = require('fs');
const path = require('path');

const PROFILES = {
  none: [],
  premium: ['DESIGN.md', 'SPEC.md', 'MOTION.md', 'QA_PLAN.md', 'DECISIONS.md'],
  spatial: ['DESIGN.md', 'SPEC.md', 'MOTION.md', '3D_SCENE.md', 'QA_PLAN.md', 'DECISIONS.md'],
  backend: ['ARCHITECTURE.md', 'API.md'],
  saas: ['DESIGN.md', 'SPEC.md', 'ARCHITECTURE.md', 'API.md', 'QA_PLAN.md', 'DECISIONS.md'],
};

const HEADINGS = {
  'DESIGN.md': [
    'Brand identity',
    'Visual direction',
    'Reference analysis',
    'Color tokens',
    'Typography system',
    'Spacing and layout',
    'Responsive breakpoints',
    'Component system',
    'Surface treatments',
    'Interaction design',
    'Motion language',
    'Scroll behaviour',
    'Cursor behaviour',
    '3D scene requirements',
    'Shader specifications',
    'Accessibility',
    'Performance budgets',
    'Acceptance criteria',
    'Visual references',
  ],
  'SPEC.md': ['Behavior', 'Out of scope', 'Acceptance'],
  'ARCHITECTURE.md': ['Data model', 'Boundaries', 'Migrations', 'Observability'],
  'MOTION.md': ['Engine', 'Scroll', 'Reduced motion', 'What must not move'],
  '3D_SCENE.md': ['Purpose', 'Camera', 'Lighting', 'Materials', 'Loading', 'Disposal', 'DPR', 'Mobile fallback', 'Reduced motion', 'WebGL failure'],
  'API.md': ['Resources', 'Auth', 'Errors', 'Validation'],
  'QA_PLAN.md': ['Viewports', 'Interactions', 'Console', 'Network', 'Accessibility', 'Reduced motion', 'Visual baseline', 'Thresholds', 'WebGL'],
};

function sectionBody(text, heading) {
  const lines = text.replace(/\r\n/g, '\n').split('\n');
  const marker = '## ' + heading;
  const start = lines.findIndex((line) => line.trim() === marker);
  if (start < 0) return null;
  const body = [];
  for (let i = start + 1; i < lines.length; i++) {
    if (lines[i].startsWith('## ')) break;
    body.push(lines[i]);
  }
  return body.join('\n').trim();
}

function checkDir(dir, profile) {
  const required = PROFILES[profile];
  if (!required) return [`unknown profile ${profile}`];
  const problems = [];
  for (const name of required) {
    const file = path.join(dir, name);
    if (!fs.existsSync(file)) {
      problems.push(`missing ${name}`);
      continue;
    }
    const text = fs.readFileSync(file, 'utf8');
    if (name === 'DECISIONS.md') {
      const lines = text.replace(/\r\n/g, '\n').split('\n').map((line) => line.trim()).filter((line) => line.startsWith('- '));
      if (!lines.length) problems.push('DECISIONS.md has no decision line');
      else if (!lines.some((line) => line.includes('ADR:'))) problems.push('DECISIONS.md decision line does not link an ADR');
      continue;
    }
    for (const heading of HEADINGS[name] || []) {
      const body = sectionBody(text, heading);
      if (body == null) problems.push(`${name} missing heading ${heading}`);
      else if (!body) problems.push(`${name} heading ${heading} is empty`);
    }
  }
  return problems;
}

function main() {
  const dir = process.argv[2];
  const profile = process.argv[3];
  if (!dir || !profile) {
    console.error('usage: node runtime/tools/check-design-contract.js <dir> <profile>');
    process.exit(2);
  }
  const problems = checkDir(dir, profile);
  console.log(`profile:  ${profile}`);
  console.log(`problems: ${problems.length}`);
  for (const problem of problems) console.log('  ' + problem);
  process.exit(problems.length ? 1 : 0);
}

if (require.main === module) main();
else module.exports = { checkDir, PROFILES };
