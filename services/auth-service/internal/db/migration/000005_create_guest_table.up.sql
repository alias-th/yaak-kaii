create table
    if not exists guests (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        token_hash text NOT NULL UNIQUE,
        ip_addr inet NOT NULL,
        user_agent text NOT NULL,
        created_at timestamptz NOT NULL DEFAULT now(),
        expires_at timestamptz,
        metadata jsonb
    );