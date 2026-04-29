package config

import "os"

type Config struct {
Port        string
DatabaseURL string
}

func Load() Config {
port := os.Getenv("PORT")
if port == "" {
port = "8080"
}

dbURL := os.Getenv("DATABASE_URL")
if dbURL == "" {
dbURL = "root:123456@tcp(127.0.0.1:3306)/campus_trade?parseTime=true&charset=utf8mb4&loc=Local"
}

return Config{
Port:        port,
DatabaseURL: dbURL,
}
}
