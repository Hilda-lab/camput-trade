-- 基本查询
-- 1) 所有未售出商品
SELECT * FROM item WHERE status = 0;

-- 2) 价格大于 30 的商品
SELECT * FROM item WHERE price > 30;

-- 3) 生活用品类商品
SELECT * FROM item WHERE category = '生活用品';

-- 4) u001 发布的商品
SELECT * FROM item WHERE seller_id = 'u001';

-- 连接查询
-- 1) 已售商品及买家姓名
SELECT i.item_name, u.user_name AS buyer_name
FROM orders o
JOIN item i ON o.item_id = i.item_id
JOIN app_user u ON o.buyer_id = u.user_id;

-- 2) 每个订单：商品名 + 买家名 + 日期
SELECT o.order_id, i.item_name, u.user_name AS buyer_name, o.order_date
FROM orders o
JOIN item i ON o.item_id = i.item_id
JOIN app_user u ON o.buyer_id = u.user_id
ORDER BY o.order_date DESC;

-- 3) 卖家是 u001 的商品是否被购买
SELECT i.item_id, i.item_name,
       CASE WHEN o.item_id IS NULL THEN '未购买' ELSE '已购买' END AS purchase_status
FROM item i
LEFT JOIN orders o ON i.item_id = o.item_id
WHERE i.seller_id = 'u001';

-- 聚合与分组
-- 1) 商品总数
SELECT COUNT(*) AS total_items FROM item;

-- 2) 每类商品数量
SELECT category, COUNT(*) AS item_count
FROM item
GROUP BY category
ORDER BY item_count DESC;

-- 3) 平均价格
SELECT AVG(price) AS avg_price FROM item;

-- 4) 发布商品数量最多的用户
SELECT u.user_id, u.user_name, COUNT(i.item_id) AS published_count
FROM app_user u
JOIN item i ON u.user_id = i.seller_id
GROUP BY u.user_id, u.user_name
ORDER BY published_count DESC
LIMIT 1;
