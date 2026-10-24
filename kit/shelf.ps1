# Clone one shelf repo. Does not run its install scripts.
param(
  [Parameter(Mandatory = $true)][string]$Url,
  [Parameter(Mandatory = $true)][string]$Id
)

$ErrorActionPreference = "Stop"
if ($Url -notmatch '^https://github\.com/') { throw "shelf clone accepts a GitHub https URL only" }
$destRoot = Join-Path "C:\projects" "_shelf"
New-Item -ItemType Directory -Force $destRoot | Out-Null
$dest = Join-Path $destRoot $Id
if (Test-Path $dest) { Write-Host "already cloned $dest"; exit 0 }
git clone --depth 1 $Url $dest
Write-Host "cloned $dest"
Write-Host "Read the README and license before you run anything."
