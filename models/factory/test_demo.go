/*
@Date: 2022/4/15 15:46
@Author: max.liu
@File : test_repo
@Desc:
*/

package factory

import (
	"github.com/maxliu9403/common/gormdb"
	"github.com/maxliu9403/go-template/models"
	"github.com/maxliu9403/go-template/models/repo"
)

type testCrudImpl struct {
}

func DemoRepoForTest() repo.DemoRepo {
	return &testCrudImpl{}
}

var testDemoData = []models.Demo{{
	Meta: models.Meta{ID: 1, CreateTime: 1010101},
	User: "test1",
}, {
	Meta: models.Meta{ID: 2, CreateTime: 1010101},
	User: "test2",
}, {
	Meta: models.Meta{ID: 3, CreateTime: 1010101},
	User: "test3",
}, {
	Meta: models.Meta{ID: 4, CreateTime: 1010101},
	User: "test4",
}, {
	Meta: models.Meta{ID: 5, CreateTime: 1010101},
	User: "test5",
}}

func (c *testCrudImpl) GetList(q gormdb.BasicQuery, model, list interface{}) (total int64, err error) {
	_, _ = q, model

	total = 5
	a := list.(*[]models.Demo)
	*a = append(*a, testDemoData...)

	return
}

func (c *testCrudImpl) GetByID(model interface{}, id int64) error {
	_ = id
	m := model.(*models.Demo)
	m.ID = testDemoData[0].ID
	m.User = testDemoData[0].User

	return nil
}

func (c *testCrudImpl) Deletes(ids []int64) error {
	_ = ids

	return nil
}
