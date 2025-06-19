package models

type Groups struct {
	Meta
	Name        string `json:"Name" gorm:"column:name;uniqueIndex;comment:组名"`
	MaxOnline   int    `json:"MaxOnline" gorm:"index;column:max_online;comment:'该分组内的IP最大同时在线模拟器数'"`
	Description string `json:"Description" gorm:"index;column:description;comment:'描述'"`
	Available   int    `json:"Available" gorm:"index;column:available;default:1;comment:'是否是激活状态'"`
}
