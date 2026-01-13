-- count_products.sql
SELECT COUNT(DISTINCT p.id)
FROM products p
LEFT JOIN product_variants pv ON pv.product_id = p.id
LEFT JOIN product_categories pc ON pc.product_id = p.id
WHERE ($1::text IS NULL OR p.name ILIKE '%' || $1 || '%')
  AND ($2::int8 IS NULL OR pv.price >= $2)
  AND ($3::int8 IS NULL OR pv.price <= $3)
  AND (COALESCE($4::uuid[], '{}') = '{}' OR pc.category_id = ANY($4));