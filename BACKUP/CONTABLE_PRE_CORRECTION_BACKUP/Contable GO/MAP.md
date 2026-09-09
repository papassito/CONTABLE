# FCOS v2.2 — Mapa lógico y de persistencia

## Estado del mapa

Este mapa define una organización objetivo. No es un inventario verificado del árbol de código. Los nombres y rutas deben contrastarse con la implementación durante una auditoría técnica.

Repositorio declarado: `https://github.com/papassito/CONTABLE.git`.
Módulo Go declarado: `github.com/klik/fcos-kernel`.

## Contextos y responsabilidades

| Contexto | Responsabilidad | Ruta objetivo relativa al backend |
| --- | --- | --- |
| Accounting | Cuentas, asientos, mayor, balances y períodos contables. | `internal/accounting` |
| Identity | Sesiones, actores y autorización. | `internal/identity` |
| Node | Identidad operativa del nodo y referencias a sus claves. | `internal/node` |
| Tenant | Empresas, acceso y configuración aislada. | `internal/tenant` |
| Normative | Versiones de parámetros normativos y su procedencia. | `internal/normative` |
| Compliance | Obligaciones y expedientes fiscales. | `internal/compliance` |
| Calculation | Cálculos fiscales sujetos a contratos propios. | `internal/calculation` |
| Reconciliation | Conciliación bancaria y documental. | `internal/reconciliation` |
| Documents | Procesamiento y almacenamiento de comprobantes. | `internal/documents` |
| Secrets | Cifrado y acceso autorizado a secretos. | `internal/secrets` |
| Audit | Eventos de auditoría, hash chain y checkpoints. | `internal/audit` |
| Automation | Integraciones externas y RPA asíncrona. | `internal/automation` |
| Health | Diagnóstico y métricas sin secretos. | `internal/health` |
| Messaging | Outbox, entrega y reintentos. | `internal/messaging` |

Se incluye Accounting explícitamente para evitar que el núcleo quede sin propietario en el mapa. Esto no ordena mover carpetas existentes. Audit mantiene evidencia de eventos; no sustituye al mayor contable.

## Relaciones contables

```text
tenant
  |-- accounts
  |-- accounting_periods
  |-- journal_entries
  |      |-- journal_lines        (líneas de borrador; selladas al contabilizar)
  |      |-- ledger_entries       (movimientos contabilizados)
  |      +-- reversal_of_entry_id (referencia a otro asiento del mismo tenant)
  |-- audit_hash_chain
  +-- outbox

accounts <--- journal_lines.account_id
accounts <--- ledger_entries.account_id
journal_lines <--- ledger_entries.journal_line_id
```

## Claves y restricciones objetivo

- `tenant`: clave primaria `id`.
- Entidades contables: identificador estable `id`, `tenant_id` obligatorio y clave única `(tenant_id, id)` para referencias compuestas.
- `accounts`: código único `(tenant_id, code)`; padre referenciado por `(tenant_id, parent_id)`; `balance_cents` entero, proyección de débitos menos créditos.
- `journal_entries`: estado restringido a BORRADOR/CONTABILIZADO; fecha contable y referencia al período; número asignado por backend y único por tenant cuando exista.
- `journal_lines`: FK compuesta al asiento y a la cuenta, siempre incluyendo tenant. Orden de línea único dentro del asiento.
- `ledger_entries`: FK compuesta al asiento, cuenta y línea origen; una sola materialización por línea. Importes enteros no negativos y exactamente un lado positivo.
- `journal_entries.reversal_of_entry_id`: FK compuesta al original con el mismo tenant. Índice único parcial `(tenant_id, reversal_of_entry_id)` cuando la referencia no sea nula. Se prohíbe autorreferencia.
- `accounting_periods`: pertenece al tenant, tiene fechas y estado abierto/cerrado. No permite períodos superpuestos que hagan ambigua la fecha de contabilización.
- `audit_hash_chain`: clave primaria compuesta **`(tenant_id, sequence_id)`**. `sequence_id` no es clave primaria global.
- `outbox`: identificador de evento único, tenant, tipo, payload, estado de entrega e información acotada de reintentos.

Las FK que contienen tenant evitan asociaciones cruzadas en persistencia. No sustituyen la autorización de consultas ni las reglas del servicio. Deben habilitarse y comprobarse en cada conexión SQLite.

La igualdad de sumas entre múltiples líneas, el control de períodos y la jerarquía requieren validaciones transaccionales; no se consideran resueltos solo por los CHECK de cada fila.

## Evidencia criptográfica

`audit_hash_chain` conserva secuencia, tenant, hash previo, payload canónico, hash de payload, hash de cadena y versión de esquema. Los nombres y formatos definitivos siguen SECURITY.

El mayor explica el impacto financiero. La cadena registra evidencia de las operaciones. Ambos se confirman en la misma Unit of Work, pero tienen funciones distintas.

## Estructura objetivo

```text
CONTABLE/
|-- README.md
|-- Contable GO/
|   |-- README.md
|   |-- ARCHITECTURE.md
|   |-- REQUIREMENTS.md
|   |-- MAP.md
|   |-- SECURITY.md
|   +-- AUDIT_GUIDE.md
|-- app.go                 # Adaptador desktop, si esta es su ubicación efectiva.
|-- main.go
|-- wails.json
|-- go-skeleton/
|   |-- go.mod
|   |-- go.sum
|   |-- cmd/api/            # Transporte HTTP opcional.
|   +-- internal/          # Contextos y adaptadores.
|-- src/                   # Frontend; confirmar ubicación efectiva.
|-- migrations/            # Migraciones versionadas; confirmar ubicación efectiva.
|-- tests/
|-- scripts/
|-- .github/workflows/
|-- package.json
|-- <archivo de bloqueo del gestor elegido>
+-- wailsjs/                # Bindings generados; confirmar ubicación efectiva.
```

No se debe crear, mover ni borrar código únicamente para hacer coincidir este dibujo. Primero se documentan las rutas reales y cualquier diferencia arquitectónica relevante. La ubicación de archivos generados y artefactos de build la determina la configuración del proyecto.
