package cmd

import (
	"context"
	"fmt"
	"go-auth-admin/internal/config"
	"go-auth-admin/internal/middleware"
	"go-auth-admin/internal/service"
	"net/http"
	"os"
	"os/signal"
	"time"

	xlog "go-auth-admin/internal/util/utillog"

	"go-auth-admin/internal/router"
	webfs "go-auth-admin/web"

	"github.com/labstack/echo/v4"
	elog "github.com/labstack/gommon/log"
)

type Command struct {
	AppService service.AppService
	WebDriver  *echo.Echo
	stop       context.CancelFunc
}

func (x *Command) Stop() {
	x.stop()
}

func (x *Command) Exec() {

	defer xlog.Sync()

	x.AppService = service.MustNewAppServiceProd()

	x.WebDriver = echo.New()
	x.WebDriver.Logger.SetLevel(elog.INFO) // has "file":"cmd.go","line":"85"

	middleware.Init(x.WebDriver, x.AppService) // 1
	router.Init(x.WebDriver, x.AppService)     // 2

	middleware.AssetsContentsMiddleware(x.WebDriver, x.AppService,
		webfs.MustAuthAdminAssetsFS(),
	)

	defer func() {
		xlog.Info("closing repository")
		_ = x.AppService.Repository().Close()
		xlog.Info("bye")
	}()

	x.startWithGracefulShutdown()

	time.Sleep(400 * time.Microsecond)
}

func applyServer(s *http.Server, c *config.AppConfig) {

	s.ReadTimeout = time.Duration(c.HTTPServer.ReadTimeout) * time.Second
	s.WriteTimeout = time.Duration(c.HTTPServer.WriteTimeout) * time.Second
	s.IdleTimeout = time.Duration(c.HTTPServer.IdleTimeout) * time.Second
	s.ReadHeaderTimeout = time.Duration(c.HTTPServer.ReadHeaderTimeout) * time.Second

}

func (x *Command) startWithGracefulShutdown() {

	appConfig := x.AppService.Config()

	listen := appConfig.HTTPServer.Listen
	// Graceful shutdown

	webDriver := x.WebDriver

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	x.stop = stop

	// Start server

	{
		applyServer(webDriver.Server, appConfig)
		applyServer(webDriver.TLSServer, appConfig)

		xlog.Info("server starting: %v", listen)

		go func() {

			defer func() {
				xlog.Info("server exiting")

				if r := recover(); r != nil {
					// Log or handle the panic
					panic(fmt.Errorf("error panic: %v", r))
				}
			}()

			if err := webDriver.Start(listen); err != nil {
				if err != http.ErrServerClosed {
					xlog.Error("%v", err)
				} else {
					xlog.Info("shutting down the server")
				}
			}

		}()

	}

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	xlog.Info("interrupt signal")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	xlog.Info("shutdown web driver")
	if err := webDriver.Shutdown(ctx); err != nil {
		xlog.Error("error on shutdown server: %v", err)
	}
}
