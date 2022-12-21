package service

import (
	"github.com/eatmoreapple/openwechat"
	qrcodeS "github.com/skip2/go-qrcode"
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
		// browser open the login url
		q, _ := qrcodeS.New(qrcodeUrl, qrcodeS.Low)
		println(q.ToSmallString(true))
		// send email
		go func() {
			println("send email")
		}()
	})
	return i
}
