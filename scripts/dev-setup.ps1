<#
.SYNOPSIS
    Install and verify local development dependencies on Windows.

.DESCRIPTION
    Windows counterpart to scripts/dev-setup.sh. The repository ships a
    devcontainer for Linux, but Windows contributors and AI agents run local
    verification directly on the host. This script installs the core
    verification baseline without a container; -WithPackaging adds the
    container's packaging toolchain.

    It also fixes three Windows-specific traps that the Makefile depends on:
      1. GNU make defaults to cmd.exe as SHELL, which breaks shell recipes.
      2. `python3` resolves to a WindowsApps stub that exits with code 9009.
      3. `$(shell date ...)` in the Makefile needs Git's usr\bin on PATH.

    Tooling only. This script never reads or prints GC_TOKEN / GITCODE_TOKEN.

.PARAMETER Check
    Verify only. Installs nothing and exits non-zero when gaps remain.

.PARAMETER WithPackaging
    Also install nfpm, goreleaser, and the Python build toolchain.

.EXAMPLE
    powershell -ExecutionPolicy Bypass -File scripts\dev-setup.ps1

.EXAMPLE
    powershell -ExecutionPolicy Bypass -File scripts\dev-setup.ps1 -Check
#>
[CmdletBinding()]
param(
    [switch]$Check,
    [switch]$WithPackaging
)

$ErrorActionPreference = 'Stop'

$RepoRoot = Split-Path -Parent $PSScriptRoot
Set-Location $RepoRoot

$GoMinMajor = 1
$GoMinMinor = 22
$NfpmVersion = 'v2.40.0'
$GolangciVersion = 'v2.12.2'
$GitleaksVersion = 'v8.30.1'
$GoreleaserVersion = 'v2.17.1'
$PreCommitVersion = '4.6.1'
$ToolsBin = Join-Path $env:USERPROFILE 'tools\bin'

$script:Missing = @()

function Write-Section { param([string]$Name) Write-Host ""; Write-Host "[$Name]" }
function Write-Ok { param([string]$Msg) Write-Host "ok   $Msg" }
function Write-Info { param([string]$Msg) Write-Host "  $Msg" }
function Write-Gap {
    param([string]$Msg)
    Write-Host "MISS $Msg"
    $script:Missing += $Msg
}

function Test-Tool { param([string]$Name) [bool](Get-Command $Name -ErrorAction SilentlyContinue) }
function Get-ToolPath { param([string]$Name) (Get-Command $Name -ErrorAction SilentlyContinue).Source }

function Test-ToolWorks {
    param([string]$Name)

    $cmd = Get-Command $Name -ErrorAction SilentlyContinue
    if (-not $cmd) { return $false }

    if ($Name -eq 'go') {
        & $cmd.Source version > $null 2>&1
    } else {
        & $cmd.Source --version > $null 2>&1
    }
    return ($LASTEXITCODE -eq 0)
}

function Update-ProcessPath {
    $machine = [Environment]::GetEnvironmentVariable('Path', 'Machine')
    $user = [Environment]::GetEnvironmentVariable('Path', 'User')
    $env:Path = "$machine;$user"
}

function Install-WingetPackage {
    param([string]$Tool, [string]$PackageId)

    if (Test-ToolWorks $Tool) { return }
    if ($Check) { return }
    if (-not (Test-Tool 'winget')) {
        Write-Gap "$Tool (winget unavailable; install $PackageId manually)"
        return
    }

    Write-Info "installing $Tool with winget ($PackageId)"
    & winget install --id $PackageId --exact --source winget `
        --accept-package-agreements --accept-source-agreements --disable-interactivity
    if ($LASTEXITCODE -ne 0) {
        Write-Gap "$Tool (winget install failed for $PackageId)"
        return
    }
    Update-ProcessPath
}

# Persist a directory onto the user PATH and make it effective in this session.
function Add-UserPath {
    param([string]$Directory, [switch]$Prepend)

    $current = [Environment]::GetEnvironmentVariable('Path', 'User')
    if (-not $current) { $current = '' }
    $parts = $current.Split(';') | Where-Object { $_ -ne '' -and $_ -ne $Directory }

    if ($Prepend) { $parts = @($Directory) + $parts } else { $parts = $parts + @($Directory) }

    [Environment]::SetEnvironmentVariable('Path', ($parts -join ';'), 'User')
    if ($Prepend) { $env:Path = "$Directory;$env:Path" } else { $env:Path = "$env:Path;$Directory" }
}

function Get-GoBinDir {
    $gobin = & go env GOBIN 2>$null
    if ($gobin) { return $gobin }
    return (Join-Path (& go env GOPATH) 'bin')
}

function Test-GoVersion {
    $raw = (& go env GOVERSION 2>$null) -replace '^go', ''
    if (-not $raw) { return $false }
    $bits = $raw.Split('.')
    if ($bits.Count -lt 2) { return $false }
    $major = [int]$bits[0]
    $minor = [int]$bits[1]
    if ($major -gt $GoMinMajor) { return $true }
    return ($major -eq $GoMinMajor -and $minor -ge $GoMinMinor)
}

function Install-GoTool {
    param([string]$Name, [string]$Module)

    if (Test-Tool $Name) {
        Write-Ok "$Name ($(Get-ToolPath $Name))"
        return
    }
    if ($Check) { Write-Gap $Name; return }

    Write-Info "installing $Name from $Module"
    & go install $Module
    if ($LASTEXITCODE -ne 0) {
        Write-Gap "$Name (go install failed)"
        return
    }

    $goBin = Get-GoBinDir
    if ($env:Path -notlike "*$goBin*") { Add-UserPath $goBin }
    if (Test-Tool $Name) {
        Write-Ok "$Name ($(Get-ToolPath $Name))"
    } else {
        Write-Gap "$Name (installed but not on PATH; add $goBin)"
    }
}

if (-not $Check) {
    Write-Section 'System dependencies'
    Install-WingetPackage 'go' 'GoLang.Go'
    Install-WingetPackage 'git' 'Git.Git'
    Install-WingetPackage 'python' 'Python.Python.3.12'
    Install-WingetPackage 'make' 'ezwinports.make'
}

Write-Section 'Go toolchain'
if (-not (Test-ToolWorks 'go')) {
    Write-Gap "go (install Go $GoMinMajor.$GoMinMinor+ from https://go.dev/dl/, or: winget install GoLang.Go)"
} elseif (-not (Test-GoVersion)) {
    Write-Gap "go $(& go env GOVERSION) is older than $GoMinMajor.$GoMinMinor"
} else {
    Write-Ok "go $(& go env GOVERSION) ($(Get-ToolPath 'go'))"
    $goBin = Get-GoBinDir
    if ($env:Path -notlike "*$goBin*") { Add-UserPath $goBin }
}

$goReady = (Test-ToolWorks 'go') -and (Test-GoVersion)

if ($goReady) {
    Write-Section 'Module dependencies'
    if ($Check) {
        $oldProxy = $env:GOPROXY
        $env:GOPROXY = 'off'
        & go list -mod=readonly -deps ./... > $null 2>&1
        if ($null -eq $oldProxy) { Remove-Item Env:GOPROXY -ErrorAction SilentlyContinue } else { $env:GOPROXY = $oldProxy }
        if ($LASTEXITCODE -eq 0) {
            Write-Ok 'module dependencies resolved'
        } else {
            Write-Gap 'module dependencies (run: go mod download)'
        }
    } else {
        Write-Info 'go mod download'
        & go mod download
        if ($LASTEXITCODE -eq 0) {
            Write-Ok 'module dependencies resolved'
        } else {
            Write-Gap 'go mod download failed'
        }
    }

    Write-Section 'Lint tooling'
    Install-GoTool 'golangci-lint' "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$GolangciVersion"

    Write-Section 'Secret scanning'
    Install-GoTool 'gitleaks' "github.com/zricethezav/gitleaks/v8@$GitleaksVersion"

    if ($WithPackaging) {
        Write-Section 'Packaging tooling'
        Install-GoTool 'nfpm' "github.com/goreleaser/nfpm/v2/cmd/nfpm@$NfpmVersion"
        Install-GoTool 'goreleaser' "github.com/goreleaser/goreleaser/v2@$GoreleaserVersion"
    }
}

Write-Section 'Git and shell utilities'
if (Test-ToolWorks 'git') {
    Write-Ok "git ($(Get-ToolPath 'git'))"
} else {
    Write-Gap 'git (winget install Git.Git)'
}

# Makefile recipes and scripts/*.sh need a real POSIX shell plus coreutils.
# Git for Windows ships both in usr\bin (bash, date, rm, sed, awk).
$gitBash = $null
foreach ($candidate in @(
    'C:\Program Files\Git\bin\bash.exe',
    'C:\Program Files (x86)\Git\bin\bash.exe',
    (Join-Path $env:LOCALAPPDATA 'Programs\Git\bin\bash.exe')
)) {
    if (Test-Path $candidate) { $gitBash = $candidate; break }
}

if ($gitBash) {
    Write-Ok "git bash ($gitBash)"
    $gitUsrBin = Join-Path (Split-Path -Parent (Split-Path -Parent $gitBash)) 'usr\bin'
    if (Test-Path (Join-Path $gitUsrBin 'date.exe')) {
        if ($env:Path -like "*$gitUsrBin*") {
            Write-Ok 'coreutils on PATH (Makefile date/rm/sed)'
        } elseif ($Check) {
            Write-Gap "coreutils not on PATH (add $gitUsrBin)"
        } else {
            # Appended, not prepended: Git's find/sort must not shadow the
            # Windows builtins that other tooling expects.
            Add-UserPath $gitUsrBin
            Write-Ok "coreutils added to PATH ($gitUsrBin)"
        }
    } else {
        Write-Gap "coreutils missing under $gitUsrBin"
    }
} else {
    Write-Gap 'git bash (required by Makefile recipes and scripts/*.sh)'
}

Write-Section 'GNU make'
$makeWrapper = Join-Path $ToolsBin 'make.cmd'
if (Test-Tool 'make') {
    & make --version > $null 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Ok "make ($(Get-ToolPath 'make'))"
        if ((Get-ToolPath 'make') -notlike '*\make.cmd') {
            Write-Info 'note: bare GNU make uses cmd.exe as SHELL; shell-based targets may fail.'
            Write-Info "      pass SHELL='C:/PROGRA~1/Git/bin/bash.exe' or use the wrapper at $makeWrapper"
        }
    } else {
        Write-Gap 'make is on PATH but fails to run (missing runtime DLLs?)'
    }
} elseif ($Check) {
    Write-Gap 'make (winget install ezwinports.make, or use go/scripts directly)'
} else {
    Write-Gap 'make was installed but is not available on PATH'
}

# GNU make on Windows spawns cmd.exe, which cannot run this repo's shell
# recipes. The wrapper forces bash via an 8.3 short path (make chokes on the
# space in "Program Files").
if ($gitBash -and -not $Check) {
    $fso = New-Object -ComObject Scripting.FileSystemObject
    $shortBash = ($fso.GetFile($gitBash).ShortPath) -replace '\\', '/'
    [void][System.Runtime.InteropServices.Marshal]::ReleaseComObject($fso)

    $realMake = $null
    foreach ($name in @('mingw32-make.exe', 'make.exe')) {
        $found = Get-Command $name -ErrorAction SilentlyContinue |
            Where-Object { $_.Source -notlike "$ToolsBin*" } |
            Select-Object -First 1
        if ($found) { $realMake = $found.Source; break }
    }

    if ($realMake) {
        New-Item -ItemType Directory -Force -Path $ToolsBin | Out-Null
        @(
            '@echo off',
            'rem Generated by scripts/dev-setup.ps1',
            'rem GNU make defaults to cmd.exe on Windows, which breaks this',
            'rem repository''s shell-based recipes. Force Git bash as SHELL.',
            "`"$realMake`" `"SHELL=$shortBash`" `".SHELLFLAGS=-c`" %*"
        ) | Set-Content -LiteralPath $makeWrapper -Encoding ASCII

        Add-UserPath $ToolsBin -Prepend
        Write-Ok "make wrapper written ($makeWrapper)"
    }
}

Write-Section 'Python tooling'
$py = $null
foreach ($name in @('python3', 'python')) {
    $cmd = Get-Command $name -ErrorAction SilentlyContinue
    if (-not $cmd) { continue }
    # WindowsApps ships a python3.exe stub that exits 9009 instead of running.
    & $cmd.Source --version > $null 2>&1
    if ($LASTEXITCODE -eq 0) { $py = $cmd.Source; break }
}
if (-not $py -and (Test-Tool 'py')) {
    $candidate = & py -3 -c "import sys; print(sys.executable)" 2>$null
    if ($LASTEXITCODE -eq 0 -and $candidate -and (Test-Path $candidate)) {
        $py = $candidate
    }
}

if (-not $py) {
    Write-Gap 'python3 (winget install Python.Python.3.12)'
} else {
    Write-Ok "python ($py, $(& $py --version 2>&1))"

    # Makefile targets call `python3` by name; the WindowsApps stub hijacks it.
    $py3 = Get-Command 'python3' -ErrorAction SilentlyContinue
    $py3Broken = $true
    if ($py3) {
        & $py3.Source --version > $null 2>&1
        $py3Broken = ($LASTEXITCODE -ne 0)
    }

    if ($py3Broken) {
        if ($Check) {
            Write-Gap 'python3 resolves to a non-functional WindowsApps stub'
        } else {
            New-Item -ItemType Directory -Force -Path $ToolsBin | Out-Null
            @(
                '@echo off',
                'rem Generated by scripts/dev-setup.ps1',
                "`"$py`" %*"
            ) | Set-Content -LiteralPath (Join-Path $ToolsBin 'python3.cmd') -Encoding ASCII
            Add-UserPath $ToolsBin -Prepend
            Write-Ok "python3 shim written ($(Join-Path $ToolsBin 'python3.cmd'))"
        }
    } else {
        Write-Ok "python3 ($($py3.Source))"
    }

    $pipPackages = @("pre-commit==$PreCommitVersion")
    if ($WithPackaging) { $pipPackages += @('build', 'wheel', 'setuptools') }

    & $py -m pre_commit --version > $null 2>&1
    $preCommitMissing = ($LASTEXITCODE -ne 0)
    $buildMissing = $false
    if ($WithPackaging) {
        & $py -m build --version > $null 2>&1
        $buildMissing = ($LASTEXITCODE -ne 0)
    }

    if (-not ($preCommitMissing -or $buildMissing)) {
        Write-Ok 'pre-commit'
    } elseif ($Check) {
        Write-Gap 'pre-commit (and packaging build tools when -WithPackaging)'
    } else {
        Write-Info "installing $($pipPackages -join ', ')"
        & $py -m pip install --user --upgrade --disable-pip-version-check @pipPackages
        if ($LASTEXITCODE -ne 0) {
            Write-Gap 'pip install failed'
        } else {
            # pip --user scripts land outside PATH on Windows by default.
            $userScripts = & $py -c "import site,os;print(os.path.join(site.USER_BASE,'Scripts'))"
            if ($userScripts -and (Test-Path $userScripts) -and ($env:Path -notlike "*$userScripts*")) {
                Add-UserPath $userScripts
            }
            Write-Ok "installed $($pipPackages -join ', ')"
        }
    }
}

Write-Section 'Optional tooling'
$optional = [ordered]@{
    'gh'     = 'needed to inspect CI runs'
    'docker' = 'needed for make docker-* targets'
}
foreach ($tool in $optional.Keys) {
    if (Test-Tool $tool) {
        Write-Ok "$tool ($(Get-ToolPath $tool))"
    } else {
        Write-Info "optional: $tool not found ($($optional[$tool]))"
    }
}

Write-Section 'Verification'
if ($goReady) {
    $smokeBin = Join-Path ([IO.Path]::GetTempPath()) "gc-dev-setup-$([guid]::NewGuid().ToString('N')).exe"
    Write-Info "go build -o $smokeBin ./cmd/gc"
    $oldProxy = $env:GOPROXY
    if ($Check) { $env:GOPROXY = 'off' }
    & go build -mod=readonly -o $smokeBin ./cmd/gc
    $buildExit = $LASTEXITCODE
    if ($Check) {
        if ($null -eq $oldProxy) { Remove-Item Env:GOPROXY -ErrorAction SilentlyContinue } else { $env:GOPROXY = $oldProxy }
    }
    if ($buildExit -eq 0) {
        Write-Ok 'build'
        & $smokeBin version > $null 2>&1
        if ($LASTEXITCODE -eq 0) { Write-Ok 'gc version' } else { Write-Gap 'gc version failed to run' }
        Remove-Item -LiteralPath $smokeBin -Force -ErrorAction SilentlyContinue
    } else {
        Write-Gap 'go build failed'
    }
} else {
    Write-Info 'skipped (Go unavailable)'
}

Write-Host ""
if ($script:Missing.Count -gt 0) {
    Write-Host "incomplete: $($script:Missing.Count) dependency gap(s)"
    foreach ($item in $script:Missing) { Write-Host "  - $item" }
    Write-Host ""
    Write-Host 'PATH changes are written to the user environment; open a new shell to pick them up.'
    exit 1
}

Write-Host 'dev environment ready.'
Write-Host ''
Write-Host 'Open a new shell so the updated PATH takes effect, then:'
Write-Host '  go test ./...'
Write-Host '  go build -o ./gc.exe ./cmd/gc; .\gc.exe version'
Write-Host '  bash scripts/regression-core.sh'
Write-Host ''
Write-Host 'Note: build ./gc.exe (not ./gc) on Windows, otherwise PowerShell refuses'
Write-Host 'to execute the extensionless binary. Set GC_BIN for regression-core.sh.'
Write-Host ''
Write-Host 'Real command verification needs GC_TOKEN (or GITCODE_TOKEN) exported by'
Write-Host 'hand and must target infra-test/* repositories only.'
