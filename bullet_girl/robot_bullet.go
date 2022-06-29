package bullet_girl

import (
	"context"
	"github.com/k-si/bili_live/config"
	"github.com/k-si/bili_live/http"
	"log"
	"strings"
)

var robot *BulletRobot

type BulletRobot struct {
	bulletRobotChan chan string
}

func PushToBulletRobot(content string) {
	log.Println("PushToBulletRobot成功", content)
	robot.bulletRobotChan <- content
}

func StartBulletRobot(ctx context.Context) {
	robot = &BulletRobot{
		bulletRobotChan: make(chan string, 1000),
	}

	var content string

	for {
		select {
		case <-ctx.Done():
			goto END
		case content = <-robot.bulletRobotChan:
			handleRobotBullet(content)
		}
	}
END:
}

func handleRobotBullet(content string) {
	var err error
	var reply string
	if reply, err = http.RequestQingyunkeRobot(content); err != nil {
		log.Println("请求机器人失败：", err)
		PushToBulletSender("不好意思，机器人坏掉了...")
		return
	}

	log.Println("机器人回复：", reply)

	bulltes := splitRobotReply(reply)
	for _, v := range bulltes {
		PushToBulletSender(v)
	}
}

// 将机器人回复语句中的 {br} 进行分割
// b站弹幕一次只能发20个字符，需要切分
func splitRobotReply(content string) []string {

	// 将机器人回复中的菲菲替换为橘子
	content = strings.ReplaceAll(content, "菲菲", config.Live.RobotName)

	var res []string
	reply := strings.Split(content, "{br}")

	for _, r := range reply {
		// 长度大于20再分割
		zh := []rune(r)
		if len(zh) > 20 {
			i := 0
			for i < len(zh) {
				if i+20 > len(zh) {
					res = append(res, string(zh[i:]))
				} else {
					res = append(res, string(zh[i:i+20]))
				}
				i += 20
			}
		} else {
			res = append(res, string(zh))
		}
	}
	return res
}

// 检查弹幕是否在@我，返回bool和@我要说的内容
func checkIsAtMe(msg string) (bool, string) {
	if strings.HasPrefix(msg, config.Live.TalkRobotCmd) {
		return true, strings.TrimPrefix(msg, config.Live.TalkRobotCmd)
	} else {
		return false, ""
	}
}
