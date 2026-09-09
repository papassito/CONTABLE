<#
.SYNOPSIS
    Generador de Documentación Técnica de Requisitos, Mapas, Seguridad y Contratos.
.DESCRIPTION
    Crea de manera automatizada las carpetas de documentación contable (\docs y \docs\phases)
    escribiendo los archivos en formato UTF-8 limpio (Sin BOM).
#>

[CmdletBinding()]
param(
    [string]$TargetDir = $PSScriptRoot
)

if ([string]::IsNullOrEmpty($TargetDir)) {
    $TargetDir = (Get-Location).Path
}

Write-Host "==================================================================" -ForegroundColor Cyan
Write-Host " 📚 GENERADOR DE MARCO DE DOCUMENTACIÓN - FCOS KERNEL v2.2       " -ForegroundColor Cyan
Write-Host "==================================================================" -ForegroundColor Cyan

# Directorios a crear
$docsPath = Join-Path -Path $TargetDir -ChildPath "docs"
$phasesPath = Join-Path -Path $docsPath -ChildPath "phases"

foreach ($path in @($docsPath, $phasesPath)) {
    if (!(Test-Path -Path $path)) {
        New-Item -ItemType Directory -Path $path -Force | Out-Null
        Write-Host "📁 Directorio creado: $($path.Replace($TargetDir, ''))" -ForegroundColor Green
    }
}

# Helper para escribir archivos en UTF-8 sin BOM (Compatible con las reglas del repositorio)
function Write-Utf8NoBom {
    param (
        [string]$FilePath,
        [string]$Content
    )
    try {
        $utf8NoBom = New-Object System.Text.UTF8Encoding($false)
        [System.IO.File]::WriteAllText($FilePath, $Content, $utf8NoBom)
        Write-Host "📝 Archivo generado: $($FilePath.Replace($TargetDir, ''))" -ForegroundColor DarkGreen
    } catch {
        Write-Error "Error escribiendo en $FilePath: $_"
    }
}

# ==========================================
# 1. REQUISITOS (requirements.md)
# ==========================================
$requirementsContent = @"
# Especificación de Requisitos Contables (FCOS v2.2)

Este documento define los requisitos funcionales mínimos e invariantes matemáticas obligatorias del núcleo contable.

## 1. Invariantes del Motor Financiero
- **Partida Doble Estricta**: No se puede registrar un asiento contable (`JournalEntry`) si la diferencia entre la suma total de débitos y créditos es distinta de cero (`∑ Débito - ∑ Crédito != 0`). El motor debe retornar `ErrUnbalancedJournal`.
- **Inmutabilidad Absoluta**: Queda estrictamente prohibida la edición (`UPDATE`) o eliminación (`DELETE`) física de asientos en estado `POSTED` (CONTABILIZADO) o de sus líneas asociadas.
- **Afectación de Cuentas Auxiliares**: Solo las cuentas de último nivel jerárquico que tengan marcado el indicador `accepts_move = true` pueden recibir imputaciones contables directas.
- **Saldos Calculados**: El campo `current_balance` es un campo puramente acumulativo y derivado. Queda prohibida su alteración por fuera del procesamiento estructurado de un asiento contable.

## 2. Precisión Decimal
- Para evitar pérdidas de precisión por redondeo del estándar IEEE 754, todos los cálculos de saldos, importes, débitos y créditos del backend de Go deben representarse mediante enteros en centavos (`int64`, donde `$10.00` = `1000` céntimos) o mediante aritmética de precisión fija utilizando la librería `github.com/shopspring/decimal`.
"@

Write-Utf8NoBom -FilePath (Join-Path -Path $docsPath -ChildPath "requirements.md") -Content $requirementsContent

# ==========================================
# 2. MAPA DE ARQUITECTURA (architecture_map.md)
# ==========================================
$mapContent = @"
# Mapa y Flujo de Arquitectura Contable

## Flujo de Datos Transaccional (De UI a Persistencia)

```text
     [ Capa de Presentación (React / TS) ]
                       │
                       │ (Tipos e importes decimales legibles)
                       ▼
        [ Wails Service Bridge (TS) ]
                       │
                       │ (Conversión: Decimales -> Céntimos Enteros)
                       ▼
     [ Capa de Transporte API (Go / HTTP / JSON) ]
                       │
                       │ (Decodifica DTOs e inyecta dependencias)
                       ▼
        [ Capa de Aplicación (Services) ]
                       │
                       │ (Aplica lógica de negocio e invariantes)
                       ▼
       [ Capa de Dominio (FCOS Kernel) ]
                       │
                       │ (Valida partida doble, tipos, jerarquía)
                       ▼
     [ Puertos / Repositorios (Interfaces) ]
                       │
                       ├──────────────────────┐
                       ▼                      ▼
             [ Adaptador Memory ]    [ Adaptador Postgres ]
```

## Mapeo de Capas y Directorios
- `internal/domain/`: Contiene el modelo enriquecido de datos del dominio contable (`account.go`, `journal.go`, `ledger.go`) y las invariantes puras de negocio.
- `internal/ports/` u `internal/repository/`: Define los contratos (interfaces de repositorio) que independizan al motor del motor de base de datos final.
- `internal/service/`: Orquesta los casos de uso transaccionales (ej: creación de asientos, cálculo del balance de comprobación).
- `internal/handler/`: Controladores de entrada (HTTP / Wails Bindings) encargados únicamente de validar formatos de entrada y enrutar las peticiones.
"@

Write-Utf8NoBom -FilePath (Join-Path -Path $docsPath -ChildPath "architecture_map.md") -Content $mapContent

# ==========================================
# 3. SEGURIDAD (security.md)
# ==========================================
$securityContent = @"
# Políticas de Seguridad y Auditoría Contable (AuditLog)

La seguridad e integridad del FCOS Kernel descansa sobre tres principios activos:

## 1. Pista de Auditoría Automática (AuditLog)
Toda acción que altere el estado contable de la organización (crear cuenta, autorizar asiento, revertir transacción) debe ser registrada de forma inmutable en la tabla `audit_logs` con los siguientes campos obligatorios:
- `timestamp`: Marca temporal en formato UTC (RFC3339).
- `user_id`: Identificación única del operador que ejecuta la acción.
- `ip_address`: Dirección IP de origen desde donde se generó la petición.
- `action`: Identificador estructurado de la acción (ej: `JOURNAL_POST`, `ACCOUNT_CREATE`).
- `details`: Payload descriptivo en formato JSON que especifica los valores originales y los modificados.

## 2. Transaccionalidad Segura (ACID)
El método `PostEntry` debe estar encapsulado en una transacción única de base de datos (`sql.Tx`). Si la actualización de los saldos de cuentas, la inserción en el libro mayor o el cambio de estado del asiento falla, la transacción completa debe ser revertida mediante un `Rollback()`, garantizando la consistencia total del ledger contable.

## 3. Control Forense del Entorno
Como parte del diagnóstico rutinario de seguridad para aislar el software de inyecciones de DNS locales maliciosas o persistencia sospechosa, se implementan utilidades como `Scanner-Agresivo.ps1` que auditan de manera no destructiva:
- Integridad y firmas del archivo de traducción DNS de Windows `C:\Windows\System32\drivers\etc\hosts`.
- Procesos de fondo que se ejecutan desde carpetas de usuario temporales (`Temp` o `AppData`).
"@

Write-Utf8NoBom -FilePath (Join-Path -Path $docsPath -ChildPath "security.md") -Content $securityContent

# ==========================================
# 4. CONTRATOS GO (contracts.md)
# ==========================================
$contractsContent = @"
# Especificación de Contratos y Puertos en Go

Las capas externas deben comunicarse con el dominio del motor únicamente a través de los puertos definidos. A continuación, se detallan las firmas de interfaces recomendadas para la consistencia de FCOS:

```go
package ports

import (
	"context"
	"github.com/contable-fix/core/internal/domain"
)

// AccountRepository define el contrato para persistir planes y catálogos de cuentas.
type AccountRepository interface {
	Create(ctx context.Context, acc *domain.Account) error
	GetByID(ctx context.Context, id string) (*domain.Account, error)
	GetByCode(ctx context.Context, code string) (*domain.Account, error)
	List(ctx context.Context) ([]*domain.Account, error)
	UpdateBalance(ctx context.Context, id string, deltaCents int64) error
}

// JournalRepository define el contrato de persistencia para el libro diario.
type JournalRepository interface {
	Create(ctx context.Context, entry *domain.JournalEntry) error
	GetByID(ctx context.Context, id string) (*domain.JournalEntry, error)
	UpdateStatus(ctx context.Context, id string, status domain.EntryStatus) error
}

// LedgerRepository gestiona la persistencia de las transacciones mayorizadas.
type LedgerRepository interface {
	PostLines(ctx context.Context, lines []domain.LedgerLine) error
	GetMovements(ctx context.Context, accountID string, periodID string) ([]domain.LedgerLine, error)
}

// UnitOfWork garantiza atomicidad completa entre múltiples repositorios.
type UnitOfWork interface {
	ExecuteTx(ctx context.Context, fn func(adapters *RepositoryAdapters) error) error
}
```
"@

Write-Utf8NoBom -FilePath (Join-Path -Path $docsPath -ChildPath "contracts.md") -Content $contractsContent

# ==========================================
# FASES DE DESARROLLO (docs\phases\*)
# ==========================================

# Fase 1
$phase1Content = @"
# FASE 1: Capa de Persistencia y Transaccionalidad

## Pasos Críticos de Implementación:
1. **Migraciones del Esquema Relacional**: Crear `migrations/schema.sql` definiendo los constraints de integridad referencial.
2. **Implementación de Repositorios**: Diseñar e implementar las consultas SQL nativas en la carpeta `internal/repository/postgres` respetando los contratos de puertos en Go.
3. **Gestión de Transacciones**: Configurar el adaptador de base de datos para manejar `BeginTx` e inyectar el contexto de transacción en los repositorios durante el flujo de posteo.
"@
Write-Utf8NoBom -FilePath (Join-Path -Path $phasesPath -ChildPath "phase1_persistence.md") -Content $phase1Content

# Fase 2
$phase2Content = @"
# FASE 2: Lógica del Plan de Cuentas (account_service.go)

## Pasos Críticos de Implementación:
1. **Jerarquías del Catálogo**: Programar la validación que evita el registro de códigos de cuentas duplicados.
2. **Regla de Mayorización**: Validar recursivamente que si una cuenta tiene subcuentas hijas, su atributo `accepts_move` quede forzado a `false`.
3. **Regla de Desactivación**: Prevenir la desactivación de cuentas contables (`status = 'INACTIVA'`) si su saldo (`current_bal`) es diferente de cero.
"@
Write-Utf8NoBom -FilePath (Join-Path -Path $phasesPath -ChildPath "phase2_chart_of_accounts.md") -Content $phase2Content

# Fase 3
$phase3Content = @"
# FASE 3: Motor de Asientos y Libro Mayor (journal_service.go)

## Pasos Críticos de Implementación:
1. **Validación de Doble Entrada**: Implementar la suma de comprobación de líneas en `CreateDraft` y `PostEntry`.
2. **Mayorización Atómica**: Al contabilizar, insertar las líneas correspondientes de forma indivisible en el libro mayor (`ledger_entries`) y acumular de forma secuencial en el balance saldo.
3. **Procedimiento de Reversión**: Programar `ReverseEntry` para copiar el asiento de origen e invertir matemáticamente las imputaciones (Débitos a Créditos y viceversa).
"@
Write-Utf8NoBom -FilePath (Join-Path -Path $phasesPath -ChildPath "phase3_accounting_engine.md") -Content $phase3Content

# Fase 4
$phase4Content = @"
# FASE 4: Capa de Entrega HTTP, Seguridad y Tests

## Pasos Críticos de Implementación:
1. **Endpoints REST/Wails**: Exponer la API e integrarla en la interfaz de React a través del servicio puente.
2. **Middleware de Trazabilidad**: Interceptar cada petición de modificación y registrarla en el log inmutable.
3. **Pruebas de Estrés Financiero**: Escribir suites de pruebas automatizadas que intenten inyectar desbalances de 0.01 centavos para validar el bloqueo del núcleo.
"@
Write-Utf8NoBom -FilePath (Join-Path -Path $phasesPath -ChildPath "phase4_security_and_testing.md") -Content $phase4Content

Write-Host "`n🎉 ¡PROCESO FINALIZADO CON ÉXITO!" -ForegroundColor Green
Write-Host "La estructura completa de documentación ha sido mapeada en: $docsPath`n" -ForegroundColor Gray
Write-Host "==================================================================" -ForegroundColor Cyan