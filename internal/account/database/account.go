package database

import (
	"github.com/real-time-chat/internal/account/model"
	"gorm.io/gorm"
)

type IDbAccount interface {
	Create(*model.Account) (*model.Account, error)
	Save(*model.Account) (*model.Account, error)
	Delete(uint) error

	Login(string, string) (*model.Account, error)
}

type DbAccount struct {
	db *gorm.DB
}

func NewDbAccount(_db *gorm.DB) *DbAccount {
	return &DbAccount{
		db: _db,
	}
}

func (da *DbAccount) Create(acc *model.Account) error {
	tx := da.db.Save(acc)
	return tx.Error
}

func (da *DbAccount) FindByEmail(email string) (*model.Account, error) {
	var res model.Account
	tx := da.db.Where("email = ?", email).First(&res)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &res, nil
}

func (da *DbAccount) Login(email string, pw string) (*model.Account, error) {
	var res model.Account
	tx := da.db.Where("email = ?", email).Where("password = ?", pw).First(&res)
	if tx.Error != nil {
		return nil, tx.Error
	} else {
		return &res, nil
	}
}
