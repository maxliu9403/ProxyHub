package subscribe

import (
	"math/rand"
	"time"

	"github.com/maxliu9403/ProxyHub/internal/logic"
	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/ProxyHub/models/factory"
	"github.com/maxliu9403/common/logger"
	"gorm.io/gorm"
)

func (s *Svc) switchProxy(emulator *models.Emulator, selected *models.Proxy) error {
	// 获取全部代理列表
	err, proxies := s.getProxies(emulator.GroupID)
	if err != nil {
		return err
	}

	// 如果当前 IP 等于推荐 IP，同时代理池大于1，则从其他 IP 中随机挑选一个
	if emulator.IP == selected.IP && len(proxies) > 1 {
		altProxies := make([]models.Proxy, 0)
		for _, p := range proxies {
			if p.IP != emulator.IP {
				altProxies = append(altProxies, p)
			}
		}

		// 随机换一个不同的 IP
		if len(altProxies) > 0 {
			rand.Seed(time.Now().UnixNano())
			selected = &altProxies[rand.Intn(len(altProxies))]
			logger.InfofWithTrace(s.Ctx, "当前已绑定IP: %s，随机更换为新IP: %s", emulator.IP, selected.IP)
		} else {
			logger.InfofWithTrace(s.Ctx, "当前已绑定IP: %s，代理池无其他可用IP", emulator.IP)
			return nil
		}
	}
	// 加锁确保并发安全
	//lockKey := fmt.Sprintf("proxy_group_lock_%d", emulator.GroupID)
	//lock, err := TryLock(s.Ctx, lockKey, 5*time.Second)
	//if err != nil {
	//	return fmt.Errorf("获取锁失败: %w", err)
	//}
	//defer lock.Release(s.Ctx)

	return logic.RetryTransaction(s.DB, func(tx *gorm.DB) error {
		proxyRepo := factory.ProxyRepo(tx)
		emulatorRepo := factory.EmulatorRepo(tx)
		// IP 发生变更，才需要更新 inuse_count 和 emulator.IP
		if emulator.IP != selected.IP {
			// 旧 IP -1 释放（若存在）
			if emulator.IP != "" {
				if err := proxyRepo.DecrementInUse(emulator.IP); err != nil {
					return err
				}
			}

			// 新 IP +1 占用
			if err := proxyRepo.IncrementInUse(selected.IP); err != nil {
				return err
			}

			// 更新 emulator 的 IP 字段
			if err := emulatorRepo.Update(emulator.UUID, map[string]interface{}{
				"ip": selected.IP,
			}); err != nil {
				return err
			}
		} else {
			logger.InfofWithTrace(s.Ctx, "IP 未变更，无需更新: %s", emulator.IP)
		}
		return nil
	}, 3)
}
