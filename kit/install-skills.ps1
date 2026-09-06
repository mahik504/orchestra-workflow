# Copy Orchestra skills listed in host-stack.json into local agent skill folders IF those folders already exist.
# Does not npx skills add --all. Does not create MCP configs with keys.
$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$Skills = Join-Path $Root "skills"
$stackPath = Join-Path $Root "registries\host-stack.json"
if (-not (Test-Path -LiteralPath $Skills)) { throw "missing skills/" }
if (-not (Test-Path -LiteralPath $stackPath)) { throw "missing registries/host-stack.json" }
$stack = Get-Content -LiteralPath $stackPath -Raw | ConvertFrom-Json

$dests = @(
  (Join-Path $env:USERPROFILE ".cursor\skills"),
  (Join-Path $env:USERPROFILE ".claude\skills"),
  (Join-Path $env:USERPROFILE ".agents\skills"),
  (Join-Path $env:USERPROFILE ".gemini\config\skills"),
  (Join-Path $env:USERPROFILE ".jcode\skills")
)

foreach ($name in $stack.skills) {
  $src = Join-Path $Skills $name
  if (-not (Test-Path -LiteralPath $src)) { throw "missing skill $name" }
  foreach ($d in $dests) {
    if (-not (Test-Path -LiteralPath $d)) { continue }
    $out = Join-Path $d $name
    if (Test-Path -LiteralPath $out) { Remove-Item -LiteralPath $out -Recurse -Force }
    Copy-Item -LiteralPath $src -Destination $out -Recurse
    Write-Host "installed $name -> $out"
  }
}

Write-Host ("Done. Orchestra {0}. Restart the agent. Missing dest folders are skipped." -f $stack.version)
