/*
@Date: 2021/1/12 下午2:23
@Author: max.liu
@File : config
@Desc:
*/

package config

import (
	"encoding/json"
	"fmt"

	"github.com/maxliu9403/common/apiserver"
)

type CronJob struct {
	ReleaseIp string `yaml:"release_ip"`
}

type MailCfg struct {
	Enable   bool     `yaml:"enable"`    // 是否启用邮件通知
	SMTPHost string   `yaml:"smtp_host"` // SMTP服务器地址
	SMTPPort int      `yaml:"smtp_port"` // SMTP服务器端口
	Username string   `yaml:"username"`  // SMTP用户名（发件人邮箱）
	Password string   `yaml:"password"`  // SMTP密码或授权码
	To       []string `yaml:"to"`        // 收件人列表
}

type CustomCfg struct {
	IntervalTime int `yaml:"interval_time"`
}

type Config struct {
	apiserver.APIConfig `yaml:"base"`
	CronJob             CronJob   `yaml:"cron_job"`
	CustomCfg           CustomCfg `yaml:"custom_cfg"`
	Mail                MailCfg   `yaml:"mailer"`
}

func (c *Config) String() string {
	configData, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
	}

	return string(configData)
}

var G = &Config{}
