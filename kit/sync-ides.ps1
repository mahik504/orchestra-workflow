# Sync canonical skills from this clone onto local hosts.
#
#   powershell -NoProfile -ExecutionPolicy Bypass -File kit/sync-ides.ps1
#   powershell -NoProfile -ExecutionPolicy Bypass -File kit/sync-ides.ps1 -Sync
#   powershell -NoProfile -ExecutionPolicy Bypass -File kit/sync-ides.ps1 -Sync -Prune
#
# With no switch, this only prints host folders that are not in the skill set.
# -Sync copies skills/ and skills-local/ onto the hosts.
# -Prune moves extra folders to the backup dir. It does not delete them.
# Paseo app skills are kept. They are a host app, not a second conductor.

param(
  [switch]$Sync,
  [switch]$Prune
)

$ErrorActionPreference = "Stop"
$User = $env:USERPROFILE
$Root = Split-Path -Parent $PSScriptRoot
$PublicSkills = Join-Path $Root "skills"
$LocalSkills = Join-Path $Root "skills-local"
$Backup = "C:\backups\orchestra-pre-3.3.5\pruned"

$names = @()
$names += Get-ChildItem $PublicSkills -Directory | Where-Object { Test-Path (Join-Path $_.FullName "SKILL.md") } | ForEach-Object Name
if (Test-Path $LocalSkills) {
  $names += Get-ChildItem $LocalSkills -Directory | Where-Object { Test-Path (Join-Path $_.FullName "SKILL.md") } | ForEach-Object Name
}
$keep = [System.Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)
foreach ($n in $names) { [void]$keep.Add($n) }
foreach ($n in @("paseo","paseo-advisor","paseo-committee","paseo-handoff","paseo-help","paseo-plugin",".system")) {
  [void]$keep.Add($n)
}

$dests = @(
  (Join-Path $User ".cursor\skills"),
  (Join-Path $User ".gemini\config\skills"),
  (Join-Path $User ".agents\skills"),
  (Join-Path $User ".claude\skills"),
  (Join-Path $User ".jcode\skills"),
  (Join-Path $User ".codex\skills")
)

Write-Host ("skill sources: {0}" -f $keep.Count)
Write-Host "---- extras (not removed unless -Prune) ----"
$extraCount = 0
foreach ($d in $dests) {
  if (-not (Test-Path $d)) { Write-Host "missing host dir $d"; continue }
  Get-ChildItem $d -Directory | ForEach-Object {
    if (-not $keep.Contains($_.Name)) {
      $extraCount++
      $label = $d.Substring($User.Length).TrimStart('\')
      Write-Host ("{0}\{1}" -f $label, $_.Name)
      if ($Prune) {
        $bucket = Join-Path $Backup ($label -replace '[\\:]','_')
        New-Item -ItemType Directory -Force $bucket | Out-Null
        $target = Join-Path $bucket $_.Name
        if (Test-Path $target) { Remove-Item $target -Recurse -Force }
        Move-Item $_.FullName $target
      }
    }
  }
}
Write-Host ("extras: {0}" -f $extraCount)

if (-not $Sync) {
  Write-Host "report only. Pass -Sync to copy. Pass -Prune to move extras."
  exit 0
}

foreach ($d in $dests) {
  if (-not (Test-Path $d)) { New-Item -ItemType Directory -Force $d | Out-Null }
  foreach ($srcRoot in @($PublicSkills, $LocalSkills)) {
    if (-not (Test-Path $srcRoot)) { continue }
    Get-ChildItem $srcRoot -Directory | Where-Object { Test-Path (Join-Path $_.FullName "SKILL.md") } | ForEach-Object {
      $dst = Join-Path $d $_.Name
      if (Test-Path $dst) { Remove-Item $dst -Recurse -Force }
      Copy-Item $_.FullName $dst -Recurse
    }
  }
}
Write-Host "synced skills and skills-local onto hosts"
