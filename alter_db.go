package main

import (
"database/sql"
"log"

"campus-trade/internal/config"

_ "github.com/go-sql-driver/mysql"
)

func main() {
cfg := config.Load()
db, err := sql.Open("mysql", cfg.DatabaseURL)
if err != nil {
log.Fatal(err)
}
defer db.Close()

// Add password column to app_user
_, err = db.Exec("ALTER TABLE app_user ADD COLUMN password VARCHAR(100) NOT NULL DEFAULT '123456'")
if err != nil {
log.Println("Column might already exist or error:", err)
} else {
log.Println("Successfully added password column to app_user table.")
}
}
