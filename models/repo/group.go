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

type GroupsRepo interface {
	gormdb.GetByIDCrud
	GetList(q models.GetGroupListParams, model, list interface{}) (total int64, err error)
	Deletes([]int64) (err error)
	Update(id int64, fields map[string]interface{}) error
	ExistsGroup(groupId int64) (bool, error)
	CreateBatch(groups []*models.Groups) error
	GetByIDs(ids []int64) (map[int64]*models.Groups, error)
}
