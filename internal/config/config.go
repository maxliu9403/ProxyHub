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

type CustomCfg struct {
	IntervalTime int `yaml:"interval_time"`
}

type Config struct {
	apiserver.APIConfig `yaml:"base"`
	CronJob             CronJob   `yaml:"cron_job"`
	CustomCfg           CustomCfg `yaml:"custom_cfg"`
}

func (c *Config) String() string {
	configData, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
	}

	return string(configData)
}

var G = &Config{}
