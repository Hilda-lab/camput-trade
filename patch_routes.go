package main

import (
"os"
"strings"
)

func main() {
content, _ := os.ReadFile("cmd/server/main.go")
str := string(content)
str = strings.Replace(str, "r.POST(\"/login\", h.Login)", "r.POST(\"/login\", h.Login)\n\t\tr.GET(\"/register\", h.RegisterForm)\n\t\tr.POST(\"/register\", h.Register)", 1)
os.WriteFile("cmd/server/main.go", []byte(str), 0644)
}
