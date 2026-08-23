# Install Orchestra skills + merge Antigravity MCP (no secrets printed).
# Run: powershell -ExecutionPolicy Bypass -File install.ps1
# Optional: -VaultPath "D:\my-brain"

param(
  [string]$VaultPath = "<ORCHESTRA_HOME>"
)

$ErrorActionPreference = "Stop"
$SkillsGlobal = Join-Path $env:USERPROFILE ".gemini\config\skills"
$AgentsSkills = Join-Path $env:USERPROFILE ".agents\skills"
$CursorSkills = Join-Path $env:USERPROFILE ".cursor\skills"
$McpPath = Join-Path $env:USERPROFILE ".gemini\config\mcp_config.json"
$KitRoot = $PSScriptRoot

New-Item -ItemType Directory -Force -Path $SkillsGlobal | Out-Null
New-Item -ItemType Directory -Force -Path (Split-Path $McpPath) | Out-Null

function Copy-Skill([string]$Name) {
  $src = Join-Path $CursorSkills $Name
  if (-not (Test-Path $src)) { $src = Join-Path $KitRoot "skills\$Name" }
  if (-not (Test-Path $src)) { Write-Host "skip missing $Name"; return }
  $dest = Join-Path $SkillsGlobal $Name
  if (Test-Path $dest) { Remove-Item -LiteralPath $dest -Recurse -Force }
  Copy-Item -LiteralPath $src -Destination $dest -Recurse
  if (Test-Path $AgentsSkills) {
    $d2 = Join-Path $AgentsSkills $Name
    if (Test-Path $d2) { Remove-Item -LiteralPath $d2 -Recurse -Force }
    Copy-Item -LiteralPath $src -Destination $d2 -Recurse
  }
  Write-Host "copied $Name"
}

foreach ($n in @("orchestra-conductor","orchestra-vault","orchestra-ship","orchestra-docs","ship-safe")) {
  Copy-Skill $n
}

function Add-Skills([string]$Repo, [string[]]$SkillNames) {
  $skillArgs = @()
  foreach ($s in $SkillNames) { $skillArgs += @("--skill", $s) }
  & npx --yes skills@latest add $Repo -g -a antigravity -y --copy @skillArgs
}

Write-Host "Installing Expo (curated)..."
Add-Skills "expo/skills" @(
  "expo-router","expo-project-structure","expo-data-fetching",
  "expo-native-ui","expo-upgrade","expo-dev-client","eas-app-stores"
)

Write-Host "Installing Stitch (curated)..."
Add-Skills "google-labs-code/stitch-skills" @(
  "stitch::generate-design","stitch::manage-design-system","stitch::extract-design-md",
  "stitch::extract-static-html","stitch::code-to-design","stitch::upload-to-stitch",
  "stitch::react-components","taste-design"
)

Write-Host "Installing Impeccable..."
& npx --yes skills@latest add pbakaus/impeccable -g -a antigravity -y --copy

Write-Host "Installing Emil motion..."
Add-Skills "emilkowalski/skills" @("animate","review-animations")

Write-Host "Installing Strix (own apps only)..."
Add-Skills "usestrix/strix" @(
  "penetration-testing-with-strix",
  "fix-security-vulnerabilities-with-strix",
  "ci-security-scanning-with-strix"
)

# Never mirror all of ~/.agents/skills — that folder still holds old dumps.
$allow = @(
  "orchestra-conductor","orchestra-vault","orchestra-ship","orchestra-docs","ship-safe",
  "eas-app-stores","expo-data-fetching","expo-dev-client","expo-native-ui","expo-project-structure","expo-router","expo-upgrade",
  "impeccable","animate","review-animations","emil-design-eng","taste-design",
  "ci-security-scanning-with-strix","fix-security-vulnerabilities-with-strix","penetration-testing-with-strix",
  "stitch-code-to-design","stitch-extract-design-md","stitch-extract-static-html","stitch-generate-design",
  "stitch-manage-design-system","stitch-react-components","stitch-upload-to-stitch",
  "generate-design","manage-design-system","extract-design-md","extract-static-html",
  "code-to-design","upload-to-stitch","react-components"
)
if (Test-Path $AgentsSkills) {
  Get-ChildItem $AgentsSkills -Directory | Where-Object { $allow -contains $_.Name } | ForEach-Object {
    $dest = Join-Path $SkillsGlobal $_.Name
    if (-not (Test-Path $dest)) {
      Copy-Item -LiteralPath $_.FullName -Destination $dest -Recurse
      Write-Host "mirrored $($_.Name) -> gemini config skills"
    }
  }
}

$py = Join-Path $KitRoot "merge_mcp.py"
python $py --mcp $McpPath --vault $VaultPath
Write-Host "Done. Restart Antigravity. Do not commit mcp_config.json."
