package main

import (
"html/template"
"log"
)

func main() {
_, err := template.ParseFiles("templates/home.html")
if err != nil {
log.Fatal(err)
}
}
