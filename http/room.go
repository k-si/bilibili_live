package http

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/k-si/bili_live/config"
	"github.com/k-si/bili_live/entity"
	"log"
	"strconv"
)

func RoomInit() (*entity.RoomInitInfo, error) {
	var err error
	var resp *resty.Response
	var url = "https://api.live.bilibili.com/room/v1/Room/room_init?id=" + strconv.Itoa(config.Live.RoomId)

	if resp, err = cli.R().
		SetHeader("user-agent", userAgent).
		Get(url); err != nil {
		log.Println("请求room_init失败：", err)
		return nil, err
	}
	r := &entity.RoomInitInfo{}
	if err = json.Unmarshal(resp.Body(), r); err != nil {
		log.Println("Unmarshal失败：", err, "body:", string(resp.Body()))
		return nil, err
	}

	return r, err
}
