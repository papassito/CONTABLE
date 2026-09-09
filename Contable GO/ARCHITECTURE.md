# ARCHITECTURE.md — Arquitectura, Patrones y Principios de Diseño (FCOS v2.2)

---

## 1. Visión Holística y Core Offline Boundary

**FCOS v2.2** está diseñado bajo los principios de **Domain-Driven Design (DDD)** y **Clean Architecture**.

El principio arquitectónico fundamental es el **Core Offline Boundary**: Los casos de uso financieros del Core jamás realizan llamadas síncronas a sistemas externos. Toda integración externa se ejecuta de forma asíncrona mediante el patrón *Transactional Outbox* y procesos desacoplados.

La indisponibilidad de Internet, del SAT, de redes Web3 o de herramientas RPA **no puede impedir**:
* Crear borradores contables.
* Contabilizar asientos.
* Consultar el catálogo de cuentas y el libro mayor.
* Generar la cadena de auditoría criptográfica local.
* Ejecutar reversiones y contraasientos.

```text
                     +-----------------------------------+
                     |          FCOS KERNEL              |
                     +-----------------------------------+
                                       |
                   +-------------------+-------------------+
                   |                                       |
                   v                                       v
     +---------------------------+           +---------------------------+
     |      CORE ACCOUNTING      |           |     AUTOMATION & RPA      |
     |   (100% Offline / Local)  |           |   (External Systems / SAT)|
     +---------------------------+           +---------------------------+
                   |                                       |
                   +-------------------+-------------------+
                                       |
                                       v
                     +-----------------------------------+
                     |        TRANSACTIONAL OUTBOX       |
                     +-----------------------------------+
2. Invariantes del Núcleo Financiero (Reglas Formales)A. Dominio Entero Absoluto (Integer/Cents Rule)Queda estrictamente prohibido el uso de tipos de punto flotante (float32, float64) para el dominio financiero en Go.Regla de Equivalencia: $\$1.00 \text{ MXN} = 100 \text{ céntimos}$.Almacenamiento: Todo monto se procesa como int64.B. Partida Doble Inviolable (Especificación Matemática Ejecutable)Un asiento contable $E$ está compuesto por un conjunto ordenado de $n$ líneas:$L = \{l_1, l_2, \dots, l_n\}$Un asiento puede transicionar a estado CONTABILIZADO si y solo si cumple simultáneamente las siguientes restricciones:1. Cantidad mínima de líneas$$n \ge 2$$2. No negatividadPara toda línea $l_i$:$$\text{debit\_cents}(l_i) \ge 0$$$$\text{credit\_cents}(l_i) \ge 0$$3. Aislamiento de débito y créditoPara toda línea $l_i$:$$\text{debit\_cents}(l_i) = 0 \quad \lor \quad \text{credit\_cents}(l_i) = 0$$(Por tanto, una misma línea nunca puede contener simultáneamente un débito y un crédito mayores que cero).4. Balance de partida doble$$\sum_{i=1}^{n} \text{debit\_cents}(l_i) = \sum_{i=1}^{n} \text{credit\_cents}(l_i)$$5. Rechazo de asiento nulo$$\sum_{i=1}^{n} \text{debit\_cents}(l_i) > 0$$En consecuencia, la condición fundamental e indivisible es:$$\sum_{i=1}^{n} \text{debit\_cents}(l_i) = \sum_{i=1}^{n} \text{credit\_cents}(l_i) > 0$$C. Máquina de Estados del Asiento ContableEl estado persistente de un asiento contable sigue un flujo unidireccional y sin retorno:$$\text{BORRADOR} \longrightarrow \text{CONTABILIZADO}$$Inmutabilidad Absoluta: Un asiento en estado CONTABILIZADO queda sellado permanentemente en la base de datos. Se prohíbe cualquier operación SQL UPDATE sobre sus montos/líneas o DELETE sobre su registro.Prohibición de Estado Perturbador (ANULADO): Un asiento contabilizado jamás cambia su columna status a ANULADO. La corrección se realiza mediante la creación de un nuevo asiento contable de reversión, también en estado CONTABILIZADO, que vincula formalmente el ID del asiento original mediante la columna reversal_of_entry_id.Presentación UI: Si la interfaz gráfica requiere indicar visualmente "Anulado", esta bandera se deriva dinámicamente mediante la presencia de un contraasiento asociado, sin alterar el estado persistido inmutable del registro original.D. Ledger como Única Fuente de Verdad (Source of Truth)La tabla ledger_entries es la única e inmutable fuente de verdad contable del sistema.El campo balance_cents en la tabla accounts es exclusivamente una proyección materializada calculada para optimizar lecturas.Regla: balance_cents jamás podrá utilizarse para reconstruir movimientos históricos ni contradecir al ledger.Toda modificación en ledger_entries y actualización de balance_cents ocurre estrictamente en la misma transacción ACID.3. Persistencia y Concurrencia (SQLite WAL)Driver Cgo-Free: Se utiliza modernc.org/sqlite (Go puro).Modo WAL (Write-Ahead Logging): Inicialización con pragmas:file:fcos_local.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)Patrón Unit of Work (UoW): Transacciones atómicas sql.Tx garantizan integridad cruzada (Account Balance + Ledger + Audit Chain).
---

### 3. `CONTRACTS.md`

```markdown
# CONTRACTS.md — Especificación de Contratos, Interfaces y DTOs

---

## 1. Contratos del Puente IPC (Wails Frontend $\leftrightarrow$ Backend)

### A. Límites Numéricos de Transporte (IPC ↔ Go Domain)
En JavaScript/TypeScript, los campos `debit_cents` y `credit_cents` cruzan el puente IPC como `number` (IEEE 754 de doble precisión), mientras que en el dominio Go operan como `int64`.

* Límite seguro en IPC (JavaScript): `Number.MAX_SAFE_INTEGER` ($9,007,199,254,740,991$)
* Límite nativo en Go: `math.MaxInt64` ($9,223,372,036,854,775,807$)

**Contrato Arquitectónico:**
El límite efectivo del sistema es $\min(\text{IPC\_SAFE\_INTEGER}, \text{INT64\_MAX})$. Cualquier valor entrante que exceda `Number.MAX_SAFE_INTEGER` es rechazado en la frontera IPC antes de invocar la capa Go.

### B. DTOs de Dominio (TypeScript)

```typescript
export type EntryStatus = 'BORRADOR' | 'CONTABILIZADO';

export interface JournalLine {
  account_id: string;
  account_code: string;
  description: string;
  debit_cents: number;  // Entero en céntimos (v ∈ ℤ, v ≥ 0)
  credit_cents: number; // Entero en céntimos (v ∈ ℤ, v ≥ 0)
  third_party_id?: string;
}

export interface JournalEntry {
  id?: string;
  number: string;
  date: string; // ISO 8601 / RFC 3339 string
  concept: string;
  reference: string;
  status?: EntryStatus;
  reversal_of_entry_id?: string;
  lines: JournalLine[];
}
C. Firmas de Métodos IPC y Aislamiento por Contexto (Go Exportados)Todos los casos de uso reciben obligatoriamente ctx context.Context, el cual transporta la identidad validada del tenant_id.Gotype JournalService interface {
    CreateDraft(ctx context.Context, entry domain.JournalEntry) (*domain.JournalEntry, error)
    PostEntry(ctx context.Context, entryID string) error
    ReverseEntry(ctx context.Context, entryID string, reason string) (*domain.JournalEntry, error)
    GetEntry(ctx context.Context, entryID string) (*domain.JournalEntry, error)
}
D. Protocolo Transaccional Atómico de ReverseEntryLa invocación de ReverseEntry ejecuta un pipeline indivisible dentro de una única transacción UnitOfWork (ACID). Si cualquier paso falla, se ejecuta ROLLBACK completo:Validar Contexto de Seguridad: Extraer y verificar la presencia de tenant_id en ctx.Validar Existencia y Tenant Match: WHERE id = entryID AND tenant_id = tenant_id.Validar Estado: Confirmar que el asiento esté CONTABILIZADO.Validar No Reversión Previa: Verificar que no exista ya un contraasiento que apunte a entryID en reversal_of_entry_id.Validar Justificación: Garantizar que reason no sea vacío.Invertir Líneas: debit_cents_new = credit_cents_orig y credit_cents_new = debit_cents_orig.Garantía de Pertenece al Mismo Tenant: Asignar explícitamente el tenant_id extraído del ctx al nuevo asiento.Generar Asiento de Reversión: Crear la entidad JournalEntry con un nuevo ID, status = CONTABILIZADO y asignando reversal_of_entry_id = entryID.Asentar en Libro Mayor: Insertar en ledger_entries.Actualizar Proyección: Recalcular accounts.balance_cents.Registrar Evento de Auditoría: Construir el bloque en audit_hash_chain.Commit Transaccional.2. Contratos REST API Operacional (/api/v1)Todas las llamadas exigen la propagación del header de inquilino o JWT Bearer validado.POST /api/v1/accountsGET /api/v1/accountsPOST /api/v1/journal-entriesPOST /api/v1/journal-entries/{id}/postPOST /api/v1/journal-entries/{id}/reverse3. Contratos REST Debug / Diagnóstico (/debug/architecture)Ruta: GET /debug/architecture/filesPolíticas de Acceso: LOCAL_ONLY | DEBUG_ONLY | AUTHENTICATED | DISABLED_IN_PRODUCTIONQueda estrictamente prohibido exponer endpoints de inspección de código fuente dentro de la API operacional /api/v1.