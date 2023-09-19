/*
@Date: 2021/1/12 下午2:16
@Author: max.liu
@File : main
@Desc:
*/

package main

import (
	"github.com/maxliu9403/common/version"
	_ "github.com/maxliu9403/go-template/docs"
	"github.com/maxliu9403/go-template/internal/app"
)

// @title        go-template
// @version      1.0
// @description  a template
// @BasePath     /

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

var Build string

func main() {
	// 发版时定义此处的版本号
	version.AppVersion.Major = "0"
	version.AppVersion.Minor = "0"
	version.AppVersion.Patch = "1"

	if Build != "" {
		version.AppVersion.Build = Build
	}

	app.Execute()
}
