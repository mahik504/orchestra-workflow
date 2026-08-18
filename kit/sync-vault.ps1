# Layer-2-lite backup for C:\projects\orchestra-brain
# Commits local changes. Pulls --rebase and pushes only if a git remote exists.
# Does not skip hooks. Does not force-push. Does not change git config.
# Obsidian Git only runs while Obsidian is open; schedule this with Task Scheduler instead.
#
# Register (once, in an elevated or same-user PowerShell):
#   schtasks /Create /TN "OrchestraBrainVaultSync" /SC HOURLY /MO 12 /F /TR "powershell.exe -NoProfile -ExecutionPolicy Bypass -File C:\projects\orchestra-brain\kit\sync-vault.ps1"
#
# Secrets must stay out of git (.gitignore). Never commit mcp_config.json, .env, or keys.

$ErrorActionPreference = "Stop"
$Vault = "C:\projects\orchestra-brain"
$LogDir = Join-Path $Vault "kit"
$Log = Join-Path $LogDir "sync-vault.log"

function Write-Log([string]$Message) {
  $line = "{0} {1}" -f (Get-Date -Format "yyyy-MM-dd HH:mm:ss"), $Message
  Add-Content -LiteralPath $Log -Value $line
}

if (-not (Test-Path (Join-Path $Vault ".git"))) {
  Write-Log "skip: vault is not a git repo"
  exit 0
}

Set-Location -LiteralPath $Vault
Write-Log "start"

git add -A
$porcelain = git status --porcelain
if ($porcelain) {
  $stamp = Get-Date -Format "yyyy-MM-dd HH:mm"
  git commit -m "Vault auto-sync $stamp"
  Write-Log "committed"
} else {
  Write-Log "no local changes"
}

$remote = git remote
if (-not $remote) {
  Write-Log "no remote; local commit only"
  exit 0
}

git pull --rebase
if ($LASTEXITCODE -ne 0) {
  Write-Log "pull --rebase failed; not pushing"
  exit 1
}

git push
if ($LASTEXITCODE -ne 0) {
  Write-Log "push failed"
  exit 1
}

Write-Log "pushed"
exit 0
