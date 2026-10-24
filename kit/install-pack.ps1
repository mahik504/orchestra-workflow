# Copy a pack's skills into one product repo. Never into global host folders.
param(
  [Parameter(Mandatory = $true)][string]$Pack,
  [Parameter(Mandatory = $true)][string]$Repo
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$packFile = Join-Path $Root "templates\packs\$Pack.md"
if (-not (Test-Path $packFile)) { throw "missing pack $packFile" }
if (-not (Test-Path $Repo)) { throw "missing repo $Repo" }

$members = @()
Get-Content $packFile | ForEach-Object {
  if ($_ -match '^- ([a-z0-9-]+)\s*$') { $members += $Matches[1] }
}
$dest = Join-Path $Repo ".cursor\skills"
New-Item -ItemType Directory -Force $dest | Out-Null
foreach ($name in $members) {
  $src = Join-Path $Root "skills\$name"
  if (-not (Test-Path $src)) { Write-Host "skip $name (not a global skill)"; continue }
  $dst = Join-Path $dest $name
  if (Test-Path $dst) { Remove-Item $dst -Recurse -Force }
  Copy-Item $src $dst -Recurse
  Write-Host "installed $name into $dst"
}
