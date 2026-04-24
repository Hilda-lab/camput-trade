package main

import (
	"log"
	"net/http"

	"campus-trade/internal/config"
	"campus-trade/internal/db"
	"campus-trade/internal/handlers"

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

	// Data operations required by the assignment.
	r.POST("/items", h.CreateItem)
	r.POST("/items/:id/price", h.UpdateItemPrice)
	r.POST("/items/:id/delete", h.DeleteUnsoldItem)
	r.POST("/items/manual/price", h.UpdateItemPrice)
	r.POST("/items/manual/delete", h.DeleteUnsoldItem)
	r.POST("/purchase", h.Purchase)

	log.Printf("server listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
