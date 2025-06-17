/*
@Date: 2022/4/15 15:46
@Author: max.liu
@File : test_repo
@Desc:
*/

package factory

import (
	"github.com/maxliu9403/ProxyHub/internal/types"
	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/ProxyHub/models/repo"
)

type testCrudImpl struct {
}

func (t testCrudImpl) GetByID(model interface{}, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (t testCrudImpl) GetList(q types.BasicQuery, model, list interface{}) (total int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (t testCrudImpl) Deletes(int64s []int64) (err error) {
	//TODO implement me
	panic("implement me")
}

func (t testCrudImpl) Create(group *models.ProxyGroups) error {
	//TODO implement me
	panic("implement me")
}

func DemoRepoForTest() repo.ProxyGroupsRepo {
	return &testCrudImpl{}
}

var testDemoData = []models.ProxyGroups{{
	Meta: models.Meta{ID: 1, CreateTime: 1010101},
	Name: "test1",
}, {
	Meta: models.Meta{ID: 2, CreateTime: 1010101},
	Name: "test2",
}}
