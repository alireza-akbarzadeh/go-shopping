-- +goose Up
-- +goose StatementBegin

-- Seed categories
INSERT INTO categories (name, slug, description, parent_id, level, path, is_active, created_at, updated_at) VALUES
    ('Accessories', 'accessories', 'Elevate your style', NULL, 0, 'accessories', true, NOW(), NOW()),
    ('Home & Living', 'home', 'Design your space', NULL, 0, 'home', true, NOW(), NOW()),
    ('Electronics', 'electronics', 'Premium tech essentials', NULL, 0, 'electronics', true, NOW(), NOW()),
    ('Lifestyle', 'lifestyle', 'Curated for you', NULL, 0, 'lifestyle', true, NOW(), NOW());

-- Add subcategories
INSERT INTO categories (name, slug, description, parent_id, level, path, is_active, created_at, updated_at) VALUES
    ('Watches', 'watches', 'Timeless pieces for every occasion', (SELECT id FROM categories WHERE slug = 'accessories'), 1, 'accessories/watches', true, NOW(), NOW()),
    ('Eyewear', 'eyewear', 'See the world in style', (SELECT id FROM categories WHERE slug = 'accessories'), 1, 'accessories/eyewear', true, NOW(), NOW()),
    ('Home', 'home-decor', 'Cozy up your space', (SELECT id FROM categories WHERE slug = 'home'), 1, 'home/home-decor', true, NOW(), NOW()),
    ('Kitchen', 'kitchen', 'Cook like a chef', (SELECT id FROM categories WHERE slug = 'home'), 1, 'home/kitchen', true, NOW(), NOW()),
    ('Lighting', 'lighting', 'Illuminate your world', (SELECT id FROM categories WHERE slug = 'home'), 1, 'home/lighting', true, NOW(), NOW());

-- Seed products with colors and sizes
WITH category_map AS (
    SELECT id, name FROM categories
),
product_data AS (
    SELECT
        'Minimal Leather Bag' AS name,
        189.00 AS price,
        249.00 AS compare_at_price,
        'Accessories' AS cat_name,
        '{"https://images.unsplash.com/photo-1548036328-c9fa89d128fa?w=600&h=600&fit=crop"}'::TEXT[] AS images,
        'Handcrafted from premium Italian leather, this minimalist bag combines elegance with functionality.' AS description,
        true AS is_new,
        '["Brown","Black","Tan"]'::JSONB AS colors,
        '["One Size"]'::JSONB AS sizes
    UNION ALL SELECT 'Premium Watch Collection', 299.00, 399.00, 'Watches', '{"https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=600&h=600&fit=crop"}', 'Swiss-engineered precision meets contemporary design in this stunning timepiece.', false, '["Silver","Gold","Rose Gold"]'::JSONB, '["Strap S","Strap M","Strap L"]'::JSONB
    UNION ALL SELECT 'Designer Sunglasses', 159.00, 199.00, 'Eyewear', '{"https://images.unsplash.com/photo-1572635196237-14b3f281503f?w=600&h=600&fit=crop"}', 'UV-protective lenses paired with a lightweight titanium frame for all-day comfort.', true, '["Black","Tortoise","Silver"]'::JSONB, '["One Size"]'::JSONB
    UNION ALL SELECT 'Wireless Headphones', 249.00, 329.00, 'Electronics', '{"https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=600&h=600&fit=crop"}', 'Immersive sound quality with active noise cancellation and 40-hour battery life.', false, '["Black","White","Navy"]'::JSONB, '["One Size"]'::JSONB
    UNION ALL SELECT 'Ceramic Vase Set', 89.00, 119.00, 'Home', '{"https://images.unsplash.com/photo-1578500494198-246f612d3b3d?w=600&h=600&fit=crop"}', 'Hand-thrown ceramic vases that bring organic beauty to any living space.', true, '["White","Terracotta","Sage"]'::JSONB, '["Small","Medium","Large"]'::JSONB
    UNION ALL SELECT 'Artisan Coffee Maker', 199.00, 259.00, 'Kitchen', '{"https://images.unsplash.com/photo-1517668808822-9ebb02f2a0e6?w=600&h=600&fit=crop"}', 'Precision brewing technology for the perfect cup every morning.', false, '["Stainless Steel","Black","Copper"]'::JSONB, '["Standard"]'::JSONB
    UNION ALL SELECT 'Minimalist Desk Lamp', 129.00, 169.00, 'Lighting', '{"https://images.unsplash.com/photo-1507473885765-e6ed057f782c?w=600&h=600&fit=crop"}', 'Adjustable LED lighting with touch-sensitive controls and wireless charging base.', true, '["Black","White","Brass"]'::JSONB, '["One Size"]'::JSONB
    UNION ALL SELECT 'Leather Wallet', 79.00, 99.00, 'Accessories', '{"https://images.unsplash.com/photo-1627123424574-724758594e93?w=600&h=600&fit=crop"}', 'RFID-blocking technology meets timeless design in full-grain leather.', false, '["Black","Brown","Cognac"]'::JSONB, '["One Size"]'::JSONB
    UNION ALL SELECT 'Smart Speaker', 179.00, 229.00, 'Electronics', '{"https://images.unsplash.com/photo-1543512214-318c7553f230?w=600&h=600&fit=crop"}', 'Room-filling sound with built-in voice assistant and multi-room connectivity.', true, '["Charcoal","White","Sand"]'::JSONB, '["One Size"]'::JSONB
    UNION ALL SELECT 'Silk Scarf', 119.00, 149.00, 'Accessories', '{"https://images.unsplash.com/photo-1584917865442-de89df76afd3?w=600&h=600&fit=crop"}', '100% mulberry silk with hand-rolled edges and exclusive print designs.', false, '["Navy","Burgundy","Emerald"]'::JSONB, '["One Size"]'::JSONB
    UNION ALL SELECT 'Mechanical Keyboard', 169.00, 219.00, 'Electronics', '{"https://images.unsplash.com/photo-1618384887929-16ec33fab9ef?w=600&h=600&fit=crop"}', 'Premium switches with customizable RGB backlighting and aircraft-grade aluminum frame.', true, '["Black","White","Grey"]'::JSONB, '["Standard"]'::JSONB
    UNION ALL SELECT 'Marble Bookends', 69.00, 89.00, 'Home', '{"https://images.unsplash.com/photo-1544457070-4cd773b4d71e?w=600&h=600&fit=crop"}', 'Hand-polished Carrara marble with felt-lined bases to protect surfaces.', false, '["White","Grey"]'::JSONB, '["Pair"]'::JSONB
    UNION ALL SELECT 'Canvas Tote', 59.00, 79.00, 'Accessories', '{"https://images.unsplash.com/photo-1544816155-12df9643f363?w=600&h=600&fit=crop"}', 'Organic cotton canvas with reinforced handles and interior pockets.', false, '["Natural","Black","Navy"]'::JSONB, '["One Size"]'::JSONB
    UNION ALL SELECT 'Copper Pendant Light', 219.00, 279.00, 'Lighting', '{"https://images.unsplash.com/photo-1524484485831-a92ffc0de03f?w=600&h=600&fit=crop"}', 'Hand-spun copper shade with adjustable cord length and warm Edison bulb.', true, '["Polished Copper","Brushed Nickel"]'::JSONB, '["One Size"]'::JSONB
    UNION ALL SELECT 'Chronograph Watch', 449.00, 549.00, 'Watches', '{"https://images.unsplash.com/photo-1539874754764-5a96559165b0?w=600&h=600&fit=crop"}', 'Japanese quartz movement with sapphire crystal and water resistance to 100m.', false, '["Silver","Black","Gold"]'::JSONB, '["Strap M","Strap L"]'::JSONB
    UNION ALL SELECT 'Pour Over Set', 89.00, 109.00, 'Kitchen', '{"https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=600&h=600&fit=crop"}', 'Borosilicate glass carafe with stainless steel dripper and reusable filter.', false, '["Clear","Smoke"]'::JSONB, '["600ml","1L"]'::JSONB
)
INSERT INTO products (
    name, slug, price, compare_at_price, description, category_id,
    stock, sku, images, status, is_digital, low_stock_threshold,
    weight, colors, sizes, created_at, updated_at
)
SELECT
    pd.name,
    LOWER(REPLACE(pd.name, ' ', '-')) || '-' || floor(random()*10000)::text AS slug,
    pd.price,
    pd.compare_at_price,
    pd.description,
    cm.id,
    floor(random() * 191 + 10)::int AS stock,
    UPPER(LEFT(pd.name, 3)) || '-' || floor(random()*10000)::text AS sku,
    pd.images,
    'active' AS status,
    false AS is_digital,
    5 AS low_stock_threshold,
    (random() * 2 + 0.5)::numeric(8,2) AS weight,
    pd.colors,
    pd.sizes,
    NOW(),
    NOW()
FROM product_data pd
JOIN category_map cm ON cm.name = pd.cat_name;

-- Seed reviews for products (using some user IDs – adjust if you have existing users)
-- We'll assume there is at least one admin user (id = 1) and a few regular users (2,3,4)
-- If no users exist, this section will do nothing but not fail.
DO $$
DECLARE
    product_record RECORD;
    user_ids INTEGER[] := ARRAY[1, 2, 3, 4]; -- adjust based on actual users
    comment_templates TEXT[] := ARRAY[
        'Absolutely love this product! High quality and fast shipping.',
        'Good value for money. Would recommend.',
        'Not bad, but expected better. Still satisfied.',
        'Excellent! Exactly as described.',
        'Great product, came quickly. Will order again.',
        'Decent quality, but a bit overpriced.',
        'Perfect! Exceeded my expectations.',
        'The quality is amazing. Five stars!',
        'Works as expected. No complaints.',
        'Beautiful design and durable materials.'
    ];
BEGIN
    FOR product_record IN SELECT id FROM products LOOP
        -- Insert between 1 and 3 random reviews per product
        FOR i IN 1..floor(random()*3+1) LOOP
            INSERT INTO reviews (product_id, user_id, rating, comment, created_at, updated_at)
            VALUES (
                product_record.id,
                user_ids[floor(random()*array_length(user_ids,1)+1)],
                floor(random()*5+1)::int,
                comment_templates[floor(random()*array_length(comment_templates,1)+1)],
                NOW() - (random() * interval '30 days'),
                NOW()
            ) ON CONFLICT DO NOTHING;
        END LOOP;
    END LOOP;
END $$;

-- Update product ratings and reviews_count based on inserted reviews
UPDATE products SET
    rating = sub.avg_rating,
    reviews_count = sub.review_count
FROM (
    SELECT product_id, COALESCE(AVG(rating), 0) AS avg_rating, COUNT(*) AS review_count
    FROM reviews
    GROUP BY product_id
) sub
WHERE products.id = sub.product_id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove seeded data (optional, but rollback will drop tables later anyway)
DELETE FROM reviews WHERE product_id IN (SELECT id FROM products);
DELETE FROM products;
DELETE FROM categories;
-- +goose StatementEnd