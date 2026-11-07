# Sync canonical skills onto local hosts through runtime/tools/ow-sync.mjs.
# Dry-run unless -Sync. A skill you edited locally is left in place.
# -Prune is refused. User-installed skills stay where they are.
#
#   powershell -NoProfile -ExecutionPolicy Bypass -File kit/sync-ides.ps1
#   powershell -NoProfile -ExecutionPolicy Bypass -File kit/sync-ides.ps1 -Sync

param(
  [switch]$Sync,
  [switch]$Prune,
  [switch]$BackupAndReplace
)

$ErrorActionPreference = "Stop"
if ($Prune) {
  Write-Host "Prune is refused. User-installed skills stay on the host."
  exit 2
}

$Root = Split-Path -Parent $PSScriptRoot
$args = @()
if ($Sync) { $args += "--apply" }
if ($BackupAndReplace) { $args += "--backup-and-replace" }
& node (Join-Path $Root "runtime\tools\ow-sync.mjs") @args
exit $LASTEXITCODE
