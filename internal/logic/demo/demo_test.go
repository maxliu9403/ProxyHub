/*
@Date: 2022/4/15 16:20
@Author: max.liu
@File : demo_test
@Desc:
*/

package demo

import (
	"context"
	"testing"

	"github.com/maxliu9403/common/gormdb"
)

var s = Svc{
	ID:          1,
	Ctx:         context.TODO(),
	RunningTest: true,
}

func TestGetList(t *testing.T) {
	list, err := s.GetList(gormdb.BasicQuery{})
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(list)
}

func TestDetail(t *testing.T) {
	resp, err := s.Detail()
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(resp)
}

func TestDelete(t *testing.T) {
	p := DeleteParams{
		IDs: []int64{1, 2, 3},
	}

	err := s.Delete(p)
	if err != nil {
		t.Fatal(err.Error())
	}
}
