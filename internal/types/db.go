/*
@Date: 2021/12/9 11:08
@Author: max.liu
@File : db
*/

package types

// 仅用于 swagger 文档
type BasicQuery struct {
	IDList     []int64           `json:"IdList"`     // id数组
	FuzzyField map[string]string `json:"FuzzyField"` // 模糊查询字段
	Fields     []string          `json:"Fields"`     // 指定返回字段
	Keyword    string            `json:"Keyword"`    // 关键词(全局模糊搜索)
	Order      string            `json:"Order"`      // 排序，支持desc和asc exp:inuse_count desc,created_at asc
	Limit      int               `json:"Limit"`      // 分页条数
	Offset     int               `json:"Offset"`     // 分页偏移量
	Query      string            `json:"Query"`      // 自定义查询语句；使用RSQL语法
}
