# FCOS v2.2 — Modelo de seguridad

## 1. Alcance y amenazas

Este documento especifica controles objetivo; no certifica su implementación. Se consideran entradas manipuladas, accesos entre tenants, escrituras concurrentes, modificación del almacenamiento, exposición de secretos y fallos de integraciones.

Una cadena hash local permite detectar inconsistencias, pero un atacante que pueda reemplazar toda la base y sus referencias locales podría reconstruir una cadena coherente. La confianza exige checkpoints verificables y una referencia protegida independiente. No se promete invulnerabilidad frente al control completo del sistema operativo.

## 2. Identidad y autorización

- La sesión autenticada determina actor y tenants autorizados.
- El tenant seleccionado se valida en backend; ningún header o campo del frontend acredita acceso por sí solo.
- Todas las consultas y referencias contables están limitadas al tenant.
- La autorización se aplica en los casos de uso, también cuando se invocan por Wails.
- Debe existir una matriz de permisos revisada antes de publicar el producto: administración, operación contable y consulta de auditoría son capacidades separadas.
- La API HTTP es opcional y por defecto local o deshabilitada. Toda exposición adicional necesita configuración explícita y protección de transporte e identidad.

Los endpoints de diagnóstico están autenticados, limitados a entorno de desarrollo y deshabilitados en producción. No se exponen fuentes, rutas privadas ni secretos en respuestas operacionales.

## 3. Escritura de la cadena de auditoría

Cada tenant tiene una secuencia independiente, protegida por la clave primaria `(tenant_id, sequence_id)`. SQLite serializa al escritor de la base; esto no es un bloqueo independiente por empresa ni una protección contra reescritura offline.

Dentro de la misma transacción escritora del cambio contable:

1. Adquirir la transacción mediante BEGIN IMMEDIATE o mecanismo equivalente del adaptador.
2. Leer el último bloque del tenant.
3. Para una cadena vacía, asignar secuencia 1 y hash previo de 32 bytes cero.
4. En los demás casos, asignar secuencia anterior más uno y utilizar exactamente el hash del último bloque. Rechazar desbordamiento.
5. Crear el payload canónico y calcular sus hashes.
6. Insertar el evento junto con el resto de escrituras del caso de uso.
7. Confirmar todo o ejecutar rollback completo.

El servicio obtiene la secuencia y el hash previo desde el estado transaccional; no acepta esos valores de un cliente. Las inconsistencias producen `AUDIT_INTEGRITY_FAILURE`. Un bloqueo temporal de almacenamiento se reporta separadamente como `STORAGE_BUSY`.

## 4. Formato determinista

El payload incluye versión de esquema, identificador de evento, actor, tenant, operación, entidad afectada, timestamp UTC y datos suficientes para vincular el cambio. Se excluyen contraseñas, claves privadas y documentos secretos completos.

Los importes y la secuencia se codifican como cadenas decimales canónicas. Así se evita depender de representaciones numéricas distintas entre lenguajes. El esquema define también qué campos son obligatorios: omitir un campo y enviarlo como null no son equivalentes por defecto.

```text
canonical_payload = JCS(payload)              # RFC 8785; bytes UTF-8.
payload_hash = SHA256(canonical_payload)

chain_object = {
  "schema_version": "1",
  "tenant_id": tenant_id,
  "sequence_id": secuencia_decimal,
  "previous_hash": hash_previo_hex,
  "payload_hash": hash_payload_hex
}

chain_hash = SHA256(JCS(chain_object))
```

Los hashes hexadecimales utilizan 64 caracteres en minúsculas. Se conserva el payload canónico necesario para volver a verificar el hash. Las implementaciones de canonicalización deben contrastarse con vectores de prueba compartidos entre lenguajes.

Este formato es una decisión de esta edición. No se debe aplicarlo a cadenas existentes sin versión de formato y procedimiento de compatibilidad: no se reescribe historia para adaptarla.

## 5. Checkpoints y límites de confianza

Un checkpoint contiene versión, nodo, tenant, secuencia, hash de cadena, instante UTC e identificador de clave de firma. Se firma una representación canónica de esos campos. El protocolo de firma, algoritmo, codificación y vectores de prueba deben fijarse antes de declarar implementado este control.

La clave pública confiable debe aprovisionarse fuera de la base y de los archivos de checkpoint que protege. El diseño debe elegir un mecanismo verificable, por ejemplo aprovisionamiento firmado o almacén de certificados protegido, y documentar sus permisos y renovación.

Una firma de un checkpoint antiguo no demuestra que sea el último. Para detectar truncamiento o restauración a una versión anterior se necesita una referencia confiable de la última secuencia aceptada, mantenida en un dominio que no pueda revertirse junto con la base. Si no existe esa referencia, el límite debe declararse expresamente.

La verificación recorre secuencia, enlaces, payloads y hashes; comprueba firma e identidad confiable de checkpoints y compara el último estado con la referencia independiente disponible. Ante una inconsistencia se preserva la evidencia y se impiden nuevas escrituras sobre la cadena afectada hasta aplicar un procedimiento de recuperación autorizado. No se reconstruyen hashes silenciosamente.

El anclaje externo es opcional y asíncrono. Una caída del proveedor no impide contabilizar localmente, aunque reduce la actualidad de la evidencia externa disponible.

## 6. Bóveda de secretos

- Cada objeto o versión secreta utiliza una DEK aleatoria de 256 bits.
- El contenido se cifra con AES-256-GCM y un nonce de 96 bits que no se reutiliza con la misma clave.
- Los datos autenticados adicionales vinculan tenant, identificador del secreto y versión.
- La DEK se almacena envuelta mediante RSA-4096 con OAEP-SHA256, identificando la clave utilizada.
- La clave privada de descifrado se protege con permisos del sistema y DPAPI en Windows cuando se use almacenamiento software.
- Una clave de firma se trata como propósito separado; OAEP es un mecanismo de cifrado, no el algoritmo de firma de checkpoints.

La disponibilidad de CNG por sí sola no prueba respaldo de hardware. Solo se declara custodia hardware cuando el proveedor y la configuración verificados realmente la proporcionan.

Antes de producción deben existir procedimientos probados de aprovisionamiento, rotación, recuperación y revocación. Una copia de la base cifrada sin acceso recuperable a las claves no garantiza restauración. La recuperación de material secreto debe limitarse a identidades expresamente autorizadas.

## 7. Operación y mantenimiento

Los backups deben usar un procedimiento consistente con SQLite WAL y conservar los elementos necesarios para verificar auditoría y recuperar secretos. Deben probarse restauraciones en un entorno separado.

Las actualizaciones deben verificar integridad y autenticidad antes de instalarse y contemplar migraciones y recuperación ante fallos. La entrega OTA es opcional.

Los logs no contienen secretos ni payloads sensibles completos. Deben permitir identificar actor, tenant, operación y resultado sin exponer información de otras empresas.

## 8. Reporte de vulnerabilidades

No publicar secretos ni datos contables en issues. El proyecto debe designar y publicar un canal privado de reporte, responsables y versiones soportadas antes de distribuirse públicamente. Esta edición no inventa una dirección de contacto ni plazos de respuesta no acordados.

## 9. Protocolo de convivencia con antivirus y firewalls

Este protocolo es un requisito de ingeniería, no una certificación de compatibilidad ni una promesa de ausencia de detecciones. El entorno prioritario es Windows con Kaspersky; producto, edición y versión exactos del usuario están pendientes de identificar. Otros sistemas requieren su propia matriz de validación antes de declararse soportados.

### 9.1. Comportamiento del programa

No desactivar antivirus, protección de comportamiento o firewall; no crear exclusiones generales, excepciones silenciosas ni mecanismos evasivos. No modificar políticas del sistema automáticamente. Evitar elevación permanente, inyección en otros procesos y descarga o ejecución oculta de herramientas. Documentar procesos hijos legítimos, rutas de instalación, datos, logs y temporales.

Guardar datos mutables en una ubicación de datos de aplicación con permisos del usuario, separada del ejecutable. No almacenar SQLite dentro de carpetas de sistema ni distribuir una base activa mediante carpetas compartidas de red. Las escrituras que requiera el instalador deben ser explícitas y justificadas; no se prohíben por error sus registros normales de desinstalación o accesos directos.

### 9.2. Firma de distribución

La release distribuida en Windows exige Authenticode para ejecutables e instaladores propios, con identidad del editor, cadena de confianza y sellado temporal verificables. Conservar la clave de firma fuera del repositorio y de los equipos de usuario.

Authenticode permite verificar editor e integridad, pero un estado `Valid` no demuestra ausencia de malware ni aprobación de Kaspersky. Las compilaciones de desarrollo sin firma se identifican como tales y no cuentan como evidencia de aceptación de la release. [Microsoft: Authenticode](https://learn.microsoft.com/en-us/windows-hardware/drivers/install/authenticode).

### 9.3. Detección y posible falso positivo en Kaspersky

1. Conservar la alerta: producto y versión, fecha de firmas, módulo que detectó, nombre de detección, acción, ruta, versión del programa y SHA-256 del artefacto.
2. No restaurar automáticamente un archivo de cuarentena ni ejecutarlo para eludir el bloqueo. Una detección no se clasifica de antemano como falsa.
3. Contrastar el hash con el artefacto de distribución y revisar procedencia, firma y comportamiento.
4. El responsable puede consultar el hash y solicitar revisión mediante el portal oficial de Kaspersky o el soporte correspondiente a su producto. Enviar un binario requiere revisar la confidencialidad y autorizar expresamente la subida; no enviar bases contables, credenciales ni documentos de clientes.
5. Registrar la respuesta y repetir las pruebas con la versión corregida o las bases antivirus actualizadas, sin convertir una exclusión en requisito permanente.

Consultar los canales vigentes, sin prometer inclusión en listas de confianza. [Kaspersky OpenTIP](https://opentip.kaspersky.com/howto/) y [documentación de análisis y revisión](https://opentip.kaspersky.com/Help/Doc_data/all-in-one.htm?sectionUrl=SandboxDetectionNamesURL.htm).

### 9.4. Red y firewalls del sistema operativo

Separar tres superficies: servidor de desarrollo, escritorio distribuido y API opcional. Inventariar los sockets reales por proceso y versión; no presuponer que todo Wails abre un puerto TCP ni que nunca lo hace.

- Desarrollo: cualquier servidor de recarga o API de desarrollo debe ligarse explícitamente a loopback IPv4/IPv6 cuando proceda. No usar `0.0.0.0` ni `::` por defecto. El árbol actual requiere contrastar esa configuración durante la reconstrucción.
- Escritorio: no necesita aceptar conexiones LAN/WAN para realizar contabilidad local. Su funcionamiento con el firewall activo debe demostrarse, incluyendo componentes WebView2 relacionados.
- API externa: fuera de la base inicial. Si se habilita posteriormente, documentar proceso, protocolo, puerto, dirección, perfil de red y orígenes autorizados; exigir autenticación y protección del transporte antes de exponerla.
- Salidas externas: solo para funciones opcionales documentadas, después del commit y sin dependencia del core. Bloquearlas debe dejar las operaciones locales utilizables y mostrar el estado pendiente de la integración.

No crear reglas generales de permitir todo ni reglas de bloqueo redundantes para un supuesto puerto efímero. En Windows, evaluar la política efectiva de perfiles público, privado y dominio y cualquier política corporativa. Las excepciones, cuando sean necesarias, deben limitarse al programa y tráfico concretos y ser administradas por el usuario autorizado. [Microsoft: configuración de reglas](https://learn.microsoft.com/en-us/windows/security/operating-system-security/network-security/windows-firewall/configure).

Si Kaspersky administra un firewall propio, identificarlo y registrar su política junto con la del sistema; no asumir que solo hay un filtro. Para macOS o Linux, aplicar el mismo principio de mínimo acceso mediante las herramientas nativas correspondientes, con instrucciones por versión aún pendientes de validación. No trasladar comandos PowerShell a esos sistemas.

### 9.5. Instalación y funcionamiento sin red

El core debe iniciar y operar sin licencias remotas, DNS, telemetría obligatoria ni anclaje externo. Esto no afirma que las consultas de reputación del antivirus o la renovación de certificados funcionen sin red.

Windows requiere WebView2 para Wails: prever runtime ya instalado o un paquete offline de procedencia verificable y arquitectura correcta. No usar un bootstrapper que necesita descargar como única vía para equipos aislados. [Wails: instalación](https://wails.io/docs/gettingstarted/installation/).

Las actualizaciones verifican firma del paquete o manifiesto contra confianza aprovisionada; el manifiesto vincula versión, plataforma y hashes de los artefactos. Rechazar firmas inválidas y retrocesos no autorizados, preparar el archivo antes de sustituirlo y probar recuperación ante interrupción. No cambiar ejecutables silenciosamente ni eludir bloqueos de seguridad durante la actualización.