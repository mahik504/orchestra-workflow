// Antigravity budget + MCP health check. Same rules as orchestra doctor.
// Does not print keys. Does not write config.
//
//   node runtime/tools/doctor-ag.js
const fs = require('fs');
const path = require('path');
const os = require('os');

const BANNED = ['science', 'data-agent-kit-plugin'];
const ALIAS = {
  stitchmcp: 'stitch',
  ['orchestra' + '-brain']: 'vault-memory',
  'vault-memory': 'vault-memory',
  supabase: 'supabase',
};

function classifyMCP(name) {
  const key = ALIAS[String(name).toLowerCase()] || String(name).toLowerCase();
  if (['playwright', 'context7', 'stitch', 'vault-memory', 'browser'].includes(key)) return 'HEALTHY';
  if (key === 'supabase') return 'AUTH_REQUIRED';
  return 'OPTIONAL';
}

function pluginEnabled(name, prefs, pluginDir) {
  if (prefs[name] && typeof prefs[name].enabled === 'boolean') return prefs[name].enabled;
  const man = path.join(pluginDir, name, 'plugin.json');
  if (fs.existsSync(man)) {
    try {
      const j = JSON.parse(fs.readFileSync(man, 'utf8'));
      if (j.disabled === true) return false;
    } catch (_) { /* treat as enabled */ }
  }
  return true;
}

const home = process.env.USERPROFILE || process.env.HOME || os.homedir();
const cfgDir = path.join(home, '.gemini', 'config');
const cfgPath = path.join(cfgDir, 'config.json');
const pluginsDir = path.join(cfgDir, 'plugins');
const skillsDir = path.join(cfgDir, 'skills');
const mcpPath = path.join(cfgDir, 'mcp_config.json');

let prefs = {};
if (fs.existsSync(cfgPath)) {
  try { prefs = JSON.parse(fs.readFileSync(cfgPath, 'utf8')).plugins || {}; } catch (_) {}
}

const present = [];
const enabled = [];
if (fs.existsSync(pluginsDir)) {
  for (const name of fs.readdirSync(pluginsDir)) {
    if (!fs.statSync(path.join(pluginsDir, name)).isDirectory()) continue;
    if (!BANNED.includes(name)) continue;
    present.push(name);
    if (pluginEnabled(name, prefs, pluginsDir)) enabled.push(name);
  }
}

const skills = fs.existsSync(skillsDir)
  ? fs.readdirSync(skillsDir).filter((n) => fs.statSync(path.join(skillsDir, n)).isDirectory()).sort()
  : [];

let mcp = [];
if (fs.existsSync(mcpPath)) {
  try {
    const j = JSON.parse(fs.readFileSync(mcpPath, 'utf8'));
    const servers = j.mcpServers || j.servers || {};
    mcp = Object.keys(servers).sort().map((n) => ({ name: n, health: classifyMCP(n) }));
  } catch (_) {}
}

const pin = process.env.ORCHESTRA_CONTRACT || '(unset, using VERSION 3.1.0)';
const warn = enabled.length > 0;
const headroomGone = warn && skills.length >= 30;

console.log('=== Orchestra 3.1.0 Antigravity doctor ===');
console.log(`Contract pin:     ${pin}`);
console.log(`AG Global skills: ${skills.length}  ${skills.join(', ') || '(none)'}`);
if (enabled.length) {
  console.log(`AG plugins:       WARN banned Global enabled: ${enabled.join(', ')}`);
  console.log('                  Disable in Antigravity Settings > Customizations (or config.json plugins.<name>.enabled=false).');
} else if (present.length) {
  console.log(`AG plugins:       OK ${present.join(', ')} installed but disabled`);
} else {
  console.log('AG plugins:       OK science / data-agent-kit not installed as Global');
}
if (headroomGone) console.log('AG headroom:      WARN customization budget is gone');
if (!mcp.length) {
  console.log('AG MCP:           (no mcp_config.json)');
} else {
  console.log('AG MCP:');
  for (const s of mcp) console.log(`  ${s.name.padEnd(22)} ${s.health}`);
}

process.exit(warn ? 1 : 0);
