/*
@Date: 2022/4/15 15:43
@Author: max.liu
@File : demo_repo
@Desc:
*/

package repo

import (
	"github.com/maxliu9403/common/gormdb"
	"github.com/maxliu9403/go-template/internal/types"
)

type DemoRepo interface {
	GetList(q types.BasicQuery, model, list interface{}) (total int64, err error)
	gormdb.GetByIDCrud
	Deletes([]int64) (err error)
}
