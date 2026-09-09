<#
.SYNOPSIS
    Script para renombrar físicamente go-skeleton a 'Contable GO' de forma segura.
.DESCRIPTION
    Valida pre-condiciones, evita colisiones de nombres, advierte sobre bloqueos de archivos
    y realiza el renombrado atómico del directorio.
#>

[CmdletBinding()]
param(
    [string]$Path = $PSScriptRoot
)

if ([string]::IsNullOrEmpty($Path)) { $Path = (Get-Location).Path }

$sourceDir = Join-Path -Path $Path -ChildPath "go-skeleton"
$targetDir = Join-Path -Path $Path -ChildPath "Contable GO"

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host " 📂 RENOMBRADO SEGURO DE DIRECTORIO CONTABLE             " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

# 1. Verificar si el origen existe
if (-not (Test-Path -Path $sourceDir -PathType Container)) {
    if (Test-Path -Path $targetDir -PathType Container) {
        Write-Host "ℹ️ El directorio ya ha sido renombrado a: '$targetDir'" -ForegroundColor Yellow
        return 0
    }
    Write-Host "❌ Error: No se encontró el directorio de origen: '$sourceDir'" -ForegroundColor Red
    return 1
}

# 2. Verificar colisión con el destino
if (Test-Path -Path $targetDir -PathType Container) {
    Write-Host "❌ Error: El directorio de destino ya existe: '$targetDir'" -ForegroundColor Red
    Write-Host "Asegúrese de borrarlo o fusionarlo manualmente antes de continuar." -ForegroundColor Yellow
    return 1
}

# 3. Advertencia de procesos abiertos (IDE, Terminal, Wails)
Write-Host "⚠️  Asegúrese de cerrar editores de código (VS Code, GoLand) o terminales" -ForegroundColor Yellow
Write-Host "   que puedan tener bloqueada la carpeta 'go-skeleton'." -ForegroundColor Yellow
Read-Host "Presione ENTER para continuar con el renombrado o CTRL+C para cancelar"

# 4. Renombrado seguro
try {
    Write-Host "`n🔄 Renombrando '$sourceDir' a 'Contable GO'..." -ForegroundColor Yellow
    Rename-Item -Path $sourceDir -NewName "Contable GO" -ErrorAction Stop
    Write-Host "✅ Directorio renombrado con éxito." -ForegroundColor Green
}
catch {
    Write-Host "❌ Error crítico al renombrar la carpeta: $_" -ForegroundColor Red
    Write-Host "Intente ejecutar este script como Administrador o cierre aplicaciones activas." -ForegroundColor Yellow
    return 1
}

# 5. Validación automática del resultado
$validationScript = Join-Path -Path $Path -ChildPath "validate-rename.ps1"
if (Test-Path $validationScript) {
    Write-Host "`n🔍 Iniciando validación automática de consistencia..." -ForegroundColor Yellow
    & $validationScript -Path $Path
}

Write-Host "`n🎉 ¡Proceso finalizado con éxito!" -ForegroundColor Green
Write-Host "==========================================================" -ForegroundColor Cyan
return 0