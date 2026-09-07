<#
.SYNOPSIS
    Buscador de anomalías de codificación, Mojibake y UTF-8 BOM.
.DESCRIPTION
    Escanea archivos de código fuente en busca de marcas de orden de bytes (BOM),
    caracteres de control corruptos y mojibakes comunes.
#>

[CmdletBinding()]
param(
    [string]$Path = ".",
    [string[]]$Extensions = @("*.go", "*.ts", "*.tsx", "*.json", "*.yaml", "*.yml", "*.md", "*.ps1"),
    [string[]]$ExcludeDirs = @("node_modules", ".git", "dist", "build", "wailsjs"),
    [switch]$Fix
)

# Obtener ruta absoluta
$targetPath = Convert-Path $Path
Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "     🔎 BUSCADOR DE CARACTERES RAROS Y CORRUPTOS          " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan
if ($Fix) {
    Write-Host "🔧 ¡MODO REPARACIÓN ACTIVO! Se corregirán las anomalías encontradas." -ForegroundColor Green
}
Write-Host "📂 Escaneando: $targetPath`n" -ForegroundColor Yellow

# Obtener archivos respetando exclusiones
$files = Get-ChildItem -Path $targetPath -Recurse -Include $Extensions -File -ErrorAction SilentlyContinue | ForEach-Object {
    # Normalizar la ruta para que la coincidencia sea independiente del SO
    $normalizedPath = $_.FullName.Replace('\', '/')

    # Verificar si la ruta contiene alguno de los directorios a excluir
    $isExcluded = $false
    foreach ($dir in $ExcludeDirs) {
        if ($normalizedPath -like "*/$dir/*") {
            $isExcluded = $true
            break
        }
    }

    # Devolver el archivo solo si no está en un directorio excluido y no es el script mismo
    if (-not $isExcluded -and $_.Name -ne 'buscar_caracteres_raros.ps1') {
        $_
    }
}

Write-Host "📝 Archivos a escanear: $($files.Count)" -ForegroundColor Gray
Write-Host "----------------------------------------------------------"

$bomCount = 0
$corruptCount = 0
$fixedCount = 0

# Tabla de traducciones comunes de Mojibake
$mojibakeReplacements = @{
    "transacciÃ³n" = "transacción"
    "provocarÃ¡"   = "provocará"
    "ejecuciÃ³n"  = "ejecución"
    "transacciÃ³"  = "transacción"
    "â"          = "—"
}

foreach ($file in $files) {
    $relativePath = $file.FullName.Replace($targetPath, "")
    $fileContentChanged = $false

    # 1. Detectar marcas UTF-8 BOM leyendo los primeros 3 bytes reales
    try {
        $stream = New-Object System.IO.FileStream($file.FullName, [System.IO.FileMode]::Open, [System.IO.FileAccess]::Read)
        $bytes = New-Object byte[] 3
        $bytesRead = $stream.Read($bytes, 0, 3)
        $stream.Close()

        if ($bytesRead -eq 3 -and $bytes[0] -eq 0xEF -and $bytes[1] -eq 0xBB -and $bytes[2] -eq 0xBF) {
            Write-Host "⚠️  [UTF-8 BOM] $relativePath" -ForegroundColor Yellow
            $bomCount++
            
            if ($Fix) {
                # Leer archivo ignorando el BOM y marcar para guardar sin BOM
                $content = [System.IO.File]::ReadAllText($file.FullName)
                $fileContentChanged = $true
                Write-Host "   ↳ 🔧 Removiendo marca BOM..." -ForegroundColor Green
            } else {
                Write-Host "   ↳ Explicación: Posee marca de orden de bytes al inicio (común en Windows, problemático en Go)." -ForegroundColor DarkYellow
            }
        }
    } catch {
        Write-Host "❌ Error al leer bytes de $($file.Name): $_" -ForegroundColor Red
        continue
    }

    # 2. Análisis línea por línea en busca de Mojibake o caracteres de control
    $lines = [System.IO.File]::ReadAllLines($file.FullName)
    $newLines = [System.Collections.Generic.List[string]]::new()
    
    for ($i = 0; $i -lt $lines.Length; $i++) {
        $line = $lines[$i]
        $lineNum = $i + 1

        # Buscar secuencias rotas de Mojibake como â, â, Ã, etc.
        if ($line -match "â|â€|Ã[^\s]|ï»¿") {
            Write-Host "❌ [Mojibake] $relativePath (Línea $lineNum)" -ForegroundColor Red
            Write-Host "   ↳ Texto corrupto: '$($line.Trim())'" -ForegroundColor DarkRed
            $corruptCount++
            
            if ($Fix) {
                foreach ($key in $mojibakeReplacements.Keys) {
                    if ($line -contains $key) {
                        $line = $line.Replace($key, $mojibakeReplacements[$key])
                    }
                }
                # Limpiar cualquier residuo de marcas de control mojibake comunes
                $line = $line -replace "â", "—"
                $line = $line -replace "â", ""
                $fileContentChanged = $true
            }
        }

        # Buscar caracteres de control ilegales (excluyendo tab, CR, LF normales)
        if ($line -match "[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]") {
            Write-Host "⚠️  [Control Char] $relativePath (Línea $lineNum): Contiene carácter de control oculto" -ForegroundColor Magenta
        }

        if ($Fix) {
            $newLines.Add($line)
        }
    }

    # Guardar cambios si el modo reparación está activo y el archivo cambió
    if ($Fix -and $fileContentChanged) {
        $utf8NoBom = New-Object System.Text.UTF8Encoding($false)
        if ($newLines.Count -gt 0) {
            [System.IO.File]::WriteAllLines($file.FullName, $newLines, $utf8NoBom)
        } else {
            $content = [System.IO.File]::ReadAllText($file.FullName)
            [System.IO.File]::WriteAllText($file.FullName, $content, $utf8NoBom)
        }
        $fixedCount++
        Write-Host "   ✔️ [REPARADO] Archivo guardado correctamente en UTF-8 puro." -ForegroundColor Green
    }
}

$bomColor = if ($bomCount -gt 0) { "Yellow" } else { "Green" }
$corruptColor = if ($corruptCount -gt 0) { "Red" } else { "Green" }

Write-Host "----------------------------------------------------------"
Write-Host "📊 RESULTADOS DEL DIAGNÓSTICO:" -ForegroundColor Yellow
Write-Host "  ✔️ Archivos con marcas BOM: $bomCount" -ForegroundColor $bomColor
Write-Host "  ✔️ Archivos con Mojibake/Corrupción: $corruptCount" -ForegroundColor $corruptColor
if ($Fix) {
    Write-Host "  ✔️ Archivos saneados y escritos en disco: $fixedCount" -ForegroundColor Green
}
Write-Host "==========================================================" -ForegroundColor Cyan