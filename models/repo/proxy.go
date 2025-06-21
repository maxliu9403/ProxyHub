/*
@Date: 2022/4/15 15:43
@Author: max.liu
@File : demo_repo
@Desc:
*/

package repo

import (
	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/common/gormdb"
	"gorm.io/gorm"
)

type ProxyRepo interface {
	gormdb.GetByIDCrud
	GetList(q models.GetListParams, model, list interface{}) (total int64, err error)
	Deletes([]int64) (err error)
	Create(group *models.Proxy) error
	Update(ID int64, fields map[string]interface{}) error
	CreateBatch(proxies []*models.Proxy) error
	DeletesByIps(IPs []string) error
	GetByIP(ip string) (*models.Proxy, error)
	IncrementInUseTx(tx *gorm.DB, ip string, count int) error
	DecrementInUseTx(tx *gorm.DB, ip string, count int) error
	GetByIPForUpdate(ip string) (*models.Proxy, error)
	ListByGroupID(groupID int64) ([]*models.ProxyBrief, error)
}
