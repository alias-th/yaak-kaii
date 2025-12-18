CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE
    IF NOT EXISTS users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        email text NOT NULL UNIQUE,
        email_verified boolean NOT NULL DEFAULT false,
        password_hash text NOT NULL,
        first_name varchar(100) NOT NULL,
        last_name varchar(100) NOT NULL,
        phone_number varchar(50) NOT NULL,
        is_active boolean NOT NULL DEFAULT true,
        created_at timestamptz NOT NULL DEFAULT now (),
        updated_at timestamptz NOT NULL DEFAULT now (),
        deleted_at timestamptz
    );