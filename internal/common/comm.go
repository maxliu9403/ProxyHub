/*
@Date: 2021/1/12 下午2:44
@Author: max.liu
@File : comm
@Desc:
*/

package common

type (
	ListData struct {
		Counts int64       `json:"counts"`
		Data   interface{} `json:"data"`
	}
)

type Test struct {
	Enable bool `json:"Enable" swaggerignore:"true"`
}
