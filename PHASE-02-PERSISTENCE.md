# PHASE-02 — Persistencia y Aislamiento SQLite

## 1. Objetivo y Requisitos Cubiertos
Implementar la persistencia de datos contables garantizando integridad referencial compuesta y aislamiento estricto multi-tenant en SQLite.

Requisitos cubiertos:
- **OPS-01 & OPS-03:** Configuración robusta de SQLite en modo WAL y habilitación de claves foráneas.
- **RF-01:** Garantías de aislamiento por TenantID en todas las operaciones.
- **ACID:** Implementación del patrón Unit of Work (`sql.Tx`).

## 2. Tareas por Realizar
- [ ] Crear migraciones de base de datos para SQLite local (`migrations/schema.sql`).
- [ ] Diseñar triggers de protección inmutable contra `UPDATE`/`DELETE` en asientos contabilizados.
- [ ] Crear adaptador `UnitOfWork` que asegure transacciones con aislamiento y bloqueo `BEGIN IMMEDIATE` para evitar concurrencia conflictiva.

## 3. Criterios de Aceptación
- Un fallo inyectado a mitad de una transacción de posteo debe revertir todos los cambios y balances mediante `Rollback()`.
- Consultas externas maliciosas que envíen un `TenantID` alternativo deben ser rechazadas con error de autorización.
- Múltiples hilos escribiendo concurrentemente no deben corromper la base de datos SQLite WAL.

## 4. Pruebas de Verificación Previstas
- Pruebas de integración concurrentes con hilos Go.
- Tests unitarios con bases de datos transaccionales en memoria (`:memory:`).

## 5. Bloqueos Técnicos
- *Ninguno.*