<#
.SYNOPSIS
    Scanner Agresivo de Diagnóstico (Solo Lectura)
.DESCRIPTION
    Inspecciona procesos, red, persistencia y estado de seguridad sin modificar el sistema.
#>

$ErrorActionPreference = "SilentlyContinue"
Clear-Host

Write-Host "===============================================================================" -ForegroundColor Red
Write-Host " 🛡️  INICIANDO SCANNER AGRESIVO DE DIAGNÓSTICO (MODO: SOLO LECTURA)  🛡️ " -ForegroundColor Red
Write-Host "===============================================================================" -ForegroundColor Red
Write-Host "Ejecutando escrutinio profundo del sistema... Ningún archivo será modificado.`n" -ForegroundColor Yellow

# -------------------------------------------------------------------------
# 1. ESTADO DEL ANTIVIRUS Y SEGURIDAD BASE
# -------------------------------------------------------------------------
Write-Host "[1/6] VERIFICANDO ESTADO DE WINDOWS DEFENDER / SEGURIDAD..." -ForegroundColor Cyan
$defender = Get-MpComputerStatus
if ($defender) {
    if ($defender.RealTimeProtectionEnabled) {
        Write-Host "  [+] Protección en Tiempo Real: ACTIVA" -ForegroundColor Green
    } else {
        Write-Host "  [!] Protección en Tiempo Real: DESACTIVADA (Peligro)" -ForegroundColor Red
    }
    Write-Host "  [*] Última actualización de firmas: $($defender.AntivirusSignatureLastUpdated)" -ForegroundColor DarkGray
} else {
    Write-Host "  [!] No se pudo obtener el estado de Defender. ¿Antivirus de terceros bloqueando?" -ForegroundColor Red
}
Write-Host ""

# -------------------------------------------------------------------------
# 2. PROCESOS SOSPECHOSOS (Ejecutándose desde Temp o AppData)
# -------------------------------------------------------------------------
Write-Host "[2/6] CAZANDO PROCESOS SOSPECHOSOS (Paths Inusuales)..." -ForegroundColor Cyan
$suspiciousProcesses = Get-Process | Where-Object { 
    $_.Path -and ($_.Path -match "\\AppData\\" -or $_.Path -match "\\Temp\\")
}

if ($suspiciousProcesses) {
    Write-Host "  [!] Se detectaron procesos ejecutándose desde carpetas temporales o de usuario:" -ForegroundColor Red
    $suspiciousProcesses | Select-Object Id, ProcessName, Path | Format-Table -AutoSize
} else {
    Write-Host "  [+] No se detectaron procesos activos en rutas sospechosas." -ForegroundColor Green
}
Write-Host ""

# -------------------------------------------------------------------------
# 3. CONEXIONES DE RED ACTIVAS (Quién está hablando hacia afuera)
# -------------------------------------------------------------------------
Write-Host "[3/6] ANALIZANDO CONEXIONES DE RED ACTIVAS (ESTABLISHED)..." -ForegroundColor Cyan
$connections = Get-NetTCPConnection | Where-Object { $_.State -eq 'Established' -and $_.RemoteAddress -notmatch "^127\.|^192\.168\.|^10\.|^172\.(1[6-9]|2[0-9]|3[0-1])\.|^::1" }

if ($connections) {
    Write-Host "  [*] Conexiones externas activas encontradas:" -ForegroundColor Yellow
    $netOutput = @()
    foreach ($conn in $connections) {
        $proc = Get-Process -Id $conn.OwningProcess -ErrorAction SilentlyContinue
        $netOutput += [PSCustomObject]@{
            PID = $conn.OwningProcess
            Proceso = if ($proc) { $proc.ProcessName } else { "Desconocido" }
            IPLocal = $conn.LocalAddress
            IPRemota = $conn.RemoteAddress
            PuertoRemoto = $conn.RemotePort
        }
    }
    $netOutput | Format-Table -AutoSize
} else {
    Write-Host "  [+] No hay conexiones externas activas o inusuales en este momento." -ForegroundColor Green
}
Write-Host ""

# -------------------------------------------------------------------------
# 4. MECANISMOS DE PERSISTENCIA (Arranque Automático)
# -------------------------------------------------------------------------
Write-Host "[4/6] REVISANDO REGISTROS DE ARRANQUE (PERSISTENCIA)..." -ForegroundColor Cyan
$runPaths = @(
    "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Run",
    "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Run"
)

$foundSuspiciousRun = $false
foreach ($path in $runPaths) {
    if (Test-Path $path) {
        $items = Get-ItemProperty -Path $path -ErrorAction SilentlyContinue
        if ($items) {
            foreach ($item in @($items)) {
                foreach ($prop in $item.psobject.properties) {
                    if ($prop.Name -notmatch "PSPath|PSParentPath|PSChildName|PSDrive|PSProvider") {
                        $val = $prop.Value
                        if ($val -match "\\Temp\\" -or $val -match "\\AppData\\") {
                            Write-Host "  [!] Alerta de Arranque Sospechoso: $($prop.Name) -> $val" -ForegroundColor Red
                            $foundSuspiciousRun = $true
                        }
                    }
                }
            }
        }
    }
}
if (-not $foundSuspiciousRun) {
    Write-Host "  [+] Registro de arranque limpio (sin ejecutables en Temp/AppData)." -ForegroundColor Green
}
Write-Host ""

# -------------------------------------------------------------------------
# 5. INTEGRIDAD DEL ARCHIVO HOSTS (Redirecciones maliciosas)
# -------------------------------------------------------------------------
Write-Host "[5/6] VERIFICANDO INTEGRIDAD DEL ARCHIVO HOSTS..." -ForegroundColor Cyan
$hostsPath = "$env:windir\System32\drivers\etc\hosts"
if (Test-Path $hostsPath) {
    $hostsContent = Get-Content $hostsPath | Where-Object { 
        $_ -notmatch "^\s*#" -and 
        $_ -match "\w" -and 
        $_ -notmatch "^\s*(127\.0\.0\.1|::1|fe80::1%lo0)\s+localhost\s*$"
    }
    if ($hostsContent.Count -gt 0) {
        Write-Host "  [!] El archivo HOSTS tiene redirecciones activas:" -ForegroundColor Yellow
        $hostsContent | ForEach-Object { Write-Host "      $_" -ForegroundColor DarkGray }
    } else {
        Write-Host "  [+] Archivo HOSTS limpio (solo contiene comentarios)." -ForegroundColor Green
    }
}
Write-Host ""

# -------------------------------------------------------------------------
# 6. DETECCIÓN DE PROCESOS HÚNGAROS (Alto Consumo / Mineros / Bloqueadores)
# -------------------------------------------------------------------------
Write-Host "[6/6] DETECTANDO PROCESOS DE ALTO CONSUMO (CPU/RAM)..." -ForegroundColor Cyan
$highHogs = Get-Process | Where-Object { $_.CPU -ne $null } | Sort-Object CPU -Descending | Select-Object -First 5
Write-Host "  [*] Top 5 Procesos con mayor tiempo acumulado de CPU (Segundos):" -ForegroundColor Yellow
$highHogs | Select-Object Id, ProcessName, @{Name="CPU Acumulado (s)";Expression={[math]::Round($_.CPU, 2)}}, Path | Format-Table -AutoSize

Write-Host "===============================================================================" -ForegroundColor Red
Write-Host " SCAN DE DIAGNÓSTICO FINALIZADO " -ForegroundColor Green
Write-Host "===============================================================================" -ForegroundColor Red