<#
.SYNOPSIS
    Script de reparación y ejecución de pruebas para el proyecto Contable Go.
.DESCRIPTION
    Detecta automáticamente la ubicación de 'go.mod' (incluso si se ejecuta desde
    la carpeta raíz de GitHub), repara contratos y ejecuta las pruebas de Go.
#>

[CmdletBinding()]
param(
    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$TestArgs
)

# 1. Determinar la ruta de origen
$basePath = $PSScriptRoot
if ([string]::IsNullOrWhiteSpace($basePath)) {
    $basePath = (Get-Location).Path
}

# 2. Definir las rutas posibles donde puede habitar 'go.mod'
$candidatePaths = @(
    $basePath,
    (Join-Path -Path $basePath -ChildPath "go-skeleton"),
    (Join-Path -Path $basePath -ChildPath "CONTABLE\go-skeleton"),
    (Join-Path -Path $basePath -ChildPath "CONTABLE")
)

$projectDir = $null

# Recorrer las rutas para encontrar la que contiene go.mod
foreach ($path in $candidatePaths) {
    $goModFile = Join-Path -Path $path -ChildPath "go.mod"
    if (Test-Path -Path $goModFile) {
        $projectDir = $path
        break
    }
}

# Si no se encuentra el archivo go.mod, se detiene la ejecución informando las rutas probadas
if (-not $projectDir) {
    Write-Host "❌ Error: No se encontró el archivo 'go.mod'. Se buscaron las siguientes rutas:" -ForegroundColor Red
    foreach ($path in $candidatePaths) {
        Write-Host "   - $path" -ForegroundColor Yellow
    }
    return
}

Write-Host "🚀 Módulo de Go localizado en: $projectDir" -ForegroundColor Green

# 3. Eliminar el archivo duplicado y vacío config\uow.go
$uowFile = Join-Path -Path $projectDir -ChildPath "config\uow.go"
if (Test-Path -Path $uowFile) {
    Remove-Item -Path $uowFile -Force
    Write-Host "  [+] Archivo vacío eliminado: config\uow.go" -ForegroundColor Green
}

# 4. Asegurar existencia de internal\domain y escribir document_parser.go
$domainDir = Join-Path -Path $projectDir -ChildPath "internal\domain"
if (-not (Test-Path -Path $domainDir)) {
    New-Item -ItemType Directory -Force -Path $domainDir | Out-Null
}

$parserFile = Join-Path -Path $domainDir -ChildPath "document_parser.go"
$parserCode = @'
package domain

import "io"

// DocumentParser define la interfaz para el análisis de documentos contables.
type DocumentParser interface {
	ParseCFDI40(xmlReader io.Reader) (*ParsedCFDIDTO, error)
	ParseSIPARELine(rawText string) (*ParsedSIPAREDTO, error)
}
'@

Set-Content -Path $parserFile -Value $parserCode -Encoding UTF8
Write-Host "  [+] Contrato restaurado en: internal\domain\document_parser.go" -ForegroundColor Green

# 5. Mover anchor.go a internal\domain\ si estaba en internal\service\
$wrongAnchor = Join-Path -Path $projectDir -ChildPath "internal\service\anchor.go"
$targetAnchor = Join-Path -Path $domainDir -ChildPath "anchor.go"

if (Test-Path -Path $wrongAnchor) {
    Move-Item -Path $wrongAnchor -Destination $targetAnchor -Force
    Write-Host "  [+] Reubicado 'anchor.go' a internal\domain\" -ForegroundColor Green
}

# 6. Cambiar a la carpeta del proyecto, sincronizar e invocar go test
Push-Location -Path $projectDir
try {
    Write-Host "`n🔄 Sincronizando dependencias con 'go mod tidy'..." -ForegroundColor Yellow
    go mod tidy

    Write-Host "🏃 Ejecutando suite de pruebas contables..." -ForegroundColor Yellow
    go test -v ./... @TestArgs
}
catch {
    Write-Host "`n❌ Ocurrió un error en la ejecución de Go:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}
finally {
    # Restaurar la ubicación original de la consola siempre
    Pop-Location
}