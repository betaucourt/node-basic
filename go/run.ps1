# PowerShell script to run the Go music application

Write-Host "Starting Go music application..." -ForegroundColor Green

# Change to the Go project directory
Set-Location $PSScriptRoot

# Check if the executable exists
if (-not (Test-Path "music-app.exe")) {
    Write-Host "Error: music-app.exe not found!" -ForegroundColor Red
    Write-Host "Please run build.ps1 first to build the application" -ForegroundColor Yellow
    exit 1
}

# Display application info
Write-Host "Found music-app.exe" -ForegroundColor Cyan
$fileSize = (Get-Item "music-app.exe").Length
Write-Host "File size: $($fileSize) bytes" -ForegroundColor Cyan

# Start the application
Write-Host "Starting server on port 8081..." -ForegroundColor Blue
Write-Host "API endpoint available at: http://localhost:8081/test" -ForegroundColor Yellow
Write-Host "Press Ctrl+C to stop the server" -ForegroundColor Magenta
Write-Host "----------------------------------------" -ForegroundColor Gray

try {
    # Run the application
    ./music-app.exe
} catch {
    Write-Host "Error running application: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}