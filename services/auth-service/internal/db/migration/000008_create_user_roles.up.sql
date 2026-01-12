CREATE TABLE
    IF NOT EXISTS user_roles (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id UUID NOT NULL,
        role_id UUID NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id),
        FOREIGN KEY (role_id) REFERENCES roles (id)
    );