# Tomba CLI installer for Windows
# Usage 1: irm https://raw.githubusercontent.com/tomba-io/tomba/master/res/package/scripts/install.ps1 | iex
# Usage 2: irm https://releases.tomba.io/install.ps1 | iex
# Options:
#   -Version    Install a specific version (e.g., -Version "1.1.5")
#   -InstallDir Custom install directory

param(
    [string]$Version = "",
    [string]$InstallDir = ""
)

$ErrorActionPreference = "Stop"

$REPO = "tomba-io/tomba"
$BINARY_NAME = "tomba"

function Write-Info($msg) { Write-Host "i " -ForegroundColor Cyan -NoNewline; Write-Host $msg }
function Write-Success($msg) { Write-Host "✓ " -ForegroundColor Green -NoNewline; Write-Host $msg }
function Write-Warn($msg) { Write-Host "! " -ForegroundColor Yellow -NoNewline; Write-Host $msg }
function Write-Err($msg) { Write-Host "x " -ForegroundColor Red -NoNewline; Write-Host $msg; exit 1 }

Write-Host ""
Write-Host "    ████████╗ ██████╗ ███╗   ███╗██████╗  █████╗    ██╗ ██████╗ " -ForegroundColor Green
Write-Host "    ╚══██╔══╝██╔═══██╗████╗ ████║██╔══██╗██╔══██╗   ██║██╔═══██╗" -ForegroundColor Green
Write-Host "       ██║   ██║   ██║██╔████╔██║██████╔╝███████║   ██║██║   ██║" -ForegroundColor Green
Write-Host "       ██║   ██║   ██║██║╚██╔╝██║██╔══██╗██╔══██║   ██║██║   ██║" -ForegroundColor Green
Write-Host "       ██║   ╚██████╔╝██║ ╚═╝ ██║██████╔╝██║  ██║██╗██║╚██████╔╝" -ForegroundColor Green
Write-Host "       ╚═╝    ╚═════╝ ╚═╝     ╚═╝╚═════╝ ╚═╝  ╚═╝╚═╝╚═╝ ╚═════╝ " -ForegroundColor Green
Write-Host ""
Write-Host "    CLI Installer" -ForegroundColor White -NoNewline
Write-Host " — search or verify email addresses in seconds" -ForegroundColor DarkGray
Write-Host ""

# Detect architecture
$ARCH = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

# Check for ARM
$cpuArch = (Get-CimInstance -ClassName Win32_Processor).Architecture
if ($cpuArch -eq 12) { $ARCH = "arm64" }

Write-Info "Detected: windows/$ARCH"

# Get version
if ($Version -eq "") {
    try {
        $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest"
        $Version = $release.tag_name -replace '^v', ''
    } catch {
        Write-Err "Failed to get latest version: $_"
    }
    Write-Info "Latest version: v$Version"
} else {
    Write-Info "Requested version: v$Version"
}

# Check if already installed and up to date
if ($InstallDir -eq "") {
    $INSTALL_DIR = Join-Path $env:LOCALAPPDATA "tomba"
} else {
    $INSTALL_DIR = $InstallDir
}

$installPath = Join-Path $INSTALL_DIR "$BINARY_NAME.exe"

if (Test-Path $installPath) {
    try {
        $currentVer = & $installPath version 2>&1 | Select-String -Pattern '\d+\.\d+\.\d+' | ForEach-Object { $_.Matches[0].Value }
        if ($currentVer -eq $Version) {
            Write-Success "tomba v$Version is already installed and up to date"
            exit 0
        }
        if ($currentVer) {
            Write-Info "Upgrading from v$currentVer to v$Version"
        }
    } catch {
        # Could not determine current version, proceed with install
    }
}

# Build download URL
$FILENAME = "${BINARY_NAME}_windows_${ARCH}.zip"
$DOWNLOAD_URL = "https://github.com/$REPO/releases/download/v$Version/$FILENAME"
$CHECKSUMS_URL = "https://github.com/$REPO/releases/download/v$Version/${BINARY_NAME}_checksums.txt"

Write-Info "Downloading $FILENAME..."

# Create temp directory
$TMP_DIR = Join-Path $env:TEMP "tomba-install-$(Get-Random)"
New-Item -ItemType Directory -Path $TMP_DIR -Force | Out-Null

try {
    # Download archive
    $zipPath = Join-Path $TMP_DIR $FILENAME
    Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $zipPath -UseBasicParsing

    if (-not (Test-Path $zipPath)) {
        Write-Err "Download failed"
    }

    Write-Success "Downloaded successfully"

    # Download and verify checksum
    $checksumsPath = Join-Path $TMP_DIR "checksums.txt"
    try {
        Invoke-WebRequest -Uri $CHECKSUMS_URL -OutFile $checksumsPath -UseBasicParsing
        if (Test-Path $checksumsPath) {
            $expectedLine = Get-Content $checksumsPath | Where-Object { $_ -match $FILENAME }
            if ($expectedLine) {
                $expectedHash = ($expectedLine -split '\s+')[0]
                $actualHash = (Get-FileHash -Path $zipPath -Algorithm SHA256).Hash.ToLower()
                if ($expectedHash -eq $actualHash) {
                    Write-Success "Checksum verified (SHA256)"
                } else {
                    Write-Err "Checksum mismatch! Expected: $expectedHash, Got: $actualHash"
                }
            } else {
                Write-Warn "Checksum not found for $FILENAME, skipping verification"
            }
        }
    } catch {
        Write-Warn "Could not download checksums, skipping verification"
    }

    # Extract
    Write-Info "Extracting..."
    Expand-Archive -Path $zipPath -DestinationPath $TMP_DIR -Force

    $binaryPath = Join-Path $TMP_DIR "$BINARY_NAME.exe"
    if (-not (Test-Path $binaryPath)) {
        Write-Err "Binary not found in archive"
    }

    Write-Success "Extracted successfully"

    # Install directory
    if (-not (Test-Path $INSTALL_DIR)) {
        New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
    }

    Copy-Item -Path $binaryPath -Destination $installPath -Force

    Write-Success "Installed tomba v$Version to $installPath"

    # Add to PATH if not already there
    $currentPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
    if ($currentPath -notlike "*$INSTALL_DIR*") {
        [System.Environment]::SetEnvironmentVariable("Path", "$currentPath;$INSTALL_DIR", "User")
        $env:Path = "$env:Path;$INSTALL_DIR"
        Write-Success "Added $INSTALL_DIR to user PATH"
        Write-Warn "Restart your terminal for PATH changes to take effect"
    }

    # Verify
    try {
        $ver = & $installPath version 2>&1
        Write-Success "Verification: $ver"
    } catch {
        Write-Warn "Could not verify installation"
    }

    # Setup PowerShell completions
    try {
        $completionScript = & $installPath completion powershell 2>&1
        if ($LASTEXITCODE -eq 0 -and $completionScript) {
            $profileDir = Split-Path $PROFILE -Parent
            if (-not (Test-Path $profileDir)) {
                New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
            }
            $completionFile = Join-Path $profileDir "tomba.completion.ps1"
            $completionScript | Out-File -FilePath $completionFile -Encoding UTF8

            # Add sourcing line to profile if not present
            $sourceLine = ". `"$completionFile`""
            if (-not (Test-Path $PROFILE) -or -not (Select-String -Path $PROFILE -Pattern "tomba\.completion\.ps1" -Quiet)) {
                Add-Content -Path $PROFILE -Value "`n$sourceLine"
            }
            Write-Success "PowerShell completions installed to $completionFile"
        }
    } catch {
        Write-Warn "Could not set up PowerShell completions. Run: tomba completion powershell --help"
    }

    Write-Host ""
    Write-Host "Getting started:" -ForegroundColor White
    Write-Host "  tomba login          Sign in with your API key"
    Write-Host "  tomba search -t      Search emails by domain"
    Write-Host "  tomba chat           Start AI chat"
    Write-Host "  tomba --help         Show all commands"
    Write-Host ""

} finally {
    # Cleanup
    Remove-Item -Path $TMP_DIR -Recurse -Force -ErrorAction SilentlyContinue
}
