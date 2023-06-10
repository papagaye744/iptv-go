package handler

import (
  "Live/liveurls"
  "fmt"
  "net/http"
  "strings"
  "log"
)

func defaultQuery(r *http.Request, name string, defaultValue string) string {
  param := r.URL.Query().Get(name)
  if param == "" {
    return defaultValue
  }
  return param
}

func duanyan(adurl string, realurl any) string {
	var liveurl string
	if str, ok := realurl.(string); ok {
		liveurl = str
	} else {
		liveurl = adurl
	}
  log.Println("Redirect url:", liveurl)
	return liveurl
}

// vercel 平台会将请求传递给该函数，这个函数名随意，但函数参数必须按照该规则。
func Handler(w http.ResponseWriter, r *http.Request)  {
  adurl := "http://159.75.85.63:5680/d/ad/roomad/playlist.m3u8"
  path := r.URL.Path
  params := strings.Split(path, "/")
  // fmt.Fprintf(w, "request url: %s", path)
  if len(params) >= 4 {
    // 解析成功
    // 直播平台
    platform := params[2]
    // 房间号
    rid := params[3]
    // fmt.Fprintf(w, "parsed platform=%s, room=%s", platform, rid)
    switch platform {
      case "douyin":
        // 抖音
        douyinobj := &liveurls.Douyin{}
        douyinobj.Rid = rid
        http.Redirect(w, r, duanyan(adurl, douyinobj.GetDouYinUrl()), http.StatusMovedPermanently)
      case "douyu":
        // 斗鱼
        douyuobj := &liveurls.Douyu{}
        douyuobj.Rid = rid
        douyuobj.Stream_type = defaultQuery(r, "stream", "hls")
        douyuobj.Cdn_type = defaultQuery(r, "cdn", "openhls-tct")
        http.Redirect(w, r, duanyan(adurl, douyuobj.GetRealUrl()), http.StatusMovedPermanently)
      case "huya":
        // 虎牙
        huyaobj := &liveurls.Huya{}
        huyaobj.Rid = rid
        huyaobj.Cdn = defaultQuery(r, "cdn", "hwcdn")
        huyaobj.Media = defaultQuery(r, "media", "flv")
        huyaobj.Type = defaultQuery(r, "type", "nodisplay")
        if huyaobj.Type == "display" {
          fmt.Fprintf(w, huyaobj.GetLiveUrl().(string))
        } else {
          http.Redirect(w, r, duanyan(adurl, huyaobj.GetLiveUrl()), http.StatusMovedPermanently)
        }
      case "bilibili":
        // B站
        biliobj := &liveurls.BiliBili{}
        biliobj.Rid = rid
        biliobj.Platform = defaultQuery(r, "platform", "web")
        biliobj.Quality = defaultQuery(r, "quality", "10000")
        biliobj.Line = defaultQuery(r, "line", "second")
        http.Redirect(w, r, duanyan(adurl, biliobj.GetPlayUrl()), http.StatusMovedPermanently)
      case "youtube":
        // 油管
        ytbObj := &liveurls.Youtube{}
        ytbObj.Rid = rid
        ytbObj.Quality = defaultQuery(r, "quality", "1080")
        http.Redirect(w, r, duanyan(adurl, ytbObj.GetLiveUrl()), http.StatusMovedPermanently)
      case "yy":
        // YY直播
        yyObj := &liveurls.Yy{}
        yyObj.Rid = rid
        yyObj.Quality = defaultQuery(r, "quality", "4")
        http.Redirect(w, r, duanyan(adurl, yyObj.GetLiveUrl()), http.StatusMovedPermanently)
      default:
        fmt.Fprintf(w, "Unknown platform=%s, room=%s", platform, rid)
    }
  } else {
    log.Println("Invalid path:", path)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  }
}
