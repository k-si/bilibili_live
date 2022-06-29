package entity

const (
	NotStarted = 0 // 未开播
	Live       = 1 // 直播中
	Carousel   = 2 // 轮播中
)

type RoomInitInfo struct {
	Data struct {
		LiveStatus int `json:"live_status"`
	} `json:"data"`
}
