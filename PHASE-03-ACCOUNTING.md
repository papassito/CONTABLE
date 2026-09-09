# PHASE-03 — Motor de Contabilización y Mayor

## 1. Objetivo y Requisitos Cubiertos
Implementar los casos de uso transaccionales clave del negocio contable, vinculando la actividad financiera a la cadena inmutable de auditoría forense local.

Requisitos cubiertos:
- **RF-07, RF-08 & RF-09:** Reglas de posteo definitivo, reversión matemática con motivo explícito e inmutabilidad absoluta.
- **RF-10, RF-11 & RF-13:** Reconstrucción íntegra de saldos a partir del libro mayor, control de períodos cerrados y cadena de bloques local en `audit_hash_chain`.

## 2. Tareas por Realizar
- [ ] Crear caso de uso `PostEntry(entryID)` que transicione borradores a contabilizados, actualizando proyecciones de saldos de forma indivisible.
- [ ] Crear caso de uso `ReverseEntry(entryID, reason)` para generar un contraasiento exacto.
- [ ] Implementar el motor canónico de hash (`audit_hash_chain`) siguiendo las reglas JCS (RFC 8785) y hashing SHA-256.

## 3. Criterios de Aceptación
- Intentar editar un asiento con estado `CONTABILIZADO` debe retornar un error fatal.
- Un desbalance detectado en los hashes de la cadena de auditoría debe abortar de inmediato futuras escrituras contables (`AUDIT_INTEGRITY_FAILURE`).
- Intentar contabilizar en fechas fuera de los rangos de períodos activos debe bloquearse.

## 4. Pruebas de Verificación Previstas
- Simulación forense de inyección y manipulación de datos en la base de datos para probar la detección automática de roturas en la cadena.