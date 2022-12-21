package server

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/go-kratos/kratos/v2/log"
	"starlights/app/wechatbot/internal/conf"
)

func NewBotServer(c *conf.Bot, logger log.Logger, register []interface{}) *openwechat.Bot {
	var bot *openwechat.Bot
	switch c.GetMode() {
	case 0:
		bot = openwechat.DefaultBot(openwechat.Normal)
	case 1:
		bot = openwechat.DefaultBot(openwechat.Normal)
	}
	for _, o := range register {
		if v, ok := o.(openwechat.MessageHandler); ok {
			bot.MessageHandler = v
		}
		if v, ok := o.(func(string)); ok {
			bot.UUIDCallback = v
		}
	}
	return bot
}
