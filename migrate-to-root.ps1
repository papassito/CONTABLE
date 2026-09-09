<#
.SYNOPSIS
    Script de unificación de espacio de trabajo y generación de documentación.
.DESCRIPTION
    1. Crea carpetas \docs y \docs\phases.
    2. Genera requirements.md y maps.md en UTF-8 puro (sin BOM).
    3. Mueve de forma segura el módulo de Go de 'Contable GO' (o go-skeleton) a la raíz del repositorio.
    4. Limpia carpetas residuales vacías sin romper el historial ni el árbol de trabajo.
#>

[CmdletBinding()]
param(
    [string]$Path = $PSScriptRoot
)

if ([string]::IsNullOrEmpty($Path)) { $Path = (Get-Location).Path }

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host " 🚀 INICIANDO MIGRACIÓN A LA RAÍZ Y GENERACIÓN DE DOCS     " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

# --- FASE 1: CREACIÓN DE CARPETAS Y DOCUMENTACIÓN ---
Write-Host "`n📁 1. Creando estructura de documentación..." -ForegroundColor Yellow
$docsPath = Join-Path -Path $Path -ChildPath "docs"
$phasesPath = Join-Path -Path $docsPath -ChildPath "phases"

foreach ($dir in @($docsPath, $phasesPath)) {
    if (-not (Test-Path -Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force | Out-Null
        Write-Host "   - Directorio creado: $($dir.Replace($Path, ''))" -ForegroundColor Green
    }
}

# Helper para escribir UTF-8 sin BOM
function Write-TextNoBom {
    param ([string]$FilePath, [string]$Content)
    $utf8NoBom = New-Object System.Text.UTF8Encoding($false)
    [System.IO.File]::WriteAllText($FilePath, $Content, $utf8NoBom)
    Write-Host "   - Archivo generado: $($FilePath.Replace($Path, ''))" -ForegroundColor DarkGreen
}

# requirements.md
$reqs = @"
# Especificación de Requisitos Contables (FCOS v2.2)

## Invariantes del Motor Financiero
- **Partida Doble**: Todo asiento contable debe cumplir de forma estricta con la condición: `∑ Débitos = ∑ Créditos`.
- **Inmutabilidad**: Prohibida la edición o borrado físico de asientos contables en estado 'POSTED'.
- **Aritmética Exacta**: Todos los cálculos monetarios se representan internamente en centavos enteros (`int64`).
"@
Write-TextNoBom -FilePath (Join-Path -Path $docsPath -ChildPath "requirements.md") -Content $reqs

# maps.md
$maps = @"
# Mapa de la Solución (Arquitectura Unificada)

El backend en Go y el frontend en React ahora conviven directamente en la raíz del repositorio, simplificando la compilación nativa de Wails.

```text
CONTABLE/ (Raíz unificada)
├── go.mod / wails.json  <-- Conviven en el mismo nivel
├── internal/            <-- Dominio, servicios y adaptadores del FCOS
├── src/                 <-- Código fuente React de la interfaz de usuario
└── docs/                <-- Requisitos, mapas y fases
```
"@
Write-TextNoBom -FilePath (Join-Path -Path $docsPath -ChildPath "maps.md") -Content $maps

# --- FASE 2: MIGRACIÓN SEGURA DEL BACKEND A LA RAÍZ ---
Write-Host "`n📦 2. Buscando carpeta contenedora de Go..." -ForegroundColor Yellow
$sourceDirName = "Contable GO"
$sourcePath = Join-Path -Path $Path -ChildPath $sourceDirName

if (-not (Test-Path -Path $sourcePath)) {
    $sourceDirName = "go-skeleton"
    $sourcePath = Join-Path -Path $Path -ChildPath $sourceDirName
}

if (-not (Test-Path -Path $sourcePath)) {
    if (Test-Path -Path (Join-Path -Path $Path -ChildPath "go.mod")) {
        Write-Host "ℹ️ El proyecto Go ya reside directamente en la raíz. No se requiere migración." -ForegroundColor Yellow
        Write-Host "==========================================================" -ForegroundColor Cyan
        return 0
    }
    Write-Host "❌ Error: No se encontró la carpeta 'Contable GO' ni 'go-skeleton'." -ForegroundColor Red
    return 1
}

Write-Host "   - Carpeta origen detectada: '$sourceDirName'" -ForegroundColor Cyan
Write-Host "⚠️  Se moverán todos los subdirectorios y archivos a la raíz de forma segura." -ForegroundColor Yellow
Read-Host "Presione ENTER para proceder con la unificación o CTRL+C para cancelar"

try {
    # Mover archivos uno por uno para no colisionar con archivos de la raíz ya existentes (.gitignore, etc)
    Get-ChildItem -Path $sourcePath | ForEach-Object {
        $destPath = Join-Path -Path $Path -ChildPath $_.Name
        if (Test-Path -Path $destPath) {
            Write-Host "   - El archivo/directorio '$($_.Name)' ya existe en la raíz, se omite o integra." -ForegroundColor Gray
        } else {
            Move-Item -Path $_.FullName -Destination $destPath -Force
            Write-Host "   [MOVIDO] $($_.Name) -> Raíz" -ForegroundColor Green
        }
    }
    
    # Remover carpeta temporal ahora vacía
    Remove-Item -Path $sourcePath -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host "`n✅ Carpeta contenedora anterior '$sourceDirName' purgada con éxito." -ForegroundColor Green
} catch {
    Write-Host "❌ Error durante la migración física: $_" -ForegroundColor Red
    return 1
}

Write-Host "`n🎉 ¡PROCESO DE UNIFICACIÓN FINALIZADO CON ÉXITO!" -ForegroundColor Green
Write-Host "==========================================================" -ForegroundColor Cyan