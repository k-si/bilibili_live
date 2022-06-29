package util

import (
	"github.com/k-si/bili_live/config"
	"github.com/skip2/go-qrcode"
	"log"
)

// 根据url生成二维码
func GenerateQrcode(url string) error {
	if err := qrcode.WriteFile(url, qrcode.Medium, 256, config.Live.QrCodePath); err != nil {
		log.Println("生成二维码失败：", err)
		return err
	}
	return nil
}
