# Install Orchestra skills into Antigravity/Gemini folders. No secrets printed.
# Run from this folder: powershell -ExecutionPolicy Bypass -File install.ps1
# Optional: -WorkspacePath "D:\my-orchestra-workspace"

param(
  [string]$WorkspacePath = ""
)

$ErrorActionPreference = "Stop"
$SkillsGlobal = Join-Path $env:USERPROFILE ".gemini\config\skills"
$AgentsSkills = Join-Path $env:USERPROFILE ".agents\skills"
$CursorSkills = Join-Path $env:USERPROFILE ".cursor\skills"
$McpPath = Join-Path $env:USERPROFILE ".gemini\config\mcp_config.json"
$KitRoot = $PSScriptRoot
$RepoRoot = Split-Path (Split-Path $KitRoot -Parent) -Parent
if (-not $WorkspacePath) {
  $WorkspacePath = Join-Path (Split-Path $RepoRoot -Parent) "orchestra-workspace"
}

New-Item -ItemType Directory -Force -Path $SkillsGlobal | Out-Null
New-Item -ItemType Directory -Force -Path (Split-Path $McpPath) | Out-Null

function Copy-Skill([string]$Name) {
  $src = Join-Path $CursorSkills $Name
  if (-not (Test-Path $src)) { $src = Join-Path $KitRoot "skills\$Name" }
  if (-not (Test-Path $src)) { $src = Join-Path $RepoRoot "skills\$Name" }
  if (-not (Test-Path $src)) { Write-Host "skip missing $Name"; return }
  $dest = Join-Path $SkillsGlobal $Name
  if (Test-Path $dest) { Remove-Item -LiteralPath $dest -Recurse -Force }
  Copy-Item -LiteralPath $src -Destination $dest -Recurse
  if (Test-Path $AgentsSkills) {
    $d2 = Join-Path $AgentsSkills $Name
    if (Test-Path $d2) { Remove-Item -LiteralPath $d2 -Recurse -Force }
    Copy-Item -LiteralPath $src -Destination $d2 -Recurse
  }
  Write-Host "copied $Name"
}

foreach ($n in @("orchestra-conductor", "orchestra-vault", "orchestra-ship", "orchestra-docs", "ship-safe")) {
  Copy-Skill $n
}

Write-Host "Curated extras are optional. This script does not run skills add --all."
Write-Host "Workspace default: $WorkspacePath"
Write-Host "Put API keys in the Antigravity MCP UI. Do not commit mcp_config.json."
Write-Host "Restart Antigravity."
