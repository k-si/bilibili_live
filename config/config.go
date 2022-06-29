package config

import (
	"flag"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"log"
)

const (
	RoomId       = 25198571
	WsServerUrl  = "wss://broadcastlv.chat.bilibili.com:2245/sub"
	QrCodePath   = "login_bili.png"
	TalkRobotCmd = "橘子，"
	RobotName    = "橘子"
)

var Live LiveConfig

type LiveConfig struct {
	RoomId       int    `toml:"room_id"`
	WsServerUrl  string `toml:"ws_server_url"`
	QrCodePath   string `toml:"qr_code_path"`
	TalkRobotCmd string `toml:"talk_robot_cmd"`
	RobotName    string `toml:"robot_name"`
}

func DefaultLiveConfig() LiveConfig {
	return LiveConfig{
		RoomId:       RoomId,
		WsServerUrl:  WsServerUrl,
		QrCodePath:   QrCodePath,
		TalkRobotCmd: TalkRobotCmd,
		RobotName:    RobotName,
	}
}

func LoadLiveConfig(path string) error {
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(path); err != nil {
		return err
	}
	if err = toml.Unmarshal(data, &Live); err != nil {
		return err
	}
	return err
}

func InitConfig() error {
	var err error
	c := flag.String("c", "", "configuration profile path")
	flag.Parse()

	// 加载配置文件
	if *c == "" {
		Live = DefaultLiveConfig()
	} else {
		if err = LoadLiveConfig(*c); err != nil {
			log.Fatal(err)
		}
	}
	return err
}
