//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"starlights/app/wechatbot/internal/biz"
	"starlights/app/wechatbot/internal/conf"
	"starlights/app/wechatbot/internal/data"
	"starlights/app/wechatbot/internal/server"
	"starlights/app/wechatbot/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Bot, log.Logger, []interface{}) (ServiceApp, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
