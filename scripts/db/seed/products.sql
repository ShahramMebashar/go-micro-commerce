
-- scripts/db/seed/products.sql
INSERT INTO products (id, name, description, price, sku, category_id, created_at, updated_at) 
VALUES 
(gen_random_uuid(), 'Smartphone X', 'Latest smartphone with advanced features', 999.99, 'PHONE-001', '550e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'Laptop Pro', 'High-performance laptop for professionals', 1499.99, 'LAPTOP-001', '550e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'Wireless Earbuds', 'Premium wireless earbuds with noise cancellation', 199.99, 'AUDIO-001', '550e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'Smart Watch', 'Fitness and health tracking smartwatch', 299.99, 'WATCH-001', '550e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'Tablet Ultra', 'Lightweight tablet with high-resolution display', 699.99, 'TABLET-001', '550e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'T-Shirt', 'Cotton t-shirt', 29.99, 'SHIRT-001', '650e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'Jeans', 'Denim jeans', 59.99, 'PANTS-001', '650e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'Novel', 'Bestselling fiction novel', 14.99, 'BOOK-001', '750e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'Cookbook', 'Recipe collection', 24.99, 'BOOK-002', '750e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'Coffee Maker', 'Automatic coffee machine', 89.99, 'KITCHEN-001', '850e8400-e29b-41d4-a716-446655440000', NOW(), NOW());

