/*
@Date: 2022/4/15 15:43
@Author: max.liu
@File : demo_repo
@Desc:
*/

package repo

import "github.com/maxliu9403/common/gormdb"

type DemoRepo interface {
	gormdb.GetListCrud
	gormdb.GetByIDCrud
	Deletes([]int64) (err error)
}
