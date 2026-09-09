# FCOS v2.2 — Mapa lógico y de persistencia

## Estado del mapa

Este mapa distingue una organización objetivo de una sección de inventario observado. La organización objetivo no certifica el árbol de código. Los nombres y rutas deben contrastarse con la implementación durante una auditoría técnica.

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

## Estructura observada y destino documental

Inventario de lectura realizado el 2026-09-09: seis Markdown en raíz y 86 archivos Go. Esto confirma presencia, no compilación ni corrección.

```text
CONTABLE/
|-- README.md
|-- ARCHITECTURE.md
|-- REQUIREMENTS.md
|-- MAP.md
|-- SECURITY.md
|-- AUDIT_GUIDE.md
|-- app.go / main.go
|-- go.mod / go.sum / go.work
|-- wails.json / package.json / vite.config.ts
|-- cmd/api/
|-- config/
|-- internal/                 # Dominio, servicios, repositorios y contextos existentes.
|-- pkg/database/
|-- pkg/response/
|-- pkg/validator/
|-- src/                     # React, componentes, datos y utilidades.
|-- public/
|-- wailsService.ts          # Puente presente en raíz; referencias aún por reconciliar.
|-- *.go                    # Incluye archivos de prueba en raíz.
+-- *.ps1                   # Herramientas existentes; no todas aprobadas para ejecutar.
```

Los directorios objetivo de la tabla de contextos se interpretan respecto de la raíz del módulo, no de `go-skeleton/`. No debe recrearse esa carpeta para resolver referencias antiguas.

Se observaron imports de componentes React hacia `src/services/wailsService`, mientras el archivo visible está en raíz. Es una discrepancia pendiente de código, no una orden de moverlo durante esta entrega.

Hay modelos con nombres similares en `internal/domain` y en contextos especializados. No se clasifican como duplicados eliminables sin analizar imports y responsabilidades.

No se detectó una carpeta de migraciones en el árbol revisado. Su ubicación debe definirse en la reconstrucción; los SQL en `src/data` no se consideran migraciones oficiales por su nombre.

## Navegación

- [README](README.md): alcance y preparación.
- [Arquitectura](ARCHITECTURE.md): contratos y transacciones.
- [Requisitos](REQUIREMENTS.md): aceptación.
- [Seguridad](SECURITY.md): controles y compatibilidad AV/firewall.
- [Auditoría](AUDIT_GUIDE.md): evidencias exigidas.