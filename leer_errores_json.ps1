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
$resolvedPath = $JsonPath
if (-not (Test-Path -Path $resolvedPath) -and $PSScriptRoot) {
    $scriptRelativePath = Join-Path -Path $PSScriptRoot -ChildPath $JsonPath
    if (Test-Path -Path $scriptRelativePath) {
        $resolvedPath = $scriptRelativePath
    }
}

if (-not (Test-Path -Path $resolvedPath)) {
    Write-Host "❌ Error: No se encontró el archivo JSON en la ruta: '$JsonPath'" -ForegroundColor Red
    $global:LASTEXITCODE = 1
    return
}

Write-Host "📁 Leyendo archivo: $resolvedPath`n" -ForegroundColor Gray

# 2. Leer y parsear el contenido JSON
try {
    $rawContent = Get-Content -Path $resolvedPath -Raw -Encoding UTF8 -ErrorAction Stop
    $jsonData = ConvertFrom-Json -InputObject $rawContent -ErrorAction Stop
}
catch {
    Write-Host "❌ Error al procesar el archivo JSON: La sintaxis es inválida." -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    $global:LASTEXITCODE = 1
    return
}

$detectedErrors = [System.Collections.Generic.List[PSCustomObject]]::new()

# 3. Analizador Recursivo Profundo para Estructuras Complejas y de Telemetría
function Find-JsonElement ($element, $pathContext = "$") {
    if ($element -is [Array]) {
        for ($i = 0; $i -lt $element.Count; $i++) {
            Find-JsonElement $element[$i] ($pathContext + "[" + $i + "]")
        }
    } elseif ($element -is [PSCustomObject] -or $element -is [System.Management.Automation.PSCustomObject]) {
        # Evaluar si este objeto en sí mismo es un nodo de error
        $hasErrorPattern = $false
        $errorProps = @{}
        
        foreach ($prop in $element.PSObject.Properties) {
            $name = $prop.Name.ToLower()
            $valStr = [string]$prop.Value
            
            # Validar coincidencias de patrones de error únicamente en nodos terminales (primitivos)
            $isPrimitive = $prop.Value -isnot [PSCustomObject] -and $prop.Value -isnot [Array] -and $prop.Value -isnot [System.Collections.IList]
            
            if ($isPrimitive -and (
                ($name -eq "status" -and $valStr -eq "error") -or
                ($name -eq "level" -and ($valStr -eq "error" -or $valStr -eq "fail")) -or
                ($name -match "error|exception|fail" -and -not [string]::IsNullOrWhiteSpace($valStr))
            )) {
                $hasErrorPattern = $true
            }
            $errorProps[$prop.Name] = $prop.Value
        }
        
        if ($hasErrorPattern) {
            $errorProps["JSONPath"] = $pathContext
            $detectedErrors.Add([PSCustomObject]$errorProps)
        }
        
        # Continuar la búsqueda recursivamente por propiedades hijas
        foreach ($prop in $element.PSObject.Properties) {
            Find-JsonElement $prop.Value "$pathContext.$($prop.Name)"
        }
    }
}

Find-JsonElement $jsonData

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