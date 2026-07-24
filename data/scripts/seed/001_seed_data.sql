-- Deterministic seed data for local/dev/testing.
-- Safe to re-run: truncates all tables and reseeds with fixed IDs.

BEGIN;

TRUNCATE TABLE
    accountlogin,
    cart_item,
    order_item,
    shipment,
    payment,
    orders,
    cart,
    product_image,
    product,
    category,
    supplier,
    account
RESTART IDENTITY CASCADE;

-- password_hash is bcrypt hash of "test" for all accounts, for testing purposes only.
INSERT INTO account (account_id, full_name, email, password_hash, role, login_attempts, is_active)
VALUES
    (1, 'Thandi Dlamini', 'thandi@example.com', '$2a$12$JiDC9kTZz45tVjzEep9lt.nGa2TPujGOAggwHMrSRIGK.jBtg9TGm', 'customer', 0, TRUE),
    (2, 'Sipho Maseko', 'sipho.admin@example.com', '$2a$12$JiDC9kTZz45tVjzEep9lt.nGa2TPujGOAggwHMrSRIGK.jBtg9TGm', 'admin', 0, TRUE),
    (3, 'Naledi Khumalo', 'naledi.supplier@example.com', '$2a$12$JiDC9kTZz45tVjzEep9lt.nGa2TPujGOAggwHMrSRIGK.jBtg9TGm', 'supplier', 0, TRUE),
    (4, 'Busi Ndlovu', 'busi@example.com', '$2a$12$JiDC9kTZz45tVjzEep9lt.nGa2TPujGOAggwHMrSRIGK.jBtg9TGm', 'customer', 1, TRUE),
    (5, 'Inactive Customer', 'inactive@example.com', '$2a$12$JiDC9kTZz45tVjzEep9lt.nGa2TPujGOAggwHMrSRIGK.jBtg9TGm', 'customer', 0, FALSE);

INSERT INTO accountlogin (token_id, account_id, token_string)
VALUES
    (1, 2, 'token-sipho-1');

INSERT INTO supplier (supplier_id, name, contact_email, contact_phone)
VALUES
    (1, 'Mzansi Gadgets', 'support@mzansigadgets.example', '+27-10-555-0101'),
    (2, 'EcoHome Supplies', 'hello@ecohome.example', '+27-10-555-0102');

INSERT INTO category (category_id, name, description, is_active)
VALUES
    (1, 'Electronics', 'Electronic devices and accessories', TRUE),
    (2, 'Home', 'Home and living products', TRUE),
    (3, 'Fashion', 'Fashion and clothing', TRUE),
    (4, 'Beauty', 'Beauty and accessories', TRUE),
    (5, 'Sports', 'Sport apparel and equipment', TRUE),
    (6, 'Books', 'Books and collections', TRUE),
    (7, 'Tools', 'Tools and machinery', TRUE),
    (8, 'Gaming, Movies & Music', 'Entertainment', TRUE),
    (9, 'Stationary & Office', 'Stationary and officeware', TRUE),
    (10, 'Groceries & Household', 'Food and household products', TRUE);

INSERT INTO product (product_id, supplier_id, category_id, name, description, price, discount_percent, stock_quantity, is_active)
VALUES
    (1, 1, 1, 'Smartphone X', 'Latest model smartphone', 4999.00, 0.00, 25, TRUE),
    (2, 1, 1, 'Wireless Earbuds', 'Noise-cancelling earbuds', 899.00, 10.00, 0, TRUE),
    (3, 2, 2, 'Reusable Bottle', 'Insulated 750ml bottle', 199.00, 15.00, 5, TRUE),
    (4, 2, 1, 'Laptop', 'Macbook Air', 39999.00, 5.00, 2, TRUE),
    (5, 2, 1, 'Drone', 'Sky Surfer', 8199.00, 2.00, 5, TRUE),
    (6, 1, 1, 'Cellphone', 'iPhone 17 Pro Max', 45999.00, 15.00, 51, TRUE),
    (7, 2, 1, 'Samsung', 'A57', 5999.00, 17.00, 5, TRUE),
    (8, 1, 1, 'Camera', 'Olympus Trip', 13000.00, 2.00, 19, TRUE),
    (9, 2, 1, 'Dj Set', 'DDJ 5000', 22199.00, 15.00, 5, TRUE),
    (10, 2, 4, 'Mac Lipsticks', 'Assorted Mac Lipstick colours', 459.00, 37.00, 5, TRUE),
    (11, 2, 4, 'Makeup', 'Makeup pallet', 100.00, 1.00, 5, TRUE),
    (12, 1, 4, 'Lotion', 'Assorted Lotions', 80.00, 15.00, 5, TRUE),
    (13, 2, 4, 'Hair brushes', 'Assorted Hairbrushes', 69.00, 15.00, 5, TRUE),
    (14, 1, 4, 'Hydrating spray', 'Hydrating mist spray', 199.00, 15.00, 5, TRUE),
    (15, 2, 4, 'Lotion Bottle', 'Lotion Bottle', 199.00, 15.00, 5, TRUE);

INSERT INTO product_image (product_image_id, product_id, image_url, alt_text, is_primary, sort_order)
VALUES
    (1, 1, 'Smartphone X - Front.jpg', 'Smartphone X front', TRUE, 1),
    (2, 1, 'Smartphone X - Back.jpg', 'Smartphone X back', FALSE, 2),
    (3, 2, 'Wireless Earbuds.jpg', 'Wireless Earbuds case', TRUE, 1),
    (4, 3, 'Reusable Bottle.jpg', 'Reusable Bottle', TRUE, 1),
    (5, 4, 'electronics1.jpg', 'Macbook Air', TRUE, 1),
    (6, 5, 'electronics2.jpg', 'Sky Surfer', TRUE, 1),
    (7, 6, 'electronics3.jpg', 'iPhone 17 Pro Max', TRUE, 1),
    (8, 7, 'electronics4.jpg', 'A57', TRUE, 1),
    (9, 8, 'electronics5.jpg', 'Olympus Trip', TRUE, 1),
    (10, 9, 'electronics6.jpg', 'DDJ 5000', TRUE, 1),
    (11, 10, 'beauty1.jpg', 'Assorted Mac Lipstick colours', TRUE, 1),
    (12, 11, 'beauty2.jpg', 'Makeup pallet', TRUE, 1),
    (13, 12, 'beauty3.jpg', 'Assorted Lotions', TRUE, 1),
    (14, 13, 'beauty4.jpg', 'Assorted Hairbrushes', TRUE, 1),
    (15, 14, 'beauty5.jpg', 'Hydrating mist spray', TRUE, 1),
    (16, 15, 'beauty6.jpg', 'Lotion Bottle', TRUE, 1);

INSERT INTO cart (cart_id, account_id, session_key, total_price)
VALUES
    (1, 1, 'session-thandi-1', 5397);

INSERT INTO cart_item (cart_id, product_id, quantity, unit_price)
VALUES
    (1, 1, 1, 4999.00),
    (1, 3, 2, 199.00);

INSERT INTO orders (order_id, account_id, order_number, status, subtotal_amount, discount_amount, total_amount, placed_at)
VALUES
    (1, 1, 'ORD-1001', 'paid', 5397.00, 0.00, 5397.00, NOW() - INTERVAL '2 days'),
    (2, 1, 'ORD-1002', 'delivered', 899.00, 100.00, 799.00, NOW() - INTERVAL '10 days'),
    (3, 4, 'ORD-1003', 'cancelled', 899.00, 0.00, 899.00, NOW() - INTERVAL '1 day');

INSERT INTO order_item (order_id, product_id, quantity, unit_price, discount_amount)
VALUES
    (1, 1, 1, 4999.00, 0.00),
    (1, 3, 2, 199.00, 0.00),
    (2, 2, 1, 899.00, 100.00),
    (3, 2, 1, 899.00, 0.00);

INSERT INTO payment (payment_id, order_id, payment_method, payment_status, transaction_ref, amount, paid_at)
VALUES
    (1, 1, 'card', 'paid', 'TXN-1001', 5397.00, NOW() - INTERVAL '2 days'),
    (2, 2, 'eft', 'paid', 'TXN-1002', 799.00, NOW() - INTERVAL '9 days'),
    (3, 3, 'card', 'failed', 'TXN-1003', 899.00, NULL);

INSERT INTO shipment (
    shipment_id,
    order_id,
    courier_name,
    courier_type,
    delivery_reference,
    delivery_address_line1,
    delivery_city,
    delivery_state,
    delivery_postal_code,
    delivery_country,
    recipient_name,
    recipient_phone,
    shipment_status,
    dispatched_at,
    delivered_at
)
VALUES
    (1, 1, 'FastShip', 'standard', 'SHIP-1001', '12 Bree St', 'Cape Town', 'WC', '8001', 'ZA', 'Thandi Dlamini', '+27-82-000-0001', 'in_transit', NOW() - INTERVAL '1 day', NULL),
    (2, 2, 'FastShip', 'express', 'SHIP-1002', '12 Bree St', 'Cape Town', 'WC', '8001', 'ZA', 'Thandi Dlamini', '+27-82-000-0001', 'delivered', NOW() - INTERVAL '9 days', NOW() - INTERVAL '8 days'),
    (3, 3, 'PostX', 'standard', 'SHIP-1003', '55 Market St', 'Johannesburg', 'GP', '2001', 'ZA', 'Busi Ndlovu', '+27-82-000-0002', 'failed', NULL, NULL);

SELECT setval(pg_get_serial_sequence('account', 'account_id'), 5, TRUE);
SELECT setval(pg_get_serial_sequence('accountlogin', 'token_id'), 1, TRUE);
SELECT setval(pg_get_serial_sequence('supplier', 'supplier_id'), 2, TRUE);
SELECT setval(pg_get_serial_sequence('product', 'product_id'), 15, TRUE);
SELECT setval(pg_get_serial_sequence('product_image', 'product_image_id'), 16, TRUE);
SELECT setval(pg_get_serial_sequence('cart', 'cart_id'), 1, TRUE);
SELECT setval(pg_get_serial_sequence('orders', 'order_id'), 3, TRUE);
SELECT setval(pg_get_serial_sequence('payment', 'payment_id'), 3, TRUE);
SELECT setval(pg_get_serial_sequence('shipment', 'shipment_id'), 3, TRUE);

COMMIT;
