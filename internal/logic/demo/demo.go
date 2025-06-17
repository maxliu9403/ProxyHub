package demo

import (
	"context"
	"fmt"
	"github.com/maxliu9403/ProxyHub/internal/types"

	"github.com/maxliu9403/ProxyHub/internal/common"
	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/ProxyHub/models/factory"
	"github.com/maxliu9403/ProxyHub/models/repo"
	"github.com/maxliu9403/common/gormdb"
	"github.com/maxliu9403/common/logger"
	"gorm.io/gorm"
)

type Svc struct {
	ID          int64
	AccountID   string
	Ctx         context.Context
	RunningTest bool
	DB          *gorm.DB
}

func (s *Svc) getRepo() repo.DemoRepo {
	if s.RunningTest {
		return factory.DemoRepoForTest()
	}

	s.DB = gormdb.Cli(s.Ctx)
	return factory.DemoRepo(s.DB)
}

func (s *Svc) GetList(q types.BasicQuery) (data *common.ListData, err error) {
	data = &common.ListData{}

	crud := s.getRepo()
	table := &models.Demo{}
	demoList := make([]models.Demo, 0)
	total, err := crud.GetList(q, table, &demoList)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "query list failed: %s", err.Error())
		return data, common.NewErrorCode(common.ErrGetList, err)
	}

	data.Counts = total
	data.Data = demoList

	return data, err
}

type IDParams struct {
	common.Test
	ID int64 `json:"Id" binding:"required"` // 主键 ID
}

func (s *Svc) Detail() (resp *models.Demo, err error) {
	crud := s.getRepo()
	demoInfo := &models.Demo{}

	err = crud.GetByID(demoInfo, s.ID)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "get %d from db failed: %s", s.ID, err.Error())
		return resp, common.NewErrorCode(common.ErrGetDetail, fmt.Errorf("查询 ID [%d] 详情失败", s.ID))
	}

	return demoInfo, err
}

type DeleteParams struct {
	common.Test
	IDs []int64 `json:"ids" binding:"required"` // 待删除 ID 列表
}

func (s *Svc) Delete(params DeleteParams) (err error) {
	crud := s.getRepo()
	err = crud.Deletes(params.IDs)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "delete %d failed: %s", params.IDs, err.Error())
		return common.NewErrorCode(common.ErrDelete, err)
	}

	return err
}
