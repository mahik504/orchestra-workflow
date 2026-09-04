# Copy Orchestra skills into local agent skill folders IF those folders already exist.
# Does not npx skills add --all. Does not create MCP configs with keys.
$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$Skills = Join-Path $Root "skills"
if (-not (Test-Path -LiteralPath $Skills)) { throw "missing skills/" }

$dests = @(
  (Join-Path $env:USERPROFILE ".cursor\skills"),
  (Join-Path $env:USERPROFILE ".claude\skills"),
  (Join-Path $env:USERPROFILE ".agents\skills"),
  (Join-Path $env:USERPROFILE ".gemini\config\skills")
)

Get-ChildItem -LiteralPath $Skills -Directory | ForEach-Object {
  $name = $_.Name
  $src = $_.FullName
  foreach ($d in $dests) {
    if (-not (Test-Path -LiteralPath $d)) { continue }
    $out = Join-Path $d $name
    if (-not (Test-Path -LiteralPath $out)) { New-Item -ItemType Directory -Path $out | Out-Null }
    Copy-Item -LiteralPath (Join-Path $src "*") -Destination $out -Recurse -Force
    Write-Host "installed $name -> $out"
  }
}

Write-Host "Done. Restart the agent. If a dest folder was missing, that product is not installed - skip it."
