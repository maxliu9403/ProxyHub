/*
@Date: 2022/4/15 16:33
@Author: max.liu
@File : repo
@Desc:
*/

package factory

import (
	"github.com/maxliu9403/common/gormdb"
	"github.com/maxliu9403/go-template/models"
	"github.com/maxliu9403/go-template/models/repo"
	"gorm.io/gorm"
)

type demoCrudImpl struct {
	Conn *gorm.DB
}

func DemoRepo(db *gorm.DB) repo.DemoRepo {
	return &demoCrudImpl{Conn: db}
}

func (r *demoCrudImpl) GetList(q gormdb.BasicQuery, model, list interface{}) (total int64, err error) {
	crud := gormdb.NewCRUD(r.Conn)
	total, err = crud.GetList(q, model, list)
	return
}

func (r *demoCrudImpl) GetByID(model interface{}, id int64) error {
	crud := gormdb.NewCRUD(r.Conn)
	err := crud.GetByID(model, id)
	return err
}

func (r *demoCrudImpl) Deletes(ids []int64) (err error) {
	err = r.Conn.Delete(&models.Demo{}, ids).Error
	return err
}
