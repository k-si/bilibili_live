package bullet_girl

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/k-si/bili_live/entity"
	"io"
	"log"
	"strings"
)

var handler *BulletHandler

type BulletHandler struct {
	BulletChan chan []byte
}

func pushToBulletHandler(message []byte) {
	handler.BulletChan <- message
}

func HandleBullet(ctx context.Context) {
	handler = &BulletHandler{
		BulletChan: make(chan []byte, 1000),
	}

	var message []byte
	for {
		select {
		case <-ctx.Done():
			goto END
		case message = <-handler.BulletChan:
			handle(message)
		}
	}
END:
}

func handle(message []byte) {
	var err error

	// 一个正文可能包含多个数据包，需要逐个解析
	index := 0
	for index < len(message) {

		// 读出包长
		var length uint32
		if err = binary.Read(bytes.NewBuffer(message[index:index+headLengthOffset]), binary.BigEndian, &length); err != nil {
			log.Println("解析包长度失败", err)
			return
		}

		// 读出正文协议版本
		var ver Version
		if err = binary.Read(bytes.NewBuffer(message[index+versionOffset:index+opcodeOffset]), binary.BigEndian, &ver); err != nil {
			log.Println("解析正文协议版本失败", err)
			return
		}

		// 读出操作码
		var op Opcode
		if err = binary.Read(bytes.NewBuffer(message[index+opcodeOffset:index+magicOffset]), binary.BigEndian, &op); err != nil {
			log.Println("解析操作码失败", err)
			return
		}

		// 读出正文内容
		body := message[index+packageLength : index+int(length)]

		// 解析正文内容
		switch ver {
		case normalJson:
			log.Println("普通json包：", string(body), ver, op)
			text := &entity.CmdText{}
			_ = json.Unmarshal(body, text)
			if op == command {
				switch Cmd(text.Cmd) {

				// 处理弹幕
				case DanmuMsg:
					danmu := &entity.DanmuMsgText{}
					_ = json.Unmarshal(body, danmu)
					from := danmu.Info[2].([]interface{})

					// 如果发现弹幕在@我，那么调用机器人进行回复
					y, content := checkIsAtMe(danmu.Info[1].(string))
					if y {
						PushToBulletRobot(content)
					}

					log.Println(from[0].(float64), from[1], "：", danmu.Info[1])

				// 欢迎舰长
				case entryEffect:
					entry := &entity.EntryEffectText{}
					_ = json.Unmarshal(body, entry)
					PushToBulletSender(welcomeCaptain(entry.Data.CopyWriting))

				// 欢迎进入房间
				case interactWord:
					interact := &entity.InteractWordText{}
					_ = json.Unmarshal(body, interact)
					PushToBulletSender(welcomeInteract(interact.Data.Uname))
				}
			}
		case heartOrCertification:
			log.Println("心跳回复包")
		case normalZlib:
			b := bytes.NewReader(body)
			r, _ := zlib.NewReader(b)
			var out bytes.Buffer
			_, _ = io.Copy(&out, r)
			handle(out.Bytes()) // zlib解压后再进行格式解析
		}
		index += int(length)
	}
}

// 欢迎舰长语句
func welcomeCaptain(s string) string {

	// 判断是否舰长
	b := strings.Contains(s, "舰长")

	// 获取名称
	zh := []rune(s)
	zh = []rune(strings.TrimPrefix(s, "<% "))
	zh = []rune(strings.TrimRight(s, " %>"))

	if b {
		return "欢迎舰长：" + string(zh)
	} else {
		return "欢迎" + string(zh)
	}
}

func welcomeInteract(name string) string {
	if strings.Contains(name, "欢迎") {
		return name
	} else {
		return "欢迎" + name
	}
}
