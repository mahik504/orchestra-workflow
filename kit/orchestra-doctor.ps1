# Orchestra doctor (no Go required).
# Prints Antigravity Global skill names, banned-plugin status, and MCP health.
# Never prints keys. Exit 1 if science or data-agent-kit are Global-enabled.
#
#   powershell -NoProfile -ExecutionPolicy Bypass -File kit\orchestra-doctor.ps1

$ErrorActionPreference = "Stop"
$here = Split-Path -Parent $MyInvocation.MyCommand.Path
$candidates = @(
  (Join-Path $here "..\runtime\tools\doctor-ag.js"),
  (Join-Path $here "..\..\orchestra-workflow\runtime\tools\doctor-ag.js")
)
$script = $candidates | Where-Object { Test-Path $_ } | Select-Object -First 1
if (-not $script) { throw "doctor-ag.js not found. Clone orchestra-workflow next to this kit, or run from the public repo." }
node $script
exit $LASTEXITCODE
