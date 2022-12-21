package main

import (
	"flag"
	"github.com/eatmoreapple/openwechat"
	"os"
	"starlights/app/wechatbot/internal/conf"
	"starlights/app/wechatbot/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

// 服务列表
type serviceApp struct {
	app *kratos.App
	bot *openwechat.Bot
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, bot *openwechat.Bot) serviceApp {
	return serviceApp{
		app: kratos.New(
			kratos.ID(id),
			kratos.Name(Name),
			kratos.Version(Version),
			kratos.Metadata(map[string]string{}),
			kratos.Logger(logger),
			kratos.Server(
				gs,
				hs,
			),
		),
		bot: bot,
	}
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	serviceApp, cleanup, err := wireApp(bc.Server, bc.Data, bc.Bot, logger, service.BotRegister())
	if err != nil {
		panic(err)
	}
	defer cleanup()
	serviceApp.bot.Login()
	// start and wait for stop signal
	if err := serviceApp.app.Run(); err != nil {
		panic(err)
	}
}
