# Generate the vault library from registries/resources.json.
# Set ORCHESTRA_VAULT to the private notes folder. This script does not invent rows.
param()

$ErrorActionPreference = "Stop"
$vault = $env:ORCHESTRA_VAULT
if (-not $vault) { throw "Set ORCHESTRA_VAULT before generating the library." }
$Root = Split-Path -Parent $PSScriptRoot
$src = Join-Path $Root "registries\resources.json"
if (-not (Test-Path $src)) { throw "missing $src" }
$lib = Join-Path $vault "library"
New-Item -ItemType Directory -Force $lib | Out-Null
Write-Host "library target: $lib"
Write-Host "source rows: $((Get-Content $src -Raw | ConvertFrom-Json).Count)"
Write-Host "P6 writes one note per row. This script only checks that the source parses."
