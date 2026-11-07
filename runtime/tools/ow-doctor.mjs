// Host doctor. Prints names, never secret values.
//
//   node runtime/tools/ow-doctor.mjs --hosts
//   node runtime/tools/ow-doctor.mjs --private
//
// OW_SYNC_HOME overrides home (tests). Vault checks run only with --private.
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { hashTree, hostDirs } from './ow-sync.mjs';

const REPO = path.join(import.meta.dirname, '..', '..');

export function inspectHost(expected, destRoot) {
  const have = fs.existsSync(destRoot)
    ? fs.readdirSync(destRoot).filter((name) => fs.statSync(path.join(destRoot, name)).isDirectory())
    : [];
  const missing = expected.filter((id) => !have.includes(id));
  const extra = have.filter((id) => !expected.includes(id));
  let manifest = { skills: {} };
  const manifestPath = path.join(destRoot, '.ow-manifest.json');
  if (fs.existsSync(manifestPath)) manifest = JSON.parse(fs.readFileSync(manifestPath, 'utf8'));
  const mismatch = [];
  for (const id of expected) {
    const recorded = manifest.skills && manifest.skills[id] && manifest.skills[id].sha256;
    const dir = path.join(destRoot, id);
    if (!recorded || !fs.existsSync(dir)) continue;
    if (hashTree(dir) !== recorded) mismatch.push(id);
  }
  return { missing, extra, mismatch, present: have.length };
}

function credentialReport(file) {
  if (!fs.existsSync(file)) return [{ file, state: 'absent' }];
  let json;
  try {
    json = JSON.parse(fs.readFileSync(file, 'utf8'));
  } catch {
    return [{ file, state: 'unreadable' }];
  }
  const servers = json.mcpServers || json.servers || {};
  const rows = [];
  for (const name of Object.keys(servers)) {
    const env = servers[name].env || {};
    const keys = Object.keys(env);
    if (keys.length === 0) rows.push({ server: name, state: 'no-env' });
    for (const key of keys) {
      const value = env[key];
      rows.push({ server: name, key, state: value ? 'set' : 'empty' });
    }
  }
  return rows;
}

function descriptionBudget(skillsDir) {
  const problems = [];
  let total = 0;
  for (const name of fs.readdirSync(skillsDir)) {
    const file = path.join(skillsDir, name, 'SKILL.md');
    if (!fs.existsSync(file)) continue;
    const text = fs.readFileSync(file, 'utf8');
    const match = text.match(/^description:\s*(.+)$/m);
    const chars = match ? match[1].trim().length : 0;
    total += chars;
    if (name.startsWith('ow-') && chars > 300) problems.push(`${name} description is ${chars} characters`);
  }
  return { total, problems };
}

function privateChecks(vault) {
  const problems = [];
  if (!vault) {
    problems.push('ORCHESTRA_VAULT unset');
    return problems;
  }
  if (!fs.existsSync(vault)) problems.push('vault path missing');
  const routes = path.join(vault, 'routes.md');
  if (fs.existsSync(routes)) {
    const text = fs.readFileSync(routes, 'utf8');
    if (text.includes('docs/3.3.2')) problems.push('routes.md still points at docs/3.3.2');
    if (/every 16 hours/i.test(text)) problems.push('routes.md still says a 16 hour heartbeat');
  }
  return problems;
}

function run(argv) {
  const hosts = argv.includes('--hosts') || argv.includes('--private') === false;
  const priv = argv.includes('--private');
  const home = process.env.OW_SYNC_HOME || os.homedir();
  const stack = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'host-stack.json'), 'utf8'));
  const problems = [];
  if (argv.includes('--hosts') || hosts) {
    for (const [host, dest] of Object.entries(hostDirs(home))) {
      if (!fs.existsSync(dest)) {
        console.log(`${host} absent`);
        continue;
      }
      const report = inspectHost(stack.skills, dest);
      console.log(`${host} present=${report.present} missing=${report.missing.length} extra=${report.extra.length} hash-mismatch=${report.mismatch.length}`);
      if (report.missing.length) problems.push(`${host} missing ${report.missing.length}`);
      if (report.mismatch.length) problems.push(`${host} hash mismatch ${report.mismatch.join(', ')}`);
    }
  }
  const budget = descriptionBudget(path.join(REPO, 'skills'));
  console.log(`description chars: ${budget.total}`);
  problems.push(...budget.problems);
  const mcpFile = path.join(home, '.cursor', 'mcp.json');
  for (const row of credentialReport(mcpFile)) {
    if (row.server) console.log(`credential ${row.server} ${row.key || ''} ${row.state}`.replace(/\s+/g, ' ').trim());
    else console.log(`mcp file ${row.state}`);
  }
  if (priv) {
    for (const problem of privateChecks(process.env.ORCHESTRA_VAULT || '')) problems.push(problem);
  } else {
    console.log('private checks skipped (pass --private)');
  }
  console.log(`problems: ${problems.length}`);
  for (const problem of problems) console.log(`FAIL ${problem}`);
  return problems.length ? 1 : 0;
}

if (process.argv[1] && path.resolve(process.argv[1]).endsWith('ow-doctor.mjs')) {
  process.exit(run(process.argv.slice(2)));
}
