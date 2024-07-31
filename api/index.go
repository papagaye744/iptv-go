package handler
 
import (
  "Golang/liveurls"
  "Golang/list"
  "Golang/utils"
  "fmt"
  "net/http"
  "strings"
  "encoding/json"
  "strconv"
  "log"
  "github.com/gin-gonic/gin"	
)



// vercel 平台会将请求传递给该函数，这个函数名随意，但函数参数必须按照该规则。
func Handler(w http.ResponseWriter, r *http.Request) {
  

  gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.HEAD("/", func(c *gin.Context) {
		c.String(http.StatusOK, "请求成功！")
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "请求成功！")
	})

	r.GET("/tv.m3u", func(c *gin.Context) {
		if enableTV {
			itvm3uobj := &list.Tvm3u{}
			c.Writer.Header().Set("Content-Type", "application/octet-stream")
			c.Writer.Header().Set("Content-Disposition", "attachment; filename=tv.m3u")
			itvm3uobj.GetTvM3u(c)
		} else {
			c.String(http.StatusForbidden, "公共服务不提供TV直播")
		}
	})

	r.GET("/huyayqk.m3u", func(c *gin.Context) {
		yaobj := &list.HuyaYqk{}
		res, _ := yaobj.HuYaYqk("https://live.cdn.huya.com/liveHttpUI/getLiveList?iGid=2135")
		var result list.YaResponse
		json.Unmarshal(res, &result)
		pageCount := result.ITotalPage
		pageSize := result.IPageSize
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Disposition", "attachment; filename=huyayqk.m3u")
		getTestVideoUrl(c)

		for i := 1; i <= pageCount; i++ {
			apiRes, _ := yaobj.HuYaYqk(fmt.Sprintf("https://live.cdn.huya.com/liveHttpUI/getLiveList?iGid=2135&iPageNo=%d&iPageSize=%d", i, pageSize))
			var res list.YaResponse
			json.Unmarshal(apiRes, &res)
			data := res.VList
			for _, value := range data {
				fmt.Fprintf(c.Writer, "#EXTINF:-1 tvg-logo=\"%s\" group-title=\"%s\", %s\n", value.SAvatar180, value.SGameFullName, value.SNick)
				fmt.Fprintf(c.Writer, "%s/huya/%v\n", getLivePrefix(c), value.LProfileRoom)
			}
		}
	})

	r.GET("/douyuyqk.m3u", func(c *gin.Context) {
		yuobj := &list.DouYuYqk{}
		resAPI, _ := yuobj.Douyuyqk("https://www.douyu.com/gapi/rkc/directory/mixList/2_208/list")

		var result list.DouYuResponse
		json.Unmarshal(resAPI, &result)
		pageCount := result.Data.Pgcnt

		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Disposition", "attachment; filename=douyuyqk.m3u")
		getTestVideoUrl(c)

		for i := 1; i <= pageCount; i++ {
			apiRes, _ := yuobj.Douyuyqk("https://www.douyu.com/gapi/rkc/directory/mixList/2_208/" + strconv.Itoa(i))

			var res list.DouYuResponse
			json.Unmarshal(apiRes, &res)
			data := res.Data.Rl

			for _, value := range data {
				fmt.Fprintf(c.Writer, "#EXTINF:-1 tvg-logo=\"https://apic.douyucdn.cn/upload/%s_big.jpg\" group-title=\"%s\", %s\n", value.Av, value.C2name, value.Nn)
				fmt.Fprintf(c.Writer, "%s/douyu/%v\n", getLivePrefix(c), value.Rid)
			}
		}
	})

	r.GET("/yylunbo.m3u", func(c *gin.Context) {
		yylistobj := &list.Yylist{}
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Disposition", "attachment; filename=yylunbo.m3u")
		getTestVideoUrl(c)

		i := 1
		for {
			apiRes := yylistobj.Yylb(fmt.Sprintf("http://rubiks-ipad.yy.com/nav/other/idx/213?channel=appstore&ispType=0&model=iPad8,6&netType=2&os=iOS&osVersion=17.2&page=%d&uid=0&yyVersion=6.17.0", i))
			var res list.ApiResponse
			json.Unmarshal([]byte(apiRes), &res)
			for _, value := range res.Data.Data {
				fmt.Fprintf(c.Writer, "#EXTINF:-1 tvg-logo=\"%s\" group-title=\"%s\", %s\n", value.Avatar, value.Biz, value.Desc)
				fmt.Fprintf(c.Writer, "%s/yy/%v\n", getLivePrefix(c), value.Sid)
			}
			if res.Data.IsLastPage == 1 {
				break
			}
			i++
		}
	})

	r.GET("/:path/:rid", func(c *gin.Context) {
		path := c.Param("path")
		rid := c.Param("rid")
		ts := c.Query("ts")
		switch path {
		case "itv":
			if enableTV {
				itvobj := &liveurls.Itv{}
				cdn := c.Query("cdn")
				if ts == "" {
					itvobj.HandleMainRequest(c, cdn, rid)
				} else {
					itvobj.HandleTsRequest(c, ts)
				}
			} else {
				c.String(http.StatusForbidden, "公共服务不提供TV直播")
			}
		case "ysptp":
			if enableTV {
				ysptpobj := &liveurls.Ysptp{}
				if ts == "" {
					ysptpobj.HandleMainRequest(c, rid)
				} else {
					ysptpobj.HandleTsRequest(c, ts, c.Query("wsTime"))
				}
			} else {
				c.String(http.StatusForbidden, "公共服务不提供TV直播")
			}
		case "douyin":
			douyinobj := &liveurls.Douyin{}
			douyinobj.Rid = rid
			douyinobj.Stream = c.DefaultQuery("stream", "flv")
			c.Redirect(http.StatusMovedPermanently, duanyan(adurl, douyinobj.GetDouYinUrl()))
		case "douyu":
			douyuobj := &liveurls.Douyu{}
			douyuobj.Rid = rid
			douyuobj.Stream_type = c.DefaultQuery("stream", "flv")
			c.Redirect(http.StatusMovedPermanently, duanyan(adurl, douyuobj.GetRealUrl()))
		case "huya":
			huyaobj := &liveurls.Huya{}
			huyaobj.Rid = rid
			huyaobj.Cdn = c.DefaultQuery("cdn", "hwcdn")
			huyaobj.Media = c.DefaultQuery("media", "flv")
			huyaobj.Type = c.DefaultQuery("type", "nodisplay")
			if huyaobj.Type == "display" {
				c.JSON(200, huyaobj.GetLiveUrl())
			} else {
				c.Redirect(http.StatusMovedPermanently, duanyan(adurl, huyaobj.GetLiveUrl()))
			}
		case "bilibili":
			biliobj := &liveurls.BiliBili{}
			biliobj.Rid = rid
			biliobj.Platform = c.DefaultQuery("platform", "web")
			biliobj.Quality = c.DefaultQuery("quality", "10000")
			biliobj.Line = c.DefaultQuery("line", "first")
			c.Redirect(http.StatusMovedPermanently, duanyan(adurl, biliobj.GetPlayUrl()))
		case "youtube":
			ytbObj := &liveurls.Youtube{}
			ytbObj.Rid = rid
			ytbObj.Quality = c.DefaultQuery("quality", "1080")
			c.Redirect(http.StatusMovedPermanently, duanyan(adurl, ytbObj.GetLiveUrl()))
		case "yy":
			yyObj := &liveurls.Yy{}
			yyObj.Rid = rid
			yyObj.Quality = c.DefaultQuery("quality", "4")
			c.Redirect(http.StatusMovedPermanently, duanyan(adurl, yyObj.GetLiveUrl()))
		}
	})
	return r
}
