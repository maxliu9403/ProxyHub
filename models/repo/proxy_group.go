/*
@Date: 2022/4/15 15:43
@Author: max.liu
@File : demo_repo
@Desc:
*/

package repo

import (
	"github.com/maxliu9403/ProxyHub/internal/types"
	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/common/gormdb"
)

type ProxyGroupsRepo interface {
	gormdb.GetByIDCrud
	GetList(q types.BasicQuery, model, list interface{}) (total int64, err error)
	Deletes([]int64) (err error)
	Create(group *models.ProxyGroups) error
}
