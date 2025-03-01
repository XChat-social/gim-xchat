package repo

import (
	"gim/internal/business/domain/user/model"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"time"
)

type tokenDao struct{}

var TokenDao = new(tokenDao)

func (d *tokenDao) Create(token *model.Token) error {
	now := time.Now()
	token.CreatedAt = now
	token.UpdatedAt = now
	err := db.DB.Create(token).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}
func (d *tokenDao) GetByAddress(address string) (*model.Token, error) {
	var token model.Token
	return &token, db.DB.Where("token_address = ?", address).First(&token).Error
}
