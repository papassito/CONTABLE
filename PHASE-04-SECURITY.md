# PHASE-04 — Seguridad, Bóveda de Secretos y Cortafuegos

## 1. Objetivo y Requisitos Cubiertos
Robustecer las defensas físicas del programa, garantizando aislamiento offline absoluto, protección de claves contables y convivencia estable en el host.

Requisitos cubiertos:
- **SEG-02, SEG-04:** Bóveda de secretos y prevención de fugas de datos en registros.
- **AV-01, AV-02, NET-01 & NET-03:** Protocolo antivirus (System Watcher) y restricción estricta de red local.

## 2. Tareas por Realizar
- [ ] Implementar cifrado de datos críticos (AES-256-GCM) utilizando claves manejadas de manera externa (DPAPI / variables de entorno protegidas).
- [ ] Configurar listeners de API para forzar loopback `127.0.0.1` en desarrollo y desactivarlos por completo en producción.
- [ ] Implementar un manejador de logs seguro que filtre payloads e información sensible de tenants.

## 3. Criterios de Aceptación
- La aplicación debe iniciar y realizar asientos locales en modo avión (sin conectividad de internet).
- Ningún puerto o socket TCP/UDP WAN debe ser abierto por el programa de escritorio en su fase de producción.
- Ningún secreto del kernel debe ser registrado en archivos temporales o de texto plano.

## 4. Decisiones Pendientes de Autorización
- Firma digital e infraestructura de claves públicas para la autenticación de checkpoints.
- Protocolo criptográfico definitivo de validación de licencias corporativas.