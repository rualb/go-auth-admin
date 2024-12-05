package service

type AuthAdminService interface {
	UserAccounts() *UserAccountDAO
}

type defaultAuthAdminService struct {
	appService AppService
	account    UserAccountDAO
}

func newAuthAdminService(appService AppService) AuthAdminService {

	res := &defaultAuthAdminService{

		appService: appService,
		account: UserAccountDAO{
			appService: appService,
		},
	}

	return res
}

func (x *defaultAuthAdminService) UserAccounts() *UserAccountDAO {
	return &x.account
}
