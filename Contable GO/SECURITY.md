# SECURITY.md — Modelo de Seguridad, Cifrado y Cadena de Auditoría (FCOS v2.2)

---

## 1. Protocolo Anti-Fork y Serialización de Escritura en SQLite

La cadena de auditoría criptográfica local (Local Hash Chain) es un mecanismo **Tamper-Evident** que permite detectar modificaciones no autorizadas en el almacenamiento en reposo.

### A. Protocolo de Transacción de Escritura
La escritura de la cadena se serializa mediante las garantías transaccionales de SQLite. `BEGIN IMMEDIATE` adquiere anticipadamente la reserva de escritura necesaria para impedir que otra transacción escritora modifique simultáneamente el estado de la base de datos.

*La serialización de SQLite constituye una garantía a nivel de motor de base de datos, no un bloqueo independiente por `tenant_id`.*

El aislamiento lógico por tenant se garantiza adicionalmente mediante:
* `tenant_id` obligatorio.
* Consultas filtradas por `tenant_id`.
* Secuencia independiente por tenant.
* `UNIQUE (tenant_id, sequence_id)`.
* Validación estricta del `previous_hash`.

1. **Iniciar Transacción:** `BEGIN IMMEDIATE`.
2. **Obtener Último Bloque:** `SELECT sequence_id, chain_hash FROM audit_hash_chain WHERE tenant_id = ? ORDER BY sequence_id DESC LIMIT 1`.
3. **Validar Invariante de Secuencia:** Verificar que el `sequence_id` a insertar sea exactamente $\text{sequence\_id}_{\text{last}} + 1$.
4. **Validar Invariante de Enlace:** Verificar que `previous_hash` coincida exactamente con $\text{chain\_hash}_{\text{last}}$.
5. **Canonicalización del Payload y Cálculo del Hash:** El payload del evento se serializa mediante RFC 8785 (JCS) antes de calcular `payload_hash = SHA256(canonical_payload)`. Posteriormente, `chain_hash` se calcula sobre una representación UTF-8 determinista que incorpora la secuencia, tenant, hash previo y el hash del payload:
   $$\text{chain\_input} = \text{UTF-8}\Big(\text{decimal}(\text{sequence\_id}) \parallel \text{":"} \parallel \text{tenant\_id} \parallel \text{":"} \parallel \text{hex}(\text{previous\_hash}) \parallel \text{":"} \parallel \text{hex}(\text{payload\_hash})\Big)$$
   $$\text{chain\_hash}_n = \text{SHA256}(\text{chain\_input})$$
6. **Insertar Evento:** Ejecutar `INSERT INTO audit_hash_chain ...`.
7. **Commit:** Confirmar transacción. En caso de inconsistencia, ejecutar `ROLLBACK` y retornar `ERR_CHAIN_FORK_ATTEMPT`.

---

## 2. Arquitectura de Confianza de 3 Niveles (Trust Anchors)

Para fortalecer la propiedad **Tamper-Evident** sin vulnerar la autonomía offline, FCOS v2.2 implementa una arquitectura de confianza en 3 niveles:

```text
[Nivel 1: Local Hash Chain] ──> Detecta alteraciones en las filas de SQLite.
              │
              v
[Nivel 2: Signed Chain Head Checkpoint] ──> Firma periódica del tip de la cadena por el Nodo.
              │
              v (Opcional / Asíncrono)
[Nivel 3: External Anchoring] ──> Publicación asíncrona del punto de control en Web3/TSA.

### A. Anclaje de Confianza de la Clave Pública del Nodo
La clave pública de identidad del nodo deberá estar anclada en un dominio de confianza independiente de la base SQLite y de los archivos de punto de control almacenados junto a ella.
Los mecanismos válidos podrán incluir:
* Windows Certificate Store para la identidad/certificado confiable.
* Payload de aprovisionamiento firmado durante la instalación inicial.
* TPM / Windows CNG cuando exista soporte de hardware.

*(Nota: DPAPI se reserva para la protección del material secreto y la clave privada local, no para establecer la confianza de la clave pública).*

## 3. Bóveda de Secretos y Cifrado de Sobres (Secrets Vault)
Las credenciales sensibles se cifran utilizando *Envelope Encryption*:

### A. Definición de Llaves Criptográficas
* **DEK (Data Encryption Key):** Llave simétrica única por objeto/versión secreta para cifrar el texto plano con AES-256-GCM, utilizando un IV/Nonce único y aleatorio por operación de cifrado. La DEK se almacena exclusivamente en forma envuelta mediante la clave pública del nodo.
* **Node Key Pair (Par de Claves RSA del Nodo):** Par asimétrico RSA-4096 / RSA-OAEP-SHA256 asignado al nodo.
  - *Llave Pública:* Actúa como KEK (Key Encryption Key) para envolver (*wrap*) la DEK.
  - *Llave Privada:* Se utiliza para desenvolver (*unwrap*) la DEK al leer el secreto.

### B. Nivel de Protección de la Clave Privada del Nodo
* **Software-Protected Key:** Por omisión, la clave privada del nodo se almacena protegida por ACLs del sistema operativo y DPAPI en Windows.
* **Hardware-Backed Key:** Únicamente cuando el hardware subyacente cuente con un módulo TPM (Trusted Platform Module) o Windows CNG configurado, la clave privada delegará su custodia directamente al hardware.