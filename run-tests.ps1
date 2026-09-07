<#
.SYNOPSIS
    Script de automatización para la ejecución de pruebas contables en Go y generación de reportes de cobertura.
.DESCRIPTION
    Busca automáticamente el módulo de Go (en la carpeta actual o subcarpetas conocidas),
    ejecuta 'go test' en su contexto y compila/abre el reporte gráfico en formato HTML
    si se genera el perfil de cobertura.
#>

function Invoke-ContableTests {
    [CmdletBinding()]
    param(
        # Permite recibir cualquier argumento adicional de la consola (-race, -coverprofile, etc.)
        [Parameter(ValueFromRemainingArguments = $true)]
        [string[]]$TestArgs
    )

    # 1. Definir la raíz del script con respaldo para ejecución interactiva
    $scriptRoot = $PSScriptRoot
    if ([string]::IsNullOrWhiteSpace($scriptRoot)) {
        $scriptRoot = (Get-Location).Path
    }

    # 2. Buscar el directorio del proyecto Go automáticamente
    $possiblePaths = @(".", "go-skeleton", "go-contable", "contable-go", "GO-contable", "contable_go", "CONTABLE\go-skeleton")
    $projectDir = $null

    foreach ($path in $possiblePaths) {
        $fullPath = Join-Path -Path $scriptRoot -ChildPath $path
        if (Test-Path -Path (Join-Path -Path $fullPath -ChildPath "go.mod")) {
            $projectDir = $fullPath
            break
        }
    }

    # Si no se encuentra el archivo go.mod, se informa sin cerrar la terminal
    if (-not $projectDir) {
        Write-Host "❌ Error: No se encontró un módulo de Go válido con 'go.mod' en '$scriptRoot'." -ForegroundColor Red
        $global:LASTEXITCODE = 1
        return
    }

    Write-Host "🚀 Proyecto contable localizado en: $projectDir" -ForegroundColor Green
    
    # Cambiar temporalmente la ubicación a la carpeta del proyecto de Go
    Push-Location -Path $projectDir

    try {
        # 3. Sincronizar dependencias y limpiar archivos de cobertura previos
        Write-Host "`n🔄 Sincronizando dependencias con 'go mod tidy'..." -ForegroundColor Yellow
        go mod tidy
        if (Test-Path "coverage.out") { Remove-Item "coverage.out" -Force }
        if (Test-Path "coverage.html") { Remove-Item "coverage.html" -Force }

        # 4. Ejecución de la suite de pruebas en Go
        Write-Host "🏃 Ejecutando suites de pruebas contables..." -ForegroundColor Yellow
        
        # Mejora: Añadir automáticamente el perfil de cobertura si no se especifica uno.
        $finalTestArgs = $TestArgs
        if (-not ($TestArgs -like "*-coverprofile*")) {
            $finalTestArgs += "-coverprofile=coverage.out"
        }

        go test -v ./... @finalTestArgs
        
        # 5. Verificación del código de retorno de 'go test'
        if ($LASTEXITCODE -ne 0) {
            throw "La ejecución de 'go test' falló con código de salida: $LASTEXITCODE"
        }

        # 6. Generación y apertura automática del reporte HTML si existe archivo de cobertura
        if (Test-Path "coverage.out") {
            Write-Host "`n📊 Detectado 'coverage.out'. Compilando y abriendo reporte de cobertura HTML..." -ForegroundColor Cyan
            go tool cover -html=coverage.out -o coverage.html
            Start-Process "coverage.html" # Abre el reporte en el navegador predeterminado
        }

        Write-Host "`n✅ ¡Todas las pruebas contables han finalizado en verde!" -ForegroundColor Green
        $global:LASTEXITCODE = 0
    }
    catch {
        # Captura y despliegue del error en consola de manera segura
        Write-Host "`n❌ Ocurrió un error durante la ejecución de las pruebas:" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
        $global:LASTEXITCODE = 1
    }
    finally {
        # Regresar siempre al directorio original de la consola
        Pop-Location
    }
}

# Ejecución automática de la función pasando todos los argumentos recibidos por consola
Invoke-ContableTests @args