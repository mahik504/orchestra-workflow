// 4.1 contracts: domains point at existing capabilities, technology rows
// do not claim usage, MCP rows are not install instructions, and the
// design-contract templates satisfy the packs.
const fs = require('fs');
const path = require('path');
const { checkDir } = require('./check-design-contract');

const REPO = path.join(__dirname, '..', '..');
const problems = [];

const graph = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'design-resource-graph.json'), 'utf8'));
const domains = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'design-domains.json'), 'utf8'));
const policy = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'technology-policy.json'), 'utf8'));
const mcps = JSON.parse(fs.readFileSync(path.join(REPO, 'registries', 'mcp-contracts.json'), 'utf8'));
const caps = new Set(Object.keys(graph.capabilities || {}));

const seenDomains = new Set();
let defaults = {};
for (const domain of domains.domains || []) {
  if (seenDomains.has(domain.id)) problems.push(`duplicate domain ${domain.id}`);
  seenDomains.add(domain.id);
  if (!caps.has(domain.parent)) problems.push(`domain ${domain.id} parent ${domain.parent} is not a capability`);
  for (const key of ['triggers', 'exclusions', 'documents', 'quality_gates']) {
    if (!Array.isArray(domain[key]) || domain[key].length === 0) problems.push(`domain ${domain.id} missing ${key}`);
  }
  if (domain.default) {
    if (defaults[domain.parent]) problems.push(`two defaults for ${domain.parent}`);
    defaults[domain.parent] = domain.id;
  }
}
for (const parent of ['premium-website', '3d-portfolio', 'saas-dashboard', 'operator-hud']) {
  if (!defaults[parent]) problems.push(`no default domain for ${parent}`);
}

const contractFields = [
  'id', 'description', 'triggers', 'exclusions', 'execution_phase', 'dependencies',
  'install_instructions', 'license', 'cost_class', 'compatibility', 'maintenance_status',
  'verification_status', 'expected_output', 'failure_fallback', 'usage_evidence', 'lifecycle_stage',
];
const seenLibs = new Set();
for (const lib of policy.libraries || []) {
  if (seenLibs.has(lib.id)) problems.push(`duplicate library ${lib.id}`);
  seenLibs.add(lib.id);
  for (const field of contractFields) {
    const value = lib[field];
    if (value == null || value === '' || (Array.isArray(value) && value.length === 0 && field !== 'dependencies')) {
      if (field === 'dependencies' && Array.isArray(value)) continue;
      problems.push(`library ${lib.id} missing ${field}`);
    }
  }
  if (!String(lib.usage_evidence).includes('none recorded')) problems.push(`library ${lib.id} claims usage`);
  if (lib.lifecycle_stage !== 'discovered') problems.push(`library ${lib.id} lifecycle is ${lib.lifecycle_stage}`);
  if (!String(lib.install_instructions).toLowerCase().includes('do not install globally') && !String(lib.install_instructions).toLowerCase().includes('no package')) {
    problems.push(`library ${lib.id} install line does not block a global install`);
  }
}

for (const row of mcps.contracts || []) {
  for (const field of ['id', 'kind', 'invoke_when', 'skip_when', 'install', 'fallback', 'usage_evidence']) {
    if (!row[field]) problems.push(`mcp ${row.id || '?'} missing ${field}`);
  }
  if (!String(row.install).includes('Do not install globally')) problems.push(`mcp ${row.id} install is not refused`);
  if (!String(row.usage_evidence).includes('none recorded')) problems.push(`mcp ${row.id} claims usage`);
}

const templateDir = path.join(REPO, 'templates', 'contracts');
for (const profile of ['premium', 'spatial', 'backend', 'saas']) {
  for (const problem of checkDir(templateDir, profile)) problems.push(`templates ${profile}: ${problem}`);
}
const os = require('os');
const bad = fs.mkdtempSync(path.join(os.tmpdir(), 'orchestra-contract-'));
fs.writeFileSync(path.join(bad, 'DESIGN.md'), '## Brand identity\n\nSomeone.\n');
const badProblems = checkDir(bad, 'premium');
if (badProblems.length === 0) problems.push('incomplete design contract was accepted');
fs.rmSync(bad, { recursive: true, force: true });

const empty = fs.mkdtempSync(path.join(os.tmpdir(), 'orchestra-none-'));
if (checkDir(empty, 'none').length !== 0) problems.push('empty bug pack failed');
fs.rmSync(empty, { recursive: true, force: true });

console.log(`domains:   ${(domains.domains || []).length}`);
console.log(`libraries: ${(policy.libraries || []).length}`);
console.log(`mcp rows:  ${(mcps.contracts || []).length}`);
console.log(`problems:  ${problems.length}`);
for (const problem of problems) console.log('  ' + problem);
process.exit(problems.length ? 1 : 0);
