-- Drop indexes
DROP INDEX IF EXISTS idx_products_sku;
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_category_id;

-- Drop tables
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
