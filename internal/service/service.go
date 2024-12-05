package service

import (
	"go-auth-admin/internal/config"
	"go-auth-admin/internal/i18n"
	"go-auth-admin/internal/messenger"
	"go-auth-admin/internal/repository"
	xlog "go-auth-admin/internal/util/utillog"
	"net/http"
	"os"
	"strings"
	"time"
)

// AppService all services
type AppService interface {
	Account() AccountService

	Auth() AuthService
	AuthAdmin() AuthAdminService

	Config() *config.AppConfig
	// Logger() logger.AppLogger

	UserLang(code string) i18n.UserLang
	HasLang(code string) bool
	Messenger() messenger.AppMessenger

	Vault() VaultService

	Repository() repository.AppRepository
}

type defaultAppService struct {
	accountService AccountService
	// container      container.AppContainer
	vaultService VaultService

	configSource *config.AppConfigSource
	repository   repository.AppRepository
	lang         i18n.AppLang
	messenger    messenger.AppMessenger

	authService      AuthService
	authAdminService AuthAdminService
}

func (x *defaultAppService) mustConfig() {

	d, _ := os.Getwd()

	xlog.Info("Current work dir: %v", d)

	x.configSource = config.MustNewAppConfigSource()

	appConfig := x.configSource.Config() // first call, init

	mustConfigRuntime(appConfig)

}

func (x *defaultAppService) mustBuild() {

	var err error

	appConfig := x.configSource.Config() // first call, init

	x.messenger = messenger.NewAppMessenger(appConfig) // , appLogger)
	//
	x.lang = i18n.NewAppLang(appConfig)

	x.repository = repository.MustNewRepository(appConfig) // , appLogger)

	//

	mustCreateRepository(x)

	x.vaultService, err = newVaultService(x)

	if err != nil {
		panic(err)
	}

	x.authService = newAuthService(x)

	x.authAdminService = newAuthAdminService(x)

	x.accountService = newAccountService(x)

}

func mustConfigRuntime(appConfig *config.AppConfig) {
	t, ok := http.DefaultTransport.(*http.Transport)

	if ok {
		x := appConfig.HTTPTransport

		if x.MaxIdleConns > 0 {
			xlog.Info("Http.Transport.MaxIdleConns=%v", x.MaxIdleConns)
			t.MaxIdleConns = x.MaxIdleConns
		}
		if x.IdleConnTimeout > 0 {
			xlog.Info("Http.Transport.IdleConnTimeout=%v", x.IdleConnTimeout)
			t.IdleConnTimeout = time.Duration(x.IdleConnTimeout) * time.Second
		}
		if x.MaxConnsPerHost > 0 {
			xlog.Info("Http.Transport.MaxConnsPerHost=%v", x.MaxConnsPerHost)
			t.MaxConnsPerHost = x.MaxConnsPerHost
		}

		if x.MaxIdleConnsPerHost > 0 {
			xlog.Info("Http.Transport.MaxIdleConnsPerHost=%v", x.MaxIdleConnsPerHost)
			t.MaxIdleConnsPerHost = x.MaxIdleConnsPerHost
		}

	} else {
		xlog.Error("Cannot init http.Transport")
	}
}

func MustNewAppServiceProd() AppService {

	appService := &defaultAppService{}

	appService.mustConfig()
	appService.mustBuild()

	return appService
}
func MustNewAppServiceTesting() AppService {

	env := map[string]string{
		"identity_is_auth_phone":       "true",
		"identity_is_auth_email":       "true",
		"identity_phone_number_prefix": "+123",
	}

	for k, v := range env {
		_ = os.Setenv(strings.ToUpper("app_"+k), v)
	}

	return MustNewAppServiceProd()
}
func (x *defaultAppService) Account() AccountService     { return x.accountService }
func (x *defaultAppService) Auth() AuthService           { return x.authService }
func (x *defaultAppService) AuthAdmin() AuthAdminService { return x.authAdminService }

func (x *defaultAppService) Config() *config.AppConfig { return x.configSource.Config() }

// func (x *appService) Logger() logger.AppLogger       { return x.container.Logger() }

func (x *defaultAppService) UserLang(code string) i18n.UserLang { return x.lang.UserLang(code) }
func (x *defaultAppService) HasLang(code string) bool           { return x.lang.HasLang(code) }
func (x *defaultAppService) Messenger() messenger.AppMessenger  { return x.messenger }

func (x *defaultAppService) Vault() VaultService { return x.vaultService }

func (x *defaultAppService) Repository() repository.AppRepository { return x.repository }
