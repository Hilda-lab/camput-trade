package main

import (
"log"
"os"
"path/filepath"
"regexp"
)

func main() {
navStr := `    <div class="collapse navbar-collapse">
      <div class="navbar-nav me-auto">
        <a class="nav-link" href="/items">商品</a>
        <a class="nav-link" href="/users">用户</a>
        <a class="nav-link" href="/orders">订单</a>
        <a class="nav-link" href="/reports">查询统计</a>
      </div>
      <div class="navbar-nav ms-auto">
        {{if .loggedIn}}
        <span class="navbar-text text-light me-3">
          当前登录: {{.user_name}} ({{.user_id}})
        </span>
        <a class="nav-link btn btn-sm btn-danger text-white px-3" style="height: 32px; padding-top:2px" href="/logout">登出</a>
        {{else}}
        <a class="nav-link btn btn-sm btn-success text-white px-3" style="height: 32px; padding-top:2px" href="/login">登录身份</a>
        {{end}}
      </div>
    </div>
  </div>
</nav>`

rx := regexp.MustCompile(`(?s)<div class="navbar-nav[^>]*>.*?</nav>`)

files, err := filepath.Glob("templates/*.html")
if err != nil {
log.Fatal(err)
}

for _, f := range files {
content, err := os.ReadFile(f)
if err != nil {
log.Fatal(err)
}
newContent := rx.ReplaceAllString(string(content), navStr)
if newContent != string(content) {
os.WriteFile(f, []byte(newContent), 0644)
log.Printf("Updated %s", f)
}
}
}
