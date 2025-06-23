package cron

import (
	"context"

	"github.com/maxliu9403/ProxyHub/internal/config"
	"github.com/maxliu9403/common/cronjob"
)

func RegisterCronJobs(ctx context.Context) {
	job := &ScanExpiredEmulatorJob{
		Svc: NewScanEmulatorsTaskSvc(ctx),
	}

	// 每小时执行一次（示例）："@hourly" or "0 * * * *"
	if _, err := cronjob.CronJobs.AddJob(config.G.CronJob.ReleaseIpPeriod, job); err != nil {
		panic("注册 ScanExpiredEmulatorJob 失败: " + err.Error())
	}
}
