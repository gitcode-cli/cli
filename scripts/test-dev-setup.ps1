# Exercise the Windows setup workflow with a fully stubbed local toolchain.
[CmdletBinding()]
param()

$ErrorActionPreference = 'Stop'
$scriptPath = Join-Path $PSScriptRoot 'dev-setup.ps1'
$powerShell = Join-Path $PSHOME 'powershell.exe'
$tempRoot = Join-Path ([IO.Path]::GetTempPath()) "gc-dev-setup-test-$([guid]::NewGuid().ToString('N'))"
$stubBin = Join-Path $tempRoot 'system-bin'
$emptyBin = Join-Path $tempRoot 'empty-bin'
$fakeLocalAppData = Join-Path $tempRoot 'local-app-data'
$dispatcher = Join-Path $tempRoot 'stub-command.ps1'
$stubExecutable = Join-Path $tempRoot 'stub.exe'
$tokenSentinel = 'gc-token-must-not-reach-child'
$legacyTokenSentinel = 'legacy-token-must-not-reach-child'
$userPathBefore = [Environment]::GetEnvironmentVariable('Path', 'User')

function Assert-True {
    param([bool]$Condition, [string]$Message)
    if (-not $Condition) { throw $Message }
}

function Assert-Equal {
    param($Expected, $Actual, [string]$Message)
    if ($Expected -ne $Actual) {
        throw "$Message (expected=$Expected actual=$Actual)"
    }
}

function Set-TestFile {
    param([string]$Path, [string[]]$Lines)
    $parent = Split-Path -Parent $Path
    if ($parent) { New-Item -ItemType Directory -Force -Path $parent | Out-Null }
    Set-Content -LiteralPath $Path -Value $Lines -Encoding ASCII
}

function Restore-TestEnvironment {
    param([hashtable]$Saved)
    foreach ($name in $Saved.Keys) {
        if ($null -eq $Saved[$name]) {
            Remove-Item "Env:$name" -ErrorAction SilentlyContinue
        } else {
            Set-Item "Env:$name" $Saved[$name]
        }
    }
}

function Invoke-Setup {
    param(
        [string]$ToolsRoot,
        [string]$LogPath,
        [string]$SystemPath = $stubBin,
        [string]$GoVersion = 'go1.22.5',
        [switch]$Check
    )

    $variables = @{
        GC_DEV_TOOLS_DIR = $ToolsRoot
        GC_TOKEN = $tokenSentinel
        GITCODE_TOKEN = $legacyTokenSentinel
        GC_TEST_CHILD_LOG = $LogPath
        GC_TEST_STUB_EXE = $stubExecutable
        GC_TEST_GO_VERSION = $GoVersion
        GC_TEST_GO_LIST_EXIT = '0'
        GOCACHE = 'parent-gocache-sentinel'
        GOTMPDIR = 'parent-gotmp-sentinel'
        GOPROXY = 'http://127.0.0.1:1'
        LOCALAPPDATA = $fakeLocalAppData
        Path = "$SystemPath;$env:SystemRoot\System32;$env:SystemRoot"
    }
    $saved = @{}
    foreach ($name in $variables.Keys) {
        $saved[$name] = [Environment]::GetEnvironmentVariable($name, 'Process')
        Set-Item "Env:$name" $variables[$name]
    }

    try {
        $arguments = @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', $scriptPath)
        if ($Check) { $arguments += '-Check' }
        $oldErrorAction = $ErrorActionPreference
        $ErrorActionPreference = 'Continue'
        $output = & $powerShell @arguments 2>&1 | Out-String
        $status = $LASTEXITCODE
        $ErrorActionPreference = $oldErrorAction
        return [pscustomobject]@{ Status = $status; Output = $output }
    } finally {
        $ErrorActionPreference = 'Stop'
        Restore-TestEnvironment $saved
    }
}

function Get-NewLogLines {
    param([string]$Path, [int]$Start)
    if (-not (Test-Path $Path)) { return @() }
    return @(Get-Content -LiteralPath $Path | Select-Object -Skip $Start)
}

try {
    New-Item -ItemType Directory -Force -Path $tempRoot, $stubBin, $emptyBin | Out-Null

    $executableSource = @'
using System;
using System.IO;

public static class Program
{
    public static int Main(string[] args)
    {
        string log = Environment.GetEnvironmentVariable("GC_TEST_CHILD_LOG");
        if (!String.IsNullOrEmpty(log))
        {
            string token = Environment.GetEnvironmentVariable("GC_TOKEN") ?? "";
            string legacy = Environment.GetEnvironmentVariable("GITCODE_TOKEN") ?? "";
            string cache = Environment.GetEnvironmentVariable("GOCACHE") ?? "";
            string tmp = Environment.GetEnvironmentVariable("GOTMPDIR") ?? "";
            File.AppendAllText(log, "exe|GC=" + token + "|LEGACY=" + legacy +
                "|CACHE=" + cache + "|TMP=" + tmp + Environment.NewLine);
        }
        if (args.Length > 0 && args[0] == "--version")
        {
            Console.WriteLine("pre-commit 4.6.1");
        }
        return 0;
    }
}
'@
    Add-Type -TypeDefinition $executableSource -Language CSharp `
        -OutputAssembly $stubExecutable -OutputType ConsoleApplication

    $dispatcherSource = @'
param(
    [Parameter(Position = 0)]
    [string]$Tool,
    [Parameter(Position = 1, ValueFromRemainingArguments = $true)]
    [string[]]$Arguments,
    [Parameter()]
    [Alias('o')]
    [string]$OutputPath
)

$ErrorActionPreference = 'Stop'
$record = "$Tool|$($Arguments -join ' ')|GC=$env:GC_TOKEN|LEGACY=$env:GITCODE_TOKEN" +
    "|CACHE=$env:GOCACHE|TMP=$env:GOTMPDIR|GOBIN=$env:GOBIN"
Add-Content -LiteralPath $env:GC_TEST_CHILD_LOG -Value $record

if ($Tool -eq 'go') {
    if ($Arguments[0] -eq 'version' -and $Arguments.Count -eq 1) {
        Write-Output "go version $env:GC_TEST_GO_VERSION windows/amd64"
        exit 0
    }
    if ($Arguments[0] -eq 'env' -and $Arguments[1] -eq 'GOVERSION') {
        Write-Output $env:GC_TEST_GO_VERSION
        exit 0
    }
    if ($Arguments[0] -eq 'version' -and $Arguments[1] -eq '-m') {
        $binary = $Arguments[2]
        if (-not (Test-Path -LiteralPath $binary)) { exit 1 }
        $content = Get-Content -Raw -LiteralPath $binary
        $module = [regex]::Match($content, '(?m)^module=(.+)$').Groups[1].Value.Trim()
        $version = [regex]::Match($content, '(?m)^version=(.+)$').Groups[1].Value.Trim()
        if (-not $module -or -not $version) { exit 1 }
        Write-Output $binary
        Write-Output "`tmod`t$module`t$version"
        exit 0
    }
    if ($Arguments[0] -eq 'list') { exit [int]$env:GC_TEST_GO_LIST_EXIT }
    if ($Arguments[0] -eq 'mod' -and $Arguments[1] -eq 'download') { exit 0 }
    if ($Arguments[0] -eq 'install') {
        $spec = $Arguments[1]
        $at = $spec.LastIndexOf('@')
        $package = $spec.Substring(0, $at)
        $version = $spec.Substring($at + 1)
        if ($package -eq 'github.com/golangci/golangci-lint/v2/cmd/golangci-lint') {
            $name = 'golangci-lint'
            $module = 'github.com/golangci/golangci-lint/v2'
        } elseif ($package -eq 'github.com/zricethezav/gitleaks/v8') {
            $name = 'gitleaks'
            $module = $package
        } else {
            exit 12
        }
        New-Item -ItemType Directory -Force -Path $env:GOBIN | Out-Null
        Set-Content -LiteralPath (Join-Path $env:GOBIN "$name.exe") -Encoding ASCII -Value @(
            'GC_STUB'
            "module=$module"
            "version=$version"
        )
        exit 0
    }
    if ($Arguments[0] -eq 'build') {
        if (-not $OutputPath) { exit 13 }
        Copy-Item -LiteralPath $env:GC_TEST_STUB_EXE -Destination $OutputPath -Force
        exit 0
    }
    if ($Arguments[0] -eq 'test') { exit 0 }
    exit 14
}

if ($Tool -in @('python', 'python3')) {
    if ($Arguments[0] -eq '--version') {
        Write-Output 'Python 3.12.4'
        exit 0
    }
    if ($Arguments[0] -eq '-m' -and $Arguments[1] -eq 'venv') {
        $scripts = Join-Path $Arguments[2] 'Scripts'
        New-Item -ItemType Directory -Force -Path $scripts | Out-Null
        Copy-Item -LiteralPath $env:GC_TEST_STUB_EXE -Destination (Join-Path $scripts 'python.exe')
        Copy-Item -LiteralPath $env:GC_TEST_STUB_EXE -Destination (Join-Path $scripts 'pre-commit.exe')
        exit 0
    }
    exit 0
}

if ($Arguments[0] -eq '--version') { Write-Output "$Tool stub 1.0" }
exit 0
'@
    Set-Content -LiteralPath $dispatcher -Value $dispatcherSource -Encoding UTF8
    foreach ($tool in @('go', 'git', 'python', 'python3', 'make', 'mingw32-make', 'gcc')) {
        Set-TestFile (Join-Path $stubBin "$tool.cmd") @(
            '@echo off'
            "`"$powerShell`" -NoProfile -ExecutionPolicy Bypass -File `"$dispatcher`" $tool %*"
            'exit /b %ERRORLEVEL%'
        )
    }

    $gitBash = Join-Path $fakeLocalAppData 'Programs\Git\bin\bash.exe'
    New-Item -ItemType Directory -Force -Path (Split-Path -Parent $gitBash) | Out-Null
    Copy-Item -LiteralPath $stubExecutable -Destination $gitBash

    Write-Host '[test] installation and idempotency'
    $installTools = Join-Path $tempRoot 'install-tools'
    $installLog = Join-Path $tempRoot 'install.log'
    $first = Invoke-Setup -ToolsRoot $installTools -LogPath $installLog
    Assert-Equal 0 $first.Status "first installation failed`n$($first.Output)"
    Assert-True ([bool]($first.Output -match [regex]::Escape(
        "`$env:Path = `"$installTools\bin;`$env:Path`""))) `
        'setup output did not provide the current-shell PATH command'
    $installLines = @(Get-Content -LiteralPath $installLog)
    Assert-True ([bool]($installLines -match
        '^go\|install github\.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2\.12\.2\|')) `
        'golangci-lint was not installed from its cmd package'
    Assert-True ([bool]($installLines -match
        '^go\|install github\.com/zricethezav/gitleaks/v8@v8\.30\.1\|')) `
        'gitleaks install package was not used'
    foreach ($name in @('golangci-lint.exe', 'gitleaks.exe', 'make.cmd', 'python3.cmd', 'pre-commit.cmd')) {
        Assert-True (Test-Path (Join-Path $installTools "bin\$name")) "missing managed artifact: $name"
    }
    $goInstallCount = @($installLines -match '^go\|install ').Count
    $second = Invoke-Setup -ToolsRoot $installTools -LogPath $installLog
    Assert-Equal 0 $second.Status "second installation failed`n$($second.Output)"
    $secondInstallCount = @((Get-Content -LiteralPath $installLog) -match '^go\|install ').Count
    Assert-Equal $goInstallCount $secondInstallCount 'second installation was not idempotent'
    $leftovers = @(Get-ChildItem -Force -Recurse -LiteralPath $installTools |
        Where-Object { $_.Name -like '.*-install-*' -or $_.Name -like '*.tmp' })
    Assert-Equal 0 $leftovers.Count 'temporary publish files were left behind'

    Write-Host '[test] complete and incomplete check mode'
    $beforeCheck = @(Get-Content -LiteralPath $installLog).Count
    $complete = Invoke-Setup -ToolsRoot $installTools -LogPath $installLog -Check
    Assert-Equal 0 $complete.Status "complete check mode failed`n$($complete.Output)"
    $checkLines = Get-NewLogLines $installLog $beforeCheck
    $goCheckLines = @($checkLines | Where-Object { $_ -like 'go|*' })
    Assert-True ($goCheckLines.Count -gt 0) 'check mode did not invoke the Go stub'
    foreach ($line in $goCheckLines) {
        Assert-True ($line -match '\|CACHE=([^|]*gc-dev-check-[^|]*\\gocache)\|TMP=([^|]*gc-dev-check-[^|]*\\gotmp)') `
            "check mode did not isolate Go cache/temp: $line"
        Assert-True (-not (Test-Path $Matches[1])) "isolated GOCACHE was not removed: $($Matches[1])"
        Assert-True (-not (Test-Path $Matches[2])) "isolated GOTMPDIR was not removed: $($Matches[2])"
    }
    $missingTools = Join-Path $tempRoot 'missing-tools'
    $missing = Invoke-Setup -ToolsRoot $missingTools -LogPath (Join-Path $tempRoot 'missing.log') `
        -SystemPath $emptyBin -Check
    Assert-Equal 1 $missing.Status "missing dependency check did not fail`n$($missing.Output)"
    Assert-True (-not (Test-Path $missingTools)) 'check mode created the managed tools directory'

    Write-Host '[test] child credential isolation'
    $allLog = Get-Content -Raw -LiteralPath $installLog
    Assert-True ($allLog -notmatch [regex]::Escape($tokenSentinel)) 'GC_TOKEN reached a child process'
    Assert-True ($allLog -notmatch [regex]::Escape($legacyTokenSentinel)) `
        'GITCODE_TOKEN reached a child process'

    Write-Host '[test] unmanaged artifact protection and system PATH filtering'
    $unmanagedTools = Join-Path $tempRoot 'unmanaged-tools'
    $unmanagedBin = Join-Path $unmanagedTools 'bin'
    New-Item -ItemType Directory -Force -Path $unmanagedBin | Out-Null
    Set-TestFile (Join-Path $unmanagedBin 'golangci-lint.exe') @('do-not-overwrite')
    Set-TestFile (Join-Path $unmanagedBin 'make.cmd') @('@echo unmanaged-wrapper')
    $hijackMarker = Join-Path $tempRoot 'managed-go-was-run'
    Set-TestFile (Join-Path $unmanagedBin 'go.cmd') @(
        '@echo off'
        "echo bad>>`"$hijackMarker`""
        'exit /b 77'
    )
    $unmanaged = Invoke-Setup -ToolsRoot $unmanagedTools -LogPath (Join-Path $tempRoot 'unmanaged.log') `
        -SystemPath "$unmanagedBin;$stubBin"
    Assert-Equal 1 $unmanaged.Status "unmanaged overwrite scenario unexpectedly passed`n$($unmanaged.Output)"
    Assert-Equal 'do-not-overwrite' ((Get-Content -Raw -LiteralPath `
        (Join-Path $unmanagedBin 'golangci-lint.exe')).Trim()) 'unmanaged Go tool was overwritten'
    Assert-Equal '@echo unmanaged-wrapper' ((Get-Content -Raw -LiteralPath `
        (Join-Path $unmanagedBin 'make.cmd')).Trim()) 'unmanaged wrapper was overwritten'
    Assert-True (-not (Test-Path $hijackMarker)) 'ManagedBin hijacked system Go resolution'
    Assert-True ($unmanaged.Output -match 'refusing to overwrite unmanaged Go tool') `
        'unmanaged Go tool refusal was not reported'
    Assert-True ($unmanaged.Output -match 'refusing to overwrite unmanaged wrapper') `
        'unmanaged wrapper refusal was not reported'

    Write-Host '[test] reparse-point protection'
    $reparseTools = Join-Path $tempRoot 'reparse-tools'
    $reparseTarget = Join-Path $tempRoot 'reparse-target'
    New-Item -ItemType Directory -Force -Path $reparseTools, $reparseTarget | Out-Null
    New-Item -ItemType Junction -Path (Join-Path $reparseTools 'bin') -Target $reparseTarget | Out-Null
    $reparse = Invoke-Setup -ToolsRoot $reparseTools -LogPath (Join-Path $tempRoot 'reparse.log')
    Assert-Equal 1 $reparse.Status "reparse-point scenario unexpectedly passed`n$($reparse.Output)"
    Assert-Equal 0 @(Get-ChildItem -Force -LiteralPath $reparseTarget).Count `
        'script wrote through the ManagedBin reparse point'
    Assert-True ($reparse.Output -match 'refusing to write through reparse point') `
        'reparse-point refusal was not reported'

    Write-Host '[test] old Go without winget guidance'
    $oldGoTools = Join-Path $tempRoot 'old-go-tools'
    $oldGo = Invoke-Setup -ToolsRoot $oldGoTools -LogPath (Join-Path $tempRoot 'old-go.log') `
        -GoVersion 'go1.21.13'
    Assert-Equal 1 $oldGo.Status "old Go scenario unexpectedly passed`n$($oldGo.Output)"
    Assert-True ($oldGo.Output -match 'winget unavailable; download and install from https://go.dev/dl/') `
        'old Go without winget did not provide actionable guidance'

    Assert-Equal $userPathBefore ([Environment]::GetEnvironmentVariable('Path', 'User')) `
        'tests or setup script changed the user PATH'
    Write-Host 'Windows dev setup tests passed.'
} finally {
    if (Test-Path $tempRoot) {
        Remove-Item -LiteralPath $tempRoot -Recurse -Force -ErrorAction SilentlyContinue
    }
}
