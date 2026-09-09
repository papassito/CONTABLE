# PHASE-01 — Capa de Dominio Contable Pura

## 1. Objetivo y Requisitos Cubiertos
Implementar las invariantes contables elementales y la precisión matemática exigida por el núcleo operativo, asegurando la no regresión financiera.

Requisitos normativos cubiertos (REQUIREMENTS.md):
- **NUM-01 & NUM-03:** Representación exacta en centavos enteros (`int64`).
- **NUM-02:** Deserialización y serialización JSON determinista basada exclusivamente en cadenas de texto que representan el entero de centavos sin formato.
- **RF-02 & RF-03:** Restricciones de unicidad jerárquica de cuentas (detección de bucles/ciclos) y protección de cuentas padre acumuladoras.
- **RF-06:** Regla de oro de partida doble inmutable en asientos de diario.

## 2. Estado Real del Código Existente
- El código se encuentra estructurado en el paquete `internal/domain/` como lógica pura libre de dependencias de persistencia (SQL) o de red.
- Se incorporó validación sintáctica estricta para garantizar la canonicidad del formato de entrada en JSON.

## 3. Criterios de Aceptación de la Fase
- [ ] Todas las pruebas en `money_test.go` pasan con éxito, incluyendo el rechazo de formatos no canónicos (`-0`, ceros iniciales, fracciones, exponentes).
- [ ] Las pruebas de catálogo de cuentas detectan de forma recursiva cualquier ciclo infinito de relaciones padre-hijo.
- [ ] El validador del diario bloquea cualquier asimetría de partida doble, importes negativos o líneas sin imputación definida.

## 4. Archivos Modificados
- `internal/domain/money.go` (Validación canónica NUM-02)
- `internal/domain/money_test.go` (Tests de robustez sintáctica)

## 5. Pruebas Ejecutadas y Resultados
- *Pruebas Unitarias:* Pendientes de ejecución física en la terminal local por parte del operador del host.

## 6. Pendientes o Bloqueos
- **Sincronización:** Confirmación táctil y ejecución del comando de pruebas en el entorno de desarrollo local.