package main

import (
"os"
"path/filepath"
"regexp"
"strings"
)

func main() {
files, _ := filepath.Glob("templates/*.html")
rxDefine := regexp.MustCompile(`(?s)^\{\{define\s+"[^"]+"\}\}\s*`)
rxEnd := regexp.MustCompile(`(?s)\{\{end\}\}\s*$`)

for _, f := range files {
if filepath.Base(f) == "base.html" || filepath.Base(f) == "login.html" {
continue // skip base, and we'll manually rewrite login
}
content, _ := os.ReadFile(f)
s := string(content)
s = rxDefine.ReplaceAllString(s, "")
s = rxEnd.ReplaceAllString(s, "")
os.WriteFile(f, []byte(strings.TrimSpace(s)), 0644)
}
}
