package main

import (
	"context"
	"github.com/k-si/bili_live/live_func"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// 开启定时弹幕任务
	bulletCtx, bulletCancel := context.WithCancel(context.Background())
	defer bulletCancel()
	go live_func.StartTimingBullet(bulletCtx)

	// 推送一条弹幕定时任务
	b := live_func.NewBullet("ios请到哔哩哔哩直播姬公众号进行投喂哦", "* */2 * * * * *")
	bt := live_func.NewBulletTask(b)
	be := live_func.NewBulletEvent(live_func.Save, bt)
	live_func.PushBulletEvent(be)

	// 阻塞等待
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
