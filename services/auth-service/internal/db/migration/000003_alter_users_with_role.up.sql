ALTER TABLE users
ADD COLUMN role_id uuid;

INSERT INTO
    roles (name, description)
values
    ('user', 'user สามารถสะสมแต้มได้');

UPDATE users
SET
    role_id = (
        select
            id
        from
            roles r
        where
            r.name = 'user'
    );

ALTER TABLE users
ALTER COLUMN role_id
SET
    NOT NULL;

ALTER TABLE users ADD CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles (id);