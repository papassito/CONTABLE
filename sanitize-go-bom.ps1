<#
.SYNOPSIS
    Script de Eliminación Masiva de UTF-8 BOM para Proyectos Go.
.DESCRIPTION
    Escanea recursivamente los archivos .go, remueve los 3 bytes iniciales (0xEF 0xBB 0xBF)
    si existen y recompila la verificación de sintaxis con 'go vet'.
#>

[CmdletBinding()]
param(
    [string]$TargetDir = (Get-Location).Path
)

# 1. Localización dinámica del directorio del proyecto Go
$candidateGoMod = Join-Path -Path $TargetDir -ChildPath "go-skeleton\go.mod"
if (Test-Path -Path $candidateGoMod) {
    $projectDir = Join-Path -Path $TargetDir -ChildPath "go-skeleton"
} elseif (Test-Path -Path (Join-Path -Path $TargetDir -ChildPath "go.mod")) {
    $projectDir = $TargetDir
} else {
    $projectDir = $TargetDir
}

Write-Host "==================================================================" -ForegroundColor Cyan
Write-Host " 🧹 ELIMINADOR DE MARCAS DE ORDEN DE BYTES (UTF-8 BOM) - GO       " -ForegroundColor Cyan
Write-Host "==================================================================" -ForegroundColor Cyan
Write-Host "📁 Proyecto analizado: $projectDir`n" -ForegroundColor Gray

# 2. Buscar todos los archivos .go
$goFiles = Get-ChildItem -Path $projectDir -Recurse -Filter "*.go"
$cleanedCount = 0

foreach ($file in $goFiles) {
    # Optimización: Leer solo los primeros 3 bytes para detectar el BOM
    $stream = $null
    $hasBom = $false
    try {
        $stream = New-Object System.IO.FileStream($file.FullName, [System.IO.FileMode]::Open, [System.IO.FileAccess]::Read, [System.IO.FileShare]::ReadWrite)
        $bomBytes = New-Object byte[] 3
        $bytesRead = $stream.Read($bomBytes, 0, 3)
        if ($bytesRead -ge 3 -and $bomBytes[0] -eq 0xEF -and $bomBytes[1] -eq 0xBB -and $bomBytes[2] -eq 0xBF) {
            $hasBom = $true
        }
    } catch {
        Write-Host "❌ Error al verificar BOM en $($file.Name): $_" -ForegroundColor Red
    } finally {
        if ($stream -ne $null) {
            $stream.Close()
            $stream.Dispose()
        }
    }

    # Evaluar si los primeros 3 bytes son 0xEF, 0xBB, 0xBF
    if ($hasBom) {
        
        # Robustez: Usar los manejadores de texto de .NET para reescribir sin BOM
        $content = [System.IO.File]::ReadAllText($file.FullName) # ReadAllText maneja y omite el BOM automáticamente
        $utf8NoBom = New-Object System.Text.UTF8Encoding($false) # Especificar UTF-8 sin BOM
        [System.IO.File]::WriteAllText($file.FullName, $content, $utf8NoBom)
        
        $relativePath = $file.FullName.Replace($projectDir, "")
        Write-Host "✅ [BOM REMOVIDO] $relativePath" -ForegroundColor Green
        $cleanedCount++
    }
}

Write-Host "`n------------------------------------------------------------------" -ForegroundColor DarkGray
if ($cleanedCount -gt 0) {
    Write-Host "🎉 Se limpió la marca UTF-8 BOM en $cleanedCount archivo(s)." -ForegroundColor Green
} else {
    Write-Host "✅ Todos los archivos .go están libres de marcas UTF-8 BOM." -ForegroundColor Green
}

# 3. Validar sintaxis con go vet (sin cambios)
Push-Location -Path $projectDir
try {
    Write-Host "`n🔬 Validando lectura de código con 'go vet ./...'..." -ForegroundColor Yellow
    $vetOutput = go vet ./... 2>&1

    if ($LASTEXITCODE -eq 0) {
        Write-Host "🎉 ¡ANÁLISIS COMPLETADO! Cero errores de codificación reportados." -ForegroundColor Green
        $global:LASTEXITCODE = 0
    } else {
        Write-Host "⚠️ Advertencias detectadas por 'go vet':" -ForegroundColor Yellow
        Write-Host $vetOutput -ForegroundColor Red
        $global:LASTEXITCODE = 1
    }
}
finally {
    Pop-Location
}

Write-Host "`n==================================================================" -ForegroundColor Cyan