# Sync the 30 canonical skills from this clone onto local hosts.
#
#   powershell -NoProfile -ExecutionPolicy Bypass -File kit/sync-ides.ps1
#
# Does not copy quarantined libraries or host extras. Does not print MCP secrets.

$ErrorActionPreference = "Stop"
$User = $env:USERPROFILE
$Root = Split-Path -Parent $PSScriptRoot
$PublicSkills = Join-Path $Root "skills"

$Allow = @(
  "animate", "ci-security-scanning-with-strix", "eas-app-stores", "emil-design-eng",
  "expo-data-fetching", "expo-dev-client", "expo-native-ui", "expo-project-structure",
  "expo-router", "expo-upgrade", "fix-security-vulnerabilities-with-strix", "impeccable",
  "orchestra-conductor", "orchestra-docs", "orchestra-ship", "orchestra-vault",
  "penetration-testing-with-strix", "review-animations", "semgrep-adapter", "ship-safe",
  "stitch-code-to-design", "stitch-extract-design-md", "stitch-extract-static-html",
  "stitch-generate-design", "stitch-manage-design-system", "stitch-react-components",
  "stitch-upload-to-stitch", "superpowers-planning", "taste-design", "web-design-guidelines"
)

$dests = @(
  (Join-Path $User ".cursor\skills"),
  (Join-Path $User ".gemini\config\skills"),
  (Join-Path $User ".agents\skills"),
  (Join-Path $User ".claude\skills")
)

foreach ($n in $Allow) {
  $src = Join-Path $PublicSkills $n
  if (-not (Test-Path -LiteralPath $src)) { throw "missing skill $n in $PublicSkills" }
  foreach ($d in $dests) {
    if (-not (Test-Path -LiteralPath $d)) { continue }
    $dst = Join-Path $d $n
    if (Test-Path -LiteralPath $dst) { Remove-Item -LiteralPath $dst -Recurse -Force }
    Copy-Item -LiteralPath $src -Destination $dst -Recurse
  }
}

Write-Host ("synced {0} canonical skills from {1}" -f $Allow.Count, $PublicSkills)
Write-Host "rollback: set ORCHESTRA_CONTRACT (see kit/ROLLBACK.md)"
