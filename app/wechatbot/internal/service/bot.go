package service

import (
	"github.com/eatmoreapple/openwechat"
)

func BotRegister() []interface{} {
	var i []interface{}
	i = append(i, func(msg *openwechat.Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	})
	i = append(i, func(uuid string) {
		qrcodeUrl := openwechat.GetQrcodeUrl(uuid)
		println(qrcodeUrl)
		// send email
		go func() {
			println("send email")
		}()
	})
	return i
}

func msgTemplate() {

}
