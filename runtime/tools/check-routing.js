// The routing table is copied by hand into the host files.
// This check fails when a Load cell names a skill that is not in skills/,
// when Skip text is missing, or when the four copies drift.
// It does not ask a model anything. Document checks are not model trials.
const fs = require('fs');
const path = require('path');

const REPO = path.join(__dirname, '..', '..');
const FILES = [
  'AGENTS.md',
  'CLAUDE.md',
  'kit/antigravity/GEMINI.md',
  'skills/orchestra-conductor/SKILL.md',
  '.cursorrules',
];
const START = '<!-- routing:start -->';
const END = '<!-- routing:end -->';
const RESERVED = new Set(['not ready', 'host-provided', 'none']);

// Exact job cell, exact Load tokens, required Skip phrases.
const EXPECTED = [
  { job: 'Award-style site', load: ['impeccable', 'taste-design', 'visual-golden-still', 'gsap-core', 'gsap-scrolltrigger', 'gsap-timeline'], skip: ['copying a reference site\'s files', 'a second motion stack'] },
  { job: 'Site to echo', load: ['visual-forensics', 'impeccable', 'taste-design'], skip: ['copying assets, copy, or a logo'] },
  { job: 'Scroll', load: ['gsap-core', 'gsap-scrolltrigger', 'gsap-timeline'], skip: ['gsap-react unless the motion lives in React components'] },
  { job: 'One 3D scene', load: ['r3f-shaders'], skip: ['3D as decoration', 'a second motion engine'] },
  { job: 'SaaS or API', load: ['backend-architecture', 'test-driven-development'], skip: ['the design lab'] },
  { job: 'Bug', load: ['systematic-debugging'], skip: ['the design lab'] },
  { job: 'Paper', load: ['research-paper', 'humanizer'], skip: ['invented citations'] },
  { job: 'Deck', load: ['orchestra-docs'], skip: ['host-provided pptx unless that host already has it'] },
  { job: 'Deck or slides', load: ['orchestra-docs', 'orchestra-presentation'], skip: ['copying a pptx skill file', 'purple template dumps'] },
  { job: 'Architecture diagram', load: ['diagram-generator', 'pretty-mermaid'], skip: ['a second diagram system', 'a README full of charts'] },
  { job: 'One control', load: ['none'], skip: ['loading every component kit', 'duplicating a control the project already has'] },
  { job: 'Next.js speed', load: ['react-best-practices'], skip: ['a redesign', 'a bulk Vercel skill install'] },
  { job: 'TypeScript review', load: ['orchestra-typescript'], skip: ['a visual redesign', 'a new framework'] },
  { job: 'Ship a deployed app', load: ['ship-safe', 'orchestra-observability'], skip: ['scanning a third-party URL', 'a global Sentry install'] },
  { job: 'What people said recently', load: ['orchestra-research'], skip: ['a last-30-days pack', 'scraping', 'posting', 'DMs'] },
  { job: 'Marketing copy', load: ['orchestra-marketing', 'humanizer'], skip: ['invented metrics', 'a design system rewrite'] },
  { job: 'Rendered video', load: ['orchestra-video'], skip: ['the GSAP route', 'installing a video suite into this method repo'] },
  { job: 'Career notes', load: ['orchestra-career'], skip: ['job-board scraping', 'applying on your behalf'] },
  { job: 'README and post', load: ['orchestra-ship', 'humanizer'], skip: ['posting for the human'] },
  { job: 'Binary or APK', load: ['not ready'], skip: ['rea until Java is installed'] },
  { job: 'Hackathon or college', load: ['prompt-forge'], skip: ['the 23-card survey when a reference is named'] },
  { job: 'Docs on disk', load: ['host-provided'], skip: ['docx, pptx, pdf, and xlsx stay on the host'] },
];

function section(file) {
  const text = fs.readFileSync(path.join(REPO, file), 'utf8').replace(/\r\n/g, '\n');
  const starts = text.split(START).length - 1;
  const ends = text.split(END).length - 1;
  if (starts !== 1 || ends !== 1) return { block: null, starts, ends };
  const a = text.indexOf(START);
  const b = text.indexOf(END);
  if (b < a) return { block: null, starts, ends };
  return { block: text.slice(a, b + END.length), starts, ends };
}

function rowsOf(block) {
  const rows = [];
  for (const line of block.split('\n')) {
    if (!line.startsWith('|') || line.includes('---') || line.includes('Job |')) continue;
    const cells = line.split('|').slice(1, -1).map((cell) => cell.trim());
    if (cells.length < 3) continue;
    rows.push({
      job: cells[0],
      load: cells[1].split(',').map((item) => item.trim()).filter(Boolean),
      skip: cells[2],
    });
  }
  return rows;
}

const problems = [];
const blocks = FILES.map((file) => {
  const found = section(file);
  if (found.starts !== 1 || found.ends !== 1 || !found.block) {
    problems.push(`${file} needs one routing start and one end (starts ${found.starts}, ends ${found.ends})`);
  }
  return found.block;
});

if (blocks[0]) {
  for (let i = 1; i < blocks.length; i++) {
    if (blocks[i] && blocks[i] !== blocks[0]) problems.push(`${FILES[i]} routing block differs from AGENTS.md`);
  }

  const rows = rowsOf(blocks[0]);
  const jobs = rows.map((row) => row.job);
  const seen = new Set();
  for (const job of jobs) {
    const key = job.toLowerCase();
    if (seen.has(key)) problems.push(`duplicate job "${job}"`);
    seen.add(key);
  }
  let skillPaths = 0;
  for (const row of rows) {
    for (const id of row.load) {
      if (RESERVED.has(id)) continue;
      if (!/^[a-z0-9-]+$/.test(id)) {
        problems.push(`bad Load token "${id}"`);
        continue;
      }
      skillPaths += 1;
      const skill = path.join(REPO, 'skills', id, 'SKILL.md');
      if (!fs.existsSync(skill)) problems.push(`Load names missing skill ${id}`);
    }
    for (const id of row.load) {
      if (!RESERVED.has(id) && !/^[a-z0-9-]+$/.test(id)) continue;
      if (![...RESERVED].includes(id) && !/^[a-z0-9-]+$/.test(id)) {
        problems.push(`Load token "${id}" is neither reserved nor a skill id`);
      }
    }
  }

  let assertions = 0;
  for (const expected of EXPECTED) {
    assertions += 1;
    const row = rows.find((item) => item.job.toLowerCase() === expected.job.toLowerCase());
    if (!row) {
      problems.push(`missing row "${expected.job}"`);
      continue;
    }
    if (row.load.join(', ') !== expected.load.join(', ')) {
      problems.push(`"${expected.job}" load is "${row.load.join(', ')}", expected "${expected.load.join(', ')}"`);
    }
    for (const phrase of expected.skip) {
      assertions += 1;
      if (!row.skip.toLowerCase().includes(phrase.toLowerCase())) {
        problems.push(`"${expected.job}" skip missing "${phrase}"`);
      }
    }
  }

  const loose = rows.filter((row) => row.job.toLowerCase() === 'app');
  if (loose.length) problems.push('a bare job named "app" would steal "Ship a deployed app"');
  const ship = rows.find((row) => row.job === 'Ship a deployed app');
  if (ship && rows.some((row) => row.job.toLowerCase() === 'app')) {
    problems.push('substring route "app" matches a real job');
  }

  const conductor = fs.readFileSync(path.join(REPO, 'skills/orchestra-conductor/SKILL.md'), 'utf8');
  const outside = conductor.replace(/<!-- routing:start -->[\s\S]*<!-- routing:end -->/, '');
  if (!outside.includes('interaction-components')) {
    problems.push('conductor body outside the routing block does not name interaction-components');
  }
  if (!/One control loads `none`/.test(outside)) {
    problems.push('conductor body does not say One control loads none');
  }

  console.log('routing copies:    ' + FILES.length);
  console.log('required routes:   ' + EXPECTED.length);
  console.log('skill-path checks: ' + skillPaths);
  console.log('load/skip asserts: ' + assertions);
  console.log('document check:    not a model trial');
  console.log('problems:          ' + problems.length);
  for (const problem of problems) console.log('  ' + problem);
  if (problems.length) process.exit(1);
} else {
  console.log('problems:          ' + problems.length);
  for (const problem of problems) console.log('  ' + problem);
  process.exit(1);
}
