# FCOS v2.2 — Plan de Reconstrucción y Desarrollo Operativo

Este plan define el mapa de ruta de ingeniería para la estabilización y despliegue del core financiero de Contable Fix, subordinado en su totalidad al baseline normativo sellado en la raíz del repositorio.

## Metodología de Integración Incremental
Para garantizar la consistencia, el desarrollo se divide en 6 fases secuenciales. Cada fase se apoya sobre las garantías lógicas y pruebas aprobadas de la anterior:

```text
[ Fase 1: Dominio ] ──► [ Fase 2: Persistencia ] ──► [ Fase 3: Transaccional ]
                                                            │
 [ Fase 6: Aceptación ] ◄── [ Fase 5: Desktop ] ◄── [ Fase 4: Seguridad ]
```

## Fases del Plan Operativo

### PHASE-01-DOMAIN.md — Dominio y Reglas Puras
- **Objetivo:** Implementación del tipo monetario en centavos exactos (`int64`), estructura jerárquica de cuentas sin ciclos, asientos de diario en borrador y puertos de repositorio.
- **Verificación:** Pruebas unitarias puras en Go sin base de datos, red ni Wails.

### PHASE-02-PERSISTENCE.md — Persistencia y Aislamiento SQLite
- **Objetivo:** Configuración de SQLite local en modo WAL, claves foráneas, aislamiento de datos por `TenantID` y transaccionalidad mediante Unit of Work.
- **Verificación:** Pruebas de rollback y escrituras concurrentes.

### PHASE-03-ACCOUNTING.md — Libro Mayor y Motor Contable
- **Objetivo:** Posteo definitivo, reversión matemática cruzada, estados de periodos por tenant y auditoría local encadenada en `audit_hash_chain`.
- **Verificación:** Test de balance de comprobación coincidente bit a bit con movimientos del diario.

### PHASE-04-SECURITY.md — Capa de Seguridad y Compatibilidad
- **Objetivo:** Cifrado de secretos (AES-256-GCM), custodia de DEKs envueltas (RSA-4096), mitigación de falsos positivos en System Watcher y aislamiento offline total.
- **Verificación:** Auditoría de sockets de red locales y simulación de comportamiento host.

### PHASE-05-DESKTOP.md — Integración Frontend y Wails
- **Objetivo:** Sincronización del IPC con React/TypeScript mediante bindings de Wails, aplicando el transporte monetario exacto de centavos como cadenas de texto.
- **Verificación:** Ciclo completo de entrada/salida de datos desde la interfaz de usuario.

### PHASE-06-ACCEPTANCE.md — Pruebas de Aceptación y Despliegue
- **Objetivo:** Cumplimiento de la matriz de validación detallada en el documento de auditoría.
- **Verificación:** Pruebas de resistencia ante cortes bruscos, restauración forense y empaquetamiento NSIS con firmas activas.

---
*Gobernanza: Si una fase entra en contradicción con REQUIREMENTS.md o ARCHITECTURE.md, la fase debe corregirse para alinearse al baseline.*