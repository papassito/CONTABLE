# --- CONFIGURACIÓN DE RUTAS ---
$SourceDir = "C:\Users\Radio 2027\Documents\Codex\2026-09-09\p\outputs\CONTABLE-MD-FINALES"
$TargetDir = "C:\Users\Radio 2027\Desktop\Soft Septioembre 2026\CONTABLE"
$BackupDir = "C:\Users\Radio 2027\Desktop\Soft Septioembre 2026\CONTABLE_DOCUMENTOS_BACKUP"

$Files = @(
    "README.md",
    "ARCHITECTURE.md",
    "REQUIREMENTS.md",
    "MAP.md",
    "SECURITY.md",
    "AUDIT_GUIDE.md"
)

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host " 🛡️  PROCESO DE RECONCILIACIÓN DOCUMENTAL DE CONTABLE      " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

# 1. Crear Carpeta de Respaldo Externa si no existe
if (-not (Test-Path $BackupDir)) {
    New-Item -ItemType Directory -Path $BackupDir -Force | Out-Null
    Write-Host "📁 Carpeta de respaldo creada en: $BackupDir" -ForegroundColor Green
}

# 2. Respaldar Documentos Actuales
Write-Host "`n💾 1. Respaldando documentos actuales..." -ForegroundColor Yellow
foreach ($File in $Files) {
    $CurrentPath = Join-Path -Path $TargetDir -ChildPath $File
    if (Test-Path $CurrentPath) {
        $BackupPath = Join-Path -Path $BackupDir -ChildPath $File
        Copy-Item -Path $CurrentPath -Destination $BackupPath -Force
        Write-Host "   - Respaldado: $File" -ForegroundColor DarkGray
    } else {
        Write-Host "   - [!] No se encontró versión anterior de $File en destino." -ForegroundColor Yellow
    }
}

# 3. Copiar y Verificar Hashes SHA-256
Write-Host "`n🚚 2. Copiando archivos y validando integridad criptográfica..." -ForegroundColor Yellow
$AllHashesMatch = $true

foreach ($File in $Files) {
    $SrcFile = Join-Path -Path $SourceDir -ChildPath $File
    $DstFile = Join-Path -Path $TargetDir -ChildPath $File

    if (-not (Test-Path $SrcFile)) {
        Write-Host "❌ Error: No se localiza el archivo de origen '$SrcFile'" -ForegroundColor Red
        $AllHashesMatch = $false
        continue
    }

    # Copia física
    Copy-Item -Path $SrcFile -Destination $DstFile -Force
    Write-Host "   - Copiado: $File -> Raíz del Proyecto" -ForegroundColor Green

    # Obtención de hashes para verificación
    $SrcHash = (Get-FileHash -Path $SrcFile -Algorithm SHA256).Hash
    $DstHash = (Get-FileHash -Path $DstFile -Algorithm SHA256).Hash

    if ($SrcHash -eq $DstHash) {
        Write-Host "     ✅ Hash SHA-256 Coincide: $SrcHash" -ForegroundColor DarkGreen
    } else {
        Write-Host "     ❌ ERROR de Integridad en $File (Hashes no coinciden)" -ForegroundColor Red
        Write-Host "        Origen: $SrcHash" -ForegroundColor Red
        Write-Host "        Destino: $DstHash" -ForegroundColor Red
        $AllHashesMatch = $false
    }
}

Write-Host "`n==========================================================" -ForegroundColor Cyan
if ($AllHashesMatch) {
    Write-Host "🎉 PROCESO COMPLETADO: Todos los archivos se han copiado y verificado con éxito." -ForegroundColor Green
} else {
    Write-Host "⚠️ El proceso finalizó con discrepancias o errores de copia." -ForegroundColor Yellow
}
Write-Host "==========================================================" -ForegroundColor Cyan
