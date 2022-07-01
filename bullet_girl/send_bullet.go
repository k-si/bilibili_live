package bullet_girl

import (
	"context"
	"github.com/k-si/bili_live/http"
	"log"
	"time"
)

var sender *BulletSender

type BulletSender struct {
	bulletChan chan string
}

func PushToBulletSender(bullet string) {
	log.Println("PushToBulletSender成功", bullet)
	sender.bulletChan <- bullet
}

func StartSendBullet(ctx context.Context) {
	var err error

	sender = &BulletSender{
		bulletChan: make(chan string, 1000),
	}

	var msg string
	for {
		select {
		case <-ctx.Done():
			goto END
		case msg = <-sender.bulletChan:
			if err = http.Send(msg); err != nil {
				log.Println("弹幕发送失败：", err, "msg:", msg)
			} else {
				log.Println("弹幕发送成功：", msg)
			}
		}
		time.Sleep(3 * time.Second) // 防止弹幕发送过快
	}
END:
}
