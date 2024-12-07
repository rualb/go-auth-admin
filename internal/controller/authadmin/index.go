package authadmin

import (
	"fmt"
	"go-auth-admin/internal/config"
	controller "go-auth-admin/internal/controller"
	"go-auth-admin/internal/service"
	"time"

	"go-auth-admin/internal/i18n"
	"go-auth-admin/internal/mvc"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthAdminIndexController struct {
	appService service.AppService
	appConfig  *config.AppConfig
	userLang   i18n.UserLang

	IsGET  bool
	IsPOST bool

	webCtxt echo.Context // webCtxt

	DTO struct {
		Input struct {
		}
		Meta struct {
			IsFragment bool `json:"-"`
		}
		Output struct {
			mvc.ModelBaseDTO
			LangCode  string
			AppConfig struct {
				AppTitle string `json:"app_title,omitempty"`
				TmTitle  string `json:"tm_title,omitempty"`
			}
			Title     string
			LangWords map[string]string
		}
	}
}

func (x *AuthAdminIndexController) Handler() error {

	err := x.createDTO()
	if err != nil {
		return err
	}

	err = x.handleDTO()
	if err != nil {
		return err
	}

	err = x.responseDTO()
	if err != nil {
		return err
	}

	return nil
}

func NewAuthAdminIndexController(appService service.AppService, c echo.Context) *AuthAdminIndexController {

	return &AuthAdminIndexController{
		appService: appService,
		appConfig:  appService.Config(),
		userLang:   controller.UserLang(c, appService),
		IsGET:      controller.IsGET(c),
		IsPOST:     controller.IsPOST(c),
		webCtxt:    c,
	}
}

func (x *AuthAdminIndexController) validateFields() {

}

func (x *AuthAdminIndexController) createDTO() error {

	dto := &x.DTO
	c := x.webCtxt

	if err := c.Bind(dto); err != nil {
		return err
	}

	x.validateFields()

	return nil
}

func (x *AuthAdminIndexController) handleDTO() error {

	dto := &x.DTO
	// input := &dto.Input
	output := &dto.Output
	// meta := &dto.Meta
	// c := x.webCtxt

	userLang := x.userLang
	output.LangCode = userLang.LangCode()
	output.Title = userLang.Lang("Auth admin") // TODO /*Lang*/
	output.LangWords = userLang.LangWords()

	cfg := &output.AppConfig

	cfg.AppTitle = x.appConfig.Title
	cfg.TmTitle = fmt.Sprintf("© %v %s", time.Now().Year(), x.appConfig.Title)

	return nil
}

func (x *AuthAdminIndexController) responseDTOAsMvc() (err error) {

	dto := &x.DTO
	meta := &dto.Meta
	output := &dto.Output
	appConfig := x.appConfig
	lang := x.userLang
	c := x.webCtxt

	data, err := mvc.NewModelWrap(c, output, meta.IsFragment, "Auth admin" /*Lang*/, appConfig, lang)
	if err != nil {
		return err
	}

	err = c.Render(http.StatusOK, "index.html", data)

	if err != nil {
		return err
	}

	return nil
}

func (x *AuthAdminIndexController) responseDTO() (err error) {

	return x.responseDTOAsMvc()

}
