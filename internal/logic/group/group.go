package group

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

func (s *Svc) getRepo() repo.GroupsRepo {
	s.DB = gormdb.Cli(s.Ctx)
	return factory.GroupsRepo(s.DB)
}

func (s *Svc) GetList(q types.BasicQuery) (data *common.ListData, err error) {
	data = &common.ListData{}

	crud := s.getRepo()
	table := &models.Groups{}
	demoList := make([]models.Groups, 0)
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

func (s *Svc) Detail() (resp *models.Groups, err error) {
	crud := s.getRepo()
	pg := &models.Groups{}
	err = crud.GetByID(pg, s.ID)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "get %d from db failed: %s", s.ID, err.Error())
		return resp, common.NewErrorCode(common.ErrGetDetail, fmt.Errorf("查询 ID [%d] 详情失败", s.ID))
	}

	return pg, err
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
		return common.NewErrorCode(common.ErrDeleteGroup, err)
	}

	return err
}

type CreateParams struct {
	common.Test
	Name        string `json:"Name" binding:"required"`                // 组名
	MaxOnline   int    `json:"MaxOnline" binding:"required,gt=0"`      // 该分组内的IP最大同时在线模拟器数，必须大于0
	Description string `json:"Description"`                            // 描述
	Available   int    `json:"Available" binding:"required,oneof=1 2"` // 是否激活，1:激活 2:不激活，默认是激活状态
}

func (p CreateParams) ToModel() *models.Groups {
	return &models.Groups{
		Name:        p.Name,
		MaxOnline:   p.MaxOnline,
		Description: p.Description,
		Available:   p.Available,
	}
}

func (s *Svc) Create(params CreateParams) (*models.Groups, error) {
	group := params.ToModel()
	err := s.getRepo().Create(group)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "create group failed: %s", err.Error())
		return nil, common.NewErrorCode(common.ErrCreateGroup, err)
	}

	return group, nil
}

type UpdateParams struct {
	common.Test
	ID          int64   `json:"ID" binding:"required"`                             // 分组 ID，必填
	Name        *string `json:"Name,omitempty"`                                    // 组名
	MaxOnline   *int    `json:"MaxOnline,omitempty" binding:"omitempty,gt=0"`      // 该分组内的IP最大同时在线模拟器数，必须大于0
	Description *string `json:"Description,omitempty"`                             // 描述
	Available   *int    `json:"Available,omitempty" binding:"omitempty,oneof=1 2"` // 是否激活，1:激活 2:不激活
}

func (s *Svc) Update(params UpdateParams) error {
	updateFields := map[string]interface{}{}
	if params.Name != nil {
		updateFields["name"] = *params.Name
	}
	if params.Description != nil {
		updateFields["description"] = *params.Description
	}
	if params.Available != nil {
		updateFields["available"] = *params.Available
	}
	if params.MaxOnline != nil {
		updateFields["max_online"] = *params.MaxOnline
	}

	err := s.getRepo().Update(params.ID, updateFields)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "update group failed: %s", err.Error())
		return common.NewErrorCode(common.ErrUpdateGroup, err)
	}

	return nil
}
