CREATE TABLE tenants (
    -- base
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    name VARCHAR (255),
    state SMALLINT NOT NULL,
    timezone VARCHAR (63) NOT NULL DEFAULT 'UTC',
    language VARCHAR (2) NOT NULL DEFAULT 'EN',
    created_time TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    modified_time TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);
CREATE UNIQUE INDEX idx_tenant_uuid
ON tenants (uuid);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    first_name VARCHAR (50) NULL,
    last_name VARCHAR (50) NULL,
    email VARCHAR (255) UNIQUE NOT NULL,
    name VARCHAR (255) NULL DEFAULT '',
    lexical_name VARCHAR (255) NULL DEFAULT '',
    password_algorithm VARCHAR (63) NOT NULL,
    password_hash VARCHAR (511) NOT NULL,
    state SMALLINT NOT NULL DEFAULT 0,
    role_id SMALLINT NOT NULL DEFAULT 0,
    timezone VARCHAR (63) NOT NULL DEFAULT 'UTC',
    language VARCHAR (2) NOT NULL DEFAULT 'EN',
    created_time TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    modified_time TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    joined_time TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    salt VARCHAR (127) NOT NULL DEFAULT '',
    was_email_activated BOOLEAN NOT NULL DEFAULT FALSE,
    pr_access_code VARCHAR (127) NOT NULL DEFAULT '',
    pr_expiry_time TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_user_uuid
ON users (uuid);
CREATE UNIQUE INDEX idx_user_email
ON users (email);
CREATE INDEX idx_user_tenant_id
ON users (tenant_id);

CREATE TABLE applications (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    name VARCHAR (255) NULL DEFAULT '',
    description TEXT NOT NULL,
    scope VARCHAR (255) NULL DEFAULT '',
    redirect_url VARCHAR (255) NULL DEFAULT '',
    image_url VARCHAR (255) NULL DEFAULT '',
    state SMALLINT NOT NULL DEFAULT 0,
    client_id VARCHAR (511) NULL DEFAULT '',
    client_secret TEXT NULL DEFAULT '',
    created_time TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    modified_time TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_application_uuid
ON applications (uuid);
CREATE INDEX idx_application_tenant_id
ON applications (tenant_id);

CREATE TABLE authorized_applications (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    application_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    state SMALLINT NOT NULL DEFAULT 0,
    created_time TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    modified_time TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_authorized_application_uuid
ON authorized_applications (uuid);
CREATE INDEX idx_authorized_application_tenant_id
ON authorized_applications (tenant_id);
CREATE INDEX idx_authorized_application_application_id
ON authorized_applications (application_id);
