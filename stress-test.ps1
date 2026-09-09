<#
.SYNOPSIS
    Script de Pruebas de Rompimiento y Estrés para FCOS v2.2 Kernel.
.DESCRIPTION
    Ejecuta pruebas de concurrencia con el detector de carreras de Go y
    realiza múltiples iteraciones para detectar fallos bajo carga.
#>

[CmdletBinding()]
param(
    [int]$Iteraciones = 5,
    [int]$Hilos = 8
)

# 1. Localizar el script principal de pruebas
$scriptRoot = $PSScriptRoot
if ([string]::IsNullOrWhiteSpace($scriptRoot)) { $scriptRoot = (Get-Location).Path }
$testRunnerScript = Join-Path -Path $scriptRoot -ChildPath "run-tests.ps1"

if (-not (Test-Path -Path $testRunnerScript)) {
    Write-Host "❌ Error: No se encontró el script 'run-tests.ps1'. Asegúrate de que esté en el mismo directorio." -ForegroundColor Red
    return
}

Write-Host "==================================================================" -ForegroundColor Red
Write-Host " 💥 PRUEBA DE ROMPIMIENTO DE LÓGICA Y ESTRÉS - FCOS v2.2 KERNEL  " -ForegroundColor Red
Write-Host "==================================================================" -ForegroundColor Red
Write-Host "📁 Proyecto objetivo : Raíz del Repositorio ($scriptRoot)"
Write-Host "🔄 Iteraciones      : $Iteraciones"
Write-Host "⚡ Hilos Paralelos   : $Hilos`n"

# 2. Fase 1: Detección de Condiciones de Carrera (Race Detector)
Write-Host "🏃 1. Evaluando integridad de memoria y condiciones de carrera (-race)..."
# Corrección: Habilitar CGO para que el detector de carreras funcione.
$env:CGO_ENABLED = 1
& $testRunnerScript -race -test.parallel $Hilos

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ [FALLO] Se detectaron condiciones de carrera en el Kernel." -ForegroundColor Red
    Write-Host "`n⚠️  SE DETECTÓ RUPTURA LÓGICA O FALLO DE CONCURRENCIA BAJO CARGA." -ForegroundColor Red
    Write-Host "`n==================================================================" -ForegroundColor Red
    # Restaurar variable de entorno y salir
    $env:CGO_ENABLED = 0
    return
} else {
    Write-Host "✅ [ÉXITO] No se detectaron condiciones de carrera." -ForegroundColor Green
}
$env:CGO_ENABLED = 0 # Deshabilitar CGO para las siguientes pruebas si no es necesario

# 3. Fase 2: Pruebas de Estrés por Iteración
Write-Host "`n🔥 2. Ejecutando pruebas de estrés y saturación ($Iteraciones iteraciones)..."
# Corrección: Pasar el argumento -count correctamente.
& $testRunnerScript -count=$Iteraciones -test.parallel $Hilos

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ [FALLO] Las pruebas de estrés fallaron en alguna de las $Iteraciones iteraciones." -ForegroundColor Red
    Write-Host "`n⚠️  SE DETECTÓ RUPTURA LÓGICA O FALLO DE CONCURRENCIA BAJO CARGA." -ForegroundColor Red
} else {
    Write-Host "✅ [ÉXITO] Las $Iteraciones iteraciones de las pruebas de estrés se completaron sin errores." -ForegroundColor Green
}

Write-Host "`n==================================================================" -ForegroundColor Red