# FCOS v2.2 — Guía técnica

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

Identificar dos rutas reales: la carpeta desktop que contiene `wails.json` y el módulo backend que contiene `go.mod`. El mapa propone `go-skeleton/` para el backend, pero debe verificarse su presencia.

En los comandos siguientes, reemplazar las rutas de ejemplo:

```powershell
$desktopRoot = 'C:\ruta\CONTABLE'
$backendRoot = Join-Path $desktopRoot 'go-skeleton'
```

## Instalar dependencias

Desde la carpeta que realmente contiene `package.json`, usar el gestor correspondiente al archivo de bloqueo del proyecto. Si existe `package-lock.json`:

```powershell
Set-Location $desktopRoot
npm ci
```

Si el frontend vive en una subcarpeta, ejecutar ese paso allí. No crear un archivo de bloqueo de otro gestor. Si no hay archivo de bloqueo, resolver y registrar esa carencia antes de declarar la instalación reproducible.

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

Las pruebas y la construcción del frontend deben utilizar los scripts existentes en `package.json`. No se presupone un script que no esté declarado.

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
