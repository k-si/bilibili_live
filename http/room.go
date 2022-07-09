package http

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/k-si/bili_live/config"
	"github.com/k-si/bili_live/entity"
	"github.com/k-si/bili_live/errs"
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

	// 先解析响应状态
	status := &entity.RoomInitStatus{}
	if err = json.Unmarshal(resp.Body(), status); err != nil {
		log.Println("Unmarshal失败：", err, "body:", string(resp.Body()))
		return nil, err
	}

	// 在解析房间状态
	r := &entity.RoomInitInfo{}
	if status.Code == 0 {
		if err = json.Unmarshal(resp.Body(), r); err != nil {
			log.Println("Unmarshal失败：", err, "body:", string(resp.Body()))
			return nil, err
		}
	}

	// 太长时间下播，房间号可能会消失，请求响应的code=60004
	if status.Code == 60004 {
		return nil, errs.RoomIdNotExistErr
	}

	return r, err
}
