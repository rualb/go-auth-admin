package authadmin

import (
	"go-auth-admin/internal/config"
	"go-auth-admin/internal/config/consts"
	controller "go-auth-admin/internal/controller"
	"go-auth-admin/internal/mvc"

	"strings"

	"go-auth-admin/internal/i18n"
	"go-auth-admin/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AccountsPasswordDTO struct {
	Input struct {
		ID          string `param:"id"`
		NewPassword string `json:"new_password"`
	}
	Meta struct {
		Status int
	}
	Output struct {
		mvc.ModelBaseDTO
	}
}
type AccountsPasswordAPIController struct {
	appService service.AppService
	appConfig  *config.AppConfig
	userLang   i18n.UserLang
	IsPOST     bool
	webCtxt    echo.Context // webCtxt
	DTO        AccountsPasswordDTO
}

func (x *AccountsPasswordAPIController) Handler() error {
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
func NewAccountsPasswordAPIController(appService service.AppService, c echo.Context) *AccountsPasswordAPIController {

	appConfig := appService.Config()
	return &AccountsPasswordAPIController{
		appService: appService,
		appConfig:  appConfig,
		userLang:   controller.UserLang(c, appService),
		IsPOST:     controller.IsPOST(c),
		webCtxt:    c,
	}
}

func (x *AccountsPasswordAPIController) validateDTOFields() (err error) {

	dto := &x.DTO
	input := &dto.Input
	output := &dto.Output
	meta := &dto.Meta
	srv := x.appService.AuthAdmin()

	if x.IsPOST {

		// validate input: add update

		{
			input.NewPassword = strings.TrimSpace(input.NewPassword)
		}

		{
			v := output.NewModelValidatorStr(x.userLang, "new_password", "New password", input.NewPassword, consts.PasswordMaxLength)
			v.Password(consts.PasswordMinLength)
		}

	}

	if !output.IsModelValid() {
		meta.Status = http.StatusUnprocessableEntity // 422 validation
		return nil
	}

	if x.IsPOST {
		// exists: delete, update

		id, err := srv.UserAccounts().ID(input.ID)
		if err != nil {
			return err
		}

		if id == "" {
			meta.Status = http.StatusNotFound // 404
			return nil
		}

	}

	return nil

}

func (x *AccountsPasswordAPIController) validateDTO() error {

	dto := &x.DTO
	input := &dto.Input

	c := x.webCtxt

	if err := c.Bind(input); err != nil {
		return err
	}

	return x.validateDTOFields()

}

func (x *AccountsPasswordAPIController) handlePOST() (err error) {
	userLang := x.userLang
	dto := &x.DTO
	output := &dto.Output
	input := &dto.Input
	srv := x.appService.AuthAdmin()

	if err = srv.UserAccounts().UpdatePassword(input.ID, input.NewPassword); err != nil {
		return err
	}

	output.Message = userLang.Lang("Password changed")
	output.Status = consts.StatusSuccess
	return nil

}

func (x *AccountsPasswordAPIController) handleDTO() error {

	dto := &x.DTO
	meta := &dto.Meta
	output := &dto.Output

	if meta.Status > 0 {
		return nil // stop processing
	}

	switch {
	case x.IsPOST:
		return x.handlePOST()
	default:
		{
			meta.Status = http.StatusMethodNotAllowed
			output.Message = "Method action undef"
		}
	}

	return nil
}
func (x *AccountsPasswordAPIController) responseDTOAsAPI() (err error) {

	dto := &x.DTO
	meta := &dto.Meta
	output := &dto.Output
	c := x.webCtxt

	if meta.Status == 0 {
		meta.Status = http.StatusOK
	}

	return c.JSON(meta.Status, output)

}

func (x *AccountsPasswordAPIController) responseDTO() (err error) {
	return x.responseDTOAsAPI()
}
