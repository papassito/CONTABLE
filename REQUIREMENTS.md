# FCOS v2.2 — Requisitos y aceptación

## Estado

Los siguientes requisitos son normativos para la base documental. Su estado de implementación es **no verificado**. Aprobarlos requiere evidencias identificadas en AUDIT_GUIDE; no basta con que aparezcan en un README.

## Requisitos funcionales

| ID | Requisito | Criterio de aceptación |
| --- | --- | --- |
| RF-01 | Cada operación requiere identidad y tenant autorizado. | Una identidad sin acceso no puede leer, escribir ni inferir objetos de otro tenant. |
| RF-02 | Códigos de cuenta únicos por tenant; jerarquía sin ciclos ni padres de otro tenant. | Rechazar duplicados locales, ciclos y relaciones cruzadas; admitir igual código en tenants distintos. |
| RF-03 | Solo cuentas activas, auxiliares y sin hijos aceptan movimientos. | Rechazar contabilización en cuentas inactivas o agrupadoras; impedir crear hijos bajo una cuenta afectable sin una transición válida. |
| RF-04 | No desactivar cuentas con saldo distinto de cero. | La validación y el cambio se ejecutan sin carreras con contabilizaciones. |
| RF-05 | Borradores editables sin efecto contable. | Guardar un borrador incompleto no altera mayor ni saldo; sus referencias presentes pertenecen al tenant. |
| RF-06 | Contabilización equilibrada, positiva y con al menos dos líneas. | Rechazar asiento nulo, una sola línea, importes negativos, línea nula, doble importe positivo o desbalance de un centavo. |
| RF-07 | Solo estados BORRADOR y CONTABILIZADO. | Ningún flujo normal modifica o elimina un asiento contabilizado ni sus líneas. |
| RF-08 | Reversión mediante nuevo asiento. | Exigir motivo, fecha válida y mismo tenant; invertir importes; impedir segunda reversión directa y conservar intacto el original. |
| RF-09 | Períodos controlados por tenant. | Rechazar contabilización o reversión con fecha en período cerrado o inexistente; probar carrera entre cierre y posteo. |
| RF-10 | Mayor como fuente de verdad y saldo como proyección. | Reconstruir saldos desde movimientos produce exactamente los saldos proyectados. |
| RF-11 | Mayor y balance de comprobación consistentes. | Filtros de tenant, cuenta y fecha correctos; débitos y créditos globales coinciden sin desbordamiento. |
| RF-12 | Atomicidad de contabilización y reversión. | Un fallo inyectado en cualquier escritura deja todas las tablas sin efectos parciales. |
| RF-13 | Auditoría local atómica y verificable. | Cada cambio auditado añade evento; una falla de cadena revierte el cambio asociado y genera un error explícito. |
| RF-14 | Outbox para integración externa. | La caída externa no bloquea contabilidad; el reintento no duplica el efecto del consumidor. |

## Requisitos numéricos y de transporte

| ID | Requisito | Criterio de aceptación |
| --- | --- | --- |
| NUM-01 | Importes de MXN en centavos enteros `int64`. | No hay aritmética contable de punto flotante; entradas y operaciones fuera de rango se rechazan. |
| NUM-02 | Importes IPC/JSON como cadenas decimales canónicas. | Ida y vuelta exacta, incluidos valores mayores que el entero seguro de JavaScript; rechazar fracciones, exponentes y ceros iniciales. |
| NUM-03 | Resultados acumulados sujetos al mismo rango. | Sumas, saldos y reversión de extremos nunca envuelven ni truncan el valor. |
| NUM-04 | Fechas contables separadas de timestamps. | Una fecha `YYYY-MM-DD` no cambia por conversión de zona; eventos registran UTC con RFC 3339. |

## Requisitos de seguridad y operación

| ID | Requisito | Criterio de aceptación |
| --- | --- | --- |
| SEG-01 | Cadena auditada por tenant con secuencia y enlace. | Secuencias independientes, detección de modificaciones y verificación contra checkpoints confiables según SECURITY. |
| SEG-02 | Claves y secretos fuera del repositorio y los logs. | Inspección sin material sensible; lectura restringida y cifrado autenticado verificado. |
| SEG-03 | Autorización en cada caso de uso. | El transporte no puede omitir controles; selección de tenant no equivale a autenticación. |
| SEG-04 | Diagnóstico limitado. | Endpoints de inspección ausentes o deshabilitados en producción; acceso local y autenticado en desarrollo. |
| OPS-01 | SQLite WAL y claves foráneas activas por conexión. | Verificación efectiva de configuración y prueba de rechazo de referencias inválidas. |
| OPS-02 | Contabilidad local autónoma. | Con red desconectada funcionan borradores, posteo, reversión, consultas y auditoría local. |
| OPS-03 | Backup y restauración consistentes. | Restaurar una copia soportada recupera mayor, proyecciones y auditoría verificables; copiar solo el archivo principal con WAL activo no se acepta como procedimiento. |
| OPS-04 | Actualización offline verificable. | Paquete con autenticidad e integridad comprobables; fallo de verificación impide instalar; migración y recuperación ensayadas. |
| OPS-05 | Identidad de nodo de mínima recopilación. | Atributos operacionales mínimos documentados; no se presenta fingerprinting como atestación de hardware. |

## Rendimiento

PERF-01 es un objetivo pendiente de validación: al menos 250 transacciones contables confirmadas por segundo en un entorno de referencia con SSD NVMe.

El benchmark debe definir hardware, SO, versiones de Go y SQLite/driver, tamaño inicial de base, tenants, concurrencia, líneas por asiento, pragmas, duración y calentamiento. Una operación incluye mayor, saldo, estado y auditoría; no se cuentan INSERT aislados como contabilizaciones.

Registrar throughput confirmado, errores, reintentos y latencias p50/p95/p99. No se aprueba capacidad sin evidencia reproducible y comprobación de integridad posterior. Este objetivo no es una promesa comercial de rendimiento.

## Fuera de la base y pendientes de especificación

Multimoneda, reglas fiscales detalladas, facturación, conciliación, conectores SAT, anclaje externo específico, reapertura de períodos y administración de actualizaciones requieren contratos propios antes de desarrollarse. El mapa de contextos no los convierte en funciones aceptadas.

También deben fijarse antes de publicar la API: límites de entrada, paginación, matriz exhaustiva de errores y mecanismo concreto de sesión/autenticación. No se deben completar estas decisiones mediante supuestos silenciosos.

## Compatibilidad con antivirus y firewalls

| ID | Requisito | Aceptación |
| --- | --- | --- |
| AV-01 | Protecciones activas y configuración identificada. | Registrar SO, producto/edición/versión Kaspersky, firmas y política; completar instalación, inicio y ciclo contable sin exclusiones generales ni desactivar módulos. |
| AV-02 | Release firmada e íntegra. | Verificar firma, editor esperado, sellado temporal y hash del ejecutable e instalador; no interpretar firma válida como aprobación antivirus. |
| AV-03 | Gestión de detecciones. | Evidencia y responsable definidos; no restaurar cuarentena ni subir archivos automáticamente; seguimiento del dictamen y repetición de prueba. |
| NET-01 | Sin exposición externa por defecto. | Identificar procesos propios y relacionados y comprobar listeners; no aceptar LAN/WAN sin función y configuración explícitas. |
| NET-02 | Respeto de políticas. | Instalación, ejecución, actualización y desinstalación no crean excepciones ni desactivan el firewall silenciosamente. |
| NET-03 | Core independiente de salidas externas. | Con salida bloqueada siguen disponibles inicio local, asientos, reversión, mayor y auditoría; integraciones quedan pendientes. |
| OFF-01 | Instalación offline con prerrequisitos. | Ensayo en equipo sin red con WebView2 disponible o instalable offline; documentar cualquier limitación de validación de certificados. |
| UPD-01 | Actualizaciones auténticas y recuperables. | Paquete alterado o firma inválida se rechaza; corte durante actualización no deja datos contables irrecuperables. |

Todos estos controles están pendientes de pruebas de producto. La entrega de los Markdown solo valida su definición documental. No se extiende la aprobación a otra versión del programa, antivirus o sistema sin evidencia.