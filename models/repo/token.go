package repo

import (
	"github.com/maxliu9403/ProxyHub/internal/types"
	"github.com/maxliu9403/ProxyHub/models"
)

type TokenRepo interface {
	Create(token *models.Token) error
	Deletes(token []string) error
	Get(token string) (*models.Token, error)
	GetList(q types.BasicQuery, model, list interface{}) (total int64, err error)
	IsValid(token string) (bool, error)
}
