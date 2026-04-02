# Tomba CLI installer for Windows
# Usage 1: irm https://raw.githubusercontent.com/tomba-io/tomba/master/res/package/scripts/install.ps1 | iex
# Usage 2: irm https://releases.tomba.io/install.ps1 | iex

$ErrorActionPreference = "Stop"

$REPO = "tomba-io/tomba"
$BINARY_NAME = "tomba"

function Write-Info($msg) { Write-Host "i " -ForegroundColor Cyan -NoNewline; Write-Host $msg }
function Write-Success($msg) { Write-Host "✓ " -ForegroundColor Green -NoNewline; Write-Host $msg }
function Write-Warn($msg) { Write-Host "! " -ForegroundColor Yellow -NoNewline; Write-Host $msg }
function Write-Error($msg) { Write-Host "X " -ForegroundColor Red -NoNewline; Write-Host $msg; exit 1 }

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

# Get latest version
try {
    $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest"
    $VERSION = $release.tag_name -replace '^v', ''
} catch {
    Write-Error "Failed to get latest version: $_"
}

Write-Info "Latest version: v$VERSION"

# Build download URL
$FILENAME = "${BINARY_NAME}_windows_${ARCH}.zip"
$DOWNLOAD_URL = "https://github.com/$REPO/releases/download/v$VERSION/$FILENAME"

Write-Info "Downloading $FILENAME..."

# Create temp directory
$TMP_DIR = Join-Path $env:TEMP "tomba-install-$(Get-Random)"
New-Item -ItemType Directory -Path $TMP_DIR -Force | Out-Null

try {
    # Download
    $zipPath = Join-Path $TMP_DIR $FILENAME
    Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $zipPath -UseBasicParsing

    if (-not (Test-Path $zipPath)) {
        Write-Error "Download failed"
    }

    Write-Success "Downloaded successfully"

    # Extract
    Write-Info "Extracting..."
    Expand-Archive -Path $zipPath -DestinationPath $TMP_DIR -Force

    $binaryPath = Join-Path $TMP_DIR "$BINARY_NAME.exe"
    if (-not (Test-Path $binaryPath)) {
        Write-Error "Binary not found in archive"
    }

    Write-Success "Extracted successfully"

    # Install directory
    $INSTALL_DIR = Join-Path $env:LOCALAPPDATA "tomba"
    if (-not (Test-Path $INSTALL_DIR)) {
        New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
    }

    $installPath = Join-Path $INSTALL_DIR "$BINARY_NAME.exe"
    Copy-Item -Path $binaryPath -Destination $installPath -Force

    Write-Success "Installed tomba v$VERSION to $installPath"

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
