package handler

import (
  "Live/liveurls"
  "Live/util"
  "fmt"
  "net/http"
  "strings"
  "log"
)

// vercel 平台会将请求传递给该函数，这个函数名随意，但函数参数必须按照该规则。
func Handler(w http.ResponseWriter, r *http.Request)  {
  adurl := "http://159.75.85.63:5680/d/ad/roomad/playlist.m3u8"
  path := r.URL.Path
  params := strings.Split(path, "/")

  // fmt.Fprintf(w, "request url: %s", path)

  if len(params) == 3 {
    // 平台
    platform := params[2]
    switch platform {
      case "douyin":
        // 处理抖音手机直播间
        vrurl := r.URL.Query().Get("url")
		    douyinobj := &liveurls.Douyin{}
		    douyinobj.Shorturl = vrurl
        http.Redirect(w, r, util.Duanyan(adurl, douyinobj.GetRealurl()), http.StatusMovedPermanently)
        return
    }
  }

  if len(params) >= 4 {
    // 解析成功
    // 平台
    platform := params[2]
    // 房间号
    rid := params[3]
    // fmt.Fprintf(w, "parsed platform=%s, room=%s", platform, rid)
    switch platform {
      case "douyin":
        // 抖音
        douyinobj := &liveurls.Douyin{}
        douyinobj.Rid = rid
        http.Redirect(w, r, util.Duanyan(adurl, douyinobj.GetDouYinUrl()), http.StatusMovedPermanently)
      case "douyu":
        // 斗鱼
        douyuobj := &liveurls.Douyu{}
        douyuobj.Rid = rid
        douyuobj.Stream_type = util.DefaultQuery(r, "stream", "flv")
        http.Redirect(w, r, util.Duanyan(adurl, douyuobj.GetRealUrl()), http.StatusMovedPermanently)
      case "huya":
        // 虎牙
        huyaobj := &liveurls.Huya{}
        huyaobj.Rid = rid
        huyaobj.Cdn = util.DefaultQuery(r, "cdn", "hwcdn")
        huyaobj.Media = util.DefaultQuery(r, "media", "flv")
        huyaobj.Type = util.DefaultQuery(r, "type", "nodisplay")
        if huyaobj.Type == "display" {
          fmt.Fprintf(w, huyaobj.GetLiveUrl().(string))
        } else {
          http.Redirect(w, r, util.Duanyan(adurl, huyaobj.GetLiveUrl()), http.StatusMovedPermanently)
        }
      case "bilibili":
        // B站
        biliobj := &liveurls.BiliBili{}
        biliobj.Rid = rid
        biliobj.Platform = util.DefaultQuery(r, "platform", "web")
        biliobj.Quality = util.DefaultQuery(r, "quality", "10000")
        biliobj.Line = util.DefaultQuery(r, "line", "first")
        http.Redirect(w, r, util.Duanyan(adurl, biliobj.GetPlayUrl()), http.StatusMovedPermanently)
      case "youtube":
        // 油管
        ytbObj := &liveurls.Youtube{}
        ytbObj.Rid = rid
        ytbObj.Quality = util.DefaultQuery(r, "quality", "1080")
        http.Redirect(w, r, util.Duanyan(adurl, ytbObj.GetLiveUrl()), http.StatusMovedPermanently)
      case "yy":
        // YY直播
        yyObj := &liveurls.Yy{}
        yyObj.Rid = rid
        yyObj.Quality = util.DefaultQuery(r, "quality", "4")
        http.Redirect(w, r, util.Duanyan(adurl, yyObj.GetLiveUrl()), http.StatusMovedPermanently)
      default:
        fmt.Fprintf(w, "Unknown platform=%s, room=%s", platform, rid)
    }
  } else {
    log.Println("Invalid path:", path)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  }
}
