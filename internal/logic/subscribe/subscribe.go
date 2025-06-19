package subscribe

import (
	"context"
	"errors"

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
	active, err := groupRepo.IsGroupActive(int64(groupId))
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

func (s *Svc) Subscribe(token string, groupId int64) error {
	err := s.check(token, groupId)
	if err != nil {
		return err
	}

	// 第一步：查询groupId下关联的Proxy，同时Proxy的状态是可用的，同时按照InUseCount递增排序返回
	// 第二步：选择其中InUseCount值最小的Proxy
	// 第三步：数据库加锁操作，

	// 第一步：在redis中查找出当前模拟器的正在使用的ProxyIP
	// 第二步：查找出合适的ProxyIP
	//		1：查询groupId下关联的Proxy，同时Proxy的状态是可用的，同时按照InUseCount递增排序返回
	// 		2. 查找到不等于当前使用的ProxyIP，而且InUseCount数最小的那个IP
	// 第三步：渲染base_proxy.yaml，保存返回结果[]byte(config)
	// 第四步：
	return nil
}
