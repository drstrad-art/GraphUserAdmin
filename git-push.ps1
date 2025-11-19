# Git Commit and Push Script

# Prompt for commit message
$commitMessage = Read-Host "Enter commit message"

# Check if commit message is empty
if ([string]::IsNullOrWhiteSpace($commitMessage)) {
    Write-Host "Error: Commit message cannot be empty!" -ForegroundColor Red
    exit 1
}

# Add all changes
Write-Host "Adding changes..." -ForegroundColor Yellow
git add .

# Commit changes
Write-Host "Committing changes..." -ForegroundColor Yellow
git commit -m $commitMessage

# Check if commit was successful
if ($LASTEXITCODE -eq 0) {
    # Push to remote
    Write-Host "Pushing to remote..." -ForegroundColor Yellow
    git push -u origin main
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Successfully pushed to remote!" -ForegroundColor Green
    } else {
        Write-Host "Failed to push to remote!" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "Failed to commit changes!" -ForegroundColor Red
    exit 1
}
