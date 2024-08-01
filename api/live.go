package handler

import (
  "Golang/liveurls"
  "Golang/utils"
  "fmt"
  "net/http"
  "strings"
  "log"
  "os"
)

// vercel 平台会将请求传递给该函数，这个函数名随意，但函数参数必须按照该规则。
func Handler(w http.ResponseWriter, r *http.Request)  {
  adurl := "https://cdn.jsdelivr.net/gh/feiyangdigital/testvideo/sdr1080pvideo/index.m3u8"
  path := r.URL.Path
  params := strings.Split(path, "/")

  // 是否禁用TV
  enableTV := os.Getenv("TV") != "false" 

  // fmt.Fprintf(w, "request url: %s", path)

  if len(params) >= 4 {
    // 解析成功
    // 平台
    platform := params[2]
    // 房间号
    rid := params[3]
    ts := utils.DefaultQuery(r, "ts", "")
    // fmt.Fprintf(w, "parsed platform=%s, room=%s", platform, rid)
    switch platform {
      case "itv":
        if enableTV {
          itvobj := &liveurls.Itv{}
          cdn := utils.DefaultQuery(r, "cdn", "")
          if ts == "" {
            itvobj.HandleMainRequest(w, r, cdn, rid)
          } else {
            itvobj.HandleTsRequest(w, ts)
          }
        } else {
          http.Error(w, "公共服务不提供TV直播", http.StatusForbidden)
        }
      case "ysptp":
        if enableTV {
          ysptpobj := &liveurls.Ysptp{}
          if ts == "" {
            ysptpobj.HandleMainRequest(w, r, rid)
          } else {
            ysptpobj.HandleTsRequest(w, ts, utils.DefaultQuery(r, "wsTime", ""))
          }
        } else {
          http.Error(w, "公共服务不提供TV直播", http.StatusForbidden)
        }
      case "douyin":
        // 抖音
        douyinobj := &liveurls.Douyin{}
        douyinobj.Rid = rid
        douyinobj.Stream = utils.DefaultQuery(r, "stream", "flv")
        http.Redirect(w, r, utils.Duanyan(adurl, douyinobj.GetDouYinUrl()), http.StatusMovedPermanently)
      case "douyu":
        // 斗鱼
        douyuobj := &liveurls.Douyu{}
        douyuobj.Rid = rid
        douyuobj.Stream_type = utils.DefaultQuery(r, "stream", "flv")
        http.Redirect(w, r, utils.Duanyan(adurl, douyuobj.GetRealUrl()), http.StatusMovedPermanently)
      case "huya":
        // 虎牙
        huyaobj := &liveurls.Huya{}
        huyaobj.Rid = rid
        huyaobj.Cdn = utils.DefaultQuery(r, "cdn", "hwcdn")
        huyaobj.Media = utils.DefaultQuery(r, "media", "flv")
        huyaobj.Type = utils.DefaultQuery(r, "cdntype", "nodisplay")
        if huyaobj.Type == "display" {
          fmt.Fprintf(w, huyaobj.GetLiveUrl().(string))
        } else {
          http.Redirect(w, r, utils.Duanyan(adurl, huyaobj.GetLiveUrl()), http.StatusMovedPermanently)
        }
      case "bilibili":
        // B站
        biliobj := &liveurls.BiliBili{}
        biliobj.Rid = rid
        biliobj.Platform = utils.DefaultQuery(r, "platform", "web")
        biliobj.Quality = utils.DefaultQuery(r, "quality", "10000")
        biliobj.Line = utils.DefaultQuery(r, "line", "first")
        http.Redirect(w, r, utils.Duanyan(adurl, biliobj.GetPlayUrl()), http.StatusMovedPermanently)
      case "youtube":
        // 油管
        ytbObj := &liveurls.Youtube{}
        ytbObj.Rid = rid
        ytbObj.Quality = utils.DefaultQuery(r, "quality", "1080")
        http.Redirect(w, r, utils.Duanyan(adurl, ytbObj.GetLiveUrl()), http.StatusMovedPermanently)
      case "yy":
        // YY直播
        yyObj := &liveurls.Yy{}
        yyObj.Rid = rid
        yyObj.Quality = utils.DefaultQuery(r, "quality", "4")
        http.Redirect(w, r, utils.Duanyan(adurl, yyObj.GetLiveUrl()), http.StatusMovedPermanently)
      default:
        fmt.Fprintf(w, "Unknown platform=%s, room=%s", platform, rid)
    }
  } else {
    log.Println("Invalid path:", path)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  }
}
