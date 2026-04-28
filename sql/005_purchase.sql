-- 购买商品（事务模板）
START TRANSACTION;

-- 1) 锁定商品并校验未售
-- SELECT status FROM item WHERE item_id = ? FOR UPDATE;

-- 2) 新增订单
-- INSERT INTO orders (order_id, buyer_id, item_id, order_date)
-- VALUES (?, ?, ?, NOW());

-- 3) 更新商品状态
-- UPDATE item SET status = 1 WHERE item_id = ?;

COMMIT;

-- 如任一步失败则 ROLLBACK;
