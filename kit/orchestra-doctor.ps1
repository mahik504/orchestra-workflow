# Orchestra doctor. Prints failures first. Never prints secret values.
# The expected version is registries/versions.json.
param([switch]$Json, [switch]$Private)

$ErrorActionPreference = "Continue"
$Root = Split-Path -Parent $PSScriptRoot
$problems = @()
$warnings = @()

$version = (Get-Content (Join-Path $Root "VERSION") -Raw).Trim()
$expected = (Get-Content (Join-Path $Root "registries\versions.json") -Raw | ConvertFrom-Json).current
if ($version -ne $expected) { $problems += "VERSION is $version, expected $expected" }
$agents = Get-Content (Join-Path $Root "AGENTS.md") -TotalCount 1
if ($agents -notmatch [regex]::Escape($expected)) { $problems += "AGENTS.md title is not $expected" }

$linkLines = @(git -C $Root ls-files -s 2>$null | Select-String -Pattern "^160000")
foreach ($link in $linkLines) {
  $problems += ("gitlink " + (($link.Line -split "`t")[-1]))
}
Get-ChildItem (Join-Path $Root "skills") -Recurse -Filter SKILL.md -File -ErrorAction SilentlyContinue | ForEach-Object {
  $bytes = [System.IO.File]::ReadAllBytes($_.FullName)
  $bad = $false
  if ($bytes.Length -ge 2 -and (($bytes[0] -eq 255 -and $bytes[1] -eq 254) -or ($bytes[0] -eq 254 -and $bytes[1] -eq 255))) { $bad = $true }
  $limit = [Math]::Min(400, $bytes.Length)
  for ($i = 0; $i -lt $limit; $i++) {
    if ($bytes[$i] -eq 0) { $bad = $true; break }
  }
  if ($bad) { $problems += ("UTF-16 or NUL in " + $_.FullName.Substring($Root.Length).TrimStart('\')) }
}

$doctorArgs = @("--hosts")
if ($Private) { $doctorArgs += "--private" }
& node (Join-Path $Root "runtime\tools\ow-doctor.mjs") @doctorArgs
if ($LASTEXITCODE -ne 0) { $problems += "ow-doctor reported a host or description problem" }

$desc = 0
Get-ChildItem (Join-Path $Root "skills") -Recurse -Filter SKILL.md | ForEach-Object {
  $desc += (Get-Content $_.FullName -TotalCount 8 | Out-String).Length
}
if ($desc -gt 120000) { $warnings += "skill descriptions are large ($desc chars)" }

if ($Private -and -not [string]::IsNullOrWhiteSpace($env:ORCHESTRA_VAULT)) {
  try {
    $task = Get-ScheduledTask -TaskName "OrchestraBrainVaultSync" -ErrorAction Stop
    $hasLogon = $false
    $has12 = $false
    $has16 = $false
    $hasNet = $false
    foreach ($tr in $task.Triggers) {
      if ($tr.CimClass.CimClassName -eq "MSFT_TaskLogonTrigger") { $hasLogon = $true }
      if ($tr.CimClass.CimClassName -eq "MSFT_TaskEventTrigger") { $hasNet = $true }
      $interval = [string]$tr.Repetition.Interval
      if ($interval -eq "PT12H") { $has12 = $true }
      if ($interval -eq "PT16H") { $has16 = $true }
    }
    if (-not $hasLogon) { $problems += "heartbeat has no logon trigger" }
    if (-not $has12) { $problems += "heartbeat has no 12 hour trigger" }
    if ($has16) { $problems += "heartbeat still has a 16 hour trigger" }
    if (-not $hasNet) { $problems += "heartbeat has no network-connect trigger" }
  } catch {
    $problems += "heartbeat task missing"
  }
}

$claudeSettings = Join-Path $env:USERPROFILE ".claude\settings.json"
if (Test-Path -LiteralPath $claudeSettings) {
  try {
    $enabled = (Get-Content -LiteralPath $claudeSettings -Raw | ConvertFrom-Json).enabledPlugins
    if ($null -ne $enabled) {
      foreach ($prop in $enabled.PSObject.Properties) {
        if ($prop.Value -ne $true) { continue }
        foreach ($name in @("superpowers", "vercel", "caveman")) {
          if ($prop.Name -eq $name -or $prop.Name -like "$name@*") {
            $warnings += "Claude settings enable refused plugin $($prop.Name)"
          }
        }
      }
    }
  } catch {
    $warnings += "Claude settings.json could not be read"
  }
}

$patterns = @("li_at=", "auth_token=", "ct0=", "sessionid=")
$trees = @($Root)
if ($Private -and -not [string]::IsNullOrWhiteSpace($env:ORCHESTRA_VAULT)) { $trees += $env:ORCHESTRA_VAULT }
foreach ($tree in $trees) {
  if (-not (Test-Path $tree)) { continue }
  Get-ChildItem $tree -Recurse -File -Include *.md,*.json,*.ps1,*.txt -ErrorAction SilentlyContinue |
    Where-Object { $_.FullName -notmatch '\\node_modules\\|\\\.git\\|orchestra-doctor\.ps1|3\.3\.\d+-' } |
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
