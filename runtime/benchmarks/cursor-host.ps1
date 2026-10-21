$ErrorActionPreference = "Stop"
$agent = $env:ORCHESTRA_BENCH_AGENT
$workspace = $env:ORCHESTRA_BENCH_WORKSPACE
$promptFile = $env:ORCHESTRA_BENCH_PROMPT
$model = $env:ORCHESTRA_BENCH_MODEL_LABEL
if ([string]::IsNullOrWhiteSpace($agent) -or [string]::IsNullOrWhiteSpace($workspace) -or [string]::IsNullOrWhiteSpace($promptFile) -or [string]::IsNullOrWhiteSpace($model)) {
  Write-Error "bench host env is incomplete"
  exit 2
}
if (-not (Test-Path -LiteralPath $agent)) { exit 2 }
if (-not (Test-Path -LiteralPath $promptFile)) { exit 2 }
$pointer = "Read prompt.txt in the workspace root and complete the task it describes, including its acceptance lines."
$help = & $agent --help 2>&1 | Out-String
$argv = New-Object System.Collections.Generic.List[string]
$argv.Add("--print")
$argv.Add("--workspace")
$argv.Add($workspace)
$argv.Add("--model")
$argv.Add($model)
if ($help -match '(^|[\s,])--force(?=[\s,]|$)') { $argv.Add("--force") }
if ($help -match '(^|[\s,])--trust(?=[\s,]|$)') { $argv.Add("--trust") }
$argv.Add("--")
$argv.Add($pointer)
& $agent @argv
exit $LASTEXITCODE
