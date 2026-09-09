# Contable Fix by KLIK

Contable Fix es un proyecto de aplicación contable de escritorio basado en un núcleo financiero escrito en Go, con interfaz React y TypeScript integrada mediante Wails v2.

Este documento es la entrada general del repositorio. Los seis documentos normativos se encuentran en la raíz del proyecto.

## Estado de esta documentación

Esta edición unifica la documentación de FCOS v2.2. Describe la arquitectura objetivo y sus criterios de aceptación; no certifica que las funciones estén implementadas, que las pruebas pasen o que exista una versión lista para producción.

La decisión adoptada para esta edición es **persistencia local SQLite en modo WAL**. Las referencias anteriores a PostgreSQL describían una alternativa que no forma parte de la base vigente. Incorporarla requiere una decisión arquitectónica explícita y pruebas de conformidad; no debe implementarse por inferencia a partir de documentación antigua.

FCOS significa **Financial Core Operating System**. El nombre del producto es **Contable Fix by KLIK**.

## Documentos y autoridad

| Documento | Responsabilidad |
| --- | --- |
| [Guía técnica](#guía-técnica) | Preparación del entorno, ejecución y navegación. |
| [Arquitectura](ARCHITECTURE.md) | Límites de componentes, persistencia, transacciones y contratos IPC/HTTP. |
| [Requisitos](REQUIREMENTS.md) | Comportamientos exigidos y criterios de aceptación. |
| [Mapa](MAP.md) | Organización lógica, relaciones de datos y estructura objetivo. |
| [Seguridad](SECURITY.md) | Identidad, aislamiento, secretos y evidencia criptográfica. |
| [Guía de auditoría](AUDIT_GUIDE.md) | Procedimiento y evidencias necesarias para evaluar una implementación. |

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


# Guía técnica

FCOS es el núcleo de Contable Fix by KLIK. La base documental vigente utiliza Go, SQLite WAL y una interfaz React/TypeScript mediante Wails v2.

## Alcance y estado

Esta guía describe el entorno y la estructura objetivo. No se ha validado contra el código ni se han ejecutado sus comandos durante la revisión documental. Las rutas deben comprobarse en la copia de trabajo antes de utilizarlas.

- Repositorio declarado: `https://github.com/papassito/CONTABLE.git`.
- Identificador de módulo declarado: `github.com/klik/fcos-kernel`.

El identificador de un módulo Go no tiene que coincidir con la URL del repositorio. El valor efectivo debe consultarse en `go.mod`; no debe cambiarse solo para hacer coincidir ambos nombres.

## Preparación del entorno

Instalar Git, Go, Node.js con npm y Wails CLI v2. En Windows, preparar también los requisitos de construcción y ejecución nativa exigidos por la versión de Wails utilizada.

Las versiones efectivas deben obtenerse de `go.mod`, `package.json`, el archivo de bloqueo y la configuración de CI. Si no existe una versión de herramienta fijada, debe acordarse y registrarse antes de declarar reproducible la instalación. Esta guía no impone versiones antiguas de Node ni presupone compatibilidad con cualquier versión superior.

Comprobar las herramientas desde PowerShell:

```powershell
git --version
go version
node --version
npm --version
wails doctor
```

El objetivo de evitar Cgo corresponde al driver SQLite del backend. No constituye una garantía de que el escritorio, sus herramientas o todas las modalidades de prueba carezcan de requisitos nativos.

## Obtener el proyecto

Si ya existe una copia local, no es necesario clonar de nuevo. Para una copia nueva:

```powershell
git clone https://github.com/papassito/CONTABLE.git
Set-Location .\CONTABLE
```

Identificar dos rutas reales: la carpeta desktop que contiene `wails.json` y el módulo backend que contiene `go.mod`. En la carpeta inspeccionada, `go.mod` y `wails.json` están en la raíz. La reconstrucción conserva esa ubicación.

En los comandos siguientes, reemplazar las rutas de ejemplo:

```powershell
$desktopRoot = 'C:\ruta\CONTABLE'
$backendRoot = $desktopRoot
```

## Instalar dependencias

Desde la carpeta que realmente contiene `package.json`, usar el gestor correspondiente al archivo de bloqueo del proyecto. Solo cuando npm haya sido seleccionado como gestor y su bloqueo validado:

```powershell
Set-Location $desktopRoot
npm ci
```

Si el frontend vive en una subcarpeta, ejecutar ese paso allí. El inventario actual contiene bloqueos de npm, pnpm y Bun; la configuración Wails apunta a pnpm. La elección y sincronización del gestor están pendientes: no eliminar bloqueos ni instalar dependencias como parte de esta entrega documental. No crear un archivo de bloqueo de otro gestor. Si no hay archivo de bloqueo, resolver y registrar esa carencia antes de declarar la instalación reproducible.

Para descargar las dependencias Go declaradas:

```powershell
Set-Location $backendRoot
go mod download
```

`go mod tidy` es una operación de mantenimiento que puede modificar manifiestos; no es un requisito automático de instalación.

## Ejecutar el escritorio

Desde la carpeta que contiene `wails.json`:

```powershell
Set-Location $desktopRoot
wails dev
```

## Ejecutar la API opcional

Solo si la copia contiene `cmd/api` y su configuración está documentada:

```powershell
Set-Location $backendRoot
go run ./cmd/api
```

La dirección, el puerto, la autenticación y el archivo SQLite deben obtenerse de la configuración efectiva. No se presupone el puerto 8080. La API debe quedar deshabilitada o ligada a loopback por defecto y aplicar los controles de SECURITY.

## Comprobar y construir

```powershell
Set-Location $backendRoot
go test ./...
go vet ./...
```

Ejecutar `go test -race ./...` cuando la plataforma y la cadena de herramientas lo soporten. Registrar expresamente si no se ejecuta; no sustituye las pruebas de concurrencia con SQLite.

Las pruebas y la construcción del frontend deben utilizar los scripts existentes en `package.json`. No se presupone un script que no esté declarado. Los comandos de esta guía son instrucciones para una fase posterior; no fueron ejecutados durante esta verificación documental.

Para el escritorio:

```powershell
Set-Location $desktopRoot
wails build
```

Consultar la configuración y la salida de Wails para localizar el artefacto real. No se garantiza un nombre de ejecutable desde esta guía.

## Documentación relacionada

- [ARCHITECTURE.md](ARCHITECTURE.md): incluye los contratos de datos y operaciones.
- [REQUIREMENTS.md](REQUIREMENTS.md): requisitos y aceptación.
- [MAP.md](MAP.md): mapa lógico y estructura objetivo.
- [SECURITY.md](SECURITY.md): controles de seguridad.
- [AUDIT_GUIDE.md](AUDIT_GUIDE.md): procedimiento de evaluación.

Los scripts auxiliares de integridad solo deben ejecutarse después de verificar que existen y revisar su alcance. Un inventario de archivos no prueba la corrección contable.

## Orden de reconstrucción y decisiones pendientes

Leer los seis documentos antes de editar código. Implementar primero dominio y pruebas, luego SQLite/UoW, casos de uso, seguridad definida y finalmente interfaz. Reutilizar el código conforme a los contratos; la reorganización no autoriza borrar implementaciones.

Los checkpoints firmados, el mecanismo concreto de autenticación y recuperación de claves, las reglas fiscales y el gestor definitivo del frontend requieren decisiones adicionales. Se puede avanzar con componentes independientes; no se declarará producción aprobada mientras falten los controles exigidos. No se han aprobado BIP-39, firma por asiento ni una base separada por empresa.

El protocolo de convivencia con Kaspersky y firewalls está en SECURITY, con requisitos AV y NET en REQUIREMENTS y evidencia de aceptación en AUDIT_GUIDE. No se garantiza certificación por antivirus.
## Inicio autorizado de la implementación

La base documental está lista para comenzar el núcleo; las decisiones de funciones futuras no bloquean esta etapa. La primera entrega debe implementar o adaptar dinero exacto, cuentas y jerarquía, borradores, validación de posteo y contratos de repositorio, con pruebas unitarias. Después se aborda SQLite y la atomicidad de posteo/reversión, mayor, períodos y auditoría local.

Usar un contexto de actor y tenant inyectable en pruebas; no convertirlo en un acceso de producción sin autenticación. No abrir API externa, implementar recuperación de claves ni firmar checkpoints hasta cerrar sus contratos específicos. Si una decisión pendiente afecta una tarea, continuar con las tareas independientes y reportar solo el bloqueo concreto.

Conservar los seis documentos como fuente principal. No añadir versiones abreviadas, reconstruir documentación desde scripts antiguos ni eliminar código porque no coincida con nombres de carpetas propuestos. Registrar para cada entrega archivos cambiados, pruebas realmente ejecutadas y requisitos cubiertos.