-- products_seed.sql
INSERT INTO products (id, name, description, price, sku, category_id, created_at, updated_at)
VALUES 
(gen_random_uuid(), 'Smartphone X', 'Latest smartphone with advanced features', 999.99, 'PHONE-001', '550e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'Laptop Pro', 'High-performance laptop for professionals', 1499.99, 'LAPTOP-001', '550e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'Wireless Earbuds', 'Premium wireless earbuds with noise cancellation', 199.99, 'AUDIO-001', '550e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'Smart Watch', 'Fitness and health tracking smartwatch', 299.99, 'WATCH-001', '550e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'Tablet Ultra', 'Lightweight tablet with high-resolution display', 699.99, 'TABLET-001', '550e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'Bluetooth Speaker', 'Portable speaker with rich sound', 129.99, 'AUDIO-002', '550e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'Gaming Console', 'Next-gen gaming console with 4K support', 499.99, 'GAME-001', '550e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'Wireless Mouse', 'Ergonomic wireless mouse', 49.99, 'PC-001', '550e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'External SSD', '1TB external solid-state drive', 159.99, 'STORAGE-001', '550e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
(gen_random_uuid(), 'Wireless Keyboard', 'Mechanical wireless keyboard', 89.99, 'PC-002', '550e8400-e29b-41d4-a716-446655440000', NOW(), NOW());