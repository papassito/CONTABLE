<#
.SYNOPSIS
    Script de automatización para pruebas en Go y generación de reportes de cobertura.
#>

function Invoke-ContableTests {
    [CmdletBinding()]
    param(
        [Parameter(ValueFromRemainingArguments = $true)]
        [string[]]$TestArgs
    )

    $scriptRoot = $PSScriptRoot
    if ([string]::IsNullOrWhiteSpace($scriptRoot)) { 
        $scriptRoot = (Get-Location).Path 
    }

    # 1. Búsqueda inteligente del módulo
    $possiblePaths = @(".", "go-skeleton", "go-contable", "contable-go", "GO-contable", "contable_go", "CONTABLE\go-skeleton")
    $projectDir = $null

    foreach ($path in $possiblePaths) {
        $fullPath = Join-Path -Path $scriptRoot -ChildPath $path
        if (Test-Path -Path (Join-Path -Path $fullPath -ChildPath "go.mod")) {
            $projectDir = $fullPath
            break
        }
    }

    if (-not $projectDir) {
        Write-Host "❌ Error: No se encontró 'go.mod' en '$scriptRoot' o sus subcarpetas." -ForegroundColor Red
        return 1
    }

    Write-Host "🚀 Proyecto localizado en: $projectDir" -ForegroundColor Green
    
    # Control explicito de la pila de directorios
    $pushedDir = $false
    $testExitCode = 1

    try {
        Push-Location -Path $projectDir
        $pushedDir = $true

        Write-Host "`n🔄 Sincronizando dependencias..." -ForegroundColor Yellow
        go mod tidy
        if ($LASTEXITCODE -ne 0) {
            Write-Host "⚠️ Advertencia: 'go mod tidy' finalizó con errores." -ForegroundColor Red
        }

        # 2. Resolución dinámica del archivo de cobertura
        $coverFile = "coverage.out"
        $finalTestArgs = @()
        if ($TestArgs) { $finalTestArgs += $TestArgs }

        # Búsqueda segura del argumento -coverprofile
        $customCoverArg = $finalTestArgs | Where-Object { $_ -like "-coverprofile=*" } | Select-Object -First 1

        if ($customCoverArg) {
            $coverFile = ($customCoverArg -split "=", 2)[1].Trim('"''')
        } else {
            $finalTestArgs += "-coverprofile=$coverFile"
        }

        # Limpiar cualquier reporte anterior
        if (Test-Path $coverFile) { Remove-Item $coverFile -Force -ErrorAction SilentlyContinue }
        if (Test-Path "coverage.html") { Remove-Item "coverage.html" -Force -ErrorAction SilentlyContinue }

        Write-Host "🏃 Ejecutando suites de pruebas..." -ForegroundColor Yellow
        
        # 3. Ejecución de pruebas
        go test -v $finalTestArgs ./... 
        $testExitCode = $LASTEXITCODE

        # 4. Generar reporte HTML SOLO si go test logró crear el archivo de cobertura
        if (Test-Path -Path $coverFile) {
            Write-Host "`n📊 Compilando reporte de cobertura HTML..." -ForegroundColor Cyan
            go tool cover "-html=$coverFile" -o coverage.html
            
            if ($LASTEXITCODE -eq 0 -and [Environment]::UserInteractive -and -not $env:CI) {
                Start-Process "coverage.html"
            }
        } else {
            Write-Host "`n⚠️ No se pudo generar el reporte porque no se encontró '$coverFile' (¿Las pruebas fallaron antes de ejecutarse?)." -ForegroundColor Yellow
        }

        # 5. Evaluación de resultado final
        if ($testExitCode -ne 0) {
            Write-Host "`n❌ Las pruebas finalizaron con errores." -ForegroundColor Red
        } else {
            Write-Host "`n✅ ¡Todas las pruebas finalizaron en verde!" -ForegroundColor Green
        }
    }
    catch {
        Write-Host "`n❌ Error crítico de ejecución: $($_.Exception.Message)" -ForegroundColor Red 
        $testExitCode = 1
    }
    finally {
        # Solo restaura la ubicación si Push-Location llegó a ejecutarse
        if ($pushedDir) {
            Pop-Location
        }
    }

    return $testExitCode
}

# Invocación directa y propagación del exit code
$exitCode = Invoke-ContableTests @args
if ($exitCode -ne 0) {
    exit $exitCode
}