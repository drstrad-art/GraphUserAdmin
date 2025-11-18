# Build script for Windows PowerShell

Write-Host "Building GraphUserAdmin (gua.exe)..." -ForegroundColor Cyan

$outputName = "./bin/gua.exe"
$mainPath = "./cmd/msgraph"

go build -o $outputName $mainPath

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Build successful!" -ForegroundColor Green
    Write-Host "Executable created: $outputName" -ForegroundColor Green
    Write-Host ""
    Write-Host "To run the application:" -ForegroundColor Yellow
    Write-Host "  .\$outputName --help" -ForegroundColor White
} else {
    Write-Host "✗ Build failed!" -ForegroundColor Red
    exit 1
}
