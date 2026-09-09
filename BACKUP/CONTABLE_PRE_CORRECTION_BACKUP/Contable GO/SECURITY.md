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
