// Validate every registries/*.json against registries/schemas/<name>.schema.json.
// Fails when a registry has no schema, or when a schema file does not parse.
//   node runtime/tools/check-schemas.mjs
import fs from 'node:fs';
import path from 'node:path';
import Ajv from 'ajv';
import addFormats from 'ajv-formats';

const REPO = path.join(import.meta.dirname, '..', '..');
const registries = path.join(REPO, 'registries');
const schemaDir = path.join(registries, 'schemas');
const schemaFiles = fs.readdirSync(schemaDir).filter((name) => name.endsWith('.json'));
if (schemaFiles.length === 0) {
  console.error('no schemas found');
  process.exit(1);
}

const ajv = new Ajv({ allErrors: true, strict: false });
addFormats(ajv);
const problems = [];
const dataFiles = fs.readdirSync(registries).filter((name) => name.endsWith('.json'));

for (const name of dataFiles) {
  const schemaName = name.replace(/\.json$/, '.schema.json');
  const schemaPath = path.join(schemaDir, schemaName);
  if (!fs.existsSync(schemaPath)) {
    problems.push(`${name} has no schema ${schemaName}`);
    continue;
  }
  const schema = JSON.parse(fs.readFileSync(schemaPath, 'utf8'));
  const data = JSON.parse(fs.readFileSync(path.join(registries, name), 'utf8'));
  const validate = ajv.compile(schema);
  if (!validate(data)) {
    for (const err of validate.errors.slice(0, 8)) {
      problems.push(`${name} ${err.instancePath || '/'} ${err.message}`);
    }
  }
}

console.log(`registries: ${dataFiles.length}`);
console.log(`schemas:    ${schemaFiles.length}`);
console.log(`problems:   ${problems.length}`);
for (const problem of problems) console.log(`  ${problem}`);
process.exit(problems.length ? 1 : 0);
