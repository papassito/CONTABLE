# FCOS v2.2 — Arquitectura y contratos

## 1. Estado y decisiones

Esta es la especificación objetivo, no una descripción verificada de la implementación. Sus decisiones vigentes son SQLite WAL local, importes de MXN en centavos enteros, Go en el dominio y Wails v2 como transporte desktop. PostgreSQL no forma parte de la base exigida.

Los contratos IPC y HTTP se incluyen aquí; no se necesita un `CONTRACTS.md` separado.

## 2. Capas y dependencias

```text
React / TypeScript
        |
Wails IPC o HTTP opcional
        |
Adaptadores de entrada: identidad, DTO, validación sintáctica
        |
Casos de uso: autorización, reglas, coordinación transaccional
        |
Dominio + puertos
        |
Adaptadores: SQLite, memoria de pruebas, reloj, identificadores
```

El dominio no importa Wails, HTTP ni drivers SQL. Los casos de uso dependen de interfaces de repositorio y Unit of Work. Los adaptadores de entrada traducen datos y errores, sin convertirse en autoridad de reglas contables.

Los repositorios en memoria sirven para pruebas de dominio y casos de uso. No prueban por sí solos restricciones, aislamiento ni transacciones del adaptador SQLite.

## 3. Autonomía local

Crear borradores, contabilizar, revertir, consultar cuentas y mayor y generar auditoría local no exige conexión a SAT, servicios de anclaje u otros proveedores externos.

Los eventos destinados a integraciones se insertan en una outbox en la misma transacción del cambio que los origina. Un worker los entrega después del commit. Debe soportar reintentos y consumidores idempotentes: no se promete entrega exactamente una vez.

No se realizan llamadas de red dentro de la transacción contable. Un fallo de la auditoría local sí impide confirmar el cambio; un fallo de entrega externa posterior no lo revierte.

## 4. Dinero y límites

- Moneda base: MXN. Una unidad monetaria equivale a 100 centavos.
- Dominio y persistencia: enteros con signo de 64 bits, con operaciones aritméticas comprobadas.
- Líneas: débito y crédito no negativos; exactamente uno debe ser positivo.
- Saldos: pueden ser negativos. Convención de proyección: suma de débitos menos suma de créditos.
- Se rechaza cualquier suma, resta o inversión que no pueda representarse en `int64`.
- No se utilizan `float32` ni `float64` para calcular importes contables.

**Contrato de transporte de esta edición:** todos los campos monetarios cruzan IPC/JSON como cadenas decimales de centavos, tanto de entrada como de salida. Ejemplo: `"12345"` representa MXN 123.45. El frontend no debe convertirlos a `number` para operar o acumularlos.

Se acepta una representación canónica sin espacios, exponentes, separadores, signo positivo ni ceros iniciales: `0`, un entero positivo, o un entero negativo para campos que permitan saldos negativos. Se rechaza `-0`. El backend valida sintaxis y rango antes de construir el valor de dominio.

Esta decisión sustituye los DTO anteriores que usaban `number`; su adopción en código requiere migrar conjuntamente productores, consumidores y pruebas. No se afirma compatibilidad retroactiva.

Los cálculos fiscales con precisión inferior al centavo requieren una especificación propia de escala y redondeo antes de implementarse. Solo sus resultados cuantizados ingresan al núcleo; no se inventa una política fiscal en esta base.

## 5. Asientos, estados y períodos

```text
BORRADOR -> CONTABILIZADO
```

La validación es una operación, no un estado persistido. La inmutabilidad es una propiedad de CONTABILIZADO. No existe estado persistido ANULADO.

Se permite guardar borradores incompletos o desequilibrados para edición, siempre que los campos presentes sean sintácticamente válidos y sus referencias estén autorizadas. No afectan mayor ni saldos.

Para contabilizar se exige simultáneamente:

1. Al menos dos líneas.
2. Cada línea tiene exactamente un importe positivo y el otro en cero.
3. Total de débitos igual a total de créditos y mayor que cero, sin desbordamiento.
4. Todas las cuentas existen, están activas, pertenecen al tenant y aceptan movimientos.
5. La fecha contable pertenece a un período abierto del tenant.
6. El asiento todavía es BORRADOR y el actor está autorizado.

Las fechas contables usan `YYYY-MM-DD`; los timestamps de eventos usan RFC 3339 en UTC. No se convierte una fecha contable entre husos horarios.

La reversión crea un nuevo asiento CONTABILIZADO, invierte las líneas, conserva la moneda y el tenant, exige motivo y fecha dentro de un período abierto y referencia al original. No cambia ninguna columna del original. Solo se admite una reversión directa por original. La UI puede mostrar “revertido” como propiedad derivada de esa relación.

## 6. Persistencia y Unit of Work

SQLite utiliza WAL, `busy_timeout` configurado y claves foráneas habilitadas en cada conexión. El driver objetivo del backend es `modernc.org/sqlite`; la versión efectiva se fija en `go.mod`.

El adaptador Unit of Work debe adquirir la transacción escritora antes de leer el estado que va a modificar, mediante `BEGIN IMMEDIATE` o una configuración equivalente documentada del driver. No debe emitir un BEGIN dentro de otra transacción ya iniciada.

Una contabilización confirma conjuntamente:

1. Revalidación del asiento, cuentas, autorización y período.
2. Inserción de movimientos en `ledger_entries`.
3. Actualización de las proyecciones de saldo.
4. Cambio de BORRADOR a CONTABILIZADO.
5. Inserción del evento en la cadena local de auditoría.
6. Inserción de eventos outbox, cuando corresponda.

Todo fallo produce rollback completo. Cerrar un período utiliza la misma disciplina de escritura, evitando carreras con contabilizaciones. Los reintentos por bloqueo son acotados y no deben repetir efectos ya confirmados.

`ledger_entries` contiene solo movimientos contabilizados. Las líneas editables del borrador se almacenan separadamente en `journal_lines`. Al contabilizar, se conserva la relación entre línea original y movimiento y se impide duplicarla.

El adaptador debe proteger asientos y líneas contabilizados y movimientos del mayor frente a UPDATE/DELETE en el flujo normal de aplicación. Las migraciones o reparaciones excepcionales requieren procedimientos explícitos y trazables.

## 7. Identidad y contratos de casos de uso

Cada caso de uso recibe un contexto con identidad autenticada, tenant autorizado y permisos comprobables. El frontend no decide la identidad por enviar un `tenant_id` o un header.

Firmas conceptuales del servicio, independientes del mecanismo de bindings:

```go
type JournalService interface {
    CreateDraft(ctx context.Context, input CreateDraftInput) (JournalEntry, error)
    UpdateDraft(ctx context.Context, id string, input UpdateDraftInput) (JournalEntry, error)
    PostEntry(ctx context.Context, id string) (JournalEntry, error)
    ReverseEntry(ctx context.Context, id string, input ReverseEntryInput) (JournalEntry, error)
    GetEntry(ctx context.Context, id string) (JournalEntry, error)
}
```

Los adaptadores Wails construyen el contexto desde una sesión validada; no se exige que JavaScript suministre un `context.Context`. Los tipos Go anteriores son nombres de contrato, no afirmaciones de símbolos existentes.

DTO mínimo de entrada para un borrador:

```typescript
type Cents = string;
type EntryStatus = 'BORRADOR' | 'CONTABILIZADO';

interface JournalLineInput {
  account_id: string;
  description: string;
  debit_cents: Cents;
  credit_cents: Cents;
  third_party_id?: string;
}

interface CreateDraftInput {
  date: string; // YYYY-MM-DD
  concept: string;
  reference?: string;
  lines: JournalLineInput[];
}

interface ReverseEntryInput {
  date: string; // Fecha dentro de un período abierto.
  reason: string;
}
```

ID, número, estado, tenant y referencia de reversión los asigna o controla el backend. La respuesta incluye esos campos y los importes en el mismo formato de cadena. Los límites de longitud, paginación y tamaño de solicitudes deben fijarse antes de publicar cada endpoint.

Repetir PostEntry sobre un asiento ya contabilizado devuelve un conflicto sin agregar movimientos. Repetir una reversión directa también devuelve conflicto. El transporte debe permitir consultar el resultado tras una respuesta perdida.

## 8. HTTP opcional y errores

Superficie objetivo bajo `/api/v1`:

| Método | Ruta | Operación |
| --- | --- | --- |
| POST / GET | `/accounts` | Crear / listar cuentas autorizadas. |
| POST | `/journal-entries` | Crear borrador. |
| GET | `/journal-entries/{id}` | Consultar asiento. |
| PATCH | `/journal-entries/{id}` | Editar exclusivamente borrador. |
| POST | `/journal-entries/{id}/post` | Contabilizar. |
| POST | `/journal-entries/{id}/reverse` | Crear reversión con fecha y motivo. |

Toda ruta aplica identidad y autorización. Un header puede seleccionar un tenant entre los autorizados, nunca acreditarlo por sí mismo. El diagnóstico se mantiene fuera de `/api/v1`, autenticado, local y deshabilitado en producción.

Errores de contrato: `INVALID_INPUT`, `AMOUNT_OUT_OF_RANGE`, `UNBALANCED_JOURNAL`, `PERIOD_CLOSED`, `ACCOUNT_NOT_POSTABLE`, `CONFLICT`, `UNAUTHENTICATED`, `FORBIDDEN`, `NOT_FOUND`, `STORAGE_BUSY` y `AUDIT_INTEGRITY_FAILURE`.

HTTP utiliza las categorías 400, 401, 403, 404, 409, 422, 503 o 500 según la causa; esta lista no establece una correspondencia posicional con los errores anteriores. Debe fijarse una tabla exhaustiva por endpoint antes de su publicación. No se exponen SQL, secretos ni existencia de objetos de otros tenants.
