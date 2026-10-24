# Enable or disable one optional MCP server in a host config.
# Does not print secret values. A missing key is left empty and reported as AUTH_REQUIRED.
param(
  [Parameter(Mandatory = $true)][string]$Name,
  [ValidateSet("enable", "disable")][string]$Action = "enable",
  [string]$Repo = ""
)

$ErrorActionPreference = "Stop"
Write-Host "mcp $Action $Name"
if ($Repo) { Write-Host "repo scope: $Repo" }
Write-Host "This script records intent. Live host edits happen in kit/sync-mcp.ps1."
if ($Action -eq "enable" -and $Name -in @("github", "supabase", "vercel", "firecrawl", "stitch")) {
  Write-Host "AUTH_REQUIRED until the host env has the key. No key is written here."
}
