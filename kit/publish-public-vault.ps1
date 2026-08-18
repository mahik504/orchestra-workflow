# Publish allowlisted files to the public orchestra template.
# Does not change visibility of mahik504/orchestra-brain (that remote stays private).
# Does not copy career.md, unpublished idea folders, local-notes, logs, or secrets.
# Does not force-push. Does not skip hooks. Does not change git config.
#
# Staging clone: C:\projects\orchestra-workflow
# Remote:        https://github.com/mahik504/orchestra-workflow.git
#
# Run after the private vault commit:
#   powershell -NoProfile -ExecutionPolicy Bypass -File C:\projects\orchestra-brain\kit\publish-public-vault.ps1

$ErrorActionPreference = "Stop"
$Vault = "C:\projects\orchestra-brain"
$PublicRoot = "C:\projects\orchestra-workflow"
$PublicRemote = "https://github.com/mahik504/orchestra-workflow.git"

function Copy-Rel([string]$RelPath) {
  $src = Join-Path $Vault $RelPath
  $dst = Join-Path $PublicRoot $RelPath
  if (-not (Test-Path -LiteralPath $src)) {
    throw "missing allowlisted source: $RelPath"
  }
  $parent = Split-Path -Parent $dst
  if ($parent -and -not (Test-Path -LiteralPath $parent)) {
    New-Item -ItemType Directory -Path $parent | Out-Null
  }
  Copy-Item -LiteralPath $src -Destination $dst -Force
}

function Copy-Tree([string]$RelDir, [string[]]$ExcludeNames) {
  $src = Join-Path $Vault $RelDir
  $dst = Join-Path $PublicRoot $RelDir
  if (-not (Test-Path -LiteralPath $src)) {
    throw "missing allowlisted directory: $RelDir"
  }
  if (Test-Path -LiteralPath $dst) {
    Remove-Item -LiteralPath $dst -Recurse -Force
  }
  Copy-Item -LiteralPath $src -Destination $dst -Recurse -Force
  foreach ($name in $ExcludeNames) {
    Get-ChildItem -LiteralPath $dst -Recurse -Force |
      Where-Object { $_.Name -eq $name } |
      ForEach-Object { Remove-Item -LiteralPath $_.FullName -Recurse -Force }
  }
}

if (-not (Test-Path -LiteralPath $PublicRoot)) {
  New-Item -ItemType Directory -Path $PublicRoot | Out-Null
}

# Wipe previous export except .git so we never leave a withdrawn file tracked.
Get-ChildItem -LiteralPath $PublicRoot -Force |
  Where-Object { $_.Name -ne ".git" } |
  ForEach-Object { Remove-Item -LiteralPath $_.FullName -Recurse -Force }

$files = @(
  ".gitignore",
  "README.md",
  "START HERE.md",
  "WORKFLOW.md",
  "Preferences.md",
  "CONNECT.md",
  "STACK.md",
  "STITCH.md",
  "kit\sync-vault.ps1",
  "kit\publish-public-vault.ps1",
  "memory\README.md",
  "memory\decisions.md",
  "projects\airlens\idea.md",
  "projects\astroverse\idea.md"
)

foreach ($f in $files) {
  Copy-Rel $f
}

Copy-Tree "templates" @()
Copy-Tree "kit\antigravity" @("mcp_config.json")
Copy-Tree ".obsidian" @("workspace.json", "workspace-mobile.json", "cache")

# Guard: these must never appear in the public tree
$forbidden = @(
  "memory\career.md",
  "memory\local-notes.md",
  "projects\yumit",
  "projects\odyss",
  "projects\portfolio-penfight",
  "mcp_config.json"
)
foreach ($rel in $forbidden) {
  $hit = Join-Path $PublicRoot $rel
  if (Test-Path -LiteralPath $hit) {
    throw "refusing to publish forbidden path: $rel"
  }
}

$scan = Get-ChildItem -LiteralPath $PublicRoot -Recurse -File |
  Where-Object { $_.Extension -in ".md", ".ps1", ".json", ".py", ".txt", ".gitignore" }
foreach ($file in $scan) {
  $text = Get-Content -LiteralPath $file.FullName -Raw -ErrorAction SilentlyContinue
  if ($null -eq $text) { continue }
  if ($text -match "ghp_[A-Za-z0-9]+" -or $text -match "github_pat_[A-Za-z0-9_]+" -or $text -match "sk-[A-Za-z0-9]{10,}") {
    throw "refusing to publish possible secret in $($file.FullName)"
  }
}

Set-Location -LiteralPath $PublicRoot

if (-not (Test-Path -LiteralPath (Join-Path $PublicRoot ".git"))) {
  git init -b main
  git remote add origin $PublicRemote
} else {
  $existing = git remote get-url origin 2>$null
  if (-not $existing) {
    git remote add origin $PublicRemote
  }
}

git add -A
$porcelain = git status --porcelain
if ($porcelain) {
  git commit -m "Publish public orchestra workflow template without personal career or unpublished ideas."
} else {
  Write-Host "public tree unchanged"
}

git push -u origin main
if ($LASTEXITCODE -ne 0) {
  Write-Host "push failed. Create the public repo first, then re-run this script:"
  Write-Host "  gh repo create mahik504/orchestra-workflow --public --description `"Public orchestra workflow template (no personal career or unpublished ideas).`" --source `"$PublicRoot`" --remote origin --push"
  exit 1
}

Write-Host "published $PublicRemote"
exit 0
