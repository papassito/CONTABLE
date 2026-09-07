<#
.SYNOPSIS
    Script de Tratamiento Médico para FCOS v2.2 Kernel.
.DESCRIPTION
    Pobla los archivos vacíos 'secret_envelope.go' y 'vault.go' en la capa de dominio,
    reparando los errores 'expected package, found EOF' en go vet.
#>

[CmdletBinding()]
param(
    [string]$TargetDir = (Get-Location).Path
)

# 1. Localización automática del módulo de Go
$candidateGoMod = Join-Path -Path $TargetDir -ChildPath "go-skeleton\go.mod"
if (Test-Path -Path $candidateGoMod) {
    $projectDir = Join-Path -Path $TargetDir -ChildPath "go-skeleton"
} elseif (Test-Path -Path (Join-Path -Path $TargetDir -ChildPath "go.mod")) {
    $projectDir = $TargetDir
} else {
    $projectDir = $TargetDir
}

Write-Host "==================================================================" -ForegroundColor Cyan
Write-Host " 💉 TRATAMIENTO DE REMEDIACIÓN MÉDICA - FCOS v2.2 KERNEL" -ForegroundColor Cyan
Write-Host "==================================================================" -ForegroundColor Cyan
Write-Host "📁 Proyecto objetivo: $projectDir`n" -ForegroundColor Gray

$domainDir = Join-Path -Path $projectDir -ChildPath "internal\domain"

# 2. Restauración de secret_envelope.go
$secretEnvelopePath = Join-Path -Path $domainDir -ChildPath "secret_envelope.go"
$secretEnvelopeCode = @'
package domain

import "time"

// SecretEnvelope representa la estructura de un contenedor cifrado de secretos.
type SecretEnvelope struct {
	ID           string    `json:"id"`
	KeyID        string    `json:"key_id"`
	Ciphertext   []byte    `json:"ciphertext"`
	CreatedAtUTC time.Time `json:"created_at_utc"`
}
'@

Set-Content -Path $secretEnvelopePath -Value $secretEnvelopeCode -Encoding UTF8
Write-Host "✅ Archivo 'secret_envelope.go' curado e inyectado exitosamente." -ForegroundColor Green

# 3. Restauración de vault.go
$vaultPath = Join-Path -Path $domainDir -ChildPath "vault.go"
$vaultCode = @'
package domain

import "context"

// SecretsVault define el contrato para el cifrado y manejo de secretos.
type SecretsVault interface {
	Encrypt(ctx context.Context, tenantID string, plaintext []byte) (*SecretEnvelope, error)
}
'@

Set-Content -Path $vaultPath -Value $vaultCode -Encoding UTF8
Write-Host "✅ Archivo 'vault.go' curado e inyectado exitosamente." -ForegroundColor Green

# 4. Sincronización y Alta Médica
Push-Location -Path $projectDir
try {
    Write-Host "`n🔄 Sincronizando módulos con 'go mod tidy'..." -ForegroundColor Yellow
    go mod tidy

    Write-Host "🫀 Tomando el pulso al código con 'go vet ./...'..." -ForegroundColor Yellow
    $vetOutput = go vet ./... 2>&1

    if ($LASTEXITCODE -eq 0) {
        Write-Host "`n🎉 ¡ALTA MÉDICA OTORGADA! 'GO VET' FINALIZÓ CON CERO ERRORES." -ForegroundColor Green
    } else {
        Write-Host "`n⚠️ Se detectaron observaciones restantes en 'go vet':" -ForegroundColor Yellow
        Write-Host $vetOutput -ForegroundColor Red
    }
}
finally {
    Pop-Location
}

Write-Host "`n==================================================================" -ForegroundColor Cyan