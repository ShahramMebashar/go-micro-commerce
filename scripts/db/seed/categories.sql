
-- scripts/db/seed/categories.sql
INSERT INTO categories (id, name, description, created_at, updated_at) 
VALUES 
('550e8400-e29b-41d4-a716-446655440000', 'Electronics', 'Electronic devices and accessories', NOW(), NOW()),
('650e8400-e29b-41d4-a716-446655440000', 'Clothing', 'Apparel and fashion items', NOW(), NOW()),
('750e8400-e29b-41d4-a716-446655440000', 'Books', 'Books and publications', NOW(), NOW()),
('850e8400-e29b-41d4-a716-446655440000', 'Home & Kitchen', 'Home and kitchen products', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

