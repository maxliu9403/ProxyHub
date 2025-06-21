package subscribe

import (
	"context"
	"fmt"

	"github.com/maxliu9403/ProxyHub/internal/logic/token"

	"github.com/maxliu9403/ProxyHub/internal/types"
	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/ProxyHub/models/factory"
	"github.com/maxliu9403/ProxyHub/models/repo"
	"github.com/maxliu9403/common/gormdb"
	"github.com/maxliu9403/common/logger"
	"gorm.io/gorm"
)

type Svc struct {
	ID             int64
	Ctx            context.Context
	DB             *gorm.DB
	TokenValidator token.Validator // 依赖注入
}

func (s *Svc) getGroupRepo() repo.GroupsRepo {
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

func (s *Svc) getProxies(groupId int64) (err error, proxies []models.Proxy) {
	query := models.GetListParams{
		GroupIDs: []int64{groupId},
		BasicQuery: types.BasicQuery{
			Order: "inuse_count desc",
		},
	}
	proxies = make([]models.Proxy, 0)
	_, err = s.getProxyRepo().GetList(query, &models.Proxy{}, &proxies)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "get proxies: %s", err.Error())
		return
	}
	return
}

func (s *Svc) prepareAndSelectProxy(token, uuid string) (*models.Emulator, *models.Groups, error) {
	// 校验 token
	if isValid, err := s.TokenValidator.ValidateToken(token); err != nil || !isValid {
		return nil, nil, fmt.Errorf("token 无效或检查失败: %w", err)
	}

	// 获取模拟器
	emulator := &models.Emulator{}
	if err := s.getEmulatorRepo().GetByUuid(emulator, uuid); err != nil {
		return nil, nil, fmt.Errorf("模拟器获取失败: %w", err)
	}

	// 获取分组
	group := &models.Groups{}
	if err := s.getGroupRepo().GetByID(group, emulator.GroupID); err != nil {
		return nil, nil, fmt.Errorf("分组获取失败: %w", err)
	}

	return emulator, group, nil
}

func (s *Svc) Subscribe(token string, uuid string) (clashCfg string, err error) {
	// Step 1: 校验并准备数据
	emulator, group, err := s.prepareAndSelectProxy(token, uuid)
	if err != nil {
		return
	}

	// Step 2: 代理切换（幂等 + 原子 + 事务 + 锁）
	selectedProxy, err := s.switchProxy(emulator, group)
	if err != nil {
		return
	}
	// Step 3: 渲染 Clash 配置
	clashCfg, err = s.renderClashConfig(selectedProxy)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "渲染Clash配置失败: %s", err.Error())
		return
	}
	return
}
