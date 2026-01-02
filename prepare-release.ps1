# Script to prepare release binaries with simple names
# Creates a releases/ directory with platform-specific folders

$ReleasesDir = "releases"
if (Test-Path $ReleasesDir) {
    Remove-Item -Recurse -Force $ReleasesDir
}
New-Item -ItemType Directory -Path $ReleasesDir | Out-Null

# Windows
New-Item -ItemType Directory -Path "$ReleasesDir\windows-amd64" | Out-Null
Copy-Item "build\todo-windows-amd64.exe" "$ReleasesDir\windows-amd64\todo.exe"

New-Item -ItemType Directory -Path "$ReleasesDir\windows-386" | Out-Null
Copy-Item "build\todo-windows-386.exe" "$ReleasesDir\windows-386\todo.exe"

New-Item -ItemType Directory -Path "$ReleasesDir\windows-arm64" | Out-Null
Copy-Item "build\todo-windows-arm64.exe" "$ReleasesDir\windows-arm64\todo.exe"

# macOS
New-Item -ItemType Directory -Path "$ReleasesDir\darwin-amd64" | Out-Null
Copy-Item "build\todo-darwin-amd64" "$ReleasesDir\darwin-amd64\todo"
# Make executable
icacls "$ReleasesDir\darwin-amd64\todo" /grant Everyone:RX

New-Item -ItemType Directory -Path "$ReleasesDir\darwin-arm64" | Out-Null
Copy-Item "build\todo-darwin-arm64" "$ReleasesDir\darwin-arm64\todo"
icacls "$ReleasesDir\darwin-arm64\todo" /grant Everyone:RX

# Linux
New-Item -ItemType Directory -Path "$ReleasesDir\linux-amd64" | Out-Null
Copy-Item "build\todo-linux-amd64" "$ReleasesDir\linux-amd64\todo"

New-Item -ItemType Directory -Path "$ReleasesDir\linux-386" | Out-Null
Copy-Item "build\todo-linux-386" "$ReleasesDir\linux-386\todo"

New-Item -ItemType Directory -Path "$ReleasesDir\linux-arm64" | Out-Null
Copy-Item "build\todo-linux-arm64" "$ReleasesDir\linux-arm64\todo"

New-Item -ItemType Directory -Path "$ReleasesDir\linux-arm" | Out-Null
Copy-Item "build\todo-linux-arm" "$ReleasesDir\linux-arm\todo"

Write-Host "Release binaries prepared in $ReleasesDir/" -ForegroundColor Green
Write-Host "Each platform folder contains a simple 'todo' or 'todo.exe' binary" -ForegroundColor Green

