<#
.SYNOPSIS
    Script de compilación local para la aplicación Contable Fix (Wails).
.DESCRIPTION
    Automatiza la compilación del frontend y backend, generando un ejecutable
    y un instalador para Windows, similar al pipeline de release.
#>

function Invoke-ContableBuild {
    [CmdletBinding()]
    param()

    # --- CONFIGURACIÓN ---
    $ErrorActionPreference = "Stop"

    # --- FASE 1: VERIFICACIÓN DE PRERREQUISITOS ---
Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host " 🚀 INICIANDO COMPILACIÓN LOCAL - Contable Fix FCOS v2.2  " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

# Búsqueda inteligente del directorio del proyecto Go
$goProjectDir = $null
$scriptRoot = $PSScriptRoot
if ([string]::IsNullOrWhiteSpace($scriptRoot)) { $scriptRoot = (Get-Location).Path }

$possiblePaths = @(".", "go-skeleton", "go-contable", "contable-go", "GO-contable", "contable_go", "CONTABLE\go-skeleton")
foreach ($path in $possiblePaths) {
    $fullPath = Join-Path -Path $scriptRoot -ChildPath $path
    if (Test-Path -Path (Join-Path -Path $fullPath -ChildPath "go.mod")) {
        $goProjectDir = $fullPath
        break
    }
}

if (-not $goProjectDir) {
    Write-Host "❌ Error: No se pudo encontrar el directorio del proyecto Go (con go.mod)." -ForegroundColor Red
    return 1
}

Write-Host "`n📁 Directorio del proyecto Go localizado en: $goProjectDir" -ForegroundColor Green

Write-Host "`n🔍 1. Verificando prerrequisitos..." -ForegroundColor Yellow

# Verificar Go
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Error: Go no está instalado o no se encuentra en el PATH." -ForegroundColor Red
    return 1
}
Write-Host "   - Go detectado: $((go version).Split(' ')[2])" -ForegroundColor Green

# Verificar pnpm
if (-not (Get-Command pnpm -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Error: pnpm no está instalado. Ejecuta: npm install -g pnpm" -ForegroundColor Red
    return 1
}
Write-Host "   - pnpm detectado." -ForegroundColor Green

# Verificar Wails
if (-not (Get-Command wails -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Error: Wails CLI no está instalado. Ejecuta: go install github.com/wailsapp/wails/v2/cmd/wails@latest" -ForegroundColor Red
    return 1
}
Write-Host "   - Wails detectado: $((wails version).Trim())" -ForegroundColor Green

# --- FASE 2: INSTALACIÓN DE DEPENDENCIAS ---
Write-Host "`n📦 2. Instalando dependencias del frontend (pnpm)..." -ForegroundColor Yellow
pnpm install
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Error durante 'pnpm install'. Revisa los logs." -ForegroundColor Red
    return 1
}
Write-Host "📦 Asegurando tipos de React (@types/react, @types/react-dom)..." -ForegroundColor Yellow
pnpm add -D @types/react@19 @types/react-dom@19
pnpm build
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Error durante 'pnpm build' de Vite. Revisa los logs." -ForegroundColor Red
    return 1
}
Write-Host "   ✅ Frontend instalado y compilado correctamente en .\dist" -ForegroundColor Green

# --- FASE 2.5: ASEGURAR DEPENDENCIA DE WAILS ---
Write-Host "`n🔧 2.5. Sincronizando dependencias del módulo Go y resolviendo go.sum..." -ForegroundColor Yellow
$pushedDir = $false
try {
    Push-Location -Path $goProjectDir
    $pushedDir = $true
    
    Write-Host "   - Ejecutando 'go mod tidy' para resolver dependencias directas e indirectas..." -ForegroundColor Cyan
    go mod tidy
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Error durante 'go mod tidy'. Revisa los logs." -ForegroundColor Red
        return 1
    }
    Write-Host "   ✅ Módulo sincronizado y go.sum actualizado con éxito." -ForegroundColor Green
} finally {
    if ($pushedDir) { Pop-Location }
}

# --- FASE 3: COMPILACIÓN CON WAILS ---
Write-Host "`n🏗️  3. Compilando aplicación con Wails para windows/amd64..." -ForegroundColor Yellow

# Wails es estricto: espera que go.mod y wails.json estén en el mismo directorio.
# La solución más limpia es copiar wails.json al directorio de Go, ajustar su ruta de frontend,
# compilar desde allí y luego limpiar.
$originalWailsJsonPath = Join-Path -Path $scriptRoot -ChildPath "wails.json"
$tempWailsJsonPath = Join-Path -Path $goProjectDir -ChildPath "wails.json"
$targetDistPath = Join-Path -Path $goProjectDir -ChildPath "dist"
$sourceDistPath = Join-Path -Path $scriptRoot -ChildPath "dist"
$pushedDir = $false
$buildSucceeded = $false

try {
    # Copiar temporalmente la carpeta dist dentro de go-skeleton para evadir restricción de embed en Go
    if (Test-Path $targetDistPath) { Remove-Item -Path $targetDistPath -Recurse -Force }
    Copy-Item -Path $sourceDistPath -Destination $targetDistPath -Recurse -Force

    # Modificar y copiar wails.json
    $wailsConfig = (Get-Content -Path $originalWailsJsonPath -Raw | ConvertFrom-Json)
    $wailsConfig."frontend:dir" = ".." # Apuntar al directorio padre para el frontend
    
    # Ajustar rutas de bindings y assets para que apunten a la raíz desde el subdirectorio de Go
    if ($wailsConfig.wailsjsdir) { $wailsConfig.wailsjsdir = "../" + $wailsConfig.wailsjsdir }
    # Al copiar dist localmente, el assetdir para Wails pasa a ser relativo en go-skeleton
    $wailsConfig.assetdir = "dist" 
    $wailsConfig."build:dir" = "build" # Asegurar que el directorio de build sea relativo al directorio de Go
    $wailsConfig | ConvertTo-Json -Depth 10 | Set-Content -Path $tempWailsJsonPath

    # Navegar al directorio de Go y compilar
    Push-Location -Path $goProjectDir
    $pushedDir = $true

    # Usamos -skipfrontend para compilar la aplicación Go con los assets estáticos que pre-copiamos
    $buildArgs = @("-platform", "windows/amd64", "-clean", "-nsis", "-skipfrontend")
    wails build $buildArgs
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Error durante 'wails build'. Revisa los logs." -ForegroundColor Red
        return 1
    }
    $buildSucceeded = $true
    Write-Host "   ✅ Aplicación compilada exitosamente en directorio temporal." -ForegroundColor Green

    Write-Host "   ✅ Aplicación compilada exitosamente." -ForegroundColor Green
} finally {
    # Limpiar sin importar el resultado
    if ($pushedDir) { Pop-Location }
    if (Test-Path $tempWailsJsonPath) {
        Remove-Item -Path $tempWailsJsonPath -Force -ErrorAction SilentlyContinue
    }
    if (Test-Path $targetDistPath) {
        Remove-Item -Path $targetDistPath -Recurse -Force -ErrorAction SilentlyContinue
    }
        # Recrear la carpeta 'dist' con un archivo placeholder permanente para evitar errores de //go:embed en el IDE
        New-Item -ItemType Directory -Force -Path $targetDistPath | Out-Null
        New-Item -ItemType File -Force -Path (Join-Path $targetDistPath "placeholder.txt") -Value "Placeholder to satisfy go:embed" | Out-Null
}

# --- FASE 3.5: MOVER ARTEFACTOS DE COMPILACIÓN ---
if ($buildSucceeded) {
    Write-Host "`n🚚 3.5. Moviendo artefactos al directorio de compilación principal..." -ForegroundColor Yellow
    $sourceBuildDir = Join-Path -Path $goProjectDir -ChildPath "build\bin"
    $targetBuildDir = Join-Path -Path $scriptRoot -ChildPath "build\bin"
    
    if (-not (Test-Path $targetBuildDir)) { New-Item -Path $targetBuildDir -ItemType Directory -Force | Out-Null }
    Move-Item -Path "$sourceBuildDir\*" -Destination $targetBuildDir -Force
    Write-Host "   ✅ Artefactos movidos a '.\build\bin\'" -ForegroundColor Green
}

# --- FASE 4: FIRMA DIGITAL (OPCIONAL) ---
Write-Host "`n✍️  4. Verificando firma digital (opcional)..." -ForegroundColor Yellow
$executablePath = Join-Path -Path $scriptRoot -ChildPath "build\bin\ContableFixByKlik.exe"
if (-not (Test-Path $executablePath)) {
    Write-Host "   ⚠️ No se encontró el ejecutable en '$executablePath'. Omitiendo firma." -ForegroundColor Yellow
} elseif ($env:LOCAL_CERT_PATH -and $env:LOCAL_CERT_PASS) {
    Write-Host "   - Se encontraron variables de entorno para firma digital. Intentando firmar..." -ForegroundColor Cyan
    # Lógica de firma similar a release.yml
    $signtool = (Get-ChildItem -Path "C:\Program Files (x86)\Windows Kits" -Filter "signtool.exe" -Recurse -ErrorAction SilentlyContinue | Sort-Object LastWriteTime -Descending | Select-Object -First 1).FullName
    if (-not $signtool) { $signtool = "signtool.exe" } # Fallback al PATH
    & $signtool sign /f $env:LOCAL_CERT_PATH /p $env:LOCAL_CERT_PASS /tr http://timestamp.digicert.com /td sha256 /fd sha256 $executablePath
    Write-Host "   ✅ Ejecutable firmado digitalmente con éxito." -ForegroundColor Green
} else {
    Write-Host "   - El ejecutable se generó sin firma digital (variables 'LOCAL_CERT_PATH' y 'LOCAL_CERT_PASS' no configuradas)." -ForegroundColor Gray
}

# --- FINALIZACIÓN ---
Write-Host "`n==========================================================" -ForegroundColor Cyan
Write-Host "🎉 ¡COMPILACIÓN FINALIZADA!" -ForegroundColor Green
Write-Host "   - Ejecutable: $($executablePath.Replace($scriptRoot, '.\'))"
Write-Host "   - Instalador: $($executablePath.Replace($scriptRoot, '.\').Replace('.exe', '-installer.exe'))"
Write-Host "==========================================================" -ForegroundColor Cyan
    return 0
}

# Invocación controlada del script
$exitCode = Invoke-ContableBuild
if ($exitCode -ne 0) {
    # En entornos de CI (como GitHub Actions), forzamos la salida para romper el build
    if ($env:CI) {
        exit $exitCode
    } else {
        return $exitCode
    }
}