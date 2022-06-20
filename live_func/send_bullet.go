package live_func

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

var cookie = "buvid3=7D3E791A-51D3-9950-66B5-70EC2859E37C34776infoc; rpdid=|(k|~uJJR||R0J'uYkRk|lkJY; LIVE_BUVID=AUTO8516339587249653; fingerprint_s=f4b605d68b1443938c5fb47d8eaa8127; video_page_version=v_old_home; buvid_fp_plain=A0E922DE-45D1-E8F3-DD84-2BCBEF26F2F482616infoc; DedeUserID=431786226; DedeUserID__ckMd5=ad84921fc5cb759b; i-wanna-go-back=-1; b_ut=5; CURRENT_BLACKGAP=0; buvid4=2F62FA57-3CEA-E0A5-BBEF-BC25F19EB04556345-022012422-c9eDNxIrTV/XPspg2ZQn3w%3D%3D; buvid_fp=fbf8ae46ddecf8162b4318193599f39f; nostalgia_conf=-1; CURRENT_QUALITY=80; hit-dyn-v2=1; is-2022-channel=1; fingerprint3=1941bf3bbb7a5b4c50a3288441f7ac6d; fingerprint=c2de2dc3e4b8ffad1d11f628e67cb15b; blackside_state=0; Hm_lvt_8a6e55dbd2870f0f5bc9194cddf32a02=1653625372,1655216528; Hm_lpvt_8a6e55dbd2870f0f5bc9194cddf32a02=1655652561; bp_video_offset_431786226=673736637856874500; innersign=1; CURRENT_FNVAL=4048; b_lsid=911BFC8C_18181D7A3B7; _uuid=859D6A8E-44AF-B9F2-87B4-101826E52A18424559infoc; _dfcaptcha=a24e4b89a61ba7a90fa111271c741c10; b_timer=%7B%22ffp%22%3A%7B%22333.1007.fp.risk_7D3E791A%22%3A%221817FAE8AE5%22%2C%22333.967.fp.risk_7D3E791A%22%3A%221817C6A0DA7%22%2C%22444.8.fp.risk_7D3E791A%22%3A%2218181F9C49E%22%2C%22333.788.fp.risk_7D3E791A%22%3A%22181820D3FF3%22%2C%22333.937.fp.risk_7D3E791A%22%3A%221816295BD93%22%2C%22333.337.fp.risk_7D3E791A%22%3A%22181620AAF6E%22%2C%22333.999.fp.risk_7D3E791A%22%3A%2218181F88397%22%2C%22444.41.fp.risk_7D3E791A%22%3A%221817CFB0F40%22%2C%22444.1001.fp.risk_7D3E791A%22%3A%2218181EBD24B%22%2C%22333.976.fp.risk_7D3E791A%22%3A%22181677EDD0B%22%2C%22444.45.fp.risk_7D3E791A%22%3A%2218177DA055A%22%2C%22444.60.0.0.fp.risk_7D3E791A%22%3A%221816DAF34F2%22%2C%22444.42.fp.risk_7D3E791A%22%3A%221817C90104A%22%2C%22333.52.fp.risk_7D3E791A%22%3A%221817C91A320%22%7D%7D; SESSDATA=e3eefbd2%2C1671296318%2C612eb%2A61; bili_jct=b6e5be9bc8dc11572a52d3828b0a0616; sid=izvvfwc5; bsource=search_baidu; PVID=5"
var origin = "https://live.bilibili.com"
var referer = "https://live.bilibili.com/25198571?broadcast_type=0&is_room_feed=1&spm_id_from=333.999.0.0"
var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Safari/537.36"

func SendBullet(msg string) {
	cli := resty.New()

	m := make(map[string]string)
	m["bubble"] = "5"
	m["msg"] = msg
	m["color"] = "4546550"
	m["mode"] = "4"
	m["fontsize"] = "25"
	m["rnd"] = strconv.FormatInt(time.Now().Unix(), 10)
	m["roomid"] = "25198571"
	m["csrf"] = "b6e5be9bc8dc11572a52d3828b0a0616"
	m["csrf_token"] = "b6e5be9bc8dc11572a52d3828b0a0616"
	fmt.Println(strconv.FormatInt(time.Now().Unix(), 10))

	_, err := cli.R().
		SetHeader("cookie", cookie).
		SetHeader("origin", origin).
		SetHeader("referer", referer).
		SetHeader("user_agent", userAgent).
		SetFormData(m).
		Post("https://api.live.bilibili.com/msg/send")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("send success")
}
