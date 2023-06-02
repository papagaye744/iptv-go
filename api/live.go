// Package api
package api

import (
  "liveurls"
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

func returnJson(w http.ResponseWriter, code int, message string) {
  // 返回数据
  resp := Response {
    Code: code,
    Message: message
  }
  // 格式化为json字符串
  json, err := json.Marshal(resp)
  if err != nil {
    log.Println("Failed to encode JSON:", err)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    return
  }
  // 设置响应头
  w.Header().Set("Content-Type", "application/json")
  // 写入json
  w.Write(json)
}

func defaultQuery(r *http.Request, name string, defaultValue string) string {
  param := r.URL.Query().Get(name)
  if param == "" {
    return defaultValue
  }
  return param
}

// vercel 平台会将请求传递给该函数，这个函数名随意，但函数参数必须按照该规则。
func Live(w http.ResponseWriter, r *http.Request)  {
    path := r.URL.Path
    params := strings.Split(path, "/")
    if len(params) >= 3 {
      // 解析成功
      // 直播平台
      platform := params[1]
      // 房间号
      rid := params[2]
      switch platform {
      case "douyin":
        douyinobj := &liveurls.Douyin{}
        douyinobj.Rid = rid
        http.Redirect(w, r, http.StatusMovedPermanently, duanyan(adurl, douyinobj.GetDouYinUrl()))
      case "douyu":
        douyuobj := &liveurls.Douyu{}
        douyuobj.Rid = rid
        douyuobj.Stream_type = defaultQuery(r, "stream", "hls")
        douyuobj.Cdn_type = defaultQuery(r, "cdn", "akm-tct")
        http.Redirect(http.StatusMovedPermanently, duanyan(adurl, douyuobj.GetRealUrl()))
      case "huya":
        huyaobj := &liveurls.Huya{}
        huyaobj.Rid = rid
        huyaobj.Cdn = defaultQuery(r, "cdn", "hwcdn")
        huyaobj.Media = defaultQuery(r, "media", "flv")
        huyaobj.Type = defaultQuery(r, "type", "nodisplay")
        if huyaobj.Type == "display" {
          returnJson(w, 200, huyaobj.GetLiveUrl())
        } else {
          http.Redirect(http.StatusMovedPermanently, duanyan(adurl, huyaobj.GetLiveUrl()))
        }
      case "bilibili":
        biliobj := &liveurls.BiliBili{}
        biliobj.Rid = rid
        biliobj.Platform = defaultQuery(r, "platform", "web")
        biliobj.Quality = defaultQuery(r, "quality", "10000")
        biliobj.Line = defaultQuery(r, "line", "second")
        http.Redirect(http.StatusMovedPermanently, duanyan(adurl, biliobj.GetPlayUrl()))
      case "youtube":
        ytbObj := &liveurls.Youtube{}
        ytbObj.Rid = rid
        ytbObj.Quality = defaultQuery(r, "quality", "1080")
        http.Redirect(http.StatusMovedPermanently, duanyan(adurl, ytbObj.GetLiveUrl()))
      }
    } else {
      log.Println("Invalid path:", path)
      returnJson(w, 500, path)
    }
}
