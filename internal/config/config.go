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

type Config struct {
	apiserver.APIConfig `yaml:"base"`
}

func (c *Config) String() string {
	configData, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
	}

	return string(configData)
}

var G = &Config{}
