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

type UserAccountDTO struct {
	service.UserAccount
}
type AccountsEntityDTO struct {
	Input struct {
		Code string         `param:"code"`
		ID   string         `param:"id"`
		Data UserAccountDTO `json:"data,omitempty"`
	}
	Meta struct {
		Status int
	}
	Output struct {
		mvc.ModelBaseDTO
		Data UserAccountDTO `json:"data,omitempty"`
	}
}
type AccountsEntityAPIController struct {
	appService service.AppService
	appConfig  *config.AppConfig
	userLang   i18n.UserLang

	IsGET    bool
	IsPOST   bool
	IsPUT    bool
	IsDELETE bool

	webCtxt echo.Context // webCtxt

	DTO AccountsEntityDTO
}

func (x *AccountsEntityAPIController) Handler() error {
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
func NewAccountsEntityAPIController(appService service.AppService, c echo.Context) *AccountsEntityAPIController {

	appConfig := appService.Config()

	return &AccountsEntityAPIController{
		appService: appService,
		appConfig:  appConfig,
		userLang:   controller.UserLang(c, appService),
		IsGET:      controller.IsGET(c),
		IsPOST:     controller.IsPOST(c),
		IsPUT:      controller.IsPUT(c),
		IsDELETE:   controller.IsDELETE(c),
		webCtxt:    c,
	}
}

func (x *AccountsEntityAPIController) validateDTOFields() (err error) {

	dto := &x.DTO
	input := &dto.Input
	output := &dto.Output
	meta := &dto.Meta
	srv := x.appService.AuthAdmin()

	if x.IsPOST || x.IsPUT {

		// validate input: add update

		{
			input.Data.Username = strings.TrimSpace(input.Data.Username)
			input.Data.Roles = strings.TrimSpace(input.Data.Roles)
			// input.Data.ContentHTML = "" // reset
		}

		{
			_ = output.NewModelValidatorStr(x.userLang, "username", "Username", input.Data.Username, consts.DefaultTextLength)
			// v.Required()
		}
		{
			_ = output.NewModelValidatorStr(x.userLang, "roles", "Roles", input.Data.Roles, consts.DefaultTextLength)
			// v.Required()
		}
		{
			_ = output.NewModelValidatorStr(x.userLang, "phone_number", "Phone number", input.Data.PhoneNumber, consts.DefaultTextLength)
			// v.Required()
		}
		{
			_ = output.NewModelValidatorStr(x.userLang, "email", "Email", input.Data.Email, consts.DefaultTextLength)
			// v.Required()
		}

	}

	if !output.IsModelValid() {
		meta.Status = http.StatusUnprocessableEntity // 422 validation
		return nil
	}

	if x.IsDELETE || x.IsPUT {
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

	if x.IsPOST || x.IsPUT {
		// code dupl: add, update
		{
			val := input.Data.Username
			id, err := srv.UserAccounts().Username(val)
			if err != nil {
				return err
			}
			if !(id == "" || input.ID == id) {
				meta.Status = http.StatusConflict                                         // e.g., duplicate data 409
				output.AddError("username", x.userLang.Lang("Duplicate entry {0}.", val)) // Lang
				return nil
			}
		}
		{
			val := input.Data.PhoneNumber
			id, err := srv.UserAccounts().PhoneNumber(val)
			if err != nil {
				return err
			}
			if !(id == "" || input.ID == id) {
				meta.Status = http.StatusConflict                                             // e.g., duplicate data 409
				output.AddError("phone_number", x.userLang.Lang("Duplicate entry {0}.", val)) // Lang
				return nil
			}
		}
		{
			val := input.Data.Email
			id, err := srv.UserAccounts().Email(val)
			if err != nil {
				return err
			}
			if !(id == "" || input.ID == id) {
				meta.Status = http.StatusConflict                                      // e.g., duplicate data 409
				output.AddError("email", x.userLang.Lang("Duplicate entry {0}.", val)) // Lang
				return nil
			}
		}
	}

	if x.IsPOST || x.IsPUT {

		input.Data.SetPhoneNumber(input.Data.PhoneNumber)
		input.Data.SetEmail(input.Data.Email)
		input.Data.SetUsername(input.Data.Username)

	}

	return nil

}

func (x *AccountsEntityAPIController) validateDTO() error {

	dto := &x.DTO
	input := &dto.Input

	c := x.webCtxt

	if err := c.Bind(input); err != nil {
		return err
	}

	return x.validateDTOFields()

}
func (x *AccountsEntityAPIController) handleGET() (err error) {
	dto := &x.DTO
	input := &dto.Input
	meta := &dto.Meta
	output := &dto.Output
	srv := x.appService.AuthAdmin()

	var res *service.UserAccount

	if input.Code != "" { // /code/:code
		res, err = srv.UserAccounts().FindByCode(input.Code)
	} else { // /:id
		res, err = srv.UserAccounts().FindByID(input.ID)
	}

	if err != nil {
		return err
	}

	if res == nil {
		meta.Status = http.StatusNotFound
	} else {
		output.Data.UserAccount = *res // copy
		// output.Data.ContentMD = ""
	}

	return nil
}

func (x *AccountsEntityAPIController) handlePOST() (err error) {
	dto := &x.DTO
	input := &dto.Input
	output := &dto.Output
	srv := x.appService.AuthAdmin()

	output.Data = input.Data
	output.Data.ID = "" // reset ID

	return srv.UserAccounts().Create(&output.Data.UserAccount)

}
func (x *AccountsEntityAPIController) handlePUT() (err error) {
	dto := &x.DTO
	input := &dto.Input
	output := &dto.Output
	srv := x.appService.AuthAdmin()

	output.Data = input.Data
	output.Data.ID = input.Data.ID // reset ID
	return srv.UserAccounts().Update(&output.Data.UserAccount)

}
func (x *AccountsEntityAPIController) handleDELETE() error {

	dto := &x.DTO
	input := &dto.Input
	output := &dto.Output
	srv := x.appService.AuthAdmin()

	output.Data.ID = input.ID // reset ID
	return srv.UserAccounts().Delete(output.Data.ID)

}
func (x *AccountsEntityAPIController) handleDTO() error {

	dto := &x.DTO
	meta := &dto.Meta
	output := &dto.Output

	if meta.Status > 0 {
		return nil // stop processing
	}

	switch {

	case x.IsGET:
		return x.handleGET()
	case x.IsPOST:
		return x.handlePOST()
	case x.IsPUT:
		return x.handlePUT()
	case x.IsDELETE:
		return x.handleDELETE()
	default:
		{
			meta.Status = http.StatusMethodNotAllowed
			output.AddError("", "Method GET only")
		}
	}

	return nil
}
func (x *AccountsEntityAPIController) responseDTOAsAPI() (err error) {

	dto := &x.DTO
	meta := &dto.Meta
	output := &dto.Output
	c := x.webCtxt

	if meta.Status == 0 {
		meta.Status = http.StatusOK
	}

	if x.IsPOST || x.IsPUT {
		output.Data.PasswordHash = "" // stop direct update
	}

	return c.JSON(meta.Status, output)

}

func (x *AccountsEntityAPIController) responseDTO() (err error) {

	return x.responseDTOAsAPI()

}
