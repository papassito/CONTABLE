# MAP.md — Mapeo Maestro de Arquitectura y Bounded Contexts

---

## 1. Identificadores Oficiales

* **Repositorio de Código (Git):** `https://github.com/papassito/CONTABLE.git`
* **Módulo Principal de Go (`go.mod`):** `github.com/klik/fcos-kernel`

---

## 2. Mapeo de los 13 Bounded Contexts (`go-skeleton/internal/`)

| Bounded Context | Directorio Físico | Responsabilidad Semántica de Dominio |
| :--- | :--- | :--- |
| **Identity** | `internal/identity` | Claves asimétricas de firma y credenciales de acceso. |
| **Node** | `internal/node` | Identidad del nodo local, fingerprinting de entorno y métricas. |
| **Tenant** | `internal/tenant` | Identidad del inquilino, ciclo de vida de aislación y configuración general. |
| **Normative** | `internal/normative` | Parámetros oficiales inmutables y reglas normativas. |
| **Compliance** | `internal/compliance` | Expedientes, obligaciones fiscales y periodos de declaración. |
| **Calculation** | `internal/calculation` | Motor algebraico de cálculo de impuestos y retenciones. |
| **Reconciliation** | `internal/reconciliation` | Conciliación de estados de cuenta bancarios contra XMLs. |
| **Documents** | `internal/documents` | Procesamiento, indexación y almacenamiento local de comprobantes. |
| **Secrets** | `internal/secrets` | Encriptación de credenciales mediante sobres cifrados (DEK/KEK). |
| **Audit** | `internal/audit` | Ledger inmutable local y anclaje externo opcional. |
| **Automation** | `internal/automation` | Orquestación RPA asíncrona, conectores SAT, parsing Offline. |
| **Health** | `internal/health` | Telemetría del sistema y health checks. |
| **Messaging** | `internal/messaging` | Outbox transaccional para comunicación asíncrona. |

---

## 3. Esquema Relacional de Persistencia y Claves Foráneas (SQLite - WAL)

```text
                           +-------------------+
                           |      tenant       |
                           +-------------------+
                           | id (PK)           |
                           +-------------------+
                             |       |       |
         +-------------------+       |       +-------------------+
         | (1:N)                     | (1:N)                     | (1:N)
         v                           v                           v
+-------------------+       +-------------------+       +-------------------+
|     accounts      |       |  journal_entries  |       | audit_hash_chain  |
|   (Projection)    |       +-------------------+       +-------------------+
+-------------------+       | id (PK)           |       | sequence_id (PK)  |
| id (PK)           |       | tenant_id (FK)    |       | tenant_id (FK)    |
| tenant_id (FK)    |       | status (VARCHAR)  |       | canonical_hash    |
| balance_cents     |       | reversal_of...    |       | chain_hash        |
+-------------------+       +-------------------+       +-------------------+
         |                           |                    UNIQUE(tenant_id,
         | (1:N)                     | (1:N)                     sequence_id)
         +-------------------+       |
                             |       |
                             v       v
                    +---------------------------+
                    |      ledger_entries       |
                    |     (Source of Truth)     |
                    +---------------------------+
                    | id (PK)                   |
                    | journal_entry_id (FK)     |
                    | account_id (FK)           |
                    | debit_cents (BIGINT)      |
                    | credit_cents (BIGINT)     |
                    +---------------------------+
Declaración Explícita de Claves Foráneas (FK Constraints)ledger_entries.account_id $\longrightarrow$ accounts.idledger_entries.journal_entry_id $\longrightarrow$ journal_entries.idaccounts.tenant_id $\longrightarrow$ tenant.idjournal_entries.tenant_id $\longrightarrow$ tenant.idaudit_hash_chain.tenant_id $\longrightarrow$ tenant.idjournal_entries.reversal_of_entry_id $\longrightarrow$ journal_entries.idDistinción Semántica de ComponentesJournalEntry (Documento Contable): Representa el hecho económico formal. Mantiene el estado inmutable (BORRADOR/CONTABILIZADO) y el concepto.LedgerEntry (Movimiento Contable Efectivo): Representa el impacto numérico individual en el Libro Mayor. Es la fuente de verdad.Account (Catálogo y Proyección): Representa la estructura de cuentas. balance_cents es una proyección materializada derivada.AuditHashChain (Evidencia Criptográfica): Prueba criptográfica Tamper-Evident de los eventos ocurridos.4. Árbol Físico del ProyectoPlaintextCONTABLE/
├── app.go                          # IPC Bindings Wails v2
├── main.go                         # Launcher de la aplicación desktop
├── wails.json                      # Configuración de compilación global
├── README.md                       
├── ARCHITECTURE.md                 
├── CONTRACTS.md                    
├── REQUIREMENTS.md                 
├── MAP.md                          
├── SECURITY.md                     
├── verificar_integridad.ps1        # Script PowerShell de auditoría de integridad del filesystem
├── go-skeleton/                    # Backend Kernel Go ([github.com/klik/fcos-kernel](https://github.com/klik/fcos-kernel))
│   ├── go.mod                      
│   ├── go.sum
│   ├── run-tests.ps1               
│   ├── cmd/api/main.go             
│   └── internal/                   
└── src/                            # Interfaz React 18
    ├── App.tsx                     
    ├── services/wailsService.ts    
    └── wailsjs/                    
