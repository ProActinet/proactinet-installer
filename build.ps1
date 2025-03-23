# PowerShell script for cross-platform Go builds

# Define available platforms
$platforms = @(
    @{OS = "linux"; Arch = "amd64"; Display = "Linux (64-bit x86)"},
    @{OS = "linux"; Arch = "arm64"; Display = "Linux (64-bit ARM)"},
    @{OS = "windows"; Arch = "amd64"; Display = "Windows (64-bit x86)"},
    @{OS = "darwin"; Arch = "amd64"; Display = "macOS (64-bit x86)"},
    @{OS = "darwin"; Arch = "arm64"; Display = "macOS (64-bit ARM)"}
)

# Get the current directory name as the app name
$appName = Split-Path -Leaf (Get-Location)

# Create builds directory if it doesn't exist
$buildDir = "builds"
if (-not (Test-Path $buildDir)) {
    New-Item -ItemType Directory -Path $buildDir | Out-Null
}

function Build-GoBinary {
    param (
        [string]$OS,
        [string]$Arch,
        [string]$OutputName
    )

    Write-Host "`nBuilding for $OS/$Arch..." -ForegroundColor Cyan

    # Set environment variables for cross-compilation
    $env:GOOS = $OS
    $env:GOARCH = $Arch

    # Build the binary
    try {
        $output = go build -o $OutputName 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "Successfully built binary: $OutputName" -ForegroundColor Green
        } else {
            Write-Host "Build failed: $output" -ForegroundColor Red
        }
    }
    catch {
        Write-Host "Error during build: $_" -ForegroundColor Red
    }
}

# Display menu
Write-Host "Available platforms:" -ForegroundColor Yellow
for ($i = 0; $i -lt $platforms.Count; $i++) {
    Write-Host "$($i + 1). $($platforms[$i].Display)"
}
Write-Host "A. Build for all platforms"
Write-Host "Q. Quit"

# Get user choice
$choice = Read-Host "`nEnter your choice (1-$($platforms.Count), A for all, Q to quit)"

switch ($choice.ToUpper()) {
    "Q" {
        Write-Host "Exiting script..." -ForegroundColor Yellow
        exit
    }
    "A" {
        Write-Host "`nBuilding for all platforms..." -ForegroundColor Cyan
        foreach ($platform in $platforms) {
            $extension = if ($platform.OS -eq "windows") { ".exe" } else { "" }
            $outputPath = Join-Path $buildDir "$appName`_$($platform.OS)_$($platform.Arch)$extension"
            Build-GoBinary -OS $platform.OS -Arch $platform.Arch -OutputName $outputPath
        }
    }
    default {
        $index = [int]$choice - 1
        if ($index -ge 0 -and $index -lt $platforms.Count) {
            $platform = $platforms[$index]
            $extension = if ($platform.OS -eq "windows") { ".exe" } else { "" }
            $outputPath = Join-Path $buildDir "$appName`_$($platform.OS)_$($platform.Arch)$extension"
            Build-GoBinary -OS $platform.OS -Arch $platform.Arch -OutputName $outputPath
        }
        else {
            Write-Host "Invalid choice. Please run the script again." -ForegroundColor Red
        }
    }
}

# Display build directory contents after completion
Write-Host "`nContents of builds directory:" -ForegroundColor Yellow
Get-ChildItem $buildDir | Format-Table Name, Length