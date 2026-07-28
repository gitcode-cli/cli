# Verify that Windows check mode is offline and side-effect free.
[CmdletBinding()]
param()

$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent $PSScriptRoot
$scriptPath = Join-Path $PSScriptRoot 'dev-setup.ps1'
$tempRoot = Join-Path ([IO.Path]::GetTempPath()) "gc-dev-setup-test-$([guid]::NewGuid().ToString('N'))"
$userPathBefore = [Environment]::GetEnvironmentVariable('Path', 'User')

$oldTools = $env:GC_DEV_TOOLS_DIR
$oldToken = $env:GC_TOKEN
$oldLegacyToken = $env:GITCODE_TOKEN
$oldProxy = $env:GOPROXY

try {
    $env:GC_DEV_TOOLS_DIR = $tempRoot
    $env:GC_TOKEN = 'must-not-appear'
    $env:GITCODE_TOKEN = 'must-not-appear'
    $env:GOPROXY = 'http://127.0.0.1:1'

    $output = & powershell -NoProfile -ExecutionPolicy Bypass -File $scriptPath -Check 2>&1 | Out-String
    $status = $LASTEXITCODE

    if ($status -notin @(0, 1)) { throw "unexpected check-mode exit code: $status`n$output" }
    if (Test-Path $tempRoot) { throw "check mode created managed tools directory: $tempRoot" }
    if ([Environment]::GetEnvironmentVariable('Path', 'User') -ne $userPathBefore) {
        throw 'check mode changed the user PATH'
    }
    if ($output -match 'must-not-appear') { throw 'check mode printed a credential sentinel' }

    "Windows dev setup check-mode test passed (exit=$status)."
} finally {
    if ($null -eq $oldTools) { Remove-Item Env:GC_DEV_TOOLS_DIR -ErrorAction SilentlyContinue } else { $env:GC_DEV_TOOLS_DIR = $oldTools }
    if ($null -eq $oldToken) { Remove-Item Env:GC_TOKEN -ErrorAction SilentlyContinue } else { $env:GC_TOKEN = $oldToken }
    if ($null -eq $oldLegacyToken) { Remove-Item Env:GITCODE_TOKEN -ErrorAction SilentlyContinue } else { $env:GITCODE_TOKEN = $oldLegacyToken }
    if ($null -eq $oldProxy) { Remove-Item Env:GOPROXY -ErrorAction SilentlyContinue } else { $env:GOPROXY = $oldProxy }
    if (Test-Path $tempRoot) { Remove-Item -LiteralPath $tempRoot -Recurse -Force }
    Set-Location $repoRoot
}
