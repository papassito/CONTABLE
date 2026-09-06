<#
.SYNOPSIS
    Script de automatización para la ejecución de pruebas contables en Go y generación de reportes de cobertura.
.DESCRIPTION
    Busca automáticamente el módulo de Go (go-skeleton o go-contable), ejecuta 'go test' en su contexto
    y compila/abre el reporte gráfico en formato HTML si se genera el perfil de cobertura.
#>

function Invoke-ContableTests {
    [CmdletBinding()]
    param(
        [Parameter(ValueFromRemainingArguments = $true)]
        [string[]]$TestArgs
    )

    # 1. Definir la raíz del script para búsquedas relativas
    $scriptRoot = $PSScriptRoot

    # 2. Buscar el directorio del proyecto Go automáticamente
    $possiblePaths = @("go-skeleton", "go-contable")
    $projectDir = $null

    foreach ($path in $possiblePaths) {
        $fullPath = Join-Path $scriptRoot $path
        if (Test-Path (Join-Path $fullPath "go.mod")) {
            $projectDir = $fullPath
            break
        }
    }

    # Si no se encuentra el archivo go.mod, se detiene la ejecución
    if (-not $projectDir) {
        Write-Host "❌ Error: No se encontró un módulo de Go válido en 'go-skeleton' ni 'go-contable' dentro de '$scriptRoot'." -ForegroundColor Red
        exit 1
    }

    Write-Host "🚀 Proyecto contable localizado en: $projectDir" -ForegroundColor Green
    
    # Cambiar temporalmente la ubicación a la carpeta del proyecto de Go
    Push-Location $projectDir

    try {
        # 3. Limpieza de archivos de cobertura previos
        if (Test-Path "coverage.out") { Remove-Item "coverage.out" -Force }
        if (Test-Path "coverage.html") { Remove-Item "coverage.html" -Force }

        # 4. Ejecución de la suite de pruebas en Go
        Write-Host "🏃 Ejecutando suites de pruebas contables..." -ForegroundColor Yellow
        
        go test -v ./... @TestArgs

        # 5. Verificación crítica del código de retorno ($LASTEXITCODE)
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
    }
    catch {
        # Captura y despliegue del error en consola con código de salida no cero
        Write-Host "`n❌ Ocurrió un error durante la ejecución de las pruebas:" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
        exit 1
    }
    finally {
        # Regresar siempre al directorio original de la consola, sin importar el resultado
        Pop-Location
    }
}

# Ejecución automática de la función pasando todos los argumentos recibidos por consola
Invoke-ContableTests @args