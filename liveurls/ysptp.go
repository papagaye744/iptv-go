package liveurls

import (
	"Golang/utils"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Ysptp struct{}

var cache sync.Map

type CacheItem struct {
	Value      string
	Expiration int64
}

var cctvList = map[string]string{
	"cctv1.m3u8":        "http://liveali-tpgq.cctv.cn/live/cctv1.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv2.m3u8":        "http://liveali-tpgq.cctv.cn/live/cctv2.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv3.m3u8":        "http://liveali-tpgq.cctv.cn/live/cctv3.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv4.m3u8":        "http://liveali-tpgq.cctv.cn/live/cctv4.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv5.m3u8":        "http://liveali-tpgq.cctv.cn/live/cctv5.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv5p.m3u8":       "http://liveali-tpgq.cctv.cn/live/cctv5p.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv6.m3u8":        "http://liveali-tpgq.cctv.cn/live/cctv6.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv7.m3u8":        "http://liveali-tpgq.cctv.cn/live/cctv7.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv8.m3u8":        "http://liveali-tpgq.cctv.cn/live/cctv8.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv9.m3u8":        "http://liveali-tpgq.cctv.cn/live/cctv9.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv10.m3u8":       "http://liveali-tpgq.cctv.cn/live/cctv10.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv11.m3u8":       "http://liveali-tpgq.cctv.cn/live/cctv11.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv12.m3u8":       "http://liveali-tpgq.cctv.cn/live/cctv12.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv13.m3u8":       "http://liveali-tpgq.cctv.cn/live/cctv13.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv14.m3u8":       "http://liveali-tpgq.cctv.cn/live/cctv14.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv15.m3u8":       "http://liveali-tpgq.cctv.cn/live/cctv15.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv16.m3u8":       "http://liveali-tpgq.cctv.cn/live/cctv16.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv17.m3u8":       "http://liveali-tpgq.cctv.cn/live/cctv17.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cgtnar.m3u8":       "http://liveali-tpgq.cctv.cn/live/cgtnar.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cgtndoc.m3u8":      "http://liveali-tpgq.cctv.cn/live/cgtndoc.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cgtnen.m3u8":       "http://liveali-tpgq.cctv.cn/live/cgtnen.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cgtnfr.m3u8":       "http://liveali-tpgq.cctv.cn/live/cgtnfr.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cgtnru.m3u8":       "http://liveali-tpgq.cctv.cn/live/cgtnru.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cgtnsp.m3u8":       "http://liveali-tpgq.cctv.cn/live/cgtnsp.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv4k.m3u8":       "http://liveali-tpgq.cctv.cn/live/cctv4k.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv4k_10m.m3u8":   "http://liveali-tpgq.cctv.cn/live/cctv4k10m.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv4k16.m3u8":     "http://liveali-tpgq.cctv.cn/live/cctv4k16.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv4k16_10m.m3u8": "http://liveali-tpgq.cctv.cn/live/cctv4k1610m.m3u8,http://liveali-tpgq.cctv.cn/live/",
	"cctv8k_36m.m3u8":   "http://liveali-tp4k.cctv.cn/live/4K36M/playlist.m3u8,http://liveali-tp4k.cctv.cn/live/4K36M/",
	"cctv8k_120m.m3u8":  "http://liveali-tp4k.cctv.cn/live/8K120M/playlist.m3u8,http://liveali-tp4k.cctv.cn/live/8K120M/",
}

func (y *Ysptp) HandleMainRequest(w http.ResponseWriter, r *http.Request, id string) {
	uid := utils.DefaultQuery(r, "uid", "1234123122")

	if _, ok := cctvList[id]; !ok {
		http.Error(w, "id not found!", http.StatusNotFound)
		return
	}

	urls := strings.Split(cctvList[id], ",")
	data := getURL(id, urls[0], uid, urls[1])
	golang := "http://" + r.Host + r.URL.Path
	re := regexp.MustCompile(`((?i).*?\.ts)`)
	data = re.ReplaceAllString(data, golang+"?ts="+urls[1]+"$1")

	w.Header().Set("Content-Disposition", "attachment;filename="+id)
	w.WriteHeader(http.StatusOK) // Set the status code to 200
    w.Write([]byte(data)) // Write the response body
}

func (y *Ysptp) HandleTsRequest(w http.ResponseWriter, ts, wsTime string) {
	data := ts + "&wsTime=" + wsTime
	w.Header().Set("Content-Type", "video/MP2T")
	w.WriteHeader(http.StatusOK) // Set the status code to 200
    w.Write([]byte(getTs(data))) // Write the response body
}

func getURL(id, url, uid, path string) string {
	cacheKey := id + uid
	if playURL, found := getCache(cacheKey); found {
		return fetchData(playURL, path, uid)
	}

	bstrURL := "https://ytpvdn.cctv.cn/cctvmobileinf/rest/cctv/videoliveUrl/getstream"
	postData := `appcommon={"ap":"cctv_app_tv","an":"央视投屏助手","adid":" ` + uid + `","av":"1.1.7"}&url=` + url

	req, _ := http.NewRequest("POST", bstrURL, strings.NewReader(postData))
	req.Header.Set("User-Agent", "cctv_app_tv")
	req.Header.Set("Referer", "api.cctv.cn")
	req.Header.Set("UID", uid)

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	var body strings.Builder
	_, _ = io.Copy(&body, resp.Body)

	var result map[string]interface{}
	json.Unmarshal([]byte(body.String()), &result)
	playURL := result["url"].(string)

	setCache(cacheKey, playURL)

	return fetchData(playURL, path, uid)
}

func fetchData(playURL, path, uid string) string {
	client := &http.Client{}
	for {
		req, _ := http.NewRequest("GET", playURL, nil)
		req.Header.Set("User-Agent", "cctv_app_tv")
		req.Header.Set("Referer", "api.cctv.cn")
		req.Header.Set("UID", uid)

		resp, _ := client.Do(req)
		defer resp.Body.Close()
		var body strings.Builder
		_, _ = io.Copy(&body, resp.Body)

		data := body.String()
		re := regexp.MustCompile(`(.*\.m3u8\?.*)`)
		matches := re.FindStringSubmatch(data)
		if len(matches) > 0 {
			playURL = path + matches[0]
		} else {
			return data
		}
	}
}

func getTs(url string) string {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "cctv_app_tv")
	req.Header.Set("Referer", "https://api.cctv.cn/")
	req.Header.Set("UID", "1234123122")
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-encoding", "gzip, deflate")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	var body strings.Builder
	_, _ = io.Copy(&body, resp.Body)

	return body.String()
}

func getCache(key string) (string, bool) {
	if item, found := cache.Load(key); found {
		cacheItem := item.(CacheItem)
		if time.Now().Unix() < cacheItem.Expiration {
			return cacheItem.Value, true
		}
	}
	return "", false
}

func setCache(key, value string) {
	cache.Store(key, CacheItem{
		Value:      value,
		Expiration: time.Now().Unix() + 3600,
	})
}