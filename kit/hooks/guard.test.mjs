import { test } from 'node:test';
import { spawnSync } from 'node:child_process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const script = path.join(path.dirname(fileURLToPath(import.meta.url)), 'guard.mjs');

function run(id, payload, env = {}) {
  const res = spawnSync(process.execPath, [script, id], {
    input: JSON.stringify(payload),
    encoding: 'utf8',
    env: { ...process.env, ORCHESTRA_HOOK_PROFILE: 'standard', ...env },
  });
  const out = res.stdout ? JSON.parse(res.stdout) : {};
  return { code: res.status, out };
}

test('skills add --all is denied', () => {
  const { out } = run('shell', { command: 'npx skills add --all' });
  if (out.permission !== 'deny') throw new Error(JSON.stringify(out));
});

test('node_modules delete is allowed', () => {
  const { out } = run('shell', { command: 'rm -rf node_modules' });
  if (out.permission !== 'allow') throw new Error(JSON.stringify(out));
});

test('drive root delete is denied', () => {
  const { out } = run('shell', { command: 'rm -rf C:\\' });
  if (out.permission !== 'deny') throw new Error(JSON.stringify(out));
});

test('git push asks', () => {
  const { out } = run('shell', { command: 'git push origin main' });
  if (out.permission !== 'ask') throw new Error(JSON.stringify(out));
});

test('minimal profile allows a push', () => {
  const { out } = run('shell', { command: 'git push origin main' }, { ORCHESTRA_HOOK_PROFILE: 'minimal' });
  if (out.permission !== 'allow') throw new Error(JSON.stringify(out));
});

test('.env read is denied', () => {
  const { out } = run('read', { file_path: 'C:/proj/.env' });
  if (out.permission !== 'deny') throw new Error(JSON.stringify(out));
});

test('off profile allows everything', () => {
  const { out } = run('shell', { command: 'npx skills add --all' }, { ORCHESTRA_HOOK_PROFILE: 'off' });
  if (out.permission !== 'allow') throw new Error(JSON.stringify(out));
});

test('strict stop without evidence exits 2', () => {
  const { code, out } = run('stop', { last_message: 'DONE' }, { ORCHESTRA_HOOK_PROFILE: 'strict' });
  if (code !== 2) throw new Error(`code ${code} ${JSON.stringify(out)}`);
});
