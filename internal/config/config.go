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
	ReleaseIpPeriod string `yaml:"release_ip" env:"ReleaseIpPeriod"`
}

type MailCfg struct {
	Enable   bool     `yaml:"enable" env:"MailCfgEnable"`      // 是否启用邮件通知
	SMTPHost string   `yaml:"smtp_host" env:"MailCfgSMTPHost"` // SMTP服务器地址
	SMTPPort int      `yaml:"smtp_port" env:"MailCfgSMTPPort"` // SMTP服务器端口
	Username string   `yaml:"username" env:"MailCfgUsername"`  // SMTP用户名（发件人邮箱）
	Password string   `yaml:"password" env:"MailCfgPassword"`  // SMTP密码或授权码
	SendTo   []string `yaml:"send_to" env:"MailCfgSendTo"`     // 收件人列表
}

type CustomCfg struct {
	IntervalTime int `yaml:"interval_time" env:"IntervalTime"`
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
