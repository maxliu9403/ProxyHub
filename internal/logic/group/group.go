package group

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/maxliu9403/ProxyHub/internal/common"
	"github.com/maxliu9403/ProxyHub/internal/logic"
	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/ProxyHub/models/factory"
	"github.com/maxliu9403/ProxyHub/models/repo"
	"github.com/maxliu9403/common/gormdb"
	"github.com/maxliu9403/common/logger"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Svc struct {
	ID  int64
	Ctx context.Context
	DB  *gorm.DB
}

func (s *Svc) getRepo() repo.GroupsRepo {
	s.DB = gormdb.Cli(s.Ctx)
	return factory.GroupsRepo(s.DB)
}

func (s *Svc) getProxyRepo() repo.ProxyRepo {
	s.DB = gormdb.Cli(s.Ctx)
	return factory.ProxyRepo(s.DB)
}

func (s *Svc) getEmulatorRepo() repo.EmulatorRepo {
	s.DB = gormdb.Cli(s.Ctx)
	return factory.EmulatorRepo(s.DB)
}

func (s *Svc) getTokenRepo() repo.TokenRepo {
	s.DB = gormdb.Cli(s.Ctx)
	return factory.TokenRepo(s.DB)
}

func (s *Svc) CheckGroupID(groupID int64) (hasGroup bool, err error) {
	groupRepo := s.getRepo()

	// 校验是否存在激活分组
	hasGroup, err = groupRepo.ExistsGroup(groupID)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "check active group failed: %s", err.Error())
		return hasGroup, errors.New("校验分组ID有效性失败")
	}
	return
}

func (s *Svc) GetList(q models.GetGroupListParams) (data *common.ListData, err error) {
	data = &common.ListData{}

	crud := s.getRepo()
	table := &models.Groups{}
	groups := make([]models.Groups, 0)
	total, err := crud.GetList(q, table, &groups)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "query list failed: %s", err.Error())
		return data, common.NewErrorCode(common.ErrGetList, err)
	}

	data.Counts = total
	data.Data = groups

	return data, err
}

type IDParams struct {
	ID int64 `json:"Id" binding:"required"` // 主键 ID
}

type GroupsDetailResp struct {
	models.Groups                         // 包含原始 Group 字段
	Token         string                  `json:"Token"`
	Proxies       []*models.ProxyBrief    `json:"Proxies"`
	Emulators     []*models.EmulatorBrief `json:"Emulators"`
}

func (s *Svc) buildSubscribeLink(token string, emulators []*models.EmulatorBrief) {
	for _, emulator := range emulators {
		emulator.SubscribeLink = fmt.Sprintf("/api/subscribe/%s/%s", token, emulator.UUID)
	}
	return
}

func (s *Svc) Detail() (resp *GroupsDetailResp, err error) {
	groupRepo := s.getRepo()
	group := &models.Groups{}
	if err = groupRepo.GetByID(group, s.ID); err != nil {
		return nil, common.NewErrorCode(common.ErrGetDetail, fmt.Errorf("查询 ID [%d] 详情失败: %w", s.ID, err))
	}

	var (
		proxies    []*models.ProxyBrief
		emulators  []*models.EmulatorBrief
		tokenModel *models.Token
	)

	// 设置 context 超时时间，例如 3 秒
	ctx, cancel := context.WithTimeout(s.Ctx, 5*time.Second)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var e error
		proxies, e = s.getProxyRepo().ListByGroupID(s.ID)
		if e != nil {
			return fmt.Errorf("查询代理失败: %w", e)
		}
		return nil
	})

	eg.Go(func() error {
		var e error
		emulators, e = s.getEmulatorRepo().ListBriefByGroupID(s.ID)
		if e != nil {
			return fmt.Errorf("查询模拟器失败: %w", e)
		}
		return nil
	})

	eg.Go(func() error {
		var e error
		tokenModel, e = s.getTokenRepo().GetByGroupID(s.ID)
		if e != nil {
			return fmt.Errorf("查询 Token 失败: %w", e)
		}
		return nil
	})

	// 等待所有 goroutine 完成或超时
	if err := eg.Wait(); err != nil {
		logger.ErrorfWithTrace(s.Ctx, "group detail query failed: %s", err.Error())
		return nil, common.NewErrorCode(common.ErrGetDetail, err)
	}

	// 构造订阅链接
	s.buildSubscribeLink(tokenModel.Token, emulators)

	// 构建响应
	resp = &GroupsDetailResp{
		Groups:    *group,
		Token:     tokenModel.Token,
		Proxies:   proxies,
		Emulators: emulators,
	}
	return resp, nil
}

type DeleteParams struct {
	common.Test
	IDs []int64 `json:"GroupIds" binding:"required"` // 待删除 ID 列表
}

func (s *Svc) Delete(params DeleteParams) error {
	return gormdb.Cli(s.Ctx).Transaction(func(tx *gorm.DB) error {
		groupIDs := params.IDs

		// 删除 proxies
		if err := tx.Where("group_id IN ?", groupIDs).Delete(&models.Proxy{}).Error; err != nil {
			logger.ErrorfWithTrace(s.Ctx, "delete proxies failed: %s", err.Error())
			return common.NewErrorCode(common.ErrDeleteGroup, fmt.Errorf("删除代理失败: %w", err))
		}

		// 删除 emulators
		if err := tx.Where("group_id IN ?", groupIDs).Delete(&models.Emulator{}).Error; err != nil {
			logger.ErrorfWithTrace(s.Ctx, "delete emulators failed: %s", err.Error())
			return common.NewErrorCode(common.ErrDeleteGroup, fmt.Errorf("删除模拟器失败: %w", err))
		}

		// 删除 tokens
		if err := tx.Where("group_id IN ?", groupIDs).Delete(&models.Token{}).Error; err != nil {
			logger.ErrorfWithTrace(s.Ctx, "delete tokens failed: %s", err.Error())
			return common.NewErrorCode(common.ErrDeleteGroup, fmt.Errorf("删除Token失败: %w", err))
		}

		// 删除 group 本身
		if err := tx.Where("id IN ?", groupIDs).Delete(&models.Groups{}).Error; err != nil {
			logger.ErrorfWithTrace(s.Ctx, "delete groups failed: %s", err.Error())
			return common.NewErrorCode(common.ErrDeleteGroup, fmt.Errorf("删除分组失败: %w", err))
		}

		return nil
	})
}

type CreateParams struct {
	Name        string `json:"Name" binding:"required"`           // 组名，必须唯一
	MaxOnline   int    `json:"MaxOnline" binding:"required,gt=0"` // 该分组内的IP最大同时在线模拟器数，必须大于0
	Description string `json:"Description"`                       // 描述
}

type CreateGroupBatchParams struct {
	Groups []CreateParams `json:"Groups" binding:"required"` // 分组列表，不能为空
}

func (p CreateParams) ToModel() *models.Groups {
	return &models.Groups{
		Name:        p.Name,
		MaxOnline:   p.MaxOnline,
		Description: p.Description,
	}
}

type CreateBatchResp struct {
	CreatedCount int                `json:"CreatedCount"` // 创建成功数
	CreatedList  []CreatedGroupInfo `json:"CreatedList"`  // 创建成功的详情
	InvalidGroup []CreateParams     `json:"InvalidGroup"` // 无效Group
}

type CreatedGroupInfo struct {
	GroupID   int64  `json:"GroupID"`   // 创建成功的组ID
	GroupName string `json:"GroupName"` // 组名
	Token     string `json:"Token"`     // token
}

func (s *Svc) CreateBatch(params CreateGroupBatchParams) (*CreateBatchResp, error) {
	invalidGroup := make([]CreateParams, 0)
	toCreate := make([]*models.Groups, 0)
	createdList := make([]CreatedGroupInfo, 0)
	paramMap := make(map[string]CreateParams)

	// 预处理：过滤非法分组（name为空）
	for _, p := range params.Groups {
		if p.Name == "" {
			invalidGroup = append(invalidGroup, p)
			continue
		}
		model := p.ToModel()
		toCreate = append(toCreate, model)
		paramMap[p.Name] = p
	}

	if len(toCreate) == 0 {
		return &CreateBatchResp{
			CreatedCount: 0,
			CreatedList:  createdList,
			InvalidGroup: invalidGroup,
		}, nil
	}

	// 启动事务
	tx := gormdb.Cli(s.Ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 1. 提取所有分组名
	groupNames := make([]string, 0, len(toCreate))
	for _, g := range toCreate {
		groupNames = append(groupNames, g.Name)
	}
	// 2. 删除所有 deleted_at IS NOT NULL 的冲突记录（避免唯一索引报错）
	if err := tx.
		Unscoped().
		Where("name IN ?", groupNames).
		Where("delete_time IS NOT NULL").
		Delete(&models.Groups{}).Error; err != nil {
		tx.Rollback()
		return nil, common.NewErrorCode(common.ErrCreateGroup, fmt.Errorf("清理软删除记录失败: %w", err))
	}

	// 批量插入分组（不使用 IGNORE，出现重复时整批失败）
	if err := tx.Create(&toCreate).Error; err != nil {
		tx.Rollback()
		logger.ErrorfWithTrace(s.Ctx, "batch create group failed: %s", err.Error())

		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, common.NewErrorCode(common.ErrCreateGroup, fmt.Errorf("存在重复分组名，创建失败"))
		}
		return nil, common.NewErrorCode(common.ErrCreateGroup, err)
	}

	tokenRepo := factory.TokenRepo(tx)
	// 遍历插入成功的分组，为每个分组生成一个 Token
	for _, group := range toCreate {
		tokenStr, err := logic.GenerateSecureToken(32)
		if err != nil {
			tx.Rollback()
			logger.ErrorfWithTrace(s.Ctx, "generate token failed: %s", err.Error())
			return nil, common.NewErrorCode(common.ErrBuildToken, err)
		}

		model := &models.Token{
			Token:       tokenStr,
			Description: fmt.Sprintf("BatchGenerated: %s", group.Name),
			GroupID:     group.ID,
		}
		if err := tokenRepo.Create(model); err != nil {
			tx.Rollback()
			logger.ErrorfWithTrace(s.Ctx, "create token failed: %s", err.Error())
			return nil, common.NewErrorCode(common.ErrCreateGroup, err)
		}

		createdList = append(createdList, CreatedGroupInfo{
			GroupName: group.Name,
			Token:     tokenStr,
			GroupID:   group.ID,
		})
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &CreateBatchResp{
		CreatedCount: len(createdList),
		CreatedList:  createdList,
		InvalidGroup: invalidGroup,
	}, nil
}

type UpdateParams struct {
	common.Test
	ID          int64   `json:"ID" binding:"required"`                        // 分组 ID，必填
	Name        *string `json:"Name,omitempty"`                               // 组名
	MaxOnline   *int    `json:"MaxOnline,omitempty" binding:"omitempty,gt=0"` // 该分组内的IP最大同时在线模拟器数，必须大于0
	Description *string `json:"Description,omitempty"`                        // 描述
}

func (s *Svc) Update(params UpdateParams) error {
	updateFields := map[string]interface{}{}
	if params.Name != nil {
		updateFields["name"] = *params.Name
	}
	if params.Description != nil {
		updateFields["description"] = *params.Description
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
