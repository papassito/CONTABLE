# CONTABLE FIX

## FCOS Kernel v2.2

**Contable Fix** es una plataforma contable profesional de escritorio construida alrededor de un **kernel contable transaccional en Go**.

El proyecto combina:

- Go
- FCOS Kernel
- Clean Architecture
- React
- TypeScript
- Tailwind CSS
- Wails v2
- PostgreSQL
- Git
- GitHub
- CI/CD

La arquitectura está diseñada para separar completamente las reglas contables de la interfaz, la persistencia y los mecanismos de transporte.

---

# 1. FCOS

**FCOS — Financial / Contable Operating System Kernel** es el núcleo contable de Contable Fix.

Su responsabilidad es proteger las reglas fundamentales de la operación financiera.

FCOS controla:

- cuentas
- cuentas auxiliares
- asientos
- líneas contables
- débitos
- créditos
- saldos
- libro mayor
- balance de comprobación
- períodos contables
- facturas
- impuestos
- reversiones
- estados contables
- integridad transaccional

La interfaz no contiene estas reglas.

Las consume.

---

# 2. Principio arquitectónico fundamental

> **El kernel gobierna la contabilidad.**

React no gobierna la contabilidad.

Wails no gobierna la contabilidad.

HTTP no gobierna la contabilidad.

PostgreSQL no gobierna la contabilidad.

GitHub tampoco gobierna la contabilidad.

Todos ellos son componentes periféricos alrededor de FCOS.

```text
                 CONTABLE FIX
                      │
                      ▼
                ┌───────────┐
                │    FCOS   │
                │  KERNEL   │
                └─────┬─────┘
                      │
       ┌──────────────┼──────────────┐
       ▼              ▼              ▼
   Persistence     Transport      Desktop
       │              │              │
 PostgreSQL          HTTP        Wails/React
```

---

# 3. Arquitectura

```text
┌──────────────────────────────────────────────┐
│              REACT / WAILS                   │
│                 UI                           │
└──────────────────────┬───────────────────────┘
                       │
                       ▼
┌──────────────────────────────────────────────┐
│              DTO / HANDLERS                  │
│              TRANSPORT                       │
└──────────────────────┬───────────────────────┘
                       │
                       ▼
┌──────────────────────────────────────────────┐
│                 SERVICES                     │
│                USE CASES                     │
└──────────────────────┬───────────────────────┘
                       │
                       ▼
┌──────────────────────────────────────────────┐
│                   FCOS                       │
│                  DOMAIN                      │
│                                              │
│ Account │ Journal │ Ledger │ Invoice         │
│ Tax     │ Period  │ Money  │ Company         │
└──────────────────────┬───────────────────────┘
                       │
                       ▼
┌──────────────────────────────────────────────┐
│                   PORTS                      │
│            REPOSITORY CONTRACTS              │
└──────────────────────┬───────────────────────┘
                       │
                ┌──────┴──────┐
                ▼             ▼
           IN-MEMORY      POSTGRESQL
```

---

# 4. Separación de responsabilidades

## Domain

Contiene:

- entidades
- value objects
- estados
- reglas
- invariantes
- errores de dominio

El dominio no conoce infraestructura.

---

## Services

Implementa los casos de uso.

Ejemplos:

```text
CreateAccount
CreateDraft
AddJournalLine
ValidateJournal
PostJournal
ReverseJournal
GetAccountLedger
GenerateTrialBalance
CreateInvoice
ClosePeriod
```

---

## Ports

Define contratos para servicios externos.

Ejemplos:

```text
AccountRepository
JournalRepository
LedgerRepository
InvoiceRepository
PeriodRepository
UnitOfWork
```

---

## Repository

Contiene las implementaciones de persistencia.

```text
repository/
├── memory/
└── postgres/
```

El kernel puede utilizar cualquiera de ellas sin modificar las reglas contables.

---

# 5. Persistencia In-Memory

El repositorio In-Memory es parte importante de la arquitectura y no solamente una herramienta temporal.

Permite:

- desarrollo rápido
- pruebas unitarias
- pruebas de integración
- pruebas del kernel
- validación de casos de uso
- ejecución sin PostgreSQL

El objetivo es poder demostrar:

```text
Crear cuenta
      ↓
Crear asiento
      ↓
Agregar líneas
      ↓
Validar
      ↓
Contabilizar
      ↓
Consultar mayor
      ↓
Generar Trial Balance
```

sin depender de una base de datos externa.

---

# 6. PostgreSQL

PostgreSQL proporciona persistencia real para los entornos que lo requieran.

Su incorporación no debe modificar las reglas del dominio.

La arquitectura buscada es:

```text
              FCOS
               │
              PORT
          ┌────┴────┐
          ▼         ▼
       MEMORY    POSTGRES
```

Ambas implementaciones deben respetar los mismos contratos.

---

# 7. Invariantes contables

FCOS debe proteger como mínimo las siguientes invariantes.

## Partida doble

```text
TOTAL DEBE = TOTAL HABER
```

Un asiento desequilibrado no puede contabilizarse.

---

## Inmutabilidad

Una operación contabilizada no puede modificarse arbitrariamente.

```text
DRAFT
  ↓
VALIDATED
  ↓
POSTED
  ↓
IMMUTABLE
```

Las correcciones se realizan mediante operaciones contables válidas, como la reversión.

---

## Cuentas auxiliares

Una cuenta que no acepta movimientos no puede recibir movimientos directamente.

---

## Precisión monetaria

No se utiliza `float64` para representar dinero.

Los importes utilizan representación exacta.

---

## Períodos

Un período cerrado no debe aceptar operaciones que alteren su información contable.

---

## Atomicidad

Las operaciones contables deben ser atómicas.

No se debe permitir que una operación falle dejando el sistema en un estado financiero parcialmente aplicado.

---

# 8. Ledger

El libro mayor representa la evolución de las cuentas.

Conceptualmente:

```text
SALDO INICIAL
      +
MOVIMIENTOS
      =
SALDO ACUMULADO
```

Debe ser posible consultar los movimientos de una cuenta dentro de un período y reconstruir su saldo de manera consistente.

---

# 9. Trial Balance

El balance de comprobación permite verificar la integridad global de la contabilidad.

Debe proporcionar información suficiente para determinar:

```text
Total Débitos
Total Créditos
Saldos Deudores
Saldos Acreedores
```

La condición fundamental es:

```text
Σ DÉBITOS = Σ CRÉDITOS
```

Esta condición forma parte de las pruebas automáticas del sistema.

---

# 10. Desktop

La aplicación de escritorio utiliza:

```text
Wails v2
   ↓
React
   ↓
TypeScript
   ↓
Tailwind CSS
```

La UI proporciona la experiencia de usuario.

FCOS proporciona la autoridad contable.

---

# 11. Módulos principales de la interfaz

La aplicación contempla, entre otros:

```text
Dashboard
Accounts
Journal
Ledger
Trial Balance
Invoices
Periods
Audit
Settings
```

Cada módulo debe consumir casos de uso reales del backend.

No se implementará lógica contable crítica directamente en componentes React.

---

# 12. API

La capa HTTP utiliza:

```text
Handlers
DTO
Mapper
Validation
Response
```

Los handlers son responsables de:

```text
Request
   ↓
Decode
   ↓
Validate
   ↓
Service
   ↓
Response
```

No contienen reglas contables.

---

# 13. Testing

La compilación no constituye por sí sola una prueba de corrección.

El proyecto debe validar:

```text
go test ./...
```

y adicionalmente comprobar:

- partida doble
- inmutabilidad
- reversión
- cuentas auxiliares
- precisión monetaria
- saldos
- ledger
- trial balance
- períodos
- atomicidad
- repositorios
- integración

---

# 14. Git y GitHub

GitHub forma parte del flujo oficial de ingeniería de Contable Fix.

El repositorio debe utilizar:

```text
Git
  ↓
Branch
  ↓
Commit
  ↓
Pull Request
  ↓
CI
  ↓
Review
  ↓
Merge
  ↓
Release
```

GitHub no sustituye las pruebas locales.

Las pruebas locales y las pruebas automatizadas deben complementarse.

---

# 15. CI/CD

El pipeline debe actuar como una barrera de calidad.

Conceptualmente:

```text
PUSH
  ↓
CHECKOUT
  ↓
DEPENDENCIES
  ↓
FORMAT
  ↓
VET
  ↓
UNIT TESTS
  ↓
ACCOUNTING TESTS
  ↓
INTEGRATION TESTS
  ↓
FRONTEND TESTS
  ↓
BUILD
  ↓
ARTIFACT
```

Un cambio que rompa una invariante contable no debe llegar a release.

---

# 16. GitHub Workflows

La estructura prevista es:

```text
.github/
│
├── workflows/
│   ├── ci.yml
│   ├── test.yml
│   ├── build.yml
│   └── release.yml
│
├── ISSUE_TEMPLATE/
│   ├── bug_report.md
│   └── feature_request.md
│
├── PULL_REQUEST_TEMPLATE.md
│
└── dependabot.yml
```

Los workflows deben mantenerse simples, reproducibles y auditables.

No deben depender de rutas locales específicas de una máquina de desarrollo.

---

# 17. Releases

Las versiones deben ser identificables y reproducibles.

Una release debe incluir:

- versión
- changelog
- artefactos
- checksums
- información de build
- notas de versión

La generación de releases debe realizarse mediante CI/CD cuando sea apropiado.

---

# 18. Seguridad

El proyecto debe considerar seguridad desde el diseño.

Entre otros aspectos:

- validación de entradas
- control de errores
- protección de credenciales
- configuración mediante variables de entorno
- dependencias controladas
- revisión de cambios
- auditoría
- trazabilidad
- mínima exposición de servicios

Nunca deben almacenarse secretos directamente en el repositorio.

---

# 19. Estructura principal

```text
CONTABLE-FIX/
│
├── .github/
├── cmd/
├── internal/
│   ├── domain/
│   ├── ports/
│   ├── service/
│   ├── repository/
│   │   ├── memory/
│   │   └── postgres/
│   ├── handler/
│   ├── dto/
│   ├── mapper/
│   ├── validation/
│   ├── config/
│   └── bootstrap/
│
├── pkg/
├── migrations/
├── tests/
├── frontend/
├── scripts/
├── docs/
├── build/
│
├── go.mod
├── go.sum
├── package.json
├── pnpm-lock.yaml
├── wails.json
├── README.md
├── CHANGELOG.md
├── SECURITY.md
├── CONTRIBUTING.md
└── LICENSE
```

---

# 20. Reglas de desarrollo

### Regla 1

No colocar lógica contable en React.

### Regla 2

No colocar lógica contable en handlers.

### Regla 3

No permitir que el dominio dependa de PostgreSQL.

### Regla 4

No utilizar `float64` para dinero.

### Regla 5

No modificar operaciones contabilizadas de manera destructiva.

### Regla 6

Toda nueva operación contable debe tener pruebas.

### Regla 7

Los contratos de los Ports deben permanecer estables.

### Regla 8

Los repositorios son reemplazables.

### Regla 9

Un cambio que rompa una invariante contable debe bloquear el merge.

### Regla 10

El código debe poder probarse sin PostgreSQL utilizando los repositorios In-Memory.

---

# 21. Criterio de terminado

Contable Fix no se considera terminado simplemente porque:

```text
go build ./...
```

sea exitoso.

El sistema debe demostrar:

```text
CUENTA
   ↓
ASIENTO
   ↓
LÍNEAS
   ↓
VALIDACIÓN
   ↓
PARTIDA DOBLE
   ↓
CONTABILIZACIÓN
   ↓
LEDGER
   ↓
SALDOS
   ↓
TRIAL BALANCE
   ↓
AUDITORÍA
```

manteniendo las invariantes del kernel.

---

# 22. Definition of Done

Una funcionalidad contable está terminada cuando:

- el dominio está definido
- las reglas están implementadas
- los Ports están definidos
- existe implementación In-Memory
- existe prueba automatizada
- el caso de uso funciona
- el Ledger permanece consistente
- Trial Balance permanece consistente
- la API está integrada cuando corresponde
- la UI está integrada cuando corresponde
- CI pasa
- no existen regresiones
- la documentación está actualizada

---

# 23. Filosofía

Contable Fix sigue una filosofía simple:

> **Correcto antes que bonito.**
>
> **Determinista antes que conveniente.**
>
> **Auditable antes que opaco.**
>
> **El kernel primero.**
>
> **La interfaz después.**

La tecnología puede cambiar.

La persistencia puede cambiar.

La interfaz puede cambiar.

El sistema operativo puede cambiar.

El mecanismo de distribución puede cambiar.

Pero las reglas fundamentales de FCOS deben permanecer protegidas.

---

# CONTABLE FIX

## FCOS Kernel v2.2

```text
DOMAIN FIRST
DOUBLE ENTRY
TRANSACTION SAFE
IMMUTABLE ACCOUNTING
AUDITABLE BY DESIGN
TESTED CONTINUOUSLY
RELEASED REPRODUCIBLY
```

**Contable Fix**
