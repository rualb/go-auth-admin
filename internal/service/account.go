package service

import (
	"go-auth-admin/internal/util/utilaccess"
	"go-auth-admin/internal/util/utilcrypto"
	utilstring "go-auth-admin/internal/util/utilstring"

	"time"

	"github.com/google/uuid"
)

const (
	IssuerConfirmPhoneNumber = "confirm_phone"
	IssuerConfirmEmail       = "confirm_email"
)
const (
	TokenLifetimeSignupWithPhoneNumber = time.Minute * 30 // 30 minutes

	TokenLifetimeSignupWithEmail = time.Minute * 30 // 30 minutes
)
const (
	SecurityStampLenDefault = 16
)

// UserAccount Username,Email,NormalizedEmail are uniqueIndex with condition "not empty"
type UserAccount struct {
	ID              string `json:"id" gorm:"size:255;primaryKey"`
	Username        string `json:"username,omitempty" gorm:"size:255;uniqueIndex:,where:username != ''"`
	PhoneNumber     string `json:"phone_number,omitempty" gorm:"size:255;uniqueIndex:,where:phone_number != ''"`
	Email           string `json:"email,omitempty" gorm:"size:255"`                             // use this on emailing and show
	NormalizedEmail string `json:"-" gorm:"size:255;uniqueIndex:,where:normalized_email != ''"` // use this on search
	// SecurityStamp   string // Key := Base32(Random(32))  HMACSHA1(Key)  Key == VTOQQ2PQKD7A2KTSXU7OFLKUNI7QEZRJ
	PasswordHash string    `json:"-" gorm:"size:255"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"` // auto-updated
	Roles        string    `json:"roles,omitempty" gorm:"size:255"`
}

func (x *UserAccount) HasAnyOfRoles(roles ...string) bool {
	return utilaccess.HasAnyOfRoles(x.Roles, roles...)
}

func (x *UserAccount) SetUsername(value string) {
	valueNorm := utilstring.NormalizeText(value)
	x.Username = valueNorm
}

func (x *UserAccount) SetPhoneNumber(value string) {
	valueNorm := utilstring.NormalizePhoneNumber(value)
	x.PhoneNumber = valueNorm

	// x.Username = valueNorm
}

func (x *UserAccount) SetEmail(value string) {
	valueNorm := utilstring.NormalizeEmail(value)

	x.Email = value
	x.NormalizedEmail = valueNorm

	// x.Username = valueNorm
}

func (x *UserAccount) SetPassword(pw string) error {

	hash, err := utilcrypto.HashPassword(pw) // bcrypt inside

	if err != nil {
		return err
	}

	x.PasswordHash = hash

	// x.RefreshSecurityStamp()

	return nil
}
func (x *UserAccount) CompareHashAndPassword(str string) bool {

	return utilcrypto.CompareHashAndPassword(x.PasswordHash, str)

}

func NewUserAccount() (*UserAccount, error) {

	now := time.Now().UTC() // now

	id := uuid.New().String()
	res := &UserAccount{
		CreatedAt: now,
		ID:        id,
	}

	// err := res.RefreshSecurityStamp()
	// if err != nil {
	// 	return nil, err
	// }

	return res, nil
}

// AccountService is a service for managing user account.
type AccountService interface {
	FindByID(id string) (*UserAccount, error)
}

type defaultAccountService struct {
	appService AppService
}

// NewAccountService is constructor.
func newAccountService(appService AppService) AccountService {

	return &defaultAccountService{
		appService: appService,
	}
}

func (x defaultAccountService) FindByID(id string) (*UserAccount, error) {

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
