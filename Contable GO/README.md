# FCOS v2.2 Kernel — Contable Fix by KLIK

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![Framework](https://img.shields.io/badge/Wails-v2-1E1E1E?style=flat&logo=wails)
![Frontend](https://img.shields.io/badge/React-18.x-61DAFB?style=flat&logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?style=flat&logo=typescript)
![Database](https://img.shields.io/badge/SQLite-WAL--Mode-003B57?style=flat&logo=sqlite)
![Architecture](https://img.shields.io/badge/Architecture-DDD%20%2F%20Clean%20Architecture-purple)

**FCOS (Financial Core Operating System) v2.2 Kernel** es el motor contable modular de alto rendimiento para **Contable Fix by KLIK**. Combina procesamiento inmutable de libros diarios, garantía estricta de partida doble en céntimos enteros, arquitectura multi-inquilino aislada, telemetría sintética *Day-2* y una cadena de auditoría criptográfica local (Tamper-Evident).

---

## ⚠️ Identificación de Módulos y Repositorio

* **Repositorio de Código (Git):** `https://github.com/papassito/CONTABLE.git`
* **Nombre de Módulo Go (`go.mod`):** `github.com/klik/fcos-kernel`

---

## 🚀 Guía Rápida de Inicio y Configuración

### Requisitos Previos del Sistema
* **Go (Golang):** Versión 1.22.0 o superior.
* **Node.js:** Versión LTS 18.x o 20.x.
* **Wails CLI v2:** Herramienta de construcción nativa.
* **PowerShell:** Versión 7.x o Windows PowerShell integrado.
* **Git:** Para control de versiones.

### Instalación Paso a Paso

```powershell
git clone [https://github.com/papassito/CONTABLE.git](https://github.com/papassito/CONTABLE.git)
cd CONTABLE

npm install

cd go-skeleton
go mod tidy
cd ..

Ejecución en Entorno de Desarrollo (Hot-Reload)
PowerShell

wails dev

Ejecución Standalone del Kernel API (Modo REST HTTP)
PowerShell

cd go-skeleton
go run cmd/api/main.go

El servidor HTTP estará disponible y escuchando en http://localhost:8080.
Compilación de Ejecutable para Producción
PowerShell

wails build

El ejecutable de producción se generará en la ruta: build/bin/ContableFixByKlik.exe.
🧪 Automatización de Pruebas y Diagnóstico Forense
PowerShell

cd go-skeleton
go test -v -race ./...

cd ..
.\verificar_integridad.ps1

📂 Visión General del Repositorio
Plaintext

CONTABLE/
├── app.go                      # Enlaces IPC (Bindings) exportados entre React y Wails
├── main.go                     # Punto de entrada base de la aplicación Desktop
├── wails.json                  # Archivo de configuración global de compilación Wails
├── README.md                   # Documentación principal
├── ARCHITECTURE.md             # Visión, principios y patrones arquitectónicos
├── CONTRACTS.md                # Contratos DTO, firmas IPC e interfaces REST/Debug
├── REQUIREMENTS.md             # Especificación técnica de requerimientos (RF/RNF/PERF)
├── MAP.md                      # Mapeo físico del sistema y registro de Bounded Contexts
├── SECURITY.md                 # Modelo de seguridad, criptografía y hash chain
├── verificar_integridad.ps1    # Script PowerShell de auditoría de integridad del filesystem
├── go-skeleton/                # Módulo Go Backend ([github.com/klik/fcos-kernel](https://github.com/klik/fcos-kernel))
│   ├── cmd/api/main.go         # Entrypoint Standalone REST
│   ├── internal/               # 13 Bounded Contexts bajo Clean Architecture / DDD
│   └── pkg/database/           # Unit of Work y administración de SQLite WAL
└── src/                        # Aplicación Frontend React 18 + TypeScript
    └── services/
        └── wailsService.ts     # Puente de comunicación IPC tipado con el kernel