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
)

type ProxyRepo interface {
	gormdb.GetByIDCrud
	GetList(q models.GetListParams, model, list interface{}) (total int64, err error)
	Deletes([]int64) (err error)
	Create(group *models.Proxy) error
	Update(id int64, fields map[string]interface{}) error
	CreateBatch(proxies []*models.Proxy) error
	DeletesByIps(IPs []string) error
}
