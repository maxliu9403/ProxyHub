package subscribe

import (
	"context"
	"errors"
	"fmt"

	"github.com/maxliu9403/ProxyHub/internal/types"
	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/ProxyHub/models/factory"
	"github.com/maxliu9403/ProxyHub/models/repo"
	"github.com/maxliu9403/common/gormdb"
	"github.com/maxliu9403/common/logger"
	"gorm.io/gorm"
)

type Svc struct {
	ID          int64
	Ctx         context.Context
	RunningTest bool
	DB          *gorm.DB
}

func (s *Svc) getGroupRepo() repo.GroupsRepo {
	s.DB = gormdb.Cli(s.Ctx)
	return factory.GroupsRepo(s.DB)
}

func (s *Svc) getProxyRepo() repo.ProxyRepo {
	s.DB = gormdb.Cli(s.Ctx)
	return factory.ProxyRepo(s.DB)
}

func (s *Svc) getTokenRepo() repo.TokenRepo {
	s.DB = gormdb.Cli(s.Ctx)
	return factory.TokenRepo(s.DB)
}
func (s *Svc) getEmulatorRepo() repo.EmulatorRepo {
	s.DB = gormdb.Cli(s.Ctx)
	return factory.EmulatorRepo(s.DB)
}

func (s *Svc) check(token string, groupId int64) error {
	// 校验Token是否有效
	tokenRepo := s.getTokenRepo()
	isValid, err := tokenRepo.IsValid(token)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "校验Token失败: %s", err.Error())
		return err
	}
	if !isValid {
		err = errors.New("token已经过期")
		logger.ErrorfWithTrace(s.Ctx, err.Error())
		return err
	}

	// 校验group是否有效
	groupRepo := s.getGroupRepo()
	active, err := groupRepo.IsGroupActive(groupId)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "校验GroupId失败: %s", err.Error())
		return err
	}
	if !active {
		err = errors.New("GroupId 不是激活的或者不存在")
		logger.ErrorfWithTrace(s.Ctx, err.Error())
		return err
	}
	return nil
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

func (s *Svc) getEmulator(uuid string) (err error, emulator *models.Emulator) {
	err = s.getEmulatorRepo().GetByUuid(emulator, uuid)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, err.Error())
		return
	}
	return
}

func (s *Svc) getGroup(groupId int64) (err error, group *models.Groups) {
	err = s.getGroupRepo().GetByID(group, groupId)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, err.Error())
		return
	}
	return
}

func (s *Svc) prepareAndSelectProxy(token, uuid string) (*models.Emulator, *models.Proxy, error) {
	// 校验 token
	if isValid, err := s.getTokenRepo().IsValid(token); err != nil || !isValid {
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

	// 查找 Proxy，找到负载最小的那个IP，但是可能就是当前模拟器正在使用的IP
	// 如果是这种逻辑，那么在后面switchProxy时会实现强行切换（如有有得选的前提下）
	selectedProxy, err := s.getProxyRepo().GetOneForUpdate(emulator.GroupID, group.MaxOnline)
	if err != nil {
		return nil, nil, fmt.Errorf("查找代理失败: %w", err)
	}

	return emulator, selectedProxy, nil
}

func (s *Svc) Subscribe(token string, uuid string) (clashCfg string, err error) {
	// Step 1: 校验并准备数据
	emulator, selectedProxy, err := s.prepareAndSelectProxy(token, uuid)
	if err != nil {
		return
	}

	// Step 2: 代理切换（幂等 + 原子 + 事务 + 锁）
	if err = s.switchProxy(emulator, selectedProxy); err != nil {
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
