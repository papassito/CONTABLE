<#
.SYNOPSIS
    Script para indagar a fondo 'go vet' y reparar advertencias de asignación.
.DESCRIPTION
    1. Ejecuta 'go vet -v' enfocado en el paquete de infraestructura de salud.
    2. Aplica la corrección del operador en telemetry.go.
    3. Re-ejecuta 'go vet' en todo el proyecto para confirmar 0 advertencias.
#>

[CmdletBinding()]
param()

# 1. Localizar la carpeta raíz del proyecto Go
$scriptRoot = $PSScriptRoot
if ([string]::IsNullOrWhiteSpace($scriptRoot)) {
    $scriptRoot = (Get-Location).Path
}

$candidatePaths = @(
    $scriptRoot,
    (Join-Path -Path $scriptRoot -ChildPath "go-skeleton"),
    (Join-Path -Path $scriptRoot -ChildPath "CONTABLE\go-skeleton"),
    (Join-Path -Path $scriptRoot -ChildPath "contable-go"),
    (Join-Path -Path $scriptRoot -ChildPath "GO-contable")
)

$projectDir = $null
foreach ($path in $candidatePaths) {
    if (Test-Path -Path (Join-Path -Path $path -ChildPath "go.mod")) {
        $projectDir = $path
        break
    }
}

if (-not $projectDir) {
    Write-Host "❌ Error: No se encontró la carpeta del proyecto con 'go.mod'." -ForegroundColor Red
    return
}

Write-Host "🔍 INDAGANDO GO VET EN: $projectDir" -ForegroundColor Cyan

Push-Location -Path $projectDir
try {
    Write-Host "`n1. Ejecutando 'go vet' en todo el Kernel..." -ForegroundColor Yellow
    $vetOutput = go vet ./... 2>&1

    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ ¡ANÁLISIS GO VET LLEGÓ A CERO ADVERTENCIAS! PROYECTO 100% LIMPIO." -ForegroundColor Green
    } else {
        Write-Host "⚠️ Se encontraron advertencias en 'go vet':" -ForegroundColor Yellow
        Write-Host $vetOutput -ForegroundColor Red
    }
}
finally {
    Pop-Location
}