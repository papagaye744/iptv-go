package util

func GetTestVideoUrl(w http.ResponseWriter) {
  str_time := time.Now().Format("2006-01-02 15:04:05")
  fmt.Fprintln(w, "#EXTM3U")
  fmt.Fprintln(w, "#EXTINF:-1 tvg-name=\""+str_time+"\" tvg-logo=\"https://cdn.jsdelivr.net/gh/youshandefeiyang/IPTV/logo/tg.jpg\" group-title=\"列表更新时间\","+str_time)
  fmt.Fprintln(w, "https://cdn.jsdelivr.net/gh/youshandefeiyang/testvideo/time/time.mp4")
  fmt.Fprintln(w, "#EXTINF:-1 tvg-name=\"4K60PSDR-H264-AAC测试\" tvg-logo=\"https://cdn.jsdelivr.net/gh/youshandefeiyang/IPTV/logo/tg.jpg\" group-title=\"4K频道\",4K60PSDR-H264-AAC测试")
  fmt.Fprintln(w, "http://159.75.85.63:5680/d/ad/h264/playad.m3u8")
  fmt.Fprintln(w, "#EXTINF:-1 tvg-name=\"4K60PHLG-HEVC-EAC3测试\" tvg-logo=\"https://cdn.jsdelivr.net/gh/youshandefeiyang/IPTV/logo/tg.jpg\" group-title=\"4K频道\",4K60PHLG-HEVC-EAC3测试")
  fmt.Fprintln(w, "http://159.75.85.63:5680/d/ad/playad.m3u8")
}

func GetLivePrefix(r *http.Request) string {
  firstUrl := DefaultQuery(r, "url", "https://www.goodiptv.club")
  realUrl, _ := url.QueryUnescape(firstUrl)
  return realUrl
}

func DefaultQuery(r *http.Request, name string, defaultValue string) string {
  param := r.URL.Query().Get(name)
  if param == "" {
    return defaultValue
  }
  return param
}
  
func Duanyan(adurl string, realurl any) string {
  var liveurl string
  if str, ok := realurl.(string); ok {
    liveurl = str
  } else {
	liveurl = adurl
  }
  log.Println("Redirect url:", liveurl)
 return liveurl
}