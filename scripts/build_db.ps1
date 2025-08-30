# scripts/build_db.ps1
$Root = (Resolve-Path "$PSScriptRoot\..").Path
$Mig  = Join-Path $Root "migrations"
$DB   = Join-Path $Root "app.db"

if (-not (Get-Command goose -ErrorAction SilentlyContinue)) {
  Write-Error "goose not found. Install: go install github.com/pressly/goose/v3/cmd/goose@latest"
  exit 1
}

Remove-Item -Force $DB -ErrorAction SilentlyContinue
goose -dir $Mig sqlite3 $DB up
Write-Host "Built $DB"
