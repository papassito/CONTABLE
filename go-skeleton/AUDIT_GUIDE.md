# DICTAMEN DE AUDITORÍA Y GUÍA TÉCNICA
## Proyecto: Contable Fix by KLIK
**Rol:** Auditor Técnico & Arquitecto de Software
**Destinatario:** Equipo de Desarrollo Backend en Go
**Fecha de Emisión:** 2026

---

### 1. REGLAS CONTABLES NO NEGOCIABLES (INVARIANTES DEL SISTEMA)

1. **Invariante de Partida Doble:**
   - Para todo comprobante o asiento contable (`JournalEntry`), la suma total de débitos DEBE ser idéntica a la suma total de créditos (`∑ Débitos == ∑ Créditos`).
   - El sistema debe rechazar cualquier intento de contabilizar un asiento desbalanceado con `ErrUnbalancedJournal`.

2. **Inmutabilidad de Asientos Contabilizados:**
   - Una vez que un asiento pasa al estado `CONTABILIZADO`, está terminantemente **PROHIBIDO** ejecutar sentencias `UPDATE` o `DELETE` sobre el asiento o sus líneas.
   - Si se requiere corregir un error, se debe generar un asiento de reversión/anulación con contrapartida (`ReverseEntry()`), dejando trazabilidad del motivo y referencia al asiento original.

3. **Restricción de Cuentas Auxiliares:**
   - Las transacciones en las líneas del asiento (`JournalLine`) solo pueden asociarse a cuentas auxiliares que tengan `accepts_move = true`.
   - Queda prohibido imputar movimientos a cuentas de nivel superior o mayorizadoras (Cuentas, Subcuentas no auxiliares).

4. **Precisión Numérica y Manejo Monetario:**
   - Evitar `float64` para saldos acumulados en producción debido a errores de redondeo de punto flotante de IEEE 754.
   - Recomendación: Utilizar la librería `github.com/shopspring/decimal` o enteros de centavos (`int64`), mapeados a columnas `NUMERIC(18, 4)` en PostgreSQL.

---

### 2. HOJA DE RUTA DE IMPLEMENTACIÓN EN 4 FASES

#### FASE 1: Capa de Persistencia y Transaccionalidad
- **Objetivo:** Implementar la interfaz de repositorios en `internal/repository` utilizando PostgreSQL (o base de datos relacional elegida).
- **Entregables:**
  - Migraciones DDL (`schema.sql`) para tablas `accounts`, `journal_entries`, `journal_lines`, `ledger_entries` y `audit_logs`.
  - Soporte de transacciones atómicas `*sql.Tx`: la operación `PostEntry()` debe actualizar el asiento, insertar en el mayor y actualizar saldos dentro de la misma transacción.

#### FASE 2: Lógica del Plan de Cuentas (`account_service.go`)
- **Objetivo:** Garantizar la jerarquía del plan contable (Clase, Grupo, Cuenta, Subcuenta, Auxiliar).
- **Entregables:**
  - Validación de unicidad de códigos (`code`).
  - Verificación de árbol: si una cuenta tiene hijos, no puede tener `accepts_move = true`.
  - Algoritmo de desactivación: no permitir desactivar cuentas que posean saldo distinto de cero.

#### FASE 3: Motor de Asientos y Libro Mayor (`journal_service.go` & `ledger_service.go`)
- **Objetivo:** Ejecutar la lógica contable central.
- **Entregables:**
  - `CreateDraft`: Almacena el borrador validando que posea al menos dos líneas contables y cuentas válidas.
  - `PostEntry`: Valida partida doble, transfiere movimientos a `LedgerRepository`, actualiza el saldo de cada cuenta y cambia el estado a `CONTABILIZADO`.
  - `ReverseEntry`: Crea un nuevo asiento con los débitos y créditos invertidos y concepto de reversión.
  - `GenerateTrialBalance`: Consulta el balance de comprobación sumas y saldos.

#### FASE 4: Capa de Entrega HTTP, Seguridad y Tests
- **Objetivo:** Exponer endpoints REST robustos y protegidos.
- **Entregables:**
  - Middleware de autenticación y autorización por roles (Administrador, Contador, Auditor).
  - Middleware de auditoría automática que alimente `internal/domain/audit.go` en cada operación de escritura.
  - Suite de pruebas unitarias (`journal_service_test.go`, `validator_test.go`) con cobertura de casos de borde: desbalanceos de 0.01 centavos, periodos cerrados, cuentas bloqueadas.

---

### 3. CHECKLIST DE APROBACIÓN DEL AUDITOR

- [ ] Todas las funciones marcadas con `// TODO:` implementadas sin fallos de compilación (`go build`).
- [ ] Validación de partida doble probada con test unitario exhaustivo.
- [ ] Transacciones atómicas (rollback ante fallos parciales durante el posteo de asientos).
- [ ] No existen sentencias `DELETE` o `UPDATE` destructivas en asientos contabilizados.
- [ ] Pista de auditoría (`AuditLog`) registrando usuario, IP y timestamp en cada cambio.
