# Creates a PRIVATE Orchestra workspace (empty brain). Does not copy any personal vault.
param(
  [string]$Target = ""
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
if (-not $Target) {
  if ($env:ORCHESTRA_HOME) {
    $Target = $env:ORCHESTRA_HOME
  } else {
    $Target = Join-Path (Split-Path -Parent $Root) "orchestra-workspace"
  }
}

$template = Join-Path $Root "workspace-template"
if (-not (Test-Path -LiteralPath $template)) { throw "missing workspace-template" }

if (-not (Test-Path -LiteralPath $Target)) {
  New-Item -ItemType Directory -Path $Target | Out-Null
}

Copy-Item -LiteralPath (Join-Path $template "*") -Destination $Target -Recurse -Force

foreach ($dir in @("protocols", "registries", "templates", "docs")) {
  $src = Join-Path $Root $dir
  if (Test-Path -LiteralPath $src) {
    $dst = Join-Path $Target $dir
    if (Test-Path -LiteralPath $dst) { Remove-Item -LiteralPath $dst -Recurse -Force }
    Copy-Item -LiteralPath $src -Destination $dst -Recurse -Force
  }
}

foreach ($f in @("AGENTS.md", "Preferences.md", "WORKFLOW.md", "START.md")) {
  $src = Join-Path $Root $f
  if (Test-Path -LiteralPath $src) {
    Copy-Item -LiteralPath $src -Destination (Join-Path $Target $f) -Force
  }
}

$pref = Join-Path $Target "Preferences.md"
if (Test-Path -LiteralPath $pref) {
  # already template defaults
}

Write-Host "Private workspace: $Target"
Write-Host "Point your agent at that folder + your app repo. Do not commit secrets."
Write-Host "Optional private git: cd `"$Target`"; git init; git remote add origin <your-private-repo>"
