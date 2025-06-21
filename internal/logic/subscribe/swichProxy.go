package subscribe

import (
	"fmt"

	"github.com/maxliu9403/ProxyHub/internal/logic"
	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/ProxyHub/models/factory"
	"github.com/maxliu9403/common/logger"
	"gorm.io/gorm"
)

func (s *Svc) bindEmulatorToProxyIP(tx *gorm.DB, emulator *models.Emulator, selected *models.Proxy) error {
	proxyRepo := factory.ProxyRepo(tx)
	emulatorRepo := factory.EmulatorRepo(tx)
	if emulator.IP == selected.IP {
		logger.InfofWithTrace(s.Ctx, "模拟器 %s 绑定IP未变更: %s", emulator.UUID, emulator.IP)
		return nil
	}

	// 解绑旧 IP
	if emulator.IP != "" {
		if err := proxyRepo.DecrementInUseTx(tx, emulator.IP, 1); err != nil {
			return fmt.Errorf("旧IP %s 减少使用数失败: %w", emulator.IP, err)
		}
	}

	// 绑定新 IP
	if err := proxyRepo.IncrementInUseTx(tx, selected.IP, 1); err != nil {
		return fmt.Errorf("新IP %s 增加使用数失败: %w", selected.IP, err)
	}

	// 更新 Emulator 表
	if err := emulatorRepo.Update(emulator.UUID, map[string]interface{}{"ip": selected.IP}); err != nil {
		return fmt.Errorf("更新模拟器绑定IP失败: %w", err)
	}

	logger.InfofWithTrace(s.Ctx, "模拟器 %s IP 已更新为: %s", emulator.UUID, selected.IP)
	return nil
}

func (s *Svc) switchProxy(emulator *models.Emulator, group *models.Groups) (selected *models.Proxy, err error) {
	// 获取全部代理列表
	err, proxies := s.getProxies(emulator.GroupID)
	if err != nil {
		return nil, err
	}

	// 初始候选列表：负载最小、未超限
	// 筛选出当前使用数 (InUseCount) 小于最大在线限制 group.MaxOnline，且使用数最小的代理列表，作为候选集合。
	// 如果候选集合为空（即所有代理都已满载），
	// 直接从所有代理池中随机选一个，不考虑负载限制，作为备选。
	// 否则，使用候选集合。
	candidates := selectLeastUsedProxies(proxies, int64(group.MaxOnline))
	if len(candidates) == 0 {
		logger.WarnfWithTrace(s.Ctx, "代理池全部已满，UUID: %s，将从所有代理中随机选一个", emulator.UUID)
		selected = pickRandomProxy(proxies, emulator.IP)
	} else {
		selected = pickRandomProxy(candidates, emulator.IP)
	}

	logger.InfofWithTrace(s.Ctx, "模拟器 %s 原IP: %s，初始选中IP: %s", emulator.UUID, emulator.IP, selected.IP)

	// 开始事务（带最多3次尝试更换代理IP）
	const maxRetries = 3

	err = logic.RetryTransaction(s.DB, func(tx *gorm.DB) error {
		proxyRepo := factory.ProxyRepo(tx)

		tried := map[string]bool{} // 已尝试 IP
		for i := 0; i < maxRetries; i++ {
			tried[selected.IP] = true

			// 加锁查询当前选中代理最新使用数（乐观锁机制）。
			selectedLatest, err := proxyRepo.GetByIPForUpdate(selected.IP)
			if err != nil {
				return fmt.Errorf("获取代理最新信息失败: %w", err)
			}

			if selectedLatest.InUseCount+1 <= int64(group.MaxOnline) {
				// 合法，执行切换逻辑
				return s.bindEmulatorToProxyIP(tx, emulator, selected)
			}

			// 当前 IP 已满，尝试重新选择一个未尝试过的 IP
			logger.WarnfWithTrace(s.Ctx, "代理 %s 超载（%d），尝试重新选择", selected.IP, selectedLatest.InUseCount)
			untried := filterUntriedProxies(proxies, tried)
			if len(untried) == 0 {
				logger.WarnfWithTrace(s.Ctx, "无其他可选代理，强制继续使用 %s", selected.IP)
				break // 最后一次容忍
			}
			selected = pickRandomProxy(untried, emulator.IP)
			logger.InfofWithTrace(s.Ctx, "重新选择代理，尝试新IP: %s", selected.IP)
		}
		return s.bindEmulatorToProxyIP(tx, emulator, selected)
	}, 3)

	if err != nil {
		return nil, err
	}

	return selected, nil
}
