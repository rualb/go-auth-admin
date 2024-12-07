package router

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"

	"go-auth-admin/internal/config/consts"

	"go-auth-admin/internal/controller/authadmin"

	"go-auth-admin/internal/service"
	webfs "go-auth-admin/web"

	xlog "go-auth-admin/internal/util/utillog"
	xweb "go-auth-admin/internal/web"

	"github.com/labstack/echo/v4/middleware"
)

func Init(e *echo.Echo, appService service.AppService) {

	e.Renderer = mustNewRenderer()

	initCORSConfig(e, appService)

	initAuthAdminController(e, appService)
	initDebugController(e, appService)

	initSys(e, appService)
}

func initSys(e *echo.Echo, appService service.AppService) {

	// !!! DANGER for private(non-public) services only
	// or use non-public port via echo.New()

	appConfig := appService.Config()

	listen := appConfig.HTTPServer.Listen
	listenSys := appConfig.HTTPServer.ListenSys
	sysMetrics := appConfig.HTTPServer.SysMetrics
	hasAnyService := sysMetrics
	sysAPIKey := appConfig.HTTPServer.SysAPIKey
	hasAPIKey := sysAPIKey != ""
	hasListenSys := listenSys != ""
	startNewListener := listenSys != listen

	if !hasListenSys {
		return
	}

	if !hasAnyService {
		return
	}

	if !hasAPIKey {
		xlog.Panic("Sys api key is empty")
		return
	}

	if startNewListener {

		e = echo.New() // overwrite override

		e.Use(middleware.Recover())
		// e.Use(middleware.Logger())
	} else {
		xlog.Warn("Sys api serve in main listener: %v", listen)
	}

	sysAPIAccessAuthMW := middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "query:api-key,header:Authorization",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == sysAPIKey, nil
		},
	})

	if sysMetrics {
		// may be eSys := echo.New() // this Echo will run on separate port
		e.GET(
			consts.PathSysMetricsAPI,
			echoprometheus.NewHandler(),
			sysAPIAccessAuthMW,
		) // adds route to serve gathered metrics

	}

	if startNewListener {

		// start as async task
		go func() {
			xlog.Info("Sys api serve on: %v main: %v", listenSys, listen)

			if err := e.Start(listenSys); err != nil {
				if err != http.ErrServerClosed {
					xlog.Error("%v", err)
				} else {
					xlog.Info("shutting down the server")
				}
			}
		}()

	} else {
		xlog.Info("Sys api server serve on main listener: %v", listen)
	}

}

type tmplRenderer struct {
	indexHTML *template.Template
}

func (x *tmplRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {

	if name == "index.html" {

		return x.indexHTML.ExecuteTemplate(w, name, data)
	}

	return fmt.Errorf("render page not found: %v", name)

}

func mustNewRenderer() echo.Renderer {

	indexHTML, err := template.New("index.html").Parse(webfs.MustAuthAdminIndexHTML())

	if err != nil {
		panic(err)
	}

	//	err := t.templates.ExecuteTemplate(w, "layout_header", data)

	handler := &tmplRenderer{
		// viewsMvc:  mvc.NewTemplateRenderer(controller.ViewsFs(), "views/auth-admin/*.html"),
		indexHTML: indexHTML,
	}

	return handler

}

func initCORSConfig(e *echo.Echo, _ service.AppService) {

	// CorsEnabled := true
	// if CorsEnabled {
	// 	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 		AllowCredentials:                         true,
	// 		UnsafeWildcardOriginWithAllowCredentials: true,
	// 		AllowOrigins:                             []string{"*"},
	// 		MaxAge:                                   86400,
	// 	}))
	// }
}

func initDebugController(e *echo.Echo, _ service.AppService) {

	e.GET(consts.PathAuthAdminPingDebugAPI, func(c echo.Context) error { return c.String(http.StatusOK, "pong") })
	// publicly-available-no-sensitive-data
	e.GET("/health", func(c echo.Context) error { return c.JSON(http.StatusOK, struct{}{}) })

}

// ///////////////////////////////////////////////////
func initAuthAdminController(e *echo.Echo, appService service.AppService) {

	{

		xlog.Warn("Adding auth admin controllers")

		prefix := consts.PathAuthAdmin
		group := e.Group(prefix)

		path := func(s string) string {
			xlog.Info("Route: %s", s)
			return strings.TrimPrefix(s, prefix)
		}

		{
			{
				// auth
				group.Use(xweb.AuthorizeMiddlewareWithConfig(xweb.AuthorizeMiddlewareConfig{
					Service:      appService,
					IfAnyOfRoles: authadmin.RolesForAPI,
				}))
			}

			{
				group.GET(path(consts.PathAuthAdminStatusAPI), func(c echo.Context) error {
					ctrl := authadmin.NewStatusAPIController(appService, c)
					return ctrl.Handler()
				})
				group.GET(path(consts.PathAuthAdminConfigAPI), func(c echo.Context) error {
					ctrl := authadmin.NewConfigAPIController(appService, c)
					return ctrl.Handler()
				})
			}

			{
				// return UI
				handler := func(c echo.Context) error {
					ctrl := authadmin.NewAuthAdminIndexController(appService, c)
					return ctrl.Handler()
				}

				group.GET(path(consts.PathAuthAdminAccounts), handler)
				group.GET(path(consts.PathAuthAdminAccountsEntity), handler)
			}

			{
				{
					handler := func(c echo.Context) error {
						ctrl := authadmin.NewAccountsAPIController(appService, c)
						return ctrl.Handler()
					}

					group.GET(path(consts.PathAuthAdminAccountsAPI), handler)

				}
				{
					handler := func(c echo.Context) error {
						ctrl := authadmin.NewAccountsEntityAPIController(appService, c)
						return ctrl.Handler()
					}

					group.GET(path(consts.PathAuthAdminAccountsEntityByCodeAPI), handler)
					group.GET(path(consts.PathAuthAdminAccountsEntityAPI), handler)
					group.POST(path(consts.PathAuthAdminAccountsAPI), handler) // no :id
					group.PUT(path(consts.PathAuthAdminAccountsEntityAPI), handler)
					group.DELETE(path(consts.PathAuthAdminAccountsEntityAPI), handler)

				}
				{

					handler := func(c echo.Context) error {
						ctrl := authadmin.NewAccountsPasswordAPIController(appService, c)
						return ctrl.Handler()
					}

					group.POST(path(consts.PathAuthAdminAccountsEntityPasswordAPI), handler)

				}

			}
		}

	}
}
