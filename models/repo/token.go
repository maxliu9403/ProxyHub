package repo

import (
	"time"

	"github.com/maxliu9403/ProxyHub/models"
)

type TokenRepo interface {
	Create(token *models.Token) error
	Deletes(token []string) error
	Get(token string) (*models.Token, error)
	GetList(q models.GetTokenListParams, model, list interface{}) (total int64, err error)
	IsValid(token string) (bool, error)
	GetValidTokensByGroup(groupID int64, now time.Time) ([]models.Token, error)
	GetByGroupID(groupID int64) (*models.Token, error)
}
