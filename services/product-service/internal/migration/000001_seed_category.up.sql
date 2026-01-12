INSERT INTO
    categories (id, name, description, created_at, updated_at)
VALUES
    (
        "00000000-0000-0000-0000-000000000001",
        'Electronics',
        'Devices and gadgets',
        NOW (),
        NOW ()
    ),
    (
        "00000000-0000-0000-0000-000000000002",
        'Books',
        'Printed and digital books',
        NOW (),
        NOW ()
    ),
    (
        "00000000-0000-0000-0000-000000000003",
        'Clothing',
        'Apparel and accessories',
        NOW (),
        NOW ()
    ),
    (
        "00000000-0000-0000-0000-000000000004",
        'Home & Kitchen',
        'Household items and kitchenware',
        NOW (),
        NOW ()
    ),
    (
        "00000000-0000-0000-0000-000000000005",
        'Sports & Outdoors',
        'Sporting goods and outdoor equipment',
        NOW (),
        NOW ()
    );

INSERT INTO
    category_attributes (
        id,
        category_id,
        key,
        label,
        type,
        required,
        options,
        created_at,
        updated_at
    )
VALUES
    (
        "10000000-0000-0000-0000-000000000001",
        "00000000-0000-0000-0000-000000000003",
        'size',
        'Size',
        'select',
        TRUE,
        '["XS", "S", "M", "L", "XL", "XXL"]',
        NOW (),
        NOW ()
    ),
    (
        "10000000-0000-0000-0000-000000000002",
        "00000000-0000-0000-0000-000000000003",
        'color',
        'Color',
        'select',
        TRUE,
        '["Red", "Blue", "Green", "Black", "White"]',
        NOW (),
        NOW ()
    ),
    (
        "10000000-0000-0000-0000-000000000003",
        "00000000-0000-0000-0000-000000000003",
        'material',
        'Material',
        'text',
        false,
        NULL,
        NOW (),
        NOW ()
    ),
    (
        "10000000-0000-0000-0000-000000000004",
        "00000000-0000-0000-0000-000000000001",
        'gender',
        'Gender',
        'select',
        false,
        '["Male", "Female", "Unisex"]',
        NOW (),
        NOW ()
    ),;