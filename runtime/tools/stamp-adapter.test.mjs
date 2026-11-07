import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { stampMarked } from './stamp-adapter.mjs';

test('creates a marked file and updates inside the markers', () => {
  const dir = fs.mkdtempSync(path.join(os.tmpdir(), 'ow-stamp-'));
  const file = path.join(dir, 'AGENTS.md');
  assert.equal(stampMarked(file, 'one'), 'created');
  assert.match(fs.readFileSync(file, 'utf8'), /<!-- orchestra:start -->\none/);
  assert.equal(stampMarked(file, 'two'), 'updated');
  const text = fs.readFileSync(file, 'utf8');
  assert.match(text, /two/);
  assert.equal(text.includes('one'), false);
});

test('an unmarked file is a conflict and stays', () => {
  const dir = fs.mkdtempSync(path.join(os.tmpdir(), 'ow-stamp-'));
  const file = path.join(dir, 'AGENTS.md');
  fs.writeFileSync(file, 'user contract\n');
  assert.equal(stampMarked(file, 'orchestra'), 'conflict');
  assert.equal(fs.readFileSync(file, 'utf8'), 'user contract\n');
});
