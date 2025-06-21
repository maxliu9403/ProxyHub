package emulator

import (
	"context"
	"errors"
	"fmt"

	"github.com/maxliu9403/ProxyHub/internal/common"
	"github.com/maxliu9403/ProxyHub/internal/logic/group"
	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/ProxyHub/models/factory"
	"github.com/maxliu9403/ProxyHub/models/repo"
	"github.com/maxliu9403/common/gormdb"
	"github.com/maxliu9403/common/logger"
	"gorm.io/gorm"
)

type Svc struct {
	ID  int64
	Ctx context.Context
	DB  *gorm.DB
}

func (s *Svc) getRepo() repo.EmulatorRepo {
	s.DB = gormdb.Cli(s.Ctx)
	return factory.EmulatorRepo(s.DB)
}

func (s *Svc) getProxyRepo() repo.ProxyRepo {
	s.DB = gormdb.Cli(s.Ctx)
	return factory.ProxyRepo(s.DB)
}

func (s *Svc) GetList(q models.GetEmulatorListParams) (data *common.ListData, err error) {
	data = &common.ListData{}

	crud := s.getRepo()
	table := &models.Emulator{}
	demoList := make([]models.Emulator, 0)
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

func (s *Svc) Detail(uuid string) (resp *models.Emulator, err error) {
	crud := s.getRepo()
	pg := &models.Emulator{}
	err = crud.GetByUuid(pg, uuid)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "get %d from db failed: %s", s.ID, err.Error())
		return resp, common.NewErrorCode(common.ErrGetDetail, fmt.Errorf("查询 ID [%d] 详情失败", s.ID))
	}

	return pg, err
}

type DeleteParams struct {
	common.Test
	Uuids []string `json:"Uuids" binding:"required"` // 待删除 ID 列表
}

type ReleaseIPDetail struct {
	IP           string `json:"IP"`           // 释放IP
	ReleaseCount int    `json:"ReleaseCount"` // 释放次数
}

type DeleteResp struct {
	ReleaseIPsDetail []*ReleaseIPDetail `json:"ReleaseIPsDetail"`
}

func (s *Svc) Delete(params DeleteParams) (resp *DeleteResp, err error) {
	proxyRepo := s.getProxyRepo()

	resp = &DeleteResp{ReleaseIPsDetail: []*ReleaseIPDetail{}}

	err = gormdb.Cli(s.Ctx).Transaction(func(tx *gorm.DB) error {
		// 查询要删除的 emulator 的 IP（排除空 IP）
		var ipList []string
		if err := tx.Model(&models.Emulator{}).
			Where("uuid IN ?", params.Uuids).
			Where("ip != ''").
			Pluck("ip", &ipList).Error; err != nil {
			logger.ErrorfWithTrace(s.Ctx, "query emulator IPs failed: %s", err.Error())
			return common.NewErrorCode(common.ErrDeleteEmulator, fmt.Errorf("查询模拟器 IP 失败: %w", err))
		}

		// 释放 IP 计数（去重 + 计数）
		releaseIPMap := make(map[string]int)
		for _, ip := range ipList {
			releaseIPMap[ip] += 1
		}

		// 批量递减 inuse_count
		for ip, count := range releaseIPMap {
			if err := proxyRepo.DecrementInUseTx(tx, ip, count); err != nil {
				logger.ErrorfWithTrace(s.Ctx, "decrement inuse_count for IP [%s] failed: %s", ip, err.Error())
				return common.NewErrorCode(common.ErrDeleteEmulator, fmt.Errorf("更新 proxy 使用数失败 (IP=%s): %w", ip, err))
			}
			// 添加返回详情
			resp.ReleaseIPsDetail = append(resp.ReleaseIPsDetail, &ReleaseIPDetail{
				IP:           ip,
				ReleaseCount: count,
			})
		}

		// 删除模拟器
		if err := tx.Where("uuid IN ?", params.Uuids).Delete(&models.Emulator{}).Error; err != nil {
			logger.ErrorfWithTrace(s.Ctx, "delete emulator failed: %s", err.Error())
			return common.NewErrorCode(common.ErrDeleteEmulator, fmt.Errorf("删除模拟器失败: %w", err))
		}

		return nil
	})

	return resp, err
}

type CreateParams struct {
	common.Test
	BrowserID string `json:"BrowserID" binding:"required,gt=0"` // 窗口ID
	UUID      string `json:"Uuid" binding:"required"`           // 模拟器uuid
	GroupID   int64  `json:"GroupID" binding:"required"`        // 组ID
}

type CreateBatchParams struct {
	Emulators []CreateParams `json:"Emulators" binding:"required"` // 多个模拟器
}

func (p CreateParams) ToModel() *models.Emulator {
	return &models.Emulator{
		UUID:      p.UUID,
		BrowserID: p.BrowserID,
		GroupID:   p.GroupID,
	}
}

type Invalid struct {
	BrowserID string `json:"BrowserID"` // 窗口ID
	UUID      string `json:"Uuid"`      // 模拟器uuid
	GroupID   int64  `json:"GroupID"`   // 组ID
	Message   string
}

type CreateBatchResp struct {
	CreatedCount    int       `json:"CreatedCount"`
	InvalidEmulator []Invalid `json:"InvalidEmulator"` // 校验失败
}

// CreateBatch 只创建新的，如果uuid已经存在则跳过处理
func (s *Svc) CreateBatch(params CreateBatchParams) (*CreateBatchResp, error) {
	invalidEmulator := make([]Invalid, 0)
	toCreate := make([]*models.Emulator, 0)

	// 收集所有 UUID
	uuidMap := make(map[string]CreateParams)
	uuidList := make([]string, 0, len(params.Emulators))
	for _, p := range params.Emulators {

		if p.UUID == "" {
			invalidEmulator = append(invalidEmulator,
				Invalid{BrowserID: p.BrowserID,
					UUID:    p.UUID,
					GroupID: p.GroupID,
					Message: "uuid为空",
				})
			continue
		}
		// 校验group_id是否合法
		groupAPI := group.NewGroupAPI(s.Ctx)
		hasGroup, err := groupAPI.CheckGroupID(p.GroupID)
		if err != nil {
			invalidEmulator = append(invalidEmulator,
				Invalid{BrowserID: p.BrowserID,
					UUID:    p.UUID,
					GroupID: p.GroupID,
					Message: "校验group id合法性失败",
				})
			logger.ErrorfWithTrace(s.Ctx, "check group id failed: %s", err.Error())
			continue
		}
		if !hasGroup {
			invalidEmulator = append(invalidEmulator,
				Invalid{BrowserID: p.BrowserID,
					UUID:    p.UUID,
					GroupID: p.GroupID,
					Message: "group id不是有效值，可能不存在或者未激活",
				})
			continue
		}
		uuidMap[p.UUID] = p
		uuidList = append(uuidList, p.UUID)
	}

	// 查询已存在的 UUID
	existUUIDs, err := s.getRepo().GetExistingUUIDs(uuidList)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "query exist UUIDs failed: %s", err.Error())
		return nil, common.NewErrorCode(common.ErrQueryExistEmulatorUUID, err)
	}
	existSet := make(map[string]struct{}, len(existUUIDs))
	for _, uuid := range existUUIDs {
		existSet[uuid] = struct{}{}
	}

	// 构建待创建模型
	for uuid, p := range uuidMap {
		if _, exists := existSet[uuid]; exists {
			continue // 已存在则跳过
		}
		toCreate = append(toCreate, p.ToModel())
	}

	if len(toCreate) == 0 {
		return &CreateBatchResp{
			CreatedCount:    0,
			InvalidEmulator: invalidEmulator,
		}, nil
	}

	err = s.getRepo().CreateBatch(toCreate)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "batch create failed: %s", err.Error())
		return nil, common.NewErrorCode(common.ErrCreateEmulator, err)
	}

	return &CreateBatchResp{
		CreatedCount:    len(toCreate),
		InvalidEmulator: invalidEmulator,
	}, nil
}

type UpdateParams struct {
	common.Test
	UUID    string  `json:"UUID" binding:"required"`
	IP      *string `json:"IP,omitempty" binding:"omitempty,ip"`
	GroupID *int64  `json:"GroupID,omitempty"  binding:"omitempty,gt=0"`
}

func (s *Svc) Update(params UpdateParams) error {
	var emulator models.Emulator
	err := s.getRepo().GetByUuid(&emulator, params.UUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewErrorCode(common.ErrUpdateEmulator, errors.New("模拟器不存在"))
		}
		logger.ErrorfWithTrace(s.Ctx, "query emulator failed: %s", err.Error())
		return common.NewErrorCode(common.ErrUpdateEmulator, err)
	}

	updateFields := map[string]interface{}{}

	if params.IP != nil {
		updateFields["ip"] = *params.IP
	}

	if params.GroupID != nil {
		groupAPI := group.NewGroupAPI(s.Ctx)
		hasGroup, err := groupAPI.CheckGroupID(*params.GroupID)
		if err != nil {
			return err
		}
		if !hasGroup {
			return errors.New("当前分组ID不是激活状态或者不存在")
		}
		updateFields["group_id"] = *params.GroupID
	}

	err = s.getRepo().Update(params.UUID, updateFields)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "update emulator failed: %s", err.Error())
		return common.NewErrorCode(common.ErrUpdateGroup, err)
	}

	return nil
}
