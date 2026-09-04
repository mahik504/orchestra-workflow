# Wire Orchestra 3.1.0 onto chosen hosts. No secrets. No skills add --all.
#
#   powershell -NoProfile -ExecutionPolicy Bypass -File kit/bootstrap.ps1
#   powershell -NoProfile -ExecutionPolicy Bypass -File kit/bootstrap.ps1 -Hosts cursor,antigravity -Target D:\work\my-orchestra
param(
    [string]$Hosts = "",
    [string]$Target = ""
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$User = $env:USERPROFILE
$stackPath = Join-Path $Root "registries\host-stack.json"
if (-not (Test-Path -LiteralPath $stackPath)) { throw "missing registries/host-stack.json" }
$stack = Get-Content -LiteralPath $stackPath -Raw | ConvertFrom-Json

function Resolve-Hosts([string]$raw) {
    if ($raw) {
        return @($raw.Split(",") | ForEach-Object { $_.Trim().ToLower() } | Where-Object { $_ -ne "" })
    }
    Write-Host "Which hosts should use Orchestra 3.1.0?"
    Write-Host "  [1] Cursor"
    Write-Host "  [2] Antigravity"
    Write-Host "  [3] Claude Code"
    Write-Host "  [4] Codex / Hermes / OpenCode (AGENTS.md only)"
    $ans = Read-Host "Enter numbers (e.g. 1 2)"
    $map = @{ "1" = "cursor"; "2" = "antigravity"; "3" = "claude"; "4" = "agents" }
    $out = @()
    foreach ($tok in ($ans -split "[\s,]+")) {
        if ($map.ContainsKey($tok)) { $out += $map[$tok] }
    }
    if ($out.Count -eq 0) { throw "No hosts selected." }
    return $out
}

$chosen = Resolve-Hosts $Hosts
Write-Host ("Hosts: " + ($chosen -join ", "))

if (-not $Target) {
    if ($env:ORCHESTRA_HOME) { $Target = $env:ORCHESTRA_HOME }
    else { $Target = Join-Path (Split-Path -Parent $Root) "orchestra-workspace" }
}
$init = Join-Path $Root "kit\init-workspace.ps1"
$needsInit = -not (Test-Path -LiteralPath $Target)
if (-not $needsInit) {
    $existing = @(Get-ChildItem -LiteralPath $Target -Force -ErrorAction SilentlyContinue)
    if ($existing.Count -eq 0) { $needsInit = $true }
}
if ($needsInit) {
    & $init -Target $Target
} else {
    Write-Host "workspace already exists, leaving $Target untouched"
}

$skillDest = @{
    cursor      = Join-Path $User ".cursor\skills"
    antigravity = Join-Path $User ".gemini\config\skills"
    claude      = Join-Path $User ".claude\skills"
    agents      = Join-Path $User ".agents\skills"
}

foreach ($h in $chosen) {
    if (-not $skillDest.ContainsKey($h)) { continue }
    $d = $skillDest[$h]
    if (-not (Test-Path -LiteralPath $d)) {
        New-Item -ItemType Directory -Path $d | Out-Null
        Write-Host "created $d"
    }
    foreach ($n in $stack.skills) {
        $src = Join-Path $Root "skills\$n"
        if (-not (Test-Path -LiteralPath $src)) { throw "missing skill $n" }
        $dst = Join-Path $d $n
        if (Test-Path -LiteralPath $dst) { Remove-Item -LiteralPath $dst -Recurse -Force }
        Copy-Item -LiteralPath $src -Destination $dst -Recurse
    }
    Write-Host ("copied {0} skills -> {1}" -f $stack.skills.Count, $d)
}

$adapters = Join-Path $Target "kit"
if (-not (Test-Path -LiteralPath $adapters)) { New-Item -ItemType Directory -Path $adapters | Out-Null }

if ($chosen -contains "cursor") {
    Copy-Item (Join-Path $Root ".cursorrules") (Join-Path $Target ".cursorrules") -Force
    Write-Host "wrote workspace .cursorrules"
}
if ($chosen -contains "claude") {
    Copy-Item (Join-Path $Root "CLAUDE.md") (Join-Path $Target "CLAUDE.md") -Force
    Write-Host "wrote workspace CLAUDE.md"
}
if ($chosen -contains "antigravity") {
    $ag = Join-Path $Target "kit\antigravity"
    if (-not (Test-Path -LiteralPath $ag)) { New-Item -ItemType Directory -Path $ag | Out-Null }
    Copy-Item (Join-Path $Root "kit\antigravity\MASTER-PROMPT.md") (Join-Path $ag "MASTER-PROMPT.md") -Force
    Copy-Item (Join-Path $Root "kit\antigravity\ALWAYS-ON.md") (Join-Path $ag "ALWAYS-ON.md") -Force
    Copy-Item (Join-Path $Root "kit\antigravity\mcp_config.example.json") (Join-Path $ag "mcp_config.example.json") -Force
    $live = Join-Path $User ".gemini\mcp_config.json"
    $exDest = Join-Path $User ".gemini\mcp_config.example.json"
    $geminiDir = Join-Path $User ".gemini"
    if (-not (Test-Path -LiteralPath $geminiDir)) { New-Item -ItemType Directory -Path $geminiDir | Out-Null }
    if (Test-Path -LiteralPath $live) {
        Write-Host "left existing $live untouched (may contain keys)"
    } else {
        Copy-Item (Join-Path $Root "kit\antigravity\mcp_config.example.json") $exDest -Force
        Write-Host "wrote $exDest (template only; paste keys in the host UI)"
    }
}

Copy-Item (Join-Path $Root "mcp_config.example.json") (Join-Path $Target "mcp_config.example.json") -Force
Copy-Item (Join-Path $Root "AGENTS.md") (Join-Path $Target "AGENTS.md") -Force

Write-Host ""
Write-Host "Marketplace plugins (install yourself; we cannot log into Google, Stripe, or Stitch for you):"
if ($chosen -contains "cursor") {
    Write-Host "  Cursor:"
    foreach ($p in $stack.plugins.cursor) { Write-Host "    - $p" }
}
if ($chosen -contains "antigravity") {
    Write-Host "  Antigravity:"
    foreach ($p in $stack.plugins.antigravity) { Write-Host "    - $p" }
}

Write-Host ""
Write-Host "Restart the IDE. Next chat in this workspace uses Orchestra 3.1.0 (AGENTS.md)."
Write-Host "Optional Cursor User Rule for chats outside this folder: GLOBAL until I say skip orchestra."
Write-Host "We cannot log into Google, Stripe, or Stitch for you."
Write-Host "Private workspace: $Target"
