// Transactional skill sync. Dry-run unless --apply.
// A skill this tool did not install stays, except old Orchestra names, which move to a backup.
//
//   node runtime/tools/ow-sync.mjs
//   node runtime/tools/ow-sync.mjs --apply --host cursor
//   node runtime/tools/ow-sync.mjs --apply --backup-and-replace
//   node runtime/tools/ow-sync.mjs --restore <stamp> --host cursor --apply
//
// OW_SYNC_HOME overrides the home directory (tests).
// OW_BACKUP_ROOT overrides the backup directory.
import { createHash } from 'node:crypto';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';

const REPO = path.join(import.meta.dirname, '..', '..');

export function hostDirs(home) {
  return {
    cursor: path.join(home, '.cursor', 'skills'),
    claude: path.join(home, '.claude', 'skills'),
    antigravity: path.join(home, '.gemini', 'config', 'skills'),
    agents: path.join(home, '.agents', 'skills'),
    jcode: path.join(home, '.jcode', 'skills'),
    codex: path.join(home, '.codex', 'skills'),
  };
}

export function listSkillIds(sourceRoot) {
  if (!fs.existsSync(sourceRoot)) return [];
  return fs.readdirSync(sourceRoot).filter((name) => {
    const dir = path.join(sourceRoot, name);
    const st = fs.lstatSync(dir);
    if (st.isSymbolicLink()) return false;
    return st.isDirectory() && fs.existsSync(path.join(dir, 'SKILL.md'));
  }).sort();
}

function refuseEscape(child, root) {
  const st = fs.lstatSync(child);
  if (st.isSymbolicLink()) {
    const error = new Error(`symlink refused: ${child}`);
    error.code = 'SYMLINK';
    throw error;
  }
  const real = fs.realpathSync(child);
  const realRoot = fs.realpathSync(root);
  const rel = path.relative(realRoot, real);
  if (rel.startsWith('..') || path.isAbsolute(rel)) {
    const error = new Error(`path escapes ${root}: ${child}`);
    error.code = 'ESCAPE';
    throw error;
  }
  return st;
}

function listFiles(dir) {
  const out = [];
  const walk = (abs, rel) => {
    for (const name of fs.readdirSync(abs).sort()) {
      const child = path.join(abs, name);
      const next = rel ? `${rel}/${name}` : name;
      const st = refuseEscape(child, dir);
      if (st.isDirectory()) walk(child, next);
      else out.push(next);
    }
  };
  walk(dir, '');
  return out;
}

function hashTreeWith(dir, transform) {
  const hash = createHash('sha256');
  for (const rel of listFiles(dir)) {
    hash.update(rel);
    hash.update('\0');
    const buf = fs.readFileSync(path.join(dir, ...rel.split('/')));
    hash.update(transform ? transform(buf) : buf);
    hash.update('\0');
  }
  return hash.digest('hex');
}

export function hashTree(dir) {
  return hashTreeWith(dir, hashBytes);
}

export function hashTreeRaw(dir) {
  return hashTreeWith(dir, null);
}

function readManifest(destRoot) {
  const file = path.join(destRoot, '.ow-manifest.json');
  if (!fs.existsSync(file)) return { skills: {} };
  return JSON.parse(fs.readFileSync(file, 'utf8'));
}

function writeManifest(destRoot, manifest) {
  fs.mkdirSync(destRoot, { recursive: true });
  const file = path.join(destRoot, '.ow-manifest.json');
  const tmp = `${file}.tmp`;
  fs.writeFileSync(tmp, JSON.stringify(manifest, null, 2) + '\n');
  fs.renameSync(tmp, file);
}

export function planSync({ sourceIds, sourceRoot, destRoot, backupAndReplace, knownHashes = {} }) {
  const manifest = readManifest(destRoot);
  const actions = [];
  for (const id of sourceIds) {
    const src = path.join(sourceRoot, id);
    const dst = path.join(destRoot, id);
    if (!fs.existsSync(src)) {
      actions.push({ id, action: 'missing-source' });
      continue;
    }
    const srcHash = hashTree(src);
    if (!fs.existsSync(dst)) {
      actions.push({ id, action: 'copy', srcHash });
      continue;
    }
    const dstHash = hashTree(dst);
    const rawHash = hashTreeRaw(dst);
    const recorded = manifest.skills && manifest.skills[id] && manifest.skills[id].sha256;
    const known = new Set(knownHashes[id] || []);
    if (dstHash === srcHash) {
      actions.push({ id, action: 'unchanged', srcHash });
    } else if (recorded === dstHash || recorded === rawHash || known.has(dstHash) || known.has(rawHash) || backupAndReplace) {
      actions.push({ id, action: 'replace', srcHash, dstHash });
    } else {
      actions.push({ id, action: 'conflict', srcHash, dstHash });
    }
  }
  return { manifest, actions };
}

export const RETIRED_EXTRA = ['creative-3d-web', 'fullstack-mastery'];

export function aliasIds(repo) {
  const file = path.join(repo, 'registries', 'skill-aliases.json');
  if (!fs.existsSync(file)) return [...RETIRED_EXTRA];
  const doc = JSON.parse(fs.readFileSync(file, 'utf8'));
  return [...new Set([...Object.keys(doc.aliases || {}), ...RETIRED_EXTRA])];
}

export function loadKnownHashes(repo) {
  const file = path.join(repo, 'registries', 'known-skill-hashes.json');
  if (!fs.existsSync(file)) return {};
  return JSON.parse(fs.readFileSync(file, 'utf8')).hashes || {};
}

export function planRetire(destRoot, ids) {
  return ids.filter((id) => {
    const dir = path.join(destRoot, id);
    return fs.existsSync(dir) && fs.statSync(dir).isDirectory();
  }).map((id) => ({ id, action: 'retire' }));
}

export function applyRetire({ destRoot, ids, backupRoot }) {
  const present = planRetire(destRoot, ids).map((item) => item.id);
  if (!present.length) return [];
  const mark = stamp();
  const applied = [];
  for (const id of present) {
    const dst = path.join(destRoot, id);
    const bucket = path.join(backupRoot, mark, id);
    fs.mkdirSync(path.dirname(bucket), { recursive: true });
    moveDir(dst, bucket);
    applied.push({ id, backup: bucket });
  }
  if (fs.existsSync(path.join(destRoot, '.ow-manifest.json'))) {
    const manifest = readManifest(destRoot);
    for (const id of present) {
      if (manifest.skills) delete manifest.skills[id];
    }
    writeManifest(destRoot, manifest);
  }
  return applied;
}

function copyTree(src, dst) {
  fs.mkdirSync(dst, { recursive: true });
  for (const rel of listFiles(src)) {
    const from = path.join(src, ...rel.split('/'));
    const to = path.join(dst, ...rel.split('/'));
    fs.mkdirSync(path.dirname(to), { recursive: true });
    fs.copyFileSync(from, to);
  }
}

export function hashBytes(buf) {
  if (!buf.includes(0)) return Buffer.from(buf.toString('utf8').replace(/\r\n/g, '\n'), 'utf8');
  return buf;
}

export function defaultBackupRoot(home) {
  if (process.env.OW_BACKUP_ROOT) return process.env.OW_BACKUP_ROOT;
  if (process.env.OW_SYNC_HOME) return path.join(process.env.OW_SYNC_HOME, '.orchestra-backups');
  if (process.platform === 'win32') {
    const base = process.env.LOCALAPPDATA || path.join(home, 'AppData', 'Local');
    return path.join(base, 'orchestra', 'backups');
  }
  return path.join(home, '.local', 'state', 'orchestra', 'backups');
}

function stamp() {
  return new Date().toISOString().replace(/[:.]/g, '-');
}

function moveDir(from, to) {
  fs.mkdirSync(path.dirname(to), { recursive: true });
  let last;
  for (let attempt = 0; attempt < 6; attempt++) {
    try {
      fs.renameSync(from, to);
      return;
    } catch (error) {
      last = error;
      if (error.code === 'EXDEV') break;
      if (!['EPERM', 'EBUSY', 'EACCES'].includes(error.code)) throw error;
      Atomics.wait(new Int32Array(new SharedArrayBuffer(4)), 0, 0, 200);
    }
  }
  if (!last || last.code !== 'EXDEV') throw last;
  fs.cpSync(from, to, { recursive: true });
  fs.rmSync(from, { recursive: true, force: true });
  if (fs.existsSync(from)) throw last;
}

export function applySync({ sourceRoot, destRoot, actions, backupRoot, source = 'orchestra' }) {
  const manifest = readManifest(destRoot);
  manifest.skills = manifest.skills || {};
  const applied = [];
  for (const item of actions) {
    if (item.action !== 'copy' && item.action !== 'replace') continue;
    const src = path.join(sourceRoot, item.id);
    const dst = path.join(destRoot, item.id);
    const stage = fs.mkdtempSync(path.join(os.tmpdir(), 'ow-sync-'));
    const staged = path.join(stage, item.id);
    copyTree(src, staged);
    if (hashTree(staged) !== item.srcHash) {
      fs.rmSync(stage, { recursive: true, force: true });
      throw new Error(`staged copy of ${item.id} did not match the planned hash`);
    }
    let bucket = '';
    try {
      if (item.action === 'replace') {
        bucket = path.join(backupRoot, stamp(), item.id);
        fs.mkdirSync(path.dirname(bucket), { recursive: true });
        moveDir(dst, bucket);
        applied.push({ id: item.id, backup: bucket, previous: item.dstHash });
      }
      fs.mkdirSync(path.dirname(dst), { recursive: true });
      moveDir(staged, dst);
    } catch (error) {
      if (bucket && !fs.existsSync(dst) && fs.existsSync(bucket)) moveDir(bucket, dst);
      fs.rmSync(stage, { recursive: true, force: true });
      throw error;
    }
    fs.rmSync(stage, { recursive: true, force: true });
    manifest.skills[item.id] = { sha256: item.srcHash, source };
  }
  for (const item of actions) {
    if (item.action === 'unchanged') {
      manifest.skills[item.id] = { sha256: item.srcHash, source };
    }
  }
  writeManifest(destRoot, manifest);
  return applied;
}

export function restoreStamp({ destRoot, backupStampDir, apply }) {
  const entries = fs.readdirSync(backupStampDir).filter((name) => {
    return fs.statSync(path.join(backupStampDir, name)).isDirectory();
  });
  const actions = [];
  for (const id of entries) {
    const src = path.join(backupStampDir, id);
    const dst = path.join(destRoot, id);
    const srcHash = hashTree(src);
    if (!fs.existsSync(dst)) {
      actions.push({ id, action: 'copy', srcHash });
      continue;
    }
    const dstHash = hashTree(dst);
    if (dstHash === srcHash) actions.push({ id, action: 'unchanged', srcHash });
    else actions.push({ id, action: 'replace', srcHash, dstHash });
  }
  if (apply) {
    applySync({
      sourceRoot: backupStampDir,
      destRoot,
      actions,
      backupRoot: path.join(backupStampDir, '..', 'restore-previous'),
    });
  }
  return actions;
}

export function canonicalIds(repo) {
  const stack = JSON.parse(fs.readFileSync(path.join(repo, 'registries', 'host-stack.json'), 'utf8'));
  return [...stack.skills].sort();
}

export function renderAgent(markdown, host) {
  const text = markdown.replace(/\r\n/g, '\n');
  const match = text.match(/^---\n([\s\S]*?)\n---\n?/);
  let body = text;
  let name = 'agent';
  let description = 'Orchestra agent';
  if (match) {
    body = text.slice(match[0].length).replace(/^\n/, '');
    const nameMatch = match[1].match(/^name:\s*(.+)$/m);
    const descMatch = match[1].match(/^description:\s*(.+)$/m);
    if (nameMatch) name = nameMatch[1].trim();
    if (descMatch) description = descMatch[1].trim();
  }
  const lines = ['---', `name: ${name}`, `description: ${description}`];
  if (host === 'cursor') lines.push('model: inherit', 'readonly: true', 'is_background: false');
  lines.push('---', '');
  return lines.join('\n') + (body.endsWith('\n') ? body : `${body}\n`);
}

export function syncAgents({ repo, destDir, host, apply }) {
  const src = path.join(repo, 'agents');
  if (!fs.existsSync(src)) return [];
  const actions = [];
  for (const name of fs.readdirSync(src).filter((item) => item.endsWith('.md')).sort()) {
    const body = renderAgent(fs.readFileSync(path.join(src, name), 'utf8'), host);
    const dest = path.join(destDir, name);
    const current = fs.existsSync(dest) ? fs.readFileSync(dest, 'utf8') : '';
    const action = !current ? 'copy' : current === body ? 'unchanged' : 'replace';
    actions.push({ id: name, action });
    if (apply && action !== 'unchanged') {
      fs.mkdirSync(destDir, { recursive: true });
      fs.writeFileSync(dest, body);
    }
  }
  return actions;
}

export function syncCommands({ repo, destDir, apply }) {
  const src = path.join(repo, 'commands');
  if (!fs.existsSync(src)) return [];
  const actions = [];
  for (const name of fs.readdirSync(src).filter((item) => item.endsWith('.md')).sort()) {
    const body = fs.readFileSync(path.join(src, name));
    const dest = path.join(destDir, name);
    const current = fs.existsSync(dest) ? fs.readFileSync(dest) : null;
    const action = !current ? 'copy' : current.equals(body) ? 'unchanged' : 'replace';
    actions.push({ id: name, action });
    if (apply && action !== 'unchanged') {
      fs.mkdirSync(destDir, { recursive: true });
      fs.writeFileSync(dest, body);
    }
  }
  return actions;
}

export function orchestraAdapter(body) {
  const text = body.replace(/\r\n/g, '\n');
  if (text.includes('<!-- orchestra:start -->')) return text.endsWith('\n') ? text : `${text}\n`;
  const inner = text.endsWith('\n') ? text : `${text}\n`;
  return `<!-- orchestra:start -->\n${inner}<!-- orchestra:end -->\n`;
}

export function adapterLooksLikeOrchestra(text) {
  return text.includes('ORCHESTRA = CONTROL PLANE')
    || text.includes('BEGIN ORCHESTRA')
    || text.includes('<!-- orchestra:start -->');
}

function spliceAdapter(current, body) {
  const start = '<!-- orchestra:start -->';
  const end = '<!-- orchestra:end -->';
  const text = current.replace(/\r\n/g, '\n');
  const inner = body.replace(/\r\n/g, '\n');
  const at = text.indexOf(start);
  const close = text.indexOf(end);
  if (at < 0 || close < at) return inner;
  return text.slice(0, at) + inner + text.slice(close + end.length);
}

export function planAdapter({ file, body, maxBytes = 0 }) {
  if (maxBytes && Buffer.byteLength(body) > maxBytes) {
    throw new Error(`${path.basename(file)} exceeds ${maxBytes} bytes`);
  }
  if (!fs.existsSync(file)) return { action: 'copy' };
  const current = fs.readFileSync(file, 'utf8');
  if (!current.trim()) return { action: 'replace', body };
  if (!adapterLooksLikeOrchestra(current)) return { action: 'leave' };
  const next = spliceAdapter(current, body);
  if (current.replace(/\r\n/g, '\n') === next) return { action: 'unchanged' };
  return { action: 'replace', body: next };
}

export function applyAdapter({ file, body, backupRoot, maxBytes = 0 }) {
  const decision = planAdapter({ file, body, maxBytes });
  if (decision.action === 'unchanged' || decision.action === 'leave') return decision;
  body = decision.body || body;
  if (fs.existsSync(file)) {
    const bucket = path.join(backupRoot, stamp(), path.basename(file));
    fs.mkdirSync(path.dirname(bucket), { recursive: true });
    fs.copyFileSync(file, bucket);
    decision.backup = bucket;
  }
  fs.mkdirSync(path.dirname(file), { recursive: true });
  fs.writeFileSync(file, body);
  return decision;
}

const CLAUDE_MCP = {
  shadcn: { command: 'npx', args: ['shadcn@4.23.0', 'mcp'] },
  'magic-ui': { command: 'npx', args: ['-y', '@magicuidesign/mcp@2.1.0'] },
};

export function claudeMcpPath(home) {
  return path.join(home, '.claude.json');
}

export function planClaudeMcp(file) {
  if (!fs.existsSync(file)) return { action: 'add', added: Object.keys(CLAUDE_MCP) };
  const json = JSON.parse(fs.readFileSync(file, 'utf8'));
  const servers = json.mcpServers || {};
  const added = Object.keys(CLAUDE_MCP).filter((name) => !servers[name]);
  return { action: added.length ? 'add' : 'unchanged', added };
}

export function applyClaudeMcp({ file, backupRoot }) {
  const decision = planClaudeMcp(file);
  if (decision.action !== 'add') return decision;
  const existed = fs.existsSync(file);
  const json = existed ? JSON.parse(fs.readFileSync(file, 'utf8')) : {};
  json.mcpServers = json.mcpServers || {};
  for (const name of decision.added) json.mcpServers[name] = CLAUDE_MCP[name];
  if (existed) {
    const bucket = path.join(backupRoot, stamp(), 'claude.json');
    fs.mkdirSync(path.dirname(bucket), { recursive: true });
    fs.copyFileSync(file, bucket);
    decision.backup = bucket;
  }
  fs.mkdirSync(path.dirname(file), { recursive: true });
  fs.writeFileSync(file, JSON.stringify(json, null, 2) + '\n');
  return decision;
}

function syncSurface(home, selected, apply, backupRoot) {
  const want = new Set(selected);
  if (want.has('cursor')) {
    printActions('cursor-agents', syncAgents({
      repo: REPO,
      destDir: path.join(home, '.cursor', 'agents'),
      host: 'cursor',
      apply,
    }));
    printActions('cursor-commands', syncCommands({
      repo: REPO,
      destDir: path.join(REPO, '.cursor', 'commands'),
      apply,
    }));
    printActions('cursor-user-commands', syncCommands({
      repo: REPO,
      destDir: path.join(home, '.cursor', 'commands'),
      apply,
    }));
  }
  if (want.has('claude')) {
    printActions('claude-agents', syncAgents({
      repo: REPO,
      destDir: path.join(home, '.claude', 'agents'),
      host: 'claude',
      apply,
    }));
    printActions('claude-commands', syncCommands({
      repo: REPO,
      destDir: path.join(home, '.claude', 'commands'),
      apply,
    }));
    const claudeBody = orchestraAdapter(fs.readFileSync(path.join(REPO, 'CLAUDE.md'), 'utf8'));
    const claudeFile = path.join(home, '.claude', 'CLAUDE.md');
    const claudePlan = planAdapter({ file: claudeFile, body: claudeBody });
    console.log(`claude-adapter  ${claudePlan.action}  CLAUDE.md`);
    if (apply) applyAdapter({ file: claudeFile, body: claudeBody, backupRoot });
    const mcpFile = path.join(home, '.claude.json');
    const mcpPlan = planClaudeMcp(mcpFile);
    console.log(`claude-mcp  ${mcpPlan.action}  ${mcpPlan.added.join(',')}`);
    if (apply) applyClaudeMcp({ file: mcpFile, backupRoot });
  }
  if (want.has('antigravity')) {
    const body = orchestraAdapter(fs.readFileSync(path.join(REPO, 'kit', 'antigravity', 'GEMINI.md'), 'utf8'));
    const file = path.join(home, '.gemini', 'GEMINI.md');
    const decision = planAdapter({ file, body, maxBytes: 24 * 1024 });
    console.log(`gemini-adapter  ${decision.action}  GEMINI.md`);
    if (apply) applyAdapter({ file, body, backupRoot, maxBytes: 24 * 1024 });
  }
  if (want.has('codex')) {
    const body = orchestraAdapter(fs.readFileSync(path.join(REPO, 'kit', 'codex', 'AGENTS.md'), 'utf8'));
    const file = path.join(home, '.codex', 'AGENTS.md');
    const decision = planAdapter({ file, body });
    console.log(`codex-adapter  ${decision.action}  AGENTS.md`);
    if (apply) applyAdapter({ file, body, backupRoot });
  }
}

function parseArgs(argv) {
  const out = { apply: false, backupAndReplace: false, hosts: [], restore: '', retire: [] };
  for (let i = 0; i < argv.length; i++) {
    const arg = argv[i];
    if (arg === '--apply') out.apply = true;
    else if (arg === '--backup-and-replace') out.backupAndReplace = true;
    else if (arg === '--host') out.hosts.push(argv[++i]);
    else if (arg === '--restore') out.restore = argv[++i];
    else if (arg === '--retire') out.retire.push(argv[++i]);
    else throw new Error(`unknown argument ${arg}`);
  }
  return out;
}

function printActions(host, actions) {
  for (const item of actions) {
    console.log(`${host}  ${item.action}  ${item.id}`);
  }
}

function runCli(argv) {
  const args = parseArgs(argv);
  const home = process.env.OW_SYNC_HOME || os.homedir();
  const dirs = hostDirs(home);
  const selected = args.hosts.length ? args.hosts : Object.keys(dirs);
  const backupRoot = defaultBackupRoot(home);
  const ids = canonicalIds(REPO);
  const localRoot = path.join(REPO, 'skills-local');
  const localIds = listSkillIds(localRoot);
  for (const host of selected) {
    const dest = dirs[host];
    if (!dest) throw new Error(`unknown host ${host}`);
    if (args.restore) {
      const stampDir = path.join(backupRoot, args.restore);
      const actions = restoreStamp({ destRoot: dest, backupStampDir: stampDir, apply: args.apply });
      printActions(host, actions);
      continue;
    }
    if (args.retire.length) {
      const manifest = readManifest(dest);
      for (const id of args.retire) {
        const recorded = manifest.skills && manifest.skills[id];
        const dst = path.join(dest, id);
        if (!recorded || recorded.source !== 'orchestra' || !fs.existsSync(dst)) {
          console.log(`${host}  skip-retire  ${id}`);
          continue;
        }
        if (hashTree(dst) !== recorded.sha256) {
          console.log(`${host}  conflict  ${id}`);
          continue;
        }
        if (!args.apply) {
          console.log(`${host}  retire  ${id}`);
          continue;
        }
        const bucket = path.join(backupRoot, stamp(), id);
        fs.mkdirSync(path.dirname(bucket), { recursive: true });
        moveDir(dst, bucket);
        delete manifest.skills[id];
        writeManifest(dest, manifest);
        console.log(`${host}  retired  ${id}`);
      }
      continue;
    }
    if (args.apply) fs.mkdirSync(dest, { recursive: true });
    const canonical = planSync({
      sourceIds: ids,
      sourceRoot: path.join(REPO, 'skills'),
      destRoot: dest,
      backupAndReplace: args.backupAndReplace,
      knownHashes: loadKnownHashes(REPO),
    });
    printActions(host, canonical.actions);
    if (args.apply) {
      applySync({
        sourceRoot: path.join(REPO, 'skills'),
        destRoot: dest,
        actions: canonical.actions,
        backupRoot,
      });
    }
    if (localIds.length) {
      const local = planSync({
        sourceIds: localIds,
        sourceRoot: localRoot,
        destRoot: dest,
        backupAndReplace: args.backupAndReplace,
      });
      printActions(host, local.actions);
      if (args.apply) {
        applySync({
          sourceRoot: localRoot,
          destRoot: dest,
          actions: local.actions,
          backupRoot,
          source: 'local',
        });
      }
    }
    const retired = planRetire(dest, aliasIds(REPO));
    printActions(host, retired);
    if (args.apply && retired.length) {
      applyRetire({ destRoot: dest, ids: retired.map((item) => item.id), backupRoot });
    }
  }
  if (!args.restore && !args.retire.length) syncSurface(home, selected, args.apply, backupRoot);
  if (!args.apply) console.log('dry-run. Pass --apply to write.');
}

if (process.argv[1] && path.resolve(process.argv[1]) === path.resolve(path.join(import.meta.dirname, 'ow-sync.mjs'))) {
  try {
    runCli(process.argv.slice(2));
  } catch (err) {
    console.error(err.message || err);
    process.exit(1);
  }
}
