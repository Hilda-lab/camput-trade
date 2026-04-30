package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"campus-trade/internal/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{db: db}
}

// getCurrentUser returns user_id, user_name, true if logged in, else "", "", false
func (h *Handler) getCurrentUser(c *gin.Context) (string, string, bool) {
	session := sessions.Default(c)
	uid := session.Get("user_id")
	uname := session.Get("user_name")
	if uid != nil && uname != nil {
		return uid.(string), uname.(string), true
	}
	return "", "", false
}

func (h *Handler) LoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "用户登录",
	})
}

func (h *Handler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	if email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "账号或密码不能为空", "title": "登录"})
		return
	}

	var userID, userName string
	err := h.db.QueryRowContext(context.Background(), "SELECT user_id, user_name FROM app_user WHERE (email = ? OR user_id = ?) AND password = ?", email, email, password).Scan(&userID, &userName)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "账号或密码错误", "title": "登录"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", userID)
	session.Set("user_name", userName)
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

func (h *Handler) RegisterForm(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "用户注册",
	})
}

func (h *Handler) Register(c *gin.Context) {
	userName := strings.TrimSpace(c.PostForm("user_name"))
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")

	if userName == "" || email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "请填写所有必填字段", "title": "注册"})
		return
	}

	var count int
	h.db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM app_user WHERE email = ?", email).Scan(&count)
	if count > 0 {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "该邮箱已被注册", "title": "注册"})
		return
	}

	userID := fmt.Sprintf("u%d", time.Now().Unix())
	
	_, err := h.db.ExecContext(context.Background(), "INSERT INTO app_user (user_id, user_name, email, password) VALUES (?, ?, ?, ?)", userID, userName, email, password)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "注册失败，请稍后再试", "title": "注册"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", userID)
	session.Set("user_name", userName)
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

func (h *Handler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

// baseContext returns a gin.H with common parameters for all views
func (h *Handler) baseContext(c *gin.Context, title string) gin.H {
	uid, uname, loggedIn := h.getCurrentUser(c)
	return gin.H{
		"title":     title,
		"user_id":   uid,
		"user_name": uname,
		"loggedIn":  loggedIn,
	}
}

func (h *Handler) Home(c *gin.Context) {
	c.Redirect(http.StatusFound, "/items")
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

	ctx := h.baseContext(c, "查询与统计")
	ctx["rows"] = rows
	ctx["message"] = message
	ctx["queryType"] = q
	ctx["selectedTitle"] = selected.Title

	c.HTML(http.StatusOK, "reports.html", ctx)
}

func (h *Handler) Users(c *gin.Context) {
	rows := []map[string]any{}
	message := h.queryRows(&rows, "SELECT user_id, user_name, email FROM app_user ORDER BY user_id")
	
	ctx := h.baseContext(c, "用户列表")
	ctx["rows"] = rows
	ctx["message"] = message

	c.HTML(http.StatusOK, "users.html", ctx)
}

func (h *Handler) Items(c *gin.Context) {
	queryType := c.Query("q")
	queryMinPrice := c.Query("min_price")
	queryMaxPrice := c.Query("max_price")
	queryCategory := c.Query("category")
	querySeller := c.Query("seller")

	qContext := "SELECT item_id, item_name, category, price, seller_id, status FROM item"
	var conditions []string
	var args []any

	if queryType == "unsold" {
		conditions = append(conditions, "status = 0")
	}

	if queryMinPrice != "" {
		conditions = append(conditions, "price >= ?")
		args = append(args, queryMinPrice)
	}

	if queryMaxPrice != "" {
		conditions = append(conditions, "price <= ?")
		args = append(args, queryMaxPrice)
	}

	if queryCategory != "" {
		conditions = append(conditions, "category LIKE ?")
		args = append(args, "%"+queryCategory+"%")
	}

	if querySeller != "" {
		conditions = append(conditions, "seller_id = ?")
		args = append(args, querySeller)
	}

	if len(conditions) > 0 {
		qContext += " WHERE " + strings.Join(conditions, " AND ")
	}
	qContext += " ORDER BY item_id DESC"

	rows := []map[string]any{}
	
	var err error
	var sqlRows *sql.Rows

	if h.db != nil {
		sqlRows, err = h.db.QueryContext(c.Request.Context(), qContext, args...)
		if err == nil {
			defer sqlRows.Close()
			cols, _ := sqlRows.Columns()
			for sqlRows.Next() {
				colsMap := make([]any, len(cols))
				colsVal := make([]any, len(cols))
				for i := range colsMap {
					colsVal[i] = &colsMap[i]
				}
				sqlRows.Scan(colsVal...)
	
				rowMap := make(map[string]any)
				for i, col := range cols {
					val := colsMap[i]
					if b, ok := val.([]byte); ok {
						rowMap[col] = string(b)
					} else {
						rowMap[col] = val
					}
				}
				rows = append(rows, rowMap)
			}
		}
	}

	var message string
	if err != nil {
		message = "query failed: " + err.Error()
	}

	ctx := h.baseContext(c, "商品列表")
	ctx["rows"] = rows
	ctx["message"] = message
	ctx["queryType"] = queryType
	ctx["qMinPrice"] = queryMinPrice
	ctx["qMaxPrice"] = queryMaxPrice
	ctx["qCategory"] = queryCategory
	ctx["qSeller"] = querySeller

	c.HTML(http.StatusOK, "items.html", ctx)
}

func (h *Handler) Orders(c *gin.Context) {
	rows := []map[string]any{}
	message := h.queryRows(&rows, "SELECT o.order_id, o.item_id, i.item_name, o.buyer_id, u.user_name AS buyer_name, o.order_date FROM orders o JOIN item i ON o.item_id = i.item_id JOIN app_user u ON o.buyer_id = u.user_id ORDER BY o.order_date DESC")
	
	ctx := h.baseContext(c, "订单列表")
	ctx["rows"] = rows
	ctx["message"] = message

	c.HTML(http.StatusOK, "orders.html", ctx)
}

func (h *Handler) CreateItem(c *gin.Context) {
	if h.db == nil {
		c.String(http.StatusBadRequest, "database not connected")
		return
	}

	uid, _, loggedIn := h.getCurrentUser(c)
	if !loggedIn {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id := strings.TrimSpace(c.PostForm("item_id"))
	name := strings.TrimSpace(c.PostForm("item_name"))
	category := strings.TrimSpace(c.PostForm("category"))
	price := c.PostForm("price")

	_, err := h.db.ExecContext(c.Request.Context(), "INSERT INTO item (item_id, item_name, category, price, seller_id, status, created_at) VALUES (?, ?, ?, ?, ?, 0, NOW())", id, name, category, price, uid)
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

	uid, _, loggedIn := h.getCurrentUser(c)
	if !loggedIn {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	itemID := strings.TrimSpace(c.Param("id"))
	if itemID == "manual" || itemID == "" {
		itemID = strings.TrimSpace(c.PostForm("item_id"))
	}
	if itemID == "" {
		c.String(http.StatusBadRequest, "item_id is required")
		return
	}

	var sellerID string
	err := h.db.QueryRowContext(c.Request.Context(), "SELECT seller_id FROM item WHERE item_id = ?", itemID).Scan(&sellerID)
	if err != nil || sellerID != uid {
		c.String(http.StatusForbidden, "unauthorized: you are not the seller of this item")
		return
	}

	price := c.PostForm("price")
	_, err = h.db.ExecContext(c.Request.Context(), "UPDATE item SET price = ? WHERE item_id = ?", price, itemID)
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

	uid, _, loggedIn := h.getCurrentUser(c)
	if !loggedIn {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	itemID := strings.TrimSpace(c.Param("id"))
	if itemID == "manual" || itemID == "" {
		itemID = strings.TrimSpace(c.PostForm("item_id"))
	}
	if itemID == "" {
		c.String(http.StatusBadRequest, "item_id is required")
		return
	}

	var sellerID string
	err := h.db.QueryRowContext(c.Request.Context(), "SELECT seller_id FROM item WHERE item_id = ?", itemID).Scan(&sellerID)
	if err != nil || sellerID != uid {
		c.String(http.StatusForbidden, "unauthorized: you are not the seller of this item")
		return
	}

	_, err = h.db.ExecContext(c.Request.Context(), "DELETE FROM item WHERE item_id = ? AND status = 0", itemID)
	if err != nil {
		c.String(http.StatusBadRequest, "delete unsold item failed: %v", err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/items")
}

func (h *Handler) Purchase(c *gin.Context) {
	uid, _, loggedIn := h.getCurrentUser(c)
	if !loggedIn {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	itemID := c.PostForm("item_id")
	if itemID == "" {
		c.String(http.StatusBadRequest, "item_id is required")
		return
	}

	var sellerID string
	err := h.db.QueryRowContext(c.Request.Context(), "SELECT seller_id FROM item WHERE item_id = ?", itemID).Scan(&sellerID)
	if err != nil {
		c.String(http.StatusBadRequest, "item not found")
		return
	}
	if sellerID == uid {
		c.String(http.StatusForbidden, "you cannot purchase your own item")
		return
	}

	orderID := fmt.Sprintf("o%s", time.Now().Format("20060102150405"))
	
	if err := service.PurchaseItem(h.db, orderID, itemID, uid); err != nil {
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
