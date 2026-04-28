package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"campus-trade/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{"title": "校园二手交易平台"})
}

func (h *Handler) Reports(c *gin.Context) {
	q := c.DefaultQuery("q", "sold_with_buyer")

	queries := map[string]struct {
		Title string
		SQL   string
	}{
		"sold_with_buyer": {
			Title: "所有已售商品及其买家姓名",
			SQL:   "SELECT i.item_id, i.item_name, u.user_name AS buyer_name FROM orders o JOIN item i ON o.item_id = i.item_id JOIN app_user u ON o.buyer_id = u.user_id ORDER BY i.item_id",
		},
		"order_full": {
			Title: "每个订单：商品名 + 买家名 + 日期",
			SQL:   "SELECT o.order_id, i.item_name, u.user_name AS buyer_name, o.order_date FROM orders o JOIN item i ON o.item_id = i.item_id JOIN app_user u ON o.buyer_id = u.user_id ORDER BY o.order_date DESC",
		},
		"seller_u001_purchase": {
			Title: "卖家是 u001 的商品是否被购买",
			SQL:   "SELECT i.item_id, i.item_name, CASE WHEN o.item_id IS NULL THEN '未购买' ELSE '已购买' END AS purchase_status FROM item i LEFT JOIN orders o ON i.item_id = o.item_id WHERE i.seller_id = 'u001' ORDER BY i.item_id",
		},
		"count_items": {
			Title: "统计商品总数",
			SQL:   "SELECT COUNT(*) AS total_items FROM item",
		},
		"count_by_category": {
			Title: "统计每类商品数量",
			SQL:   "SELECT category, COUNT(*) AS item_count FROM item GROUP BY category ORDER BY item_count DESC",
		},
		"avg_price": {
			Title: "计算所有商品平均价格",
			SQL:   "SELECT ROUND(AVG(price), 2) AS avg_price FROM item",
		},
		"top_seller": {
			Title: "发布商品数量最多的用户",
			SQL:   "SELECT u.user_id, u.user_name, COUNT(i.item_id) AS published_count FROM app_user u JOIN item i ON u.user_id = i.seller_id GROUP BY u.user_id, u.user_name ORDER BY published_count DESC LIMIT 1",
		},
		"sold_view": {
			Title: "已售商品视图（商品名 + 买家ID）",
			SQL:   "SELECT * FROM sold_items_view",
		},
		"unsold_view": {
			Title: "未售商品视图",
			SQL:   "SELECT * FROM unsold_items_view ORDER BY item_id",
		},
	}

	selected, ok := queries[q]
	if !ok {
		q = "sold_with_buyer"
		selected = queries[q]
	}

	rows := []map[string]any{}
	message := h.queryRows(&rows, selected.SQL)
	c.HTML(http.StatusOK, "reports.html", gin.H{
		"title":         "查询与统计",
		"rows":          rows,
		"message":       message,
		"queryType":     q,
		"selectedTitle": selected.Title,
	})
}

func (h *Handler) Users(c *gin.Context) {
	rows := []map[string]any{}
	message := h.queryRows(&rows, "SELECT user_id, user_name, email FROM app_user ORDER BY user_id")
	c.HTML(http.StatusOK, "users.html", gin.H{"title": "用户列表", "rows": rows, "message": message})
}

func (h *Handler) Items(c *gin.Context) {
	queryType := c.Query("q")
	sql := "SELECT item_id, item_name, category, price, seller_id, status FROM item ORDER BY item_id"

	switch queryType {
	case "unsold":
		sql = "SELECT item_id, item_name, category, price, seller_id, status FROM item WHERE status = 0 ORDER BY item_id"
	case "price_gt_30":
		sql = "SELECT item_id, item_name, category, price, seller_id, status FROM item WHERE price > 30 ORDER BY item_id"
	case "daily":
		sql = "SELECT item_id, item_name, category, price, seller_id, status FROM item WHERE category = '生活用品' ORDER BY item_id"
	case "seller_u001":
		sql = "SELECT item_id, item_name, category, price, seller_id, status FROM item WHERE seller_id = 'u001' ORDER BY item_id"
	}

	rows := []map[string]any{}
	message := h.queryRows(&rows, sql)
	c.HTML(http.StatusOK, "items.html", gin.H{"title": "商品列表", "rows": rows, "message": message, "queryType": queryType})
}

func (h *Handler) Orders(c *gin.Context) {
	rows := []map[string]any{}
	message := h.queryRows(&rows, "SELECT o.order_id, o.item_id, i.item_name, o.buyer_id, u.user_name AS buyer_name, o.order_date FROM orders o JOIN item i ON o.item_id = i.item_id JOIN app_user u ON o.buyer_id = u.user_id ORDER BY o.order_date DESC")
	c.HTML(http.StatusOK, "orders.html", gin.H{"title": "订单列表", "rows": rows, "message": message})
}

func (h *Handler) CreateItem(c *gin.Context) {
	if h.db == nil {
		c.String(http.StatusBadRequest, "database not connected")
		return
	}

	id := strings.TrimSpace(c.PostForm("item_id"))
	name := strings.TrimSpace(c.PostForm("item_name"))
	category := strings.TrimSpace(c.PostForm("category"))
	sellerID := strings.TrimSpace(c.PostForm("seller_id"))
	price := c.PostForm("price")

	_, err := h.db.ExecContext(c.Request.Context(), "INSERT INTO item (item_id, item_name, category, price, seller_id, status, created_at) VALUES (?, ?, ?, ?, ?, 0, NOW())", id, name, category, price, sellerID)
	if err != nil {
		c.String(http.StatusBadRequest, "create item failed: %v", err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/items")
}

func (h *Handler) UpdateItemPrice(c *gin.Context) {
	if h.db == nil {
		c.String(http.StatusBadRequest, "database not connected")
		return
	}
	itemID := strings.TrimSpace(c.Param("id"))
	if itemID == "manual" {
		itemID = strings.TrimSpace(c.PostForm("item_id"))
	}
	if itemID == "" {
		c.String(http.StatusBadRequest, "item_id is required")
		return
	}
	price := c.PostForm("price")
	_, err := h.db.ExecContext(c.Request.Context(), "UPDATE item SET price = ? WHERE item_id = ?", price, itemID)
	if err != nil {
		c.String(http.StatusBadRequest, "update item price failed: %v", err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/items")
}

func (h *Handler) DeleteUnsoldItem(c *gin.Context) {
	if h.db == nil {
		c.String(http.StatusBadRequest, "database not connected")
		return
	}
	itemID := strings.TrimSpace(c.Param("id"))
	if itemID == "manual" {
		itemID = strings.TrimSpace(c.PostForm("item_id"))
	}
	if itemID == "" {
		c.String(http.StatusBadRequest, "item_id is required")
		return
	}
	_, err := h.db.ExecContext(c.Request.Context(), "DELETE FROM item WHERE item_id = ? AND status = 0", itemID)
	if err != nil {
		c.String(http.StatusBadRequest, "delete unsold item failed: %v", err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/items")
}

func (h *Handler) Purchase(c *gin.Context) {
	orderID := fmt.Sprintf("o%s", time.Now().Format("20060102150405"))
	itemID := c.PostForm("item_id")
	buyerID := c.PostForm("buyer_id")

	if err := service.PurchaseItem(h.db, orderID, itemID, buyerID); err != nil {
		c.String(http.StatusBadRequest, "purchase failed: %v", err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/orders")
}

func (h *Handler) queryRows(target *[]map[string]any, sql string) string {
	if h.db == nil {
		return "DATABASE_URL 未配置，当前仅展示页面骨架。"
	}
	rows, err := h.db.QueryContext(context.Background(), sql)
	if err != nil {
		return "query failed: " + err.Error()
	}
	defer rows.Close()

	result, err := rowsToMap(rows)
	if err != nil {
		return "decode rows failed: " + err.Error()
	}
	*target = result
	return ""
}

func rowsToMap(rows *sql.Rows) ([]map[string]any, error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	items := make([]map[string]any, 0)

	values := make([]any, len(fields))
	valuePtrs := make([]any, len(fields))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		line := make(map[string]any, len(fields))
		for i, field := range fields {
			line[field] = normalizeValue(values[i])
		}
		items = append(items, line)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func normalizeValue(value any) any {
	switch v := value.(type) {
	case []byte:
		return string(v)
	case time.Time:
		return v.Format("2006-01-02 15:04:05")
	default:
		return v
	}
}
