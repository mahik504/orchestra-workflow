// Cross-platform hook. Cursor and Claude Code both pass a JSON payload on stdin.
//   node kit/hooks/guard.mjs <shell|read|mcp|edit|session|stop|subagent-start|subagent-stop>
//
// ORCHESTRA_HOOK_PROFILE=off|minimal|standard|strict
// ORCHESTRA_DISABLED_HOOKS=shell,read
import fs from 'node:fs';
import path from 'node:path';

const id = process.argv[2] || 'shell';
const profile = process.env.ORCHESTRA_HOOK_PROFILE || 'standard';
const disabled = new Set(
  (process.env.ORCHESTRA_DISABLED_HOOKS || '').split(',').map((s) => s.trim()).filter(Boolean),
);

function emit(obj, code = 0) {
  process.stdout.write(JSON.stringify(obj));
  process.exit(code);
}

function allow() {
  emit({ permission: 'allow' });
}

function decide(permission, message) {
  if (profile === 'minimal' && permission === 'ask') permission = 'allow';
  if (profile === 'strict' && permission === 'ask') permission = 'deny';
  emit({
    permission,
    user_message: message,
    agent_message: message,
  });
}

if (profile === 'off' || disabled.has(id)) allow();

const raw = fs.readFileSync(0, 'utf8');
const input = raw.trim() ? JSON.parse(raw) : {};
const command = String(input.command || input.tool_input?.command || '');
const filePath = String(
  input.file_path || input.path || input.filePath || input.tool_input?.file_path || input.tool_input?.path || '',
);

if (id === 'shell') {
  const deny = shellDeny(command);
  if (deny) decide('deny', deny);
  const ask = shellAsk(command);
  if (ask) decide('ask', ask);
  allow();
}

if (id === 'read') {
  const why = readDeny(filePath);
  if (why) decide('deny', why);
  allow();
}

if (id === 'mcp') {
  const server = String(input.server || input.mcp_server || input.tool_input?.server || '');
  if (server && !mcpAllowed(server)) {
    decide('deny', `MCP server ${server} is not in kit/mcp-core.json`);
  }
  allow();
}

if (id === 'edit') {
  const note = editNote(filePath);
  emit({ permission: 'allow', additional_context: note || '' });
}

if (id === 'session') {
  emit({ additional_context: sessionNote() });
}

if (id === 'stop') {
  const text = JSON.stringify(input);
  const claims = /\b(DONE|FIXED|VERIFIED|PASSED|SHIPPED)\b/.test(text);
  const evidence = /exit code|\$ |npm test|go test|node --test|EXIT /.test(text);
  if (claims && !evidence) {
    const message = 'Completion claim has no command output in the hook payload.';
    if (profile === 'strict') emit({ followup_message: message, permission: 'deny' }, 2);
    emit({ followup_message: message, permission: 'allow' });
  }
  allow();
}

if (id === 'subagent-start' || id === 'subagent-stop') {
  const dir = path.join(process.cwd(), '.orchestra', 'handoff');
  fs.mkdirSync(dir, { recursive: true });
  const record = {
    hook: id,
    at: new Date().toISOString(),
    files: input.files || input.owned_files || [],
  };
  fs.writeFileSync(path.join(dir, 'subagent.json'), JSON.stringify(record, null, 2));
  if (id === 'subagent-stop' && profile === 'strict' && input.wrote_outside) {
    decide('deny', 'subagent wrote outside its file list');
  }
  allow();
}

allow();

function shellDeny(cmd) {
  if (!cmd) return '';
  if (/npx\s+skills\s+add\s+--all\b/.test(cmd) || /\bskills\s+add\s+-g\s+-y\b/.test(cmd)) {
    return 'skills add --all and global yes-installs are refused';
  }
  if (/\bgit\s+push\b[^\n]*--force\b/.test(cmd) || /\bgit\s+push\b[^\n]*(^|\s)-f(\s|$)/.test(cmd)) {
    return 'force-push is refused';
  }
  if (/(?:curl|wget)\b[^|\n]*\|\s*(?:sh|bash)\b/i.test(cmd)) {
    return 'curl piped to a shell is refused';
  }
  if (/\b(cat|type|Get-Content)\b[^\n]*\.env\b/i.test(cmd)) {
    return 'reading .env from the shell is refused';
  }
  if (isRootDelete(cmd)) return 'deleting a drive root or home directory is refused';
  return '';
}

function shellAsk(cmd) {
  if (!cmd) return '';
  if (/\b(npm|pnpm|yarn)\s+(install|i|add)\s+-g\b/.test(cmd)) return 'a global install needs approval';
  if (/\bgit\s+reset\s+--hard\b/.test(cmd)) return 'git reset --hard needs approval';
  if (/\bgit\s+push\b/.test(cmd)) return 'git push needs approval';
  return '';
}

function isRootDelete(cmd) {
  if (!/\brm\s+-[^\n]*r/.test(cmd) && !/Remove-Item\b[^\n]*-Recurse/i.test(cmd)) return false;
  if (/\b(node_modules|dist|\.next)\b/.test(cmd)) return false;
  return /(?:\s|^)(?:\/|~|\/\*|C:\\|C:\/|\$HOME|%USERPROFILE%)(?:\s|$)/.test(cmd);
}

function readDeny(file) {
  const norm = file.replace(/\\/g, '/');
  if (!norm) return '';
  if (/(^|\/)\.env($|\.)/.test(norm)) return 'reading a .env file is refused';
  if (/id_rsa|\.pem$|credentials\.json$/i.test(norm)) return 'reading key material is refused';
  if (/mcp[^/]*\.json$/i.test(norm) && !/mcp_config\.example|kit\/mcp/.test(norm)) {
    return 'reading a live MCP config is refused';
  }
  return '';
}

function mcpAllowed(server) {
  const file = path.join(process.cwd(), 'kit', 'mcp-core.json');
  if (!fs.existsSync(file)) return true;
  const text = fs.readFileSync(file, 'utf8');
  return text.includes(`"${server}"`);
}

function editNote(file) {
  if (!file || !fs.existsSync(file)) return '';
  const buf = fs.readFileSync(file);
  if (buf.length >= 3 && buf[0] === 0xef && buf[1] === 0xbb && buf[2] === 0xbf) {
    return `BOM in ${file}`;
  }
  const text = buf.toString('utf8');
  if (/sk-[A-Za-z0-9]{20,}/.test(text) || /ghp_[A-Za-z0-9]{20,}/.test(text)) {
    return `secret-shaped string in ${file}`;
  }
  if (file.endsWith('SKILL.md') && !text.startsWith('---')) {
    return `${file} is missing frontmatter`;
  }
  return '';
}

function sessionNote() {
  const state = path.join(process.cwd(), '.orchestra', 'STATE.md');
  const limits = 'Hard limits: no skills add --all, no secrets in git, no second conductor, evidence before a completion claim.';
  if (!fs.existsSync(state)) return limits.slice(0, 1000);
  const body = fs.readFileSync(state, 'utf8').slice(0, 700);
  return `${limits}\n${body}`.slice(0, 1000);
}
