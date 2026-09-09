# Contable Fix by KLIK

Contable Fix es un proyecto de aplicación contable de escritorio basado en un núcleo financiero escrito en Go, con interfaz React y TypeScript integrada mediante Wails v2.

Este documento es la entrada general del repositorio. La especificación técnica se encuentra en `Contable GO/`.

## Estado de esta documentación

Esta edición unifica la documentación de FCOS v2.2. Describe la arquitectura objetivo y sus criterios de aceptación; no certifica que las funciones estén implementadas, que las pruebas pasen o que exista una versión lista para producción.

La decisión adoptada para esta edición es **persistencia local SQLite en modo WAL**. Las referencias anteriores a PostgreSQL describían una alternativa que no forma parte de la base vigente. Incorporarla requiere una decisión arquitectónica explícita y pruebas de conformidad; no debe implementarse por inferencia a partir de documentación antigua.

FCOS significa **Financial Core Operating System**. El nombre del producto es **Contable Fix by KLIK**.

## Documentos y autoridad

| Documento | Responsabilidad |
| --- | --- |
| [Guía técnica](Contable%20GO/README.md) | Preparación del entorno, ejecución y navegación. |
| [Arquitectura](Contable%20GO/ARCHITECTURE.md) | Límites de componentes, persistencia, transacciones y contratos IPC/HTTP. |
| [Requisitos](Contable%20GO/REQUIREMENTS.md) | Comportamientos exigidos y criterios de aceptación. |
| [Mapa](Contable%20GO/MAP.md) | Organización lógica, relaciones de datos y estructura objetivo. |
| [Seguridad](Contable%20GO/SECURITY.md) | Identidad, aislamiento, secretos y evidencia criptográfica. |
| [Guía de auditoría](Contable%20GO/AUDIT_GUIDE.md) | Procedimiento y evidencias necesarias para evaluar una implementación. |

REQUIREMENTS define qué debe cumplirse; ARCHITECTURE define cómo se organiza la solución; SECURITY desarrolla sus controles de seguridad. README y MAP no sustituyen esos contratos. Una contradicción requiere corregir los documentos implicados antes de implementar el comportamiento afectado.

Los contratos se mantienen dentro de ARCHITECTURE. Esta edición no requiere un archivo `CONTRACTS.md` adicional.

## Principios del núcleo

- La lógica contable reside en el backend. La interfaz y los transportes invocan casos de uso.
- El dinero contabilizado se representa como centavos enteros de MXN con `int64` en Go.
- Solo se contabilizan asientos equilibrados, no nulos y con al menos dos líneas válidas.
- Los estados persistidos son `BORRADOR` y `CONTABILIZADO`.
- Un asiento contabilizado se corrige mediante otro asiento de reversión; el original permanece intacto.
- El mayor es la fuente de verdad de los movimientos contabilizados. Los saldos almacenados son proyecciones reconstruibles.
- El cambio de estado, mayor, saldos, auditoría y mensajes de salida asociados se confirman o revierten juntos.
- Cada operación se ejecuta dentro de una identidad autenticada y un tenant autorizado.
- La indisponibilidad de servicios externos no bloquea las operaciones contables locales.

## Alcance

La base comprende cuentas, borradores, contabilización, reversión, mayor, balance de comprobación y control de períodos. Facturación, cálculo fiscal, conciliación y automatización son contextos previstos: su presencia en el mapa no demuestra implementación ni cumplimiento normativo.

La base documental no define contabilidad multimoneda. No deben mezclarse monedas ni asumir conversiones implícitas.

## Desarrollo y entrega

Antes de ejecutar comandos, consultar la guía técnica y confirmar la ubicación real de los módulos en la copia de trabajo.

Un cambio contable debe incluir evidencia de sus invariantes, fallos parciales, aislamiento y compatibilidad con los contratos. Compilar no equivale a aprobar una funcionalidad. Una release requiere una revisión identificable, pruebas relevantes aprobadas, versión, notas de cambios, artefactos e información verificable de construcción.

No se deben almacenar credenciales ni material privado en el repositorio.
