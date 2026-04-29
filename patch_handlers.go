package main

import (
"os"
"strings"
)

func main() {
content, _ := os.ReadFile("internal/handlers/handlers.go")
str := string(content)

oldBlock := `func (h *Handler) LoginForm(c *gin.Context) {
if h.db == nil {
c.String(http.StatusBadRequest, "database not connected")
return
}

var users []map[string]any
h.queryRows(&users, "SELECT user_id, user_name FROM app_user ORDER BY user_id")

c.HTML(http.StatusOK, "login.html", gin.H{
"title": "”√ªßµ«¬º",
"users": users,
})
}

func (h *Handler) Login(c *gin.Context) {
userID := c.PostForm("user_id")
if userID == "" {
c.String(http.StatusBadRequest, "missing user_id")
return
}

var userName string
err := h.db.QueryRowContext(context.Background(), "SELECT user_name FROM app_user WHERE user_id = ?", userID).Scan(&userName)
if err != nil {
c.String(http.StatusBadRequest, "invalid user")
return
}

session := sessions.Default(c)
session.Set("user_id", userID)
session.Set("user_name", userName)
session.Save()

c.Redirect(http.StatusFound, "/")
}`

newBlock := `func (h *Handler) LoginForm(c *gin.Context) {
c.HTML(http.StatusOK, "login.html", gin.H{
"title": "”√ªßµ«¬º",
})
}

func (h *Handler) Login(c *gin.Context) {
email := c.PostForm("email")
password := c.PostForm("password")
if email == "" || password == "" {
c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "’À∫≈ªÚ√‹¬Î≤ªƒ‹Œ™ø’", "title": "µ«¬º"})
return
}

var userID, userName string
// Check auth
err := h.db.QueryRowContext(context.Background(), "SELECT user_id, user_name FROM app_user WHERE (email = ? OR user_id = ?) AND password = ?", email, email, password).Scan(&userID, &userName)
if err != nil {
c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "’À∫≈ªÚ√‹¬Î¥ÌŒÛ", "title": "µ«¬º"})
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
"title": "”√ªß◊¢≤·",
})
}

func (h *Handler) Register(c *gin.Context) {
userName := c.PostForm("user_name")
email := c.PostForm("email")
password := c.PostForm("password")

if userName == "" || email == "" || password == "" {
c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "«ÎÃÓ–¥À˘”–±ÿÃÓ◊÷∂Œ", "title": "◊¢≤·"})
return
}

// Basic check
var count int
h.db.QueryRow("SELECT COUNT(*) FROM app_user WHERE email = ?", email).Scan(&count)
if count > 0 {
c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "∏√” œ‰“—±ª◊¢≤·", "title": "◊¢≤·"})
return
}

// Generate a user string ID
userID := fmt.Sprintf("u%d", time.Now().Unix())

_, err := h.db.Exec("INSERT INTO app_user (user_id, user_name, email, password) VALUES (?, ?, ?, ?)", userID, userName, email, password)
if err != nil {
c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "◊¢≤· ß∞‹£¨«Î…‘∫Û‘Ÿ ‘", "title": "◊¢≤·"})
return
}

// Auto login
session := sessions.Default(c)
session.Set("user_id", userID)
session.Set("user_name", userName)
session.Save()

c.Redirect(http.StatusFound, "/")
}`

if !strings.Contains(str, oldBlock) {
println("Error: Block not found")
}
str = strings.Replace(str, oldBlock, newBlock, 1)

// Since we import fmt and time it might already be imported or need adding
// Let's rely on goimports or the existing imports (fmt, time are already imported based on read_file tool results)
os.WriteFile("internal/handlers/handlers.go", []byte(str), 0644)
}
