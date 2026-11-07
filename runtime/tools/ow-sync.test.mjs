import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { applySync, hashTree, planSync, restoreStamp } from './ow-sync.mjs';

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
