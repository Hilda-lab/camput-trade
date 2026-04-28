CREATE DATABASE IF NOT EXISTS campus_trade CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE campus_trade;

CREATE TABLE IF NOT EXISTS app_user (
    user_id VARCHAR(20) PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS item (
    item_id VARCHAR(20) PRIMARY KEY,
    item_name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    seller_id VARCHAR(20) NOT NULL,
    status TINYINT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_item_price CHECK (price >= 0),
    CONSTRAINT chk_item_status CHECK (status IN (0, 1)),
    CONSTRAINT fk_item_seller FOREIGN KEY (seller_id) REFERENCES app_user(user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS orders (
    order_id VARCHAR(30) PRIMARY KEY,
    buyer_id VARCHAR(20) NOT NULL,
    item_id VARCHAR(20) NOT NULL,
    order_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_orders_item_id (item_id),
    CONSTRAINT fk_orders_buyer FOREIGN KEY (buyer_id) REFERENCES app_user(user_id),
    CONSTRAINT fk_orders_item FOREIGN KEY (item_id) REFERENCES item(item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_item_seller_id ON item(seller_id);
CREATE INDEX idx_item_status ON item(status);
CREATE INDEX idx_orders_buyer_id ON orders(buyer_id);
