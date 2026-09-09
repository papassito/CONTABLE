# PHASE-06 — Pruebas de Aceptación Integral y Entrega

## 1. Objetivo y Requisitos Cubiertos
Validar la entrega de software bajo escenarios reales de estrés, entornos corporativos restringidos y auditoría independiente.

Requisitos cubiertos:
- **Matriz de Aceptación Completa de AUDIT_GUIDE.md.**
- Verificación de estabilidad y resiliencia local.

## 2. Tareas por Realizar
- [ ] Ejecutar simulación de posteo contable intensivo para cumplir el hito de rendimiento **PERF-01** (mínimo 250 transacciones por segundo sobre SSD).
- [ ] Realizar análisis forense completo de restauración de backups SQLite WAL.
- [ ] Verificar compatibilidad del instalador NSIS empaquetado bajo firmas Authenticode activas.

## 3. Criterios de Aceptación
- Firma del binario validada en el entorno Windows host.
- Pruebas automatizadas de cobertura general superadas en Go (`go test -race ./...`).
- Aprobación inequívoca del reporte forense e informe de auditoría técnica local.

---
*El software se considerará listo para liberación únicamente tras el cumplimiento pleno de esta fase.*