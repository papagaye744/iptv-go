package handler
 
import (
  "fmt"
  "net/http"
)

// 入口函数
func Handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
}