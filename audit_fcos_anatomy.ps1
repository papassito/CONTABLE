<#
.SYNOPSIS
    Script de auditoría anatómica y milimétrica para FCOS v2.2 Kernel.
.DESCRIPTION
    Inspecciona todos los Bounded Contexts, verifica la integridad de los archivos .go
    y ejecuta el analizador estático nativo de Go ('go vet').
#>

[CmdletBinding()]
param(
    [switch]$Fix
)

# 1. Localización dinámica de la raíz del proyecto Go
$scriptRoot = $PSScriptRoot
$scriptRoot = $PSScriptRoot
if ([string]::IsNullOrWhiteSpace($scriptRoot)) { $scriptRoot = (Get-Location).Path }

$projectDir = $scriptRoot

if (-not $projectDir) {
    Write-Host "❌ Error Crítico: No se encontró un módulo Go válido con 'go.mod'." -ForegroundColor Red
    return
}

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host " 🩺 AUDITORÍA MILIMÉTRICA Y ANATÓMICA - FCOS v2.2 KERNEL  " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "📁 Directorio Auditado: $projectDir`n" -ForegroundColor Yellow
if ($Fix) {
    Write-Host "🔧 ¡MODO REPARACIÓN ACTIVO! Se eliminarán los archivos .go vacíos." -ForegroundColor Green
}

# 2. Matriz anatómica adaptable a Bounded Contexts
$anatomyMap = [ordered]@{
    "👁️  OJOS (Visibilidad y Telemetría)"    = @("health", "architecture", "telemetry")
    "🦵 PIERNAS (Persistencia y UoW)"         = @("database", "migrations", "repository")
    "🦾 BRAZOS (Ejecutores y RPA)"           = @("automation", "documents")
    "✋ MANOS (Interfaces HTTP y Entradas)"  = @("handler", "validator", "cmd/api")
    "🦶 PIES (Estabilidad y Pruebas)"        = @("_test.go", "tests")
    "👂 OÍDOS (Mensajería y Eventos)"         = @("messaging")
    "🧠 CEREBRO (Dominio y Reglas Contables)" = @("domain", "service", "application", "normative", "compliance", "calculation", "tenant", "audit", "identity", "secrets")
}

# 3. Lectura recursiva de todos los archivos .go del proyecto usando listas optimizadas
$allGoFiles = Get-ChildItem -Path $projectDir -Recurse -Filter "*.go"
$emptyFiles = [System.Collections.Generic.List[string]]::new()
$invalidPackageFiles = [System.Collections.Generic.List[string]]::new()

foreach ($file in $allGoFiles) {
    try {
        # Optimización con lectura nativa de .NET para máxima velocidad y resiliencia
        $content = [System.IO.File]::ReadAllText($file.FullName)
        if ([string]::IsNullOrWhiteSpace($content)) {
            $emptyFiles.Add($file.FullName)
        } elseif (-not ($content -match "(?m)^\s*package\s+\w+")) {
            $invalidPackageFiles.Add($file.FullName)
        }
    } catch {
        # Ignorar de forma segura si el archivo está bloqueado por otro proceso
    }
}

# 4. Reporte de Integridad Sintáctica de Archivos
Write-Host "📊 1. INTEGRIDAD DE ARCHIVOS FUENTE (.GO)" -ForegroundColor Yellow
Write-Host "----------------------------------------------------------"
Write-Host "Total de archivos .go detectados: $($allGoFiles.Count)" -ForegroundColor White

$fixedCount = 0
if ($emptyFiles.Count -gt 0) {
    Write-Host "❌ Archivos vacíos detectados: $($emptyFiles.Count)" -ForegroundColor Red
    foreach ($f in $emptyFiles) {
        Write-Host "   - $($f.Replace($projectDir, ''))" -ForegroundColor Red
    }
    if ($Fix) {
        Write-Host "   ↳ 🔧 Eliminando archivos vacíos..." -ForegroundColor Green
        foreach ($f in $emptyFiles) {
            try {
                Remove-Item -Path $f -Force -ErrorAction Stop
                $fixedCount++
            } catch {
                Write-Host "     - Error al eliminar $f" -ForegroundColor DarkRed
            }
        }
    }
} else {
    Write-Host "✅ Cero archivos vacíos detectados. Estructura sintáctica limpia." -ForegroundColor Green
}

if ($invalidPackageFiles.Count -gt 0) {
    Write-Host "⚠️ Archivos sin cabecera 'package' válida: $($invalidPackageFiles.Count)" -ForegroundColor Yellow
    foreach ($f in $invalidPackageFiles) {
        Write-Host "   - $($f.Replace($projectDir, ''))" -ForegroundColor Yellow
    }
} else {
    Write-Host "✅ Todos los archivos contienen la cabecera 'package' correspondiente." -ForegroundColor Green
}

if ($Fix -and $fixedCount -gt 0) {
    Write-Host "🔧 Total de archivos vacíos eliminados: $fixedCount" -ForegroundColor Green
}

# 5. Evaluación Anatómica de Subsistemas con rutas normalizadas
Write-Host "`n🧩 2. EVALUACIÓN DE SUBSISTEMAS ANATÓMICOS" -ForegroundColor Yellow
Write-Host "----------------------------------------------------------"

foreach ($system in $anatomyMap.Keys) {
    $keywords = $anatomyMap[$system]
    $matchedFiles = [System.Collections.Generic.List[string]]::new()

    foreach ($file in $allGoFiles) {
        # Normalizamos barras de Windows (\) a estilo Unix (/)
        $normalizedPath = $file.FullName.Replace("\", "/")
        foreach ($kw in $keywords) {
            # Uso correcto del operador -ilike (sin guion intermedio)
            if ($normalizedPath -ilike "*$kw*") {
                $matchedFiles.Add($file.FullName)
                break
            }
        }
    }

    $uniqueCount = @($matchedFiles | Select-Object -Unique).Count

    if ($uniqueCount -gt 0) {
        Write-Host "  $system : $uniqueCount componentes activos [OK]" -ForegroundColor Green
    } else {
        Write-Host "  $system : 0 componentes detectados [INCOMPLETO]" -ForegroundColor Red
    }
}

# 6. Ejecución del Análisis Estático de Código con 'go vet'
Write-Host "`n🔬 3. ANÁLISIS ESTÁTICO DE CÓDIGO (GO VET)" -ForegroundColor Yellow
Write-Host "----------------------------------------------------------"

Push-Location -Path $projectDir
try {
    # Asegurar el uso de la redirección correcta de PowerShell evitando escapes HTML
    $vetOutput = go vet ./... 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Análisis 'go vet' completado sin errores ni advertencias." -ForegroundColor Green
    } else {
        Write-Host "❌ Advertencias detectadas por 'go vet':" -ForegroundColor Red
        Write-Host $vetOutput -ForegroundColor Red
    }
}
catch {
    Write-Host "❌ No se pudo ejecutar 'go vet'. Verifica que Go esté correctamente instalado." -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}
finally {
    Pop-Location
}

Write-Host "`n==========================================================" -ForegroundColor Cyan
Write-Host "             FIN DE LA AUDITORÍA ANATÓMICA                " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan