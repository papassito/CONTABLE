PRAGMA foreign_keys = ON;

-- =============================================================================
-- 1. BOUNDED CONTEXT: IDENTITY & ACCESS (SQLITE)
-- =============================================================================

CREATE TABLE identity_users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    is_superadmin INTEGER NOT NULL DEFAULT 0,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at_utc TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%SZ', 'NOW'))
);

-- =============================================================================
-- 2. BOUNDED CONTEXT: TENANT & TAXPAYER (SQLITE)
-- =============================================================================

CREATE TABLE tenant_tenants (
    id TEXT PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    legal_name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'ACTIVE',
    created_at_utc TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%SZ', 'NOW'))
);

CREATE TABLE tenant_taxpayers (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL REFERENCES tenant_tenants(id) ON DELETE RESTRICT,
    rfc TEXT NOT NULL,
    legal_name TEXT NOT NULL,
    tax_person_type TEXT NOT NULL,
    state_code TEXT NOT NULL,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at_utc TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%SZ', 'NOW')),
    UNIQUE(tenant_id, id),
    UNIQUE(tenant_id, rfc)
);

-- =============================================================================
-- 3. BOUNDED CONTEXT: COMPLIANCE & EXPEDIENTES (SQLITE)
-- =============================================================================

CREATE TABLE compliance_periods (
    id TEXT PRIMARY KEY,
    period_year INTEGER NOT NULL,
    period_month INTEGER NOT NULL,
    period_type TEXT NOT NULL DEFAULT 'MONTHLY',
    period_number INTEGER NOT NULL,
    created_at_utc TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%SZ', 'NOW')),
    UNIQUE(period_year, period_type, period_number)
);

CREATE TABLE compliance_expedientes (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    taxpayer_id TEXT NOT NULL,
    period_id TEXT NOT NULL REFERENCES compliance_periods(id) ON DELETE RESTRICT,
    authority_domain TEXT NOT NULL,
    expediente_code TEXT NOT NULL,
    current_stage TEXT NOT NULL DEFAULT 'CREATED',
    is_closed INTEGER NOT NULL DEFAULT 0,
    created_at_utc TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%SZ', 'NOW')),
    updated_at_utc TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%SZ', 'NOW')),
    UNIQUE(tenant_id, id),
    FOREIGN KEY(tenant_id, taxpayer_id) REFERENCES tenant_taxpayers(tenant_id, id) ON DELETE RESTRICT,
    UNIQUE(tenant_id, taxpayer_id, period_id, authority_domain)
);

-- =============================================================================
-- 4. BOUNDED CONTEXT: SECRETS & AUDIT LEDGER (SQLITE)
-- =============================================================================

CREATE TABLE secret_envelopes (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL REFERENCES tenant_tenants(id) ON DELETE RESTRICT,
    secret_type TEXT NOT NULL,
    encrypted_payload BLOB NOT NULL,
    wrapped_tek_node BLOB NOT NULL,
    wrapped_tek_recovery BLOB,
    payload_nonce BLOB NOT NULL,
    key_version TEXT NOT NULL,
    created_at_utc TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%SZ', 'NOW')),
    updated_at_utc TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%SZ', 'NOW')),
    UNIQUE(tenant_id, secret_type)
);

CREATE TABLE audit_hash_chain (
    sequence_id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id TEXT NOT NULL UNIQUE,
    tenant_id TEXT NOT NULL REFERENCES tenant_tenants(id) ON DELETE RESTRICT,
    node_id TEXT NOT NULL,
    actor_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    canonical_payload_hash TEXT NOT NULL,
    previous_hash TEXT NOT NULL,
    chain_hash TEXT NOT NULL,
    created_at_utc TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%SZ', 'NOW'))
);

-- Disparadores de Inmutabilidad en SQLite
CREATE TRIGGER trg_prevent_audit_update
BEFORE UPDATE ON audit_hash_chain
BEGIN
    SELECT RAISE(ABORT, 'VIOLACIÓN CONSTITUCIONAL FCOS v2.2: La cadena de auditoría es inmutable.');
END;

CREATE TRIGGER trg_prevent_audit_delete
BEFORE DELETE ON audit_hash_chain
BEGIN
    SELECT RAISE(ABORT, 'VIOLACIÓN CONSTITUCIONAL FCOS v2.2: Prohibido eliminar eventos de auditoría.');
END;