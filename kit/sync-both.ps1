# 12h: private vault backup, then public template if allowlisted files changed.
# Does not empty-commit the public repo. Does not force-push. Does not skip hooks.
# Does not change git config. Secrets stay out of git.
#
# Task Scheduler (replace the old private-only task):
#   schtasks /Create /TN "OrchestraBrainVaultSync" /SC HOURLY /MO 12 /F /TR "powershell.exe -NoProfile -ExecutionPolicy Bypass -File C:\projects\orchestra-brain\kit\sync-both.ps1"

$ErrorActionPreference = "Stop"
$Vault = "C:\projects\orchestra-brain"
$Private = Join-Path $Vault "kit\sync-vault.ps1"
$Public = Join-Path $Vault "kit\publish-public-vault.ps1"
$Log = Join-Path $Vault "kit\sync-vault.log"

function Write-Log([string]$Message) {
  $line = "{0} {1}" -f (Get-Date -Format "yyyy-MM-dd HH:mm:ss"), $Message
  Add-Content -LiteralPath $Log -Value $line
}

& $Private
Write-Log "private sync exit $LASTEXITCODE"
if ($LASTEXITCODE -ne 0) {
  Write-Log "skip public publish; private sync failed"
  exit $LASTEXITCODE
}

& $Public
Write-Log "public publish exit $LASTEXITCODE"
exit $LASTEXITCODE
