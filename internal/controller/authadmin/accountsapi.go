package authadmin

import (
	"go-auth-admin/internal/config"
	controller "go-auth-admin/internal/controller"
	"go-auth-admin/internal/util/utilaccess"
	"go-auth-admin/internal/util/utilpaging"

	"go-auth-admin/internal/i18n"
	"go-auth-admin/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

var userAccountOmit = []string{
	"password_hash",
	"created_at",
}

type AccountsDTO struct {
	Input struct {
		utilpaging.PagingInputDTO
	}
	Meta struct {
		Status int
	}
	Output struct {
		Message string `json:"message,omitempty"`
		utilpaging.PagingOutputDTO[service.UserAccount]
		Permissions utilaccess.PermissionsDTO `json:"permissions,omitempty"`
	}
}

type AccountsAPIController struct {
	appService service.AppService
	appConfig  *config.AppConfig
	userLang   i18n.UserLang

	IsGET bool

	webCtxt echo.Context // webCtxt

	userAccount *service.UserAccount
	DTO         AccountsDTO
}

func (x *AccountsAPIController) Handler() error {
	// TODO sign out force

	err := x.validateDTO()
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

// NewAccountController is constructor.
func NewAccountsAPIController(appService service.AppService, c echo.Context) *AccountsAPIController {

	appConfig := appService.Config()

	return &AccountsAPIController{
		appService:  appService,
		appConfig:   appConfig,
		userLang:    controller.UserLang(c, appService),
		IsGET:       controller.IsGET(c),
		userAccount: controller.GetAccount(c),
		webCtxt:     c,
	}
}

func (x *AccountsAPIController) validateDTO() error {

	dto := &x.DTO
	input := &dto.Input

	c := x.webCtxt

	if err := c.Bind(input); err != nil {
		return err
	}

	// input.Filter = c.QueryParams() //

	return nil
}

func (x *AccountsAPIController) handleDTO() error {

	dto := &x.DTO
	input := &dto.Input
	meta := &dto.Meta
	output := &dto.Output
	// userLang := x.userLang
	// c := x.webCtxt
	// isInputValid := output.IsModelValid()

	if x.IsGET {

		bs := x.appService.AuthAdmin()

		{
			bs.UserAccounts().Permissions(x.userAccount, &output.Permissions)
		}

		//

		if err := bs.UserAccounts().Query(&input.PagingInputDTO, &output.PagingOutputDTO, userAccountOmit); err != nil {
			return err
		}

	} else {
		meta.Status = http.StatusMethodNotAllowed
		output.Message = "method action undef"
	}

	return nil
}
func (x *AccountsAPIController) responseDTOAsAPI() (err error) {

	dto := &x.DTO
	meta := &dto.Meta
	output := &dto.Output
	c := x.webCtxt

	if meta.Status == 0 {
		meta.Status = http.StatusOK
	}

	return c.JSON(meta.Status, output)

}

func (x *AccountsAPIController) responseDTO() (err error) {

	return x.responseDTOAsAPI()

}
