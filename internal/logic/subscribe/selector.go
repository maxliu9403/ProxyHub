package subscribe

import (
	"math/rand"
	"time"

	"github.com/maxliu9403/ProxyHub/models"
)

// 选负载最低
func selectLeastUsedProxies(proxies []models.Proxy, maxOnline int64) []models.Proxy {
	var candidates []models.Proxy
	minCount := maxOnline + 1

	for _, p := range proxies {
		if p.InUseCount >= maxOnline {
			continue
		}
		if p.InUseCount < minCount {
			minCount = p.InUseCount
			candidates = []models.Proxy{p}
		} else if p.InUseCount == minCount {
			candidates = append(candidates, p)
		}
	}
	return candidates
}

// 随机选择
func pickRandomProxy(candidates []models.Proxy, currentIP string) *models.Proxy {
	rand.Seed(time.Now().UnixNano())

	// 首先尝试从非当前IP中选
	var altProxies []models.Proxy
	for _, p := range candidates {
		if p.IP != currentIP {
			altProxies = append(altProxies, p)
		}
	}
	if len(altProxies) > 0 {
		return &altProxies[rand.Intn(len(altProxies))]
	}

	// 若都与当前IP重复，只能重用当前
	return &candidates[rand.Intn(len(candidates))]
}

func filterUntriedProxies(all []models.Proxy, tried map[string]bool) []models.Proxy {
	var result []models.Proxy
	for _, p := range all {
		if !tried[p.IP] {
			result = append(result, p)
		}
	}
	return result
}
