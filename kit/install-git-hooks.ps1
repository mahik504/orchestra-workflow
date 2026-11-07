# Point this repo at kit/git-hooks, or copy the commit-msg hook into another checkout.
# Repo-local only. Does not change global git config or the commit author.
param(
  [string]$Repo = ""
)

$ErrorActionPreference = "Stop"
$Kit = $PSScriptRoot
$HookSrc = Join-Path $Kit "git-hooks\commit-msg"
if (-not (Test-Path -LiteralPath $HookSrc)) {
  throw "Missing $HookSrc"
}

$Orchestra = Split-Path -Parent $Kit
if ($Repo -eq "") {
  $Repo = $Orchestra
}
$Repo = (Resolve-Path -LiteralPath $Repo).Path
if (-not (Test-Path -LiteralPath (Join-Path $Repo ".git"))) {
  throw "$Repo is not a git checkout"
}

$OrchestraResolved = (Resolve-Path -LiteralPath $Orchestra).Path
if ($Repo -eq $OrchestraResolved) {
  # Relative to the repo root so the hook travels with the checkout.
  git -C $Repo config --local core.hooksPath "kit/git-hooks"
  Write-Output "core.hooksPath=kit/git-hooks in $Repo"
  exit 0
}

$destDir = Join-Path $Repo ".git\hooks"
New-Item -ItemType Directory -Force -Path $destDir | Out-Null
Copy-Item -LiteralPath $HookSrc -Destination (Join-Path $destDir "commit-msg") -Force
Write-Output "copied commit-msg into $destDir"
