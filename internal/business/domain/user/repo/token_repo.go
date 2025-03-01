package repo

import (
	"gim/internal/business/domain/user/model"
)

type tokenRepo struct{}

var TokenRepo = new(tokenRepo)

func (r *tokenRepo) Create(token *model.Token) error {
	return TokenDao.Create(token)
}

func (r *tokenRepo) GetByAddress(address string) (*model.Token, error) {
	return TokenDao.GetByAddress(address)
}
