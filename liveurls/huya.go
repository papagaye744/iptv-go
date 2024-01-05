// Package liveurls
// @Time:2023/02/05 23:34
// @File:huya.go
// @SoftWare:Goland
// @Author:feiyang
// @Contact:TG@feiyangdigital

package liveurls

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/hr3lxphr6j/requests"
	uuid "github.com/satori/go.uuid"
	"github.com/tidwall/gjson"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const userAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1"

type Huya struct {
	Rid     string
	Cdn     string
	CdnType string
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func parseAntiCode(anticode string, uid int64, streamName string) (string, error) {
	qr, err := url.ParseQuery(anticode)
	if err != nil {
		return "", err
	}
	qr.Set("ver", "1")
	qr.Set("sv", "2110211124")
	qr.Set("seqid", strconv.FormatInt(time.Now().Unix()*1000+uid, 10))
	qr.Set("uid", strconv.FormatInt(uid, 10))
	reluuid, _ := uuid.NewV4()
	qr.Set("uuid", reluuid.String())
	ss := GetMD5Hash(fmt.Sprintf("%s|%s|%s", qr.Get("seqid"), qr.Get("ctype"), qr.Get("t")))
	wsTime := strconv.FormatInt(time.Now().Add(6*time.Hour).Unix(), 16)
	decodeString, _ := base64.StdEncoding.DecodeString(qr.Get("fm"))
	fm := string(decodeString)
	fm = strings.ReplaceAll(fm, "$0", qr.Get("uid"))
	fm = strings.ReplaceAll(fm, "$1", streamName)
	fm = strings.ReplaceAll(fm, "$2", ss)
	fm = strings.ReplaceAll(fm, "$3", wsTime)
	qr.Set("wsSecret", GetMD5Hash(fm))
	qr.Set("ratio", "0")
	qr.Set("wsTime", wsTime)
	return qr.Encode(), nil
}

func (h *Huya) GetLiveUrl() any {
	mobileUrl := "https://m.huya.com/" + h.Rid
	resp, err := requests.Get(mobileUrl, requests.UserAgent(userAgent))
	if err != nil {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	body, err := resp.Text()
	if err != nil {
		return nil
	}
	tmpStrings := strings.Split(body, `"tLiveInfo":`)
	if len(tmpStrings) < 2 {
		return nil
	}
	liveInfoJsonRawString := strings.Split(tmpStrings[1], `,"_classname":"LiveRoom.LiveInfo"}`)[0] + "}"
	if !gjson.Valid(liveInfoJsonRawString) {
		return nil
	}
	liveInfoJson := gjson.Parse(liveInfoJsonRawString)
	streamInfoJsons := liveInfoJson.Get("tLiveStreamInfo.vStreamInfo.value")
	if len(streamInfoJsons.Array()) == 0 {
		return nil
	}
	var cdnSlice []string
	var finalurl string
	streamInfoJsons.ForEach(func(key, value gjson.Result) bool {
		var cdnType = gjson.Get(value.String(), "sCdnType").String()
		cdnSlice = append(cdnSlice, cdnType)
		if cdnType == h.Cdn {
			sStreamName := gjson.Get(value.String(), "sStreamName").String()
			sFlvAntiCode := gjson.Get(value.String(), "sFlvAntiCode").String()
			sFlvUrl := gjson.Get(value.String(), "sFlvUrl").String()
			uid := rand.Int63n(99999999999) + 1400000000000
			query, _ := parseAntiCode(sFlvAntiCode, uid, sStreamName)
			finalurl = strings.Replace(fmt.Sprintf("%s/%s.flv?%s", sFlvUrl, sStreamName, query), "http://", "https://", 1)
		}
		return true
	})
	if h.CdnType == "display" {
		return cdnSlice
	}
	return finalurl
}