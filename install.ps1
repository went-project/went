param(
    [string]$Version = $env:WENT_VERSION,
    [string]$Channel = $env:WENT_CHANNEL
)

if (-not $Channel) {
    $Channel = 'stable'
}

$ErrorActionPreference = 'Stop'
$repo = if ($env:WENT_REPOSITORY) { $env:WENT_REPOSITORY } else { 'went-project/went' }
$binDir = Join-Path $env:LOCALAPPDATA 'went\bin'
$binaryName = 'went.exe'

function Write-Info {
    param([string]$Message)
    Write-Host $Message
}

function Throw-InstallError {
    param([string]$Message)
    throw $Message
}

function Download-File {
    param(
        [string]$Url,
        [string]$OutputPath
    )

    try {
        Invoke-WebRequest -Uri $Url -OutFile $OutputPath -UseBasicParsing
    }
    catch {
        Throw-InstallError "Download failed: $Url"
    }
}

function Get-LatestTag {
    if ($Channel -eq 'beta') {
        $apiUrl = "https://api.github.com/repos/$repo/releases?per_page=100"
        try {
            $releases = Invoke-RestMethod -Uri $apiUrl -UseBasicParsing
            $betaRelease = $releases | Where-Object { $_.prerelease -and $_.tag_name -match '^v[0-9]+\.[0-9]+\.[0-9]+-beta\.[0-9]+$' } | Select-Object -First 1
            if (-not $betaRelease) {
                Throw-InstallError 'Unable to determine the latest beta release tag.'
            }
            return [string]$betaRelease.tag_name
        }
        catch {
            Throw-InstallError "Failed to read latest beta release tag from GitHub."
        }
    }
    else {
        $apiUrl = "https://api.github.com/repos/$repo/releases/latest"
        try {
            $release = Invoke-RestMethod -Uri $apiUrl -UseBasicParsing
            if (-not $release.tag_name) {
                Throw-InstallError 'Unable to determine the latest release tag.'
            }
            return [string]$release.tag_name
        }
        catch {
            Throw-InstallError "Failed to read latest release tag from GitHub."
        }
    }
}

function Get-ArchName {
    switch ($env:PROCESSOR_ARCHITECTURE) {
        'AMD64' { return 'amd64' }
        'x86_64' { return 'amd64' }
        'ARM64' { return 'arm64' }
        default { Throw-InstallError "Unsupported architecture: $($env:PROCESSOR_ARCHITECTURE)" }
    }
}

if (-not (Test-Path $binDir)) {
    New-Item -ItemType Directory -Force -Path $binDir | Out-Null
}

$archName = Get-ArchName

if ($Version) {
    $versionTag = [string]$Version
}
else {
    $versionTag = Get-LatestTag
}

if ($versionTag.StartsWith('v')) {
    $versionTag = $versionTag.Substring(1)
}

$assetName = "went-windows-$archName-v.$versionTag.exe"

$downloadUrl = "https://github.com/$repo/releases/download/v$versionTag/$assetName"

$tempFile = Join-Path $env:TEMP $assetName
Write-Info "Downloading $assetName"
Download-File -Url $downloadUrl -OutputPath $tempFile

$installPath = Join-Path $binDir $binaryName
Move-Item -Force $tempFile $installPath

$currentUserPath = [Environment]::GetEnvironmentVariable('Path', 'User')
$parts = @()
if ($currentUserPath) {
    $parts = $currentUserPath -split ';' | Where-Object { $_ -and $_ -ne $binDir }
}
$newUserPath = ($parts + $binDir) -join ';'
[Environment]::SetEnvironmentVariable('Path', $newUserPath, 'User')
$env:Path = "$binDir;$env:Path"

Write-Host "Kurulum başarılı! Aramıza hoşgeldin! Hemen başlamak için terminali yeniden açabilir veya şu komutu çalıştırabilirsin:"
Write-Host "  went --help"
