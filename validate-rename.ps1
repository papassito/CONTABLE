<#
.SYNOPSIS
    Script de validación para confirmar la correcta transición de 'go-skeleton' a 'Contable GO'.
.DESCRIPTION
    1. Verifica la existencia de 'Contable GO' y la ausencia de 'go-skeleton'.
    2. Asegura la integridad del archivo go.mod.
    3. Escanea el espacio de trabajo en busca de referencias residuales a 'go-skeleton'.
#>

[CmdletBinding()]
param(
    [string]$Path = $PSScriptRoot
)

if ([string]::IsNullOrEmpty($Path)) { $Path = (Get-Location).Path }

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host " 🔍 VERIFICADOR DE TRANSICIÓN A 'Contable GO'            " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

$oldDir = Join-Path -Path $Path -ChildPath "go-skeleton"
$newDir = Join-Path -Path $Path -ChildPath "Contable GO"
$errors = 0

# 1. Validar Directorio Anterior
Write-Host "`n📁 1. Comprobando estado del directorio anterior..." -ForegroundColor Yellow
if (Test-Path -Path $oldDir -PathType Container) {
    Write-Host "❌ Alerta: La carpeta anterior 'go-skeleton' todavía existe." -ForegroundColor Red
    Write-Host "   Se recomienda eliminarla o moverla para evitar duplicidad lógica." -ForegroundColor DarkYellow
    $errors++
} else {
    Write-Host "✅ Excelente: El directorio anterior 'go-skeleton' ya no existe físicamente." -ForegroundColor Green
}

# 2. Validar Directorio Nuevo
Write-Host "`n📁 2. Comprobando existencia de 'Contable GO'..." -ForegroundColor Yellow
if (-not (Test-Path -Path $newDir -PathType Container)) {
    Write-Host "❌ Error: No se encuentra la nueva carpeta 'Contable GO'." -ForegroundColor Red
    Write-Host "   Ejecute .\rename-folder.ps1 antes de ejecutar esta validación." -ForegroundColor DarkYellow
    $errors++
} else {
    Write-Host "✅ Confirmado: La carpeta 'Contable GO' está presente." -ForegroundColor Green
    
    # Validar go.mod
    $goModPath = Join-Path -Path $newDir -ChildPath "go.mod"
    if (Test-Path -Path $goModPath) {
        Write-Host "✅ Confirmado: 'go.mod' encontrado en '$newDir'." -ForegroundColor Green
    } else {
        Write-Host "❌ Error: Falta 'go.mod' dentro de '$newDir'." -ForegroundColor Red
        $errors++
    }
}

# 3. Escaneo de Referencias Residuales
Write-Host "`n🔍 3. Escaneando referencias residuales a 'go-skeleton'..." -ForegroundColor Yellow
$extensions = @("*.ps1", "*.json", "*.ts", "*.tsx", "*.yml", "*.md")
$residualMatches = @()

$filesToScan = Get-ChildItem -Path $Path -Recurse -Include $extensions -File | Where-Object {
    $_.FullName -notlike "*node_modules*" -and
    $_.FullName -notlike "*.git*" -and
    $_.FullName -notlike "*dist*" -and
    $_.FullName -notlike "*build*" -and
    $_.Name -neq "validate-rename.ps1" -and
    $_.Name -neq "rename-folder.ps1"
}

foreach ($file in $filesToScan) {
    try {
        $content = Get-Content -Path $file.FullName -Raw
        if ($content -match "go-skeleton") {
            $lineNum = 1
            Get-Content -Path $file.FullName | ForEach-Object {
                if ($_ -match "go-skeleton") {
                    $residualMatches += [PSCustomObject]@{
                        Archivo = $file.FullName.Replace($Path, ".")
                        Linea   = $lineNum
                        Texto   = $_.Trim()
                    }
                }
                $lineNum++
            }
        }
    } catch {
        # Ignorar errores de acceso a archivos bloqueados
    }
}

if ($residualMatches.Count -gt 0) {
    Write-Host "⚠️  Se detectaron $($residualMatches.Count) referencias residuales en el código:" -ForegroundColor Yellow
    $residualMatches | Format-Table -AutoSize
    Write-Host "💡 Sugerencia: Reemplace estas referencias por 'Contable GO' para evitar inconsistencias." -ForegroundColor Gray
} else {
    Write-Host "✅ ¡Limpieza perfecta! Cero referencias a 'go-skeleton' detectadas en los archivos analizados." -ForegroundColor Green
}

Write-Host "`n==========================================================" -ForegroundColor Cyan
if ($errors -eq 0) {
    Write-Host "🎉 VALIDACIÓN EXITOSA: La transición a 'Contable GO' es consistente." -ForegroundColor Green
    return 0
} else {
    Write-Host "❌ VALIDACIÓN CON ERRORES: Revise los puntos señalados arriba." -ForegroundColor Red
    return 1
}