package liveurls

import (
	"context"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Itv struct{}

var (
	hostMappings = map[string]string{
		"cache.ott.ystenlive.itv.cmvideo.cn": "feiyangdigital.tg.ystenlive.ottdns.com",
		"cache.ott.bestlive.itv.cmvideo.cn":  "feiyangdigital.tg.bestlive.ottdns.com",
		"cache.ott.wasulive.itv.cmvideo.cn":  "feiyangdigital.tg.wasulive.ottdns.com",
		"cache.ott.fifalive.itv.cmvideo.cn":  "feiyangdigital.tg.fifalive.ottdns.com",
		"cache.ott.hnbblive.itv.cmvideo.cn":  "feiyangdigital.tg.hnbblive.ottdns.com",
	}
	programList = map[string]string{
		"wasusyt/6000000001000029752.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000029752/1.m3u8?channel-id=wasusyt&Contentid=6000000001000029752&livemode=1&stbId=3",
		"bestzb/5000000004000002226.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000004000002226/1.m3u8?channel-id=bestzb&Contentid=5000000004000002226&livemode=1&stbId=3",
		"ystenlive/1000000005000265001.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265001/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265001&livemode=1&stbId=3",
		"ystenlive/1000000001000023315.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000023315/1.m3u8?channel-id=ystenlive&Contentid=1000000001000023315&livemode=1&stbId=3",
		"wasusyt/6000000001000014161.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000014161/1.m3u8?channel-id=wasusyt&Contentid=6000000001000014161&livemode=1&stbId=3",
		"wasusyt/6000000001000022313.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000022313/1.m3u8?channel-id=wasusyt&Contentid=6000000001000022313&livemode=1&stbId=3",
		"ystenlive/1000000005000265003.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265003/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265003&livemode=1&stbId=3",
		"bestzb/5000000011000031102.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031102/1.m3u8?channel-id=bestzb&Contentid=5000000011000031102&livemode=1&stbId=3",
		"ystenlive/1000000005000265004.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265004/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265004&livemode=1&stbId=3",
		"ystenlive/1000000005000025222.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000025222/1.m3u8?channel-id=ystenlive&Contentid=1000000005000025222&livemode=1&stbId=3",
		"ystenlive/1000000005000265005.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265005/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265005&livemode=1&stbId=3",
		"wasusyt/6000000001000015875.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000015875/1.m3u8?channel-id=wasusyt&Contentid=6000000001000015875&livemode=1&stbId=3",
		"ystenlive/1000000005000265016.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265016/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265016&livemode=1&stbId=3",
		"ystenlive/1000000001000001737.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000001737/1.m3u8?channel-id=ystenlive&Contentid=1000000001000001737&livemode=1&stbId=3",
		"wasusyt/6000000001000004574.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000004574/1.m3u8?channel-id=wasusyt&Contentid=6000000001000004574&livemode=1&stbId=3",
		"ystenlive/1000000005000265006.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265006/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265006&livemode=1&stbId=3",
		"ystenlive/1000000001000024341.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000024341/1.m3u8?channel-id=ystenlive&Contentid=1000000001000024341&livemode=1&stbId=3",
		"wasusyt/6000000001000009055.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000009055/1.m3u8?channel-id=wasusyt&Contentid=6000000001000009055&livemode=1&stbId=3",
		"ystenlive/1000000005000265007.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265007/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265007&livemode=1&stbId=3",
		"wasusyt/6000000001000001070.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000001070/1.m3u8?channel-id=wasusyt&Contentid=6000000001000001070&livemode=1&stbId=3",
		"ystenlive/1000000005000265008.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265008/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265008&livemode=1&stbId=3",
		"ystenlive/1000000001000014583.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000014583/1.m3u8?channel-id=ystenlive&Contentid=1000000001000014583&livemode=1&stbId=3",
		"wasusyt/6000000001000032162.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000032162/1.m3u8?channel-id=wasusyt&Contentid=6000000001000032162&livemode=1&stbId=3",
		"ystenlive/1000000005000265009.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265009/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265009&livemode=1&stbId=3",
		"ystenlive/1000000001000023734.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000023734/1.m3u8?channel-id=ystenlive&Contentid=1000000001000023734&livemode=1&stbId=3",
		"bestzb/5000000004000012827.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000004000012827/1.m3u8?channel-id=bestzb&Contentid=5000000004000012827&livemode=1&stbId=3",
		"ystenlive/1000000005000265010.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265010/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265010&livemode=1&stbId=3",
		"bestzb/5000000011000031106.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031106/1.m3u8?channel-id=bestzb&Contentid=5000000011000031106&livemode=1&stbId=3",
		"ystenlive/1000000005000265011.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265011/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265011&livemode=1&stbId=3",
		"ystenlive/1000000001000032494.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000032494/1.m3u8?channel-id=ystenlive&Contentid=1000000001000032494&livemode=1&stbId=3",
		"wasusyt/6000000001000022586.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000022586/1.m3u8?channel-id=wasusyt&Contentid=6000000001000022586&livemode=1&stbId=3",
		"ystenlive/1000000005000265012.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265012/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265012&livemode=1&stbId=3",
		"bestzb/5000000011000031108.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031108/1.m3u8?channel-id=bestzb&Contentid=5000000011000031108&livemode=1&stbId=3",
		"ystenlive/1000000001000008170.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000008170/1.m3u8?channel-id=ystenlive&Contentid=1000000001000008170&livemode=1&stbId=3",
		"bestzb/5000000004000006673.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000004000006673/1.m3u8?channel-id=bestzb&Contentid=5000000004000006673&livemode=1&stbId=3",
		"ystenlive/1000000005000265013.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265013/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265013&livemode=1&stbId=3",
		"bestzb/5000000011000031109.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031109/1.m3u8?channel-id=bestzb&Contentid=5000000011000031109&livemode=1&stbId=3",
		"ystenlive/1000000005000265014.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265014/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265014&livemode=1&stbId=3",
		"ystenlive/1000000006000233002.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000006000233002/1.m3u8?channel-id=ystenlive&Contentid=1000000006000233002&livemode=1&stbId=3",
		"bestzb/5000000008000023254.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000008000023254/1.m3u8?channel-id=bestzb&Contentid=5000000008000023254&livemode=1&stbId=3",
		"ystenlive/1000000006000268004.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000006000268004/1.m3u8?channel-id=ystenlive&Contentid=1000000006000268004&livemode=1&stbId=3",
		"ystenlive/1000000005000265015.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265015/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265015&livemode=1&stbId=3",
		"hnbblive/7745129417417101820.m3u8":    "http://gslbserv.itv.cmvideo.cn:80/7745129417417101820/1.m3u8?channel-id=hnbblive&Contentid=7745129417417101820&livemode=1&stbId=3",
		"hnbblive/7114647837765104058.m3u8":    "http://gslbserv.itv.cmvideo.cn:80/7114647837765104058/1.m3u8?channel-id=hnbblive&Contentid=7114647837765104058&livemode=1&stbId=3",
		"bestzb/5000000002000002652.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000002000002652/1.m3u8?channel-id=bestzb&Contentid=5000000002000002652&livemode=1&stbId=3",
		"bestzb/5000000011000031126.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031126/1.m3u8?channel-id=bestzb&Contentid=5000000011000031126&livemode=1&stbId=3",
		"wasusyt/6000000001000020451.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000020451/1.m3u8?channel-id=wasusyt&Contentid=6000000001000020451&livemode=1&stbId=3",
		"ystenlive/1000000005000265027.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265027/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265027&livemode=1&stbId=3",
		"ystenlive/1000000001000001910.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000001910/1.m3u8?channel-id=ystenlive&Contentid=1000000001000001910&livemode=1&stbId=3",
		"ystenlive/1000000005000265020.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265020/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265020&livemode=1&stbId=3",
		"bestzb/7851974109718180595.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/7851974109718180595/1.m3u8?channel-id=bestzb&Contentid=7851974109718180595&livemode=1&stbId=3",
		"ystenlive/1000000001000030159.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000030159/1.m3u8?channel-id=ystenlive&Contentid=1000000001000030159&livemode=1&stbId=3",
		"wasusyt/6000000001000009954.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000009954/1.m3u8?channel-id=wasusyt&Contentid=6000000001000009954&livemode=1&stbId=3",
		"ystenlive/1000000005000265025.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265025/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265025&livemode=1&stbId=3",
		"bestzb/5000000004000010584.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000004000010584/1.m3u8?channel-id=bestzb&Contentid=5000000004000010584&livemode=1&stbId=3",
		"ystenlive/1000000005000265033.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265033/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265033&livemode=1&stbId=3",
		"bestzb/5000000011000031121.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031121/1.m3u8?channel-id=bestzb&Contentid=5000000011000031121&livemode=1&stbId=3",
		"ystenlive/1000000001000014176.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000014176/1.m3u8?channel-id=ystenlive&Contentid=1000000001000014176&livemode=1&stbId=3",
		"wasusyt/6000000001000031076.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000031076/1.m3u8?channel-id=wasusyt&Contentid=6000000001000031076&livemode=1&stbId=3",
		"ystenlive/1000000005000265034.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265034/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265034&livemode=1&stbId=3",
		"bestzb/5000000011000031118.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031118/1.m3u8?channel-id=bestzb&Contentid=5000000011000031118&livemode=1&stbId=3",
		"bestzb/5000000004000025843.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000004000025843/1.m3u8?channel-id=bestzb&Contentid=5000000004000025843&livemode=1&stbId=3",
		"bestzb/5000000004000006211.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000004000006211/1.m3u8?channel-id=bestzb&Contentid=5000000004000006211&livemode=1&stbId=3",
		"bestzb/5000000006000040016.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000006000040016/1.m3u8?channel-id=bestzb&Contentid=5000000006000040016&livemode=1&stbId=3",
		"bestzb/5000000011000031119.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031119/1.m3u8?channel-id=bestzb&Contentid=5000000011000031119&livemode=1&stbId=3",
		"ystenlive/1000000001000001925.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000001925/1.m3u8?channel-id=ystenlive&Contentid=1000000001000001925&livemode=1&stbId=3",
		"wasusyt/6000000001000016510.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000016510/1.m3u8?channel-id=wasusyt&Contentid=6000000001000016510&livemode=1&stbId=3",
		"ystenlive/1000000005000265029.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265029/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265029&livemode=1&stbId=3",
		"ystenlive/1000000001000024621.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000024621/1.m3u8?channel-id=ystenlive&Contentid=1000000001000024621&livemode=1&stbId=3",
		"wasusyt/6000000001000015436.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000015436/1.m3u8?channel-id=wasusyt&Contentid=6000000001000015436&livemode=1&stbId=3",
		"ystenlive/1000000005000265023.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265023/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265023&livemode=1&stbId=3",
		"bestzb/5000000004000006692.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000004000006692/1.m3u8?channel-id=bestzb&Contentid=5000000004000006692&livemode=1&stbId=3",
		"wasusyt/6000000001000018044.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000018044/1.m3u8?channel-id=wasusyt&Contentid=6000000001000018044&livemode=1&stbId=3",
		"ystenlive/1000000005000265024.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265024/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265024&livemode=1&stbId=3",
		"bestzb/5000000011000031203.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031203/1.m3u8?channel-id=bestzb&Contentid=5000000011000031203&livemode=1&stbId=3",
		"bestzb/5000000011000031206.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031206/1.m3u8?channel-id=bestzb&Contentid=5000000011000031206&livemode=1&stbId=3",
		"bestzb/5000000011000031209.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031209/1.m3u8?channel-id=bestzb&Contentid=5000000011000031209&livemode=1&stbId=3",
		"bestzb/5000000011000031117.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031117/1.m3u8?channel-id=bestzb&Contentid=5000000011000031117&livemode=1&stbId=3",
		"wasusyt/6000000001000014861.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000014861/1.m3u8?channel-id=wasusyt&Contentid=6000000001000014861&livemode=1&stbId=3",
		"ystenlive/1000000001000001828.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000001828/1.m3u8?channel-id=ystenlive&Contentid=1000000001000001828&livemode=1&stbId=3",
		"ystenlive/1000000005000265030.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265030/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265030&livemode=1&stbId=3",
		"ystenlive/1000000006000268001.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000006000268001/1.m3u8?channel-id=ystenlive&Contentid=1000000006000268001&livemode=1&stbId=3",
		"ystenlive/1000000005000265032.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265032/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265032&livemode=1&stbId=3",
		"bestzb/5000000004000011671.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000004000011671/1.m3u8?channel-id=bestzb&Contentid=5000000004000011671&livemode=1&stbId=3",
		"ystenlive/1000000005000265022.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265022/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265022&livemode=1&stbId=3",
		"ystenlive/1000000002000013359.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000002000013359/1.m3u8?channel-id=ystenlive&Contentid=1000000002000013359&livemode=1&stbId=3",
		"ystenlive/1000000001000016568.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000016568/1.m3u8?channel-id=ystenlive&Contentid=1000000001000016568&livemode=1&stbId=3",
		"wasusyt/6000000001000004134.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000004134/1.m3u8?channel-id=wasusyt&Contentid=6000000001000004134&livemode=1&stbId=3",
		"ystenlive/1000000005000265019.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265019/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265019&livemode=1&stbId=3",
		"wasusyt/6000000001000003639.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000003639/1.m3u8?channel-id=wasusyt&Contentid=6000000001000003639&livemode=1&stbId=3",
		"bestzb/5000000004000014098.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000004000014098/1.m3u8?channel-id=bestzb&Contentid=5000000004000014098&livemode=1&stbId=3",
		"ystenlive/1000000005000265018.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265018/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265018&livemode=1&stbId=3",
		"bestzb/5000000010000030951.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000010000030951/1.m3u8?channel-id=bestzb&Contentid=5000000010000030951&livemode=1&stbId=3",
		"bestzb/5000000010000027146.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000010000027146/1.m3u8?channel-id=bestzb&Contentid=5000000010000027146&livemode=1&stbId=3",
		"bestzb/5000000007000010003.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000007000010003/1.m3u8?channel-id=bestzb&Contentid=5000000007000010003&livemode=1&stbId=3",
		"bestzb/5000000010000032212.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000010000032212/1.m3u8?channel-id=bestzb&Contentid=5000000010000032212&livemode=1&stbId=3",
		"bestzb/5000000010000018926.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000010000018926/1.m3u8?channel-id=bestzb&Contentid=5000000010000018926&livemode=1&stbId=3",
		"hnbblive/2000000002000000014.m3u8":    "http://gslbserv.itv.cmvideo.cn:80/2000000002000000014/1.m3u8?channel-id=hnbblive&Contentid=2000000002000000014&livemode=1&stbId=3",
		"bestzb/5000000011000031123.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031123/1.m3u8?channel-id=bestzb&Contentid=5000000011000031123&livemode=1&stbId=3",
		"bestzb/5000000004000010282.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000004000010282/1.m3u8?channel-id=bestzb&Contentid=5000000004000010282&livemode=1&stbId=3",
		"ystenlive/1000000005000265021.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265021/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265021&livemode=1&stbId=3",
		"bestzb/5000000010000017540.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000010000017540/1.m3u8?channel-id=bestzb&Contentid=5000000010000017540&livemode=1&stbId=3",
		"bestzb/5000000011000031110.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031110/1.m3u8?channel-id=bestzb&Contentid=5000000011000031110&livemode=1&stbId=3",
		"bestzb/5000000004000007410.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000004000007410/1.m3u8?channel-id=bestzb&Contentid=5000000004000007410&livemode=1&stbId=3",
		"wasusyt/6000000001000002116.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000002116/1.m3u8?channel-id=wasusyt&Contentid=6000000001000002116&livemode=1&stbId=3",
		"ystenlive/1000000005000265028.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265028/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265028&livemode=1&stbId=3",
		"bestzb/5000000004000006119.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000004000006119/1.m3u8?channel-id=bestzb&Contentid=5000000004000006119&livemode=1&stbId=3",
		"bestzb/5000000004000006827.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000004000006827/1.m3u8?channel-id=bestzb&Contentid=5000000004000006827&livemode=1&stbId=3",
		"wasusyt/6000000001000009186.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000001000009186/1.m3u8?channel-id=wasusyt&Contentid=6000000001000009186&livemode=1&stbId=3",
		"ystenlive/1000000005000265026.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265026/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265026&livemode=1&stbId=3",
		"bestzb/5000000011000031120.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031120/1.m3u8?channel-id=bestzb&Contentid=5000000011000031120&livemode=1&stbId=3",
		"bestzb/5000000004000007275.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000004000007275/1.m3u8?channel-id=bestzb&Contentid=5000000004000007275&livemode=1&stbId=3",
		"ystenlive/1000000001000014260.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000014260/1.m3u8?channel-id=ystenlive&Contentid=1000000001000014260&livemode=1&stbId=3",
		"ystenlive/1000000005000265031.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265031/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265031&livemode=1&stbId=3",
		"ystenlive/1000000001000001096.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000001096/1.m3u8?channel-id=ystenlive&Contentid=1000000001000001096&livemode=1&stbId=3",
		"ystenlive/1000000005000265017.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000265017/1.m3u8?channel-id=ystenlive&Contentid=1000000005000265017&livemode=1&stbId=3",
		"wasusyt/6000000003000004748.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000003000004748/1.m3u8?channel-id=wasusyt&Contentid=6000000003000004748&livemode=1&stbId=3",
		"ystenlive/1000000004000011651.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000004000011651/1.m3u8?channel-id=ystenlive&Contentid=1000000004000011651&livemode=1&stbId=3",
		"FifastbLive/3000000010000005180.m3u8": "http://gslbserv.itv.cmvideo.cn:80/3000000010000005180/1.m3u8?channel-id=FifastbLive&Contentid=3000000010000005180&livemode=1&stbId=3",
		"FifastbLive/3000000010000015686.m3u8": "http://gslbserv.itv.cmvideo.cn:80/3000000010000015686/1.m3u8?channel-id=FifastbLive&Contentid=3000000010000015686&livemode=1&stbId=3",
		"FifastbLive/3000000020000031315.m3u8": "http://gslbserv.itv.cmvideo.cn:80/3000000020000031315/1.m3u8?channel-id=FifastbLive&Contentid=3000000020000031315&livemode=1&stbId=3",
		"wasusyt/6000000002000010046.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000002000010046/1.m3u8?channel-id=wasusyt&Contentid=6000000002000010046&livemode=1&stbId=3",
		"wasusyt/6000000002000032052.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000002000032052/1.m3u8?channel-id=wasusyt&Contentid=6000000002000032052&livemode=1&stbId=3",
		"wasusyt/6000000002000032344.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000002000032344/1.m3u8?channel-id=wasusyt&Contentid=6000000002000032344&livemode=1&stbId=3",
		"wasusyt/6000000002000003382.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000002000003382/1.m3u8?channel-id=wasusyt&Contentid=6000000002000003382&livemode=1&stbId=3",
		"ystenlive/1000000004000019008.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000004000019008/1.m3u8?channel-id=ystenlive&Contentid=1000000004000019008&livemode=1&stbId=3",
		"ystenlive/1000000004000013968.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000004000013968/1.m3u8?channel-id=ystenlive&Contentid=1000000004000013968&livemode=1&stbId=3",
		"ystenlive/1000000004000013730.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000004000013730/1.m3u8?channel-id=ystenlive&Contentid=1000000004000013730&livemode=1&stbId=3",
		"ystenlive/1000000004000014634.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000004000014634/1.m3u8?channel-id=ystenlive&Contentid=1000000004000014634&livemode=1&stbId=3",
		"ystenlive/1000000006000032328.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000006000032328/1.m3u8?channel-id=ystenlive&Contentid=1000000006000032328&livemode=1&stbId=3",
		"hnbblive/2000000003000000010.m3u8":    "http://gslbserv.itv.cmvideo.cn:80/2000000003000000010/1.m3u8?channel-id=hnbblive&Contentid=2000000003000000010&livemode=1&stbId=3",
		"ystenlive/1000000006000268003.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000006000268003/1.m3u8?channel-id=ystenlive&Contentid=1000000006000268003&livemode=1&stbId=3",
		"ystenlive/1000000003000012426.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000003000012426/1.m3u8?channel-id=ystenlive&Contentid=1000000003000012426&livemode=1&stbId=3",
		"ystenlive/1000000001000009601.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000009601/1.m3u8?channel-id=ystenlive&Contentid=1000000001000009601&livemode=1&stbId=3",
		"ystenlive/1000000006000268002.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000006000268002/1.m3u8?channel-id=ystenlive&Contentid=1000000006000268002&livemode=1&stbId=3",
		"hnbblive/2000000003000000018.m3u8":    "http://gslbserv.itv.cmvideo.cn:80/2000000003000000018/1.m3u8?channel-id=hnbblive&Contentid=2000000003000000018&livemode=1&stbId=3",
		"ystenlive/1000000005000266013.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000266013/1.m3u8?channel-id=ystenlive&Contentid=1000000005000266013&livemode=1&stbId=3",
		"ystenlive/1000000004000018653.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000004000018653/1.m3u8?channel-id=ystenlive&Contentid=1000000004000018653&livemode=1&stbId=3",
		"hnbblive/2000000003000000024.m3u8":    "http://gslbserv.itv.cmvideo.cn:80/2000000003000000024/1.m3u8?channel-id=hnbblive&Contentid=2000000003000000024&livemode=1&stbId=3",
		"ystenlive/1000000005000266012.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000266012/1.m3u8?channel-id=ystenlive&Contentid=1000000005000266012&livemode=1&stbId=3",
		"ystenlive/1000000004000008284.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000004000008284/1.m3u8?channel-id=ystenlive&Contentid=1000000004000008284&livemode=1&stbId=3",
		"ystenlive/1000000004000026167.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000004000026167/1.m3u8?channel-id=ystenlive&Contentid=1000000004000026167&livemode=1&stbId=3",
		"ystenlive/1000000004000024282.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000004000024282/1.m3u8?channel-id=ystenlive&Contentid=1000000004000024282&livemode=1&stbId=3",
		"hnbblive/2000000003000000014.m3u8":    "http://gslbserv.itv.cmvideo.cn:80/2000000003000000014/1.m3u8?channel-id=hnbblive&Contentid=2000000003000000014&livemode=1&stbId=3",
		"hnbblive/2000000003000000022.m3u8":    "http://gslbserv.itv.cmvideo.cn:80/2000000003000000022/1.m3u8?channel-id=hnbblive&Contentid=2000000003000000022&livemode=1&stbId=3",
		"ystenlive/1000000001000006197.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000006197/1.m3u8?channel-id=ystenlive&Contentid=1000000001000006197&livemode=1&stbId=3",
		"hnbblive/2000000003000000016.m3u8":    "http://gslbserv.itv.cmvideo.cn:80/2000000003000000016/1.m3u8?channel-id=hnbblive&Contentid=2000000003000000016&livemode=1&stbId=3",
		"hnbblive/2000000003000000003.m3u8":    "http://gslbserv.itv.cmvideo.cn:80/2000000003000000003/1.m3u8?channel-id=hnbblive&Contentid=2000000003000000003&livemode=1&stbId=3",
		"hnbblive/2000000003000000007.m3u8":    "http://gslbserv.itv.cmvideo.cn:80/2000000003000000007/1.m3u8?channel-id=hnbblive&Contentid=2000000003000000007&livemode=1&stbId=3",
		"ystenlive/1000000001000000515.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000000515/1.m3u8?channel-id=ystenlive&Contentid=1000000001000000515&livemode=1&stbId=3",
		"ystenlive/1000000005000266011.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000005000266011/1.m3u8?channel-id=ystenlive&Contentid=1000000005000266011&livemode=1&stbId=3",
		"hnbblive/2000000003000000009.m3u8":    "http://gslbserv.itv.cmvideo.cn:80/2000000003000000009/1.m3u8?channel-id=hnbblive&Contentid=2000000003000000009&livemode=1&stbId=3",
		"ystenlive/1000000004000019624.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000004000019624/1.m3u8?channel-id=ystenlive&Contentid=1000000004000019624&livemode=1&stbId=3",
		"ystenlive/1000000004000021734.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000004000021734/1.m3u8?channel-id=ystenlive&Contentid=1000000004000021734&livemode=1&stbId=3",
		"ystenlive/1000000006000032327.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000006000032327/1.m3u8?channel-id=ystenlive&Contentid=1000000006000032327&livemode=1&stbId=3",
		"ystenlive/1000000001000003775.m3u8":   "http://gslbserv.itv.cmvideo.cn:80/1000000001000003775/1.m3u8?channel-id=ystenlive&Contentid=1000000001000003775&livemode=1&stbId=3",
		"bestzb/5000000011000031113.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031113/1.m3u8?channel-id=bestzb&Contentid=5000000011000031113&livemode=1&stbId=3",
		"bestzb/5000000011000031111.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031111/1.m3u8?channel-id=bestzb&Contentid=5000000011000031111&livemode=1&stbId=3",
		"bestzb/9001547084732463424.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/9001547084732463424/1.m3u8?channel-id=bestzb&Contentid=9001547084732463424&livemode=1&stbId=3",
		"bestzb/5000000002000009455.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000002000009455/1.m3u8?channel-id=bestzb&Contentid=5000000002000009455&livemode=1&stbId=3",
		"bestzb/5000000007000010001.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000007000010001/1.m3u8?channel-id=bestzb&Contentid=5000000007000010001&livemode=1&stbId=3",
		"bestzb/5000000010000026105.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000010000026105/1.m3u8?channel-id=bestzb&Contentid=5000000010000026105&livemode=1&stbId=3",
		"bestzb/5000000002000029972.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000002000029972/1.m3u8?channel-id=bestzb&Contentid=5000000002000029972&livemode=1&stbId=3",
		"bestzb/5000000011000031112.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031112/1.m3u8?channel-id=bestzb&Contentid=5000000011000031112&livemode=1&stbId=3",
		"bestzb/5000000011000031207.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031207/1.m3u8?channel-id=bestzb&Contentid=5000000011000031207&livemode=1&stbId=3",
		"bestzb/5000000011000031116.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031116/1.m3u8?channel-id=bestzb&Contentid=5000000011000031116&livemode=1&stbId=3",
		"bestzb/5000000002000019634.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000002000019634/1.m3u8?channel-id=bestzb&Contentid=5000000002000019634&livemode=1&stbId=3",
		"bestzb/5000000011000031114.m3u8":      "http://gslbserv.itv.cmvideo.cn:80/5000000011000031114/1.m3u8?channel-id=bestzb&Contentid=5000000011000031114&livemode=1&stbId=3",
		"wasusyt/6000000006000230630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000230630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000230630&livemode=1&stbId=3",
		"wasusyt/6000000006000070630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000070630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000070630&livemode=1&stbId=3",
		"wasusyt/6000000006000280630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000280630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000280630&livemode=1&stbId=3",
		"wasusyt/6000000006000080630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000080630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000080630&livemode=1&stbId=3",
		"wasusyt/6000000006000260630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000260630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000260630&livemode=1&stbId=3",
		"wasusyt/6000000006000060630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000060630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000060630&livemode=1&stbId=3",
		"wasusyt/6000000006000020630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000020630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000020630&livemode=1&stbId=3",
		"wasusyt/6000000006000160630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000160630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000160630&livemode=1&stbId=3",
		"wasusyt/6000000006000040630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000040630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000040630&livemode=1&stbId=3",
		"wasusyt/6000000006000150630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000150630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000150630&livemode=1&stbId=3",
		"wasusyt/6000000006000250630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000250630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000250630&livemode=1&stbId=3",
		"wasusyt/6000000006000270630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000270630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000270630&livemode=1&stbId=3",
		"wasusyt/6000000006000100630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000100630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000100630&livemode=1&stbId=3",
		"wasusyt/6000000006000240630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000240630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000240630&livemode=1&stbId=3",
		"wasusyt/6000000006000290630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000290630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000290630&livemode=1&stbId=3",
		"wasusyt/6000000006000220630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000220630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000220630&livemode=1&stbId=3",
		"wasusyt/6000000006000010630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000010630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000010630&livemode=1&stbId=3",
		"wasusyt/6000000006000050630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000050630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000050630&livemode=1&stbId=3",
		"wasusyt/6000000006000180630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000180630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000180630&livemode=1&stbId=3",
		"wasusyt/6000000006000030630.m3u8":     "http://gslbserv.itv.cmvideo.cn:80/6000000006000030630/1.m3u8?channel-id=wasusyt&Contentid=6000000006000030630&livemode=1&stbId=3",
	}

	dnsCache         = sync.Map{}
	successCacheTime = 24 * time.Hour
)

type cacheEntry struct {
	ip     string
	expiry time.Time
}

func (i *Itv) HandleMainRequest(w http.ResponseWriter, r *http.Request, cdn string, id string) {
	key := cdn + "/" + id
	startUrl, ok := programList[key]
	if !ok {
		http.Error(w, "id not found!", http.StatusNotFound)
		return
	}

	data, redirectURL, err := getHTTPResponse(startUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	redirectPrefix := redirectURL[:strings.LastIndex(redirectURL, "/")+1]

	// 替换TS文件的链接
	golang := "http://" + r.Host + r.URL.Path
	re := regexp.MustCompile(`((?i).*?\.ts)`)
	data = re.ReplaceAllStringFunc(data, func(match string) string {
		return golang + "?ts=" + redirectPrefix + match
	})

	// 将&替换为$
	data = strings.ReplaceAll(data, "&", "$")

	w.Header().Set("Content-Disposition", "attachment;filename="+id)
	w.WriteHeader(http.StatusOK) // Set the status code to 200
    w.Write([]byte(data)) // Write the response body
}

func (i *Itv) HandleTsRequest(w http.ResponseWriter, ts string) {
	// 将$替换回&
	ts = strings.ReplaceAll(ts, "$", "&")

	w.Header().Set("Content-Type", "video/MP2T")
	content, _, err := getHTTPResponse(ts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK) // Set the status code to 200
    w.Write([]byte(content)) // Write the response body
}

func getHTTPResponse(requestURL string) (string, string, error) {
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	var mappedHost string

	// 自定义resolver
	resolver := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			for originalHost, host := range hostMappings {
				if strings.Contains(address, originalHost) {
					ip := resolveIP(host)
					mappedHost = host
					if ip != "" {
						address = strings.Replace(address, originalHost, ip, 1)
					}
				}
			}
			return dialer.DialContext(ctx, network, address)
		},
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: resolver.Dial,
		},
	}

	resp, err := client.Get(requestURL)
	if err != nil {
		if mappedHost != "" {
			clearCache(mappedHost) // 清除缓存失败的IP
		}
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if mappedHost != "" {
			clearCache(mappedHost) // 请求失败清除缓存
		}
		return "", "", err
	}

	redirectURL := resp.Header.Get("Location")
	if redirectURL == "" {
		redirectURL = requestURL
	}

	body, err := readResponseBody(resp)
	if err != nil {
		if mappedHost != "" {
			clearCache(mappedHost) // 读取响应体失败时清除缓存
		}
		return "", "", err
	}

	if mappedHost != "" {
		updateCacheTime(mappedHost, successCacheTime) // 成功获取响应后缓存IP
	}

	return body, redirectURL, nil
}

func resolveIP(host string) string {
	now := time.Now()
	if entry, found := dnsCache.Load(host); found {
		cachedEntry := entry.(cacheEntry)
		if now.Before(cachedEntry.expiry) {
			return cachedEntry.ip // 使用缓存中的IP
		}
		dnsCache.Delete(host) // 缓存过期，删除
	}

	ips, err := net.LookupIP(host) // DNS解析
	if err != nil || len(ips) == 0 {
		return "" // 解析失败，返回空字符串
	}

	ip := ips[0].String()
	dnsCache.Store(host, cacheEntry{ip: ip, expiry: now.Add(successCacheTime)}) // 缓存解析到的IP
	return ip
}

func updateCacheTime(host string, duration time.Duration) {
	if entry, found := dnsCache.Load(host); found {
		cachedEntry := entry.(cacheEntry)
		cachedEntry.expiry = time.Now().Add(duration) // 更新缓存过期时间
		dnsCache.Store(host, cachedEntry)
	}
}

func clearCache(host string) {
	dnsCache.Delete(host) // 删除缓存
}

func readResponseBody(resp *http.Response) (string, error) {
	var builder strings.Builder
	_, err := io.Copy(&builder, resp.Body)
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}