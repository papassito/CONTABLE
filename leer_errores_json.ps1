<#
.SYNOPSIS
    Script para lectura y análisis de errores en archivos JSON.
.DESCRIPTION
    Lee un archivo en formato JSON, convierte su estructura a objetos de PowerShell
    y filtra cualquier elemento que reporte un estado o nivel de error.
#>

[CmdletBinding()]
param(
    [Parameter(Mandatory = $false)]
    [string]$JsonPath = "reporte.json"
)

Write-Host "==================================================================" -ForegroundColor Cyan
Write-Host " 🔍 LECTOR Y DIAGNÓSTICO DE ERRORES JSON                          " -ForegroundColor Cyan
Write-Host "==================================================================" -ForegroundColor Cyan

# 1. Verificar si el archivo existe
if (-not (Test-Path -Path $JsonPath)) {
    Write-Host "❌ Error: No se encontró el archivo JSON en la ruta: '$JsonPath'" -ForegroundColor Red
    $global:LASTEXITCODE = 1
    return
}

Write-Host "📁 Leyendo archivo: $JsonPath`n" -ForegroundColor Gray

# 2. Leer y parsear el contenido JSON
try {
    $rawContent = Get-Content -Path $JsonPath -Raw -ErrorAction Stop
    $jsonData = ConvertFrom-Json -InputObject $rawContent -ErrorAction Stop
}
catch {
    Write-Host "❌ Error al procesar el archivo JSON: La sintaxis es inválida." -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    $global:LASTEXITCODE = 1
    return
}

# 3. Identificar y filtrar errores
$detectedErrors = [System.Collections.Generic.List[PSCustomObject]]::new()

# Si el JSON es una lista de elementos (Array)
if ($jsonData -is [Array]) {
    foreach ($item in $jsonData) {
        # Evaluar campos comunes de error
        if ($item.status -eq "error" -or 
            $item.level -eq "error" -or 
            $item.type -eq "error" -or 
            [string]::IsNullOrWhiteSpace($item.error) -eq $false) {
            
            $detectedErrors.Add($item)
        }
    }
}
# Si el JSON es un objeto único
else {
    # Comprobar si el objeto tiene una propiedad de errores o un estado de fallo
    if ($jsonData.status -eq "error" -or $jsonData.errors -or $jsonData.error) {
        $detectedErrors.Add($jsonData)
    }
}

# 4. Mostrar resultados
Write-Host "📊 RESULTADOS DEL ANÁLISIS:" -ForegroundColor Yellow
Write-Host "------------------------------------------------------------------" -ForegroundColor DarkGray

if ($detectedErrors.Count -gt 0) {
    Write-Host "❌ Se encontraron $($detectedErrors.Count) registro(s) de error en el archivo:`n" -ForegroundColor Red
    
    $index = 1
    foreach ($err in $detectedErrors) {
        Write-Host "  [$index] Detalle del Error:" -ForegroundColor Yellow
        
        # Imprimir todas las propiedades del objeto de error
        foreach ($prop in $err.PSObject.Properties) {
            Write-Host "      • $($prop.Name) : $($prop.Value)" -ForegroundColor White
        }
        Write-Host ""
        $index++
    }
    $global:LASTEXITCODE = 1
} else {
    Write-Host "✅ ¡Excelente! No se detectaron registros de error en el archivo JSON." -ForegroundColor Green
    $global:LASTEXITCODE = 0
}

Write-Host "==================================================================" -ForegroundColor Cyan