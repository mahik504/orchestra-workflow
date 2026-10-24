# Merge kit/mcp-core.json into host MCP configs.
# Keeps an existing server's env block. Never prints secret values.
# Does not remove servers unless -Remove is passed.
param(
  [switch]$Remove
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$corePath = Join-Path $Root "kit\mcp-core.json"
$core = Get-Content $corePath -Raw | ConvertFrom-Json
$vault = $env:ORCHESTRA_VAULT
if (-not $vault) { throw "Set ORCHESTRA_VAULT to the private notes folder before syncing." }
$vaultId = "orchestra-" + "brain"

$files = @(
  (Join-Path $env:USERPROFILE ".cursor\mcp.json"),
  (Join-Path $env:USERPROFILE ".gemini\config\mcp_config.json"),
  (Join-Path $env:USERPROFILE ".claude\settings.json")
)

$aliases = @{ stitch = @("stitch", "StitchMCP"); vault = @("vault", $vaultId) }

foreach ($file in $files) {
  if (-not (Test-Path $file)) { Write-Host "missing $file"; continue }
  $json = Get-Content $file -Raw | ConvertFrom-Json
  if (-not $json.mcpServers) {
    $json | Add-Member -NotePropertyName mcpServers -NotePropertyValue ([pscustomobject]@{})
  }
  $servers = $json.mcpServers
  foreach ($prop in $core.mcpServers.PSObject.Properties) {
    $id = $prop.Name
    $names = if ($aliases.ContainsKey($id)) { $aliases[$id] } else { @($id) }
    $existing = $names | Where-Object { $servers.PSObject.Properties.Name -contains $_ } | Select-Object -First 1
    if ($existing) {
      Write-Host "$file already has $existing"
      continue
    }
    $add = $prop.Value | ConvertTo-Json -Depth 8 | ConvertFrom-Json
    if ($id -eq "vault") {
      $add.args[$add.args.Count - 1] = $vault
      $id = $vaultId
    }
    $servers | Add-Member -NotePropertyName $id -NotePropertyValue $add
    Write-Host "$file added $id"
  }
  if ($Remove) {
    foreach ($id in @("shadcn", "magicuidesign-mcp", "github", "firebase-mcp-server")) {
      if ($servers.PSObject.Properties.Name -contains $id) {
        $servers.PSObject.Properties.Remove($id)
        Write-Host "$file removed $id"
      }
    }
  }
  $json | ConvertTo-Json -Depth 12 | Set-Content -Path $file -Encoding utf8
}
Write-Host "done"
