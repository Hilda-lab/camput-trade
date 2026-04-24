CREATE OR REPLACE VIEW sold_items_view AS
SELECT i.item_name, o.buyer_id
FROM orders o
JOIN item i ON o.item_id = i.item_id;

CREATE OR REPLACE VIEW unsold_items_view AS
SELECT item_id, item_name, category, price, seller_id, created_at
FROM item
WHERE status = 0;
