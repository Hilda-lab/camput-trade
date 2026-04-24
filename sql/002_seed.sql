INSERT INTO app_user (user_id, user_name, email) VALUES
('u001', '张三', 'zhangsan@example.com'),
('u002', '李四', 'lisi@example.com'),
('u003', '王五', 'wangwu@example.com'),
('u004', '赵六', 'zhaoliu@example.com')
ON CONFLICT (user_id) DO NOTHING;

INSERT INTO item (item_id, item_name, category, price, seller_id, status) VALUES
('i001', '高等数学教材', '学习用品', 25.00, 'u001', 0),
('i002', '二手台灯', '生活用品', 32.00, 'u001', 0),
('i003', '羽毛球拍', '体育用品', 48.00, 'u002', 0),
('i004', '马克杯', '生活用品', 15.00, 'u003', 1)
ON CONFLICT (item_id) DO NOTHING;

INSERT INTO orders (order_id, buyer_id, item_id, order_date) VALUES
('o001', 'u004', 'i004', NOW() - INTERVAL '1 day')
ON CONFLICT (order_id) DO NOTHING;
