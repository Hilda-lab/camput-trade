package main

import (
	"log"
	"net/http"

	"campus-trade/internal/config"
	"campus-trade/internal/db"
	"campus-trade/internal/handlers"

        "github.com/gin-contrib/sessions"
        "github.com/gin-contrib/sessions/cookie"
        "github.com/gin-gonic/gin"
)

func main() {
        cfg := config.Load()

        pool, err := db.Connect(cfg.DatabaseURL)
        if err != nil {
                log.Printf("database not connected: %v", err)
        }
        if pool != nil {
                defer pool.Close()
        }

        r := gin.Default()
        
        // Use cookie-based session
        store := cookie.NewStore([]byte("secret"))
        r.Use(sessions.Sessions("campus_session", store))

        r.Static("/static", "./static")
        r.LoadHTMLGlob("templates/*.html")

        h := handlers.New(pool)

        r.GET("/health", func(c *gin.Context) {
                c.JSON(http.StatusOK, gin.H{"status": "ok"})
        })
        
        r.GET("/", h.Home)
        r.GET("/items", h.Items)
        r.GET("/users", h.Users)
        r.GET("/orders", h.Orders)
        r.GET("/reports", h.Reports)

        // Login routes
        r.GET("/login", h.LoginForm)
        r.POST("/login", h.Login)
		r.GET("/register", h.RegisterForm)
		r.POST("/register", h.Register)
        r.GET("/logout", h.Logout)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
