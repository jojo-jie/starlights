package main

import (
	"errors"
	"flag"
	"github.com/eatmoreapple/openwechat"
	"golang.org/x/sync/errgroup"
	"os"
	"starlights/app/wechatbot/internal/conf"
	"starlights/app/wechatbot/internal/service"
	"time"

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

// ServiceApp 服务列表
type ServiceApp struct {
	app *kratos.App
	Bot *openwechat.Bot
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, bot *openwechat.Bot) ServiceApp {
	return ServiceApp{
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
		Bot: bot,
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

	mService, cleanup, err := wireApp(bc.Server, bc.Data, bc.Bot, logger, service.BotRegister())
	if err != nil {
		panic(err)
	}
	defer cleanup()
	g := errgroup.Group{}
	g.Go(func() error {
		storage := openwechat.NewJsonFileHotReloadStorage("storage.json")
		t := time.NewTicker(2 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				if err := mService.Bot.HotLogin(storage, true); err != nil {
					if !errors.Is(err, openwechat.ErrLoginTimeout) {
						return err
					}
				} else {
					t.Stop()
					break
				}
			}
		}
		return mService.Bot.Block()
	})
	g.Go(func() error {
		// start and wait for stop signal
		return mService.app.Run()
	})
	if err := g.Wait(); err != nil {
		panic(err)
	}
}
