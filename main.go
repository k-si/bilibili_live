package main

import (
	"context"
	"github.com/k-si/bili_live/bullet_girl"
	"github.com/k-si/bili_live/config"
	"github.com/k-si/bili_live/entity"
	"github.com/k-si/bili_live/http"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var sendBulletCtx context.Context
var sendBulletCancel context.CancelFunc
var timingBulletCtx context.Context
var timingBulletCancel context.CancelFunc
var robotBulletCtx context.Context
var robotBulletCancel context.CancelFunc
var catchBulletCtx context.Context
var catchBulletCancel context.CancelFunc
var handleBulletCtx context.Context
var handleBulletCancel context.CancelFunc
var thanksGiftCtx context.Context
var thankGiftCancel context.CancelFunc

func main() {
	var err error

	defer func() {
		if sendBulletCancel != nil {
			sendBulletCancel()
		}
		if timingBulletCancel != nil {
			timingBulletCancel()
		}
		if robotBulletCancel != nil {
			robotBulletCancel()
		}
		if catchBulletCancel != nil {
			catchBulletCancel()
		}
		if handleBulletCancel != nil {
			handleBulletCancel()
		}
		if thankGiftCancel != nil {
			thankGiftCancel()
		}
	}()

	// 初始化配置文件，http客户端
	if err = config.InitConfig(); err != nil {
		log.Fatal("配置文件错误：", err)
	}
	http.InitHttpClient()

	// 扫码登录
	if err = bullet_girl.UserLogin(); err != nil {
		log.Fatal("用户登录失败：", err)
	}

	// 准备select中用到的变量
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	var interval = time.Minute
	t := time.NewTimer(interval)
	defer t.Stop()
	var info *entity.RoomInitInfo
	var preStatus int

	log.Println("正在检测直播间是否开播...")

	// 循环监听直播间情况
	for {
		select {

		// 程序退出
		case <-sig:
			goto END

		// 每1分钟检查一次直播间是否开播
		case <-t.C:
			t.Reset(interval)
			if info, err = http.RoomInit(); err != nil {
				log.Println("http请求错误：", err)
			}
			if info.Data.LiveStatus == entity.Live && preStatus == entity.NotStarted { // 由NotStarted到Live是开播
				log.Println("开播啦！")
				preStatus = entity.Live
				// 弹幕姬各goroutine上下文
				sendBulletCtx, sendBulletCancel = context.WithCancel(context.Background())
				timingBulletCtx, timingBulletCancel = context.WithCancel(context.Background())
				robotBulletCtx, robotBulletCancel = context.WithCancel(context.Background())
				catchBulletCtx, catchBulletCancel = context.WithCancel(context.Background())
				handleBulletCtx, handleBulletCancel = context.WithCancel(context.Background())
				thanksGiftCtx, thankGiftCancel = context.WithCancel(context.Background())

				StartBulletGirl(sendBulletCtx, timingBulletCtx, robotBulletCtx, catchBulletCtx, handleBulletCtx, thanksGiftCtx) // 开启弹幕姬
			} else if info.Data.LiveStatus == entity.NotStarted && preStatus == entity.Live { // 由Live到NotStarted是下播
				log.Println("下播啦！")
				preStatus = entity.NotStarted
				sendBulletCancel()
				timingBulletCancel()
				robotBulletCancel()
				catchBulletCancel()
				handleBulletCancel()
				thankGiftCancel()// 关闭弹幕姬goroutine
			}
		}
	}
END:
}

func StartBulletGirl(sendBulletCtx, timingBulletCtx, robotBulletCtx, catchBulletCtx, handleBulletCtx, thanksGiftCtx context.Context) {

	// 开启弹幕推送
	go bullet_girl.StartSendBullet(sendBulletCtx)
	log.Println("弹幕推送已开启...")

	// 开启定时弹幕任务
	go bullet_girl.StartTimingBullet(timingBulletCtx)
	log.Println("定时弹幕已开启...")

	// 开启弹幕机器人
	go bullet_girl.StartBulletRobot(robotBulletCtx)
	log.Println("弹幕机器人已开启")

	// 开启弹幕抓取
	go bullet_girl.StartCatchBullet(catchBulletCtx)
	log.Println("弹幕抓取已开启...")

	// 开启弹幕处理
	go bullet_girl.HandleBullet(handleBulletCtx)
	log.Println("弹幕处理已开启...")

	// 开启礼物感谢
	go bullet_girl.ThanksGift(thanksGiftCtx)
	log.Println("礼物感谢已开启")

	// 指定弹幕定时任务
	time.Sleep(time.Second) // 现开启定时任务弹幕再推送，这个方法很low，暂且这样吧
	bullet_girl.PushToBulletEvent(
		bullet_girl.NewBulletEvent(
			bullet_girl.Save, bullet_girl.NewBulletTask(
				bullet_girl.NewBullet("ios请到哔哩哔哩直播姬公众号投喂哦～", "*/9 * * * * *"))))
	bullet_girl.PushToBulletEvent(
		bullet_girl.NewBulletEvent(
			bullet_girl.Save, bullet_girl.NewBulletTask(
				bullet_girl.NewBullet("喜欢主播可以加入粉丝团哦～", "*/7 * * * * *"))))
	bullet_girl.PushToBulletEvent(
		bullet_girl.NewBulletEvent(
			bullet_girl.Save, bullet_girl.NewBulletTask(
				bullet_girl.NewBullet("主播今天很可爱哦！干巴爹！", "*/17 * * * * *"))))
	bullet_girl.PushToBulletEvent(
		bullet_girl.NewBulletEvent(
			bullet_girl.Save, bullet_girl.NewBulletTask(
				bullet_girl.NewBullet("无聊的同学可以找橘子聊天喔！", "*/23 * * * * *"))))
}

func main_test() {
	sendBulletCtx, sendBulletCancel = context.WithCancel(context.Background())
	timingBulletCtx, timingBulletCancel = context.WithCancel(context.Background())
	robotBulletCtx, robotBulletCancel = context.WithCancel(context.Background())
	catchBulletCtx, catchBulletCancel = context.WithCancel(context.Background())
	handleBulletCtx, handleBulletCancel = context.WithCancel(context.Background())
	thanksGiftCtx, thankGiftCancel = context.WithCancel(context.Background())

	var err error

	defer func() {
		if sendBulletCancel != nil {
			sendBulletCancel()
		}
		if timingBulletCancel != nil {
			timingBulletCancel()
		}
		if robotBulletCancel != nil {
			robotBulletCancel()
		}
		if catchBulletCancel != nil {
			catchBulletCancel()
		}
		if handleBulletCancel != nil {
			handleBulletCancel()
		}
		if thankGiftCancel != nil {
			thankGiftCancel()
		}
	}()

	// 初始化配置文件，http客户端
	if err = config.InitConfig(); err != nil {
		log.Fatal("配置文件错误：", err)
	}
	http.InitHttpClient()

	// 扫码登录
	if err = bullet_girl.UserLogin(); err != nil {
		log.Fatal("用户登录失败：", err)
	}

	StartBulletGirl(sendBulletCtx, timingBulletCtx, robotBulletCtx, catchBulletCtx, handleBulletCtx, thanksGiftCtx) // 开启弹幕姬

	// 准备select中用到的变量
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
