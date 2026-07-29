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
$ManagedGoTools = @{
    'golangci-lint' = @{
        Module = 'github.com/golangci/golangci-lint/v2'
        Package = 'github.com/golangci/golangci-lint/v2/cmd/golangci-lint'
        Version = $GolangciVersion
    }
    'gitleaks' = @{
        Module = 'github.com/zricethezav/gitleaks/v8'
        Package = 'github.com/zricethezav/gitleaks/v8'
        Version = $GitleaksVersion
    }
}
$ManagedWrappers = @('make.cmd', 'pre-commit.cmd', 'python3.cmd')
$OldGoCache = $env:GOCACHE
$OldGoTmpDir = $env:GOTMPDIR
$OldGoProxy = $env:GOPROXY
$CheckTempRoot = $null
$script:Missing = @()

function Write-Section { param([string]$Name) Write-Host ''; Write-Host "[$Name]" }
function Write-Ok { param([string]$Message) Write-Host "ok   $Message" }
function Write-Info { param([string]$Message) Write-Host "  $Message" }
function Write-Gap { param([string]$Message) Write-Host "MISS $Message"; $script:Missing += $Message }

function Restore-EnvironmentVariable {
    param([string]$Name, [AllowNull()][string]$Value)
    if ($null -eq $Value) {
        Remove-Item "Env:$Name" -ErrorAction SilentlyContinue
    } else {
        Set-Item "Env:$Name" $Value
    }
}

function Test-IsManagedPath {
    param([string]$Path)
    if (-not $Path) { return $false }
    try {
        $root = [IO.Path]::GetFullPath($ManagedBin).TrimEnd('\', '/')
        $candidate = [IO.Path]::GetFullPath($Path.Trim('"'))
    } catch {
        return $false
    }
    return $candidate.Equals($root, [StringComparison]::OrdinalIgnoreCase) -or
        $candidate.StartsWith("$root\", [StringComparison]::OrdinalIgnoreCase)
}

function Get-SystemCommandPath {
    param([string]$Name)
    $commands = Get-Command $Name -All -CommandType Application -ErrorAction SilentlyContinue
    foreach ($command in $commands) {
        if ($command.Source -and -not (Test-IsManagedPath $command.Source)) {
            return $command.Source
        }
    }
    return $null
}

function Test-SystemCommandWorks {
    param([string]$Name, [string[]]$Arguments = @('--version'))
    $path = Get-SystemCommandPath $Name
    if (-not $path) { return $false }
    & $path @Arguments > $null 2>&1
    return ($LASTEXITCODE -eq 0)
}

function Refresh-SystemPath {
    param([switch]$IncludePersistent)
    $paths = @($env:Path)
    if ($IncludePersistent) {
        $paths += [Environment]::GetEnvironmentVariable('Path', 'Machine')
        $paths += [Environment]::GetEnvironmentVariable('Path', 'User')
    }
    $entries = $paths -join ';' -split ';' | Where-Object {
        $_ -and -not (Test-IsManagedPath $_)
    }
    $env:Path = @($entries | Select-Object -Unique) -join ';'
}

function Get-PathItem {
    param([string]$Path)
    return Get-Item -Force -LiteralPath $Path -ErrorAction SilentlyContinue
}

function Test-IsReparsePoint {
    param([string]$Path)
    $item = Get-PathItem $Path
    return $item -and (($item.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0)
}

function Test-ManagedLocationSafe {
    param([string]$Path)
    foreach ($candidate in @($ToolsRoot, $ManagedBin, $Path)) {
        if ((Get-PathItem $candidate) -and (Test-IsReparsePoint $candidate)) { return $false }
    }
    return $true
}

function Ensure-SafeManagedDirectory {
    param([string]$Path)
    foreach ($candidate in @($ToolsRoot, $Path)) {
        $item = Get-PathItem $candidate
        if ($item -and (($item.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0)) {
            Write-Gap "refusing to write through reparse point: $candidate"
            return $false
        }
        if ($item -and -not $item.PSIsContainer) {
            Write-Gap "managed directory path is not a directory: $candidate"
            return $false
        }
    }
    New-Item -ItemType Directory -Force -Path $Path | Out-Null
    if (Test-IsReparsePoint $Path) {
        Write-Gap "refusing to write through reparse point: $Path"
        return $false
    }
    return $true
}

function Publish-FileAtomically {
    param([string]$Source, [string]$Destination)
    if (Get-PathItem $Destination) {
        $backup = "$Destination.$([guid]::NewGuid().ToString('N')).bak"
        try {
            [IO.File]::Replace($Source, $Destination, $backup)
        } finally {
            Remove-Item -LiteralPath $backup -Force -ErrorAction SilentlyContinue
        }
    } else {
        [IO.File]::Move($Source, $Destination)
    }
}

function Install-WingetPackage {
    param([string]$Name, [string]$PackageId, [string[]]$VersionArguments = @('--version'))
    if (Test-SystemCommandWorks $Name $VersionArguments) { return }
    if ($Check) { return }
    $winget = Get-SystemCommandPath winget
    if (-not $winget) {
        Write-Gap "$Name (winget unavailable; install $PackageId manually)"
        return
    }
    Write-Info "installing $Name with winget ($PackageId)"
    & $winget install --id $PackageId --exact --source winget `
        --accept-package-agreements --accept-source-agreements --disable-interactivity
    if ($LASTEXITCODE -ne 0) { Write-Gap "$Name (winget install failed)"; return }
    Refresh-SystemPath -IncludePersistent
}

function Initialize-CheckEnvironment {
    if (-not $Check) { return }
    $script:CheckTempRoot = Join-Path ([IO.Path]::GetTempPath()) "gc-dev-check-$([guid]::NewGuid().ToString('N'))"
    $cache = Join-Path $script:CheckTempRoot 'gocache'
    $tmp = Join-Path $script:CheckTempRoot 'gotmp'
    New-Item -ItemType Directory -Force -Path $cache, $tmp | Out-Null
    $env:GOCACHE = $cache
    $env:GOTMPDIR = $tmp
    $env:GOPROXY = 'off'
}

function Remove-CheckEnvironment {
    Restore-EnvironmentVariable GOCACHE $OldGoCache
    Restore-EnvironmentVariable GOTMPDIR $OldGoTmpDir
    Restore-EnvironmentVariable GOPROXY $OldGoProxy
    if ($script:CheckTempRoot) {
        Remove-Item -LiteralPath $script:CheckTempRoot -Recurse -Force -ErrorAction SilentlyContinue
    }
}

function Test-GoVersion {
    $go = Get-SystemCommandPath go
    if (-not $go) { return $false }
    & $go version > $null 2>&1
    if ($LASTEXITCODE -ne 0) { return $false }
    $raw = (& $go env GOVERSION 2>$null) -replace '^go', ''
    $parts = $raw.Split('.')
    if ($parts.Count -lt 2) { return $false }
    $major = 0; $minor = 0
    if (-not [int]::TryParse($parts[0], [ref]$major)) { return $false }
    if (-not [int]::TryParse($parts[1], [ref]$minor)) { return $false }
    return ($major -gt $GoMinMajor -or ($major -eq $GoMinMajor -and $minor -ge $GoMinMinor))
}

function Get-GoModuleVersion {
    param([string]$Binary, [string]$Module)
    $go = Get-SystemCommandPath go
    if (-not $go -or -not $Binary -or -not (Get-PathItem $Binary)) { return '' }
    if (Test-IsReparsePoint $Binary) { return '' }
    $lines = & $go version -m $Binary 2>$null
    foreach ($line in $lines) {
        $fields = $line.Trim() -split '\s+'
        if ($fields.Count -ge 3 -and $fields[0] -eq 'mod' -and $fields[1] -eq $Module) { return $fields[2] }
    }
    return ''
}

function Ensure-GoTool {
    param([string]$Name, [string]$Module, [string]$Package, [string]$Version)
    $definition = $ManagedGoTools[$Name]
    if (-not $definition -or $definition.Module -ne $Module -or
        $definition.Package -ne $Package -or $definition.Version -ne $Version) {
        Write-Gap "refusing non-whitelisted managed Go tool: $Name"
        return
    }

    $binary = Join-Path $ManagedBin "$Name.exe"
    if (-not (Test-ManagedLocationSafe $binary)) {
        Write-Gap "refusing managed Go tool through reparse point: $binary"
        return
    }
    $item = Get-PathItem $binary
    if ($item -and (Test-IsReparsePoint $binary)) {
        Write-Gap "refusing to replace reparse point: $binary"
        return
    }
    $current = Get-GoModuleVersion $binary $Module
    if ($current -eq $Version) { Write-Ok "$Name $Version ($binary)"; return }
    if ($item -and -not $current) {
        Write-Gap "refusing to overwrite unmanaged Go tool: $binary"
        return
    }
    if ($Check) { Write-Gap "$Name $Version$(if($current){" (found $current)"})"; return }

    $go = Get-SystemCommandPath go
    if (-not $go -or -not (Ensure-SafeManagedDirectory $ToolsRoot)) {
        Write-Gap "$Name $Version (safe Go installation unavailable)"
        return
    }

    $stagingRoot = Join-Path $ToolsRoot ".$Name-install-$([guid]::NewGuid().ToString('N'))"
    $stagingBin = Join-Path $stagingRoot 'bin'
    $oldGoBin = $env:GOBIN
    try {
        New-Item -ItemType Directory -Force -Path $stagingBin | Out-Null
        $env:GOBIN = $stagingBin
        Write-Info "installing $Name $Version into temporary GOBIN"
        & $go install "$Package@$Version"
        if ($LASTEXITCODE -ne 0) {
            Write-Gap "$Name $Version (go install failed)"
            return
        }

        $staged = Join-Path $stagingBin "$Name.exe"
        if ((Get-GoModuleVersion $staged $Module) -ne $Version) {
            Write-Gap "$Name $Version (installed module or version mismatch)"
            return
        }
        if (-not (Ensure-SafeManagedDirectory $ManagedBin)) { return }

        $existing = Get-PathItem $binary
        if ($existing) {
            if (Test-IsReparsePoint $binary) {
                Write-Gap "refusing to replace reparse point: $binary"
                return
            }
            if (-not (Get-GoModuleVersion $binary $Module)) {
                Write-Gap "refusing to overwrite unmanaged Go tool: $binary"
                return
            }
        }

        try {
            Publish-FileAtomically $staged $binary
        } catch {
            Write-Gap "$Name $Version (atomic publish failed: $($_.Exception.Message))"
            return
        }
        if ((Get-GoModuleVersion $binary $Module) -ne $Version) {
            Write-Gap "$Name $Version (published module verification failed)"
            return
        }
        Write-Ok "$Name $Version ($binary)"
    } finally {
        Restore-EnvironmentVariable GOBIN $oldGoBin
        Remove-Item -LiteralPath $stagingRoot -Recurse -Force -ErrorAction SilentlyContinue
    }
}

function Find-GitBash {
    foreach ($candidate in @(
        (Join-Path $env:LOCALAPPDATA 'Programs\Git\bin\bash.exe'),
        'C:\Program Files\Git\bin\bash.exe',
        'C:\Program Files (x86)\Git\bin\bash.exe'
    )) { if (Test-Path $candidate) { return $candidate } }
    return $null
}

function Write-ManagedWrapper {
    param([string]$Path, [string[]]$Lines)
    $leaf = [IO.Path]::GetFileName($Path)
    $parent = [IO.Path]::GetDirectoryName([IO.Path]::GetFullPath($Path))
    if ($ManagedWrappers -notcontains $leaf -or
        -not $parent.Equals([IO.Path]::GetFullPath($ManagedBin), [StringComparison]::OrdinalIgnoreCase)) {
        Write-Gap "refusing non-whitelisted managed wrapper: $Path"
        return $false
    }
    if (-not (Ensure-SafeManagedDirectory $ManagedBin)) { return $false }

    $item = Get-PathItem $Path
    if ($item) {
        if (Test-IsReparsePoint $Path) {
            Write-Gap "refusing to replace reparse point: $Path"
            return $false
        }
        $first = Get-Content -LiteralPath $Path -TotalCount 1
        if ($first -ne $ManagedMarker) {
            Write-Gap "refusing to overwrite unmanaged wrapper: $Path"
            return $false
        }
    }

    $temporary = Join-Path $ManagedBin ".$leaf-$([guid]::NewGuid().ToString('N')).tmp"
    try {
        [IO.File]::WriteAllLines($temporary, @($ManagedMarker, '@echo off') + $Lines, [Text.Encoding]::ASCII)
        if (Test-IsReparsePoint $Path) {
            Write-Gap "refusing to replace reparse point: $Path"
            return $false
        }
        Publish-FileAtomically $temporary $Path
    } catch {
        Write-Gap "failed to publish managed wrapper atomically: $Path ($($_.Exception.Message))"
        return $false
    } finally {
        Remove-Item -LiteralPath $temporary -Force -ErrorAction SilentlyContinue
    }
    return $true
}

function Test-ManagedWrapperOwnership {
    param([string]$Path)
    if (-not (Test-ManagedLocationSafe $Path)) { return $false }
    if (-not (Get-PathItem $Path) -or (Test-IsReparsePoint $Path)) { return $false }
    if ($ManagedWrappers -notcontains [IO.Path]::GetFileName($Path)) { return $false }
    if ((Get-Content -LiteralPath $Path -TotalCount 1) -ne $ManagedMarker) { return $false }
    return $true
}

function Test-ManagedWrapper {
    param([string]$Path, [string[]]$Arguments = @('--version'))
    if (-not (Test-ManagedWrapperOwnership $Path)) { return $false }
    & $Path @Arguments > $null 2>&1
    return ($LASTEXITCODE -eq 0)
}

function Find-Python {
    foreach ($name in @('python3', 'python')) {
        $command = Get-SystemCommandPath $name
        if ($command) {
            & $command --version > $null 2>&1
            if ($LASTEXITCODE -eq 0) { return $command }
        }
    }
    $launcher = Get-SystemCommandPath py
    if ($launcher) {
        $candidate = & $launcher -3 -c "import sys; print(sys.executable)" 2>$null
        if ($LASTEXITCODE -eq 0 -and $candidate -and (Test-Path $candidate)) { return $candidate }
    }
    return $null
}

function Ensure-PreCommit {
    param([string]$Python)
    $binary = Join-Path $ManagedBin 'pre-commit.cmd'
    $current = ''
    if (Test-ManagedWrapperOwnership $binary) {
        $current = ((& $binary --version 2>$null) -split '\s+')[-1]
    }
    if ($current -eq $PreCommitVersion) { Write-Ok "pre-commit $PreCommitVersion ($binary)"; return }
    if ($Check) { Write-Gap "pre-commit $PreCommitVersion$(if($current){" (found $current)"})"; return }
    if ((Get-PathItem $PreCommitVenv) -and (Test-IsReparsePoint $PreCommitVenv)) {
        Write-Gap "refusing to write pre-commit through reparse point: $PreCommitVenv"
        return
    }
    if (-not (Ensure-SafeManagedDirectory $ToolsRoot)) { return }
    Write-Info "installing pre-commit $PreCommitVersion into $PreCommitVenv"
    & $Python -m venv $PreCommitVenv
    if ($LASTEXITCODE -eq 0) {
        $venvPython = Join-Path $PreCommitVenv 'Scripts\python.exe'
        & $venvPython -m pip install --disable-pip-version-check "pre-commit==$PreCommitVersion" > $null
    }
    if ($LASTEXITCODE -ne 0) { Write-Gap 'pre-commit venv install failed'; return }
    $target = Join-Path $PreCommitVenv 'Scripts\pre-commit.exe'
    if (-not (Get-PathItem $target) -or (Test-IsReparsePoint $target)) {
        Write-Gap 'pre-commit installation target is missing or a reparse point'
        return
    }
    if (Write-ManagedWrapper $binary @("`"$target`" %*")) {
        Write-Ok "pre-commit $PreCommitVersion ($binary)"
    }
}

$result = 0
try {
    Initialize-CheckEnvironment
    Refresh-SystemPath
    if (-not $Check) {
        Write-Section 'System dependencies'
        Install-WingetPackage go GoLang.Go @('version')
        Install-WingetPackage git Git.Git
        Install-WingetPackage python Python.Python.3.12
        Install-WingetPackage make ezwinports.make
        Install-WingetPackage gcc BrechtSanders.WinLibs.MCF.UCRT
    }

    $go = Get-SystemCommandPath go
    if (-not $Check -and $go -and -not (Test-GoVersion)) {
        $winget = Get-SystemCommandPath winget
        if ($winget) {
            Write-Info 'upgrading Go with winget'
            & $winget upgrade --id GoLang.Go --exact --source winget `
                --accept-package-agreements --accept-source-agreements --disable-interactivity
            if ($LASTEXITCODE -ne 0) { Write-Gap 'Go upgrade failed' }
            Refresh-SystemPath -IncludePersistent
        } else {
            $guidance = "Go $GoMinMajor.$GoMinMinor+ (winget unavailable; download and install from " +
                'https://go.dev/dl/, then rerun this script)'
            Write-Gap $guidance
        }
    }

    Write-Section 'Go toolchain'
    if (-not (Test-GoVersion)) {
        Write-Gap "Go $GoMinMajor.$GoMinMinor+ (install with winget or download from https://go.dev/dl/)"
    } else {
        $go = Get-SystemCommandPath go
        Write-Ok "go $(& $go env GOVERSION) ($go)"
    }
    $goReady = Test-GoVersion

    if ($goReady) {
        Write-Section 'Module dependencies'
        $go = Get-SystemCommandPath go
        & $go list -mod=readonly -deps ./... > $null 2>&1
        $listExit = $LASTEXITCODE
        if ($listExit -eq 0) { Write-Ok 'module dependencies resolved' }
        elseif ($Check) { Write-Gap 'module dependencies (run without -Check to download)' }
        else {
            & $go mod download
            if ($LASTEXITCODE -eq 0) { Write-Ok 'module dependencies resolved' }
            else { Write-Gap 'go mod download failed' }
        }

        Write-Section 'Verification tools'
        $golangci = $ManagedGoTools['golangci-lint']
        Ensure-GoTool golangci-lint $golangci.Module $golangci.Package $golangci.Version
        $gitleaks = $ManagedGoTools['gitleaks']
        Ensure-GoTool gitleaks $gitleaks.Module $gitleaks.Package $gitleaks.Version
    }

    Write-Section 'Git and shell utilities'
    if (Test-SystemCommandWorks git) { Write-Ok "git ($(Get-SystemCommandPath git))" }
    else { Write-Gap 'git' }
    $gitBash = Find-GitBash
    if ($gitBash) { Write-Ok "git bash ($gitBash)" } else { Write-Gap 'git bash' }

    $realMake = $null
    foreach ($name in @('mingw32-make', 'make')) {
        $found = Get-SystemCommandPath $name
        if ($found) { $realMake = $found; break }
    }
    $makeWrapper = Join-Path $ManagedBin 'make.cmd'
    if ($Check) {
        if ($gitBash -and (Test-ManagedWrapper $makeWrapper)) {
            Write-Ok "managed make wrapper ($makeWrapper)"
        } else { Write-Gap 'managed make wrapper' }
    } elseif ($realMake -and $gitBash) {
        $fso = New-Object -ComObject Scripting.FileSystemObject
        $shortBash = ($fso.GetFile($gitBash).ShortPath) -replace '\\', '/'
        [void][Runtime.InteropServices.Marshal]::ReleaseComObject($fso)
        if (Write-ManagedWrapper $makeWrapper @(
            "`"$realMake`" `"SHELL=$shortBash`" `".SHELLFLAGS=-c`" %*"
        )) { Write-Ok "managed make wrapper ($makeWrapper)" }
    } else { Write-Gap 'GNU make and Git bash' }

    $python = Find-Python
    if ($python) {
        Write-Ok "python ($python)"
        $pythonWrapper = Join-Path $ManagedBin 'python3.cmd'
        if (-not $Check) {
            if (Write-ManagedWrapper $pythonWrapper @("`"$python`" %*")) {
                Write-Ok "managed python3 wrapper ($pythonWrapper)"
            }
        } elseif (Test-ManagedWrapper $pythonWrapper) { Write-Ok 'managed python3 wrapper' }
        else { Write-Gap 'managed python3 wrapper' }
        Ensure-PreCommit $python
    } else { Write-Gap 'Python 3' }

    $compiler = $null
    foreach ($name in @('gcc', 'clang', 'cc')) {
        $found = Get-SystemCommandPath $name
        if ($found) { $compiler = $found; break }
    }
    if ($compiler) { Write-Ok "C compiler ($compiler)" }
    else { Write-Gap 'C compiler (required by go test -race)' }

    Write-Section 'Verification'
    if ($goReady) {
        $smokeDir = Join-Path ([IO.Path]::GetTempPath()) "gc-dev-setup-$([guid]::NewGuid().ToString('N'))"
        New-Item -ItemType Directory -Force -Path $smokeDir | Out-Null
        try {
            $go = Get-SystemCommandPath go
            $smokeBinary = Join-Path $smokeDir 'gc.exe'
            & $go build -mod=readonly -o $smokeBinary ./cmd/gc
            $buildExit = $LASTEXITCODE
            if ($buildExit -eq 0) { & $smokeBinary version > $null 2>&1 }
            if ($buildExit -eq 0 -and $LASTEXITCODE -eq 0) {
                Write-Ok 'build and gc version'
            } else { Write-Gap 'build or gc version' }
            if ($compiler) {
                & $go test -mod=readonly -race ./pkg/config > $null 2>&1
                if ($LASTEXITCODE -eq 0) { Write-Ok 'race-enabled test' }
                else { Write-Gap 'race-enabled test' }
            }
        } finally {
            Remove-Item -LiteralPath $smokeDir -Recurse -Force -ErrorAction SilentlyContinue
        }
    }

    Write-Host ''
    if ($script:Missing.Count) {
        Write-Host "incomplete: $($script:Missing.Count) dependency gap(s)"
        foreach ($item in $script:Missing) { Write-Host "  - $item" }
        Write-Host "Managed tools directory: $ManagedBin"
        Write-Host 'After resolving gaps, add it to the current shell PATH:'
        Write-Host ('  $env:Path = "{0};$env:Path"' -f $ManagedBin)
        $result = 1
    } else {
        Write-Host 'dev environment ready.'
        Write-Host "Managed tools: $ManagedBin"
        Write-Host 'No user PATH was changed. To use managed tools in this shell, run:'
        Write-Host ('  $env:Path = "{0};$env:Path"' -f $ManagedBin)
    }
} finally {
    Remove-CheckEnvironment
}
exit $result
