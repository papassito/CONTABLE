# FCOS v2.2 — Guía de auditoría técnica

## 1. Naturaleza del documento

Esta es una guía de evaluación, no un dictamen de aprobación. La revisión documental que la origina no inspeccionó implementación, ejecutó pruebas ni certificó seguridad o cumplimiento fiscal.

La auditoría técnica debe contrastar REQUIREMENTS, ARCHITECTURE, MAP y SECURITY contra una revisión identificable del código. Cualquier diferencia se registra como hallazgo; no se modifica el código automáticamente para coincidir con un dibujo o una ruta propuesta.

## 2. Alcance y evidencia inicial

Registrar fecha, auditor, repositorio, commit, cambios locales, versiones de herramientas, sistema operativo, base utilizada y exclusiones. Trabajar con datos sintéticos y una base de prueba separada de la información real.

Confirmar:

- Ubicación efectiva del módulo Go, frontend, migraciones y configuración desktop.
- Driver y pragmas reales de SQLite.
- Superficies activas de Wails, HTTP, diagnóstico y workers.
- Diferencias entre componentes implementados, parciales, previstos y ausentes.
- Compatibilidad del formato monetario y de auditoría con datos o clientes existentes.

No ejecutar scripts desconocidos ni migraciones sobre una base del usuario como parte de una auditoría de lectura.

## 3. Revisión de arquitectura

Verificar que el dominio no dependa de transporte o SQL; que las reglas y autorización se ejecuten en backend; que los adaptadores utilicen la misma Unit of Work y que no existan llamadas externas dentro de las transacciones contables.

Comprobar que la outbox se confirma junto con el evento origen y que sus consumidores admiten reintentos sin duplicar efectos.

Revisar que el mapa refleje las responsabilidades reales. No exigir una reestructuración de directorios si la separación ya es verificable con otra organización.

## 4. Matriz mínima de pruebas

| Área | Casos exigidos | Evidencia |
| --- | --- | --- |
| Partida doble | Válido; desbalance de 1 centavo; cero; una línea; línea con ambos lados positivos; línea nula; negativos. | Resultado y ausencia de cambios al rechazar. |
| Borradores | Guardado incompleto, edición y rechazo de edición tras contabilización. | Mayor y saldos intactos antes del posteo. |
| Dinero | Frontera `int64`, acumulación fuera de rango, importes mayores al entero seguro JS, cadenas inválidas. | Exactitud de ida y vuelta; rechazo sin truncamiento. |
| Cuentas | Duplicados por tenant, jerarquía cíclica, agrupadora, inactiva, saldo no cero al desactivar. | Restricciones y errores de dominio. |
| Tenants | Lectura y escritura cruzadas, referencias de cuentas y reversión ajenas, selección de tenant manipulada. | Denegación sin fuga de datos. |
| Inmutabilidad | UPDATE/DELETE por rutas normales; posteo repetido. | Original intacto y movimientos sin duplicación. |
| Reversión | Motivo vacío, fecha cerrada, asiento no contabilizado, dos solicitudes simultáneas. | Un único contraasiento directo y atomicidad. |
| Períodos | Fecha sin período, cerrado y carrera cierre/posteo. | Ninguna contabilización inválida confirmada. |
| Transacciones | Fallo inyectado en mayor, saldo, estado, auditoría y outbox. | Rollback de todos los efectos asociados. |
| Mayor y balance | Reconstrucción de saldo, filtros de fecha/tenant y sumas globales. | Igualdad exacta entre fuente y proyección. |
| Auditoría | Primer bloque, secuencias de dos tenants, enlace alterado, payload alterado, checkpoint inválido. | Verificación determinista y rechazo. |
| Restauración | Backup bajo actividad, recuperación y comparación con checkpoint independiente. | Integridad restaurada y límites de detección documentados. |
| Autonomía | Red y proveedor externo caídos. | Contabilidad local disponible y outbox pendiente. |
| Seguridad | Permisos, secretos en logs, diagnóstico en producción, clave no disponible. | Controles efectivos y fallos explícitos. |

Las pruebas de repositorios en memoria no sustituyen las de integración con SQLite. La suite debe ejercitar conexiones concurrentes reales y la configuración efectiva de producción cuando corresponda.

## 5. Comprobaciones de construcción

Desde el módulo Go real:

```powershell
go test ./...
go vet ./...
```

Añadir el detector de carreras cuando la plataforma lo soporte y los scripts reales de pruebas y build del frontend. Construir Wails desde su raíz configurada si el escritorio está incluido en el alcance.

Registrar comando, entorno, resultado y revisión del código. Una prueba no ejecutada se marca “no verificada”; no se interpreta como aprobada. Compilar no acredita partida doble, seguridad ni rollback.

## 6. Rendimiento

Utilizar PERF-01 de REQUIREMENTS. Medir transacciones completas, incluir auditoría y comprobar integridad después de la carga. Informar errores y percentiles, no solo el mejor throughput observado.

No declarar cumplido el objetivo si faltan datos para reproducirlo.

## 7. Clasificación de hallazgos

- **Crítico:** evidencia de pérdida o alteración financiera grave, extracción de secretos o acceso entre tenants de impacto inmediato.
- **Alto:** incumplimiento de una invariante, atomicidad, autorización o contrato esencial con riesgo demostrado.
- **Medio:** carencia que afecta operación, reproducibilidad o verificabilidad sin demostrar el impacto anterior.
- **Bajo:** claridad, formato o mantenibilidad sin efecto funcional directo comprobado.

La prioridad depende de evidencia y contexto. Una ambigüedad documental no demuestra una vulnerabilidad implementada.

Cada hallazgo incluye archivo y línea, requisito afectado, observación, consecuencia, evidencia reproducible, recomendación y estado. No presentar hipótesis como hechos.

## 8. Criterios de aceptación

- [ ] Revisión y entorno identificados.
- [ ] Invariantes financieras verificadas.
- [ ] Aislamiento y autorización verificados en todos los transportes.
- [ ] Rollback y concurrencia SQLite demostrados.
- [ ] Contratos numéricos consistentes entre frontend, backend y auditoría.
- [ ] Cadena, checkpoints y límites de confianza documentados y probados según alcance.
- [ ] Restauración y custodia de claves evaluadas antes de producción.
- [ ] Pruebas y construcción pertinentes aprobadas o exclusiones explicadas.
- [ ] Hallazgos críticos y altos resueltos o explícitamente pendientes; no se emite aprobación plena mientras subsistan.
- [ ] Documentación actualizada con el estado real, sin promesas no verificadas.

El informe final debe separar resultados comprobados, riesgos abiertos y exclusiones. Debe indicar si evalúa solo documentos, un componente, el backend completo o una release del escritorio.

## 9. Matriz de aceptación Kaspersky y firewall

Antes de probar, registrar SO/build, arquitectura, versión y hash de artefactos, editor, WebView2, producto/edición/versión de Kaspersky, actualización de firmas, componentes habilitados, políticas corporativas y firewalls efectivos. No instalar ni cambiar productos de seguridad automáticamente.

Ejecutar en entorno de prueba con datos sintéticos:

| Escenario | Resultado exigido | Requisitos |
| --- | --- | --- |
| Instalación y primer inicio con protecciones activas | Sin necesidad de desactivar módulos ni excluir directorios; incidencias registradas. | AV-01, NET-02 |
| Firma del instalador y ejecutable | Firma verificable y editor esperado; registro de fecha y hash. | AV-02 |
| Ciclo contable y backup/restauración | Persistencia y auditoría consistentes bajo protección activa. | AV-01 |
| Sin Internet y salidas bloqueadas | Core utilizable; integración opcional pendiente. | NET-03, OFF-01 |
| Inspección de red de desarrollo y release por separado | Ningún listener externo no documentado. | NET-01 |
| Actualización válida, manipulada e interrumpida | Validación de autenticidad, rechazo o recuperación según caso. | UPD-01 |
| Alerta real durante la prueba | Preservar evidencia y aplicar el flujo de revisión; no fabricar malware ni desactivar el antivirus. | AV-03 |

Para comprobar Authenticode, sustituir por la ruta real del artefacto:

```powershell
Get-AuthenticodeSignature -FilePath 'C:\ruta\ContableFix.exe'
Get-FileHash -LiteralPath 'C:\ruta\ContableFix.exe' -Algorithm SHA256
```

Un estado `Valid` es evidencia de validación de firma en ese entorno, no un dictamen de Kaspersky. Inspeccionar también el editor esperado y la política de confianza.

La inspección de red debe filtrar por PID del ejecutable y procesos relacionados comprobados. No atribuir al producto todos los puertos de Windows. Registrar listeners TCP/UDP y conexiones salientes con ruta del proceso; la comprobación debe hacerse durante los escenarios relevantes, no solo al iniciar.

Ejemplo de lectura TCP, reemplazando 1234 por un PID realmente identificado:

```powershell
$contableProcessIds = @(1234)
Get-NetTCPConnection | Where-Object { $_.OwningProcess -in $contableProcessIds }
Get-NetUDPEndpoint | Where-Object { $_.OwningProcess -in $contableProcessIds }
```

No se garantiza que estos comandos estén disponibles en todas las ediciones o sin permisos suficientes. Una lectura denegada se registra como no verificada, no como ausencia de sockets.

Cada caso se marca aprobado, fallido, no ejecutado o no aplicable con motivo. No declarar compatibilidad universal ni ejecución de estas pruebas por haber redactado este documento.