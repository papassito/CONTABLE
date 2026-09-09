# PHASE-05 — Integración del Escritorio (Wails & React)

## 1. Objetivo y Requisitos Cubiertos
Vincular la interfaz visual de React con la capa de servicios transaccionales del backend, manteniendo precisión numérica exacta al cruzar el puente de comunicación.

Requisitos cubiertos:
- **NUM-02:** Intercambio de importes monetarios como cadenas decimales.
- **RF-05:** Gestión interactiva de catálogos y edición ágil de borradores de diario.

## 2. Tareas por Realizar
- [ ] Implementar el adaptador e inyector de contextos de sesión de Wails en Go.
- [ ] Desarrollar en TypeScript conversores y validadores numéricos que traten los centavos como `string` o `bigint`, protegiendo la UI de desbordamientos.
- [ ] Conectar las vistas del diario, el catálogo y balances a las llamadas de servicio correspondientes del backend.

## 3. Criterios de Aceptación
- Valores monetarios extremadamente grandes representados en la UI no deben truncarse ni sufrir redondeos indeseados en el navegador local.
- El paso de variables context de Wails debe determinar dinámicamente el tenant activo sin vulnerabilidades de inyección del cliente.

## 4. Pruebas de Verificación Previstas
- Suite de pruebas unitarias en Jest/Vitest para los validadores de cadenas monetarias en TypeScript.