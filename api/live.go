// Package api
package api
import (
 "liveurls"
  "encoding/json"
	"fmt"
  "net/http"
)

// vercel 平台会将后面的请求传递给该 Tool 函数，这个函数名也随意，但函数参数必须按照该规则。
func Tool(w http.ResponseWriter, r *http.Request)  {
    // todo
}
