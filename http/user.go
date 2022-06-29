package http

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/k-si/bili_live/entity"
	"log"
	"strings"
	"time"
)

func GetLoginUrl() (*entity.LoginUrl, error) {
	var err error
	var resp *resty.Response
	var url = "https://passport.bilibili.com/qrcode/getLoginUrl"

	r := &entity.LoginUrl{}
	if resp, err = cli.R().
		SetHeader("user-agent", userAgent).
		Get(url); err != nil {
		log.Println("请求getLoginUrl失败：", err)
		return nil, err
	}
	if err = json.Unmarshal(resp.Body(), r); err != nil {
		log.Println("Unmarshal失败：", err, "body:", string(resp.Body()))
		return nil, err
	}

	log.Println("oauthKey:", r.Data.OauthKey)

	return r, err
}

// 验证登录的同时，将cookie赋值
func GetLoginInfo(oauthKey string) (*entity.LoginInfoData, error) {
	var err error
	var url = "https://passport.bilibili.com/qrcode/getLoginInfo?oauthKey=" + oauthKey
	var resp *resty.Response
	var data *entity.LoginInfoData

	pre := &entity.LoginInfoPre{}

	for {
		log.Println("等待扫码登录...")

		if resp, err = cli.R().
			SetHeader("user-agent", userAgent).
			Post(url); err != nil {
			log.Println("请求getLoginInfo失败：", err)
			return nil, err
		}

		if err = json.Unmarshal(resp.Body(), pre); err != nil {
			log.Println("Unmarshal失败：", err, "body:", string(resp.Body()))
			return nil, err
		}

		if pre.Status {
			data = &entity.LoginInfoData{}
			if err = json.Unmarshal(resp.Body(), data); err != nil {
				log.Println("Unmarshal失败：", err, "body:", string(resp.Body()))
				return nil, err
			}
			log.Println("登录成功！")
			break
		}

		time.Sleep(5 * time.Second)
	}

	for _, v := range resp.Header().Values("Set-Cookie") {
		pair := strings.Split(v, ";")
		kv := strings.Split(pair[0], "=")
		CookieList[kv[0]] = kv[1]
		CookieStr += pair[0] + ";"
	}

	return data, err
}
