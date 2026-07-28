<#
.SYNOPSIS
    Install and verify the core host toolchain without a dev container.

.DESCRIPTION
    Installs project-managed verification tools, checks Windows prerequisites,
    and runs a smoke build plus a race-enabled test. Packaging and release
    tooling remains in .devcontainer/ and the documented release workflow.
    This script removes GitCode credential variables before launching tools.

.PARAMETER Check
    Run offline checks only. Installs nothing and makes no persistent changes.
#>
[CmdletBinding()]
param([switch]$Check)

$ErrorActionPreference = 'Stop'
Remove-Item Env:GC_TOKEN -ErrorAction SilentlyContinue
Remove-Item Env:GITCODE_TOKEN -ErrorAction SilentlyContinue

$RepoRoot = Split-Path -Parent $PSScriptRoot
Set-Location $RepoRoot

$GoMinMajor = 1
$GoMinMinor = 22
$GolangciVersion = 'v2.12.2'
$GitleaksVersion = 'v8.30.1'
$PreCommitVersion = '4.6.1'
$ToolsRoot = if ($env:GC_DEV_TOOLS_DIR) { $env:GC_DEV_TOOLS_DIR } else {
    Join-Path $env:LOCALAPPDATA 'gitcode-cli\dev-tools'
}
$ManagedBin = Join-Path $ToolsRoot 'bin'
$PreCommitVenv = Join-Path $ToolsRoot 'pre-commit'
$ManagedMarker = 'rem Managed by gitcode-cli scripts/dev-setup.ps1'
$script:Missing = @()

function Write-Section { param([string]$Name) Write-Host ''; Write-Host "[$Name]" }
function Write-Ok { param([string]$Message) Write-Host "ok   $Message" }
function Write-Info { param([string]$Message) Write-Host "  $Message" }
function Write-Gap { param([string]$Message) Write-Host "MISS $Message"; $script:Missing += $Message }

function Get-CommandPath {
    param([string]$Name)
    $managed = Join-Path $ManagedBin "$Name.exe"
    if (Test-Path $managed) { return $managed }
    $managedCmd = Join-Path $ManagedBin "$Name.cmd"
    if (Test-Path $managedCmd) { return $managedCmd }
    $command = Get-Command $Name -ErrorAction SilentlyContinue
    if ($command) { return $command.Source }
    return $null
}

function Test-CommandWorks {
    param([string]$Name, [string[]]$Arguments = @('--version'))
    $path = Get-CommandPath $Name
    if (-not $path) { return $false }
    & $path @Arguments > $null 2>&1
    return ($LASTEXITCODE -eq 0)
}

function Update-ProcessPath {
    $machine = [Environment]::GetEnvironmentVariable('Path', 'Machine')
    $user = [Environment]::GetEnvironmentVariable('Path', 'User')
    $env:Path = "$ManagedBin;$machine;$user;$env:Path"
}

function Install-WingetPackage {
    param([string]$Name, [string]$PackageId, [string[]]$VersionArguments = @('--version'))
    if (Test-CommandWorks $Name $VersionArguments) { return }
    if ($Check) { return }
    if (-not (Get-Command winget -ErrorAction SilentlyContinue)) {
        Write-Gap "$Name (winget unavailable; install $PackageId manually)"
        return
    }
    Write-Info "installing $Name with winget ($PackageId)"
    & winget install --id $PackageId --exact --source winget `
        --accept-package-agreements --accept-source-agreements --disable-interactivity
    if ($LASTEXITCODE -ne 0) { Write-Gap "$Name (winget install failed)"; return }
    Update-ProcessPath
}

function Test-GoVersion {
    if (-not (Test-CommandWorks go @('version'))) { return $false }
    $raw = (& go env GOVERSION 2>$null) -replace '^go', ''
    $parts = $raw.Split('.')
    if ($parts.Count -lt 2) { return $false }
    $major = 0; $minor = 0
    if (-not [int]::TryParse($parts[0], [ref]$major)) { return $false }
    if (-not [int]::TryParse($parts[1], [ref]$minor)) { return $false }
    return ($major -gt $GoMinMajor -or ($major -eq $GoMinMajor -and $minor -ge $GoMinMinor))
}

function Get-GoModuleVersion {
    param([string]$Binary, [string]$Module)
    if (-not $Binary -or -not (Test-Path $Binary)) { return '' }
    $lines = & go version -m $Binary 2>$null
    foreach ($line in $lines) {
        $fields = $line.Trim() -split '\s+'
        if ($fields.Count -ge 3 -and $fields[0] -eq 'mod' -and $fields[1] -eq $Module) { return $fields[2] }
    }
    return ''
}

function Ensure-GoTool {
    param([string]$Name, [string]$Module, [string]$Version)
    $binary = Get-CommandPath $Name
    $current = Get-GoModuleVersion $binary $Module
    if ($current -eq $Version) { Write-Ok "$Name $Version ($binary)"; return }
    if ($Check) { Write-Gap "$Name $Version$(if($current){" (found $current)"})"; return }
    New-Item -ItemType Directory -Force -Path $ManagedBin | Out-Null
    Write-Info "installing $Name $Version into $ManagedBin"
    $oldGoBin = $env:GOBIN; $env:GOBIN = $ManagedBin
    & go install "$Module@$Version"
    $exitCode = $LASTEXITCODE
    if ($null -eq $oldGoBin) { Remove-Item Env:GOBIN -ErrorAction SilentlyContinue } else { $env:GOBIN = $oldGoBin }
    if ($exitCode -eq 0) { Write-Ok "$Name $Version ($(Join-Path $ManagedBin "$Name.exe"))" } else { Write-Gap "$Name $Version (go install failed)" }
}

function Find-GitBash {
    foreach ($candidate in @(
        'C:\Program Files\Git\bin\bash.exe',
        'C:\Program Files (x86)\Git\bin\bash.exe',
        (Join-Path $env:LOCALAPPDATA 'Programs\Git\bin\bash.exe')
    )) { if (Test-Path $candidate) { return $candidate } }
    return $null
}

function Write-ManagedWrapper {
    param([string]$Path, [string[]]$Lines)
    if (Test-Path $Path) {
        $first = Get-Content -LiteralPath $Path -TotalCount 1
        if ($first -ne $ManagedMarker) { Write-Gap "refusing to overwrite unmanaged wrapper: $Path"; return $false }
    }
    @($ManagedMarker, '@echo off') + $Lines | Set-Content -LiteralPath $Path -Encoding ASCII
    return $true
}

function Test-ManagedWrapper {
    param([string]$Path, [string[]]$Arguments = @('--version'))
    if (-not (Test-Path $Path)) { return $false }
    if ((Get-Content -LiteralPath $Path -TotalCount 1) -ne $ManagedMarker) { return $false }
    & $Path @Arguments > $null 2>&1
    return ($LASTEXITCODE -eq 0)
}

function Find-Python {
    foreach ($name in @('python3', 'python')) {
        $commands = Get-Command $name -All -ErrorAction SilentlyContinue |
            Where-Object { $_.Source -and $_.Source -notlike "$ManagedBin*" }
        foreach ($command in $commands) {
            & $command.Source --version > $null 2>&1
            if ($LASTEXITCODE -eq 0) { return $command.Source }
        }
    }
    if (Get-Command py -ErrorAction SilentlyContinue) {
        $candidate = & py -3 -c "import sys; print(sys.executable)" 2>$null
        if ($LASTEXITCODE -eq 0 -and $candidate -and (Test-Path $candidate)) { return $candidate }
    }
    return $null
}

function Ensure-PreCommit {
    param([string]$Python)
    $binary = Get-CommandPath pre-commit
    $current = ''
    if ($binary) { $current = ((& $binary --version 2>$null) -split '\s+')[-1] }
    if ($current -eq $PreCommitVersion) { Write-Ok "pre-commit $PreCommitVersion ($binary)"; return }
    if ($Check) { Write-Gap "pre-commit $PreCommitVersion$(if($current){" (found $current)"})"; return }
    New-Item -ItemType Directory -Force -Path $ToolsRoot, $ManagedBin | Out-Null
    Write-Info "installing pre-commit $PreCommitVersion into $PreCommitVenv"
    & $Python -m venv $PreCommitVenv
    if ($LASTEXITCODE -eq 0) {
        $venvPython = Join-Path $PreCommitVenv 'Scripts\python.exe'
        & $venvPython -m pip install --disable-pip-version-check "pre-commit==$PreCommitVersion" > $null
    }
    if ($LASTEXITCODE -ne 0) { Write-Gap 'pre-commit venv install failed'; return }
    $wrapper = Join-Path $ManagedBin 'pre-commit.cmd'
    $target = Join-Path $PreCommitVenv 'Scripts\pre-commit.exe'
    if (Write-ManagedWrapper $wrapper @("`"$target`" %*")) {
        Write-Ok "pre-commit $PreCommitVersion ($wrapper)"
    }
}

if (-not $Check) {
    Write-Section 'System dependencies'
    Install-WingetPackage go GoLang.Go @('version')
    Install-WingetPackage git Git.Git
    Install-WingetPackage python Python.Python.3.12
    Install-WingetPackage make ezwinports.make
    Install-WingetPackage gcc BrechtSanders.WinLibs.MCF.UCRT
}
Update-ProcessPath

if (-not $Check -and (Test-CommandWorks go @('version')) -and -not (Test-GoVersion)) {
    Write-Info 'upgrading Go with winget'
    & winget upgrade --id GoLang.Go --exact --source winget `
        --accept-package-agreements --accept-source-agreements --disable-interactivity
    if ($LASTEXITCODE -ne 0) { Write-Gap 'Go upgrade failed' }
    Update-ProcessPath
}

Write-Section 'Go toolchain'
if (-not (Test-GoVersion)) {
    Write-Gap "Go $GoMinMajor.$GoMinMinor+ (upgrade with winget or https://go.dev/dl/)"
} else {
    Write-Ok "go $(& go env GOVERSION) ($(Get-CommandPath go))"
}
$goReady = Test-GoVersion

if ($goReady) {
    Write-Section 'Module dependencies'
    $oldProxy = $env:GOPROXY
    if ($Check) { $env:GOPROXY = 'off' }
    & go list -mod=readonly -deps ./... > $null 2>&1
    $listExit = $LASTEXITCODE
    if ($Check) { if ($null -eq $oldProxy) { Remove-Item Env:GOPROXY -ErrorAction SilentlyContinue } else { $env:GOPROXY = $oldProxy } }
    if ($listExit -eq 0) { Write-Ok 'module dependencies resolved' }
    elseif ($Check) { Write-Gap 'module dependencies (run without -Check to download)' }
    else { & go mod download; if ($LASTEXITCODE -eq 0) { Write-Ok 'module dependencies resolved' } else { Write-Gap 'go mod download failed' } }

    Write-Section 'Verification tools'
    Ensure-GoTool golangci-lint github.com/golangci/golangci-lint/v2 $GolangciVersion
    Ensure-GoTool gitleaks github.com/zricethezav/gitleaks/v8 $GitleaksVersion
}

Write-Section 'Git and shell utilities'
Test-CommandWorks git | ForEach-Object { if($_){Write-Ok "git ($(Get-CommandPath git))"}else{Write-Gap 'git'} }
$gitBash = Find-GitBash
if ($gitBash) { Write-Ok "git bash ($gitBash)" } else { Write-Gap 'git bash' }

$realMake = $null
foreach ($name in @('mingw32-make.exe', 'make.exe')) {
    $found = Get-Command $name -ErrorAction SilentlyContinue | Where-Object { $_.Source -notlike "$ManagedBin*" } | Select-Object -First 1
    if ($found) { $realMake = $found.Source; break }
}
$makeWrapper = Join-Path $ManagedBin 'make.cmd'
if ($Check) {
    if ($gitBash -and (Test-ManagedWrapper $makeWrapper)) { Write-Ok "managed make wrapper ($makeWrapper)" } else { Write-Gap 'managed make wrapper' }
} elseif ($realMake -and $gitBash) {
    New-Item -ItemType Directory -Force -Path $ManagedBin | Out-Null
    $fso = New-Object -ComObject Scripting.FileSystemObject
    $shortBash = ($fso.GetFile($gitBash).ShortPath) -replace '\\', '/'
    [void][Runtime.InteropServices.Marshal]::ReleaseComObject($fso)
    if (Write-ManagedWrapper $makeWrapper @("`"$realMake`" `"SHELL=$shortBash`" `".SHELLFLAGS=-c`" %*")) { Write-Ok "managed make wrapper ($makeWrapper)" }
} else { Write-Gap 'GNU make and Git bash' }

$python = Find-Python
if ($python) {
    Write-Ok "python ($python)"
    if (-not $Check) {
        New-Item -ItemType Directory -Force -Path $ManagedBin | Out-Null
        $pythonWrapper = Join-Path $ManagedBin 'python3.cmd'
        if (Write-ManagedWrapper $pythonWrapper @("`"$python`" %*")) { Write-Ok "managed python3 wrapper ($pythonWrapper)" }
    } elseif (Test-ManagedWrapper (Join-Path $ManagedBin 'python3.cmd')) { Write-Ok 'managed python3 wrapper' }
    else { Write-Gap 'managed python3 wrapper' }
    Ensure-PreCommit $python
} else { Write-Gap 'Python 3' }

$compiler = Get-Command gcc, clang, cc -ErrorAction SilentlyContinue | Select-Object -First 1
if ($compiler) { Write-Ok "C compiler ($($compiler.Source))" } else { Write-Gap 'C compiler (required by go test -race)' }

Write-Section 'Verification'
if ($goReady) {
    $smokeDir = Join-Path ([IO.Path]::GetTempPath()) "gc-dev-setup-$([guid]::NewGuid().ToString('N'))"
    New-Item -ItemType Directory -Force -Path $smokeDir | Out-Null
    $oldProxy = $env:GOPROXY; if ($Check) { $env:GOPROXY = 'off' }
    & go build -mod=readonly -o (Join-Path $smokeDir 'gc.exe') ./cmd/gc
    $buildExit = $LASTEXITCODE
    if ($buildExit -eq 0) { & (Join-Path $smokeDir 'gc.exe') version > $null 2>&1 }
    if ($buildExit -eq 0 -and $LASTEXITCODE -eq 0) { Write-Ok 'build and gc version' } else { Write-Gap 'build or gc version' }
    if ($compiler) {
        & go test -mod=readonly -race ./pkg/config > $null 2>&1
        if ($LASTEXITCODE -eq 0) { Write-Ok 'race-enabled test' } else { Write-Gap 'race-enabled test' }
    }
    if ($Check) { if ($null -eq $oldProxy) { Remove-Item Env:GOPROXY -ErrorAction SilentlyContinue } else { $env:GOPROXY = $oldProxy } }
    Remove-Item -LiteralPath $smokeDir -Recurse -Force -ErrorAction SilentlyContinue
}

Write-Host ''
if ($script:Missing.Count) {
    Write-Host "incomplete: $($script:Missing.Count) dependency gap(s)"
    foreach ($item in $script:Missing) { Write-Host "  - $item" }
    Write-Host "Managed tools directory: $ManagedBin"
    exit 1
}

Write-Host 'dev environment ready.'
Write-Host "Managed tools: $ManagedBin"
Write-Host 'No user PATH changes were made. Add the managed tools directory manually if desired.'
