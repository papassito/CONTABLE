# REQUIREMENTS.md — Especificación de Requerimientos del Sistema (FCOS v2.2)

---

## 1. Requerimientos Funcionales (RF)

### 1.1. Catálogo y Cuentas Contables
* **RF-01.1 — Afectabilidad:** Únicamente las cuentas auxiliares aceptan asientos.
* **RF-01.2 — Aislamiento Multi-Tenant:** Cada entidad multi-tenant estará asociada a un `tenant_id` validado. Se garantiza mediante el paso explícito de `ctx` y filtros `WHERE id = ? AND tenant_id = ?`.

### 1.2. Diario y Libro Mayor
* **RF-02.1 — Partida Doble Rigurosa:** Un asiento será aceptado únicamente si $\sum \text{debit\_cents} = \sum \text{credit\_cents} > 0$ y cada línea cumple las restricciones de no negatividad y exclusión mutua de débito/crédito.
* **RF-02.2 — Flujo de Estados Inmutable:** El estado sigue el flujo `BORRADOR` $\rightarrow$ `CONTABILIZADO`. No existe la transición a `ANULADO`.
* **RF-02.3 — Consistencia de la Fuente de Verdad:** `ledger_entries` es la fuente de verdad. `accounts.balance_cents` es solo una proyección materializada. Ambos se actualizan de manera atómica (ACID).
* **RF-02.4 — Reversión Atomizada:** Toda corrección se realiza mediante `ReverseEntry` en un Unit of Work atómico.

### 1.3. Auditoría e Inmutabilidad
* **RF-03.1 — Serialización Canónica (RFC 8785 / JCS):** Los importes se normalizan como enteros en céntimos antes de convertirlos a JSON canónico y pasarlos por SHA-256.
* **RF-03.2 — Detección de Alteraciones (Tamper-Evident):** La cadena criptográfica local está diseñada para evidenciar y detectar modificaciones no autorizadas comparándola con puntos de control de confianza (Checkpoints).
* **RF-03.3 — Anclaje Externo y Autonomía:** El sistema soportará anclaje asíncrono. La indisponibilidad del anclaje externo **jamás** bloqueará la contabilidad local (Principio Core Offline).
* **RF-03.4 — Protección Anti-Bifurcación:** El sistema aplicará serialización de escritura SQLite con `BEGIN IMMEDIATE` para rechazar reescrituras.

### 1.4. Identidad del Nodo
* **RF-04.1 — Fingerprinting de Entorno:** El contexto `Node` recopilará atributos mínimos (CPU, RAM, Hostname, OS). El fingerprinting operacional utilizará únicamente atributos mínimos necesarios; no se requiere Disk UUID. Esto no constituye atestación de hardware (Hardware Attestation).

---

## 2. Requerimientos No Funcionales (RNF)

* **RNF-01.1 — Cero Aritmética de Flotantes:** Aritmética contable 100% sobre enteros.
* **RNF-01.2 — Despliegue Cgo-Free:** Compilación nativa directa vía Go/Wails.
* **RNF-01.3 — Concurrencia SQLite:** Uso de SQLite en modo WAL con `busy_timeout`.
* **RNF-01.4 — Estrategia de Actualización:** Actualizaciones offline verificadas criptográficamente. Infraestructura OTA estrictamente opcional.

---

## 3. Especificación de Rendimiento (PERF-01)

El siguiente valor constituye un objetivo de rendimiento y no una garantía de capacidad hasta que sea validado mediante benchmarks reproducibles.

* **PERF-01.1 — Target Base:** $\ge 250$ operaciones ACID de escritura por segundo sobre almacenamiento SSD NVMe en el hardware de referencia definido por el benchmark.

La aceptación del objetivo requiere que el benchmark documente como mínimo: hardware, sistema operativo, versión de Go, versión de SQLite/driver, tamaño de DB, tamaño de transacción, número de líneas por asiento, modo WAL, `busy_timeout`, duración, tasa de éxito/error, y percentiles de latencia.