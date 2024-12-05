package service

import (
	"go-auth-admin/internal/config/consts"
	"go-auth-admin/internal/util/utilaccess"
	"go-auth-admin/internal/util/utilpaging"
	"go-auth-admin/internal/util/utilstring"

	"github.com/google/uuid"
)

var userAccountOmit = []string{
	"password_hash",
	"created_at",
}

func (x *UserAccount) Fill() {
	if x.ID == "" {
		x.ID = uuid.New().String()
		// x.CreatedAt = time.Now().UTC()
	}

	// err := res.RefreshSecurityStamp()
	// if err != nil {
	// 	return nil, err
	// }
}

type UserAccountDAO struct {
	appService AppService
}

type AuthService interface {
	UserAccounts() *UserAccountDAO
}

type defaultAuthService struct {
	appService  AppService
	userAccount UserAccountDAO
}

func newAuthService(appService AppService) AuthService {

	res := &defaultAuthService{

		appService: appService,
		userAccount: UserAccountDAO{
			appService: appService,
		},
	}

	return res
}

func (x *defaultAuthService) UserAccounts() *UserAccountDAO {
	return &x.userAccount
}

func (x *UserAccountDAO) Check(filter *utilpaging.PagingInputDTO) {
	filter.Limit = min(filter.Limit, 10) // validate
}
func (x *UserAccountDAO) Permissions(userAccount *UserAccount, dto *utilaccess.PermissionsDTO) {
	dto.Fill(userAccount.Roles, consts.AuthRolePrefix)
}

func (x *UserAccountDAO) Where(filter *utilpaging.PagingInputDTO) (whereCondition string, whereArgs []any, err error) {
	whereCondition = "1=1"
	whereArgs = []any{}

	if v := filter.Search; v != "" { // filter.GetFilter("text");
		whereCondition += " and (username ilike ? or email ilike ? or phone_number ilike ?)" // " and (title ilike ? or content_md ilike ?)"
		whereArgs = append(whereArgs, "%"+v+"%", "%"+v+"%", "%"+v+"%")
	}

	if whereCondition == "1=1" {
		whereCondition = ""
	}

	return whereCondition, whereArgs, err
}

func (x *UserAccountDAO) Sort(filter *utilpaging.PagingInputDTO) (sqlSort string, err error) {

	sqlSort = "id desc"
	switch filter.Sort {
	case "-id":
		sqlSort = "id desc"
	case "id":
		sqlSort = "id asc"
	default:
		filter.Sort = "-id"
	}

	return sqlSort, err
}

func (x *UserAccountDAO) Query(filter *utilpaging.PagingInputDTO, output *utilpaging.PagingOutputDTO[UserAccount], omitColumns *[]string) (err error) {

	x.Check(filter)

	repo := x.appService.Repository()

	sqlWhere, sqlWhereArgs, _ := x.Where(filter)
	sqlSort, _ := x.Sort(filter)

	var count int64

	err = repo.Model(&UserAccount{}).
		Where(sqlWhere, sqlWhereArgs...).
		Count(&count).Error

	if err != nil {
		return err
	}

	info := filter.Info(int(count))
	output.Fill(filter, info)
	output.Data = make([]*UserAccount, 0, info.Limit)

	if omitColumns == nil {
		omitColumns = &[]string{}
	}

	err = repo.
		Where(sqlWhere, sqlWhereArgs...).
		Order(sqlSort).
		Omit(*omitColumns...). // ContentMD ContentHTML
		Limit(info.Limit).
		Offset(info.Offset).
		Find(&output.Data).Error

	if err != nil {
		return err
	}

	return err
}

func (x *UserAccountDAO) FindByID(id string) (*UserAccount, error) {
	if id == "" {
		return nil, nil // fmt.Errorf("id cannot be empty")
	}

	user := new(UserAccount)

	result := x.appService.Repository().Find(user, "id = ?", id)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	return user, nil
}
func (x *UserAccountDAO) FindByCode(code string) (*UserAccount, error) {
	if code == "" {
		return nil, nil // fmt.Errorf("id cannot be empty")
	}

	user := new(UserAccount)

	result := x.appService.Repository().Find(user, "code = ?", code)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	return user, nil
}
func (x *UserAccountDAO) ID(id string) (string, error) {
	if id == "" {
		return "", nil // fmt.Errorf("id cannot be empty")
	}

	user := new(UserAccount)

	result := x.appService.Repository().Select("id").Limit(1).Find(user, "id = ? ", id)

	if result.Error != nil || result.RowsAffected == 0 {
		return "", result.Error
	}

	return user.ID, nil
}

func (x *UserAccountDAO) Username(username string) (string, error) {
	if username == "" {
		return "", nil // fmt.Errorf("id cannot be empty")
	}

	user := new(UserAccount)

	result := x.appService.Repository().Select("id").Find(user, "username = ?", username)

	if result.Error != nil || result.RowsAffected == 0 {
		return "", result.Error
	}

	return user.ID, nil
}
func (x *UserAccountDAO) Email(email string) (string, error) {
	if email == "" {
		return "", nil // fmt.Errorf("id cannot be empty")
	}

	user := new(UserAccount)

	result := x.appService.Repository().Select("id").Find(user, "email = ?", email)

	if result.Error != nil || result.RowsAffected == 0 {
		return "", result.Error
	}

	return user.ID, nil
}
func (x *UserAccountDAO) PhoneNumber(phoneNumber string) (string, error) {

	phoneNumber = utilstring.NormalizePhoneNumber(phoneNumber) // Normalize

	if phoneNumber == "" {
		return "", nil // fmt.Errorf("id cannot be empty")
	}

	user := new(UserAccount)

	result := x.appService.Repository().Select("id").Find(user, "phone_number = ?", phoneNumber)

	if result.Error != nil || result.RowsAffected == 0 {
		return "", result.Error
	}

	return user.ID, nil
}

func (x *UserAccountDAO) Create(data *UserAccount) error {

	repo := x.appService.Repository()
	data.Fill()
	res := repo.Model(data).Omit(userAccountOmit...).Create(data)
	return res.Error

}
func (x *UserAccountDAO) Update(data *UserAccount) error {
	repo := x.appService.Repository()

	// res := repo.Model(data).Omit(userAccountOmit...).Updates(data) // .Updates() ignores zero-fileds

	// res := repo.Model(data).Omit(userAccountOmit...).Save(data)

	res := repo.Model(data).Select("*" /*over all columns*/).Omit(userAccountOmit...).Updates(data)

	return res.Error
}
func (x *UserAccountDAO) UpdatePassword(id string, pw string) error {

	data := &UserAccount{ID: id}
	if err := data.SetPassword(pw); err != nil {
		return err
	}
	//
	repo := x.appService.Repository()
	res := repo.Model(data).Select("password_hash" /*over all columns*/).Updates(data)
	return res.Error
}
func (x *UserAccountDAO) Delete(id string) error {

	if id == "" {
		return nil
	}

	repo := x.appService.Repository()
	res := repo.Delete(&UserAccount{ID: id})
	return res.Error
}
