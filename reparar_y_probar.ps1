<#
.SYNOPSIS
    Script contenedor para ejecutar la suite de pruebas y reparación.
.DESCRIPTION
    Este script ahora delega toda la lógica de reparación y ejecución de pruebas
    al script 'run-tests.ps1' para centralizar la lógica y evitar duplicación.
#>

[CmdletBinding()]
param(
    # Reenvía todos los argumentos adicionales al script de pruebas principal.
    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$TestArgs
)

$scriptRoot = $PSScriptRoot
if ([string]::IsNullOrWhiteSpace($scriptRoot)) { $scriptRoot = (Get-Location).Path }

$testRunnerScript = Join-Path -Path $scriptRoot -ChildPath "run-tests.ps1"
& $testRunnerScript @TestArgs