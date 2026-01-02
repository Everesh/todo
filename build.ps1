# PowerShell build script for Windows
# Usage: .\build.ps1 [all|current]

param(
    [string]$Target = "current"
)

$AppName = "todo"
$BuildDir = "build"
$Version = "dev"

# Try to get version from git
try {
    $Version = git describe --tags --always --dirty 2>$null
    if (-not $Version) { $Version = "dev" }
} catch {
    $Version = "dev"
}

$GoFlags = "-ldflags `"-X main.version=$Version`""

function Build-Platform {
    param(
        [string]$GOOS,
        [string]$GOARCH,
        [string]$OutputName
    )
    
    Write-Host "Building $AppName for $GOOS/$GOARCH..." -ForegroundColor Green
    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH
    go build $GoFlags -o "$BuildDir\$OutputName" ./cmd/todo
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Failed to build for $GOOS/$GOARCH" -ForegroundColor Red
        exit 1
    }
}

# Create build directory
if (-not (Test-Path $BuildDir)) {
    New-Item -ItemType Directory -Path $BuildDir | Out-Null
}

if ($Target -eq "all") {
    Write-Host "Building $AppName for all platforms..." -ForegroundColor Cyan
    
    Build-Platform "windows" "amd64" "$AppName-windows-amd64.exe"
    Build-Platform "windows" "386" "$AppName-windows-386.exe"
    Build-Platform "windows" "arm64" "$AppName-windows-arm64.exe"
    Build-Platform "darwin" "amd64" "$AppName-darwin-amd64"
    Build-Platform "darwin" "arm64" "$AppName-darwin-arm64"
    Build-Platform "linux" "amd64" "$AppName-linux-amd64"
    Build-Platform "linux" "386" "$AppName-linux-386"
    Build-Platform "linux" "arm64" "$AppName-linux-arm64"
    Build-Platform "linux" "arm" "$AppName-linux-arm"
    
    Write-Host "`nDone! Binaries are in $BuildDir\" -ForegroundColor Green
} else {
    Write-Host "Building $AppName for current platform..." -ForegroundColor Cyan
    $CurrentOS = (go env GOOS)
    $CurrentARCH = (go env GOARCH)
    $Extension = if ($CurrentOS -eq "windows") { ".exe" } else { "" }
    
    Build-Platform $CurrentOS $CurrentARCH "$AppName$Extension"
    Write-Host "`nDone! Binary is in $BuildDir\" -ForegroundColor Green
}

