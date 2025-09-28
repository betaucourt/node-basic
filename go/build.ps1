# PowerShell script to build the Go music application

Write-Host "Building Go music application..." -ForegroundColor Green

# Change to the Go project directory
Set-Location $PSScriptRoot

# Check if Go is installed
try {
    $goVersion = go version
    Write-Host "Using $goVersion" -ForegroundColor Cyan
} catch {
    Write-Host "Error: Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host "Please install Go from https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}

# Clean up previous build
if (Test-Path "music-app.exe") {
    Remove-Item "music-app.exe" -Force
    Write-Host "Removed previous build" -ForegroundColor Yellow
}

# Build the application
Write-Host "Compiling..." -ForegroundColor Blue
try {
    go build -o music-app.exe main.go
    
    if (Test-Path "music-app.exe") {
        Write-Host "Build successful! Generated music-app.exe" -ForegroundColor Green
        $fileSize = (Get-Item "music-app.exe").Length
        Write-Host "File size: $($fileSize) bytes" -ForegroundColor Cyan
    } else {
        Write-Host "Build failed - executable not found" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "Build failed with error: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Write-Host "Build completed successfully!" -ForegroundColor Green