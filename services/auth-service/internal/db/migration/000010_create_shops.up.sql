CREATE TABLE
    IF NOT EXISTS shops (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        seller_id UUID NOT NULL UNIQUE REFERENCES sellers (id) ON DELETE CASCADE,
        name TEXT NOT NULL,
        slug TEXT NOT NULL UNIQUE,
        description TEXT,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );