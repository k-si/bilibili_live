package entity

type CmdText struct {
	Cmd string `json:"cmd"`
}

type DanmuMsgText struct {
	Info []interface{} `json:"info"`
}

type EntryEffectText struct {
	Data struct{
		CopyWriting string `json:"copy_writing"`
	} `json:"data"`
}

type InteractWordText struct {
	Data struct{
		Uname string `json:"uname"`
	} `json:"data"`
}