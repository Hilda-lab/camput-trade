package main

import (
"html/template"
"log"
)

func main() {
_, err := template.ParseGlob("templates/*.html")
if err != nil {
log.Fatal(err)
}
}
