<#
.SYNOPSIS
    Script de Tratamiento Médico para FCOS v2.2 Kernel.
.DESCRIPTION
    Sanea la capa de dominio eliminando contratos duplicados obsoletos y toma el pulso de salud del compilador.
#>

[CmdletBinding()]
param(
    [string]$TargetDir = (Get-Location).Path
)

# 1. Localización automática del módulo de Go
$candidateGoMod = Join-Path -Path $TargetDir -ChildPath "go-skeleton\go.mod"
if (Test-Path -Path $candidateGoMod) {
    $projectDir = Join-Path -Path $TargetDir -ChildPath "go-skeleton"
} elseif (Test-Path -Path (Join-Path -Path $TargetDir -ChildPath "go.mod")) {
    $projectDir = $TargetDir
} else {
    $projectDir = $TargetDir
}

Write-Host "==================================================================" -ForegroundColor Cyan
Write-Host " 💉 TRATAMIENTO DE REMEDIACIÓN MÉDICA - FCOS v2.2 KERNEL" -ForegroundColor Cyan
Write-Host "==================================================================" -ForegroundColor Cyan
Write-Host "📁 Proyecto objetivo: $projectDir`n" -ForegroundColor Gray

$domainDir = Join-Path -Path $projectDir -ChildPath "internal\domain"

$secretEnvelopePath = Join-Path -Path $domainDir -ChildPath "secret_envelope.go"
$vaultPath = Join-Path -Path $domainDir -ChildPath "vault.go"

# 2. Saneamiento de contratos duplicados en internal/domain para prevenir colisiones
if (Test-Path -Path $secretEnvelopePath) {
    Remove-Item -Path $secretEnvelopePath -Force
    Write-Host "🧹 Archivo obsoleto 'secret_envelope.go' removido para evitar colisiones de dominio con internal/secrets." -ForegroundColor Green
}

if (Test-Path -Path $vaultPath) {
    Remove-Item -Path $vaultPath -Force
    Write-Host "🧹 Archivo obsoleto 'vault.go' removido para evitar colisiones de dominio con internal/secrets." -ForegroundColor Green
}

# 4. Sincronización y Alta Médica
Push-Location -Path $projectDir
try {
    Write-Host "`n🔄 Sincronizando módulos con 'go mod tidy'..." -ForegroundColor Yellow
    go mod tidy

    Write-Host "🫀 Tomando el pulso al código con 'go vet ./...'..." -ForegroundColor Yellow
    $vetOutput = go vet ./... 2>&1

    if ($LASTEXITCODE -eq 0) {
        Write-Host "`n🎉 ¡ALTA MÉDICA OTORGADA! 'GO VET' FINALIZÓ CON CERO ERRORES." -ForegroundColor Green
    } else {
        Write-Host "`n⚠️ Se detectaron observaciones restantes en 'go vet':" -ForegroundColor Yellow
        Write-Host $vetOutput -ForegroundColor Red
    }
}
finally {
    Pop-Location
}

Write-Host "`n==================================================================" -ForegroundColor Cyan