package account

import (
	"fmt"
	"regexp"
	"time"

	"github.com/real-time-chat/internal/account/database"
	"github.com/real-time-chat/internal/account/model"
	"github.com/real-time-chat/pkg/utils"
	"gorm.io/gorm"
)

type IHandler interface {
	Register(*model.AccountRegisterRequest) error
	Login(*model.LoginRequest) (*model.Account, error) // return token
	ValidateAccountField(*model.AccountRegisterRequest) (*model.Account, error)
}

type Handler struct {
	aDB *database.DbAccount
}

var _ IHandler = (*Handler)(nil)

func NewHandler(_db *gorm.DB) *Handler {
	return &Handler{
		aDB: database.NewDbAccount(_db),
	}
}

func (h *Handler) Register(aRR *model.AccountRegisterRequest) error {
	acc, err := h.ValidateAccountField(aRR)
	if err != nil {
		return err
	}

	err = h.aDB.Create(acc)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) Login(lR *model.LoginRequest) (*model.Account, error) {
	// need to check email first
	acc, err := h.aDB.Login(lR.Email, utils.Hash_SHA256(lR.Password))

	if err != nil {
		return nil, fmt.Errorf("Fail to get record: %v", err)
	}
	return acc, nil
}

func (h *Handler) ValidateAccountField(acc *model.AccountRegisterRequest) (*model.Account, error) {
	// check email
	regexPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regexPattern)
	if !re.MatchString(acc.Email) {
		return nil, fmt.Errorf("%s", "Invalid email address")
	}

	// regexPasswordPattern := "^[A-Z][a-zA-Z0-9!@#$%^&*]{5,}[!@#$%^&*]$"
	// re = regexp.MustCompile(regexPasswordPattern)
	// if !re.MatchString(acc.Password) {
	// 	return nil, fmt.Errorf("%s", "Invalid password")
	// }

	regexNamePattern := "^[a-zA-Z0-9]+$"
	re = regexp.MustCompile(regexNamePattern)
	if !re.MatchString(acc.Username) {
		return nil, fmt.Errorf("%s", "Invalid username")
	}

	now := time.Now()

	return &model.Account{
		Email:     acc.Email,
		Username:  acc.Username,
		Password:  utils.Hash_SHA256(acc.Password),
		CreatedAt: &now,
		UpdatedAt: &now,
		Online:    false,
	}, nil
}
