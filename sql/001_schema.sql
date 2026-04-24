CREATE TABLE IF NOT EXISTS app_user (
    user_id VARCHAR(20) PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS item (
    item_id VARCHAR(20) PRIMARY KEY,
    item_name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    price NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    seller_id VARCHAR(20) NOT NULL REFERENCES app_user(user_id),
    status SMALLINT NOT NULL DEFAULT 0 CHECK (status IN (0, 1)),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
    order_id VARCHAR(30) PRIMARY KEY,
    buyer_id VARCHAR(20) NOT NULL REFERENCES app_user(user_id),
    item_id VARCHAR(20) NOT NULL UNIQUE REFERENCES item(item_id),
    order_date TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_item_seller_id ON item(seller_id);
CREATE INDEX IF NOT EXISTS idx_item_status ON item(status);
CREATE INDEX IF NOT EXISTS idx_orders_buyer_id ON orders(buyer_id);
