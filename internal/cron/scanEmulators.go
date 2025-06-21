package cron

import (
	"context"
	"fmt"
	"github.com/maxliu9403/ProxyHub/internal/config"
	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/ProxyHub/models/factory"
	"github.com/maxliu9403/ProxyHub/models/repo"
	"github.com/maxliu9403/common/gormdb"
	"github.com/maxliu9403/common/logger"
	"gorm.io/gorm"
	"time"
)

type ScanEmulatorsTaskSvc struct {
	ctx          context.Context
	db           *gorm.DB
	emulatorRepo repo.EmulatorRepo
	proxyRepo    repo.ProxyRepo
	groupRepo    repo.GroupsRepo
}

type ReleaseIPDetail struct {
	IP    string `json:"IP"`
	Count int    `json:"Count"`
}

type UnbindEmulator struct {
	BrowserID string `json:"BrowserID"`
	UUID      string `json:"UUID"`
}

func (s *ScanEmulatorsTaskSvc) getGroupRepo() repo.GroupsRepo {
	s.db = gormdb.Cli(s.ctx)
	return factory.GroupsRepo(s.db)
}

func (s *ScanEmulatorsTaskSvc) getProxyRepo() repo.ProxyRepo {
	s.db = gormdb.Cli(s.ctx)
	return factory.ProxyRepo(s.db)
}

func (s *ScanEmulatorsTaskSvc) getEmulatorRepo() repo.EmulatorRepo {
	s.db = gormdb.Cli(s.ctx)
	return factory.EmulatorRepo(s.db)
}

type GroupReleaseResult struct {
	GroupName       string            `json:"GroupName"`
	MaxOnline       int               `json:"MaxOnline"`
	ReleaseIPDetail []ReleaseIPDetail `json:"ReleaseIPDetail"`
	UnbindEmulator  []UnbindEmulator  `json:"UnbindEmulator"`
}

func (s *ScanEmulatorsTaskSvc) ScanAndDeleteExpiredEmulators() ([]*GroupReleaseResult, error) {
	msgPrefix := "定时清扫模拟器失败"

	// 1. 查询过期 emulator
	emulatorRepo := s.getEmulatorRepo()
	// 多少时间之前的 Emulator 是过期的
	expiredBefore := time.Now().Add(-time.Duration(config.G.CustomCfg.IntervalTime) * time.Hour)
	emulators, err := emulatorRepo.ListExpired(expiredBefore)
	if err != nil {
		logger.ErrorfWithTrace(s.ctx, "%s：查询过期模拟器失败: %s", msgPrefix, err.Error())
		return nil, err
	}
	if len(emulators) == 0 {
		return nil, nil
	}

	// 2. 分组统计
	grouped := make(map[int64][]*models.Emulator)
	ipReleaseMap := make(map[string]int)
	groupIDSet := make(map[int64]struct{})
	for _, e := range emulators {
		grouped[e.GroupID] = append(grouped[e.GroupID], e)
		groupIDSet[e.GroupID] = struct{}{}
		if e.IP != "" {
			ipReleaseMap[e.IP]++
		}
	}

	// 3. 获取 group 元信息
	var groupIDList []int64
	for id := range groupIDSet {
		groupIDList = append(groupIDList, id)
	}
	groupRepo := s.getGroupRepo()
	groupsMap, err := groupRepo.GetByIDs(groupIDList)
	if err != nil {
		logger.ErrorfWithTrace(s.ctx, "%s：查询 Group 失败: %s", msgPrefix, err.Error())
		return nil, err
	}

	var results []*GroupReleaseResult

	// 4. 启动事务处理：更新 Proxy 使用数 + 删除 Emulator
	err = s.db.Transaction(func(tx *gorm.DB) error {
		proxyRepo := factory.ProxyRepo(tx)
		emulatorRepo := factory.EmulatorRepo(tx)

		// 4.1 遍历 IP 执行递减
		for ip, count := range ipReleaseMap {
			if err := proxyRepo.DecrementInUseTx(tx, ip, count); err != nil {
				logger.ErrorfWithTrace(s.ctx, "更新 IP [%s] 的使用数失败: %s", ip, err.Error())
				return fmt.Errorf("更新 IP %s 使用数失败: %w", ip, err)
			}
		}

		// 4.2 删除 emulator
		var uuids []string
		for _, e := range emulators {
			uuids = append(uuids, e.UUID)
		}
		if err := emulatorRepo.DeletesByUuidsTx(tx, uuids); err != nil {
			logger.ErrorfWithTrace(s.ctx, "删除 Emulator 失败: %s", err.Error())
			return err
		}

		// 5. 构造返回结果
		for groupID, emus := range grouped {
			groupInfo := groupsMap[groupID]
			result := &GroupReleaseResult{
				GroupName:       groupInfo.Name,
				MaxOnline:       groupInfo.MaxOnline,
				UnbindEmulator:  []UnbindEmulator{},
				ReleaseIPDetail: []ReleaseIPDetail{},
			}
			ipCount := make(map[string]int)
			for _, e := range emus {
				result.UnbindEmulator = append(result.UnbindEmulator, UnbindEmulator{
					BrowserID: e.BrowserID,
					UUID:      e.UUID,
				})
				if e.IP != "" {
					ipCount[e.IP]++
				}
			}
			for ip, count := range ipCount {
				result.ReleaseIPDetail = append(result.ReleaseIPDetail, ReleaseIPDetail{
					IP:    ip,
					Count: count,
				})
			}
			results = append(results, result)
		}

		return nil
	})

	if err != nil {
		logger.ErrorfWithTrace(s.ctx, "%s：事务执行失败: %s", msgPrefix, err.Error())
	}

	return results, err
}
