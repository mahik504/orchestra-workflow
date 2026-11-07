// Write an adapter between <!-- orchestra:start --> and <!-- orchestra:end -->.
// Creates the file when it is missing. Leaves an unmarked file alone.
//
//   node runtime/tools/stamp-adapter.mjs --dest path --source path
import fs from 'node:fs';
import path from 'node:path';

const START = '<!-- orchestra:start -->';
const END = '<!-- orchestra:end -->';

export function stampMarked(file, body) {
  const block = `${START}\n${body.trim()}\n${END}\n`;
  if (!fs.existsSync(file)) {
    fs.mkdirSync(path.dirname(file), { recursive: true });
    fs.writeFileSync(file, block);
    return 'created';
  }
  const text = fs.readFileSync(file, 'utf8');
  const a = text.indexOf(START);
  const b = text.indexOf(END);
  if (a >= 0 && b > a) {
    const next = text.slice(0, a) + block.trimEnd() + text.slice(b + END.length);
    fs.writeFileSync(file, next.endsWith('\n') ? next : next + '\n');
    return 'updated';
  }
  return 'conflict';
}

function arg(name) {
  const i = process.argv.indexOf(name);
  return i >= 0 ? process.argv[i + 1] : '';
}

if (process.argv[1] && path.resolve(process.argv[1]).endsWith('stamp-adapter.mjs')) {
  const dest = arg('--dest');
  const source = arg('--source');
  if (!dest || !source) {
    console.error('usage: stamp-adapter.mjs --dest file --source file');
    process.exit(2);
  }
  const result = stampMarked(dest, fs.readFileSync(source, 'utf8'));
  console.log(`${result} ${dest}`);
  if (result === 'conflict') {
    console.log('left unmarked file untouched');
  }
}
