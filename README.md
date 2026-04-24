# Campus Trade (Go + Gin + PostgreSQL)

校园二手交易平台数据库课程作业骨架项目。

## 技术栈
- Go + Gin
- PostgreSQL
- HTML Template + Bootstrap
- Render 部署

## 目录结构
- `cmd/server/main.go`: 启动入口
- `internal/config`: 配置加载
- `internal/db`: 数据库连接
- `internal/handlers`: 页面与接口处理
- `internal/service`: 业务逻辑（购买事务）
- `templates`: 页面模板（首页/商品/用户/订单）
- `sql`: 建表、初始数据、视图、查询脚本

## 本地运行
1. 复制环境变量
   - 将 `.env.example` 内容写入系统环境变量，或手动设置 `DATABASE_URL` 和 `PORT`
2. 初始化数据库
   - 依次执行：
     - `sql/001_schema.sql`
     - `sql/002_seed.sql`
     - `sql/003_views.sql`
3. 安装依赖并启动
   - `go mod tidy`
   - `go run ./cmd/server`
4. 打开页面
   - `/` 首页
   - `/items` 商品列表
   - `/users` 用户列表
   - `/orders` 订单列表
   - `/reports` 查询与统计（连接查询/聚合分组/视图）

## 已包含的作业能力
- 三表与主外键约束
- status 约束和订单唯一商品约束（orders.item_id UNIQUE）
- 视图：已售商品视图、未售商品视图
- 购买逻辑：事务 + 行锁（已售不可重复购买）
- 商品页面支持基础查询筛选
- 查询统计页面支持：连接查询、聚合分组、视图结果展示

## Render 部署
1. 推送代码到 GitHub
2. Render 新建 Web Service，连接仓库
3. 使用 `render.yaml` 自动配置
4. 在 Render 设置环境变量：
   - `DATABASE_URL`（指向 Render Postgres）
   - `PORT`（默认可用 10000）
5. 部署完成后获得可访问公网 URL

## 安全性与并发恢复说明（用于报告）
- 安全性：
  - 普通用户仅授予 `SELECT` 权限
  - 删除/修改操作只保留给管理员角色
- 并发：
  - 两人同时购买同一商品会产生竞争
  - 通过事务 + `FOR UPDATE` 行锁 + `orders.item_id` 唯一约束防止重复下单
- 恢复：
  - 依赖 PostgreSQL WAL 与定时备份
  - 故障后可通过备份 + 日志回放恢复订单数据
