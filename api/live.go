package api

import (
	"fmt"
  "net/http"
  "strings"
  "log"
)

// vercel 平台会将请求传递给该函数，这个函数名随意，但函数参数必须按照该规则。
func Handler(w http.ResponseWriter, r *http.Request)  {
  path := r.URL.Path
  params := strings.Split(path, "/")
  if len(params) >= 3 {
    platform := params[1]
    rid := params[2]
    fmt.Sprintf("platform=%s, room=%s", platform, rid)
  } else {
    log.Println("Invalid path:", path)
    fmt.Sprintf("Invalid path: %s", path)
  }
}
