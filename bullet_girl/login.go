package bullet_girl

import (
	"github.com/k-si/bili_live/entity"
	"github.com/k-si/bili_live/http"
	"github.com/k-si/bili_live/util"
	"log"
)

// 先登录，获取cookie
func UserLogin() error {
	var err error
	var loginUrl *entity.LoginUrl

	if loginUrl, err = http.GetLoginUrl(); err != nil {
		log.Println(err)
		return err
	}

	if err = util.GenerateQrcode(loginUrl.Data.Url); err != nil {
		log.Println(err)
		return err
	}

	if _, err = http.GetLoginInfo(loginUrl.Data.OauthKey); err != nil {
		log.Println(err)
		return err
	}

	return err
}
