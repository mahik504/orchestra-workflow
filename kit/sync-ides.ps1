# Sync canonical skills from this clone onto local hosts (Cursor, Antigravity, Claude, Agents, jcode).
#
#   powershell -NoProfile -ExecutionPolicy Bypass -File kit/sync-ides.ps1
#
# Allowlist is registries/host-stack.json. Does not copy quarantined libraries.
# Does not print MCP secrets.

$ErrorActionPreference = "Stop"
$User = $env:USERPROFILE
$Root = Split-Path -Parent $PSScriptRoot
$PublicSkills = Join-Path $Root "skills"
$stackPath = Join-Path $Root "registries\host-stack.json"
if (-not (Test-Path -LiteralPath $stackPath)) { throw "missing registries/host-stack.json" }
$stack = Get-Content -LiteralPath $stackPath -Raw | ConvertFrom-Json
$Allow = @($stack.skills)

$dests = @(
  (Join-Path $User ".cursor\skills"),
  (Join-Path $User ".gemini\config\skills"),
  (Join-Path $User ".agents\skills"),
  (Join-Path $User ".claude\skills"),
  (Join-Path $User ".jcode\skills")
)

foreach ($n in $Allow) {
  $src = Join-Path $PublicSkills $n
  if (-not (Test-Path -LiteralPath $src)) { throw "missing skill $n in $PublicSkills" }
  foreach ($d in $dests) {
    if (-not (Test-Path -LiteralPath $d)) { continue }
    $dst = Join-Path $d $n
    if (Test-Path -LiteralPath $dst) { Remove-Item -LiteralPath $dst -Recurse -Force }
    Copy-Item -LiteralPath $src -Destination $dst -Recurse
  }
}

Write-Host ("synced {0} canonical skills (Orchestra {1}) from {2}" -f $Allow.Count, $stack.version, $PublicSkills)
Write-Host "rollback: set ORCHESTRA_CONTRACT (see kit/ROLLBACK.md)"
