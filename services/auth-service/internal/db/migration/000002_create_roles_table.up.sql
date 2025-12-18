CREATE TABLE
    IF NOT EXISTS roles (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        name varchar(255) not null,
        description text not null
    );