CREATE TYPE order_status AS ENUM ('open', 'close');
CREATE TYPE payment_method AS ENUM ('cash', 'card', 'kaspi_qr');
CREATE TYPE item_size AS ENUM ('small', 'medium', 'large');
CREATE TYPE transaction_type AS ENUM ('added', 'written off', 'sale', 'created');

CREATE TABLE IF NOT EXISTS inventory (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    stock NUMERIC(10, 2) NOT NULL,
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    unit_type TEXT NOT NULL,
    last_updated TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    preferences JSONB
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id),
    total_amount NUMERIC(10, 2) NOT NULL CHECK (total_amount >= 0),
    status order_status NOT NULL DEFAULT 'open',
    special_instructions JSONB,
    payment_method payment_method NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS menu_items (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    allergens TEXT[],
    size item_size NOT NULL,
    CONSTRAINT unique_menu_item_size UNIQUE (name, size)
);

CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id TEXT NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    quantity NUMERIC NOT NULL CHECK (quantity > 0),
    price_at_order NUMERIC(10, 2) NOT NULL CHECK (price_at_order >= 0)
);

CREATE TABLE IF NOT EXISTS menu_item_ingredients (
    id SERIAL PRIMARY KEY,
    menu_item_id TEXT NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    ingredient_id TEXT NOT NULL REFERENCES inventory(id) ON DELETE CASCADE,
    quantity NUMERIC NOT NULL CHECK (quantity > 0),
    CONSTRAINT unique_menu_item_ingredient UNIQUE (menu_item_id, ingredient_id)
);

CREATE TABLE IF NOT EXISTS order_status_history (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    previous_status order_status NOT NULL,
    new_status order_status NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS price_history (
    id SERIAL PRIMARY KEY,
    menu_item_id TEXT NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    old_price NUMERIC(10, 2) NOT NULL CHECK (old_price >= 0),
    new_price NUMERIC(10, 2) NOT NULL CHECK (new_price >= 0),
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS inventory_transactions (
    id SERIAL PRIMARY KEY,
    inventory_id TEXT NOT NULL REFERENCES inventory(id) ON DELETE CASCADE,
    change_amount NUMERIC NOT NULL,
    transaction_type transaction_type NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

-- ALTER TABLE menu_items ADD COLUMN search_vector tsvector;
-- ALTER TABLE orders ADD COLUMN search_vector tsvector;

-- UPDATE menu_items 
-- SET search_vector = to_tsvector('english', name || ' ' || COALESCE(description, ''));
--triger for menu)items
-- CREATE OR REPLACE FUNCTION update_menu_items_search_vector() 
-- RETURNS TRIGGER AS $$
-- BEGIN
--     UPDATE menu_items 
--     SET search_vector = to_tsvector('english', 
--         NEW.name || ' ' || 
--         COALESCE(NEW.description, '') || ' ' || 
--         COALESCE(array_to_string(NEW.allergens, ' '), ''))
--     WHERE id = NEW.id;
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE TRIGGER menu_items_search_vector_trigger
-- AFTER INSERT OR UPDATE ON menu_items
-- FOR EACH ROW EXECUTE FUNCTION update_menu_items_search_vector();


-- UPDATE orders 
-- SET search_vector = to_tsvector('english', customers.name || ' ' || COALESCE(menu_items.name, ''))
-- FROM customers
-- JOIN order_items oi ON oi.order_id = orders.id
-- JOIN menu_items menu_items ON menu_items.id = oi.menu_item_id
-- WHERE orders.customer_id = customers.id;
--trigger for orders
-- CREATE OR REPLACE FUNCTION update_orders_search_vector() 
-- RETURNS TRIGGER AS $$
-- BEGIN
--     UPDATE orders 
--     SET search_vector = to_tsvector('english', 
--         customers.name || ' ' || 
--         COALESCE(CAST(NEW.total_amount AS TEXT), '') || ' ' || 
--         COALESCE(NEW.status::text, '') || ' ' || 
--         COALESCE(NEW.payment_method::text, '') || ' ' || 
--         COALESCE(NEW.special_instructions::text, '') || ' ' || 
--         COALESCE(menu_items.name, '')
--     )
--     FROM customers
--     JOIN order_items oi ON oi.order_id = NEW.id
--     JOIN menu_items menu_items ON menu_items.id = oi.menu_item_id
--     WHERE orders.id = NEW.id AND orders.customer_id = customers.id;
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE TRIGGER orders_search_vector_trigger
-- AFTER INSERT OR UPDATE ON orders
-- FOR EACH ROW EXECUTE FUNCTION update_orders_search_vector();

-- CREATE INDEX menu_items_search_idx ON menu_items USING gin(search_vector);
-- CREATE INDEX orders_search_idx ON orders USING gin(search_vector);

CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_menu_item_id ON order_items(menu_item_id);
CREATE INDEX idx_menu_items_name ON menu_items(name);
CREATE INDEX idx_inventory_name ON inventory(name);
CREATE INDEX idx_inventory_price ON inventory(price);
CREATE INDEX idx_inventory_stock_level ON inventory(stock);
CREATE INDEX idx_order_status_history_order_id ON order_status_history(order_id);
CREATE INDEX idx_price_history_menu_item_id ON price_history(menu_item_id);
