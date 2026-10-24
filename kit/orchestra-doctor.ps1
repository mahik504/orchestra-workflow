# Orchestra 3.3.5 doctor. Prints failures first. Never prints secret values.
param([switch]$Json)

$ErrorActionPreference = "Continue"
$Root = Split-Path -Parent $PSScriptRoot
$problems = @()
$warnings = @()

$version = (Get-Content (Join-Path $Root "VERSION") -Raw).Trim()
if ($version -ne "3.3.5") { $problems += "VERSION is $version" }
$agents = Get-Content (Join-Path $Root "AGENTS.md") -TotalCount 1
if ($agents -notmatch "3\.3\.5") { $problems += "AGENTS.md title is not 3.3.5" }

$want = Get-ChildItem (Join-Path $Root "skills") -Directory | Where-Object { Test-Path (Join-Path $_.FullName "SKILL.md") } | ForEach-Object Name
$localRoot = Join-Path $Root "skills-local"
if (Test-Path $localRoot) {
  $want += Get-ChildItem $localRoot -Directory | Where-Object { Test-Path (Join-Path $_.FullName "SKILL.md") } | ForEach-Object Name
}
$cursor = Join-Path $env:USERPROFILE ".cursor\skills"
if (Test-Path $cursor) {
  $have = Get-ChildItem $cursor -Directory | ForEach-Object Name
  $missing = @($want | Where-Object { $have -notcontains $_ })
  if ($missing.Count -gt 0) { $problems += ("cursor missing " + ($missing -join ", ")) }
}

$desc = 0
Get-ChildItem (Join-Path $Root "skills") -Recurse -Filter SKILL.md | ForEach-Object {
  $desc += (Get-Content $_.FullName -TotalCount 8 | Out-String).Length
}
if ($desc -gt 120000) { $warnings += "skill descriptions are large ($desc chars)" }

$vaultName = "orchestra-" + "brain"
$vault = Join-Path "C:\projects" $vaultName
$note = Join-Path $vault "protocols"
if (Test-Path $note) {
  $extra = @(Get-ChildItem $note -File | Where-Object { $_.Name -ne "NOTE.md" })
  if ($extra.Count -gt 0) { $problems += "vault protocols has files other than NOTE.md" }
}

try {
  $task = Get-ScheduledTask -TaskName "OrchestraBrainVaultSync" -ErrorAction Stop
  $hasLogon = $false
  $has16 = $false
  $has12 = $false
  foreach ($tr in $task.Triggers) {
    if ($tr.CimClass.CimClassName -eq "MSFT_TaskLogonTrigger") { $hasLogon = $true }
    $interval = [string]$tr.Repetition.Interval
    if ($interval -eq "PT16H") { $has16 = $true }
    if ($interval -eq "PT12H") { $has12 = $true }
  }
  if (-not $hasLogon) { $problems += "heartbeat has no logon trigger" }
  if (-not $has16) { $problems += "heartbeat has no 16 hour trigger" }
  if ($has12) { $problems += "heartbeat still has a 12 hour trigger" }
} catch {
  $problems += "heartbeat task missing"
}

$patterns = @("li_at=", "auth_token=", "ct0=", "sessionid=")
foreach ($tree in @($Root, $vault)) {
  if (-not (Test-Path $tree)) { continue }
  Get-ChildItem $tree -Recurse -File -Include *.md,*.json,*.ps1,*.txt -ErrorAction SilentlyContinue |
    Where-Object { $_.FullName -notmatch '\\node_modules\\|\\\.git\\|orchestra-doctor\.ps1|3\.3\.5-' } |
    ForEach-Object {
      foreach ($pat in $patterns) {
        if (Select-String -Path $_.FullName -Pattern $pat -Quiet -ErrorAction SilentlyContinue) {
          $problems += ("secret pattern in " + $_.Name)
        }
      }
    }
}

Write-Host "problems: $($problems.Count)"
$problems | ForEach-Object { Write-Host "FAIL $_" }
Write-Host "warnings: $($warnings.Count)"
$warnings | ForEach-Object { Write-Host "WARN $_" }
if ($Json) {
  @{ problems = $problems; warnings = $warnings } | ConvertTo-Json | Write-Host
}
if ($problems.Count -gt 0) { exit 1 }
exit 0
