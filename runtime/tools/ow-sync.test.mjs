import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import {
  applyAdapter,
  applyClaudeMcp,
  applyRetire,
  applySync,
  claudeMcpPath,
  defaultBackupRoot,
  hashBytes,
  hashTree,
  hashTreeRaw,
  planAdapter,
  planClaudeMcp,
  planSync,
  renderAgent,
  restoreStamp,
} from './ow-sync.mjs';

function skill(dir, id, body) {
  const root = path.join(dir, id);
  fs.mkdirSync(root, { recursive: true });
  fs.writeFileSync(path.join(root, 'SKILL.md'), body);
  return root;
}

test('dry-run plans a copy and does not write the destination', () => {
  const home = fs.mkdtempSync(path.join(os.tmpdir(), 'ow-home-'));
  const source = path.join(home, 'source');
  const dest = path.join(home, 'dest');
  skill(source, 'alpha', '---\nname: alpha\n---\n# Alpha\n');
  const plan = planSync({
    sourceIds: ['alpha'],
    sourceRoot: source,
    destRoot: dest,
    backupAndReplace: false,
  });
  assert.equal(plan.actions[0].action, 'copy');
  assert.equal(fs.existsSync(path.join(dest, 'alpha')), false);
});

test('apply then a second plan is unchanged', () => {
  const home = fs.mkdtempSync(path.join(os.tmpdir(), 'ow-home-'));
  const source = path.join(home, 'source');
  const dest = path.join(home, 'dest');
  skill(source, 'alpha', '---\nname: alpha\n---\n# Alpha\n');
  const first = planSync({ sourceIds: ['alpha'], sourceRoot: source, destRoot: dest, backupAndReplace: false });
  applySync({ sourceRoot: source, destRoot: dest, actions: first.actions, backupRoot: path.join(home, 'bak') });
  const second = planSync({ sourceIds: ['alpha'], sourceRoot: source, destRoot: dest, backupAndReplace: false });
  assert.equal(second.actions[0].action, 'unchanged');
});

test('a local edit is a conflict and the file stays', () => {
  const home = fs.mkdtempSync(path.join(os.tmpdir(), 'ow-home-'));
  const source = path.join(home, 'source');
  const dest = path.join(home, 'dest');
  skill(source, 'alpha', '---\nname: alpha\n---\n# Alpha\n');
  const first = planSync({ sourceIds: ['alpha'], sourceRoot: source, destRoot: dest, backupAndReplace: false });
  applySync({ sourceRoot: source, destRoot: dest, actions: first.actions, backupRoot: path.join(home, 'bak') });
  fs.writeFileSync(path.join(dest, 'alpha', 'SKILL.md'), 'user edit\n');
  fs.writeFileSync(path.join(source, 'alpha', 'SKILL.md'), '---\nname: alpha\n---\n# Alpha two\n');
  const plan = planSync({ sourceIds: ['alpha'], sourceRoot: source, destRoot: dest, backupAndReplace: false });
  assert.equal(plan.actions[0].action, 'conflict');
  assert.equal(fs.readFileSync(path.join(dest, 'alpha', 'SKILL.md'), 'utf8'), 'user edit\n');
});

test('backup-and-replace keeps the previous tree and restore puts it back', () => {
  const home = fs.mkdtempSync(path.join(os.tmpdir(), 'ow-home-'));
  const source = path.join(home, 'source');
  const dest = path.join(home, 'dest');
  const backupRoot = path.join(home, 'bak');
  skill(source, 'alpha', 'version one\n');
  const first = planSync({ sourceIds: ['alpha'], sourceRoot: source, destRoot: dest, backupAndReplace: false });
  applySync({ sourceRoot: source, destRoot: dest, actions: first.actions, backupRoot });
  fs.writeFileSync(path.join(source, 'alpha', 'SKILL.md'), 'version two\n');
  const next = planSync({ sourceIds: ['alpha'], sourceRoot: source, destRoot: dest, backupAndReplace: false });
  assert.equal(next.actions[0].action, 'replace');
  const applied = applySync({ sourceRoot: source, destRoot: dest, actions: next.actions, backupRoot });
  assert.equal(fs.readFileSync(path.join(dest, 'alpha', 'SKILL.md'), 'utf8'), 'version two\n');
  assert.equal(fs.readFileSync(path.join(applied[0].backup, 'SKILL.md'), 'utf8'), 'version one\n');
  const stampDir = path.dirname(applied[0].backup);
  restoreStamp({ destRoot: dest, backupStampDir: stampDir, apply: true });
  assert.equal(hashTree(path.join(dest, 'alpha')), hashTree(applied[0].backup));
});

test('a local source label is not recorded as orchestra', () => {
  const home = fs.mkdtempSync(path.join(os.tmpdir(), 'ow-home-'));
  const source = path.join(home, 'source');
  const dest = path.join(home, 'dest');
  skill(source, 'docx', 'local skill\n');
  const plan = planSync({ sourceIds: ['docx'], sourceRoot: source, destRoot: dest, backupAndReplace: false });
  applySync({
    sourceRoot: source,
    destRoot: dest,
    actions: plan.actions,
    backupRoot: path.join(home, 'bak'),
    source: 'local',
  });
  const manifest = JSON.parse(fs.readFileSync(path.join(dest, '.ow-manifest.json'), 'utf8'));
  assert.equal(manifest.skills.docx.source, 'local');
});

test('an extra skill in the destination is left alone', () => {
  const home = fs.mkdtempSync(path.join(os.tmpdir(), 'ow-home-'));
  const source = path.join(home, 'source');
  const dest = path.join(home, 'dest');
  skill(source, 'alpha', 'alpha\n');
  skill(dest, 'user-installed', 'keep me\n');
  const plan = planSync({ sourceIds: ['alpha'], sourceRoot: source, destRoot: dest, backupAndReplace: false });
  applySync({ sourceRoot: source, destRoot: dest, actions: plan.actions, backupRoot: path.join(home, 'bak') });
  assert.equal(fs.readFileSync(path.join(dest, 'user-installed', 'SKILL.md'), 'utf8'), 'keep me\n');
  assert.equal(plan.actions.some((item) => item.id === 'user-installed'), false);
});

test('a known historical hash is safe to replace', () => {
  const home = fs.mkdtempSync(path.join(os.tmpdir(), 'ow-home-'));
  const source = path.join(home, 'source');
  const dest = path.join(home, 'dest');
  skill(dest, 'alpha', 'old orchestra copy\n');
  const oldHash = hashTree(path.join(dest, 'alpha'));
  skill(source, 'alpha', 'current orchestra copy\n');
  const plan = planSync({
    sourceIds: ['alpha'],
    sourceRoot: source,
    destRoot: dest,
    backupAndReplace: false,
    knownHashes: { alpha: [oldHash] },
  });
  assert.equal(plan.actions[0].action, 'replace');
});

test('an old alias folder moves to a backup and a user skill stays', () => {
  const home = fs.mkdtempSync(path.join(os.tmpdir(), 'ow-home-'));
  const dest = path.join(home, 'dest');
  skill(dest, 'orchestra-conductor', 'old name\n');
  skill(dest, 'user-installed', 'keep me\n');
  const applied = applyRetire({
    destRoot: dest,
    ids: ['orchestra-conductor', 'creative-3d-web'],
    backupRoot: path.join(home, 'bak'),
  });
  assert.equal(applied.length, 1);
  assert.equal(fs.existsSync(path.join(dest, 'orchestra-conductor')), false);
  assert.equal(fs.readFileSync(path.join(applied[0].backup, 'SKILL.md'), 'utf8'), 'old name\n');
  assert.equal(fs.readFileSync(path.join(dest, 'user-installed', 'SKILL.md'), 'utf8'), 'keep me\n');
});

test('cursor agent frontmatter is host-specific and claude stays short', () => {
  const source = '---\nname: ow-reviewer\ndescription: Review a pass.\n---\n\n# Body\n';
  const cursor = renderAgent(source, 'cursor');
  const claude = renderAgent(source, 'claude');
  assert.match(cursor, /model: inherit/);
  assert.match(cursor, /readonly: true/);
  assert.equal(claude.includes('model: inherit'), false);
  assert.match(claude, /# Body/);
});

test('an unmarked orchestra adapter is replaced and a personal file is left', () => {
  const home = fs.mkdtempSync(path.join(os.tmpdir(), 'ow-home-'));
  const orchestraFile = path.join(home, 'CLAUDE.md');
  const personal = path.join(home, 'NOTES.md');
  fs.writeFileSync(orchestraFile, 'ORCHESTRA = CONTROL PLANE\nOrchestra 4.1.0\n');
  fs.writeFileSync(personal, 'my notes\n');
  const body = '<!-- orchestra:start -->\nOrchestra 4.6.0\n<!-- orchestra:end -->\n';
  const replaced = applyAdapter({ file: orchestraFile, body, backupRoot: path.join(home, 'bak') });
  assert.equal(replaced.action, 'replace');
  assert.equal(fs.readFileSync(orchestraFile, 'utf8'), body);
  assert.equal(fs.existsSync(replaced.backup), true);
  const left = planAdapter({ file: personal, body });
  assert.equal(left.action, 'leave');
  assert.equal(fs.readFileSync(personal, 'utf8'), 'my notes\n');
  const empty = path.join(home, 'AGENTS.md');
  fs.writeFileSync(empty, '\n');
  const filled = applyAdapter({ file: empty, body, backupRoot: path.join(home, 'bak') });
  assert.equal(filled.action, 'replace');
  assert.equal(fs.readFileSync(empty, 'utf8'), body);
  const noted = path.join(home, 'VERSION-NOTE.md');
  fs.writeFileSync(noted, 'Orchestra 4.5 was the previous pin\n');
  assert.equal(planAdapter({ file: noted, body }).action, 'leave');
  const wrapped = path.join(home, 'WRAPPED.md');
  fs.writeFileSync(wrapped, '<!-- orchestra:start -->\nold\n<!-- orchestra:end -->\nmy note\n');
  const spliced = applyAdapter({ file: wrapped, body, backupRoot: path.join(home, 'bak') });
  assert.equal(spliced.action, 'replace');
  assert.equal(fs.readFileSync(wrapped, 'utf8').endsWith('my note\n'), true);
  assert.equal(fs.readFileSync(wrapped, 'utf8').includes('old'), false);
});

test('claude mcp adds server names and does not echo a secret', () => {
  const home = fs.mkdtempSync(path.join(os.tmpdir(), 'ow-home-'));
  const file = path.join(home, 'settings.json');
  const secret = 'sk-test-secret-value';
  fs.writeFileSync(file, JSON.stringify({ mcpServers: { stitch: { env: { STITCH_API_KEY: secret } } } }));
  const planned = planClaudeMcp(file);
  assert.deepEqual(planned.added, ['shadcn', 'magic-ui']);
  assert.equal(JSON.stringify(planned).includes(secret), false);
  applyClaudeMcp({ file, backupRoot: path.join(home, 'bak') });
  const after = JSON.parse(fs.readFileSync(file, 'utf8'));
  assert.equal(after.mcpServers.shadcn.command, 'npx');
  assert.equal(after.mcpServers.stitch.env.STITCH_API_KEY, secret);
  assert.equal(claudeMcpPath(home), path.join(home, '.claude.json'));
});

test('CRLF and LF hash as the same text', () => {
  const a = hashBytes(Buffer.from('a\r\nb\n'));
  const b = hashBytes(Buffer.from('a\nb\n'));
  assert.equal(a.toString('utf8'), b.toString('utf8'));
});

test('OW_SYNC_HOME owns the backup root', () => {
  const prev = process.env.OW_SYNC_HOME;
  const prevBackup = process.env.OW_BACKUP_ROOT;
  delete process.env.OW_BACKUP_ROOT;
  process.env.OW_SYNC_HOME = 'C:\\sync-home';
  assert.equal(defaultBackupRoot('C:\\other'), path.join('C:\\sync-home', '.orchestra-backups'));
  if (prev === undefined) delete process.env.OW_SYNC_HOME;
  else process.env.OW_SYNC_HOME = prev;
  if (prevBackup === undefined) delete process.env.OW_BACKUP_ROOT;
  else process.env.OW_BACKUP_ROOT = prevBackup;
});

test('a symlink inside a skill is refused', () => {
  const root = fs.mkdtempSync(path.join(os.tmpdir(), 'ow-link-'));
  const skillDir = path.join(root, 'demo');
  fs.mkdirSync(skillDir);
  fs.writeFileSync(path.join(skillDir, 'SKILL.md'), '---\nname: demo\n---\n');
  const outside = path.join(root, 'outside.txt');
  fs.writeFileSync(outside, 'secret');
  try {
    fs.symlinkSync(outside, path.join(skillDir, 'linked.txt'));
  } catch (error) {
    if (error.code === 'EPERM') return;
    throw error;
  }
  assert.throws(() => hashTree(skillDir), /symlink refused/);
});

test('a recorded raw hash is safe to replace after the hash normalization', () => {
  const root = fs.mkdtempSync(path.join(os.tmpdir(), 'ow-raw-'));
  const source = path.join(root, 'src');
  const dest = path.join(root, 'dst');
  fs.mkdirSync(path.join(source, 'demo'), { recursive: true });
  fs.mkdirSync(path.join(dest, 'demo'), { recursive: true });
  fs.writeFileSync(path.join(source, 'demo', 'SKILL.md'), '---\nname: demo\n---\nnew\n');
  fs.writeFileSync(path.join(dest, 'demo', 'SKILL.md'), '---\r\nname: demo\r\n---\r\nold\r\n');
  const raw = hashTreeRaw(path.join(dest, 'demo'));
  fs.mkdirSync(dest, { recursive: true });
  fs.writeFileSync(path.join(dest, '.ow-manifest.json'), JSON.stringify({ skills: { demo: { sha256: raw } } }));
  const planned = planSync({
    sourceIds: ['demo'],
    sourceRoot: source,
    destRoot: dest,
    knownHashes: {},
  });
  assert.equal(planned.actions[0].action, 'replace');
});
