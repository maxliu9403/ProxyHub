package proxy_group

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

func (s *Svc) getRepo() repo.ProxyGroupsRepo {
	if s.RunningTest {
		return factory.DemoRepoForTest()
	}

	s.DB = gormdb.Cli(s.Ctx)
	return factory.ProxyGroupsRepo(s.DB)
}

func (s *Svc) GetList(q types.BasicQuery) (data *common.ListData, err error) {
	data = &common.ListData{}

	crud := s.getRepo()
	table := &models.ProxyGroups{}
	demoList := make([]models.ProxyGroups, 0)
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

func (s *Svc) Detail() (resp *models.ProxyGroups, err error) {
	crud := s.getRepo()
	demoInfo := &models.ProxyGroups{}
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

type CreateParams struct {
	common.Test
	Name        string `json:"Name"`        // 组名
	MaxOnline   int    `json:"MaxOnline"`   // 该分组内的IP最大同时在线模拟器数
	Description string `json:"Description"` // 描述
}

func (p CreateParams) ToModel() *models.ProxyGroups {
	return &models.ProxyGroups{
		Name:        p.Name,
		MaxOnline:   p.MaxOnline,
		Description: p.Description,
	}
}

func (s *Svc) Create(params CreateParams) (*models.ProxyGroups, error) {
	group := params.ToModel()
	err := s.getRepo().Create(group)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "create group failed: %s", err.Error())
		return nil, common.NewErrorCode(common.ErrCreate, err)
	}

	return group, nil
}
