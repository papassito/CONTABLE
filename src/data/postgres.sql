-- =============================================================================
-- BASE & EXTENSIONES
-- =============================================================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =============================================================================
-- 1. BOUNDED CONTEXT: IDENTITY & ACCESS
-- =============================================================================

CREATE TABLE identity_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    is_superadmin BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC')
);

CREATE TABLE identity_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT NULL
);

CREATE TABLE identity_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    module VARCHAR(50) NOT NULL
);

CREATE TABLE identity_role_permissions (
    role_id UUID NOT NULL REFERENCES identity_roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES identity_permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE identity_tenant_memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL, -- FK definida tras crear tenant_tenants
    user_id UUID NOT NULL REFERENCES identity_users(id) ON DELETE RESTRICT,
    role_id UUID NOT NULL REFERENCES identity_roles(id) ON DELETE RESTRICT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_user_tenant UNIQUE (tenant_id, user_id)
);

CREATE TABLE identity_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES identity_users(id) ON DELETE CASCADE,
    token_hash CHAR(64) NOT NULL UNIQUE,
    ip_address INET NOT NULL,
    user_agent TEXT NOT NULL,
    expires_at_utc TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC')
);

-- =============================================================================
-- 2. BOUNDED CONTEXT: TENANT & TAXPAYER
-- =============================================================================

CREATE TABLE tenant_tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(30) NOT NULL UNIQUE,
    legal_name VARCHAR(200) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT chk_tenant_status CHECK (status IN ('ACTIVE', 'SUSPENDED', 'ARCHIVED'))
);

ALTER TABLE identity_tenant_memberships 
    ADD CONSTRAINT fk_memberships_tenant 
    FOREIGN KEY (tenant_id) REFERENCES tenant_tenants(id) ON DELETE RESTRICT;

CREATE TABLE tenant_taxpayers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenant_tenants(id) ON DELETE RESTRICT,
    rfc VARCHAR(13) NOT NULL,
    legal_name VARCHAR(250) NOT NULL,
    tax_person_type VARCHAR(10) NOT NULL,
    state_code VARCHAR(5) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_taxpayer_tenant_id UNIQUE (tenant_id, id),
    CONSTRAINT uq_tenant_rfc UNIQUE (tenant_id, rfc),
    CONSTRAINT chk_rfc_format CHECK (rfc ~* '^[A-Z&Ñ]{3,4}[0-9]{6}[A-Z0-9]{3}$')
);

CREATE TABLE tenant_patronal_registers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    taxpayer_id UUID NOT NULL,
    registro_patronal VARCHAR(11) NOT NULL,
    risk_class VARCHAR(5) NOT NULL,
    risk_premium NUMERIC(8, 5) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_patronal_reg_tenant_id UNIQUE (tenant_id, id),
    CONSTRAINT fk_patronal_taxpayer FOREIGN KEY (tenant_id, taxpayer_id) 
        REFERENCES tenant_taxpayers(tenant_id, id) ON DELETE RESTRICT,
    CONSTRAINT uq_taxpayer_reg_patronal UNIQUE (tenant_id, taxpayer_id, registro_patronal)
);

CREATE TABLE tenant_workers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    patronal_register_id UUID NOT NULL,
    curp VARCHAR(18) NOT NULL,
    nss VARCHAR(11) NOT NULL,
    rfc VARCHAR(13) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    daily_integrated_salary_cents BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_worker_tenant_id UNIQUE (tenant_id, id),
    CONSTRAINT fk_worker_patronal FOREIGN KEY (tenant_id, patronal_register_id) 
        REFERENCES tenant_patronal_registers(tenant_id, id) ON DELETE RESTRICT,
    CONSTRAINT uq_register_nss UNIQUE (tenant_id, patronal_register_id, nss)
);

-- =============================================================================
-- 3. BOUNDED CONTEXT: NORMATIVE ENGINE (Global)
-- =============================================================================

CREATE TABLE normative_official_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    authority_type VARCHAR(50) NOT NULL,
    is_official BOOLEAN NOT NULL DEFAULT TRUE,
    created_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC')
);

CREATE TABLE normative_parameter_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES normative_official_sources(id) ON DELETE RESTRICT,
    parameter_code VARCHAR(50) NOT NULL,
    jurisdiction VARCHAR(20) NOT NULL DEFAULT 'FEDERAL',
    numeric_value NUMERIC(14, 6) NOT NULL,
    effective_from DATE NOT NULL,
    effective_to DATE NULL,
    publication_date DATE NOT NULL,
    document_reference VARCHAR(255) NOT NULL,
    document_hash CHAR(64) NOT NULL,
    created_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT chk_param_dates CHECK (effective_to IS NULL OR effective_to >= effective_from)
);

CREATE TABLE normative_rule_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_code VARCHAR(50) NOT NULL,
    version_label VARCHAR(20) NOT NULL,
    description TEXT NOT NULL,
    formula_json JSONB NOT NULL,
    effective_from DATE NOT NULL,
    effective_to DATE NULL,
    created_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_rule_code_version UNIQUE (rule_code, version_label)
);

-- =============================================================================
-- 4. BOUNDED CONTEXT: COMPLIANCE & EXPEDIENTES (Entidad Pivote)
-- =============================================================================

CREATE TABLE compliance_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    period_year INT NOT NULL,
    period_month INT NOT NULL,
    period_type VARCHAR(20) NOT NULL DEFAULT 'MONTHLY',
    period_number INT NOT NULL,
    created_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_fiscal_period UNIQUE (period_year, period_type, period_number),
    CONSTRAINT chk_month_range CHECK (period_month BETWEEN 1 AND 12)
);

CREATE TABLE compliance_expedientes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    taxpayer_id UUID NOT NULL,
    period_id UUID NOT NULL REFERENCES compliance_periods(id) ON DELETE RESTRICT,
    authority_domain VARCHAR(20) NOT NULL,
    expediente_code VARCHAR(50) NOT NULL,
    current_stage VARCHAR(30) NOT NULL DEFAULT 'CREATED',
    is_closed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_expediente_tenant_id UNIQUE (tenant_id, id),
    CONSTRAINT fk_expediente_taxpayer FOREIGN KEY (tenant_id, taxpayer_id) 
        REFERENCES tenant_taxpayers(tenant_id, id) ON DELETE RESTRICT,
    CONSTRAINT uq_expediente_domain UNIQUE (tenant_id, taxpayer_id, period_id, authority_domain),
    CONSTRAINT chk_expediente_stage CHECK (
        current_stage IN (
            'CREATED', 'PENDING', 'PROCESSING', 'VALIDATION', 'READY', 
            'REQUIRES_HUMAN', 'EXECUTED', 'EVIDENCE_CAPTURED', 'RECONCILED', 
            'COMPLETED', 'FAILED', 'EXTERNAL_UNAVAILABLE', 'RETRY', 'CANCELLED', 'BLOCKED'
        )
    )
);

-- =============================================================================
-- 6. BOUNDED CONTEXT: SECRETS VAULT
-- =============================================================================

CREATE TABLE secret_envelopes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenant_tenants(id) ON DELETE RESTRICT,
    secret_type VARCHAR(50) NOT NULL,
    encrypted_payload BYTEA NOT NULL,
    wrapped_tek_node BYTEA NOT NULL,
    wrapped_tek_recovery BYTEA NULL,
    payload_nonce BYTEA NOT NULL,
    key_version VARCHAR(20) NOT NULL,
    created_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_tenant_secret_type UNIQUE (tenant_id, secret_type)
);

-- =============================================================================
-- 7. BOUNDED CONTEXT: AUDIT LEDGER (Append-Only)
-- =============================================================================

CREATE TABLE audit_hash_chain (
    sequence_id BIGSERIAL PRIMARY KEY,
    event_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenant_tenants(id) ON DELETE RESTRICT,
    node_id VARCHAR(100) NOT NULL,
    actor_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    canonical_payload_hash CHAR(64) NOT NULL,
    previous_hash CHAR(64) NOT NULL,
    chain_hash CHAR(64) NOT NULL,
    created_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC')
);

-- =============================================================================
-- POLÍTICAS RLS Y TRIGGERS DE PROTECCIÓN EN POSTGRESQL
-- =============================================================================

CREATE OR REPLACE FUNCTION fn_prevent_audit_tampering()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'VIOLACIÓN CONSTITUCIONAL FCOS v2.2: La cadena de auditoría (audit_hash_chain) es inmutable.';
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Registro dinámico de selectores para robots RPA
CREATE TABLE rpa_layout_registry (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    portal_code VARCHAR(50) NOT NULL,
    element_key VARCHAR(100) NOT NULL,
    css_selector VARCHAR(255) NOT NULL,
    xpath_selector VARCHAR(255) NOT NULL,
    version_tag VARCHAR(20) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at_utc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_portal_element UNIQUE (portal_code, element_key, version_tag)
);