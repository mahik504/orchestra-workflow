# Persist ORCHESTRA_HOME and install orchestra.exe onto the user PATH.
# Does not hardcode a vault path. Pass -HomeDir (or set ORCHESTRA_HOME first).
# Does not npm install -g. Does not write secrets.
param(
    [string]$HomeDir = $env:ORCHESTRA_HOME
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$Runtime = Join-Path $Root "runtime"

if (-not $HomeDir) {
    throw "Set ORCHESTRA_HOME or pass -HomeDir to your private workspace root."
}
if (-not (Test-Path -LiteralPath $HomeDir)) {
    throw "HomeDir does not exist: $HomeDir"
}
$cmdDir = Join-Path $Runtime "cmd\orchestra"
if (-not (Test-Path -LiteralPath $cmdDir)) {
    throw "missing runtime/cmd/orchestra - run this from an orchestra-workflow clone"
}

[Environment]::SetEnvironmentVariable("ORCHESTRA_HOME", $HomeDir, "User")
$env:ORCHESTRA_HOME = $HomeDir
Write-Host "ORCHESTRA_HOME (User) = $HomeDir"

$GoBin = Join-Path $env:USERPROFILE "go\bin"
if (-not (Test-Path -LiteralPath $GoBin)) {
    New-Item -ItemType Directory -Path $GoBin | Out-Null
}
$env:GOBIN = $GoBin

Push-Location $Runtime
try {
    go install ./cmd/orchestra
} finally {
    Pop-Location
}

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (-not $userPath) { $userPath = "" }
$parts = @()
if ($userPath) { $parts = $userPath.Split(";") | Where-Object { $_ -ne "" } }
if ($parts -notcontains $GoBin) {
    if ($userPath) {
        $newPath = $userPath + ";" + $GoBin
    } else {
        $newPath = $GoBin
    }
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "Added $GoBin to User PATH (new terminals pick this up)"
} else {
    Write-Host "User PATH already contains $GoBin"
}

$exe = Join-Path $GoBin "orchestra.exe"
if (-not (Test-Path -LiteralPath $exe)) {
    throw "go install finished but $exe is missing"
}
Write-Host "installed $exe"
Write-Host "Restart Cursor / Antigravity so they see ORCHESTRA_HOME and PATH."
Write-Host "Then: orchestra doctor"
