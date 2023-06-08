package api

import (
  "encoding/json"
	"fmt"
  "net/http"
  "strings"
  "log"
)

type Response struct {
  Code int `json:"code"`
  Message string `json:"message"`
}

// 返回json数据
func returnJson(w http.ResponseWriter, code int, message string) {
  resp := Response {
    Code: code,
    Message: message
  }
  json, err := json.Marshal(resp)
  if err != nil {
    log.Println("Failed to encode JSON:", err)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(json)
}

// vercel 平台会将请求传递给该函数，这个函数名随意，但函数参数必须按照该规则。
func Handler(w http.ResponseWriter, r *http.Request)  {
  path := r.URL.Path
  params := strings.Split(path, "/")
  if len(params) >= 3 {
    platform := params[1]
    rid := params[2]
    returnJson(w, 200, fmt.Sprintf("platform=%s, room=%s", platform, rid))
  } else {
    log.Println("Invalid path:", path)
    returnJson(w, 500, path)
  }
}
