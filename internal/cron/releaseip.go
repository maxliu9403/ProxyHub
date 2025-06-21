package cron

//func test()  {
//	job := cronjob.CronJobs
//	job.AddJob()
//
//}
//

// 定时释放IP，如果模拟器所绑定的IP，超过一定时间没有被更新记录，则自动释放
func releaseIp() {
	// 1. 查找emulator表，找出所有没有被删除，同时更新时间到当前时间小于配置时间的模拟器
	// 2. 查找出

}
