// Transactional skill sync. Dry-run unless --apply.
// Never deletes a skill that this tool did not install.
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
    return fs.statSync(dir).isDirectory() && fs.existsSync(path.join(dir, 'SKILL.md'));
  }).sort();
}

function listFiles(dir) {
  const out = [];
  const walk = (abs, rel) => {
    for (const name of fs.readdirSync(abs).sort()) {
      const child = path.join(abs, name);
      const next = rel ? `${rel}/${name}` : name;
      if (fs.statSync(child).isDirectory()) walk(child, next);
      else out.push(next);
    }
  };
  walk(dir, '');
  return out;
}

export function hashTree(dir) {
  const hash = createHash('sha256');
  for (const rel of listFiles(dir)) {
    hash.update(rel);
    hash.update('\0');
    hash.update(fs.readFileSync(path.join(dir, ...rel.split('/'))));
    hash.update('\0');
  }
  return hash.digest('hex');
}

function readManifest(destRoot) {
  const file = path.join(destRoot, '.ow-manifest.json');
  if (!fs.existsSync(file)) return { skills: {} };
  return JSON.parse(fs.readFileSync(file, 'utf8'));
}

function writeManifest(destRoot, manifest) {
  fs.mkdirSync(destRoot, { recursive: true });
  fs.writeFileSync(path.join(destRoot, '.ow-manifest.json'), JSON.stringify(manifest, null, 2) + '\n');
}

export function planSync({ sourceIds, sourceRoot, destRoot, backupAndReplace }) {
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
    const recorded = manifest.skills && manifest.skills[id] && manifest.skills[id].sha256;
    if (dstHash === srcHash) {
      actions.push({ id, action: 'unchanged', srcHash });
    } else if (recorded === dstHash || backupAndReplace) {
      actions.push({ id, action: 'replace', srcHash, dstHash });
    } else {
      actions.push({ id, action: 'conflict', srcHash, dstHash });
    }
  }
  return { manifest, actions };
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

export function defaultBackupRoot(home) {
  if (process.env.OW_BACKUP_ROOT) return process.env.OW_BACKUP_ROOT;
  if (process.platform === 'win32') {
    const base = process.env.LOCALAPPDATA || path.join(home, 'AppData', 'Local');
    return path.join(base, 'orchestra', 'backups');
  }
  return path.join(home, '.local', 'state', 'orchestra', 'backups');
}

function stamp() {
  return new Date().toISOString().replace(/[:.]/g, '-');
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
    if (item.action === 'replace') {
      const bucket = path.join(backupRoot, stamp(), item.id);
      fs.mkdirSync(path.dirname(bucket), { recursive: true });
      fs.renameSync(dst, bucket);
      applied.push({ id: item.id, backup: bucket, previous: item.dstHash });
    }
    fs.mkdirSync(path.dirname(dst), { recursive: true });
    fs.renameSync(staged, dst);
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
        fs.renameSync(dst, bucket);
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
  }
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
